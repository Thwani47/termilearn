package practice

import (
	"fmt"
	"io"
	"net/http"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const baseUrl = "https://raw.githubusercontent.com/Thwani47/termilearn-sourcefiles/master"

type fileDownloadedMsg struct{}
type runTestsMsg struct{ message string }

func getPracticeFiles(concept string) tea.Cmd {
	practiceFileCmd := downloadFile(concept, "main.go")
	testFileCmd := downloadFile(concept, fmt.Sprintf("%s_test.go", concept))

	return tea.Batch(practiceFileCmd, testFileCmd)
}

func downloadFile(folder, fileName string) tea.Cmd {
	//TODO: do not download file if it already exists. This overwrites the file
	return func() tea.Msg {
		dir := fmt.Sprintf("practice/concepts/%s", folder)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return errorMsg{err}
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
