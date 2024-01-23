package service

import (
	"github.com/google/uuid"
	"github.com/justarandomlearner/WalletTransferAPI/internal/model"
	"github.com/justarandomlearner/WalletTransferAPI/internal/repository"
)

type walletBalanceService struct {
	Repository repository.WalletRepository
}

func NewWalletBalanceService(repo repository.WalletRepository) walletBalanceService {
	return walletBalanceService{Repository: repo}
}

func (s *walletBalanceService) WalletBalance(walletID uuid.UUID) (*model.WalletBalance, error) {
	balance, err := s.Repository.SelectBalanceByWalletID(walletID)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}
