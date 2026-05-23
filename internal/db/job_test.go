package db

import (
	"testing"
	"time"
)

func TestInsertAndFinishJobRun(t *testing.T) {
	db, err := Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()

	start := time.Now().Truncate(time.Second)
	id, err := InsertJobRun(db, "backup", "tar -czf /tmp/backup.tar.gz /data", start)
	if err != nil {
		t.Fatalf("InsertJobRun: %v", err)
	}
	if id <= 0 {
		t.Fatalf("expected positive id, got %d", id)
	}

	end := start.Add(3 * time.Second)
	if err := FinishJobRun(db, id, end, 0, "done"); err != nil {
		t.Fatalf("FinishJobRun: %v", err)
	}

	run, err := GetJobRun(db, id)
	if err != nil {
		t.Fatalf("GetJobRun: %v", err)
	}

	if run.Name != "backup" {
		t.Errorf("Name: got %q, want %q", run.Name, "backup")
	}
	if run.ExitCode == nil || *run.ExitCode != 0 {
		t.Errorf("ExitCode: got %v, want 0", run.ExitCode)
	}
	if run.EndedAt == nil {
		t.Error("EndedAt should not be nil after FinishJobRun")
	}
	if run.Output != "done" {
		t.Errorf("Output: got %q, want %q", run.Output, "done")
	}
}

func TestListJobRuns(t *testing.T) {
	db, err := Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()

	names := []string{"job-a", "job-b", "job-c"}
	for _, name := range names {
		_, err := InsertJobRun(db, name, "echo "+name, time.Now())
		if err != nil {
			t.Fatalf("InsertJobRun(%s): %v", name, err)
		}
	}

	runs, err := ListJobRuns(db)
	if err != nil {
		t.Fatalf("ListJobRuns: %v", err)
	}
	if len(runs) != len(names) {
		t.Fatalf("expected %d runs, got %d", len(names), len(runs))
	}
}

func TestGetJobRun_NotFound(t *testing.T) {
	db, err := Open(t.TempDir() + "/test.db")
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	defer db.Close()

	_, err = GetJobRun(db, 9999)
	if err == nil {
		t.Fatal("expected error for missing row, got nil")
	}
}
