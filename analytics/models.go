package analytics

// Storage represents the structure of the storage information returned by the analytics API.
type Storage struct {
	ManagementStorage string `json:"Management storage"`
	ActualData        string `json:"Actual data"`
}
