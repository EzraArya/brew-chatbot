package rag

import (
    "context"
    "fmt"
)

// FullContextRetriever returns ALL documents from all sources.
// The query is intentionally ignored — Gemini sees the entire knowledge
// corpus via context caching and decides what's relevant itself.
// Swap this for VectorRetriever in Phase 2 when the corpus grows large.
type FullContextRetriever struct {
    Sources []KnowledgeSource
}

func (r *FullContextRetriever) Retrieve(ctx context.Context, query string) ([]Document, error) {
    var all []Document

    for _, source := range r.Sources {
        docs, err := source.Load(ctx)
        if err != nil {
            return nil, fmt.Errorf("loading from %s: %w", source.Name(), err)
        }
        all = append(all, docs...)
    }

    return all, nil
}