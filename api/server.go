package api

import (
	"github.com/gin-gonic/gin"
	db "github.com/rjriverac/messaging-server/db/sqlc"
)

type Server struct {
	store  db.Store
	router *gin.Engine
}

func NewServer(store db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/account", server.createUser)
	router.GET("/account/:id", server.getUser)
	router.GET("/account", server.listUser)
	router.PUT("/account/", server.updateUser)

	router.POST("/message",server.sendMessage)

	server.router = router
	return server
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) StartServer(addr string) error {
	return server.router.Run(addr)
}
