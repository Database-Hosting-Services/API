package tables

// Table struct is a row record of the tables table in the database
type Table struct {
	ID          int    `json:"id" db:"id"`
	ProjectID   int64    `json:"project_id" db:"project_id"`
	OID		  	string `json:"oid" db:"oid"`
	Name 	  	string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
}

// TableData struct is the how the client will send the definition of the table to the server
// and how the server will respond with the table definition
type ClientTable struct {
	TableName string `json:"tableName"`
	Columns   []Column `json:"columns"`
}

type Column struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	IsUnique bool   `json:"isUnique"`
	IsNullable bool `json:"isNullable"`
	IsPrimaryKey bool `json:"isPrimaryKey"`
	ForeignKey ForeignKey `json:"foreignKey"`
}

type ForeignKey struct {
	ColumnName string `json:"columnName"`
	TableName string `json:"tableName"`
}