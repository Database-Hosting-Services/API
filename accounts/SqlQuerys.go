package accounts

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
)
