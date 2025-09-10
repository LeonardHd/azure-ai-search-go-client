package main

// NOTE: These constructors use unsafe reflection to populate unexported fields
// of the generated clients so we can keep the generated folders (aisearch,
// aisearchindex) untouched. This is brittle: regenerating code that changes
// internal field names will break these functions.

// TODO: probably want to implement something closer to the official way of
// doing this, see:
// https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/batch/azbatch/custom_client.go
// https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/containers/azcontainerregistry/custom_client.go
// https://github.com/Azure/azure-sdk-for-go/blob/main/sdk/storage/azqueue/queue_client.go

import (
	"net/http"
	"reflect"
	"strings"
	"unsafe"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"

	"sample-app/services/search/2025-09-01/searchindex"
	"sample-app/services/search/2025-09-01/searchservice"
)

// apiKeyPolicy adds the api-key header to every outbound request.
type apiKeyPolicy struct{ key string }

func (a *apiKeyPolicy) Do(req *policy.Request) (*http.Response, error) {
	req.Raw().Header.Set("api-key", a.key)
	return req.Next()
}

func newPerCallPolicy(key string) policy.Policy { return &apiKeyPolicy{key: key} }

func NewIndexesClient(endpoint, apiKey string) (*searchservice.IndexesClient, error) {
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = "https://" + endpoint
	}
	internal, err := azcore.NewClient("azure-search-indexes", "v0.1.0", runtime.PipelineOptions{}, &policy.ClientOptions{PerCallPolicies: []policy.Policy{newPerCallPolicy(apiKey)}})
	if err != nil {
		return nil, err
	}
	var client searchservice.IndexesClient
	// Unsafe set of unexported fields
	rv := reflect.ValueOf(&client).Elem()
	setField := func(name string, value any) {
		f := rv.FieldByName(name)
		ptr := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
		ptr.Set(reflect.ValueOf(value))
	}
	setField("internal", internal)
	setField("endpoint", endpoint)
	return &client, nil
}

func NewDocumentsClient(endpoint, indexName, apiKey string) (*searchindex.DocumentsClient, error) {
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = "https://" + endpoint
	}
	internal, err := azcore.NewClient("azure-search-documents", "v0.1.0", runtime.PipelineOptions{}, &policy.ClientOptions{PerCallPolicies: []policy.Policy{newPerCallPolicy(apiKey)}})
	if err != nil {
		return nil, err
	}
	var client searchindex.DocumentsClient
	rv := reflect.ValueOf(&client).Elem()
	setField := func(name string, value any) {
		f := rv.FieldByName(name)
		ptr := reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem()
		ptr.Set(reflect.ValueOf(value))
	}
	setField("internal", internal)
	setField("endpoint", endpoint)
	setField("indexName", indexName)
	return &client, nil
}
