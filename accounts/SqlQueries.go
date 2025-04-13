package accounts

import (
	"errors"
	"fmt"
	"strings"
)

const (
	SELECT_USER_BY_Email         = `SELECT id, oid, username, email, password, image, created_at, last_login FROM "users" WHERE email = $1`
	SELECT_USER_BY_ID            = `SELECT id, oid, username, email, password, image, created_at, last_login FROM "users" WHERE oid = $1`
	SELECT_ID_FROM_USER_BY_EMAIL = `SELECT id FROM "users" WHERE email = $1`
	luaDeleteScript              = `
	local emailDeleted = redis.call("DEL", KEYS[1])
	local usernameDeleted = redis.call("DEL", KEYS[2])

	if emailDeleted == 0 or usernameDeleted == 0 then
		return "ERROR"
	end

	return "OK"
	`
	UPDATE_USER_PASSWORD = `UPDATE "users" SET password = $1 WHERE oid = $2`
)

func BuildUserUpdateQuery(userOid string, updatedFields []string) (string, error) {
	if len(updatedFields) == 0 {
		return "", errors.New("no fields provided for update")
	}

	query := `UPDATE "users" SET `
	setClauses := []string{}

	index := 1
	for _, field := range updatedFields {
		setClauses = append(setClauses, fmt.Sprintf(`%s = $%d`, field, index))
		index++
	}

	query += strings.Join(setClauses, ", ")
	query += fmt.Sprintf(` WHERE oid = '%s'`, userOid)
	return query, nil
}
