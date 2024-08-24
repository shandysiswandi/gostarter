// Package searching provides a unified interface for interacting with various search engines,
// such as ZincSearch, MeiliSearch, Elasticsearch, and Solr. It defines a set of interfaces
// that allow for creating and deleting indexes, indexing documents, retrieving documents,
// and performing search queries. This package is designed to be flexible and extensible,
// making it easy to integrate different search engines with consistent API methods.
package searching

// SearchEngine defines the interface for a search engine that supports indexing,
// document management, and retrieval functionalities by embedding the other interfaces.
type SearchEngine interface {
	Indexer
	DocumentIndexer
	DocumentRetriever
	DocumentManager
}

// Indexer provides methods to create and delete indexes in the search engine.
type Indexer interface {
	// CreateIndex creates a new index in the search engine with the given name and settings.
	// The settings parameter allows customization of the index (e.g., sharding, replication).
	CreateIndex(indexName string, settings any) error

	// DeleteIndex removes an index from the search engine by its name.
	DeleteIndex(indexName string) error
}

// DocumentIndexer provides methods to index single or multiple documents into the search engine.
type DocumentIndexer interface {
	// IndexDocument indexes a single document with the specified document ID.
	// The document parameter can be of any type, depending on the implementation.
	IndexDocument(docID string, document any) error

	// BatchIndexDocuments indexes multiple documents at once.
	// This method is typically more efficient than indexing documents one by one.
	BatchIndexDocuments(documents []any) error
}

// DocumentRetriever provides methods to retrieve and search documents in the search engine.
type DocumentRetriever interface {
	// GetDocument retrieves a single document from the search engine by its ID.
	// It returns the document as an interface and an error if the retrieval fails.
	GetDocument(docID string) (any, error)

	// Search performs a search query with optional filters and returns the results.
	// The query parameter is the search query string, and filters allow for refining the search.
	Search(query string, filters map[string]any) (any, error)
}

// DocumentManager provides methods to delete and update documents in the search engine.
type DocumentManager interface {
	// DeleteDocument removes a document from the search engine by its ID.
	DeleteDocument(docID string) error

	// UpdateDocument updates an existing document in the search engine.
	// The document ID identifies the document to update, and the document parameter contains the new data.
	UpdateDocument(docID string, document any) error
}
