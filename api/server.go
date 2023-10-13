package api

import (
	"fmt"

	"github.com/ftvdexcz/simplebank/config"
	db "github.com/ftvdexcz/simplebank/db/sqlc"
	"github.com/ftvdexcz/simplebank/token"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type Server struct {
	config config.Config
	store  *db.Store
	router *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config config.Config, store *db.Store) (*Server, error){
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil{
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	
	server := &Server{config: config, store: store, tokenMaker: tokenMaker}
	

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok{
		v.RegisterValidation("currency", validCurrency)
	}

	
	server.setupRoute()

	
	return server, nil
}

func (server *Server) setupRoute(){
	router := gin.Default()

	router.POST("/users", server.createUser)
	router.POST("/users/login", server.LoginUser)


	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker)) // use middleware

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.POST("/transfers", server.createTransfer)
	
	authRoutes.GET("/accounts/:id", server.getAccount)
	authRoutes.GET("/accounts", server.listAccount)

	server.router = router
}

func (server *Server) Start(address string) error{
	return server.router.Run(address)
}

func errorResponse(err error) gin.H{
	return gin.H{
		"error": err.Error(),
	}
}