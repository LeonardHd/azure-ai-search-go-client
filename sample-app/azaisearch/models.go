package azaisearch

import (
	"sample-app/azaisearch/internal/services/search/2025-09-01/searchindex"
	"sample-app/azaisearch/internal/services/search/2025-09-01/searchservice"
)

type SearchIndex = searchservice.SearchIndex
type SearchField = searchservice.SearchField
type IndexBatch = searchindex.IndexBatch
type IndexAction = searchindex.IndexAction
type DocumentsClientSearchGetOptions = searchindex.DocumentsClientSearchGetOptions

const SearchFieldDataTypeString = searchservice.SearchFieldDataTypeString
const IndexActionTypeUpload = searchindex.IndexActionTypeUpload
