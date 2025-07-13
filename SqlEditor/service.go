package sqleditor

import (
	"DBHS/indexes"
	api "DBHS/utils/apiError"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

func GetQueryResponse(ctx context.Context, db *pgxpool.Pool, projectOid, query string) (ResponseBody, api.ApiError) {
	owner_id, ok := ctx.Value("user-id").(int64)
	if !ok || owner_id == 0 {
		return ResponseBody{}, *api.NewApiError("Unauthorized", 401, errors.New("user is not authorized"))
	}

	// ------------------------ Get the project pool connection ------------------------
	conn, err := indexes.ProjectPoolConnection(ctx, db, owner_id, projectOid)
	if err != nil {
		if err.Error() == "Project not found" || err.Error() == "connection pool not found" {
			return ResponseBody{}, *api.NewApiError("Project not found", 404, errors.New(err.Error()))
		}
		return ResponseBody{}, *api.NewApiError("Internal server error", 500, errors.New(err.Error()))
	}
	defer conn.Close()

	// ------------------------ Fetch the query data ------------------------

	requestBody, apiErr := FetchQueryData(ctx, conn, query)
	if apiErr.Error() != nil {
		return ResponseBody{}, apiErr
	}

	return requestBody, api.ApiError{}
}
