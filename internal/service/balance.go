package service

import (
	"github.com/google/uuid"
	"github.com/justarandomlearner/WalletTransferAPI/internal/model"
	"github.com/justarandomlearner/WalletTransferAPI/internal/repository"
)

type AccountBalanceService struct {
	Repository repository.AccountRepository
}

func NewAccountBalanceService(repo repository.AccountRepository) AccountBalanceService {
	return AccountBalanceService{Repository: repo}
}

func (s *AccountBalanceService) AccountBalance(accountID uuid.UUID) (*model.AccountBalance, error) {
	balance, err := s.Repository.SelectBalanceByAccountID(accountID)
	if err != nil {
		return nil, err
	}

	return &balance, nil
}
