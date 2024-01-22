package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/justarandomlearner/WalletTransferAPI/internal/db"
	"github.com/justarandomlearner/WalletTransferAPI/internal/lib/errors"
	"github.com/justarandomlearner/WalletTransferAPI/internal/service"
)

type transferPayload struct {
	DebtorID      string  `json:"debtor_id"`
	BeneficiaryID string  `json:"beneficiary_id"`
	Amount        float64 `json:"amount"`
}

func TransferHandler(ctx *gin.Context) {
	conn, err := db.CreateConnection()
	if err != nil {
		// responseFromError(err, w)
		return
	}
	defer conn.Close()

	var tr transferPayload
	decoder := json.NewDecoder(ctx.Request.Body)
	if err := decoder.Decode(&tr); err != nil {
		ctx.Status(errors.ResponseFromError(err))
		return
	}

	if err := validateTransfer(tr); err != nil {
		ctx.Status(errors.ResponseFromError(err))
		return
	}

	transferService := service.NewTransferService()
	debtorID, err := uuid.Parse(tr.DebtorID)
	if err != nil {
		ctx.Status(errors.ResponseFromError(err))
		return
	}
	beneficiaryID, err := uuid.Parse(tr.BeneficiaryID)
	if err != nil {
		ctx.Status(errors.ResponseFromError(err))
		return
	}
	err = transferService.Transfer(
		tr.Amount,
		debtorID,
		beneficiaryID,
	)

	if err != nil {
		ctx.Status(errors.ResponseFromError(err))
		return
	}

	ctx.Status(http.StatusCreated)
}

func Bolotinha(ctx *gin.Context) {
	// var post post.Post

	// outerContext, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	// defer cancel()

	// if err := ctx.BindJSON(&post); err != nil {
	// 	ctx.JSON(http.StatusBadRequest, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

	// if err := service.Create(outerContext, post); err != nil {
	// 	ctx.JSON(http.StatusUnprocessableEntity, gin.H{
	// 		"error": err.Error(),
	// 	})
	// 	return
	// }

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
