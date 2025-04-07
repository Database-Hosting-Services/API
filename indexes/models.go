package indexes

type IndexData struct {
	IndexName string   `json:"name"`
	IndexType string   `json:"type"`
	Columns   []string `json:"columns"`
	TableName string   `json:"table_name"`
}
