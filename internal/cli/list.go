package cli

import (
	"fmt"
	"os"
	"text/tabwriter"
	"time"

	"github.com/user/crontrace/internal/db"
)

// ListRuns prints all recorded job runs for the given job name.
// If jobName is empty, all runs across all jobs are listed.
func ListRuns(database *db.DB, jobName string) error {
	runs, err := db.ListJobRuns(database, jobName)
	if err != nil {
		return fmt.Errorf("listing job runs: %w", err)
	}

	if len(runs) == 0 {
		if jobName != "" {
			fmt.Printf("No runs found for job %q\n", jobName)
		} else {
			fmt.Println("No runs found.")
		}
		return nil
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "ID\tJOB\tSTARTED\tDURATION\tEXIT CODE\tSTATUS")
	fmt.Fprintln(w, "--\t---\t-------\t--------\t---------\t------")

	for _, r := range runs {
		duration := "-"
		exitCode := "-"
		status := "running"

		if r.FinishedAt != nil {
			duration = r.FinishedAt.Sub(r.StartedAt).Round(time.Millisecond).String()
			status = "done"
		}
		if r.ExitCode != nil {
			exitCode = fmt.Sprintf("%d", *r.ExitCode)
			if *r.ExitCode != 0 {
				status = "failed"
			}
		}

		fmt.Fprintf(w, "%d\t%s\t%s\t%s\t%s\t%s\n",
			r.ID,
			r.JobName,
			r.StartedAt.Format(time.RFC3339),
			duration,
			exitCode,
			status,
		)
	}

	return w.Flush()
}
