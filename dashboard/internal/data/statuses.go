package data

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type StatusDefinition struct {
	ID             string
	Label          string
	Aliases        []string
	DashboardGroup string
	Rank           int
	Actionable     bool
	TopFilter      bool
}

type StatusConfig struct {
	States  []StatusDefinition
	ByID    map[string]StatusDefinition
	Aliases map[string]StatusDefinition
}

func loadStatusConfig(careerOpsPath string) StatusConfig {
	config := StatusConfig{
		ByID:    map[string]StatusDefinition{},
		Aliases: map[string]StatusDefinition{},
	}

	statesPath := resolveExistingPath(careerOpsPath,
		"templates/states.yml",
		"states.yml",
	)
	if statesPath == "" {
		return fallbackStatusConfig()
	}

	data, err := os.ReadFile(statesPath)
	if err != nil {
		return fallbackStatusConfig()
	}

	lines := strings.Split(string(data), "\n")
	var current *StatusDefinition

	flush := func() {
		if current == nil || current.ID == "" || current.Label == "" {
			return
		}
		config.States = append(config.States, *current)
	}

	for _, raw := range lines {
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") || line == "states:" {
			continue
		}
		if strings.HasPrefix(line, "- id:") {
			flush()
			current = &StatusDefinition{
				ID:    strings.TrimSpace(strings.TrimPrefix(line, "- id:")),
				Rank:  999,
				Label: "",
			}
			continue
		}
		if current == nil {
			continue
		}

		switch {
		case strings.HasPrefix(line, "label:"):
			current.Label = strings.TrimSpace(strings.TrimPrefix(line, "label:"))
		case strings.HasPrefix(line, "aliases:"):
			current.Aliases = parseInlineList(strings.TrimSpace(strings.TrimPrefix(line, "aliases:")))
		case strings.HasPrefix(line, "dashboard_group:"):
			current.DashboardGroup = strings.TrimSpace(strings.TrimPrefix(line, "dashboard_group:"))
		case strings.HasPrefix(line, "rank:"):
			if rank, err := strconv.Atoi(strings.TrimSpace(strings.TrimPrefix(line, "rank:"))); err == nil {
				current.Rank = rank
			}
		case strings.HasPrefix(line, "actionable:"):
			current.Actionable = strings.EqualFold(strings.TrimSpace(strings.TrimPrefix(line, "actionable:")), "true")
		case strings.HasPrefix(line, "top_filter:"):
			current.TopFilter = strings.EqualFold(strings.TrimSpace(strings.TrimPrefix(line, "top_filter:")), "true")
		}
	}

	flush()
	if len(config.States) == 0 {
		return fallbackStatusConfig()
	}

	for _, state := range config.States {
		config.ByID[state.ID] = state
		config.Aliases[strings.ToLower(state.ID)] = state
		config.Aliases[strings.ToLower(state.Label)] = state
		for _, alias := range state.Aliases {
			config.Aliases[strings.ToLower(alias)] = state
		}
	}

	for alias, targetID := range map[string]string{
		"hold":        "evaluated",
		"condicional": "evaluated",
		"evaluar":     "evaluated",
		"verificar":   "evaluated",
		"applied":     "applied",
		"sent":        "applied",
		"geo blocker": "skip",
		"duplicado":   "discarded",
		"dup":         "discarded",
		"repost":      "discarded",
	} {
		if state, ok := config.ByID[targetID]; ok {
			config.Aliases[alias] = state
		}
	}

	return config
}

func LoadStatusConfig(careerOpsPath string) StatusConfig {
	return loadStatusConfig(careerOpsPath)
}

func ResolveCareerOpsPath(path string) string {
	candidates := []string{path}
	clean := filepath.Clean(path)
	if filepath.Base(clean) == "dashboard" {
		candidates = append(candidates, filepath.Dir(clean))
	}

	for _, candidate := range candidates {
		if resolveExistingPath(candidate, "data/applications.md", "applications.md") != "" {
			return candidate
		}
	}

	return path
}

func fallbackStatusConfig() StatusConfig {
	states := []StatusDefinition{
		{ID: "offer", Label: "Offer", DashboardGroup: "offer", Rank: 0, Actionable: true, TopFilter: true},
		{ID: "interview", Label: "Interview", DashboardGroup: "interview", Rank: 10, Actionable: true, TopFilter: true},
		{ID: "responded", Label: "Responded", DashboardGroup: "responded", Rank: 20, Actionable: true, TopFilter: true},
		{ID: "applied", Label: "Applied", DashboardGroup: "applied", Rank: 30, Actionable: true, TopFilter: true},
		{ID: "evaluated", Label: "Evaluated", DashboardGroup: "evaluated", Rank: 40, Actionable: true, TopFilter: true},
		{ID: "skip", Label: "SKIP", DashboardGroup: "skip", Rank: 50, Actionable: false, TopFilter: false},
		{ID: "rejected", Label: "Rejected", DashboardGroup: "rejected", Rank: 60, Actionable: false, TopFilter: false},
		{ID: "discarded", Label: "Discarded", DashboardGroup: "discarded", Rank: 70, Actionable: false, TopFilter: false},
	}

	config := StatusConfig{
		States:  states,
		ByID:    map[string]StatusDefinition{},
		Aliases: map[string]StatusDefinition{},
	}
	for _, state := range states {
		config.ByID[state.ID] = state
		config.Aliases[strings.ToLower(state.ID)] = state
		config.Aliases[strings.ToLower(state.Label)] = state
	}
	return config
}

func parseInlineList(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	if !strings.HasPrefix(trimmed, "[") || !strings.HasSuffix(trimmed, "]") {
		return nil
	}
	items := strings.Split(strings.TrimSuffix(strings.TrimPrefix(trimmed, "["), "]"), ",")
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func resolveExistingPath(basePath string, candidates ...string) string {
	for _, candidate := range candidates {
		fullPath := filepath.Join(basePath, filepath.FromSlash(candidate))
		if _, err := os.Stat(fullPath); err == nil {
			return fullPath
		}
	}
	return ""
}

func normalizeStatusWithConfig(config StatusConfig, raw string) StatusDefinition {
	clean := strings.TrimSpace(strings.ReplaceAll(raw, "**", ""))
	if idx := strings.Index(clean, " 202"); idx > 0 {
		clean = strings.TrimSpace(clean[:idx])
	}
	if state, ok := config.Aliases[strings.ToLower(clean)]; ok {
		return state
	}
	return config.ByID["evaluated"]
}
