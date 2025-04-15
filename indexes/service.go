package indexes

import (
	"DBHS/utils"
	"context"
	"errors"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"DBHS/config"
	api "DBHS/utils/apiError"
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
		if err.Error() == "Project not found" || err.Error() == "connection pool not found" {
			return *api.NewApiError("Project not found", 404, errors.New(err.Error()))
		}
		return *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}
	defer conn.Close()

	// ------------------------ Create the index in the database ------------------------

	query := GenerateIndexQuery(indexData)
	if _, err = conn.Exec(ctx, query); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return *api.NewApiError("the index name must be unique", 400, errors.New(err.Error()))
		}
		return *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}

	config.App.InfoLog.Println("Index created successfully for project:", projectOid)
	return *api.NewApiError("Index created successfully", 200, nil)
}

func GetIndexes(ctx context.Context, db *pgxpool.Pool, projectOid string) ([]RetrievedIndex, api.ApiError) {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return nil, *api.NewApiError("Unauthorized", 401, errors.New("user is not authorized"))
	}

	// ------------------------ Get the project pool connection ------------------------
	conn, err := ProjectPoolConnection(ctx, db, UserID, projectOid)
	if err != nil {
		if err.Error() == "Project not found" || err.Error() == "connection pool not found" {
			return nil, *api.NewApiError("Project not found", 404, errors.New(err.Error()))
		}
		return nil, *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}
	defer conn.Close()

	// ------------------------ Get the indexes from the database ------------------------

	indexes, err := GetProjectIndexes(ctx, conn)
	if err != nil {
		return nil, *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}

	config.App.InfoLog.Println("Indexes retrieved successfully for project:", projectOid)
	return indexes, *api.NewApiError("Indexes retrieved successfully", 200, nil)
}

func GetSpecificIndex(ctx context.Context, db *pgxpool.Pool, projectOid, indexOid string) (SpecificIndex, api.ApiError) {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return DefaultSpecificIndex, *api.NewApiError("Unauthorized", 401, errors.New("user is not authorized"))
	}

	// ------------------------ Get the project pool connection ------------------------
	conn, err := ProjectPoolConnection(ctx, db, UserID, projectOid)
	if err != nil {
		if err.Error() == "Project not found" || err.Error() == "connection pool not found" {
			return DefaultSpecificIndex, *api.NewApiError("Project not found", 404, errors.New(err.Error()))
		}
		return DefaultSpecificIndex, *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}

	defer conn.Close()

	// ------------------------ Get the index from the database ------------------------

	index := GetSpecificIndexFromDatabase(ctx, conn, indexOid)
	if index == (SpecificIndex{}) {
		return DefaultSpecificIndex, *api.NewApiError("Index not found", 404, errors.New("index with the given ID not found"))
	}

	// ------------------------ Close the connection ------------------------
	return index, *api.NewApiError("Index retrieved successfully", 200, nil)
}

func DeleteSpecificIndex(ctx context.Context, db *pgxpool.Pool, projectOid, indexOid string) api.ApiError {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return *api.NewApiError("Unauthorized", 401, errors.New("user is not authorized"))
	}

	// ------------------------ Get the project pool connection ------------------------
	conn, err := ProjectPoolConnection(ctx, db, UserID, projectOid)
	if err != nil {
		if err.Error() == "Project not found" || err.Error() == "connection pool not found" {
			return *api.NewApiError("Project not found", 404, errors.New(err.Error()))
		}
		return *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}

	defer conn.Close()

	// ------------------------ Delete the index from the database ------------------------

	IndexData := GetSpecificIndexFromDatabase(ctx, conn, indexOid)
	if IndexData == (SpecificIndex{}) {
		return *api.NewApiError("Index not found", 404, errors.New("index with the given ID not found"))
	}

	err = DeleteIndexFromDatabase(ctx, conn, IndexData.IndexName)

	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return *api.NewApiError("Index not found", 404, errors.New(err.Error()))
		}
		if strings.Contains(err.Error(), "cannot drop index") {
			return *api.NewApiError("Index cannot be dropped", 400, errors.New(err.Error()))
		}
		return *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}
	return *api.NewApiError("Index deleted successfully", 200, nil)
}

func UpdateSpecificIndex(ctx context.Context, db *pgxpool.Pool, projectOid, indexOid string, newName string) api.ApiError {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return *api.NewApiError("Unauthorized", 401, errors.New("user is not authorized"))
	}

	// ------------------------ Get the project pool connection ------------------------
	conn, err := ProjectPoolConnection(ctx, db, UserID, projectOid)
	if err != nil {
		if err.Error() == "Project not found" || err.Error() == "connection pool not found" {
			return *api.NewApiError("Project not found", 404, errors.New(err.Error()))
		}
		return *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}

	defer conn.Close()

	// ------------------------ Update the index in the database ------------------------
	newName = utils.ReplaceWhiteSpacesWithUnderscore(newName)

	// Get the current index name
	indexData := GetSpecificIndexFromDatabase(ctx, conn, indexOid)
	if indexData == (SpecificIndex{}) {
		return *api.NewApiError("Index not found", 404, errors.New("index with the given ID not found"))
	}

	if indexData.IndexName == newName {
		return *api.NewApiError("Index name is the same as the current name", 400, errors.New("index name is the same as the current name"))
	}

	err = UpdateIndexNameInDatabase(ctx, conn, indexData.IndexName, newName)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return *api.NewApiError("Index with the same name already exists", 400, errors.New(err.Error()))
		}
		if strings.Contains(err.Error(), "not found") {
			return *api.NewApiError("Index not found", 404, errors.New("index with the given ID not found"))
		}
		return *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}

	return *api.NewApiError("Index updated successfully", 200, nil)
}
