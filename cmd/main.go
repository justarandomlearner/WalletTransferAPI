package main

import (
	"github.com/gin-gonic/gin"
	account "github.com/justarandomlearner/WalletTransferAPI/internal/handlers"
)

func main() {
	g := gin.Default()

	g.GET("/accountinfo/:accUUID", account.AccountBalance)

	g.Run(":12345")
}
