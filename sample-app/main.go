package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"sample-app/azaisearch"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func ptr[T any](v T) *T { return &v }

func main() {
	endpoint := os.Getenv("AZSEARCH_ENDPOINT")
	apiKey := os.Getenv("AZSEARCH_API_KEY")
	indexName := os.Getenv("AZSEARCH_INDEX_NAME")
	if indexName == "" {
		indexName = "sample-index"
	}

	if endpoint == "" {
		fmt.Println("Please set AZSEARCH_ENDPOINT (and optionally AZSEARCH_API_KEY / AZSEARCH_INDEX_NAME).")
		return
	}

	ctx := context.Background()

	// 1. Create the index (if it does not already exist)
	var (
		indexesClient *azaisearch.IndexesClient
		err           error
	)

	// Decide auth strategy
	useKey := apiKey != ""
	if useKey {
		fmt.Println("Using API Key authentication.")
		indexesClient, err = azaisearch.NewIndexesClientWithSharedKey(endpoint, azcore.NewKeyCredential(apiKey), nil)
	} else {
		fmt.Println("No AZSEARCH_API_KEY provided. Falling back to Azure AD (DefaultAzureCredential).")
		cred, credErr := azidentity.NewDefaultAzureCredential(nil)
		if credErr != nil {
			fmt.Printf("Failed to create DefaultAzureCredential: %v\n", credErr)
			return
		}
		indexesClient, err = azaisearch.NewIndexesClient(endpoint, cred, nil)
	}
	if err != nil {
		panic(err)
	}

	// Define fields
	var (
		fieldKeyName      = "id"
		fieldTitleName    = "title"
		dataTypeEdmString = azaisearch.SearchFieldDataTypeString
		isKey             = true
		searchable        = true
		retrievable       = true
	)

	indexDef := azaisearch.SearchIndex{
		Name: &indexName,
		Fields: []*azaisearch.SearchField{
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
	var docsClient *azaisearch.DocumentsClient
	if useKey {
		docsClient, err = azaisearch.NewDocumentsClientWithSharedKey(endpoint, indexName, azcore.NewKeyCredential(apiKey), nil)
	} else {
		cred, credErr := azidentity.NewDefaultAzureCredential(nil)
		if credErr != nil {
			fmt.Printf("Failed to create DefaultAzureCredential: %v\n", credErr)
			return
		}
		docsClient, err = azaisearch.NewDocumentsClient(endpoint, indexName, cred, nil)
	}
	if err != nil {
		panic(err)
	}

	docKey := "1"
	sampleDoc := map[string]any{
		"id":    docKey,
		"title": "Hello Azure AI Search",
	}
	batch := azaisearch.IndexBatch{Actions: []*azaisearch.IndexAction{{
		ActionType:           ptr(azaisearch.IndexActionTypeUpload),
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
	searchResp, err := docsClient.SearchGet(ctx, &azaisearch.DocumentsClientSearchGetOptions{SearchText: &star}, nil, nil)
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
	searchResp2, err := docsClient.SearchGet(ctx, &azaisearch.DocumentsClientSearchGetOptions{SearchText: &query}, nil, nil)
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
