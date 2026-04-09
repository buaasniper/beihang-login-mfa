package ingest

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"sort"
	"strings"
)

type EventInput struct {
	Fingerprint json.RawMessage
}

type NormalizedResult struct {
	CookieHash string
	CanvasHash string
	WebglHash  string
	FontsHash  string
	RestJSON   string
	CookieJSON string
	CanvasJSON string
	WebglJSON  string
	FontsJSON  string
}

func NormalizeAndHash(in EventInput) (*NormalizedResult, error) {
	if len(in.Fingerprint) == 0 {
		return nil, errors.New("empty fingerprint payload")
	}

	var raw any
	if err := json.Unmarshal(in.Fingerprint, &raw); err != nil {
		return nil, err
	}

	obj, ok := raw.(map[string]any)
	if !ok {
		return nil, errors.New("fingerprint must be a JSON object")
	}

	canvasPart := pickPart(obj, "canvas")
	webglPart := pickPart(obj, "webgl")
	fontsPart := pickPart(obj, "fonts")
	cookiePart := pickPart(obj, "cookie")

	cookieJSON := canonicalJSON(cookiePart)

	canvasJSON := canonicalJSON(canvasPart)
	webglJSON := canonicalJSON(webglPart)
	fontsJSON := canonicalJSON(fontsPart)

	rest := removeHeavyParts(obj, []string{"cookie", "canvas", "webgl", "fonts"})
	restJSON := canonicalJSON(rest)

	return &NormalizedResult{
		CookieHash: sha256Hex(cookieJSON),
		CanvasHash: sha256Hex(canvasJSON),
		WebglHash:  sha256Hex(webglJSON),
		FontsHash:  sha256Hex(fontsJSON),
		RestJSON:   restJSON,
		CookieJSON: cookieJSON,
		CanvasJSON: canvasJSON,
		WebglJSON:  webglJSON,
		FontsJSON:  fontsJSON,
	}, nil
}

func removeHeavyParts(obj map[string]any, keys []string) map[string]any {
	out := make(map[string]any, len(obj))
	for k, v := range obj {
		if hasKeyIgnoreCase(k, keys) {
			continue
		}
		out[k] = v
	}

	if _, ok := out["platform"]; !ok {
		if p := firstNonEmpty(asString(obj["platform"]), asString(obj["os"])); p != "" {
			out["platform"] = p
		}
	}
	if _, ok := out["language"]; !ok {
		if l := firstNonEmpty(asString(obj["language"]), asString(obj["lang"])); l != "" {
			out["language"] = l
		}
	}
	if _, ok := out["timezone"]; !ok {
		if tz := firstNonEmpty(asString(obj["timezone"]), asString(obj["tz"])); tz != "" {
			out["timezone"] = tz
		}
	}

	return out
}

func hasKeyIgnoreCase(k string, keys []string) bool {
	for _, candidate := range keys {
		if strings.EqualFold(k, candidate) {
			return true
		}
	}
	return false
}

func pickPart(obj map[string]any, key string) any {
	if v, ok := pickPartWithExists(obj, key); ok {
		return v
	}
	return map[string]any{}
}

func pickPartWithExists(obj map[string]any, key string) (any, bool) {
	if v, ok := obj[key]; ok {
		return v, true
	}
	for k, v := range obj {
		if strings.EqualFold(k, key) {
			return v, true
		}
	}
	return nil, false
}

func canonicalJSON(v any) string {
	canonical := canonicalize(v)
	b, err := json.Marshal(canonical)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func canonicalize(v any) any {
	switch x := v.(type) {
	case map[string]any:
		keys := make([]string, 0, len(x))
		for k := range x {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		ordered := make(map[string]any, len(x))
		for _, k := range keys {
			ordered[k] = canonicalize(x[k])
		}
		return ordered
	case []any:
		out := make([]any, len(x))
		for i := range x {
			out[i] = canonicalize(x[i])
		}
		return out
	default:
		return x
	}
}

func sha256Hex(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

func asString(v any) string {
	s, _ := v.(string)
	return strings.TrimSpace(s)
}

func firstNonEmpty(values ...string) string {
	for _, v := range values {
		if strings.TrimSpace(v) != "" {
			return strings.TrimSpace(v)
		}
	}
	return ""
}
