package service

import (
	"fmt"

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
	fmt.Printf("ERRO RECEBIDO DE decreaseWalletBalance: %v\n", err)
	if err != nil {
		fmt.Println("ENTROU NO IF PRA VER SE decreaseWalletBalance GEROU UM ERRO NÃO NULO")
		s.Repository.Rollback()
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

	result := debtorBalance.Amount - amount
	fmt.Printf("---->>>   debtorBalance.Amount(%g) - amount(%g) = result(%g)\n", debtorBalance.Amount, amount, result)
	if result < 0 {
		fmt.Println("ENTROU NO LOOP RESULT < 0")
		return errors.ErrCodeInsufficientBalance
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
