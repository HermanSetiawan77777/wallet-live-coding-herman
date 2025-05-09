package service

import "github.com/HermanSetiawan77777/wallet-live-coding-herman/internal/wallet/repository"

type WalletService interface {
	Withdraw(userID int, amount int) error
	GetBalance(userID int) (int, error)
}

type walletService struct {
	repo repository.WalletRepository
}

func NewWalletService(r repository.WalletRepository) WalletService {
	return &walletService{r}
}

func (s *walletService) Withdraw(userID int, amount int) error {
	return s.repo.Withdraw(userID, amount)
}

func (s *walletService) GetBalance(userID int) (int, error) {
	return s.repo.GetBalance(userID)
}