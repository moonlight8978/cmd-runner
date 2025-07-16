package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/moonlight8978/cmd-runner/pkg/runner"
)

// parseArgs scans os.Args for special flags and the command, supporting flags before or after the command.
type ParsedArgs struct {
	ConfigPath string
	DryRun     bool
	Command    string
	CmdArgs    []string
}

func parseArgs() (*ParsedArgs, error) {
	args := os.Args[1:]
	var (
		configPath string
		dryRun = false
		command string
		cmdArgs []string
	)

	// Track which args are our special flags
	skip := make([]bool, len(args))

	// First pass: extract our flags
	for i := 0; i < len(args); i++ {
		if args[i] == "--config" && i+1 < len(args) {
			configPath = args[i+1]
			skip[i], skip[i+1] = true, true
			i++
		} else if strings.HasPrefix(args[i], "--config=") {
			configPath = strings.TrimPrefix(args[i], "--config=")
			skip[i] = true
		} else if args[i] == "--dry-run" {
			dryRun = true
			skip[i] = true
		}
	}

	// Second pass: find the command (first non-flag, non-skipped arg)
	for i, arg := range args {
		if skip[i] {
			continue
		}
		if command == "" && !strings.HasPrefix(arg, "-") {
			command = arg
			skip[i] = true
			continue
		}
	}

	// Third pass: collect remaining args as command args
	for i, arg := range args {
		if skip[i] {
			continue
		}
		cmdArgs = append(cmdArgs, arg)
	}

	if configPath == "" {
		return nil, fmt.Errorf("--config flag is required")
	}
	if command == "" {
		return nil, fmt.Errorf("no command specified")
	}

	return &ParsedArgs{
		ConfigPath: configPath,
		DryRun:     dryRun,
		Command:    command,
		CmdArgs:    cmdArgs,
	}, nil
}

func main() {
	parsed, err := parseArgs()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "Usage: %s <command> [args...] --config <config.yaml> [--dry-run]\n", os.Args[0])
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
