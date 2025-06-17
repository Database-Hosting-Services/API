package analytics

import (
	"DBHS/indexes"
	api "DBHS/utils/apiError"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetDatabaseStorage(ctx context.Context, db *pgxpool.Pool, projectOid string) (Storage, api.ApiError) {
	// Get user ID from context
	UserID, ok := ctx.Value("user-id").(int)
	if !ok || UserID == 0 {
		return Storage{}, *api.NewApiError("Unauthorized", 401, errors.New("user is not authorized"))
	}

	// ------------------------ Get the project pool connection ------------------------
	conn, err := indexes.ProjectPoolConnection(ctx, db, UserID, projectOid)
	if err != nil {
		if err.Error() == "Project not found" || err.Error() == "connection pool not found" {
			return Storage{}, *api.NewApiError("Project not found", 404, errors.New(err.Error()))
		}
		return Storage{}, *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}
	defer conn.Close()

	var storage Storage
	if err := conn.QueryRow(ctx, GET_CURRENT_STORAGE).Scan(&storage.ManagementStorage, &storage.ActualData); err != nil {
		return Storage{}, *api.NewApiError("Internal server error", 500, errors.New("failed to retrieve storage information: "+err.Error()))
	}

	return storage, api.ApiError{} // Return empty ApiError to indicate success
}
