package practice

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const baseUrl = "https://raw.githubusercontent.com/Thwani47/termilearn-sourcefiles/master"

type fileDownloadedMsg struct{}

//type runTestsMsg struct{ message string }

func getPracticeFiles(concept string) tea.Cmd {
	practiceFileCmd := downloadFile(concept, "main.go")
	testFileCmd := downloadFile(concept, fmt.Sprintf("%s_test.go", concept))

	return tea.Batch(practiceFileCmd, testFileCmd)
}

func downloadFile(folder, fileName string) tea.Cmd {
	return func() tea.Msg {
		dir := fmt.Sprintf("practice/concepts/%s", folder)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return errorMsg{err}
		}

		_, err := os.Stat(fmt.Sprintf("practice/concepts/%s/%s", folder, fileName))

		if os.IsExist(err) {
			return fileDownloadedMsg{}
		}

		out, err := os.Create(fmt.Sprintf("practice/concepts/%s/%s", folder, fileName))
		if err != nil {
			return errorMsg{err}
		}
		defer out.Close()

		resp, err := http.Get(fmt.Sprintf("%s/practice-questions/%s/%s", baseUrl, folder, fileName))
		if err != nil {
			return errorMsg{err}
		}
		defer resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			return errorMsg{fmt.Errorf("Error downloading file: %s", resp.Status)}
		}
		_, err = io.Copy(out, resp.Body)
		if err != nil {
			return errorMsg{err}
		}

		return fileDownloadedMsg{}
	}
}

func parseTestOutput(output string) []testResult {
	var results []testResult
	lines := strings.Split(output, "\n")

	var currentTest *testResult

	for _, line := range lines {
		if strings.HasPrefix(line, "=== RUN") {
			testName := strings.TrimSpace(strings.TrimPrefix(line, "=== RUN"))
			currentTest = &testResult{name: testName, passed: false}
			results = append(results, *currentTest)
		} else if strings.HasPrefix(line, "--- PASS") {
			testName := strings.TrimSpace(strings.TrimPrefix(line, "--- PASS"))
			for i := range results {
				if results[i].name == testName {
					results[i].passed = true
					break
				}
			}
		} else if strings.HasPrefix(line, "--- FAIL") {
			testName := strings.TrimSpace(strings.TrimPrefix(line, "--- FAIL"))
			for i := range results {
				if results[i].name == testName {
					results[i].passed = false
					break
				}
			}
		} else if currentTest != nil && !currentTest.passed {
			// Append error message to the current test result
			currentTest.errorMessage += line + "\n"
		}
	}

	return results
}
