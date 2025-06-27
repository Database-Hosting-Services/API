package analytics

// ----------------------- Storage represents the structure of the storage information returned by the analytics API --------
type Storage struct {
	ManagementStorage string `json:"Management storage"`
	ActualData        string `json:"Actual data"`
}

type StorageWithDates struct {
	Timestamp         string `json:"timestamp"`
	ManagementStorage string `json:"Management storage"`
	ActualData        string `json:"Actual data"`
}

// ------------------------------ DatabaseActivity represents statistics about database activity ------------------------------
type DatabaseActivity struct {
	TotalTimeMs  float64 `json:"total_time_ms"`
	TotalQueries int64   `json:"total_queries"`
}

type DatabaseActivityWithDates struct {
	Timestamp   string  `json:"timestamp"`
	TotalTimeMs float64 `json:"total_time_ms"`
	TotalQueries int64   `json:"total_queries"`
}

// ------------------------------ represents database read/write statistics and cost calculations ------------------------------
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
