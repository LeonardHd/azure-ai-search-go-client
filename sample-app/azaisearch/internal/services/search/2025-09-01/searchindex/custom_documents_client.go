package searchindex

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
)

func NewDocumentsClient(endpoint string, indexName string, coreclient *azcore.Client) (*DocumentsClient, error) {
	return &DocumentsClient{
		internal:  coreclient,
		endpoint:  endpoint,
		indexName: indexName,
	}, nil
}
