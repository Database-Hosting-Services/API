package accounts

var SELECT_USER_BY_Email = `SELECT id, oid, username, email, password, image, created_at, last_login FROM "users" WHERE email = $1`
var SELECT_USER_BY_ID = `SELECT id, oid, username, email, password, image, created_at, last_login FROM "users" WHERE oid = $1`
var SELECT_ID_FROM_USER_BY_EMAIL = `SELECT id FROM "users" WHERE email = $1`
