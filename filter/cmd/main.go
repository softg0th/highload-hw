package main

import (
	"errors"
	"filter/internal/core"
	"filter/internal/infra"
	"filter/internal/pkg"
	"filter/internal/service"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Config struct {
	grpcAddress     string
	bootstrapServer string
	minioEndpoint   string
	accessKeyId     string
	secretAccessKey string
	minioBucket     string
}

func initConfig() (*Config, error) {
	grpcAddress := os.Getenv("GRPC_ADDRESS")
	bootstrapServer := os.Getenv("BOOTSTRAP_SERVER")
	minioEndpoint := os.Getenv("MINIO_ENDPOINT")
	accessKeyId := os.Getenv("ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("SECRET_ACCESS_KEY")
	minioBucket := os.Getenv("MINIO_BUCKET")

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
	default:
	}

	cfg := &Config{
		grpcAddress:     grpcAddress,
		bootstrapServer: bootstrapServer,
		minioEndpoint:   minioEndpoint,
		accessKeyId:     accessKeyId,
		secretAccessKey: secretAccessKey,
		minioBucket:     minioBucket,
	}
	return cfg, nil
}

func main() {
	cfg, err := initConfig()
	log.Printf("init")
	if err != nil {
		log.Fatal(err)
		return
	}

	consumer, err := infra.NewKafkaConsumer(cfg.bootstrapServer)
	log.Printf("init consumer")
	if err != nil {
		log.Fatal(err)
		return
	}

	minio, err := infra.NewClient(cfg.minioEndpoint, cfg.accessKeyId, cfg.secretAccessKey, cfg.minioBucket)
	log.Printf("init minio")
	if err != nil {
		log.Fatal(err)
		return
	}

	rpc, err := infra.NewRPCConn(cfg.grpcAddress)
	log.Printf("init rpc")
	if err != nil {
		log.Printf("rpc failed", err)
		return
	}

	infraLayer := infra.NewInfra(consumer, rpc, minio)
	log.Printf("init layer")

	log.Printf("infra")
	pool := pkg.NewPool[[]byte](
		16,
		200,
		10*time.Second,
		service.FilterIt,
		infraLayer,
	)
	go func() {
		pool.Start()
	}()

	service, err := service.NewService(infraLayer, pool)
	log.Printf("init service")
	core.SetSpamProcessor(infraLayer.Minio.PutObject)
	core.StartSpamBatchJob(1000)
	if err != nil {
		log.Fatal(err)
	}
	service.RunLoop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	log.Println("Shutdown signal received")

}
