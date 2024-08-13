package practice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const baseUrl = "https://raw.githubusercontent.com/Thwani47/termilearn-sourcefiles/master"

type fileDownloadedMsg struct{}

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

type testEvent struct {
	Action  string `json:"Action"`
	Package string `json:"Package"`
	Test    string `json:"Test"`
	Output  string `json:"Output"`
}

func parseTestOutput(output string) []testResult {
	decoder := json.NewDecoder(strings.NewReader(output))
	resultsMap := make(map[string]testResult)
	results := make([]testResult, 0)

	for {
		var t testEvent

		if err := decoder.Decode(&t); err != nil {
			break
		}

		if t.Action == "run" && t.Test != "" {
			resultsMap[t.Test] = testResult{name: t.Test, passed: false, errorMessage: ""}
		} else if t.Action == "output" && t.Output != "" {
			if entry, ok := resultsMap[t.Test]; ok {
				entry.errorMessage += t.Output
				resultsMap[t.Test] = entry
			}

		} else if t.Action == "pass" && t.Test != "" {
			resultsMap[t.Test] = testResult{name: t.Test, passed: true, errorMessage: ""}
		}
	}

	for _, result := range resultsMap {
		results = append(results, result)
	}
	return results
}
