package main

import (
	"log"
	"log/slog"
	"os"

	httpserver "zapier-project/internal/http"
	"zapier-project/internal/http/handlers"

	"zapier-project/internal/deployments"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	seedData := deployments.SeedData()

	store := deployments.NewStore(seedData)
	logger.Info("store initialized", "seeded_events", len(seedData))

	deploymentsHandler := handlers.NewDeploymentsHandler(store, logger)
	router := httpserver.NewRouter(deploymentsHandler)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("server failed to start: %v", err)
	}
}
