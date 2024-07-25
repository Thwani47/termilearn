package concept

import (
	"embed"
	"fmt"

	"github.com/charmbracelet/glamour"
)

//go:embed notes/*
var notesReader embed.FS

func ReaderConceptNotes(concept string) (string, error) {

	filePath := fmt.Sprintf("notes/%s.md", concept)

	notesContent, err := notesReader.ReadFile(filePath)

	if err != nil {
		return "", fmt.Errorf("failed to read file %s: %w", filePath, err)
	}

	render, err := glamour.NewTermRenderer(glamour.WithStandardStyle("dark"), glamour.WithWordWrap(100))

	if err != nil {
		return "", fmt.Errorf("failed to create render %w", err)
	}
	out, err := render.Render(string(notesContent))

	if err != nil {
		return "", fmt.Errorf("failed to render content: %w", err)
	}

	return out, nil
}
