package service

import (
	"github.com/google/uuid"
	"github.com/justarandomlearner/WalletTransferAPI/internal/db"
	"github.com/justarandomlearner/WalletTransferAPI/internal/lib/errors"
	"github.com/justarandomlearner/WalletTransferAPI/internal/repository"
)

type TransferService struct {
	Repository repository.AccountRepository
}

func NewTransferService() TransferService {
	conn, _ := db.CreateConnection()
	repo := repository.PostgresRepository{Conn: conn}
	return TransferService{Repository: &repo}
}

func (s *TransferService) Transfer(amount float64, debtorID, beneficiaryID uuid.UUID) error {
	cancel, err := s.Repository.OpenTransaction()
	defer cancel()
	if err != nil {
		return err
	}

	err = s.removeFromBalance(amount, debtorID)
	if err != nil {
		return err
	}

	err = s.topUpBalance(amount, beneficiaryID)
	if err != nil {
		s.Repository.Rollback()
		return err
	}

	s.Repository.Commit()
	return nil
}

func (s *TransferService) removeFromBalance(amount float64, userID uuid.UUID) error {
	debtorBalance, err := s.Repository.SelectBalanceByUserID(userID)
	if err != nil {
		return errors.New(errors.CodeInternalDatabaseError, "error on selecting balance", err)
	}

	if debtorBalance.Amount-amount < 0 {
		return errors.New(errors.CodeInsufficientBalance, "insufficient balance on debtor account", err)
	}

	err = s.Repository.RemoveFromBalanceByUserID(amount, userID)

	return err
}

func (s *TransferService) topUpBalance(amount float64, userID uuid.UUID) error {
	_, err := s.Repository.SelectBalanceByUserID(userID)
	if err != nil {
		return errors.New(errors.CodeInternalDatabaseError, "error on selecting balance", err)
	}

	err = s.Repository.AddOnBalanceByUserID(amount, userID)

	return err
}
