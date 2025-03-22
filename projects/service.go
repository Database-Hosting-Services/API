package projects

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUserProject(ctx context.Context, db *pgxpool.Pool, projectname string) (bool, error) {
	UserId, ok := ctx.Value("user-id").(string)
	if !ok || UserId == "" {
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

	return true, nil
}
