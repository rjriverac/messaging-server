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

func MarshalNullStr(str *NullString) ([]byte, error) {
	if str.Valid {
		return json.Marshal(str.String)
	} else {
		return json.Marshal(nil)
	}
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
		str, _ := MarshalNullStr((*NullString)(&conv.Name))
		ret = append(ret, ConversationReturn{ID: conv.ID, Name: str})
	}

	g.JSON(http.StatusOK, ret)

}
