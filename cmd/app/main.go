package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/stan.go"
	"golang.org/x/sync/errgroup"

	"wb-tech/internal/api"
	"wb-tech/internal/config"
	db2 "wb-tech/internal/pkg/db"
	"wb-tech/internal/repository"
	"wb-tech/internal/services"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle system interrupts for graceful shutdown
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := db2.OpenDB(ctx, cfg.DBConfig)
	if err != nil {
		log.Fatalf("Failed to connect to DB: %v", err)
	}

	conn, err := stan.Connect(
		cfg.NatsConfig.ClusterID,
		cfg.NatsConfig.ClientID,
		stan.NatsURL(cfg.NatsConfig.NatsURL),
	)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}

	service := services.New(conn, repository.New(db), db)
	if err = service.LoadCache(ctx); err != nil {
		log.Printf("Error loading cache: %v", err)
	}

	r := api.New(service)
	srv := &http.Server{
		Addr:    ":" + cfg.ServerConfig.HTTPPort,
		Handler: r,
	}

	eg, gCtx := errgroup.WithContext(ctx)

	// Run HTTP server
	eg.Go(func() error {
		log.Printf("Server is running on port %s", cfg.ServerConfig.HTTPPort)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
		return nil
	})

	// Run NATS listener
	eg.Go(func() error {
		return service.AddFromChannel(gCtx, cfg.NatsConfig.NatsChan)
	})

	// Wait for system signals
	eg.Go(func() error {
		select {
		case <-sigCh:
			log.Println("Received interrupt, shutting down...")
			cancel()
			// Shutdown the HTTP server
			if err := srv.Shutdown(ctx); err != nil {
				return err
			}
		case <-gCtx.Done():
		}
		return nil
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Printf("Server exited with error: %v", err)
	}
}
