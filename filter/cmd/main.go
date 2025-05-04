package main

import (
	"errors"
	"filter/internal/core"
	"filter/internal/infra"
	"filter/internal/pkg"
	"filter/internal/service"
	logstash_logger "github.com/KaranJagtiani/go-logstash"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

type Config struct {
	grpcAddress      string
	bootstrapServer  string
	minioEndpoint    string
	accessKeyId      string
	secretAccessKey  string
	minioBucket      string
	logstashProtocol string
	logstashPort     int
}

func initConfig() (*Config, error) {
	grpcAddress := os.Getenv("GRPC_ADDRESS")
	bootstrapServer := os.Getenv("BOOTSTRAP_SERVER")
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("SECRET_ACCESS_KEY")
	minioBucket := os.Getenv("MINIO_BUCKET")

	logstashProtocol := os.Getenv("LOGSTASH_PROTOCOL")
	logstashPort, err := strconv.Atoi(os.Getenv("LOGSTASH_PORT"))

	if err != nil {
		return nil, errors.New("LOGSTASH_PORT not set")
	}

	switch {
	case grpcAddress == "":
		return nil, errors.New("GRPC_ADDRESS not set")
	case bootstrapServer == "":
		return nil, errors.New("KAFKA_CONNECT not set")
	case minioEndpoint == "":
		return nil, errors.New("MINIO_ENDPOINT not set")
	case accessKeyId == "":
		return nil, errors.New("ACCESS_KEY_ID not set")
	case secretAccessKey == "":
		return nil, errors.New("SECRET_ACCESS_KEY not set")
	case minioBucket == "":
		return nil, errors.New("MINIO_BUCKET not set")
	case logstashProtocol == "":
		return nil, errors.New("LOGSTASH_PROTOCOL not set")
	default:
	}

	cfg := &Config{
		grpcAddress:      grpcAddress,
		bootstrapServer:  bootstrapServer,
		minioEndpoint:    minioEndpoint,
		accessKeyId:      accessKeyId,
		secretAccessKey:  secretAccessKey,
		minioBucket:      minioBucket,
		logstashProtocol: logstashProtocol,
		logstashPort:     logstashPort,
	}
	return cfg, nil
}

func main() {
	cfg, err := initConfig()

	if err != nil {
		log.Fatal(err)
		return
	}

	logger := logstash_logger.Init("logstash", cfg.logstashPort, cfg.logstashProtocol, 5)

	consumer, err := infra.NewKafkaConsumer(cfg.bootstrapServer)
	logger.Info(map[string]interface{}{
		"message": "Init consumer",
		"error":   false,
	})

	if err != nil {
		logger.Error(map[string]interface{}{
			"message":           "Failed to setup consumer",
			"error":             true,
			"error_description": err.Error(),
		})
		return
	}

	minio, err := infra.NewClient(cfg.minioEndpoint, cfg.accessKeyId, cfg.secretAccessKey, cfg.minioBucket)
	logger.Info(map[string]interface{}{
		"message": "Init minio",
		"error":   false,
	})

	if err != nil {
		logger.Error(map[string]interface{}{
			"message":           "Failed to setup minio",
			"error":             true,
			"error_description": err.Error(),
		})
		return
	}

	rpc, err := infra.NewRPCConn(cfg.grpcAddress)
	logger.Info(map[string]interface{}{
		"message": "Init rpc",
		"error":   false,
	})

	if err != nil {
		logger.Error(map[string]interface{}{
			"message":           "Failed to setup rpc",
			"error":             true,
			"error_description": err.Error(),
		})
		return
	}

	infraLayer := infra.NewInfra(consumer, rpc, minio)

	pool := pkg.NewPool[[]byte](
		16,
		200,
		10*time.Second,
		service.FilterIt,
		infraLayer,
		logger,
	)
	go func() {
		pool.Start()
	}()

	service, err := service.NewService(infraLayer, pool, logger)

	core.SetSpamProcessor(infraLayer.Minio.PutObject)
	core.StartSpamBatchJob(1000, logger)
	if err != nil {
		log.Fatal(err)
	}
	service.RunLoop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	logger.Info(map[string]interface{}{
		"message": "Shutdown signal received",
		"error":   false,
	})
}
