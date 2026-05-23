package db

import (
	"database/sql"
	"time"
)

// JobRun represents a single execution of a cron job.
type JobRun struct {
	ID        int64
	Name      string
	Command   string
	StartedAt time.Time
	EndedAt   *time.Time
	ExitCode  *int
	Output    string
}

// InsertJobRun records the start of a job execution and returns the new row ID.
func InsertJobRun(db *sql.DB, name, command string, startedAt time.Time) (int64, error) {
	res, err := db.Exec(
		`INSERT INTO job_runs (name, command, started_at) VALUES (?, ?, ?)`,
		name, command, startedAt.UTC().Format(time.RFC3339Nano),
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}

// FinishJobRun updates an existing job run with its end time, exit code, and output.
func FinishJobRun(db *sql.DB, id int64, endedAt time.Time, exitCode int, output string) error {
	_, err := db.Exec(
		`UPDATE job_runs SET ended_at = ?, exit_code = ?, output = ? WHERE id = ?`,
		endedAt.UTC().Format(time.RFC3339Nano), exitCode, output, id,
	)
	return err
}

// GetJobRun retrieves a single job run by ID.
func GetJobRun(db *sql.DB, id int64) (*JobRun, error) {
	row := db.QueryRow(`SELECT id, name, command, started_at, ended_at, exit_code, output FROM job_runs WHERE id = ?`, id)
	return scanJobRun(row)
}

// ListJobRuns returns all recorded job runs ordered by most recent first.
func ListJobRuns(db *sql.DB) ([]JobRun, error) {
	rows, err := db.Query(`SELECT id, name, command, started_at, ended_at, exit_code, output FROM job_runs ORDER BY started_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var runs []JobRun
	for rows.Next() {
		run, err := scanJobRun(rows)
		if err != nil {
			return nil, err
		}
		runs = append(runs, *run)
	}
	return runs, rows.Err()
}

type scanner interface {
	Scan(dest ...any) error
}

func scanJobRun(s scanner) (*JobRun, error) {
	var r JobRun
	var startedAtStr string
	var endedAtStr *string

	if err := s.Scan(&r.ID, &r.Name, &r.Command, &startedAtStr, &endedAtStr, &r.ExitCode, &r.Output); err != nil {
		return nil, err
	}

	parsed, err := time.Parse(time.RFC3339Nano, startedAtStr)
	if err != nil {
		return nil, err
	}
	r.StartedAt = parsed

	if endedAtStr != nil {
		t, err := time.Parse(time.RFC3339Nano, *endedAtStr)
		if err != nil {
			return nil, err
		}
		r.EndedAt = &t
	}

	return &r, nil
}
