package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"sample-app/searchservice"
	"sample-app/searchindex"
	"sample-app/searchclients"
)

func ptr[T any](v T) *T { return &v }

func main() {
	endpoint := os.Getenv("AZSEARCH_ENDPOINT")
	apiKey := os.Getenv("AZSEARCH_API_KEY")
	indexName := os.Getenv("AZSEARCH_INDEX_NAME")
	if indexName == "" {
		indexName = "sample-index"
	}

	if endpoint == "" || apiKey == "" {
		fmt.Println("Please set AZSEARCH_ENDPOINT, AZSEARCH_API_KEY (and optional AZSEARCH_INDEX_NAME).")
		return
	}

	ctx := context.Background()

	// 1. Create the index (if it does not already exist)
	indexesClient, err := searchclients.NewIndexesClient(endpoint, apiKey)
	if err != nil {
		panic(err)
	}

	// Define fields
	var (
		fieldKeyName      = "id"
		fieldTitleName    = "title"
		dataTypeEdmString = searchservice.SearchFieldDataTypeString
		isKey             = true
		searchable        = true
		retrievable       = true
	)

	indexDef := searchservice.SearchIndex{
		Name: &indexName,
		Fields: []*searchservice.SearchField{
			{Name: &fieldKeyName, Type: &dataTypeEdmString, Key: &isKey, Filterable: ptr(true), Sortable: ptr(true), Retrievable: &retrievable},
			{Name: &fieldTitleName, Type: &dataTypeEdmString, Searchable: &searchable, Retrievable: &retrievable},
		},
	}

	// Try create; if already exists, skip
	if _, err = indexesClient.Create(ctx, indexDef, nil, nil); err != nil {
		// attempt get to see if exists
		if _, getErr := indexesClient.Get(ctx, indexName, nil, nil); getErr != nil {
			fmt.Printf("Failed to create index and it does not exist: %v\n", err)
			return
		} else {
			fmt.Printf("Index '%s' already exists, continuing.\n", indexName)
		}
	} else {
		fmt.Printf("Created index '%s'.\n", indexName)
	}

	// 2. Index a sample document
	docsClient, err := searchclients.NewDocumentsClient(endpoint, indexName, apiKey)
	if err != nil {
		panic(err)
	}

	docKey := "1"
	sampleDoc := map[string]any{
		"id":    docKey,
		"title": "Hello Azure AI Search",
	}
	batch := searchindex.IndexBatch{Actions: []*searchindex.IndexAction{{
		ActionType:           ptr(searchindex.IndexActionTypeUpload),
		AdditionalProperties: sampleDoc,
	}}}

	if _, err = docsClient.Index(ctx, batch, nil, nil); err != nil {
		fmt.Printf("Indexing failed: %v\n", err)
		return
	}
	fmt.Println("Submitted indexing batch. Waiting for propagation...")
	time.Sleep(3 * time.Second)

	// 3. Simple wildcard search to verify indexing
	star := "*"
	searchResp, err := docsClient.SearchGet(ctx, &searchindex.DocumentsClientSearchGetOptions{SearchText: &star}, nil, nil)
	if err != nil {
		fmt.Printf("Wildcard search failed: %v\n", err)
		return
	}
	fmt.Printf("Wildcard search returned %d result(s).\n", len(searchResp.SearchDocumentsResult.Results))
	for _, r := range searchResp.SearchDocumentsResult.Results {
		if r != nil && r.AdditionalProperties != nil {
			fmt.Printf("  Doc: %v\n", r.AdditionalProperties)
		}
	}

	// 4. Targeted search (user-provided or default term)
	query := os.Getenv("AZSEARCH_QUERY")
	if query == "" {
		query = "hello" // default demo term
	}
	searchResp2, err := docsClient.SearchGet(ctx, &searchindex.DocumentsClientSearchGetOptions{SearchText: &query}, nil, nil)
	if err != nil {
		fmt.Printf("Query search failed: %v\n", err)
		return
	}
	fmt.Printf("Query search ('%s') returned %d result(s).\n", query, len(searchResp2.SearchDocumentsResult.Results))
	for _, r := range searchResp2.SearchDocumentsResult.Results {
		if r != nil && r.AdditionalProperties != nil {
			fmt.Printf("  Match: %v\n", r.AdditionalProperties)
		}
	}
}
