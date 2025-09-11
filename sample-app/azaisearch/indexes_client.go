package azaisearch

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"

	"sample-app/azaisearch/internal"
	"sample-app/azaisearch/internal/services/search/2025-09-01/searchservice"
)

type IndexesClientOptions struct {
	azcore.ClientOptions
}

// NewIndexesClient creates a new instance of IndexesClient with the specified values.
//   - endpoint - the endpoint of the Azure AI Search service
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - client options, pass nil to accept the default values.
func NewIndexesClient(endpoint string, cred azcore.TokenCredential, options *IndexesClientOptions) (*searchservice.IndexesClient, error) {

	authPolicy := runtime.NewBearerTokenPolicy(cred, []string{internal.TokenScope}, nil)
	return newIndexesClient(endpoint, authPolicy, options)
}

// NewIndexesClientWithSharedKey creates a new instance of IndexesClient with the specified values.
//   - endpoint - the endpoint of the Azure AI Search service
//   - keyCred - used to authorize requests with a shared key
//   - options - client options, pass nil to accept the default values.
func NewIndexesClientWithSharedKey(endpoint string, keyCred *azcore.KeyCredential, options *IndexesClientOptions) (*searchservice.IndexesClient, error) {

	authPolicy := runtime.NewKeyCredentialPolicy(keyCred, "api-key", &runtime.KeyCredentialPolicyOptions{})
	return newIndexesClient(endpoint, authPolicy, options)
}

func newIndexesClient(endpoint string, authPolicy policy.Policy, options *IndexesClientOptions) (*searchservice.IndexesClient, error) {
	if options == nil {
		options = &IndexesClientOptions{}
	}

	c, err := azcore.NewClient(moduleName, moduleVersion, runtime.PipelineOptions{
		PerRetry: []policy.Policy{authPolicy},
	}, &options.ClientOptions)

	if err != nil {
		return nil, err
	}

	return searchservice.NewIndexesClient(endpoint, c)
}
