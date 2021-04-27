package elasticsearch

// IndexedDocument represents an Elasticsearch response of an index request with indexed document information
type IndexedDocument struct {
	Index string `json:"_index"`
	Id    string `json:"_id"`
}
