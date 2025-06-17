package analytics

const (
	GET_CURRENT_STORAGE = `
		SELECT 
			COALESCE(MAX(CASE WHEN spcname = 'pg_default' THEN pg_size_pretty(pg_tablespace_size(spcname)) END), '0 B') AS management_storage,
			COALESCE(MAX(CASE WHEN spcname = 'pg_global' THEN pg_size_pretty(pg_tablespace_size(spcname)) END), '0 B') AS actual_data
		FROM pg_tablespace
		WHERE spcname IN ('pg_default', 'pg_global')
	`

	GET_MAX_AVG_TOTAL_EXECUTION_TIME = `
		SELECT
			ROUND(SUM(total_exec_time)::numeric, 2) as total_time_ms,
			ROUND(MAX(max_exec_time)::numeric, 2) AS max_time_ms,
			ROUND(AVG(mean_exec_time)::numeric, 2) AS avg_time_ms
		FROM pg_stat_statements pss
		JOIN pg_database pd ON pss.dbid = pd.oid
		WHERE pd.datname = $1
	`
)
