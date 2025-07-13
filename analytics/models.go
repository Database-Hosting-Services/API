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
	Timestamp    string  `json:"timestamp"`
	TotalTimeMs  float64 `json:"total_time_ms"`
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

type DatabaseUsageCostWithDates struct {
	Timestamp     string  `json:"timestamp"`
	ReadWriteCost float64 `json:"read_write_cost"`
	CPUCost       float64 `json:"cpu_cost"`
	TotalCost     float64 `json:"total_cost"`
}

var (
	// Storage dummy data - showing growth over time (all in KB)
	StorageResponse = []StorageWithDates{
		{
			Timestamp:         "2025-01-01T00:00:00Z",
			ManagementStorage: "128000 kB",
			ActualData:        "2202010 kB",
		},
		{
			Timestamp:         "2025-01-02T00:00:00Z",
			ManagementStorage: "131072 kB",
			ActualData:        "2411724 kB",
		},
		{
			Timestamp:         "2025-01-03T00:00:00Z",
			ManagementStorage: "135168 kB",
			ActualData:        "2730394 kB",
		},
		{
			Timestamp:         "2025-01-04T00:00:00Z",
			ManagementStorage: "138240 kB",
			ActualData:        "2939002 kB",
		},
		{
			Timestamp:         "2025-01-05T00:00:00Z",
			ManagementStorage: "144384 kB",
			ActualData:        "3257671 kB",
		},
		{
			Timestamp:         "2025-01-06T00:00:00Z",
			ManagementStorage: "147456 kB",
			ActualData:        "3576340 kB",
		},
		{
			Timestamp:         "2025-01-07T00:00:00Z",
			ManagementStorage: "151552 kB",
			ActualData:        "3895009 kB",
		},
		{
			Timestamp:         "2025-01-08T00:00:00Z",
			ManagementStorage: "155648 kB",
			ActualData:        "4194304 kB",
		},
		{
			Timestamp:         "2025-01-09T00:00:00Z",
			ManagementStorage: "159744 kB",
			ActualData:        "4403021 kB",
		},
		{
			Timestamp:         "2025-01-10T00:00:00Z",
			ManagementStorage: "163840 kB",
			ActualData:        "4718592 kB",
		},
	}

	// Database Activity dummy data - showing realistic patterns with peaks and valleys
	DatabaseActivityResponse = []DatabaseActivityWithDates{
		{
			Timestamp:    "2025-01-01T00:00:00Z",
			TotalTimeMs:  1245.67,
			TotalQueries: 1523,
		},
		{
			Timestamp:    "2025-01-02T00:00:00Z",
			TotalTimeMs:  1789.34,
			TotalQueries: 2156,
		},
		{
			Timestamp:    "2025-01-03T00:00:00Z",
			TotalTimeMs:  2134.89,
			TotalQueries: 2834,
		},
		{
			Timestamp:    "2025-01-04T00:00:00Z",
			TotalTimeMs:  1876.23,
			TotalQueries: 2445,
		},
		{
			Timestamp:    "2025-01-05T00:00:00Z",
			TotalTimeMs:  2567.45,
			TotalQueries: 3234,
		},
		{
			Timestamp:    "2025-01-06T00:00:00Z",
			TotalTimeMs:  2234.78,
			TotalQueries: 2967,
		},
		{
			Timestamp:    "2025-01-07T00:00:00Z",
			TotalTimeMs:  3045.12,
			TotalQueries: 3789,
		},
		{
			Timestamp:    "2025-01-08T00:00:00Z",
			TotalTimeMs:  2789.56,
			TotalQueries: 3456,
		},
		{
			Timestamp:    "2025-01-09T00:00:00Z",
			TotalTimeMs:  2456.89,
			TotalQueries: 3123,
		},
		{
			Timestamp:    "2025-01-10T00:00:00Z",
			TotalTimeMs:  2998.34,
			TotalQueries: 3678,
		},
	}

	// Database Usage Stats dummy data - showing cost variations
	DatabaseUsageStatsResponse = []DatabaseUsageCostWithDates{
		{
			ReadWriteCost: 85.45,
			CPUCost:       67.23,
			TotalCost:     152.68,
			Timestamp:     "2025-01-01T00:00:00Z",
		},
		{
			ReadWriteCost: 142.78,
			CPUCost:       89.34,
			TotalCost:     232.12,
			Timestamp:     "2025-01-02T00:00:00Z",
		},
		{
			ReadWriteCost: 198.56,
			CPUCost:       112.67,
			TotalCost:     311.23,
			Timestamp:     "2025-01-03T00:00:00Z",
		},
		{
			ReadWriteCost: 176.89,
			CPUCost:       95.45,
			TotalCost:     272.34,
			Timestamp:     "2025-01-04T00:00:00Z",
		},
		{
			ReadWriteCost: 234.67,
			CPUCost:       134.23,
			TotalCost:     368.90,
			Timestamp:     "2025-01-05T00:00:00Z",
		},
		{
			ReadWriteCost: 201.45,
			CPUCost:       118.78,
			TotalCost:     320.23,
			Timestamp:     "2025-01-06T00:00:00Z",
		},
		{
			ReadWriteCost: 287.34,
			CPUCost:       156.89,
			TotalCost:     444.23,
			Timestamp:     "2025-01-07T00:00:00Z",
		},
		{
			ReadWriteCost: 256.78,
			CPUCost:       143.56,
			TotalCost:     400.34,
			Timestamp:     "2025-01-08T00:00:00Z",
		},
		{
			ReadWriteCost: 213.45,
			CPUCost:       125.67,
			TotalCost:     339.12,
			Timestamp:     "2025-01-09T00:00:00Z",
		},
		{
			ReadWriteCost: 298.67,
			CPUCost:       167.89,
			TotalCost:     466.56,
			Timestamp:     "2025-01-10T00:00:00Z",
		},
	}
)
