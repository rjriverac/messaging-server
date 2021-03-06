package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/rjriverac/messaging-server/db/sqlc"
	"github.com/rjriverac/messaging-server/token"
)

type NullString sql.NullString

func (s *NullString) MarshalNullStr() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return []byte(""), nil
}
func (s *NullString) NullStrToString() string {
	if s.Valid {
		return s.String
	}
	return ""
}

type ConversationReturn struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

func (server *Server) getConvos(g *gin.Context) {
	var ret []ConversationReturn

	authPayload := g.MustGet(authPayloadKey).(*token.Payload)

	convs, err := server.store.ListConvFromUser(context.Background(), authPayload.User)
	if err != nil {
		if err == sql.ErrNoRows {
			g.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		g.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	for _, conv := range convs {
		nstr := NullString(conv.Name)
		str := nstr.NullStrToString()
		ret = append(ret, ConversationReturn{ID: conv.ID, Name: str})
	}

	g.JSON(http.StatusOK, ret)

}

type getConvDetailRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) detailConvo(g *gin.Context) {
	var req getConvDetailRequest
	if err := g.ShouldBindUri(&req); err != nil {
		g.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	auth := g.MustGet(authPayloadKey).(*token.Payload)

	arg := db.ListConvMessagesParams{
		ConvID: req.ID,
		UserID: auth.User,
	}

	messages, err := server.store.ListConvMessages(context.Background(), arg)
	if err != nil {
		if err == sql.ErrNoRows {
			g.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		g.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	g.JSON(http.StatusOK, messages)
}

type createConvRequest struct {
	Name    ToBeNullString `json:"conv_name"`
	ToUsers []string       `json:"recipient_emails" binding:"required,gt=0,dive,required,email"`
	From    int64          `json:"from" binding:"required,min=1"`
}

//
func (server *Server) createConvo(g *gin.Context) {

	var req createConvRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	auth := g.MustGet(authPayloadKey).(*token.Payload)

	arg := db.CreateConvParams{
		Name:    req.Name.ToNstring(),
		ToUsers: req.ToUsers,
		From:    auth.User,
	}

	conv, err := server.store.CreateConvTx(context.Background(), arg)
	if err != nil {
		g.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ret := ConversationReturn{
		ID:   conv.ID,
		Name: conv.Name,
	}

	g.JSON(http.StatusAccepted, ret)
}
