package repository

import (
	"errors"

	transactionModel "github.com/HermanSetiawan77777/wallet-live-coding-herman/internal/transaction/model"
	"github.com/HermanSetiawan77777/wallet-live-coding-herman/internal/wallet/model"
	"gorm.io/gorm"
)

var (
	ErrInsufficientBalance = errors.New("insufficient balance")
	ErrWalletNotFound     = errors.New("wallet not found")
)

type WalletRepository interface {
    GetWalletByUserID(userID int) (*model.Wallet, error)
    Withdraw(userID int, amount int) error
    GetBalance(userID int) (int, error)
}

type walletRepo struct {
    db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
    return &walletRepo{db}
}

func (r *walletRepo) GetWalletByUserID(userID int) (*model.Wallet, error) {
    var wallet model.Wallet
    err := r.db.Where("user_id = ?", userID).First(&wallet).Error
    return &wallet, err
}

func (r *walletRepo) Withdraw(userID int, amount int) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        var wallet model.Wallet
        if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                return ErrWalletNotFound
            }
            return err
        }

        if wallet.Balance < amount {
            return ErrInsufficientBalance
        }

        result := tx.Model(&wallet).Update("balance", wallet.Balance-amount)
        if result.Error != nil {
            return result.Error
        }

        return tx.Create(&transactionModel.Transaction{
            UserID: userID,
            Type:   "withdraw",
            Amount: amount,
        }).Error
    })
}

func (r *walletRepo) GetBalance(userID int) (int, error) {
    var wallet model.Wallet
    err := r.db.Where("user_id = ?", userID).First(&wallet).Error
    return wallet.Balance, err
}
