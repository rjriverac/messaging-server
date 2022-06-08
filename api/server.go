package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/rjriverac/messaging-server/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	if v,ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterStructValidation(validRequest,UpdateUserRequest{})
	}

	router.POST("/account", server.createUser)
	router.GET("/account/:id", server.getUser)
	router.GET("/account", server.listUser)
	router.PUT("/account/", server.updateUser)

	router.POST("/message",server.sendMessage)

	router.GET("/conversation",server.getConvos)

	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) StartServer(addr string) error {
	return server.router.Run(addr)
}
