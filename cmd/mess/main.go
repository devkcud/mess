package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/devkcud/mess/internal/core"
	"github.com/devkcud/mess/pkg/messlog"
)

func main() {
	scriptTimeStart := time.Now()

	dir, err := os.Getwd()
	if err != nil {
		dir = "."
	}

	cli := core.NewCLI()

	base := cli.String("base", dir, "base working directory")
	dryRun := cli.Bool("dry-run", false, "simulate file/folder creation without writing anything on disk")
	toggleSummary := cli.Bool("summary", false, "print a summary after execution")
	loglevel := cli.Int("loglevel", int(messlog.LogLevelError), "logging output (0 = error | 1 = warn | 2 = info | 3 = debug | 4 = trace)")
	help := cli.Bool("help", false, "help menu")

	tokens, err := cli.Parse()
	if err != nil {
		log.Fatalf("failed to parse flags: %v", err)
	}

	if *help {
		cli.HelpExit(false)
	}

	if len(tokens) == 0 {
		cli.HelpExit(true)
	}

	logger := messlog.NewLogger(messlog.LogLevel(*loglevel))

	var summary *core.Summary
	if *toggleSummary {
		summary = &core.Summary{}
	}

	tokenIterStart := time.Now()
	builder := core.NewBuilder(*base, logger, *dryRun, summary)
	for i, token := range tokens {
		iterStart := time.Now()

		if err := builder.ProcessToken(token); err != nil {
			logger.Error("Error processing %q: %v", token, err)
		}

		logger.Trace("Loop %d/%d for token %q in %s", i+1, len(tokens), token, time.Since(iterStart))
	}

	logger.Trace("Ran all %d tokens in %s", len(tokens), time.Since(tokenIterStart))

	if *toggleSummary {
		fmt.Println("+------------------+--------+")
		fmt.Printf("| %-16s | %6s |\n", "OPERATION", "COUNT")
		fmt.Println("+------------------+--------+")
		fmt.Printf("| %-16s | %6d |\n", "Folders Created", summary.DirectoriesCreated)
		fmt.Printf("| %-16s | %6d |\n", "Files Created", summary.FilesCreated)
		fmt.Printf("| %-16s | %6d |\n", "Token Successes", summary.Successes)
		fmt.Printf("| %-16s | %6d |\n", "Token Failures", summary.Failures)
		fmt.Println("+------------------+--------+")
	}

	if *dryRun {
		builder.PrintDryRunTree()
	}

	logger.Trace("Finished in %s", time.Since(scriptTimeStart))
}
