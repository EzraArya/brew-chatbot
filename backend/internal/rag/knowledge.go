package rag

import "context"

// Document is a single unit of knowledge.
// It's the currency flowing between KnowledgeSources and Retrievers.
type Document struct {
	Title   string // e.g., "V60 Pour Over Guide"
	Content string // full text content
	Source  string // origin: "file://knowledge/coffee/v60.md" or "https://..."
}

// KnowledgeSource answers: WHERE does knowledge come from?
// Implement this for: markdown files, web scrapers, REST APIs.
type KnowledgeSource interface {
	Name() string
	Load(ctx context.Context) ([]Document, error)
}

// Retriever answers: WHICH documents are relevant for this query?
// Implement this for:
//   - FullContextRetriever: ignores query, returns everything (Phase 1)
//   - VectorRetriever:      embeds query, similarity search (Phase 2)
type Retriever interface {
	Retrieve(ctx context.Context, query string) ([]Document, error)
}