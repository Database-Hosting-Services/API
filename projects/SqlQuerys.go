package projects

var CheckUserHasProject = `SELECT EXISTS(SELECT 1 FROM "Project" WHERE owner_id = $1 AND name = $2)`
