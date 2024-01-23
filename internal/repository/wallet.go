package repository

import (
	"context"
	"fmt"

	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/justarandomlearner/WalletTransferAPI/internal/model"
)

type WalletRepository interface {
	OpenTransaction() (cancelContext context.CancelFunc, err error)
	Commit() error
	Rollback() error
	SelectBalanceByWalletID(walletID uuid.UUID) (model.WalletBalance, error)
	RemoveFromBalanceByWalletID(amount float64, walletID uuid.UUID) error
	AddOnBalanceByWalletID(amount float64, walletID uuid.UUID) error
}

const transactionTimeout = 20 * time.Second

const defaultTimeout = 10 * time.Second

type PostgresRepository struct {
	Conn *pgxpool.Pool
	tx   pgx.Tx
	ctx  context.Context
}

func (repo *PostgresRepository) OpenTransaction() (cancelContext context.CancelFunc, err error) {
	ctx, cancelContext := context.WithTimeout(context.Background(), transactionTimeout)

	repo.ctx = ctx
	options := pgx.TxOptions{}
	repo.tx, err = repo.Conn.BeginTx(ctx, options)
	fmt.Println(err)

	return
}

func (repo *PostgresRepository) Commit() error {
	if repo.tx.Conn().IsClosed() {
		return errors.New("database transaction is closed already")
	}

	return repo.tx.Commit(repo.ctx)
}

func (repo *PostgresRepository) Rollback() error {
	if repo.tx.Conn().IsClosed() {
		return errors.New("database transaction is closed already")
	}

	return repo.tx.Rollback(repo.ctx)
}

func (repo *PostgresRepository) SelectBalanceByWalletID(walletID uuid.UUID) (model.WalletBalance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	query := "SELECT id, amount, user_id FROM wallet WHERE id = $1"

	row := repo.Conn.QueryRow(ctx, query, walletID)

	var balance model.WalletBalance
	if err := row.Scan(
		&balance.ID,
		&balance.Amount,
		&balance.UserID,
	); err != nil {
		return model.WalletBalance{}, err
	}

	return balance, nil
}

func (repo *PostgresRepository) RemoveFromBalanceByWalletID(amount float64, walletID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(repo.ctx, defaultTimeout)
	defer cancel()

	sql := "UPDATE wallet SET amount = amount - $1 WHERE id = $2"

	_, err := repo.tx.Exec(ctx, sql, amount, walletID)

	return err
}

func (repo *PostgresRepository) AddOnBalanceByWalletID(amount float64, walletID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(repo.ctx, defaultTimeout)
	defer cancel()

	sql := "UPDATE wallet SET amount = amount + $1 WHERE id = $2"

	_, err := repo.tx.Exec(ctx, sql, amount, walletID)

	return err
}
