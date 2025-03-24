package projects

// Project struct is a row record of the projects table in the projects database
type Project struct {
	ID          int    `json:"id"`
	Oid         string `json:"oid"`
	OwnerID     int    `json:"owner_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
	CreatedAt   string `json:"created_at"`
	APIURL      string `json:"api_url"`
	APIKey      string `json:"api_key"`
}

type DatabaseConfig struct {
	ID        int    `json:"id"`
	Host      string `json:"host"`
	Port      string `json:"port"`
	UserID    int    `json:"user_id"`
	Password  string `json:"password"`
	DBName    string `json:"db_name"`
	SSLMode   string `json:"ssl_mode"`
	CreatedAt string `json:"created_at"`
}
