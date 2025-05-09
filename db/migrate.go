package db

import (
	"fmt"
	"log"
	"time"

	"github.com/HermanSetiawan77777/wallet-live-coding-herman/config"
	transactionModel "github.com/HermanSetiawan77777/wallet-live-coding-herman/internal/transaction/model"
	userModel "github.com/HermanSetiawan77777/wallet-live-coding-herman/internal/user/model"
	walletModel "github.com/HermanSetiawan77777/wallet-live-coding-herman/internal/wallet/model"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB(appConfig *config.AppConfig) (*gorm.DB, error) {
	// Create DSN string for GORM with proper formatting
	dsn := fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s",
		appConfig.DB.User,
		appConfig.DB.Password,
		appConfig.DB.Host,
		appConfig.DB.Port,
		appConfig.DB.DBName,
	)

	// Initialize GORM with updated config
	db, err := gorm.Open(sqlserver.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQL Server: %v", err)
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	// Check if tables exist
	if !tablesExist(db) {
		log.Println("Tables not found, running migrations...")
		if err := runMigrations(db); err != nil {
			return nil, err
		}
	} else {
		log.Println("Tables already exist, skipping migrations")
	}

	return db, nil
}

func tablesExist(db *gorm.DB) bool {
	tables := []string{"users", "wallets", "transactions"}
	for _, table := range tables {
		if !db.Migrator().HasTable(table) {
			return false
		}
	}
	return true
}

func runMigrations(db *gorm.DB) error {
	models := []interface{}{
		&userModel.User{},
		&walletModel.Wallet{},
		&transactionModel.Transaction{},
	}

	for _, model := range models {
		if err := db.AutoMigrate(model); err != nil {
			return fmt.Errorf("failed to migrate %T: %v", model, err)
		}
	}

	// Only seed data if no existing data is found
	if !hasExistingData(db) {
		if err := seedDummyData(db); err != nil {
			return fmt.Errorf("failed to seed data: %v", err)
		}
		log.Println("Database migration and seeding completed successfully")
	} else {
		log.Println("Database migration completed, skipping seeding as data already exists")
	}

	return nil
}

func seedDummyData(db *gorm.DB) error {
	// Begin transaction for all seeding operations
	return db.Transaction(func(tx *gorm.DB) error {
		// Create dummy users
		users := []userModel.User{
			{ID: 1, Username: "john_doe", Email: "john@example.com"},
			{ID: 2, Username: "jane_doe", Email: "jane@example.com"},
			{ID: 3, Username: "bob_smith", Email: "bob@example.com"},
		}

		if err := tx.Create(&users).Error; err != nil {
			return fmt.Errorf("failed to seed users: %v", err)
		}

		// Create wallets for users
		wallets := []walletModel.Wallet{
			{UserID: 1, Balance: 1000000}, // 1M initial balance
			{UserID: 2, Balance: 500000},  // 500K initial balance
			{UserID: 3, Balance: 750000},  // 750K initial balance
		}

		if err := tx.Create(&wallets).Error; err != nil {
			return fmt.Errorf("failed to seed wallets: %v", err)
		}

		// Create some sample transactions
		now := time.Now()
		transactions := []transactionModel.Transaction{
			{UserID: 1, Type: "deposit", Amount: 1000000, CreatedAt: now.Add(-24 * time.Hour)},
			{UserID: 2, Type: "deposit", Amount: 500000, CreatedAt: now.Add(-12 * time.Hour)},
			{UserID: 3, Type: "deposit", Amount: 1000000, CreatedAt: now.Add(-6 * time.Hour)},
			{UserID: 3, Type: "withdraw", Amount: 250000, CreatedAt: now.Add(-1 * time.Hour)},
		}

		if err := tx.Create(&transactions).Error; err != nil {
			return fmt.Errorf("failed to seed transactions: %v", err)
		}

		log.Println("Successfully seeded dummy data")
		return nil
	})
}

// Add helper function to check if data already exists
func hasExistingData(db *gorm.DB) bool {
	var count int64
	db.Model(&userModel.User{}).Count(&count)
	return count > 0
}