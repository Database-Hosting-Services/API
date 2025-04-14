package indexes

import (
	"DBHS/utils"
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	api "DBHS/apiError"

	"github.com/jackc/pgx/v5"

	"DBHS/config"
)

func CreateIndexInDatabase(ctx context.Context, db *pgxpool.Pool, projectOid string, indexData IndexData) api.ApiError {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return *api.NewApiError("Unauthorized", 401, errors.New("user is not authorized"))
	}

	// ------------------------ Get the project pool connection ------------------------
	conn, err := ProjectPoolConnection(ctx, db, UserID, projectOid)
	if err != nil {
		if err == pgx.ErrNoRows {
			return *api.NewApiError("Project not found", 404, errors.New("project was this id not found"))
		}
		return *api.NewApiError("Internal server error", 500, errors.New("failed to connect to the project"))
	}
	defer conn.Close()

	// ------------------------ Create the index in the database ------------------------

	query := GenerateIndexQuery(indexData)
	if _, err = conn.Exec(ctx, query); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return *api.NewApiError("the index name must be unique", 400, errors.New("index name already exists"))
		}
		return *api.NewApiError("Internal server error", 500, errors.New("failed to create the index"))
	}

	config.App.InfoLog.Println("Index created successfully for project:", projectOid)
	return *api.NewApiError("Index created successfully", 200, nil)
}

func GetIndexes(ctx context.Context, db *pgxpool.Pool, projectOid string) ([]RetrievedIndex, error) {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return nil, errors.New("Unauthorized")
	}

	// ------------------------ Get the project pool connection ------------------------
	conn, err := ProjectPoolConnection(ctx, db, UserID, projectOid)
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

	// ------------------------ Get the project pool connection ------------------------
	conn, err := ProjectPoolConnection(ctx, db, UserID, projectOid)
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

	// ------------------------ Get the project pool connection ------------------------
	conn, err := ProjectPoolConnection(ctx, db, UserID, projectOid)
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

	// ------------------------ Get the project pool connection ------------------------
	conn, err := ProjectPoolConnection(ctx, db, UserID, projectOid)
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
		return err
	}

	return nil
}
