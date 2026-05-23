package cli_test

import (
	"io"
	"os"
	"testing"
	"time"

	"github.com/user/crontrace/internal/cli"
	"github.com/user/crontrace/internal/db"
)

func openTestDB(t *testing.T) *db.DB {
	t.Helper()
	database, err := db.Open(":memory:")
	if err != nil {
		t.Fatalf("open db: %v", err)
	}
	t.Cleanup(func() { database.Close() })
	return database
}

func TestListRuns_Empty(t *testing.T) {
	database := openTestDB(t)

	// Redirect stdout to discard output
	old := os.Stdout
	os.Stdout, _ = os.Open(os.DevNull)
	defer func() { os.Stdout = old }()

	if err := cli.ListRuns(database, ""); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestListRuns_WithRuns(t *testing.T) {
	database := openTestDB(t)

	id, err := db.InsertJobRun(database, "backup", "rsync -av /src /dst")
	if err != nil {
		t.Fatalf("insert: %v", err)
	}
	exitCode := 0
	if err := db.FinishJobRun(database, id, exitCode, ""); err != nil {
		t.Fatalf("finish: %v", err)
	}

	// Capture stdout
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	if err := cli.ListRuns(database, "backup"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)

	if len(out) == 0 {
		t.Error("expected output, got none")
	}
}

func TestListRuns_FilterByJob(t *testing.T) {
	database := openTestDB(t)

	for _, name := range []string{"jobA", "jobB", "jobA"} {
		_, err := db.InsertJobRun(database, name, "echo "+name)
		if err != nil {
			t.Fatalf("insert %s: %v", name, err)
		}
	}
	_ = time.Now() // ensure time package used

	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w

	if err := cli.ListRuns(database, "jobA"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	w.Close()
	os.Stdout = old
	out, _ := io.ReadAll(r)

	if len(out) == 0 {
		t.Error("expected output for jobA, got none")
	}
}
