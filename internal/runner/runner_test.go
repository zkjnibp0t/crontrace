package runner_test

import (
	"database/sql"
	"testing"

	"github.com/user/crontrace/internal/db"
	"github.com/user/crontrace/internal/runner"
)

func openTestDB(t *testing.T) *sql.DB {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { database.Close() })
	return database
}

func TestRun_Success(t *testing.T) {
	database := openTestDB(t)

	res, err := runner.Run(database, "true", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ExitCode != 0 {
		t.Errorf("expected exit code 0, got %d", res.ExitCode)
	}
	if res.Duration < 0 {
		t.Errorf("expected non-negative duration")
	}
	if res.JobRunID <= 0 {
		t.Errorf("expected positive job run ID, got %d", res.JobRunID)
	}
}

func TestRun_Failure(t *testing.T) {
	database := openTestDB(t)

	res, err := runner.Run(database, "false", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.ExitCode == 0 {
		t.Errorf("expected non-zero exit code")
	}
	if res.Error == nil {
		t.Errorf("expected non-nil Error for failed command")
	}
}

func TestRun_RecordedInDB(t *testing.T) {
	database := openTestDB(t)

	res, err := runner.Run(database, "echo", []string{"hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	jobRun, err := db.GetJobRun(database, res.JobRunID)
	if err != nil {
		t.Fatalf("get job run: %v", err)
	}
	if jobRun.ExitCode == nil || *jobRun.ExitCode != 0 {
		t.Errorf("expected exit code 0 in db")
	}
	if jobRun.FinishedAt == nil {
		t.Errorf("expected finished_at to be set")
	}
}
