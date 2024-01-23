package service

import (
	"github.com/google/uuid"
	"github.com/justarandomlearner/WalletTransferAPI/internal/lib/errors"
	"github.com/justarandomlearner/WalletTransferAPI/internal/repository"
)

type transferService struct {
	Repository repository.WalletRepository
}

func NewTransferService(repo repository.WalletRepository) transferService {
	return transferService{Repository: repo}
}

func (s *transferService) Transfer(amount float64, debtorID, beneficiaryID uuid.UUID) error {
	cancel, err := s.Repository.OpenTransaction()
	defer cancel()
	if err != nil {
		return err
	}

	err = s.decreaseWalletBalance(amount, debtorID)
	if err != nil {
		return err
	}

	err = s.increaseWalletBalance(amount, beneficiaryID)
	if err != nil {
		s.Repository.Rollback()
		return err
	}

	s.Repository.Commit()
	return nil
}

func (s *transferService) decreaseWalletBalance(amount float64, walletID uuid.UUID) error {
	debtorBalance, err := s.Repository.SelectBalanceByWalletID(walletID)
	if err != nil {
		return errors.New(errors.CodeInternalDatabaseError, "error on selecting balance", err)
	}

	if debtorBalance.Amount-amount < 0 {
		return errors.New(errors.CodeInsufficientBalance, "insufficient balance on debtor wallet", err)
	}

	err = s.Repository.RemoveFromBalanceByWalletID(amount, walletID)

	return err
}

func (s *transferService) increaseWalletBalance(amount float64, walletID uuid.UUID) error {
	_, err := s.Repository.SelectBalanceByWalletID(walletID)
	if err != nil {
		return errors.New(errors.CodeInternalDatabaseError, "error on selecting balance", err)
	}

	err = s.Repository.AddOnBalanceByWalletID(amount, walletID)

	return err
}
