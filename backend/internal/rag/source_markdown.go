package rag

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// MarkdownSource loads documents from a directory of .md files.
// Each file becomes one Document.
type MarkdownSource struct {
	Dir string // e.g., "knowledge/coffee/"
}

func (m *MarkdownSource) Name() string {
	return "markdown:" + m.Dir
}

func (m *MarkdownSource) Load(ctx context.Context) ([]Document, error) {
	var docs []Document

	err := filepath.WalkDir(m.Dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		// Skip directories and non-markdown files
		if d.IsDir() || !strings.HasSuffix(d.Name(), ".md") {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("reading %s: %w", path, err)
		}

		// Use filename without extension as title
		title := strings.TrimSuffix(d.Name(), ".md")
		title = strings.ReplaceAll(title, "_", " ")

		docs = append(docs, Document{
			Title:   title,
			Content: string(content),
			Source:  "file://" + path,
		})

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("walking %s: %w", m.Dir, err)
	}

	return docs, nil
}