package config

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
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

// ----------------------------------- Database Connection Pooling -----------------------------------

// NewUserDbConfig creates a new instance of UserDbConfig with a base configuration.
func NewUserDbConfig(baseConnString string) (*UserDbConfig, error) {
	// Parse the base connection string to create a base configuration.
	baseConfig, err := pgxpool.ParseConfig(baseConnString)
	if err != nil {
		return nil, fmt.Errorf("unable to parse base connection string: %w", err)
	}

	return &UserDbConfig{
		baseConfig: baseConfig,
	}, nil
}

// GetPool creates a connection pool for the specified database.
func (m *UserDbConfig) GetDbConnection(ctx context.Context, dbName string) (*pgxpool.Pool, error) {
	// Clone the base configuration for the database.
	newConfig := m.baseConfig.Copy()
	newConfig.ConnConfig.Database = dbName

	// Create a new connection pool with the updated configuration.
	newPool, err := pgxpool.NewWithConfig(ctx, newConfig)
	if err != nil {
		return nil, fmt.Errorf("unable to create connection pool for database %s: %w", dbName, err)
	}

	// ping the database
	err = newPool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to ping database %s: %w", dbName, err)
	}

	return newPool, nil
}
