package analytics

// Storage represents the structure of the storage information returned by the analytics API.
type Storage struct {
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
