package indexes

type IndexData struct {
	IndexName string   `json:"name"`
	IndexType string   `json:"type"`
	Columns   []string `json:"columns"`
	TableName string   `json:"table_name"`
}

type RetrievedIndex struct {
	IndexName string `json:"index_name"`
	IndexOid  string `json:"index_oid"`
	IndexType string `json:"index_type"`
}

type SpecificIndex struct {
	IndexName string `json:"index_name"`
	IndexType string `json:"index_type"`
}

var DefaultSpecificIndex = SpecificIndex{}
