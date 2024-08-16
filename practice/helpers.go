package practice

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

const baseUrl = "https://raw.githubusercontent.com/Thwani47/termilearn-sourcefiles/master"

type Question struct {
	Title        string `json:"title"`
	QuestionType string `json:"questionType"`
}

type MCQQuestion struct {
	Question
	QuestionText string   `json:"question"`
	Answers      []string `json:"answers"`
	Answer       string   `json:"answer"`
}

type EditQuestion struct {
	Question
	File     string `json:"file"`
	TestFile string `json:"testFile"`
}

type QuestionWrapper struct {
	QuestionType string
	MCQQuestion  *MCQQuestion
	EditQuestion *EditQuestion
}

func (q *QuestionWrapper) UnmarshalJSON(data []byte) error {
	var base Question

	if err := json.Unmarshal(data, &base); err != nil {
		return err
	}

	q.QuestionType = base.QuestionType
	switch base.QuestionType {
	case "mcq":
		var mcq MCQQuestion
		if err := json.Unmarshal(data, &mcq); err != nil {
			return err
		}
		q.MCQQuestion = &mcq
	case "edit":
		var edit EditQuestion
		if err := json.Unmarshal(data, &edit); err != nil {
			return err
		}
		q.EditQuestion = &edit
	default:
		return fmt.Errorf("unknown question type: %s", base.QuestionType)
	}

	return nil
}

type fileDownloadedMsg struct {
	questions []QuestionWrapper
	err       error
}

func getPracticeFiles(concept string) tea.Cmd {
	practiceFileCmd := downloadFile(concept, "main.go")
	testFileCmd := downloadFile(concept, fmt.Sprintf("%s_test.go", concept))

	return tea.Batch(practiceFileCmd, testFileCmd)
}

// TODO: I need to find a way to download an entire folder and not specify the file
func downloadFile(folder, fileName string) tea.Cmd {
	return func() tea.Msg {
		dir := fmt.Sprintf("practice/concepts/%s", folder)
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return errorMsg{err}
		}

		_, err := os.Stat(fmt.Sprintf("practice/concepts/%s/%s", folder, fileName))

		if os.IsExist(err) {
			questions, err := readQuestions(fmt.Sprintf("practice/concepts/%s", folder))
			log.Printf("%v %v\n", questions, err)
			return fileDownloadedMsg{
				err:       err,
				questions: questions,
			}
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

		questions, err := readQuestions(fmt.Sprintf("practice/concepts/%s", folder))
		log.Printf("%v %v\n", questions, err)
		return fileDownloadedMsg{
			err:       err,
			questions: questions,
		}

	}
}

func readQuestions(folder string) ([]QuestionWrapper, error) {
	content, err := os.ReadFile(fmt.Sprintf("%s/questions.json", folder))

	if err != nil {
		return nil, fmt.Errorf("failed to open question files: %v", err)
	}

	var questions []QuestionWrapper

	if err = json.Unmarshal(content, &questions); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %v", err)
	}

	log.Println("questions ", questions)

	return questions, nil

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
