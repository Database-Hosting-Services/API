package tables_test

import (
	"DBHS/config"
	"DBHS/tables"
	"DBHS/utils"
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const (
	existingProjectOID = "066e68e9-95c1-4d9f-a3d6-7fedcb0e46a7"
)

// Repository test suite
type RepositoryTestSuite struct {
	suite.Suite
	db      *pgxpool.Pool
	ctx     context.Context
	origDB  *pgxpool.Pool
	origApp *config.Application
}

func (suite *RepositoryTestSuite) SetupSuite() {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		suite.T().Fatal("Error loading .env file:", err)
	}

	// Save original config references
	suite.origDB = config.DB
	suite.origApp = config.App

	// Initialize context with JWT token information
	suite.ctx = context.WithValue(context.Background(), "user-id", 3)
	suite.ctx = context.WithValue(suite.ctx, "user-name", "Mohamed_Fathy")
	suite.ctx = context.WithValue(suite.ctx, "user-oid", "c536343f-2e63-42af-82db-ab6c4721106c")

	// Connect to test database
	dbURL := os.Getenv("TEST_DATABASE_URL")
	if dbURL == "" {
		dbURL = os.Getenv("DATABASE_ADMIN_URL") // Fallback to admin URL if test URL not set
		if dbURL == "" {
			suite.T().Fatal("Neither TEST_DATABASE_URL nor DATABASE_ADMIN_URL environment variable is set")
		}
	}

	// Connect to database
	suite.db, err = pgxpool.New(context.Background(), dbURL)
	require.NoError(suite.T(), err, "Failed to connect to test database")

	// Initialize App with proper loggers
	config.App = &config.Application{
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}

	// Set global DB for repository functions
	config.DB = suite.db

}

func (suite *RepositoryTestSuite) TearDownSuite() {
	// Close database connection
	if suite.db != nil {
		suite.db.Close()
	}

	// Restore original config
	config.DB = suite.origDB
	config.App = suite.origApp
}

func (suite *RepositoryTestSuite) SetupTest() {
	// Create test projects using the existing schema
	_, err := suite.db.Exec(context.Background(), `
		INSERT INTO "Ptable" (oid, name, project_id, description)
		VALUES (
			'test-table-1', 
			'test_table_1', 
			(SELECT id FROM projects WHERE oid = $1), 
			'Test table'
		)
		ON CONFLICT (oid) DO NOTHING;
	`, existingProjectOID)

	require.NoError(suite.T(), err, "Failed to set up test data")

}

func (suite *RepositoryTestSuite) TearDownTest() {
	// Clean up test data
	_, err := suite.db.Exec(context.Background(), `
		DELETE FROM "Ptable";
	`)
	require.NoError(suite.T(), err, "Failed to clean up test data")
}

func (suite *RepositoryTestSuite) TestGetProjectNameID() {
	// Test getting project name and ID
	name, id, err := utils.GetProjectNameID(context.Background(), existingProjectOID, suite.db)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "ragnardbtest", name)
	assert.NotNil(suite.T(), id)
}


func (suite *RepositoryTestSuite) TestGetTableName() {
	// Test getting table name by OID
	tableName, err := tables.GetTableName(context.Background(), "test-table-1", suite.db)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test_table_1", tableName)
}

func (suite *RepositoryTestSuite) TestInsertAndDeleteTableRecord() {
	// Test inserting a new table record
	table := &tables.Table{
		OID:         "test-table-insert",
		Name:        "test_insert_table",
		Description: "Table for testing insert",
		ProjectID:   129,
	}

	var tableID int64
	err := tables.InsertNewTable(context.Background(), table, &tableID, suite.db)
	require.NoError(suite.T(), err)
	assert.Greater(suite.T(), tableID, 0)

	// Verify the table was inserted
	var count int
	err = suite.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM "Ptable" WHERE oid = $1`, "test-table-insert").Scan(&count)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	// Test deleting the table record
	err = tables.DeleteTableRecord(context.Background(), tableID, suite.db)
	require.NoError(suite.T(), err)

	// Verify the table was deleted
	err = suite.db.QueryRow(context.Background(),
		`SELECT COUNT(*) FROM "Ptable" WHERE id = $1`, tableID).Scan(&count)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)
}

func (suite *RepositoryTestSuite) TestCheckOwnershipQuery() {
	// Test checking ownership - positive case
	isOwner, err := tables.CheckOwnershipQuery(context.Background(), existingProjectOID, 3, suite.db)
	require.NoError(suite.T(), err)
	assert.True(suite.T(), isOwner)

	// Test checking ownership - negative case
	isOwner, err = tables.CheckOwnershipQuery(context.Background(), existingProjectOID, 999, suite.db)
	require.NoError(suite.T(), err)
	assert.False(suite.T(), isOwner)
}

func (suite *RepositoryTestSuite) TestPrepareQuery() {
	// Test query preparation with different parameters
	tests := []struct {
		name       string
		tableName  string
		parameters map[string][]string
		wantQuery  string
		wantErr    bool
	}{
		{
			name:      "Basic query",
			tableName: "test_table",
			parameters: map[string][]string{
				"limit": {"10"},
				"page":  {"0"},
			},
			wantQuery: "SELECT * FROM test_table LIMIT 10 OFFSET 0;",
			wantErr:   false,
		},
		{
			name:      "Query with ordering",
			tableName: "test_table",
			parameters: map[string][]string{
				"limit": {"10"},
				"page":  {"0"},
				"order": {"id:asc"},
			},
			wantQuery: "SELECT * FROM test_table ORDER BY id ASC LIMIT 10 OFFSET 0;",
			wantErr:   false,
		},
		{
			name:      "Query with filter",
			tableName: "test_table",
			parameters: map[string][]string{
				"limit":  {"10"},
				"page":   {"0"},
				"filter": {"id:eq:5"},
			},
			wantQuery: "SELECT * FROM test_table WHERE id = 5 LIMIT 10 OFFSET 0;",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		suite.T().Run(tt.name, func(t *testing.T) {
			query, err := tables.PrepareQuery(tt.tableName, tt.parameters)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.wantQuery, query)
		})
	}
}

func TestRepositorySuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping repository tests in short mode")
	}
	suite.Run(t, new(RepositoryTestSuite))
}
