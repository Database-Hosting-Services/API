package tables

import (
	"DBHS/utils"
)

// Table struct is a row record of the tables table in the database
type Table struct {
	ID          int64       `json:"id" db:"id"`
	ProjectID   int64       `json:"project_id" db:"project_id"`
	OID         string      `json:"oid" db:"oid"`
	Name        string      `json:"name" db:"name" validate:"required"`
	Description string      `json:"description" db:"description"`
	Schema      *utils.Table `json:"schema" validate:"required"`
}

type UpdateTableSchema struct {
	Table
	Renames []utils.RenameRelation `json:"renames"`
}



type ShortTable struct {
	OID  string `json:"oid" db:"oid"`
	Name string `json:"name" db:"name"`
}

// TableData struct is the how the client will send the definition of the table to the server
// and how the server will respond with the table definition
type ClientTable struct {
	TableName string   `json:"tableName"`
	Columns   []Column `json:"columns"`
}

// TableUpdate struct is the how the client will send the alteration of the table to the server
type TableUpdate struct {
	Inserts ColumnCollection `json:"insert"`
	Updates []UpdateColumn   `json:"update"`
	Deletes []string         `json:"delete"`
}

type UpdateColumn struct {
	Name   string `json:"name"`
	Update Column `json:"update"`
}

type ColumnCollection struct {
	Columns []Column `json:"columns"`
}

type Column struct {
	Name         string     `json:"name"`
	Type         string     `json:"type"`
	IsUnique     *bool      `json:"isUnique"`
	IsNullable   *bool      `json:"isNullable"`
	IsPrimaryKey *bool      `json:"isPrimaryKey"`
	ForeignKey   ForeignKey `json:"foreignKey"`
}

type ForeignKey struct {
	ColumnName string `json:"columnName"`
	TableName  string `json:"tableName"`
}

/*
  c.column_name AS column_name,
    c.data_type AS data_type,
    c.is_nullable AS is_nullable,
    c.column_default AS column_default,
    tc.constraint_name AS unique_constraint_name,
    tc.constraint_type AS unique_constraint_type,
    fk.constraint_name AS foreign_key_name,
    fk.ref_table AS referenced_table,
    fk.ref_column AS referenced_column
*/

type DbColumn struct {
	Name                 string  `json:"name" db:"column_name"`
	Type                 string  `json:"type" db:"data_type"`
	IsNullable           bool    `json:"is_nullable" db:"is_nullable"`
	ColumnDefault        *string `json:"column_default" db:"column_default"`
	UniqueConstraintName *string `json:"unique_constraint_name" db:"unique_constraint_name"`
	UniqueConstraintType *string `json:"unique_constraint_type" db:"unique_constraint_type"`
	ReferencedTable      *string `json:"referenced_table" db:"referenced_table"`
	ReferencedColumn     *string `json:"referenced_column" db:"referenced_column"`
}

type ShowColumn struct {
	Name string `json:"name"`
	Type string `json:"type"`
}
type Data struct {
	Columns []ShowColumn             `json:"columns"`
	Rows    []map[string]interface{} `json:"rows"`
}

type RowValue struct {
	ColumnName string `json: "column"`
	Value interface{} `json: "value"`
}
