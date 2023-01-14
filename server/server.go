package server

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"poll-service/polls"
	"poll-service/polls/repository/mongodb"
	"poll-service/polls/usecase"
	logger "poll-service/utils/logger"
	"time"

	httpPoll "poll-service/polls/delivery/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// mongodb config keys
const (
	mongoHost           = "mongohost"
	mongoUsername       = "mongocred.username"
	mongoPwd            = "mongocred.password"
	mongoDB             = "mongodb"
	mongoCredAuthMech   = "mongocred.authmechanism"
	mongoCredAuthSource = "mongocred.authsource"
	mongoURI            = "mongodb://%s/?maxPoolSize=20&w=majority"
)

// logger config keys
const (
	pollsLog = "logs.pollslog_filename"
	logLevel = "logs.level"
	logPath  = "logs.path"
)

const defaultPerm = 0774

type App struct {
	httpSrv     *http.Server
	pollUsecase polls.UseCase
	logger      *logger.Logger
}

func NewApp() *App {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()

	mongoDB := initMongoDB(ctx)

	// loggers
	pollsLogger := logger.NewLogger(initLogrusLogger(viper.GetString(pollsLog)))

	// repos
	pollRepo := mongodb.NewPollRepo(mongoDB, pollsLogger)

	return &App{
		pollUsecase: usecase.NewPollUsecase(pollRepo, pollsLogger),
		logger:      pollsLogger,
	}
}

func (a *App) Run(port string) error {
	// Init gin handler
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Set up http handlers
	// Middleware
	middleware := httpPoll.NewPollsMiddleware(a.pollUsecase)

	// API endpoints
	polls := r.Group("/api", middleware)
	httpPoll.RegisterMidRoutes(polls, a.pollUsecase, a.logger)

	// HTTP Server
	a.httpSrv = &http.Server{
		Addr:           ":" + port,
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// go run $(go env GOROOT)/src/crypto/tls/generate_cert.go --host=localhost
	if err := a.httpSrv.ListenAndServe(); err != nil {
		// if err := a.httpSrv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
		log.Fatalf("Failed to listen and serve: %+v", err)
	}
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	return a.httpSrv.Shutdown(ctx)
}

func initLogrusLogger(fileName string) *logrus.Logger {
	logger := logrus.New()
	level, err := logrus.ParseLevel(viper.GetString(logLevel))
	if err != nil {
		log.Fatalf("error logger init: %v", err)
	}
	logger.SetLevel(level)

	path := viper.GetString(logPath)
	if err := os.MkdirAll(path, defaultPerm); err != nil {
		log.Fatalf("error creating log dir: %s; %v", path, err)
	}

	filePath := filepath.Join(path, fileName)
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, defaultPerm)
	if err != nil {
		log.Fatalf("error creating log dir: %v", err)
	}

	logger.SetOutput(io.MultiWriter(file, os.Stderr))

	return logger
}

func initMongoDB(ctx context.Context) *mongo.Database {
	uri := fmt.Sprintf(mongoURI, viper.GetString(mongoHost))

	log.Println(uri)

	clientCred := options.Credential{
		AuthMechanism: viper.GetString(mongoCredAuthMech),
		AuthSource:    viper.GetString(mongoCredAuthSource),
		Username:      viper.GetString(mongoUsername),
		Password:      viper.GetString(mongoPwd),
	}
	clientOptions := options.Client().ApplyURI(uri).SetAuth(clientCred)

	client, err := mongo.Connect(ctx, clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	// Ping the primary
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	log.Printf("Successfully connected to MongoDB: host:%s db:%s", viper.GetString(mongoHost), viper.GetString(mongoDB))
	return client.Database(viper.GetString(mongoDB))
}
