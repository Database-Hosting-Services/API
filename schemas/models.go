package schemas

type Column struct {
	Name         string       `json:"name"`
	Type         string       `json:"type"`
	IsUnique     bool         `json:"isUnique"`
	IsNullable   bool         `json:"isNullable"`
	IsPrimaryKey bool         `json:"isPrimaryKey"`
	ForeignKeys  []ForeignKey `json:"foreignKey"`
}

type ForeignKey struct {
	TableName  string `json:"tableName"`
	ColumnName string `json:"columnName"`
}

type TableSchema struct {
	TableName string   `json:"tableName"`
	Cols      []Column `json:"columnss"`
}

type SchemaResponse struct {
	Schema []TableSchema `json:"schema"`
}
