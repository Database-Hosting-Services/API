package indexes

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

func GetProjectIndexes(ctx context.Context, conn *pgxpool.Pool) ([]RetrievedIndex, error) {
	rows, err := conn.Query(ctx, SELECT_ALL_INDEXES)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var indexes []RetrievedIndex
	for rows.Next() {
		var index RetrievedIndex
		// Scan three columns: (index_name, index_oid, index_type)
		if err := rows.Scan(&index.IndexName, &index.IndexOid, &index.IndexType); err != nil {
			return nil, err
		}
		indexes = append(indexes, index)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return indexes, nil
}

func GetSpecificIndexFromDatabase(ctx context.Context, conn *pgxpool.Pool, indexOid string) SpecificIndex {
	row := conn.QueryRow(ctx, SELECT_SPECIFIC_INDEX, indexOid)
	if row == nil {
		return DefaultSpecificIndex
	}
	var index SpecificIndex
	if err := row.Scan(&index.IndexName, &index.IndexType); err != nil {
		return DefaultSpecificIndex
	}
	return index
}
