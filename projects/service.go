package projects

import (
	"DBHS/config"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUserProject(ctx context.Context, db *pgxpool.Pool, projectname string) (bool, error) {
	UserId, ok := ctx.Value("user-id").(int)
	if !ok || UserId == 0 {
		return false, errors.New("Unauthorized")
	}

	// Check if the project already exists
	Has, err := CheckDatabaseExists(ctx, db, CheckUserHasProject, UserId, projectname)
	if err != nil {
		return false, err
	}

	if Has {
		return false, errors.New("Project already exists")
	}

	// Check if the ProjectName is valid (must not be a reserved name and other validation)
	err = ValidatePostgresDatabaseName(projectname)
	if err != nil {
		return false, err
	}

	_, err = config.AdminDB.Exec(ctx, "CREATE DATABASE " + projectname)
	if err != nil {
		return false, err
	}
	return false, nil
}
