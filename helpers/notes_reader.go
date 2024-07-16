package helpers

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/glamour"
)

func ReaderConceptNotes(concept string) (string, error) {
	baseDir, err := os.Getwd()

	if err != nil {
		return "", err
	}

	notesPath := filepath.Join(baseDir, "concepts", concept, "notes.md")
	notesContent, err := os.ReadFile(notesPath)

	if err != nil {
		return "", err
	}

	render, _ := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"), glamour.WithWordWrap(100))
	out, err := render.Render(string(notesContent))

	if err != nil {
		return "", err
	}

	return out, nil
}
