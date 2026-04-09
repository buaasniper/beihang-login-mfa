package antibot

import (
	"strings"
	"time"

	"go-server/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type compiledRule struct {
	RuleKey  string
	Source   string
	FieldPath string
	Operator string
	Value    string
	Decision string
	Priority int
}

func EnsureBootstrapRules(db *gorm.DB) error {
	now := time.Now()
	rules := []model.OnlineRuleConfig{
		{
			RuleKey:   "ua_headless",
			Enabled:   true,
			Priority:  10,
			Source:    "ua",
			FieldPath: "",
			Operator:  "contains",
			Value:     "headless",
			Decision:  "bot",
			Version:   "v0.1",
			UpdatedAt: now,
		},
		{
			RuleKey:   "ua_automation_framework",
			Enabled:   true,
			Priority:  11,
			Source:    "ua",
			FieldPath: "",
			Operator:  "contains_any",
			Value:     "playwright,puppeteer,phantomjs,selenium",
			Decision:  "bot",
			Version:   "v0.1",
			UpdatedAt: now,
		},
		{
			RuleKey:   "webdriver_true",
			Enabled:   true,
			Priority:  12,
			Source:    "rest",
			FieldPath: "webdriver|navigator.webdriver|window.webdriver",
			Operator:  "bool_true",
			Value:     "",
			Decision:  "bot",
			Version:   "v0.1",
			UpdatedAt: now,
		},
		{
			RuleKey:   "languages_empty",
			Enabled:   true,
			Priority:  100,
			Source:    "rest",
			FieldPath: "languages|navigator.languages",
			Operator:  "empty",
			Value:     "",
			Decision:  "pass",
			Version:   "v0.1",
			UpdatedAt: now,
		},
		{
			RuleKey:   "plugins_zero",
			Enabled:   true,
			Priority:  101,
			Source:    "rest",
			FieldPath: "pluginsLength|plugins_length|plugins|navigator.plugins|plugins.length",
			Operator:  "zero",
			Value:     "",
			Decision:  "pass",
			Version:   "v0.1",
			UpdatedAt: now,
		},
		{
			RuleKey:   "lang_and_plugins_suspicious",
			Enabled:   true,
			Priority:  102,
			Source:    "rest",
			FieldPath: "languages|navigator.languages|pluginsLength|plugins_length|plugins|navigator.plugins|plugins.length",
			Operator:  "lang_and_plugins_suspicious",
			Value:     "",
			Decision:  "challenge",
			Version:   "v0.1",
			UpdatedAt: now,
		},
	}

	for _, row := range rules {
		if err := db.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "rule_key"}},
			DoUpdates: clause.Assignments(map[string]any{
				"enabled":    row.Enabled,
				"priority":   row.Priority,
				"source":     row.Source,
				"field_path": row.FieldPath,
				"operator":   row.Operator,
				"value":      row.Value,
				"decision":   row.Decision,
				"version":    row.Version,
				"updated_at": row.UpdatedAt,
			}),
		}).Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

func evaluateWithDB(db *gorm.DB, input EvaluationInput) (EvaluationResult, error) {
	var rows []model.OnlineRuleConfig
	if err := db.Where("enabled = ?", true).Order("priority ASC, id ASC").Find(&rows).Error; err != nil {
		return EvaluationResult{}, err
	}

	if len(rows) == 0 {
		return Evaluate(input), nil
	}

	ua := strings.ToLower(strings.TrimSpace(input.UserAgent))
	rest := parseRestJSON(input.RestJSON)

	botReasons := make([]string, 0, 4)
	challengeReasons := make([]string, 0, 4)

	for _, row := range rows {
		rule := compiledRule{
			RuleKey:   row.RuleKey,
			Source:    row.Source,
			FieldPath: row.FieldPath,
			Operator:  row.Operator,
			Value:     row.Value,
			Decision:  strings.ToLower(strings.TrimSpace(row.Decision)),
			Priority:  row.Priority,
		}
		if !matchRule(rule, ua, rest) {
			continue
		}
		if rule.Decision == "bot" {
			botReasons = append(botReasons, rule.RuleKey)
			continue
		}
		if rule.Decision == "challenge" {
			challengeReasons = append(challengeReasons, rule.RuleKey)
		}
	}

	if len(botReasons) > 0 {
		return EvaluationResult{Decision: "bot", Reasons: append(botReasons, challengeReasons...)}, nil
	}
	if len(challengeReasons) > 0 {
		return EvaluationResult{Decision: "challenge", Reasons: challengeReasons}, nil
	}
	return EvaluationResult{Decision: "pass", Reasons: []string{}}, nil
}

func matchRule(rule compiledRule, ua string, rest map[string]any) bool {
	source := strings.ToLower(strings.TrimSpace(rule.Source))
	op := strings.ToLower(strings.TrimSpace(rule.Operator))
	value := strings.ToLower(strings.TrimSpace(rule.Value))

	var current any
	currentExists := true
	if source == "ua" {
		current = ua
	} else {
		paths := splitPaths(rule.FieldPath)
		current, currentExists = getAnyByPathWithExists(rest, paths...)
	}

	switch op {
	case "contains":
		s, ok := current.(string)
		if !ok {
			if source == "ua" {
				s = ua
			} else {
				return false
			}
		}
		return strings.Contains(strings.ToLower(s), value)
	case "contains_any":
		s, ok := current.(string)
		if !ok {
			if source == "ua" {
				s = ua
			} else {
				return false
			}
		}
		for _, token := range splitCSV(value) {
			if strings.Contains(strings.ToLower(s), token) {
				return true
			}
		}
		return false
	case "bool_true":
		if source == "rest" && !currentExists {
			return false
		}
		return asBool(current)
	case "empty":
		if source == "rest" && !currentExists {
			return false
		}
		return isEmptyLanguage(current)
	case "zero":
		if source == "rest" && !currentExists {
			return false
		}
		return isZeroPlugins(current)
	case "lang_and_plugins_suspicious":
		langVal, langOK := getAnyByPathWithExists(rest, "languages", "navigator.languages")
		pluginsVal, pluginsOK := getAnyByPathWithExists(rest, "pluginsLength", "plugins_length", "plugins", "navigator.plugins", "plugins.length")
		return langOK && pluginsOK && isEmptyLanguage(langVal) && isZeroPlugins(pluginsVal)
	default:
		return false
	}
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}

func splitPaths(s string) []string {
	parts := strings.Split(s, "|")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		v := strings.TrimSpace(p)
		if v != "" {
			out = append(out, v)
		}
	}
	return out
}
