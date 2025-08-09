package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"planner-pro/internal/api"
	"planner-pro/internal/auth"
	"planner-pro/internal/infra"

	"github.com/go-chi/chi/v5"
)

func main() {
	ctx := context.Background()
	cfg := infra.LoadConfigFromEnv()

	pool, err := infra.NewPgPool(ctx, cfg.DatabaseURL, cfg.DBMaxConns)
	if err != nil {
		log.Fatalf("failed connecting to db: %v", err)
	}
	defer pool.Close()

	r := chi.NewRouter()
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK); w.Write([]byte("ok")) })

	// OIDC middleware con retry/backoff
	var oidcMiddleware func(http.Handler) http.Handler
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		oidcMiddleware, err = auth.NewMiddleware(ctx, cfg.OIDC.Issuer, cfg.OIDC.Audience, cfg.OIDC.SkipVerify)
		if err == nil && oidcMiddleware != nil {
			break
		}
		log.Printf("OIDC middleware init failed (tentativo %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(3 * time.Second)
	}
	if oidcMiddleware == nil {
		log.Fatalf("OIDC middleware non disponibile dopo %d tentativi: %v", maxRetries, err)
	}

	api.RegisterRoutes(r, pool, oidcMiddleware)

	srv := &http.Server{
		Addr:         ":8000",
		Handler:      r,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	log.Printf("listening on %s", srv.Addr)
	if err := srv.ListenAndServe(); err != nil {
		log.Printf("server stopped: %v", err)
		os.Exit(1)
	}
}
