package indexes

import (
	"DBHS/config"
	"DBHS/projects"
	"DBHS/utils"
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateIndexInDatabase(ctx context.Context, db *pgxpool.Pool, projectOid string, indexData IndexData) error {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return errors.New("Unauthorized")
	}

	// ------------------------ Get the project database connection ------------------------
	projectDB, err := projects.GetUserSpecificProject(ctx, db, UserID, projectOid)
	if err != nil {
		return err
	}

	if projectDB == nil {
		return errors.New("project not found")
	}

	// ------------------------ Get The project connection Pool ------------------------

	DBname := strings.ToLower(projectDB.Name) + "_" + strconv.Itoa(UserID)
	conn, err := config.ConfigManager.GetDbConnection(ctx, DBname)
	if err != nil {
		return err
	}
	defer conn.Close()

	// ------------------------ Create the index in the database ------------------------

	query := GenerateIndexQuery(indexData)
	if _, err = conn.Exec(ctx, query); err != nil {
		return err
	}

	return nil
}

func GetIndexes(ctx context.Context, db *pgxpool.Pool, projectOid string) ([]RetrievedIndex, error) {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return nil, errors.New("Unauthorized")
	}

	// ------------------------ Get the project database connection ------------------------
	projectDB, err := projects.GetUserSpecificProject(ctx, db, UserID, projectOid)
	if err != nil {
		return nil, err
	}

	if projectDB == nil {
		return nil, errors.New("project not found")
	}

	// ------------------------ Get The project connection Pool ------------------------

	DBname := strings.ToLower(projectDB.Name) + "_" + strconv.Itoa(UserID)
	conn, err := config.ConfigManager.GetDbConnection(ctx, DBname)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	// ------------------------ Get the indexes from the database ------------------------

	indexes, err := GetProjectIndexes(ctx, conn)
	if err != nil {
		return nil, err
	}

	return indexes, nil
}

func GetSpecificIndex(ctx context.Context, db *pgxpool.Pool, projectOid, indexOid string) (SpecificIndex, error) {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return DefaultSpecificIndex, errors.New("Unauthorized")
	}

	// ------------------------ Get the project database connection ------------------------
	projectDB, err := projects.GetUserSpecificProject(ctx, db, UserID, projectOid)
	if err != nil {
		return DefaultSpecificIndex, err
	}

	if projectDB == nil {
		return DefaultSpecificIndex, errors.New("project not found")
	}

	// ------------------------ Get The project connection Pool ------------------------
	DBname := strings.ToLower(projectDB.Name) + "_" + strconv.Itoa(UserID)
	conn, err := config.ConfigManager.GetDbConnection(ctx, DBname)
	if err != nil {
		return DefaultSpecificIndex, err
	}
	defer conn.Close()

	// ------------------------ Get the index from the database ------------------------

	index := GetSpecificIndexFromDatabase(ctx, conn, indexOid)
	if index == (SpecificIndex{}) {
		return DefaultSpecificIndex, errors.New("index not found")
	}

	// ------------------------ Close the connection ------------------------
	return index, nil
}

func DeleteSpecificIndex(ctx context.Context, db *pgxpool.Pool, projectOid, indexOid string) error {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return errors.New("Unauthorized")
	}

	// ------------------------ Get the project database connection ------------------------
	projectDB, err := projects.GetUserSpecificProject(ctx, db, UserID, projectOid)
	if err != nil {
		return err
	}

	if projectDB == nil {
		return errors.New("project not found")
	}

	// ------------------------ Get The project connection Pool ------------------------

	DBname := strings.ToLower(projectDB.Name) + "_" + strconv.Itoa(UserID)
	conn, err := config.ConfigManager.GetDbConnection(ctx, DBname)
	if err != nil {
		return err
	}
	defer conn.Close()

	// ------------------------ Delete the index from the database ------------------------

	IndexData := GetSpecificIndexFromDatabase(ctx, conn, indexOid)
	if IndexData == (SpecificIndex{}) {
		return errors.New("index not found")
	}

	err = DeleteIndexFromDatabase(ctx, conn, IndexData.IndexName)

	if err != nil {
		return err
	}
	return nil
}

func UpdateSpecificIndex(ctx context.Context, db *pgxpool.Pool, projectOid, indexOid string, newName string) error {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return errors.New("Unauthorized")
	}

	// ------------------------ Get the project database connection ------------------------
	projectDB, err := projects.GetUserSpecificProject(ctx, db, UserID, projectOid)
	if err != nil {
		return err
	}

	if projectDB == nil {
		return errors.New("project not found")
	}

	// ------------------------ Get The project connection Pool ------------------------

	DBname := strings.ToLower(projectDB.Name) + "_" + strconv.Itoa(UserID)
	conn, err := config.ConfigManager.GetDbConnection(ctx, DBname)
	if err != nil {
		return err
	}
	defer conn.Close()

	// ------------------------ Update the index in the database ------------------------
	newName = utils.ReplaceWhiteSpacesWithUnderscore(newName)

	// Get the current index name
	indexData := GetSpecificIndexFromDatabase(ctx, conn, indexOid)
	if indexData == (SpecificIndex{}) {
		return errors.New("index not found")
	}

	if indexData.IndexName == newName {
		return errors.New("index name is the same as the current name")
	}

	err = UpdateIndexNameInDatabase(ctx, conn, indexData.IndexName, newName)
	if err != nil {
		fmt.Println("err:", err.Error())
		return err
	}

	return nil
}
