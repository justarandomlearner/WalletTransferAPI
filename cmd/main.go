package main

import (
	"github.com/gin-gonic/gin"
	"github.com/justarandomlearner/WalletTransferAPI/internal/handlers"
)

func main() {
	g := gin.Default()

	g.GET("/accountinfo/:accountID", handlers.AccountBalance)

	g.POST("/transfer/", handlers.TransferHandler)

	g.Run(":3000")
}
