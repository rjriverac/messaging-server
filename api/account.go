package api

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "github.com/rjriverac/messaging-server/db/sqlc"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	HashedPw string `json:"hashedPw" binding:"required"`
}

func (server *Server) createUser(ctx *gin.Context) {
	var req CreateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		HashedPw: req.HashedPw,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
	}
	ctx.JSON(http.StatusOK, user)
}

type GetUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req GetUserRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	user, err := server.store.GetUser(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}

type ListUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

func (server *Server) listUser(ctx *gin.Context) {
	var req ListUserRequest

	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListUsersParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	users, err := server.store.ListUsers(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, users)
}

type ToBeNullString string

func (s *ToBeNullString) Scan(value interface{}) sql.NullString {
	var res sql.NullString
	i := *s
	if i == "" {
		res = sql.NullString{String: "", Valid: false}
	} else {
		val := fmt.Sprintf("%v", value)
		res = sql.NullString{String: val, Valid: true}
	}
	return res
}

type UpdateUserRequest struct {
	Name     ToBeNullString `json:"name"`
	Email    ToBeNullString `json:"email"`
	Image    ToBeNullString `json:"image"`
	Status   ToBeNullString `json:"status"`
	HashedPw ToBeNullString `json:"hashedPw"`
}

type UpdateUserID struct {
	ID int64 `form:"uid" binding:"required,min=1"`
}

func (server *Server) updateUser(g *gin.Context) {
	var req UpdateUserRequest
	var uid UpdateUserID
	if err := g.ShouldBindQuery(&uid); err != nil {
		g.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	argToParse := struct {
		Name     sql.NullString
		Email    sql.NullString
		Image    sql.NullString
		Status   sql.NullString
		HashedPw sql.NullString
	}{
		Name:     req.Name.Scan(req.Name),
		Email:    req.Email.Scan(req.Email),
		Image:    req.Image.Scan(req.Image),
		Status:   req.Status.Scan(req.Status),
		HashedPw: req.HashedPw.Scan(req.HashedPw),
	}

	arg := db.UpdateUserInfoParams{
		Name:     argToParse.Name,
		Email:    argToParse.Email,
		Image:    argToParse.Image,
		Status:   argToParse.Status,
		HashedPw: argToParse.HashedPw,
		ID:       uid.ID,
	}
	user, err := server.store.UpdateUserInfo(g, arg)
	if err != nil {
		g.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	g.JSON(http.StatusAccepted, user)

}
