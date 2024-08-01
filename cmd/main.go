package main

import (
	"FinUslugi/internal/config"
	"FinUslugi/internal/logger"
	"FinUslugi/internal/server"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	sugar, err := logger.New()
	if err != nil {
		log.Fatalf("failed create new logger: %v", err)
	}
	cfg, err := config.Parse()
	if err != nil {
		log.Fatalf("failed parse config: %v", err)
	}

	srv, err := server.NewServer(sugar)
	if err != nil {
		sugar.Fatal(err)
	}
	defer srv.DB.Close()

	http.HandleFunc("/materials", srv.HandleMaterials)
	http.HandleFunc("/materials/", srv.HandleMaterial)

	sugar.Infof("Starting server on port %d", cfg.HTTP.Port)

	server := &http.Server{
		Addr:    ":" + fmt.Sprintf("%d", cfg.HTTP.Port),
		Handler: nil,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			sugar.Fatalf("listen: %s\n", err)
		}
	}()
	sugar.Info("Server started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	sugar.Info("Shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		sugar.Fatal("Server forced to shutdown: ", err)
	}

	sugar.Info("Server exiting")
}
