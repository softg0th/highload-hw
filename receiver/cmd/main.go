package main

import (
	"context"
	"errors"
	logstash_logger "github.com/KaranJagtiani/go-logstash"
	"log"
	"os"
	"os/signal"
	"receiver/internal/infra"
	"receiver/internal/server"
	"receiver/internal/service"
	"strconv"
	"syscall"
	"time"
)

type Config struct {
	Port             string
	KafkaConn        string
	logstashProtocol string
	logstashPort     int
}

func initConfig() (*Config, error) {
	serverPort := os.Getenv("SERVER_PORT")
	kafkaConn := os.Getenv("KAFKA_CONNECT")
	logstashProtocol := os.Getenv("LOGSTASH_PROTOCOL")

	if serverPort == "" {
		return nil, errors.New("SERVER_PORT not set")
	}

	if kafkaConn == "" {
		return nil, errors.New("KAFKA_PORT not set")
	}

	logstashPort, err := strconv.Atoi(os.Getenv("LOGSTASH_PORT"))

	if err != nil {
		return nil, errors.New("LOGSTASH_PORT not set")
	}

	cfg := &Config{
		Port:             serverPort,
		KafkaConn:        kafkaConn,
		logstashPort:     logstashPort,
		logstashProtocol: logstashProtocol,
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

	logger := logstash_logger.Init("logstash", cfg.logstashPort, cfg.logstashProtocol, 5)

	workerPool := infra.NewPool(infraLayer, 16, 200, 10000000000)
	sender := service.NewSender(workerPool, logger)
	services := service.NewService(sender)

	go func() {
		workerPool.Start()
	}()

	srv := server.Setup(infraLayer, services)

	logger.Info(map[string]interface{}{
		"message": "Server listening on port",
		"error":   false,
		"port":    cfg.Port,
	})

	go func() {
		if err := srv.Start(cfg.Port); err != nil {
			logger.Error(map[string]interface{}{"Server error": err, "error": true})
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error(map[string]interface{}{"Graceful shutdown failed:": err, "error": true})
	}

	logger.Info(map[string]interface{}{
		"message": "Server stopped",
		"error":   false,
	})
}
