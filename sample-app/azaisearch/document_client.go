package azaisearch

import (
	"sample-app/azaisearch/internal"
	"sample-app/azaisearch/internal/services/search/2025-09-01/searchindex"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
)

type DocumentClientOptions struct {
	azcore.ClientOptions
}

// NewDocumentsClient creates a new instance of DocumentsClient with the specified values.
//   - endpoint - the endpoint of the Azure AI Search service
//   - indexName - the name of the index to manage documents
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - client options, pass nil to accept the default values.
func NewDocumentsClient(endpoint string, indexName string, cred azcore.TokenCredential, options *DocumentClientOptions) (*searchindex.DocumentsClient, error) {

	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{internal.TokenScope}, nil)

	return newDocumentsClient(endpoint, indexName, authPolicy, options)
}

// NewDocumentsClientWithSharedKey creates a new instance of DocumentsClient with the specified values.
//   - endpoint - the endpoint of the Azure AI Search service
//   - indexName - the name of the index to manage documents
//   - keyCred - used to authorize requests with a shared key
//   - options - client options, pass nil to accept the default values.
func NewDocumentsClientWithSharedKey(endpoint string, indexName string, keyCred *azcore.KeyCredential, options *DocumentClientOptions) (*searchindex.DocumentsClient, error) {

	authPolicy := runtime.NewKeyCredentialPolicy(keyCred, "api-key", &runtime.KeyCredentialPolicyOptions{})

	return newDocumentsClient(endpoint, indexName, authPolicy, options)
}

func newDocumentsClient(endpoint string, indexName string, authPolicy policy.Policy, options *DocumentClientOptions) (*searchindex.DocumentsClient, error) {
	if options == nil {
		options = &DocumentClientOptions{}
	}

	c, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return searchindex.NewDocumentsClient(endpoint, indexName, c)
}
