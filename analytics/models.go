package analytics

// Storage represents the structure of the storage information returned by the analytics API.
type Storage struct {
	ManagementStorage string `json:"Management storage"`
	ActualData        string `json:"Actual data"`
}

type StorageWithDates struct {
	Timestamp        string `json:"timestamp"`
	ManagementStorage string `json:"Management storage"`
	ActualData        string `json:"Actual data"`
}

// ExecutionTimeStats represents statistics about query execution times
type ExecutionTimeStats struct {
	TotalTimeMs float64 `json:"total_time_ms"`
	MaxTimeMs   float64 `json:"max_time_ms"`
	AvgTimeMs   float64 `json:"avg_time_ms"`
}

// ------ represents database read/write statistics and cost calculations ------
type DatabaseUsageStats struct {
	ReadQueries    int64   `json:"read_queries"`
	WriteQueries   int64   `json:"write_queries"`
	TotalCPUTimeMs float64 `json:"total_cpu_time_ms"`
}

type DatabaseUsageCost struct {
	ReadWriteCost float64 `json:"read_write_cost"`
	CPUCost       float64 `json:"cpu_cost"`
	TotalCost     float64 `json:"total_cost"`
}


/*
  1 | 2025-06-26 04:19:15.048835+00 |       177 | Storage            | {"Management storage":"1596 kB","Actual data":"37 kB"}
  2 | 2025-06-26 04:19:15.120183+00 |       177 | ExecutionTimeStats | {"total_time_ms":4.02,"max_time_ms":2.61,"avg_time_ms":0.37}
  3 | 2025-06-26 04:19:15.18474+00  |       177 | DatabaseUsageStats | {"read_write_cost":0.000013,"cpu_cost":5.414527777777778e-10,"total_cost":0.000013000541452777776}
  4 | 2025-06-26 04:25:51.606299+00 |       177 | Storage            | {"Management storage":"1596 kB","Actual data":"37 kB"}
  5 | 2025-06-26 04:25:51.672108+00 |       177 | ExecutionTimeStats | {"total_time_ms":3.28,"max_time_ms":2.4,"avg_time_ms":0.82}
  6 | 2025-06-26 04:25:51.736408+00 |       177 | DatabaseUsageStats | {"read_write_cost":0.000006,"cpu_cost":4.4242222222222224e-10,"total_cost":0.000006000442422222223}
  7 | 2025-06-27 00:00:11.746163+00 |       177 | Storage            | {"Management storage":"1596 kB","Actual data":"37 kB"}
  8 | 2025-06-27 00:00:11.782553+00 |       177 | ExecutionTimeStats | {"total_time_ms":5.73,"max_time_ms":4.05,"avg_time_ms":1.43}
  9 | 2025-06-27 00:00:11.79325+00  |       177 | DatabaseUsageStats | {"read_write_cost":0.000006,"cpu_cost":7.536611111111111e-10,"total_cost":0.000006000753661111111}
*/