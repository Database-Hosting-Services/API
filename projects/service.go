package projects

import (
	"DBHS/config"
	"DBHS/utils"
	"context"
	"errors"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateUserProject(ctx context.Context, db *pgxpool.Pool, projectname, projectDescription string) (bool, error, SafeProjectData) {
	UserId, ok := ctx.Value("user-id").(int)
	if !ok || UserId == 0 {
		return false, errors.New("Unauthorized"), DefaultProjectResponse
	}

	// Replace white spaces with underscores
	projectname = utils.ReplaceWhiteSpacesWithUnderscore(projectname)

	err := validateProjectData(ctx, db, projectname, UserId)
	if err != nil {
		return false, err, SafeProjectData{}
	}

	// Begin transaction
	tx, err := db.Begin(ctx)
	if err != nil {
		return false, err, DefaultProjectResponse
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
		return false, err, DefaultProjectResponse
	}

	// --------------------------- Database Project Data --------------------------------

	oid := utils.GenerateOID()
	//fmt.Println(projectname)
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
		return false, err, DefaultProjectResponse
	}

	// --------------------------- Create the Database -----------------------------------

	DBname := projectname + "_" + strconv.Itoa(UserId)
	_, err = config.AdminDB.Exec(ctx, "CREATE DATABASE "+DBname)
	if err != nil {
		return false, err, DefaultProjectResponse
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return false, err, DefaultProjectResponse
	}

	return false, nil, projectDBData
}

// DeleteUserProject handles the business logic for deleting a project
func DeleteUserProject(ctx context.Context, db *pgxpool.Pool, projectOID string) error {
	// Get user ID from context to check if user is owner for this project
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return errors.New("Unauthorized")
	}

	// ---------------------- Begin transaction ----------------------
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// ---------------------- Get The Project Name ----------------------
	project, err := GetUserSpecificProject(ctx, tx, UserID, projectOID)
	if err != nil {
		return errors.New("Project not found")
	}

	projectName := project.Name

	// ---------------------- Delete project data from projects table ----------------------
	_, err = tx.Exec(ctx, "DELETE FROM projects WHERE oid = $1", projectOID)
	if err != nil {
		return err
	}

	// ---------------------- Delete database configuration ----------------------
	_, err = tx.Exec(ctx, "DELETE FROM database_config WHERE db_name = $1 AND user_id = $2", projectName, UserID)
	if err != nil {
		return err
	}

	// ---------------------- Drop the actual database ----------------------
	DBname := projectName + "_" + strconv.Itoa(UserID)
	_, err = config.AdminDB.Exec(ctx, "DROP DATABASE IF EXISTS "+DBname)
	if err != nil {
		return err
	}

	// Commit the transaction
	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func getUserProjects(ctx context.Context, db *pgxpool.Pool, userId int) ([]*SafeProjectData, error) {
	projects, err := getUserProjectsFromDatabase(ctx, db, userId)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func GetUserSpecificProject(ctx context.Context, db utils.Querier, userId int, projectOid string) (*SafeProjectData, error) {
	project, err := getUserSpecificProjectFromDatabase(ctx, db, userId, projectOid)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func updateProjectData(ctx context.Context, transaction pgx.Tx, query string, values []interface{}) error {
	if err := utils.UpdateDataInDatabase(ctx, transaction, query, values...); err != nil {
		return err
	}
	return nil
}
