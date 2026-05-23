package db_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/crontrace/internal/db"
)

func TestOpen_CreatesDatabase(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "crontrace.db")

	conn, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("Open() error = %v", err)
	}
	defer conn.Close()

	// Verify the file was created.
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Fatal("expected database file to exist")
	}

	// Verify the table exists by querying it.
	row := conn.QueryRow(`SELECT count(*) FROM job_runs`)
	var count int
	if err := row.Scan(&count); err != nil {
		t.Fatalf("querying job_runs: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected 0 rows, got %d", count)
	}
}

func TestOpen_Idempotent(t *testing.T) {
	dir := t.TempDir()
	dbPath := filepath.Join(dir, "crontrace.db")

	for i := 0; i < 3; i++ {
		conn, err := db.Open(dbPath)
		if err != nil {
			t.Fatalf("Open() attempt %d error = %v", i+1, err)
		}
		conn.Close()
	}
}

func TestOpen_InvalidPath(t *testing.T) {
	_, err := db.Open("/nonexistent/directory/crontrace.db")
	if err == nil {
		t.Fatal("expected error for invalid path, got nil")
	}
}
