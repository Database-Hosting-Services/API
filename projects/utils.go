package projects

import (
	"DBHS/config"
	"DBHS/utils"
	"context"
	"errors"
	"regexp"
	"strings"
	"time"
)

var validName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_$]*$`).MatchString
var reservedNames = []string{
	"postgres", "template0", "template1", "admin", "public", "system", "information_schema", "template", "global", "test", "tmp", "temp",
}

func ValidatePostgresDatabaseName(name string) error {
	if len(name) < 3 || len(name) > 63 {
		return errors.New("database name must be between 3 and 63 characters")
	}

	if !validName(name) {
		return errors.New("database name must start with a letter or underscore and contain only letters, numbers, underscores, or $")
	}

	nameLower := strings.ToLower(name)
	for _, reserved := range reservedNames {
		if nameLower == reserved {
			return errors.New("database name is reserved: " + name)
		}
	}

	if strings.HasPrefix(nameLower, "pg_") {
		return errors.New("database names starting with 'pg_' are reserved")
	}

	return nil
}

func validateProjectData(ctx context.Context, db utils.Querier, projectname string, UserId int) error {
	// Check if the project already exists
	Has, err := CheckDatabaseExists(ctx, db, CheckUserHasProject, UserId, projectname)
	if err != nil {
		return err
	}

	if Has {
		return errors.New("Project already exists")
	}

	// Check if the ProjectName is valid (must not be a reserved name and other validation)
	err = ValidatePostgresDatabaseName(projectname)
	if err != nil {
		return err
	}
	return nil
}

func CreateDatabaseConfig(dbName string, userId int) DatabaseConfig {
	return DatabaseConfig{
		Host:      config.DBConfig.Host,
		Port:      config.DBConfig.Port,
		UserID:    userId,
		Password:  config.DBConfig.Password,
		DBName:    dbName,
		SSLMode:   config.DBConfig.SSLMode,
		CreatedAt: time.Now().Format(time.RFC3339), // default time format like "2006-01-02T15:04:05Z07:00"
	}
}

// TODO: support for generating API key
func GenerateApiKey() string {
	return ""
}

// TODO: support for generating API url
func GenerateApiUrl(databaseConfig DatabaseConfig) string {
	return "https://" + databaseConfig.Host + ":" + databaseConfig.Port + "/" + databaseConfig.DBName
}

func CreateDatabaseProjectData(oid, name, description, status string, ownerID int, databaseConfig DatabaseConfig) SafeProjectData {
	return SafeProjectData{
		Oid:         oid,
		OwnerID:     ownerID,
		Name:        name,
		Description: description,
		Status:      status,
		APIURL:      GenerateApiUrl(databaseConfig),
		APIKey:      GenerateApiKey(),
		CreatedAt:   time.Time{}, //اعتمد ان اي تايم بيبقي من تايب time.Time احسن
	}
}
