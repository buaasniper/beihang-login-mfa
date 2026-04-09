package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"go-server/antibot"
	"go-server/config"
	"go-server/handler"
	"go-server/model"
)

//go:embed static/*
var staticFS embed.FS

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "*")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	db := config.InitDB()

	if err := db.AutoMigrate(
		&model.User{},
		&model.FingerprintLog{},
		&model.BfpEvent{},
		&model.CanvasHashLibrary{},
		&model.WebglHashLibrary{},
		&model.FontsHashLibrary{},
		&model.OnlineRuleConfig{},
		&model.RiskAccountProfile{},
		&model.RiskUAProfile{},
		&model.RiskHashProfile{},
		&model.OnlineBotResult{},
	); err != nil {
		log.Fatalf("auto migrate failed: %v", err)
	}

	if err := antibot.EnsureBootstrapRules(db); err != nil {
		log.Printf("failed to bootstrap online rules: %v", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("POST /login", handler.Login(db))
	mux.HandleFunc("POST /register", handler.Register(db))
	mux.HandleFunc("POST /fingerprint", handler.CollectFingerprint(db))

	staticContent, _ := fs.Sub(staticFS, "static")
	fileServer := http.FileServer(http.FS(staticContent))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/login.html", http.StatusFound)
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	srv := corsMiddleware(mux)

	fmt.Println("Server started on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", srv))
}
