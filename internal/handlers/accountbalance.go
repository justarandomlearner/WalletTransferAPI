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

func AccountBalance(ctx *gin.Context) {
	accIdStr := ctx.Param("accountID")

	accUUID, err := uuid.Parse(accIdStr)
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

	service := service.NewAccountBalanceService(repo)

	balance, err := service.AccountBalance(accUUID)
	if err != nil {
		ctx.Status(errors.ResponseFromError(err))
		return
	}

	ctx.JSON(http.StatusOK, balance)
}
