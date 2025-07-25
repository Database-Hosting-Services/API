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
			COALESCE(ROUND(SUM(total_exec_time)::numeric, 2), 0) as total_time_ms,
			COALESCE(SUM(calls), 0) as total_queries
		FROM pg_stat_statements pss
		JOIN pg_database pd ON pss.dbid = pd.oid
		WHERE pd.datname = $1
	`

	GET_READ_WRITE_CPU = `
		SELECT                                   
			COALESCE(SUM(CASE WHEN query ILIKE 'SELECT%' THEN calls ELSE 0 END), 0) as read_queries,
			COALESCE(SUM(CASE WHEN query ILIKE ANY(ARRAY['INSERT%', 'UPDATE%', 'DELETE%']) THEN calls ELSE 0 END), 0) as write_queries,
			COALESCE(ROUND(SUM(total_exec_time)::numeric, 2), 0) as total_cpu_time_ms
		FROM pg_stat_statements pss
		JOIN pg_database pd ON pss.dbid = pd.oid
		WHERE pd.datname = current_database()
		GROUP BY pd.datname;
	`

	GET_ALL_CURRENT_STORAGE = `SELECT created_at::text, data->>'Management storage', data->>'Actual data' FROM analytics WHERE type = 'Storage' and "projectId" = $1 ORDER BY created_at DESC;`

	GET_ALL_EXECUTION_TIME_STATS = `SELECT created_at::text, COALESCE((data->>'total_time_ms')::numeric, 0), COALESCE((data->>'total_queries')::bigint, 0) FROM analytics WHERE type = 'ExecutionTimeStats' AND "projectId" = $1;	`

	GET_ALL_DATABASE_USAGE_STATS = `SELECT created_at::text, COALESCE((data->>'read_write_cost')::numeric, 0), COALESCE((data->>'cpu_cost')::numeric, 0), COALESCE((data->>'total_cost')::numeric, 0) FROM analytics WHERE type = 'DatabaseUsageStats' AND "projectId" = $1;`

	// Queries to get the last records for each type of analytics

	GET_LAST_EXECUTION_TIME_STATS = `SELECT created_at::text, COALESCE((data->>'total_time_ms')::numeric, 0), COALESCE((data->>'total_queries')::bigint, 0) FROM analytics WHERE type = 'ExecutionTimeStats' AND "projectId" = $1 ORDER BY created_at DESC LIMIT 1;`

	GET_LAST_DATABASE_USAGE_STATS = `SELECT created_at::text, COALESCE((data->>'read_write_cost')::numeric, 0), COALESCE((data->>'cpu_cost')::numeric, 0), COALESCE((data->>'total_cost')::numeric, 0) FROM analytics WHERE type = 'DatabaseUsageStats' AND "projectId" = $1 ORDER BY created_at DESC LIMIT 1;`
)
