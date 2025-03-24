package projects

var CheckUserHasProject = `SELECT EXISTS(SELECT 1 FROM "Project" WHERE owner_id = $1 AND name = $2)`
var InsertDatabaseConfig = `INSERT INTO "database_config" ("host", "port", "user_id", "password", "db_name", "ssl_mode", "created_at") VALUES ($1, $2, $3, $4, $5, $6, $7)`
var InsertDatabaseProjectData = `INSERT INTO "Project" (oid, owner_id, name, description, status, created_at, "API_URL", "API_key") VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
