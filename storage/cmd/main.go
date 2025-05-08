package main

import (
	"errors"
	logstash_logger "github.com/KaranJagtiani/go-logstash"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"storage/internal/infra"
	pb "storage/internal/proto"
	"storage/internal/repository"
	"storage/internal/server"
	"strconv"
)

type Config struct {
	grpcProtocol     string
	grpcPort         string
	dbUrl            string
	dbName           string
	collectionName   string
	httpPort         string
	logstashProtocol string
	logstashPort     int
}

func initConfig() (*Config, error) {
	grpcProtocol := os.Getenv("GRPC_PROTOCOL")
	grpcPort := os.Getenv("GRPC_PORT")
	dbUrl := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")
	httpPort := os.Getenv("HTTP_PORT")
	logstashProtocol := os.Getenv("LOGSTASH_PROTOCOL")
	logstashPort, err := strconv.Atoi(os.Getenv("LOGSTASH_PORT"))

	if err != nil {
		return nil, errors.New("LOGSTASH_PORT not set")
	}

	switch {
	case grpcProtocol == "":
		return nil, errors.New("GRPC_PROTOCOL not set")
	case grpcPort == "":
		return nil, errors.New("GRPC_PORT not set")
	case dbUrl == "":
		return nil, errors.New("DB_URL not set")
	case dbName == "":
		return nil, errors.New("DB_NAME not set")
	case collectionName == "":
		return nil, errors.New("COLLECTION_NAME not set")
	case httpPort == "":
		return nil, errors.New("HTTP_PORT not set")
	case logstashProtocol == "":
		return nil, errors.New("LOGSTASH_PROTOCOL not set")

	default:
	}
	cfg := &Config{
		grpcProtocol:     grpcProtocol,
		grpcPort:         grpcPort,
		dbUrl:            dbUrl,
		dbName:           dbName,
		collectionName:   collectionName,
		httpPort:         httpPort,
		logstashProtocol: logstashProtocol,
		logstashPort:     logstashPort,
	}
	return cfg, nil
}

func getDB(cfg *Config) (*repository.DataBase, error) {
	conn, err := repository.NewMongoConnection(cfg.dbUrl)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	db := repository.NewDataBase(conn, cfg.dbName, cfg.collectionName)

	return db, nil

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

	db, err := getDB(cfg)

	if err != nil {
		return
	}

	repo := repository.NewRepository(db)

	workerPool := infra.NewPool(repo, 16, 200, 10000000000)

	go func() {
		workerPool.Start()
	}()

	infraLayer := infra.NewInfra(repo, workerPool)

	listen, err := net.Listen(cfg.grpcProtocol, cfg.grpcPort)
	serv := grpc.NewServer()
	storageServer := server.NewServer(infraLayer, logger)

	log.Println("Created server")
	logger.Info(map[string]interface{}{
		"message": "Created server",
		"error":   false,
	})

	app := server.NewApp(repo)

	log.Println("Created app")
	logger.Info(map[string]interface{}{
		"message": "Created app",
		"error":   false,
	})

	go func() {
		if err := app.SetupApp(cfg.httpPort); err != nil {
			log.Println(err)
			logger.Error(map[string]interface{}{
				"message":           "Failed to setup app",
				"error":             true,
				"error_description": err.Error(),
			})
		}
	}()

	log.Println("Started server")
	logger.Info(map[string]interface{}{
		"message": "Started server",
		"error":   false,
	})

	pb.RegisterStorageServiceServer(serv, storageServer)

	log.Println("Registred server")
	logger.Info(map[string]interface{}{
		"message": "Registred server",
		"error":   false,
	})

	if err := serv.Serve(listen); err != nil {
		log.Println(err)
		logger.Error(map[string]interface{}{
			"message":           "Failed to start server",
			"error":             true,
			"error_description": err.Error(),
		})
		return
	}

	log.Println("Running server")
	logger.Info(map[string]interface{}{
		"message": "Running server",
		"error":   false,
	})
	log.Printf("ok")
}
