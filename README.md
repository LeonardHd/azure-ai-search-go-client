# Azure AI Search Go Client via Autorest

The azure sdk for Go currently does not include a client for Azure AI Search.

As a temporary solution (until a first-party client is available), you can use the `autorest` tool to generate a Go client from the Azure Search REST API specifications.


## Prerequisites
- Go 1.20 or later
- Autorest installed `npm install -g autorest`

## Generated the clients and run sample app (without adjustments for `azure-rest-api-specs` repo)

```bash

# Clone the Azure REST API specifications repository (include API specs for autorest)
git clone https://github.com/Azure/azure-rest-api-specs.git 

# Generate the Azure Search Index client
autorest --input-file=azure-rest-api-specs/specification/search/data-plane/Azure.Search/stable/2025-09-01/searchindex.json --go --containing-module --output-folder=sample-app/services/search/2025-09-01/searchindex --clear-output-folder

# Generate the Azure Search client
autorest --input-file=azure-rest-api-specs/specification/search/data-plane/Azure.Search/stable/2025-09-01/searchservice.json --go --containing-module --output-folder=sample-app/services/search/2025-09-01/searchservice --clear-output-folder

# NOTE: These generated clients will NOT necessarily follow the official clients as
# this does not apply any customizations that the official clients might have.

# Create a .env file to store your Azure Search service details
echo "AZSEARCH_ENDPOINT='https://<your-service>.search.windows.net'" > .env
echo "AZSEARCH_API_KEY='<your-admin-or-query-key>'" >> .env
echo "AZSEARCH_INDEX_NAME='sample-index'" >> .env


cd sample-app
go mod tidy
go build ./...

# Run the sample app
set -a; source ../.env; go run .
```

## Adjustments to `readme.go.md` in `azure-rest-api-specs`

To ensure that the generated clients are more aligned with the official Azure SDK implementations,
you can adjust the `readme.go.md` file in the `azure-rest-api-specs` repository to include
the latest API version and any necessary customizations (autorest will pick up the custom directives).

```bash
cp readme.go.md azure-rest-api-specs/specification/search/data-plane/Azure.Search/readme.go.md

# Adjust the `readme.go.md` to include 2025-09-01 version for both clients (searchindex and searchservice)
autorest azure-rest-api-specs/specification/search/data-plane/Azure.Search --containing-module --tag=package-2025-09-searchindex --go --go-sdk-folder=$(pwd)/go-sdk-folder
autorest azure-rest-api-specs/specification/search/data-plane/Azure.Search --containing-module --tag=package-2025-09-searchservice --go --go-sdk-folder=$(pwd)/go-sdk-folder
```