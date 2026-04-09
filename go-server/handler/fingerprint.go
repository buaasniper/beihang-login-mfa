package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"go-server/antibot"
	"go-server/ingest"
	"go-server/model"

	"golang.org/x/time/rate"
	"gorm.io/gorm/clause"
	"gorm.io/gorm"
)

type collectFingerprintResponse struct {
	Status      string   `json:"status"`
	EventID     uint     `json:"event_id"`
	BotDecision string   `json:"bot_decision"`
	Reasons     []string `json:"reasons"`
}

// 全局限流器：比如每秒允许最多1000个请求，突发(burst)最高2000。
// 这样在日常情况下完全不会被触发，高并发突发时能保护后端数据库。
var ingestLimiter = rate.NewLimiter(1000, 2000)

type fingerprintRequest struct {
	Username    string          `json:"username"`
	Fingerprint json.RawMessage `json:"fingerprint"`
	URL         string          `json:"url"`
	DeltaTime   *float64        `json:"delta_time"`
	ClickTime   string          `json:"click_time"`
}

func CollectFingerprint(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 简单的限流：如果不允许通过，则直接返回 429 Too Many Requests
		if !ingestLimiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		var req fingerprintRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		fpStr := string(req.Fingerprint)

		log := model.FingerprintLog{
			Username:    req.Username,
			Fingerprint: fpStr,
			URL:         req.URL,
			DeltaTime:   req.DeltaTime,
			ClickTime:   req.ClickTime,
		}

		if err := db.Create(&log).Error; err != nil {
			http.Error(w, "Failed to write raw log", http.StatusInternalServerError)
			return
		}

		norm, err := ingest.NormalizeAndHash(ingest.EventInput{Fingerprint: req.Fingerprint})
		if err != nil {
			http.Error(w, "Invalid fingerprint payload", http.StatusBadRequest)
			return
		}

		event := model.BfpEvent{
			Username:   req.Username,
			URL:        req.URL,
			DeltaTime:  req.DeltaTime,
			ClickTime:  req.ClickTime,
			CookieHash: norm.CookieHash,
			CanvasHash: norm.CanvasHash,
			WebglHash:  norm.WebglHash,
			FontsHash:  norm.FontsHash,
			UserAgent:  r.UserAgent(),
			RestJSON:   norm.RestJSON,
		}
		if err := db.Create(&event).Error; err != nil {
			http.Error(w, "Failed to write bfp_event", http.StatusInternalServerError)
			return
		}

		eval := antibot.EvaluateWithDB(db, antibot.EvaluationInput{
			UserAgent: event.UserAgent,
			RestJSON:  event.RestJSON,
		})

		reasonsJSON := mustJSON(eval.Reasons)
		botResult := model.OnlineBotResult{
			BfpEventID:  event.ID,
			Username:    req.Username,
			Decision:    eval.Decision,
			ReasonsJSON: reasonsJSON,
			CreatedAt:   time.Now(),
		}
		if err := db.Create(&botResult).Error; err != nil {
			http.Error(w, "Failed to write online_bot_result", http.StatusInternalServerError)
			return
		}

		now := time.Now()
		upsertCanvasHash(db, norm.CanvasHash, norm.CanvasJSON, now)
		upsertWebglHash(db, norm.WebglHash, norm.WebglJSON, now)
		upsertFontsHash(db, norm.FontsHash, norm.FontsJSON, now)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(collectFingerprintResponse{
			Status:      "ok",
			EventID:     event.ID,
			BotDecision: eval.Decision,
			Reasons:     eval.Reasons,
		})
	}
}

func mustJSON(v any) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "[]"
	}
	return string(b)
}

func upsertCanvasHash(db *gorm.DB, hash, sample string, now time.Time) {
	if hash == "" {
		return
	}

	row := model.CanvasHashLibrary{
		Hash:       hash,
		SampleJSON: sample,
		SeenCount:  1,
		FirstSeen:  now,
		LastSeen:   now,
	}

	db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "hash"}},
		DoUpdates: clause.Assignments(map[string]any{
			"seen_count":  gorm.Expr("seen_count + 1"),
			"last_seen":   now,
			"sample_json": sample,
		}),
	}).Create(&row)
}

func upsertWebglHash(db *gorm.DB, hash, sample string, now time.Time) {
	if hash == "" {
		return
	}

	row := model.WebglHashLibrary{
		Hash:       hash,
		SampleJSON: sample,
		SeenCount:  1,
		FirstSeen:  now,
		LastSeen:   now,
	}

	db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "hash"}},
		DoUpdates: clause.Assignments(map[string]any{
			"seen_count":  gorm.Expr("seen_count + 1"),
			"last_seen":   now,
			"sample_json": sample,
		}),
	}).Create(&row)
}

func upsertFontsHash(db *gorm.DB, hash, sample string, now time.Time) {
	if hash == "" {
		return
	}

	row := model.FontsHashLibrary{
		Hash:       hash,
		SampleJSON: sample,
		SeenCount:  1,
		FirstSeen:  now,
		LastSeen:   now,
	}

	db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "hash"}},
		DoUpdates: clause.Assignments(map[string]any{
			"seen_count":  gorm.Expr("seen_count + 1"),
			"last_seen":   now,
			"sample_json": sample,
		}),
	}).Create(&row)
}

