package database

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3" // SQLite driver
	"log"
	"os"
	"path/filepath"
	"time"
	"strings"
)

const dbFileName = "finetuning.db"

type InteractionLog struct {
	ID                     string
	Timestamp              time.Time
	Endpoint               string
	RequestPayload         string
	ResponsePayload        string
	AppName                string
	AppPath                string
	TestResultsJSON        string
	AnalysisResultsJSON    string
	FeedbackJSON           string
	Status                 string
	ProcessedForFinetuning bool
}

type DB struct {
	*sql.DB
}

func NewDB(dataDir string) (*DB, error) {
	dbPath := filepath.Join(dataDir, dbFileName)
	// Ensure the directory exists
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err = createTables(db); err != nil {
		return nil, fmt.Errorf("failed to create tables: %w", err)
	}

	log.Printf("Database initialized at %s", dbPath)
	return &DB{db}, nil
}

func createTables(db *sql.DB) error {
	sqlStmt := `
	CREATE TABLE IF NOT EXISTS interactions_log (
		id TEXT PRIMARY KEY,
		timestamp TEXT NOT NULL,
		endpoint TEXT NOT NULL,
		request_payload TEXT,
		response_payload TEXT,
		app_name TEXT,
		app_path TEXT,
		test_results_json TEXT,
		analysis_results_json TEXT,
		feedback_json TEXT,
		status TEXT NOT NULL,
		processed_for_finetuning INTEGER DEFAULT 0
	);
	CREATE INDEX IF NOT EXISTS idx_timestamp ON interactions_log (timestamp);
	CREATE INDEX IF NOT EXISTS idx_endpoint ON interactions_log (endpoint);
	CREATE INDEX IF NOT EXISTS idx_processed ON interactions_log (processed_for_finetuning);
	`
	_, err := db.Exec(sqlStmt)
	return err
}

func (d *DB) InsertInteractionLog(logEntry InteractionLog) error {
	stmt, err := d.Prepare(`
	INSERT INTO interactions_log (
		id, timestamp, endpoint, request_payload, response_payload, app_name, app_path,
		test_results_json, analysis_results_json, feedback_json, status, processed_for_finetuning
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return fmt.Errorf("failed to prepare insert statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		logEntry.ID,
		logEntry.Timestamp.Format(time.RFC3339),
		logEntry.Endpoint,
		logEntry.RequestPayload,
		logEntry.ResponsePayload,
		logEntry.AppName,
		logEntry.AppPath,
		logEntry.TestResultsJSON,
		logEntry.AnalysisResultsJSON,
		logEntry.FeedbackJSON,
		logEntry.Status,
		logEntry.ProcessedForFinetuning,
	)
	return err
}

func (d *DB) GetUnprocessedLogs() ([]InteractionLog, error) {
	rows, err := d.Query(`
	SELECT id, timestamp, endpoint, request_payload, response_payload, app_name, app_path,
		test_results_json, analysis_results_json, feedback_json, status, processed_for_finetuning
	FROM interactions_log
	WHERE processed_for_finetuning = 0
	ORDER BY timestamp ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query unprocessed logs: %w", err)
	}
	defer rows.Close()

	var logs []InteractionLog
	for rows.Next() {
		var logEntry InteractionLog
		var timestampStr string
		var processedInt int
		if err := rows.Scan(
			&logEntry.ID, &timestampStr, &logEntry.Endpoint, &logEntry.RequestPayload,
			&logEntry.ResponsePayload, &logEntry.AppName, &logEntry.AppPath,
			&logEntry.TestResultsJSON, &logEntry.AnalysisResultsJSON, &logEntry.FeedbackJSON,
			&logEntry.Status, &processedInt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		logEntry.Timestamp, err = time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			return nil, fmt.Errorf("failed to parse timestamp: %w", err)
		}
		logEntry.ProcessedForFinetuning = (processedInt == 1)
		logs = append(logs, logEntry)
	}

	return logs, nil
}

func (d *DB) MarkLogsAsProcessed(ids []string) error {
	if len(ids) == 0 {
		return nil
	}
	// Create a string of \'?\' placeholders for the IN clause
	placeholders := make([]string, len(ids))
	for i := range ids {
		placeholders[i] = "?"
	}
	query := fmt.Sprintf(`
	UPDATE interactions_log
	SET processed_for_finetuning = 1
	WHERE id IN (%s)
	`, strings.Join(placeholders, ","))

	stmt, err := d.Prepare(query)
	if err != nil {
		return fmt.Errorf("failed to prepare update statement: %w", err)
	}
	defer stmt.Close()

	args := make([]interface{}, len(ids))
	for i, id := range ids {
		args[i] = id
	}

	_, err = stmt.Exec(args...)
	return err
}


