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
	// Storage dummy data - showing growth over time (all in KB, under 100)
	StorageResponse = []StorageWithDates{
		{
			Timestamp:         "2025-01-01T00:00:00Z",
			ManagementStorage: "12 kB",
			ActualData:        "45 kB",
		},
		{
			Timestamp:         "2025-01-02T00:00:00Z",
			ManagementStorage: "14 kB",
			ActualData:        "52 kB",
		},
		{
			Timestamp:         "2025-01-03T00:00:00Z",
			ManagementStorage: "16 kB",
			ActualData:        "58 kB",
		},
		{
			Timestamp:         "2025-01-04T00:00:00Z",
			ManagementStorage: "18 kB",
			ActualData:        "64 kB",
		},
		{
			Timestamp:         "2025-01-05T00:00:00Z",
			ManagementStorage: "21 kB",
			ActualData:        "71 kB",
		},
		{
			Timestamp:         "2025-01-06T00:00:00Z",
			ManagementStorage: "24 kB",
			ActualData:        "77 kB",
		},
		{
			Timestamp:         "2025-01-07T00:00:00Z",
			ManagementStorage: "27 kB",
			ActualData:        "84 kB",
		},
		{
			Timestamp:         "2025-01-08T00:00:00Z",
			ManagementStorage: "30 kB",
			ActualData:        "89 kB",
		},
		{
			Timestamp:         "2025-01-09T00:00:00Z",
			ManagementStorage: "33 kB",
			ActualData:        "95 kB",
		},
		{
			Timestamp:         "2025-01-10T00:00:00Z",
			ManagementStorage: "36 kB",
			ActualData:        "99 kB",
		},
	}

	// Database Activity dummy data - showing realistic patterns with peaks and valleys (under 100)
	DatabaseActivityResponse = []DatabaseActivityWithDates{
		{
			Timestamp:    "2025-01-01T00:00:00Z",
			TotalTimeMs:  12.45,
			TotalQueries: 15,
		},
		{
			Timestamp:    "2025-01-02T00:00:00Z",
			TotalTimeMs:  17.89,
			TotalQueries: 21,
		},
		{
			Timestamp:    "2025-01-03T00:00:00Z",
			TotalTimeMs:  21.34,
			TotalQueries: 28,
		},
		{
			Timestamp:    "2025-01-04T00:00:00Z",
			TotalTimeMs:  18.76,
			TotalQueries: 24,
		},
		{
			Timestamp:    "2025-01-05T00:00:00Z",
			TotalTimeMs:  25.67,
			TotalQueries: 32,
		},
		{
			Timestamp:    "2025-01-06T00:00:00Z",
			TotalTimeMs:  22.34,
			TotalQueries: 29,
		},
		{
			Timestamp:    "2025-01-07T00:00:00Z",
			TotalTimeMs:  30.45,
			TotalQueries: 37,
		},
		{
			Timestamp:    "2025-01-08T00:00:00Z",
			TotalTimeMs:  27.89,
			TotalQueries: 34,
		},
		{
			Timestamp:    "2025-01-09T00:00:00Z",
			TotalTimeMs:  24.56,
			TotalQueries: 31,
		},
		{
			Timestamp:    "2025-01-10T00:00:00Z",
			TotalTimeMs:  29.98,
			TotalQueries: 36,
		},
	}

	// Database Usage Stats dummy data - showing cost variations (under 100)
	DatabaseUsageStatsResponse = []DatabaseUsageCostWithDates{
		{
			ReadWriteCost: 25.45,
			CPUCost:       17.23,
			TotalCost:     42.68,
			Timestamp:     "2025-01-01T00:00:00Z",
		},
		{
			ReadWriteCost: 32.78,
			CPUCost:       19.34,
			TotalCost:     52.12,
			Timestamp:     "2025-01-02T00:00:00Z",
		},
		{
			ReadWriteCost: 38.56,
			CPUCost:       22.67,
			TotalCost:     61.23,
			Timestamp:     "2025-01-03T00:00:00Z",
		},
		{
			ReadWriteCost: 36.89,
			CPUCost:       25.45,
			TotalCost:     62.34,
			Timestamp:     "2025-01-04T00:00:00Z",
		},
		{
			ReadWriteCost: 44.67,
			CPUCost:       34.23,
			TotalCost:     78.90,
			Timestamp:     "2025-01-05T00:00:00Z",
		},
		{
			ReadWriteCost: 41.45,
			CPUCost:       28.78,
			TotalCost:     70.23,
			Timestamp:     "2025-01-06T00:00:00Z",
		},
		{
			ReadWriteCost: 47.34,
			CPUCost:       36.89,
			TotalCost:     84.23,
			Timestamp:     "2025-01-07T00:00:00Z",
		},
		{
			ReadWriteCost: 46.78,
			CPUCost:       33.56,
			TotalCost:     80.34,
			Timestamp:     "2025-01-08T00:00:00Z",
		},
		{
			ReadWriteCost: 43.45,
			CPUCost:       35.67,
			TotalCost:     79.12,
			Timestamp:     "2025-01-09T00:00:00Z",
		},
		{
			ReadWriteCost: 48.67,
			CPUCost:       37.89,
			TotalCost:     86.56,
			Timestamp:     "2025-01-10T00:00:00Z",
		},
	}
)
