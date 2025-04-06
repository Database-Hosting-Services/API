package config

import (
	"strings"
)

// ParseDatabaseURL parses a PostgreSQL connection string into its components
func ParseDatabaseURL(dbURL string) *DatabaseConfig {
	// Default values
	config := &DatabaseConfig{
		Host:    "localhost",
		Port:    "5432",
		User:    "postgres",
		DBName:  "postgres",
		SSLMode: "disable",
	}

	// Remove postgres:// prefix if present
	dbURL = strings.TrimPrefix(dbURL, "postgres://")

	// Split user:password@host:port/dbname
	parts := strings.Split(dbURL, "@")
	if len(parts) > 1 {
		// Handle user:password
		userPass := strings.Split(parts[0], ":")
		if len(userPass) > 0 {
			config.User = userPass[0]
		}
		if len(userPass) > 1 {
			config.Password = userPass[1]
		}

		// Handle host:port/dbname
		hostPortDB := parts[1]
		hostPortParts := strings.Split(hostPortDB, "/")
		if len(hostPortParts) > 0 {
			hostPort := strings.Split(hostPortParts[0], ":")
			if len(hostPort) > 0 {
				config.Host = hostPort[0]
			}
			if len(hostPort) > 1 {
				config.Port = hostPort[1]
			}
		}
		if len(hostPortParts) > 1 {
			// Handle dbname?params
			dbNameParams := strings.Split(hostPortParts[1], "?")
			if len(dbNameParams) > 0 {
				config.DBName = dbNameParams[0]
			}
			if len(dbNameParams) > 1 {
				// Extract sslmode if present
				params := strings.Split(dbNameParams[1], "&")
				for _, param := range params {
					if strings.HasPrefix(param, "sslmode=") {
						config.SSLMode = strings.TrimPrefix(param, "sslmode=")
					}
				}
			}
		}
	}
	return config
}
