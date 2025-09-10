# Azure AI Search Go Client via Autorest

Microsoft's offical Go SDK does not include clients for Azure AI Search services.

This repository demonstrates how one could use the `autorest` tool to generate Go clients from Microsoft's Azure API
specifications themselves.

## Fundamentals

* [AutoRest](https://github.com/Azure/autorest) is a tool that generates client libraries for RESTful web services 
  based on OpenAPI (formerly Swagger) specifications. It supports multiple programming languages via generators (e.g., Go, Python, C#, Java).
* [Azure REST API Specifications](https://github.com/Azure/azure-rest-api-specs) is a repository that contains the OpenAPI
  specifications for various Azure services. These include specifications for the data plane and management plane of Azure services. Each of the specifications follows a prescribed structure and format (see [directory structure](https://github.com/Azure/azure-rest-api-specs/blob/main/documentation/directory-structure.md#key-concepts)).

  **Note:** In this sample we focus on the data plane specifications for Azure AI Search, found under `specification/search/data-plane/Azure.Search` folder.
  The service folder includes the important files for generating clients, including `readme.md` and the OpenAPI specification files (e.g., `searchindex.json`, `searchservice.json`).

  Generally, API specifications used to be OpenAPI (Swagger) files, but are now definied as TypeSpec files, used to emit OpenAPI specifications.

* [Azure SDK for Go](https://github.com/Azure/azure-sdk-for-go) is the official Go SDK for Azure services.
  It provides a set of libraries that allow developers to interact with Azure services using Go.
  However, not all Azure services have official SDK support in Go (e.g., Azure AI Search).
* [Azure SDL for Go Guidelines](https://azure.github.io/azure-sdk/golang_introduction.html) provide guidelines and best
  practices for developing Azure SDKs in Go, ensuring consistency and quality across the SDKs.
  By using `autorest` to generate clients, we can ensure that the generated clients adhere to these guidelines as closely as possible.

Hence, we can use `autorest` to generate clients for Azure services based on the OpenAPI specifications found in the `azure-rest-api-specs` repository.

> IMPORTANT: Both `autorest` and the `azure-rest-api-specs` repository are maintained by Microsoft, but using them to generate clients for Azure services is your own responsibility.
> For official support and updates, it's recommended to use the first-party SDKs provided by Microsoft.
> However, for services that have no official SDK support (like Azure AI Search in Go), this approach can be a viable alternative to avoid writing custom clients from scratch.

## Expected Outcome

By following the steps in this repository, you will be able to:
* Generate Go clients for Azure AI Search services using `autorest` and the OpenAPI specifications
    from the `azure-rest-api-specs` repository.
* Have a generated client based on the `azcore` module (https://github.com/Azure/azure-sdk-for-go/tree/main/sdk/azcore).
  As a result, the generated clients will be able to integrate with other official Azure SDK for Go components that also use `azcore` (e.g., for authentication via `azidentity`).

## Prerequisites
- Go 1.20 or later
- Autorest installed `npm install -g autorest`

## Generated the clients and run sample app (without adjustments for `azure-rest-api-specs` repo)

These files use the bare OpenAPI specifications without any autorest configuration and customizations from the `azure-rest-api-specs` repo.

As a result, the generated clients will be less aligned with the official Azure SDK for Go guidelines and implementations.

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

## Using `azure-rest-api-specs` adjustments for `readme.go.md` for Azure Search clients

To ensure that the generated clients are more aligned with the official Azure SDK implementations,
you can adjust the `readme.go.md` file in the `azure-rest-api-specs` repository to include
the latest API version and any necessary customizations (autorest will pick up the custom directives).

The commands below will use the existing AutoRest configuration from the `azure-rest-api-specs` repo when
generating the clients.

```bash
cp readme.go.md azure-rest-api-specs/specification/search/data-plane/Azure.Search/readme.go.md

# Adjust the `readme.go.md` to include 2025-09-01 version for both clients (searchindex and searchservice)
autorest azure-rest-api-specs/specification/search/data-plane/Azure.Search --containing-module --tag=package-2025-09-searchindex --go --go-sdk-folder=$(pwd)/sample-app
autorest azure-rest-api-specs/specification/search/data-plane/Azure.Search --containing-module --tag=package-2025-09-searchservice --go --go-sdk-folder=$(pwd)/sample-app
```