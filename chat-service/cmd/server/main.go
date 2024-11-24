package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gorilla/mux"
    "go.uber.org/zap"
    httpSwagger "github.com/swaggo/http-swagger"
    _ "chat-service/docs"

    "chat-service/config"
    "chat-service/internal/handler"
    "chat-service/internal/repository/mysql"
    "chat-service/internal/repository/redis"
    "chat-service/internal/service"
    "chat-service/pkg/database"
    "chat-service/pkg/elasticsearch"
    "chat-service/pkg/rabbitmq"
)

// @title Chat Service API
// @version 1.0
// @description A service for managing chats and messages
// @host localhost:8080
// @BasePath /
func main() {

    logger, err := zap.NewProduction()
    if err != nil {
        log.Fatalf("Failed to create logger: %v", err)
    }
    defer logger.Sync()


    cfg, err := config.Load()
    if err != nil {
        logger.Fatal("Error loading configuration", zap.Error(err))
    }

    db, err := database.NewMySQLConnection(database.MySQLConfig{
        Host:     cfg.MySQL.Host,
        Port:     cfg.MySQL.Port,
        User:     cfg.MySQL.User,
        Password: cfg.MySQL.Password,
        Database: cfg.MySQL.Database,
    })
    if err != nil {
        logger.Fatal("Failed to connect to MySQL", zap.Error(err))
    }
    defer db.Close()

    redisClient, err := database.NewRedisConnection(database.RedisConfig{
        Host: cfg.Redis.Host,
        Port: cfg.Redis.Port,
    })
    if err != nil {
        logger.Fatal("Failed to connect to Redis", zap.Error(err))
    }
    defer redisClient.Close()

    rabbitMQ, err := rabbitmq.NewClient(rabbitmq.Config{
        Host:     cfg.RabbitMQ.Host,
        Port:     cfg.RabbitMQ.Port,
        User:     cfg.RabbitMQ.User,
        Password: cfg.RabbitMQ.Password,
    })
    if err != nil {
        logger.Fatal("Failed to connect to RabbitMQ", zap.Error(err))
    }
    defer rabbitMQ.Close()

    esClient, err := elasticsearch.NewClient(elasticsearch.Config{
        URL: cfg.Elasticsearch.URL,
    })
    if err != nil {
        logger.Fatal("Failed to connect to Elasticsearch", zap.Error(err))
    }

    chatRepo := mysql.NewChatRepository(db)
    messageRepo := mysql.NewMessageRepository(db, esClient)
    sequenceRepo := redis.NewSequenceRepository(redisClient)

    chatService := service.NewChatService(
        chatRepo,
        sequenceRepo,
        rabbitMQ,
        logger,
    )
    
    messageService := service.NewMessageService(
        messageRepo,
        chatRepo,
        sequenceRepo,
        rabbitMQ,
        esClient,
        logger,
    )

    chatHandler := handler.NewChatHandler(chatService, logger)
    messageHandler := handler.NewMessageHandler(messageService, logger)

    router := mux.NewRouter()
    
    router.HandleFunc("/applications/{token}/chats", chatHandler.Create).Methods("POST")
    router.HandleFunc("/applications/{token}/chats/{number}/messages", messageHandler.Create).Methods("POST")
    router.HandleFunc("/applications/{token}/chats/{number}/messages", messageHandler.List).Methods("GET")
    router.HandleFunc("/applications/{token}/chats/{number}/messages/search", messageHandler.Search).Methods("GET")
    router.HandleFunc("/applications/{token}/chats/", chatHandler.ListChats).Methods("GET")

    router.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
        httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
        httpSwagger.DeepLinking(true),
        httpSwagger.DocExpansion("none"),
        httpSwagger.DomID("swagger-ui"),
    ))

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      router,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }

    go func() {
        logger.Info("Starting server on :8080")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            logger.Fatal("Server failed", zap.Error(err))
        }
    }()

    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
    <-stop

    logger.Info("Shutting down gracefully...")
    
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Fatal("Failed to gracefully shutdown server", zap.Error(err))
    }

    logger.Info("Server stopped")
}
