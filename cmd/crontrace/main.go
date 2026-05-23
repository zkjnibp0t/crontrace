package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/user/crontrace/internal/cli"
	"github.com/user/crontrace/internal/db"
	"github.com/user/crontrace/internal/runner"
)

const defaultDBPath = "/var/lib/crontrace/crontrace.db"

func main() {
	dbPath := flag.String("db", defaultDBPath, "path to SQLite database")
	listJob := flag.String("list", "", "list runs for a job name (or all if empty)")
	listAll := flag.Bool("list-all", false, "list all recorded job runs")
	flag.Parse()

	database, err := db.Open(*dbPath)
	if err != nil {
		log.Fatalf("crontrace: open database: %v", err)
	}
	defer database.Close()

	switch {
	case *listAll:
		if err := cli.ListRuns(database, ""); err != nil {
			log.Fatalf("crontrace: list runs: %v", err)
		}
		return

	case *listJob != "":
		if err := cli.ListRuns(database, *listJob); err != nil {
			log.Fatalf("crontrace: list runs: %v", err)
		}
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "usage: crontrace [flags] <command> [args...]")
		fmt.Fprintln(os.Stderr, "       crontrace --list-all")
		fmt.Fprintln(os.Stderr, "       crontrace --list <job-name>")
		flag.PrintDefaults()
		os.Exit(1)
	}

	jobName := args[0]
	cmdLine := strings.Join(args, " ")

	exitCode, err := runner.Run(database, jobName, cmdLine, args[0], args[1:]...)
	if err != nil {
		log.Fatalf("crontrace: run: %v", err)
	}
	os.Exit(exitCode)
}
