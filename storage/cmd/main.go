package main

import (
	"errors"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
	"storage/internal/infra"
	pb "storage/internal/proto"
	"storage/internal/repository"
	"storage/internal/server"
)

type Config struct {
	grpcProtocol      string
	grpcPort          string
	dbUrl             string
	dbName            string
	collectionName    string
	mongoRootUserName string
	mongoRootPassword string
	httpPort          string
}

func initConfig() (*Config, error) {
	grpcProtocol := os.Getenv("GRPC_PROTOCOL")
	grpcPort := os.Getenv("GRPC_PORT")
	dbUrl := os.Getenv("DB_URL")
	dbName := os.Getenv("DB_NAME")
	collectionName := os.Getenv("COLLECTION_NAME")
	mongoRootUserName := os.Getenv("MONGO_INITDB_ROOT_USERNAME")
	mongoRootPassword := os.Getenv("MONGO_INITDB_ROOT_PASSWORD")
	httpPort := os.Getenv("HTTP_PORT")

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
	case mongoRootUserName == "":
		return nil, errors.New("MONGO_INITDB_ROOT_USERNAME not set")
	case mongoRootPassword == "":
		return nil, errors.New("MONGO_INITDB_ROOT_PASSWORD not set")
	case httpPort == "":
		return nil, errors.New("HTTP_PORT not set")
	default:
	}
	cfg := &Config{
		grpcProtocol:      grpcProtocol,
		grpcPort:          grpcPort,
		dbUrl:             dbUrl,
		dbName:            dbName,
		collectionName:    collectionName,
		mongoRootUserName: mongoRootUserName,
		mongoRootPassword: mongoRootPassword,
		httpPort:          httpPort,
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
	storageServer := server.NewServer(infraLayer)
	log.Printf("created server")
	app := server.NewApp(repo)
	log.Printf("created app")
	go func() {
		if err := app.SetupApp(cfg.httpPort); err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	log.Printf("started server")
	pb.RegisterStorageServiceServer(serv, storageServer)
	log.Printf("registred server")
	if err := serv.Serve(listen); err != nil {
		log.Printf("failed to start server: %v", err)
		return
	}
	log.Println("server running")
}
