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

type AccountRepository interface {
	OpenTransaction() (cancelContext context.CancelFunc, err error)
	Commit() error
	Rollback() error
	SelectBalanceByUserID(userID uuid.UUID) (model.AccountBalance, error)
	RemoveFromBalanceByUserID(amount float64, userID uuid.UUID) error
	AddOnBalanceByUserID(amount float64, userID uuid.UUID) error
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

func (repo *PostgresRepository) SelectBalanceByUserID(accID uuid.UUID) (model.AccountBalance, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
	defer cancel()

	sql := "SELECT id, amount, account_id FROM accounts WHERE account_id = $1"

	row := repo.Conn.QueryRow(ctx, sql, accID)

	var balance model.AccountBalance
	if err := row.Scan(
		&balance.ID,
		&balance.Amount,
		&balance.AccountID,
	); err != nil {
		return model.AccountBalance{}, err
	}

	return balance, nil
}

func (repo *PostgresRepository) RemoveFromBalanceByUserID(amount float64, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(repo.ctx, defaultTimeout)
	defer cancel()

	sql := "UPDATE accounts SET amount = amount - $1 WHERE account_id = $2"

	_, err := repo.tx.Exec(ctx, sql, amount, userID)

	return err
}

func (repo *PostgresRepository) AddOnBalanceByUserID(amount float64, userID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(repo.ctx, defaultTimeout)
	defer cancel()

	sql := "UPDATE accounts SET amount = amount + $1 WHERE account_id = $2"

	_, err := repo.tx.Exec(ctx, sql, amount, userID)

	return err
}
