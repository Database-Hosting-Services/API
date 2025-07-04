package projects

import (
	"errors"
	"fmt"
	"strings"
)

const (
	CheckUserHasProject        = `SELECT EXISTS(SELECT 1 FROM "projects" WHERE owner_id = $1 AND name = $2)`
	CheckUserHasProjectWithOid = `SELECT EXISTS(SELECT 1 FROM "projects" WHERE owner_id = $1 AND oid = $2)`

	InsertDatabaseConfig      = `INSERT INTO "database_config" ("host", "port", "user_id", "password", "db_name", "ssl_mode", "created_at") VALUES ($1, $2, $3, $4, $5, $6, $7)`
	InsertDatabaseProjectData = `INSERT INTO "projects" (oid, owner_id, name, description, status, created_at, api_url, api_key) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	// in future plans this query will return also the user's projects where the user is member not only owner
	RetrieveUserProjects = `SELECT oid, owner_id, name, description, status, created_at, api_url, api_key
							FROM "projects"
							WHERE owner_id = $1`

	RetrieveUserSpecificProject = `SELECT oid, owner_id, name, description, status, created_at, api_url, api_key
									FROM "projects"
									WHERE owner_id = $1 AND oid = $2`

	RetrieveProjectID = `
		SELECT id FROM projects WHERE owner_id = $1 AND oid = $2
	`
)

func BuildProjectUpdateQuery(projectOid string, feildsToUpdate []string) (string, error) {
	if len(feildsToUpdate) == 0 {
		return "", errors.New("no fields provided for update")
	}

	query := `UPDATE "projects" SET `
	setClauses := []string{}

	index := 1
	for _, field := range feildsToUpdate {
		setClauses = append(setClauses, fmt.Sprintf(`%s = $%d`, field, index))
		index++
	}

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(` WHERE oid = '%s'`, projectOid)
	return query, nil
}
