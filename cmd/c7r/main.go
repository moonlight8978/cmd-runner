package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/moonlight8978/cmd-runner/pkg/runner"
)

type ParsedArgs struct {
	ConfigPath string
	DryRun     bool
	Command    string
	CmdArgs    []string
}

func parseArgs() (*ParsedArgs, error) {
	args := os.Args[1:]
	if len(args) == 0 {
		return nil, fmt.Errorf("no command specified")
	}

	// First argument must be the command
	command := args[0]
	if strings.HasPrefix(command, "-") {
		return nil, fmt.Errorf("command must be the first argument")
	}

	var configPath string
	var dryRun bool

	// Parse remaining args
	i := 1
	for i < len(args) {
		switch args[i] {
		case "--config":
			if i+1 >= len(args) {
				return nil, fmt.Errorf("--config requires a value")
			}
			configPath = args[i+1]
			i += 2
		case "--dry-run":
			dryRun = true
			i++
		default:
			i++
		}
	}

	if configPath == "" {
		return nil, fmt.Errorf("--config flag is required")
	}

	return &ParsedArgs{
		ConfigPath: configPath,
		DryRun:     dryRun,
		Command:    command,
	}, nil
}

func main() {
	parsed, err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Usage: %s COMMAND [args...] --config <config.yaml> [--dry-run]\n", os.Args[0])
		os.Exit(1)
	}

	// Parse config
	configItems, err := runner.ParseConfigWithOptions(parsed.ConfigPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Build arguments from config
	commandRunner := runner.NewCommandRunner()
	configArgs := commandRunner.BuildArgs(configItems)
	finalArgs := append(parsed.CmdArgs, configArgs...)

	if parsed.DryRun {
		fmt.Printf("Command: %s %s\n", parsed.Command, strings.Join(finalArgs, " "))
		return
	}

	// Execute the command
	if err := commandRunner.RunCommand(parsed.Command, finalArgs); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		// Try to get exit code from the error
		if exitErr, ok := err.(*exec.ExitError); ok {
			os.Exit(exitErr.ExitCode())
		} else {
			// If no exit code available, use 1
			os.Exit(1)
		}
	}
}
