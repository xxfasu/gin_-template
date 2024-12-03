package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {
}

func newHttpServer(
	router *gin.Engine,
) *http.Server {

	return &http.Server{
		Addr:    ":" + "1234",
		Handler: router,
	}
}
