package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/rjriverac/messaging-server/db/sqlc"
	"github.com/rjriverac/messaging-server/token"
	"github.com/rjriverac/messaging-server/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token:%w", err)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterStructValidation(validRequest, UpdateUserRequest{})
	}
	server.createRoutes()
	return server, nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) StartServer(addr string) error {
	return server.router.Run(addr)
}

func (server *Server) createRoutes() {
	router := gin.Default()
	router.POST("/account/", server.createUser)
	router.POST("/account/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMWare(server.tokenMaker))

	authRoutes.GET("/account/:id", server.getUser)
	authRoutes.GET("/account/", server.listUser)
	authRoutes.PUT("/account/", server.updateUser)

	authRoutes.POST("/message", server.sendMessage)

	authRoutes.GET("/conversation", server.getConvos)
	authRoutes.GET("/conversation/:id", server.detailConvo)
	authRoutes.POST("/conversation", server.createConvo)

	server.router = router
}
