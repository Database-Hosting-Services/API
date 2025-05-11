package schemas

const (
	GetDatabaseByName   = `SELECT host, port, user_id, password, db_name, ssl_mode, created_at FROM database_config WHERE db_name = $1`
	GetTableNameByOID   = `SELECT name FROM "Ptable" WHERE oid = $1`
	GetTableSchemaQuery = `SELECT
					t.table_name,
					c.column_name,
					c.data_type,
					c.is_nullable = 'YES' AS is_nullable,
					EXISTS (
						SELECT 1 
						FROM information_schema.key_column_usage kcu
						JOIN information_schema.table_constraints tc
						ON kcu.constraint_name = tc.constraint_name
						WHERE tc.constraint_type = 'PRIMARY KEY'
						AND kcu.table_name = t.table_name
						AND kcu.column_name = c.column_name
						AND kcu.table_schema = 'public'
					) AS is_primary_key,
					EXISTS (
						SELECT 1 
						FROM information_schema.key_column_usage kcu
						JOIN information_schema.table_constraints tc
						ON kcu.constraint_name = tc.constraint_name
						WHERE tc.constraint_type = 'UNIQUE'
						AND kcu.table_name = t.table_name
						AND kcu.column_name = c.column_name
						AND kcu.table_schema = 'public'
					) AS is_unique,
					fk.foreign_table_name,
					fk.foreign_column_name
				FROM 
					information_schema.tables t
				JOIN 
					information_schema.columns c
					ON t.table_name = c.table_name
					AND t.table_schema = c.table_schema
				LEFT JOIN LATERAL (
					SELECT
						ccu.table_name AS foreign_table_name,
						ccu.column_name AS foreign_column_name
					FROM 
						information_schema.table_constraints tc
					JOIN 
						information_schema.key_column_usage kcu
						ON tc.constraint_name = kcu.constraint_name
					JOIN 
						information_schema.constraint_column_usage ccu
						ON ccu.constraint_name = tc.constraint_name
					WHERE 
						tc.constraint_type = 'FOREIGN KEY'
						AND tc.table_name = t.table_name
						AND kcu.column_name = c.column_name
						AND tc.table_schema = 'public'
					LIMIT 1
				) fk ON true
				WHERE 
					t.table_schema = 'public'
					AND t.table_type = 'BASE TABLE'
					AND t.table_name = $1
				ORDER BY 
					c.ordinal_position;`

	GetAllTablesSchema = `SELECT
					t.table_name,
					c.column_name,
					c.data_type,
					c.is_nullable = 'YES' AS is_nullable,
					EXISTS (
						SELECT 1 
						FROM information_schema.key_column_usage kcu
						JOIN information_schema.table_constraints tc
						ON kcu.constraint_name = tc.constraint_name
						WHERE tc.constraint_type = 'PRIMARY KEY'
						AND kcu.table_name = t.table_name
						AND kcu.column_name = c.column_name
						AND kcu.table_schema = 'public'
					) AS is_primary_key,
					EXISTS (
						SELECT 1 
						FROM information_schema.key_column_usage kcu
						JOIN information_schema.table_constraints tc
						ON kcu.constraint_name = tc.constraint_name
						WHERE tc.constraint_type = 'UNIQUE'
						AND kcu.table_name = t.table_name
						AND kcu.column_name = c.column_name
						AND kcu.table_schema = 'public'
					) AS is_unique,
					fk.foreign_table_name,
					fk.foreign_column_name
				FROM 
					information_schema.tables t
				JOIN 
					information_schema.columns c
					ON t.table_name = c.table_name
					AND t.table_schema = c.table_schema
				LEFT JOIN LATERAL (
					SELECT
						ccu.table_name AS foreign_table_name,
						ccu.column_name AS foreign_column_name
					FROM 
						information_schema.table_constraints tc
					JOIN 
						information_schema.key_column_usage kcu
						ON tc.constraint_name = kcu.constraint_name
					JOIN 
						information_schema.constraint_column_usage ccu
						ON ccu.constraint_name = tc.constraint_name
					WHERE 
						tc.constraint_type = 'FOREIGN KEY'
						AND tc.table_name = t.table_name
						AND kcu.column_name = c.column_name
						AND tc.table_schema = 'public'
					LIMIT 1
				) fk ON true
				WHERE 
					t.table_schema = 'public'
					AND t.table_type = 'BASE TABLE'
				ORDER BY 
					t.table_name, 
					c.ordinal_position;`
)
