package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"learning.com/golang_backend/auth"
	db "learning.com/golang_backend/db/sqlc/repository"
	"learning.com/golang_backend/utils"
)

type Server struct {
	config       utils.Config
	repository   db.Repository
	tokenBuilder auth.Token
	router       *gin.Engine
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func NewServer(config utils.Config, repository db.Repository) (*Server, error) {
	tokenBuilder, err := auth.NewPasetoToken(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:       config,
		repository:   repository,
		tokenBuilder: tokenBuilder,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRouter()

	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	router.POST("/auth/signup", server.sigupUser)
	router.POST("/auth/signin", server.signInUser)

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccounts)

	router.POST("/transfers", server.createTransfer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
