package tables

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row
	Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error)
}

func GetProjcetFeild(ctx context.Context, projectId int, fieldName string, db Querier) (interface{}, error) {
	query := fmt.Sprintf("SELECT %s FROM projects WHERE oid = $1", fieldName)
	var res interface{}
	err := db.QueryRow(ctx, query, projectId).Scan(&res)
	if err != nil {
		return nil, err
	}
	return res, nil
}