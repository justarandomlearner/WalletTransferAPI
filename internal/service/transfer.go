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

	err = s.decreaseAccountBalance(amount, debtorID)
	if err != nil {
		return err
	}

	err = s.increaseAccountBalance(amount, beneficiaryID)
	if err != nil {
		s.Repository.Rollback()
		return err
	}

	s.Repository.Commit()
	return nil
}

func (s *TransferService) decreaseAccountBalance(amount float64, accountID uuid.UUID) error {
	debtorBalance, err := s.Repository.SelectBalanceByAccountID(accountID)
	if err != nil {
		return errors.New(errors.CodeInternalDatabaseError, "error on selecting balance", err)
	}

	if debtorBalance.Amount-amount < 0 {
		return errors.New(errors.CodeInsufficientBalance, "insufficient balance on debtor account", err)
	}

	err = s.Repository.RemoveFromBalanceByAccountID(amount, accountID)

	return err
}

func (s *TransferService) increaseAccountBalance(amount float64, accountID uuid.UUID) error {
	_, err := s.Repository.SelectBalanceByAccountID(accountID)
	if err != nil {
		return errors.New(errors.CodeInternalDatabaseError, "error on selecting balance", err)
	}

	err = s.Repository.AddOnBalanceByAccountID(amount, accountID)

	return err
}
