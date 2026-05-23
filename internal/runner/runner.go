package runner

import (
	"database/sql"
	"os/exec"
	"time"

	"github.com/user/crontrace/internal/db"
)

// Result holds the outcome of a job execution.
type Result struct {
	JobRunID int64
	Command  string
	Args     []string
	ExitCode int
	Duration time.Duration
	Error    error
}

// Run executes the given command with args, records the run in the database,
// and returns a Result with execution details.
func Run(database *sql.DB, command string, args []string) (*Result, error) {
	full := command
	if len(args) > 0 {
		for _, a := range args {
			full += " " + a
		}
	}

	jobRunID, err := db.InsertJobRun(database, full)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(command, args...)
	start := time.Now()
	runErr := cmd.Run()
	duration := time.Since(start)

	exitCode := 0
	if runErr != nil {
		if exitErr, ok := runErr.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			exitCode = -1
		}
	}

	if err := db.FinishJobRun(database, jobRunID, exitCode, duration); err != nil {
		return nil, err
	}

	return &Result{
		JobRunID: jobRunID,
		Command:  command,
		Args:     args,
		ExitCode: exitCode,
		Duration: duration,
		Error:    runErr,
	}, nil
}
