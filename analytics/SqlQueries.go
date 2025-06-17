package analytics

const (
	GET_CURRENT_STORAGE = `
		SELECT 
			COALESCE(MAX(CASE WHEN spcname = 'pg_default' THEN pg_size_pretty(pg_tablespace_size(spcname)) END), '0 B') AS management_storage,
			COALESCE(MAX(CASE WHEN spcname = 'pg_global' THEN pg_size_pretty(pg_tablespace_size(spcname)) END), '0 B') AS actual_data
		FROM pg_tablespace
		WHERE spcname IN ('pg_default', 'pg_global')
	`
)
