package projects


import (
    "regexp"
    "strings"
	"errors"
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