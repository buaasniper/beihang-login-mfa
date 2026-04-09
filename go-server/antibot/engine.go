package antibot

import (
	"encoding/json"
	"strings"

	"gorm.io/gorm"
)

type EvaluationInput struct {
	UserAgent string
	RestJSON  string
}

type EvaluationResult struct {
	Decision string
	Reasons  []string
}

func EvaluateWithDB(db *gorm.DB, input EvaluationInput) EvaluationResult {
	if db == nil {
		return Evaluate(input)
	}
	res, err := evaluateWithDB(db, input)
	if err != nil {
		return Evaluate(input)
	}
	return res
}

func Evaluate(input EvaluationInput) EvaluationResult {
	ua := strings.ToLower(strings.TrimSpace(input.UserAgent))
	rest := parseRestJSON(input.RestJSON)

	strongHits := make([]string, 0, 4)
	weakHits := make([]string, 0, 4)

	if containsAny(ua, []string{"headless", "headlesschrome"}) {
		strongHits = append(strongHits, "ua_headless")
	}
	if containsAny(ua, []string{"playwright", "puppeteer", "phantomjs", "selenium"}) {
		strongHits = append(strongHits, "ua_automation_framework")
	}
	if asBool(getAnyByPath(rest, "webdriver", "navigator.webdriver", "window.webdriver")) {
		strongHits = append(strongHits, "webdriver_true")
	}

	langValue, langExists := getAnyByPathWithExists(rest, "languages", "navigator.languages")
	if langExists && isEmptyLanguage(langValue) {
		weakHits = append(weakHits, "languages_empty")
	}

	pluginsValue, pluginsExists := getAnyByPathWithExists(rest, "pluginsLength", "plugins_length", "plugins", "navigator.plugins", "plugins.length")
	if pluginsExists && isZeroPlugins(pluginsValue) {
		weakHits = append(weakHits, "plugins_zero")
	}

	if hasReason(weakHits, "languages_empty") && hasReason(weakHits, "plugins_zero") {
		weakHits = append(weakHits, "lang_and_plugins_suspicious")
	}

	if len(strongHits) > 0 {
		return EvaluationResult{
			Decision: "bot",
			Reasons:  append(strongHits, weakHits...),
		}
	}

	if hasReason(weakHits, "languages_empty") && hasReason(weakHits, "plugins_zero") {
		return EvaluationResult{
			Decision: "challenge",
			Reasons:  weakHits,
		}
	}

	return EvaluationResult{
		Decision: "pass",
		Reasons:  []string{},
	}
}

func parseRestJSON(s string) map[string]any {
	out := map[string]any{}
	if strings.TrimSpace(s) == "" {
		return out
	}
	_ = json.Unmarshal([]byte(s), &out)
	return out
}

func containsAny(s string, needles []string) bool {
	for _, n := range needles {
		if strings.Contains(s, n) {
			return true
		}
	}
	return false
}

func getAnyByPath(root map[string]any, paths ...string) any {
	v, _ := getAnyByPathWithExists(root, paths...)
	return v
}

func getAnyByPathWithExists(root map[string]any, paths ...string) (any, bool) {
	for _, p := range paths {
		for _, candidate := range expandCompatiblePaths(p) {
			parts := strings.Split(candidate, ".")
			var cur any = root
			ok := true
			for _, part := range parts {
				m, isMap := cur.(map[string]any)
				if !isMap {
					ok = false
					break
				}
				next, exists := m[part]
				if !exists {
					ok = false
					break
				}
				cur = next
			}
			if ok {
				return cur, true
			}
		}
	}
	return nil, false
}

func expandCompatiblePaths(path string) []string {
	base := strings.TrimSpace(path)
	if base == "" {
		return []string{}
	}

	parts := strings.Split(base, ".")
	if len(parts) > 0 {
		head := strings.ToLower(strings.TrimSpace(parts[0]))
		switch head {
		case "leve1", "level1", "leve2", "level2", "leve3", "level3", "navigator", "window":
			return []string{base}
		}
	}

	return []string{
		base,
		"leve1." + base,
		"level1." + base,
		"leve2." + base,
		"level2." + base,
		"leve3." + base,
		"level3." + base,
	}
}

func asBool(v any) bool {
	switch x := v.(type) {
	case bool:
		return x
	case string:
		s := strings.ToLower(strings.TrimSpace(x))
		return s == "1" || s == "true" || s == "yes"
	default:
		return false
	}
}

func isEmptyLanguage(v any) bool {
	switch x := v.(type) {
	case nil:
		return true
	case string:
		return strings.TrimSpace(x) == ""
	case []any:
		if len(x) == 0 {
			return true
		}
		if len(x) == 1 {
			s, _ := x[0].(string)
			return strings.TrimSpace(s) == ""
		}
		return false
	default:
		return false
	}
}

func isZeroPlugins(v any) bool {
	switch x := v.(type) {
	case nil:
		return true
	case float64:
		return x == 0
	case int:
		return x == 0
	case int64:
		return x == 0
	case string:
		s := strings.TrimSpace(x)
		return s == "" || s == "0"
	case []any:
		return len(x) == 0
	case map[string]any:
		if raw, ok := x["length"]; ok {
			return isZeroPlugins(raw)
		}
		return len(x) == 0
	default:
		return false
	}
}

func hasReason(reasons []string, target string) bool {
	for _, r := range reasons {
		if r == target {
			return true
		}
	}
	return false
}
