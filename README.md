# Cmd/Command Runner (c7r)

A Go program that wraps commands and converts YAML configuration files to command-line arguments. This allows you to use YAML configuration files with CLI programs that don't natively support YAML configuration.

## Usage

```bash
c7r <command> --config <config.yaml>
```

## Examples

### Basic Usage

```bash
c7r curl --config config.yaml
```

### Configuration File Format

The YAML configuration file supports two types of items:

1. **Positional arguments** (without `name` field):
   ```yaml
   - https://example.com/something.tar.gz
   ```

2. **Named flags** (with `name` and `value` fields):
   ```yaml
   - name: --output
     value: something.tar.gz
   - name: -H
     value: ["Authorization: Bearer token", "User-Agent: Chrome"]
   ```

3. **Named flags with custom joiner** (with `name`, `value`, and `joiner` fields):
   ```yaml
   - name: --output
     value: something.tar.gz
     joiner: =
   - name: -H
     value: "Authorization: Bearer token"
   - name: --custom
     value: value
     joiner: ": "
   ```

**Joiner options:**
- Space (default): `--flag value`
- `=`: `--flag=value`
- Custom string: `--flag: value` (with `joiner: ": "`)

### Environment Variable Support

The configuration supports environment variable expansion in all values:

```yaml
---
- https://example.com/$FILENAME
- name: --output
  value: $FILENAME
- name: -H
  value: ["Authorization: Bearer $TOKEN", "User-Agent: $USER_AGENT"]
- name: --config
  value: ${CONFIG_PATH}
```

**Supported syntax:**
- `$VAR` - Simple environment variable expansion
- `${VAR}` - Braced environment variable expansion

**Usage example:**
```bash
FILENAME=test.tar.gz TOKEN=abc123 USER_AGENT=Chrome CONFIG_PATH=/etc/config.json \
  c7r curl --config config.yaml
```

This will expand to:
```bash
curl https://example.com/test.tar.gz --output test.tar.gz \
  -H "Authorization: Bearer abc123" -H "User-Agent: Chrome" \
  --config /etc/config.json
```

### Example Configuration

```yaml
---
- https://example.com/something.tar.gz
- name: --output
  value: something.tar.gz
- name: -H
  value: ["Authorization: Bearer token", "User-Agent: Chrome"]
```

This configuration will be transformed to:
```bash
curl https://example.com/something.tar.gz --output=something.tar.gz -H "Authorization: Bearer token" -H "User-Agent: Chrome"
```

## Features

- **Positional Arguments**: Items without a `name` field are treated as positional arguments
- **Array Values**: Arrays are expanded into multiple flag instances
- **Automatic Quoting**: Values with whitespace are automatically quoted
- **Flexible Flag Format**: Control the joiner between flag name and value per item (space by default, `=` for equals, or custom)
- **Environment Variable Expansion**: Support for `$VAR` and `${VAR}` syntax in all values
- **Error Code Preservation**: Exit codes from the executed command are preserved
- **Stream Forwarding**: stdin, stdout, and stderr are forwarded to the executed command

## Command Line Options

- `--config`: Path to the YAML configuration file (required)
- `--dry-run`: Show the command that would be executed without running it

## Installation

```bash
go install github.com/moonlight8978/cmd-runner@latest
```

## Building from Source

```bash
git clone https://github.com/moonlight8978/cmd-runner.git
cd cmd-runner
go build -o c7r main.go
```

## Testing

You can test the program with the provided test configuration:

```bash
./c7r curl --config test-config.yaml --dry-run
```

Note: The `--dry-run` flag is not implemented in this version, but you can see the command that would be executed by adding debug output.
