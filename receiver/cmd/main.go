package main

import (
	"context"
	"errors"
	"log"
	"os"
	"os/signal"
	"receiver/internal/infra"
	"receiver/internal/server"
	"receiver/internal/service"
	"syscall"
	"time"
)

type Config struct {
	Port      string
	KafkaConn string
}

func initConfig() (*Config, error) {
	serverPort := os.Getenv("SERVER_PORT")
	kafkaConn := os.Getenv("KAFKA_CONNECT")

	if serverPort == "" {
		return nil, errors.New("SERVER_PORT not set")
	}

	if kafkaConn == "" {
		return nil, errors.New("KAFKA_PORT not set")
	}

	cfg := &Config{
		Port:      serverPort,
		KafkaConn: kafkaConn,
	}

	return cfg, nil
}

func main() {
	cfg, err := initConfig()
	if err != nil {
		log.Fatal(err)
		return
	}

	infraLayer, err := infra.NewInfra(cfg.KafkaConn)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer infraLayer.Producer.Close()

	workerPool := infra.NewPool(infraLayer, 16, 200, 10000000000)
	sender := service.NewSender(workerPool)
	services := service.NewService(sender)

	go func() {
		workerPool.Start()
	}()

	srv := server.Setup(infraLayer, services)

	log.Printf("server listening on port %s", cfg.Port)
	go func() {
		if err := srv.Start(cfg.Port); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Graceful shutdown failed:", err)
	}

	log.Println("Server gracefully stopped")
}
