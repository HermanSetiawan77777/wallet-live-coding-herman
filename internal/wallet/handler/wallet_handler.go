package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/HermanSetiawan77777/wallet-live-coding-herman/internal/wallet/repository"
	"github.com/HermanSetiawan77777/wallet-live-coding-herman/internal/wallet/service"
	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	svc service.WalletService
}

func NewWalletHandler(svc service.WalletService) *WalletHandler {
	return &WalletHandler{svc}
}

func (h *WalletHandler) Withdraw(c *gin.Context) {
	var req struct {
		UserID int `json:"user_id" binding:"required,min=1"`
		Amount int `json:"amount" binding:"required,min=1"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request parameters"})
		return
	}

	err := h.svc.Withdraw(req.UserID, req.Amount)
	if err != nil {
		switch {
		case errors.Is(err, repository.ErrInsufficientBalance):
			c.JSON(http.StatusBadRequest, gin.H{"error": "insufficient balance"})
		case errors.Is(err, repository.ErrWalletNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id format"})
		return
	}

	balance, err := h.svc.GetBalance(userID)
	if err != nil {
		if errors.Is(err, repository.ErrWalletNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}