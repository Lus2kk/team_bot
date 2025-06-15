package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"os/signal"
	"syscall"

	"team_bot/config"
	"team_bot/internal/handler"
	"team_bot/internal/repository/sqlrepo"
	"team_bot/internal/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	_ "github.com/lib/pq"
)

func main() {

	cfg, err := config.MustLoadConfig("config/config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	db, err := sql.Open("postgres", cfg.GetDatabaseConnectionString())
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer db.Close()

	logDB, err := sql.Open("postgres", cfg.GetLogDatabaseConnectionString())
	if err != nil {
		log.Fatalf("Error opening log database: %v", err)
	}
	defer logDB.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Error connecting to main database: %v", err)
	}
	if err := logDB.Ping(); err != nil {
		log.Fatalf("Error connecting to log database: %v", err)
	}

	authRepo := sqlrepo.NewAuthRepository(db)
	logRepo := sqlrepo.NewLogRepository(logDB)

	logService := service.NewLogService(logRepo)

	botAPI, err := tgbotapi.NewBotAPI(cfg.Bot.Token)
	if err != nil {
		log.Fatalf("Error creating bot: %v", err)
	}
	botAPI.Debug = cfg.Bot.Debug

	authHandler := handler.NewAuthHandler(botAPI, authRepo, cfg.TelegramAdmins.Usernames, logService)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		log.Println("Received shutdown signal")
		cancel()
	}()

	authHandler.Start(ctx)
}
