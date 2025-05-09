package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/HermanSetiawan77777/wallet-live-coding-herman/config"
	dbMigrate "github.com/HermanSetiawan77777/wallet-live-coding-herman/db"
	"github.com/HermanSetiawan77777/wallet-live-coding-herman/internal/wallet/handler"
	"github.com/HermanSetiawan77777/wallet-live-coding-herman/internal/wallet/repository"
	"github.com/HermanSetiawan77777/wallet-live-coding-herman/internal/wallet/service"
	"github.com/HermanSetiawan77777/wallet-live-coding-herman/routes"
	"github.com/HermanSetiawan77777/wallet-live-coding-herman/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Application error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	// Load environment variables
	envPath := utils.GetAppRootDirectory() + "/.env"
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Warning: .env file not found at %s, using system environment variables", envPath)
	}

	// Load application configuration
	appConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Initialize database with GORM
	db, err := dbMigrate.InitDB(appConfig)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	// Get the underlying *sql.DB instance
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get database instance: %v", err)
	}
	defer sqlDB.Close()

	log.Printf("Successfully connected to the database at %s\n", appConfig.DB.Host)
	log.Printf("Server starting on port %d in %s mode\n", appConfig.Port, appConfig.Environment)

	// Set Gin mode based on environment
	if appConfig.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Setup repositories and services
	walletRepo := repository.NewWalletRepository(db)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := handler.NewWalletHandler(walletService)

	// Setup router
	router := routes.SetupRouter(walletHandler)

	// Create server with Gin router
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", appConfig.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, cancelShutdown := context.WithTimeout(serverCtx, 30*time.Second)
		defer cancelShutdown() // Add this to prevent context leak

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Printf("Error shutting down server: %v", err)
		}
		serverStopCtx()
	}()

	// Run the server
	log.Printf("Server is running on port %d", appConfig.Port)
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("error starting server: %v", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	return nil
}