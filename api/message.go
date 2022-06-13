package api

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/rjriverac/messaging-server/db/sqlc"
	"github.com/rjriverac/messaging-server/token"
)

type NewMessageReq struct {
	// From    string `json:"from" binding:"required"`
	Content string `json:"content" binding:"required"`
	ConvID  int64  `json:"convID" binding:"min=1"`
	// UserID  int64  `json:"from_id" binding:"required,min=1"`
}

func (s *Server) sendMessage(ctx *gin.Context) {
	var msgReq NewMessageReq
	if err := ctx.ShouldBindJSON(&msgReq); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	auth := ctx.MustGet(authPayloadKey).(*token.Payload)

	user, err := s.store.GetUser(ctx, auth.User)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.SendMessageParams{
		CreateMessageParams: db.CreateMessageParams{
			From:    user.Name,
			Content: msgReq.Content,
			ConvID:  msgReq.ConvID,
		},
		UserID: auth.User,
	}
	sent, err := s.store.SendMessage(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusAccepted, sent)
}
