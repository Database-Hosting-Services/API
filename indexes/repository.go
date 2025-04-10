package indexes

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5"
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

func DeleteIndexFromDatabase(ctx context.Context, conn *pgxpool.Pool, indexName string) error {
	DELETE_INDEX := GenerateDeleteIndexQuery(indexName)
	_, err := conn.Exec(ctx, DELETE_INDEX)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("index not found")
		}
		return err
	}
	return nil
}

func UpdateIndexNameInDatabase(ctx context.Context, conn *pgxpool.Pool, oldName string, newName string) error {
	UPDATE_INDEX := GenerateRenameIndexQuery(oldName, newName)
	_, err := conn.Exec(ctx, UPDATE_INDEX)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("index not found")
		}
		return err
	}
	return nil
}
