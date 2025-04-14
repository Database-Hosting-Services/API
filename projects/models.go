package projects

import "time"

// Project struct is a row record of the projects table in the projects database
type Project struct {
	ID          int       `json:"id" db:"id"`
	Oid         string    `json:"oid" db:"oid"`
	OwnerID     int       `json:"owner_id" db:"owner_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	APIURL      string    `json:"API_URL" db:"api_url"`
	APIKey      string    `json:"API_key" db:"api_key"`
}

type SafeProjectData struct {
	Oid         string    `json:"oid" db:"oid"`
	OwnerID     int       `json:"owner_id" db:"owner_id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	Status      string    `json:"status" db:"status"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	APIURL      string    `json:"API_URL" db:"api_url"`
	APIKey      string    `json:"API_key" db:"api_key"`
}

type updateProjectDataModel struct {
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Status      string `json:"status" db:"status"`
}

type DatabaseConfig struct {
	ID        int       `json:"id" db:"id"`
	Host      string    `json:"host" db:"host"`
	Port      string    `json:"port" db:"port"`
	UserID    int       `json:"user_id" db:"user_id"`
	Password  string    `json:"password" db:"password"`
	DBName    string    `json:"db_name" db:"db_name"`
	SSLMode   string    `json:"ssl_mode" db:"ssl_mode"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

var DefaultProjectResponse = SafeProjectData{}
