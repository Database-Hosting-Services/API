package tables_test

import (
	"DBHS/config"
	"DBHS/tables"
	"DBHS/utils"
	"context"
	"fmt"
	"log"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

// Define project ID constants
const (
	existingProjectName = "ragnardbtest"
	testUserID          = 3
)

// Mock implementations for testing purposes
func mockCreateTable(ctx context.Context, projectOID string, table *tables.ClientTable, db *pgxpool.Pool, userDB *pgxpool.Pool, projectID int64) (string, error) {
	// Begin transaction in user database for table creation
	tx, err := userDB.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	// Create the table in the user db
	if err := tables.CreateTableIntoHostingServer(ctx, table, tx); err != nil {
		return "", err
	}

	// Generate table OID
	tableOID := utils.GenerateOID()

	// Insert table record into metadata DB
	tableRecord := tables.Table{
		Name:      table.TableName,
		ProjectID: projectID,
		OID:       tableOID,
	}

	var tableId int64
	if err := tables.InsertNewTable(ctx, &tableRecord, &tableId, db); err != nil {
		return "", err
	}

	if err := tx.Commit(ctx); err != nil {
		tables.DeleteTableRecord(ctx, tableId, db)
		return "", err
	}

	return tableRecord.OID, nil
}

func mockUpdateTable(ctx context.Context, projectOID string, tableOID string, updates *tables.TableUpdate,
	metadataDB *pgxpool.Pool, userDB *pgxpool.Pool) error {

	// Get table name from metadata DB
	var tableName string
	err := metadataDB.QueryRow(ctx, "SELECT name FROM \"Ptable\" WHERE oid = $1", tableOID).Scan(&tableName)
	if err != nil {
		return err
	}

	// Get table columns from the user database
	rows, err := userDB.Query(ctx, `
		SELECT column_name, data_type, is_nullable 
		FROM information_schema.columns 
		WHERE table_name = $1
	`, tableName)
	if err != nil {
		return err
	}
	defer rows.Close()

	// Build column map
	table := make(map[string]tables.DbColumn)
	for rows.Next() {
		var columnName, dataType, isNullable string
		if err := rows.Scan(&columnName, &dataType, &isNullable); err != nil {
			return err
		}

		isNull := isNullable == "YES"
		table[columnName] = tables.DbColumn{
			Name:       columnName,
			Type:       dataType,
			IsNullable: isNull,
		}
	}

	// Begin transaction in user database
	tx, err := userDB.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Execute update
	if err := tables.ExecuteUpdate(tableName, table, updates, tx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func mockDeleteTable(ctx context.Context, projectOID, tableOID string,
	metadataDB *pgxpool.Pool, userDB *pgxpool.Pool) error {

	// Get table name from metadata DB
	var tableName string
	err := metadataDB.QueryRow(ctx, "SELECT name FROM \"Ptable\" WHERE oid = $1", tableOID).Scan(&tableName)
	if err != nil {
		return err
	}

	// Begin transaction for user database operations
	userTx, err := userDB.Begin(ctx)
	if err != nil {
		return err
	}
	defer userTx.Rollback(ctx)

	// Delete table from user database
	if err := tables.DeleteTableFromHostingServer(ctx, tableName, userTx); err != nil {
		return err
	}

	// Commit user database transaction
	if err := userTx.Commit(ctx); err != nil {
		return err
	}

	// Now delete from metadata
	_, err = metadataDB.Exec(ctx, "DELETE FROM \"Ptable\" WHERE oid = $1", tableOID)
	return err
}

// ServiceTestSuite defines our test suite
type ServiceTestSuite struct {
	suite.Suite
	metadataDB    *pgxpool.Pool // For metadata (TEST_DATABASE_URL)
	userDB        *pgxpool.Pool // For user data (DATABASE_ADMIN_URL)
	ctx           context.Context
	origDB        *pgxpool.Pool
	origApp       *config.Application
	origConfigMgr *config.UserDbConfig
	projectID     int64
}

func (suite *ServiceTestSuite) SetupSuite() {
	// Load environment variables
	err := godotenv.Load("../.env")
	if err != nil {
		suite.T().Fatal("Error loading .env file:", err)
	}

	// Save original config references
	suite.origDB = config.DB
	suite.origApp = config.App
	suite.origConfigMgr = config.ConfigManager

	// Initialize context with JWT token information
	suite.ctx = context.WithValue(context.Background(), "user-id", testUserID)
	suite.ctx = context.WithValue(suite.ctx, "user-name", "Mohamed_Fathy")
	suite.ctx = context.WithValue(suite.ctx, "user-oid", "c536343f-2e63-42af-82db-ab6c4721106c")

	// Get database URLs
	metadataDBURL := os.Getenv("TEST_DATABASE_URL")
	if metadataDBURL == "" {
		suite.T().Fatal("TEST_DATABASE_URL environment variable is not set")
	}

	userDBURL := os.Getenv("DATABASE_ADMIN_URL")
	if userDBURL == "" {
		suite.T().Fatal("DATABASE_ADMIN_URL environment variable is not set")
	}

	// Create metadata DB connection with prepared statements disabled
	metaConfig, err := pgxpool.ParseConfig(metadataDBURL)
	if err != nil {
		suite.T().Fatal("Error parsing metadata database config:", err)
	}
	metaConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	// Create the metadata database connection
	suite.metadataDB, err = pgxpool.NewWithConfig(context.Background(), metaConfig)
	require.NoError(suite.T(), err, "Failed to connect to metadata database")

	// Set global DB for use in service functions to the metadata DB
	config.DB = suite.metadataDB

	// Initialize ConfigManager with the user database URL
	config.ConfigManager, err = config.NewUserDbConfig(userDBURL)
	require.NoError(suite.T(), err, "Failed to initialize ConfigManager")

	// Get the project ID for the test project
	err = suite.metadataDB.QueryRow(suite.ctx,
		"SELECT id FROM projects WHERE oid = $1", existingProjectOID).Scan(&suite.projectID)
	if err != nil {
		// Project doesn't exist, create it
		suite.projectID = suite.createTestProject()
	}

	// Initialize App with proper loggers
	config.App = &config.Application{
		ErrorLog: log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile),
		InfoLog:  log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime),
	}

	// Connect to the user's database
	userConfig, err := pgxpool.ParseConfig(userDBURL)
	if err != nil {
		suite.T().Fatal("Error parsing user database config:", err)
	}
	userConfig.ConnConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

	suite.userDB, err = pgxpool.NewWithConfig(context.Background(), userConfig)
	require.NoError(suite.T(), err, "Failed to connect to user database")

	suite.T().Logf("TEST_DATABASE_URL connected as metadataDB for service metadata")
	suite.T().Logf("DATABASE_ADMIN_URL connected as userDB for actual project data")
	config.ConfigManager, err = config.NewUserDbConfig(userDBURL)
	require.NoError(suite.T(), err, "Failed to initialize ConfigManager")
}

func (suite *ServiceTestSuite) createTestProject() int64 {
	var projectID int64
	err := suite.metadataDB.QueryRow(suite.ctx, `
		INSERT INTO projects (oid, owner_id, name, description, status, created_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, existingProjectOID, testUserID, existingProjectName, "Test project for service tests", "active", time.Now()).Scan(&projectID)

	require.NoError(suite.T(), err, "Failed to create test project")
	return projectID
}

func (suite *ServiceTestSuite) TearDownSuite() {
	// Restore original application settings
	config.DB = suite.origDB
	config.App = suite.origApp
	config.ConfigManager = suite.origConfigMgr

	// Close both database connections
	if suite.metadataDB != nil {
		suite.metadataDB.Close()
	}

	if suite.userDB != nil {
		suite.userDB.Close()
	}
}

func (suite *ServiceTestSuite) SetupTest() {
	// Clean up any existing test tables before each test
	_, err := suite.metadataDB.Exec(suite.ctx, `
		DELETE FROM "Ptable" WHERE project_id = $1 AND (
			name = 'test_table' OR 
			name = 'update_test_table' OR 
			name = 'delete_test_table' OR
			name = 'service_test_table'
		)
	`, suite.projectID)

	if err != nil {
		suite.T().Logf("Warning: Failed to clean up test tables: %v", err)
	}

	// Drop physical tables if they exist in user database
	tables := []string{"test_table", "update_test_table", "delete_test_table", "service_test_table"}
	for _, table := range tables {
		_, err := suite.userDB.Exec(suite.ctx, fmt.Sprintf("DROP TABLE IF EXISTS %s", table))
		if err != nil {
			suite.T().Logf("Warning: Failed to drop table %s: %v", table, err)
		}
	}
}

func (suite *ServiceTestSuite) TestGetAllTables() {
	// First insert a test table into metadata DB
	_, err := suite.metadataDB.Exec(suite.ctx, `
		INSERT INTO "Ptable" (oid, name, project_id, description)
		VALUES ('service-test-table', 'service_test_table', $1, 'Service test table')
		ON CONFLICT (oid) DO NOTHING;
	`, suite.projectID)
	require.NoError(suite.T(), err)

	// Test GetAllTables function, which should query the metadata DB
	allTables, err := tables.GetAllTables(suite.ctx, existingProjectOID, suite.metadataDB)
	require.NoError(suite.T(), err)
	assert.GreaterOrEqual(suite.T(), len(allTables), 1)

	// Verify our test table is in the response
	found := false
	for _, t := range allTables {
		if t.OID == "service-test-table" && t.Name == "service_test_table" {
			found = true
			break
		}
	}
	assert.True(suite.T(), found, "Test table not found in GetAllTables response")
}

func (suite *ServiceTestSuite) TestCreateTable() {
	// Create a client table
	table := &tables.ClientTable{
		TableName: "test_table",
		Columns: []tables.Column{
			{
				Name:         "id",
				Type:         "serial",
				IsPrimaryKey: func() *bool { b := true; return &b }(),
			},
			{
				Name:       "name",
				Type:       "varchar(100)",
				IsNullable: func() *bool { b := false; return &b }(),
			},
		},
	}

	// Use our mock function that doesn't rely on ExtractDb
	tableOID, err := mockCreateTable(suite.ctx, existingProjectOID, table, suite.metadataDB, suite.userDB, suite.projectID)

	// Assertions
	require.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), tableOID)

	// Verify the table was created in the metadata database
	var count int
	err = suite.metadataDB.QueryRow(suite.ctx, `
		SELECT COUNT(*) FROM "Ptable" 
		WHERE name = 'test_table' AND project_id = $1
	`, suite.projectID).Scan(&count)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count)

	// Also verify the actual table was created in the user database
	err = suite.userDB.QueryRow(suite.ctx, `
		SELECT COUNT(*) FROM information_schema.tables 
		WHERE table_name = 'test_table'
	`).Scan(&count)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, count, "Physical table was not created in user database")
}

func (suite *ServiceTestSuite) TestUpdateTable() {
	// First create a table to update
	table := &tables.ClientTable{
		TableName: "update_test_table",
		Columns: []tables.Column{
			{
				Name:         "id",
				Type:         "serial",
				IsPrimaryKey: func() *bool { b := true; return &b }(),
			},
			{
				Name:       "description",
				Type:       "text",
				IsNullable: func() *bool { b := true; return &b }(),
			},
		},
	}

	// Create the table using our
	tableOID, err := tables.CreateTable(suite.ctx, existingProjectOID, table, suite.metadataDB)
	require.NoError(suite.T(), err)
	// Give the system time to fully create the table
	time.Sleep(100 * time.Millisecond)

	// Create update definition
	updates := &tables.TableUpdate{
		Inserts: tables.ColumnCollection{
			Columns: []tables.Column{
				{
					Name:       "created_at",
					Type:       "timestamp",
					IsNullable: func() *bool { b := false; return &b }(),
				},
			},
		},
		Updates: []tables.UpdateColumn{
			{
				Name: "description",
				Update: tables.Column{
					Name: "title",
					Type: "varchar(200)",
				},
			},
		},
	}

	// Update the table
	err = tables.UpdateTable(suite.ctx, existingProjectOID, tableOID, updates, suite.metadataDB)
	require.NoError(suite.T(), err)

	// Assertions
	assert.NoError(suite.T(), err)

	// Verify the column was renamed and new column added in the user database
	var columnCount int
	err = suite.userDB.QueryRow(suite.ctx, `
		SELECT COUNT(*) FROM information_schema.columns 
		WHERE table_name = 'update_test_table' AND column_name = 'title'
	`).Scan(&columnCount)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, columnCount, "Column was not renamed in user database")

	err = suite.userDB.QueryRow(suite.ctx, `
		SELECT COUNT(*) FROM information_schema.columns 
		WHERE table_name = 'update_test_table' AND column_name = 'created_at'
	`).Scan(&columnCount)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 1, columnCount, "New column was not added in user database")
}

func (suite *ServiceTestSuite) TestDeleteTable() {
	// First create a table to delete
	table := &tables.ClientTable{
		TableName: "delete_test_table",
		Columns: []tables.Column{
			{
				Name:         "id",
				Type:         "serial",
				IsPrimaryKey: func() *bool { b := true; return &b }(),
			},
		},
	}

	// Create the table using our mock function
	tableOID, err := mockCreateTable(suite.ctx, existingProjectOID, table, suite.metadataDB, suite.userDB, suite.projectID)
	require.NoError(suite.T(), err)

	// Test using our mock delete function
	err = mockDeleteTable(suite.ctx, existingProjectOID, tableOID, suite.metadataDB, suite.userDB)

	// Assertions
	assert.NoError(suite.T(), err)

	// Verify the table was deleted from metadata
	var count int
	err = suite.metadataDB.QueryRow(suite.ctx,
		`SELECT COUNT(*) FROM "Ptable" WHERE oid = $1`, tableOID).Scan(&count)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count)

	// Verify the physical table was deleted from user database
	err = suite.userDB.QueryRow(suite.ctx, `
		SELECT COUNT(*) FROM information_schema.tables 
		WHERE table_name = 'delete_test_table'
	`).Scan(&count)
	require.NoError(suite.T(), err)
	assert.Equal(suite.T(), 0, count, "Physical table was not deleted from user database")
}

// Example of how to set URL variables for mux in HTTP handler tests
func (suite *ServiceTestSuite) TestURLVariables() {
	// Create a new HTTP request
	req := httptest.NewRequest("GET", "/api/projects/"+existingProjectOID+"/tables/some-table-id", nil)

	// Set the URL variables that mux would extract from the URL
	vars := map[string]string{
		"project-id": existingProjectOID,
		"table-id":   "some-table-id",
	}

	// Use mux.SetURLVars to add these variables to the request
	req = mux.SetURLVars(req, vars)

	// Now the request has the URL variables that mux.Vars(r) would extract
	// You can verify this works by extracting them again
	extractedVars := mux.Vars(req)
	suite.Equal(existingProjectOID, extractedVars["project-id"])
	suite.Equal("some-table-id", extractedVars["table-id"])
}

// Run the test suite
func TestServiceSuite(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping service tests in short mode")
	}
	suite.Run(t, new(ServiceTestSuite))
}
