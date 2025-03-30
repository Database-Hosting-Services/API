package projects

import (
	"DBHS/utils"
	"context"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CheckDatabaseExists(ctx context.Context, db *pgxpool.Pool, query string, SearchField ...interface{}) (bool, error) {
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

func getUserProjectsFromDatabase(ctx context.Context, db *pgxpool.Pool, userId int) ([]*SafeProjectData, error) {
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

func getUserSpecificProjectFromDatabase(ctx context.Context, db *pgxpool.Pool, userId int, projectOid string) (*SafeProjectData, error) {
	var project SafeProjectData
	err := pgxscan.Get(ctx, db, &project, RetrieveUserSpecificProject, userId, projectOid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &project, nil
}
