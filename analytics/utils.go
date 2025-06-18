package analytics

// CalculateCosts calculates the costs associated with database usage based on read/write queries and CPU time.
func (d *DatabaseUsageStats) CalculateCosts() DatabaseUsageCost {
	Cost := DatabaseUsageCost{
		ReadWriteCost: float64(d.ReadQueries)/1_000_000*1.00 + float64(d.WriteQueries)/1_000_000*1.50,
		CPUCost:       (d.TotalCPUTimeMs / 1000 / 3600) * 0.000463,
	}
	Cost.TotalCost = Cost.ReadWriteCost + Cost.CPUCost
	return Cost
}
