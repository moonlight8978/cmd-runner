package runner

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

// CommandRunner handles the execution of commands with YAML configuration
type CommandRunner struct{}

// NewCommandRunner creates a new command runner instance
func NewCommandRunner() *CommandRunner {
	return &CommandRunner{}
}

// BuildArgs converts the YAML configuration to command-line arguments
func (cr *CommandRunner) BuildArgs(config []ConfigItem) []string {
	var args []string
	for _, item := range config {
		if item.Name == "" {
			if str, ok := item.Value.(string); ok {
				// Expand environment variables in positional arguments
				args = append(args, expandEnvVars(str))
			}
			continue
		}

		joiner := *item.Joiner
		switch v := item.Value.(type) {
		case string:
			args = append(args, cr.formatFlag(item.Name, v, joiner)...)
		case []interface{}:
			for _, val := range v {
				if str, ok := val.(string); ok {
					args = append(args, cr.formatFlag(item.Name, str, joiner)...)
				}
			}
		case []string:
			for _, val := range v {
				args = append(args, cr.formatFlag(item.Name, val, joiner)...)
			}
		default:
			args = append(args, cr.formatFlag(item.Name, fmt.Sprintf("%v", v), joiner)...)
		}
	}
	return args
}

// formatFlag formats a flag name and value according to the configuration
func (cr *CommandRunner) formatFlag(name, value string, joiner string) []string {
	// Expand environment variables in the value
	value = expandEnvVars(value)

	if strings.Contains(value, " ") {
		value = fmt.Sprintf(`"%s"`, value)
	}

	if joiner == " " {
		return []string{name, value}
	}

	return []string{fmt.Sprintf("%s%s%s", name, joiner, value)}
}

// expandEnvVars replaces $VAR or ${VAR} with environment variable values
func expandEnvVars(s string) string {
	// Replace ${VAR} syntax
	s = os.Expand(s, func(key string) string {
		return os.Getenv(key)
	})

	// Replace $VAR syntax (simple case)
	re := regexp.MustCompile(`\$([A-Za-z_][A-Za-z0-9_]*)`)
	s = re.ReplaceAllStringFunc(s, func(match string) string {
		varName := match[1:] // Remove the $ prefix
		return os.Getenv(varName)
	})

	return s
}

// RunCommand executes the command with the given arguments
func (cr *CommandRunner) RunCommand(command string, args []string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
