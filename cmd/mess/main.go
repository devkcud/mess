package main

import (
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

	base := cli.StringP("base", "b", dir, "base working directory")
	dryRun := cli.BoolP("dry", "d", false, "simulate file/directory creation without writing anything on disk")
	echo := cli.BoolP("echo", "e", false, "print shell commands instead of creating anything")
	loglevel := cli.Int("loglevel", int(messlog.LogLevelError), "logging output (0 = error | 1 = warn | 2 = info | 3 = debug | 4 = trace)")
	help := cli.BoolP("help", "h", false, "help menu")

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

	tokenIterStart := time.Now()
	builder := core.NewBuilder(*base, logger, *dryRun, *echo)
	for i, token := range tokens {
		iterStart := time.Now()

		if err := builder.ProcessToken(token); err != nil {
			logger.Error("Error processing %q: %v", token, err)
		}

		logger.Trace("Loop %d/%d for token %q in %s", i+1, len(tokens), token, time.Since(iterStart))
	}

	logger.Trace("Ran all %d tokens in %s", len(tokens), time.Since(tokenIterStart))

	if *dryRun != false || *echo != false {
		logger.Info("Skipping file builds. Dry Run or Echo detected")

		if *dryRun {
			logger.Debug("Printing Dry Run tree")
			builder.PrintDryRunTree()
		}

		if *echo {
			logger.Debug("Printing Echo tree")
			builder.PrintEchoFiles()
		}
	} else {
		logger.Info("Building directories and files")

		if err := builder.BuildFiles(); err != nil {
			logger.Error("Couldn't write dir/file: %v", err)
		}
	}

	logger.Trace("Finished in %s", time.Since(scriptTimeStart))
}
