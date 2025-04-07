package tables

const (
	InsertNewTableRecordStmt = `INSERT INTO "Ptable" (oid, name, description, project_id) VALUES ($1, $2, $3, $4) RETURNING id;`
	DeleteTableRecordStmt = `DELETE FROM "Ptable" WHERE id = $1;`
)