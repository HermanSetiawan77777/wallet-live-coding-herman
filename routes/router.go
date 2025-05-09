package routes

import (
	"time"

	walletHandler "github.com/HermanSetiawan77777/wallet-live-coding-herman/internal/wallet/handler"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func SetupRouter(h *walletHandler.WalletHandler) *gin.Engine {
    r := gin.Default()

    // Add middleware
    r.Use(gin.Recovery())
    r.Use(cors.New(cors.Config{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
        ExposeHeaders:    []string{"Content-Length"},
        AllowCredentials: true,
        MaxAge:           12 * time.Hour,
    }))

    // Health check endpoint
    r.GET("/health", func(c *gin.Context) {
        c.JSON(200, gin.H{"status": "ok"})
    })

    // API routes
    api := r.Group("/api/v1")
    {
        wallet := api.Group("/wallet")
        {
            wallet.POST("/withdraw", h.Withdraw)
            wallet.GET("/balance", h.GetBalance)
        }
    }

    return r
}