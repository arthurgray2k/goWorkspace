package config

import (
	"strings"
)

// StringArrayFlag allows repeatable CLI flags (e.g. --pane "name:cmd").
type StringArrayFlag []string

func (s *StringArrayFlag) String() string {
	return strings.Join(*s, ", ")
}

func (s *StringArrayFlag) Set(value string) error {
	*s = append(*s, value)
	return nil
}

// ParsePaneFlag parses "name:command" or "name" string into a PaneConfig.
func ParsePaneFlag(raw string) PaneConfig {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return PaneConfig{}
	}

	parts := strings.SplitN(raw, ":", 2)
	if len(parts) == 1 {
		return PaneConfig{
			Name:    strings.TrimSpace(parts[0]),
			Command: "",
		}
	}

	return PaneConfig{
		Name:    strings.TrimSpace(parts[0]),
		Command: strings.TrimSpace(parts[1]),
	}
}
