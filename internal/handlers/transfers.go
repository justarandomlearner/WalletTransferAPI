package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/justarandomlearner/WalletTransferAPI/internal/db"
	"github.com/justarandomlearner/WalletTransferAPI/internal/lib/errors"
	"github.com/justarandomlearner/WalletTransferAPI/internal/repository"
	"github.com/justarandomlearner/WalletTransferAPI/internal/service"
)

type transferPayload struct {
	DebtorID      string  `json:"debtor_id"`
	BeneficiaryID string  `json:"beneficiary_id"`
	Amount        float64 `json:"amount"`
}

func TransferHandler(ctx *gin.Context) {
	var tr transferPayload

	decoder := json.NewDecoder(ctx.Request.Body)
	if err := decoder.Decode(&tr); err != nil {
		ctx.JSON(errors.ResponseFromError(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	if err := validateTransfer(tr); err != nil {
		ctx.JSON(errors.ResponseFromError(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	beneficiaryID, err := uuid.Parse(tr.BeneficiaryID)
	if err != nil {
		ctx.JSON(errors.ResponseFromError(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	conn, err := db.CreateConnection()
	if err != nil {
		ctx.Status(errors.ResponseFromError(err))
		return
	}

	defer conn.Close()

	repo := &repository.PostgresRepository{Conn: conn}

	transferService := service.NewTransferService(repo)

	debtorID, err := uuid.Parse(tr.DebtorID)
	if err != nil {
		ctx.JSON(errors.ResponseFromError(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	err = transferService.Transfer(
		tr.Amount,
		debtorID,
		beneficiaryID,
	)
	fmt.Printf("ERRO ENVIADO POR transferService.Transfer É: %v\n", err)
	if err != nil {
		fmt.Println("ERRO ENVIADO POR transferService.Transfer NÃO É NIL")
		ctx.JSON(errors.ResponseFromError(err), gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "transfer succeeded",
	})
}

func validateTransfer(t transferPayload) error {
	if t.Amount <= 0 {
		return errors.ErrCodeInvalidAmountToTransfer
	}

	if t.BeneficiaryID == t.DebtorID {
		return errors.ErrCodeSameDebtorAndBeneficiary
	}

	if t.BeneficiaryID == "" || t.DebtorID == "" {
		return errors.ErrCodeMissingPart
	}

	return nil
}
