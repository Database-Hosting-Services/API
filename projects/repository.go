package projects

import (
	"DBHS/utils"
	"context"
	"errors"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckDatabaseExists(ctx context.Context, db utils.Querier, query string, SearchField ...interface{}) (bool, error) {
	var exists bool
	err := db.QueryRow(ctx, query, SearchField...).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func InsertNewRecord(ctx context.Context, db utils.Querier, query string, values ...interface{}) error {
	_, err := db.Exec(ctx, query, values...)
	if err != nil {
		return err
	}
	return nil
}

func getUserProjectsFromDatabase(ctx context.Context, db *pgxpool.Pool, userId int64) ([]*SafeProjectData, error) {
	var projects []*SafeProjectData
	err := pgxscan.Select(
		ctx, db, &projects,
		RetrieveUserProjects,
		userId,
	)

	if err != nil {
		return nil, err
	}
	return projects, nil
}

func getUserSpecificProjectFromDatabase(ctx context.Context, db utils.Querier, userId int64, projectOid string) (*SafeProjectData, error) {
	var project SafeProjectData
	err := pgxscan.Get(ctx, db, &project, RetrieveUserSpecificProject, userId, projectOid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrorProjectNotFound
		}
		return nil, err
	}
	return &project, nil
}


func GetProjectNameID(ctx context.Context, projectId string, db utils.Querier) (interface{}, interface{}, error) {
	var name, id interface{}
	err := db.QueryRow(ctx, "SELECT id, name FROM projects WHERE oid = $1", projectId).Scan(&id, &name)
	if err != nil {
		return nil, nil, err
	}
	return name, id, nil
}

func GetProjectID(ctx context.Context, db utils.Querier, userId int64, projectOid string) (int64, error) {
	var projectID int64
	err := pgxscan.Get(ctx, db, &projectID, RetrieveProjectID, userId, projectOid)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, ErrorProjectNotFound
		}
		return 0, err
	}
	return projectID, nil
}