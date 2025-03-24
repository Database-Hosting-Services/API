package projects

import (
	"DBHS/config"
	"context"
	"errors"
	"DBHS/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUserProject(ctx context.Context, db *pgxpool.Pool, projectname, projectDescription string) (bool, error) {
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

	// --------------------------- Database Connection Config ---------------------------
	projectDBConfig := CreateDatabaseConfig(projectname, UserId)

	// // Insert the new project Config into the database
	err = InsertNewRecord(ctx, db, InsertDatabaseConfig,
		projectDBConfig.Host,
		projectDBConfig.Port,
		projectDBConfig.UserID,
		projectDBConfig.Password,
		projectDBConfig.DBName,
		projectDBConfig.SSLMode,
		projectDBConfig.CreatedAt,
	)

	if err != nil {
		return false, err
	}

	// --------------------------- Database Project Data --------------------------------

	oid := utils.GenerateOID()
	projectDBData := CreateDatabaseProjectData(oid, projectname, projectDescription, "active", UserId, projectDBConfig)
	
	// Insert the new project data into the database
	err = InsertNewRecord(ctx, db, InsertDatabaseProjectData,
		projectDBData.Oid,
		projectDBData.OwnerID,
		projectDBData.Name,
		projectDBData.Description,
		projectDBData.Status,
		projectDBData.CreatedAt,
		projectDBData.APIURL,
		projectDBData.APIKey,
	)

	if err != nil {
		return false, err
	}

	// --------------------------- Create the Database -----------------------------------

	_, err = config.AdminDB.Exec(ctx, "CREATE DATABASE " + projectname)
	if err != nil {
		return false, err
	}

	return false, nil
}
