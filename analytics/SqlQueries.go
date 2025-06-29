package analytics

const (
	GET_CURRENT_STORAGE = `
		SELECT 
			COALESCE(MAX(CASE WHEN spcname = 'pg_default' THEN pg_size_pretty(pg_tablespace_size(spcname)) END), '0 B') AS management_storage,
			COALESCE(MAX(CASE WHEN spcname = 'pg_global' THEN pg_size_pretty(pg_tablespace_size(spcname)) END), '0 B') AS actual_data
		FROM pg_tablespace
		WHERE spcname IN ('pg_default', 'pg_global')
	`

	GET_TOTAL_TimeAndQueries = `
		SELECT
			ROUND(SUM(total_exec_time)::numeric, 2) as total_time_ms,
			SUM(calls) as total_queries
		FROM pg_stat_statements pss
		JOIN pg_database pd ON pss.dbid = pd.oid
		WHERE pd.datname = $1
	`

	GET_READ_WRITE_CPU = `
		SELECT                                   
			SUM(CASE WHEN query ILIKE 'SELECT%' THEN calls ELSE 0 END) as read_queries,
			SUM(CASE WHEN query ILIKE ANY(ARRAY['INSERT%', 'UPDATE%', 'DELETE%']) THEN calls ELSE 0 END) as write_queries,
			ROUND(SUM(total_exec_time)::numeric, 2) as total_cpu_time_ms
		FROM pg_stat_statements pss
		JOIN pg_database pd ON pss.dbid = pd.oid
		WHERE pd.datname = current_database()
		GROUP BY pd.datname;
	`

	GET_ALL_CURRENT_STORAGE = `SELECT created_at::text, data->>'Management storage', data->>'Actual data' FROM analytics WHERE type = 'Storage' and "projectId" = $1 ORDER BY created_at DESC;`

	GET_ALL_EXECUTION_TIME_STATS = `SELECT created_at::text, data->>'total_time_ms', data->>'total_queries' FROM analytics WHERE type = 'ExecutionTimeStats' AND "projectId" = $1;	`

	GET_ALL_DATABASE_USAGE_STATS = `SELECT created_at::text, data->>'read_write_cost', data->>'cpu_cost', data->>'total_cost' FROM analytics WHERE type = 'DatabaseUsageStats' AND "projectId" = $1;`
)
