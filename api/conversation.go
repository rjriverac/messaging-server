package api

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

type getConversationRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type NullString sql.NullString

func (s *NullString) MarshalNullStr() ([]byte, error) {
	if s.Valid {
		return json.Marshal(s.String)
	}
	return []byte(""), nil
}

type ConversationReturn struct {
	ID   int64  `json:"id"`
	Name []byte `json:"name"`
}

func (server *Server) getConvos(g *gin.Context) {
	var req getConversationRequest
	var ret []ConversationReturn
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	convs, err := server.store.ListConvFromUser(context.Background(), req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			g.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		g.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	for _, conv := range convs {
		nstr := NullString(conv.Name)
		str, _:= nstr.MarshalNullStr()
		ret = append(ret, ConversationReturn{ID: conv.ID, Name: str})
	}

	g.JSON(http.StatusOK, ret)

}
