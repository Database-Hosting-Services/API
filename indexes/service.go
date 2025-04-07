package indexes

import (
	"DBHS/config"
	"DBHS/projects"
	"context"
	"errors"
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

	// ------------------------ Create the index in the database ------------------------

	query := GenerateIndexQuery(indexData)
	if _, err = conn.Exec(ctx, query); err != nil {
		return err
	}

	// ------------------------ Close the connection ------------------------
	conn.Close()
	return nil
}
