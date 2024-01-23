package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/justarandomlearner/WalletTransferAPI/internal/db"
	"github.com/justarandomlearner/WalletTransferAPI/internal/lib/errors"
	"github.com/justarandomlearner/WalletTransferAPI/internal/repository"
	"github.com/justarandomlearner/WalletTransferAPI/internal/service"
)

func WalletBalance(ctx *gin.Context) {
	walletIdStr := ctx.Param("walletID")

	walletUUID, err := uuid.Parse(walletIdStr)
	if err != nil {
		ctx.Status(http.StatusBadRequest)
		return
	}

	conn, err := db.CreateConnection()
	if err != nil {
		ctx.Status(errors.ResponseFromError(err))
		return
	}
	defer conn.Close()

	repo := &repository.PostgresRepository{Conn: conn}

	service := service.NewWalletBalanceService(repo)

	balance, err := service.WalletBalance(walletUUID)
	if err != nil {
		ctx.Status(errors.ResponseFromError(err))
		return
	}

	ctx.JSON(http.StatusOK, balance)
}
