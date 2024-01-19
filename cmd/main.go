package main

import (
	"os"

	"github.com/gin-gonic/gin"
	account "github.com/justarandomlearner/WalletTransferAPI/internal/handlers"
)

func main() {
	g := gin.Default()

	g.GET("/accountinfo/:accUUID", account.AccountBalance)

	port := os.Getenv("API_PORT")
	if port == "" {
		port = "3000"
	}

	g.Run(":" + port)
}
