package tables

const (
	InsertNewTableRecordStmt = `INSERT INTO "Ptable" (oid, name, description, project_id) VALUES ($1, $2, $3, $4) RETURNING id;`
	CheckOwnershipStmt       = `SELECT COUNT(*) FROM "projects" WHERE oid = $1 AND owner_id = $2;`
	GetTableNameStmt         = `SELECT name FROM "Ptable" WHERE oid = $1;`
	DropTableStmt            = `DROP TABLE IF EXISTS %s;`
	DeleteTableStmt          = `DELETE FROM "Ptable" WHERE %s = $1;`

	ReadTableStmt = `SELECT 
						c.column_name AS column_name,
						c.data_type AS data_type,
						(c.is_nullable = 'YES') AS is_nullable,
						c.column_default AS column_default,
						tc.constraint_name AS unique_constraint_name,
						tc.constraint_type AS unique_constraint_type,
						fk.ref_table AS referenced_table,
						fk.ref_column AS referenced_column
					FROM 
						information_schema.columns c
					LEFT JOIN 
						information_schema.constraint_column_usage ccu
						ON c.table_name = ccu.table_name 
						AND c.column_name = ccu.column_name
						AND c.table_schema = ccu.table_schema
					LEFT JOIN 
						information_schema.table_constraints tc
						ON ccu.constraint_name = tc.constraint_name
						AND tc.constraint_type IN ('UNIQUE', 'PRIMARY KEY')
					LEFT JOIN (
						SELECT 
							kc.table_name,
							kc.column_name,
							rc.constraint_name,
							ccu2.table_name AS ref_table,
							ccu2.column_name AS ref_column
						FROM 
							information_schema.key_column_usage kc
						JOIN 
							information_schema.referential_constraints rc
							ON kc.constraint_name = rc.constraint_name
						JOIN 
							information_schema.constraint_column_usage ccu2
							ON rc.unique_constraint_name = ccu2.constraint_name
					) fk
						ON c.table_name = fk.table_name 
						AND c.column_name = fk.column_name
					WHERE 
						c.table_name = $1
						AND c.table_schema = 'public'
					ORDER BY 
						c.column_name;`
	InsertNewRowStmt = `
		INSERT INTO "%s"(%s) VALUES(%s)
	`
)
