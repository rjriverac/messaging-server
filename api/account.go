package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	db "github.com/rjriverac/messaging-server/db/sqlc"
)

type CreateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	HashedPw string `json:"hashedPw" binding:"required"`
	// Image    string `json:"image"`
	// Status   string `json:"status"`
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

// type NullString sql.NullString

// func (x *NullString) MarshalJSON() ([]byte, error) {
// 	if !x.Valid {
// 		x.Valid = true
// 		x.String = ""
// 	}
// 	return json.Marshal(x.String)
// }

// func (ns *NullString) Scan(value interface{}) error {
// 	var s sql.NullString
// 	if err := s.Scan(value); err != nil {
// 		return err
// 	}
// 	if reflect.TypeOf(value) == nil {
// 		*ns = NullString{s.String, false}
// 	} else {
// 		*ns = NullString{s.String, true}
// 	}
// 	return nil
// }

type ToBeNullString string

func (s *ToBeNullString) Scan(value interface{}) sql.NullString {
	var res sql.NullString
	// if _, err := s.Scan(value); err != nil {
	// 	return res, err
	// }
	i := reflect.ValueOf(value)
	if i.Elem().Interface() == "" {
		res = sql.NullString{"", false}
	} else {
		val := fmt.Sprintf("%v", value)
		res = sql.NullString{string(val), true}
	}
	return res
}

type UpdateUserRequest struct {
	ID       *int64          `uri:"id" binding:"min=1"`
	Name     ToBeNullString `json:"name"`
	Email    ToBeNullString `json:"email"`
	Image    ToBeNullString `json:"image"`
	Status   ToBeNullString `json:"status"`
	HashedPw ToBeNullString `json:"hashedPw"`
}

func (server *Server) updateUser(g *gin.Context) {
	var req UpdateUserRequest
	fmt.Printf("%v",req)
	if err := g.ShouldBindQuery(&req); err != nil {
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
		ID:       *req.ID,
	}
	user, err := server.store.UpdateUserInfo(g, arg)
	if err != nil {
		g.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	g.JSON(http.StatusAccepted, user)

}
