package indexes

const (
	SELECT_ALL_INDEXES = `SELECT c.relname AS index_name, c.oid AS index_oid, am.amname AS index_type
    FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace JOIN pg_am am ON am.oid = c.relam
    WHERE c.relkind = 'i' AND n.nspname = 'public'`

	SELECT_SPECIFIC_INDEX = `SELECT c.relname AS index_name, am.amname AS index_type
    FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace JOIN pg_am am ON am.oid = c.relam
    WHERE c.relkind = 'i' AND n.nspname = 'public' AND c.oid = $1`
)
