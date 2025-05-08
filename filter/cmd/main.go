package main

import (
	"context"
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
	prometheusPort   string
}

func initConfig() (*Config, error) {
	grpcAddress := os.Getenv("GRPC_ADDRESS")
	bootstrapServer := os.Getenv("BOOTSTRAP_SERVER")
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("SECRET_ACCESS_KEY")
	minioBucket := os.Getenv("MINIO_BUCKET")
	prometheusPort := os.Getenv("PROMETHEUS_PORT")
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
	case prometheusPort == "":
		return nil, errors.New("PROMETHEUS_PORT not set")
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
		prometheusPort:   prometheusPort,
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

	logger.Info(map[string]interface{}{
		"message": "Successfully connected to Logstash",
		"error":   false,
	})

	prometheusPort := os.Getenv("PROMETHEUS_PORT")

	ctx := context.Background()

	go func() {
		if err := <-infra.StartPrometheus(ctx, prometheusPort, logger); err != nil {
			log.Println(err)
			logger.Error(map[string]interface{}{
				"message":           "Failed to setup prometheus server",
				"error":             true,
				"error_description": err.Error(),
			})
		}
	}()

	consumer, err := infra.NewKafkaConsumer(cfg.bootstrapServer)

	log.Printf("Init consumer")
	logger.Info(map[string]interface{}{
		"message": "Init consumer",
		"error":   false,
	})

	if err != nil {
		log.Println(err)
		logger.Error(map[string]interface{}{
			"message":           "Failed to setup consumer",
			"error":             true,
			"error_description": err.Error(),
		})
		return
	}

	minio, err := infra.NewClient(cfg.minioEndpoint, cfg.accessKeyId, cfg.secretAccessKey, cfg.minioBucket)

	log.Printf("Init minio")
	logger.Info(map[string]interface{}{
		"message": "Init minio",
		"error":   false,
	})

	if err != nil {
		log.Println(err)
		logger.Error(map[string]interface{}{
			"message":           "Failed to setup minio",
			"error":             true,
			"error_description": err.Error(),
		})
		return
	}

	var rpc *infra.RPCConn
	const maxAttempts = 10
	const retryDelay = 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		rpc, err = infra.NewRPCConn(cfg.grpcAddress)
		if err == nil {
			log.Printf("Connected to RPC on attempt %d", attempt)
			logger.Info(map[string]interface{}{
				"message": "Connected to RPC",
				"attempt": attempt,
				"error":   false,
			})
			break
		}

		log.Printf("RPC connection failed (attempt %d/%d): %v", attempt, maxAttempts, err)
		logger.Warn(map[string]interface{}{
			"message":           "RPC connection failed",
			"attempt":           attempt,
			"error":             true,
			"error_description": err.Error(),
		})
		time.Sleep(retryDelay)
	}

	if err != nil {
		log.Println("RPC setup failed after retries:", err)
		logger.Error(map[string]interface{}{
			"message":           "RPC setup failed after retries",
			"error":             true,
			"error_description": err.Error(),
		})
		return
	}

	log.Printf("Init rpc")
	logger.Info(map[string]interface{}{
		"message": "Init rpc",
		"error":   false,
	})

	if err != nil {
		log.Println(err)
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
	log.Printf("ok")
}
