package api

import (
	"github.com/gin-gonic/gin"
	db "learning.com/golang_backend/db/sqlc/repository"
)

type Server struct {
	repository db.Repository
	router     *gin.Engine
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func NewServer(repository db.Repository) *Server {
	server := &Server{repository: repository}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	server.router = router
	return server
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
