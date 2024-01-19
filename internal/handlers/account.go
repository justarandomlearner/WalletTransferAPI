package account

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AccountBalance(ctx *gin.Context) {
	accUUID, _ := uuid.Parse(ctx.Param("accUUID"))
	fmt.Println(accUUID)

	ctx.JSON(
		http.StatusOK, map[string]any{fmt.Sprintf("Account %v", accUUID): "+1000000 !!!!"},
	)

}
