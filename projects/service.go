package projects

import (
	"DBHS/config"
	"DBHS/utils"
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
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

	// Begin transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer tx.Rollback(ctx)

	// --------------------------- Database Connection Config ---------------------------
	projectDBConfig := CreateDatabaseConfig(projectname, UserId)

	// Insert the new project Config into the database using the transaction
	err = InsertNewRecord(ctx, tx, InsertDatabaseConfig,
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

	// Insert the new project data into the database using the transaction
	err = InsertNewRecord(ctx, tx, InsertDatabaseProjectData,
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

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return false, err
	}

	// --------------------------- Create the Database -----------------------------------

	DBname := projectname + "_" + strconv.Itoa(UserId)
	_, err = config.AdminDB.Exec(ctx, "CREATE DATABASE "+DBname)
	if err != nil {
		return false, err
	}

	return false, nil
}

func getUserProjects(ctx context.Context, db *pgxpool.Pool, userId int) ([]*SafeReadProject, error) {
	projects, err := getUserProjectsFromDatabase(ctx, db, userId)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func getUserSpecificProject(ctx context.Context, db *pgxpool.Pool, userId int, projectOid string) (*SafeReadProject, error) {
	project, err := getUserSpecificProjectFromDatabase(ctx, db, userId, projectOid)
	if err != nil {
		return nil, err
	}
	return project, nil
}
