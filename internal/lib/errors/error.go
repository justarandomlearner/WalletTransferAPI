package errors

import (
	"log"
	"net/http"
)

type Error struct {
	Code    int
	Message string
	err     error
}

func (e Error) Error() string {
	return e.Message
}

func New(code int, message string, err error) Error {
	return Error{
		Code:    code,
		Message: message,
		err:     err,
	}
}

var ErrCodeInvalidAmountToTransfer = New(
	CodeInvalidAmountToTransfer,
	"invalid amount to transfer, expected greater than zero",
	nil,
)

var ErrCodeSameDebtorAndBeneficiary = New(
	CodeSameDebtorAndBeneficiary,
	"cannot transfer to its own wallet",
	nil,
)

var ErrCodeMissingPart = New(
	CodeMissingPart,
	"missing beneficiary or debtor",
	nil,
)

func ResponseFromError(err error) int {
	e, ok := err.(Error)
	if !ok {
		log.Println(err)
		return http.StatusInternalServerError
	}

	switch e.Code {
	case CodeInsufficientBalance:
		return http.StatusExpectationFailed
	case CodeSameDebtorAndBeneficiary:
		return http.StatusExpectationFailed
	case CodeInvalidAmountToTransfer:
		return http.StatusExpectationFailed
	case CodeMissingPart:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
