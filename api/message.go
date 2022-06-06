package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/rjriverac/messaging-server/db/sqlc"
)

type NewMessageReq struct {
	From    string `json:"from" binding:"required"`
	Content string `json:"content" binding:"required"`
	ConvID  int64  `json:"convID" binding:"required"`
	UserID  int64  `json:"from_id" binding:"required"`
}

func (s *Server) sendMessage(ctx *gin.Context) {
	var msgReq NewMessageReq
	if err := ctx.ShouldBindJSON(&msgReq); err != nil {
		ctx.JSON(http.StatusBadRequest,errorResponse(err))
		return
	}
	arg := db.SendMessageParams{
		CreateMessageParams: &db.CreateMessageParams{
			From: msgReq.From,
			Content: msgReq.Content,
			ConvID: msgReq.ConvID,
		},
		UserID: msgReq.UserID,
	}
	sent, err := s.store.SendMessage(ctx,arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError,errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted,sent)
}
