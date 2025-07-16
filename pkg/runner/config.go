package runner

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

var defaultJoiner string = " "

// ConfigItem represents a single configuration item in the YAML
type ConfigItem struct {
	Name   string      `yaml:"name"`
	Value  interface{} `yaml:"value"`
	Joiner *string     `yaml:"joiner,omitempty"` // Per-item control for joiner (default: "space")
}

// Config represents the entire configuration structure with global options
type Config struct {
	Items []ConfigItem `yaml:"items,omitempty"` // Configuration items
}

// ParseConfig reads and parses the YAML configuration file
func ParseConfig(configPath string) ([]ConfigItem, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var rawItems []interface{}
	if err := yaml.Unmarshal(data, &rawItems); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	var configItems []ConfigItem
	for _, raw := range rawItems {
		switch v := raw.(type) {
		case map[string]interface{}:
			item := ConfigItem{}
			if name, ok := v["name"].(string); ok {
				item.Name = name
			}
			if val, ok := v["value"]; ok {
				item.Value = val
			}
			if joiner, ok := v["joiner"].(string); ok {
				item.Joiner = &joiner
			} else {
				item.Joiner = &defaultJoiner
			}
			configItems = append(configItems, item)
		case string:
			configItems = append(configItems, ConfigItem{Name: "", Value: v})
		default:
			return nil, fmt.Errorf("unsupported config item type: %T", v)
		}
	}

	return configItems, nil
}

// ParseConfigWithOptions reads and parses the YAML configuration file with global options
func ParseConfigWithOptions(configPath string) ([]ConfigItem, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	// Try to parse as structured config first
	var structuredConfig Config
	if err := yaml.Unmarshal(data, &structuredConfig); err == nil && len(structuredConfig.Items) > 0 {
		return structuredConfig.Items, nil
	}

	// Fall back to simple array format
	var rawItems []interface{}
	if err := yaml.Unmarshal(data, &rawItems); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	var configItems []ConfigItem
	for _, raw := range rawItems {
		switch v := raw.(type) {
		case map[string]interface{}:
			item := ConfigItem{}
			if name, ok := v["name"].(string); ok {
				item.Name = name
			}
			if val, ok := v["value"]; ok {
				item.Value = val
			}
			if joiner, ok := v["joiner"].(string); ok {
				item.Joiner = &joiner
			} else {
				item.Joiner = &defaultJoiner
			}

			configItems = append(configItems, item)
		case string:
			configItems = append(configItems, ConfigItem{Name: "", Value: v})
		default:
			return nil, fmt.Errorf("unsupported config item type: %T", v)
		}
	}

	return configItems, nil
}
