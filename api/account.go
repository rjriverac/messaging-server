package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "github.com/rjriverac/messaging-server/db/sqlc"
	"github.com/rjriverac/messaging-server/token"
	"github.com/rjriverac/messaging-server/util"
)

type createUserRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type userReturn struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Image     []byte    `json:"image"`
	Status    []byte    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func newUserReturn(user db.User) userReturn {
	istr := NullString(user.Image)
	str, _ := istr.MarshalNullStr()
	ststr := NullString(user.Status)
	ustr, _ := ststr.MarshalNullStr()

	return userReturn{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Image:     str,
		Status:    ustr,
		CreatedAt: user.CreatedAt,
	}
}

func (server *Server) createUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	hashedPw, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	arg := db.CreateUserParams{
		Name:     req.Name,
		Email:    req.Email,
		HashedPw: hashedPw,
	}

	user, err := server.store.CreateUser(ctx, arg)
	if err != nil {

		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ret := newUserReturn(user)

	ctx.JSON(http.StatusOK, ret)
}

type GetUserRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

type GetUserReturn struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Image     string    `json:"image"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (server *Server) getUser(ctx *gin.Context) {
	var req GetUserRequest

	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	auth := ctx.MustGet(authPayloadKey).(*token.Payload)
	if auth.User != req.ID {
		err := errors.New("access to requested information not allowed")
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
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
	ret := GetUserReturn{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}
	nullStrs := map[string]NullString{
		"Image":  NullString(user.Image),
		"Status": NullString(user.Status),
	}

	for key, nstring := range nullStrs {
		str := nstring.NullStrToString()
		switch key {
		case "Image":
			ret.Image = str
		case "Status":
			ret.Status = str
		}
	}

	ctx.JSON(http.StatusOK, ret)
}

type ListUserRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=20"`
}

type ListUserAcc struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Image  string `json:"image"`
	Status string `json:"status"`
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

	var listRet []ListUserAcc
	for _, user := range users {
		item := ListUserAcc{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
		nullStrs := map[string]NullString{
			"Image":  NullString(user.Image),
			"Status": NullString(user.Status),
		}
		for key, nstring := range nullStrs {
			str := nstring.NullStrToString()
			switch key {
			case "Image":
				item.Image = str
			case "Status":
				item.Status = str
			}
		}
		listRet = append(listRet, item)

	}
	ctx.JSON(http.StatusOK, listRet)
}

type ToBeNullString string

func (s *ToBeNullString) Scan(value interface{}) sql.NullString {
	var res sql.NullString
	i := *s
	if i == "" {
		res = sql.NullString{String: "", Valid: false}
	} else {
		val := fmt.Sprintf("%v", value)
		// val, _ := value.(string)
		res = sql.NullString{String: val, Valid: true}
	}
	return res
}
func (s *ToBeNullString) ToNstring() sql.NullString {
	if len(*s) == 0 {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: string(*s), Valid: true}
}

type UpdateUserRequest struct {
	Name     ToBeNullString `json:"name"`
	Email    ToBeNullString `json:"email"`
	Image    ToBeNullString `json:"image"`
	Status   ToBeNullString `json:"status"`
	Password string         `json:"password"`
}

type UpdateUserReturn struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Image     string    `json:"image"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"createdAt"`
}

func (server *Server) updateUser(g *gin.Context) {
	var req UpdateUserRequest

	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	auth := g.MustGet(authPayloadKey).(*token.Payload)

	var arg db.UpdateUserInfoParams
	if len(req.Password) == 0 {
		nPw := ToBeNullString(req.Password)
		arg = db.UpdateUserInfoParams{
			Name:     req.Name.Scan(req.Name),
			Email:    req.Email.Scan(req.Email),
			Image:    req.Image.Scan(req.Image),
			Status:   req.Status.Scan(req.Status),
			HashedPw: nPw.ToNstring(),
			ID:       auth.User,
		}
	} else {

		hashed, err := util.HashPassword(req.Password)
		if err != nil {
			g.JSON(http.StatusInternalServerError, errorResponse(err))
			return
		}

		nHash := ToBeNullString(hashed)

		arg = db.UpdateUserInfoParams{
			Name:     req.Name.Scan(req.Name),
			Email:    req.Email.Scan(req.Email),
			Image:    req.Image.Scan(req.Image),
			Status:   req.Status.Scan(req.Status),
			HashedPw: nHash.Scan(nHash),
			ID:       auth.User,
		}
	}

	user, err := server.store.UpdateUserInfo(g, arg)
	if err != nil {
		g.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	nullStrs := map[string]NullString{
		"Image":  NullString(user.Image),
		"Status": NullString(user.Status),
	}
	ret := UpdateUserReturn{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
	}

	for key, nstring := range nullStrs {
		str := nstring.NullStrToString()
		switch key {
		case "Image":
			ret.Image = str
		case "Status":
			ret.Status = str
		}
	}

	g.JSON(http.StatusAccepted, ret)

}

type loginUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type loginUserResponse struct {
	AccessToken string     `json:"access_token"`
	User        userReturn `json:"user"`
}

func (server *Server) loginUser(ctx *gin.Context) {
	var req loginUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := server.store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	err = util.CheckPassword(req.Password, user.HashedPw)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errorResponse(err))
		return
	}
	accessToken, err := server.tokenMaker.CreateToken(user.ID, server.config.AccessTokenDuration)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}
	res := loginUserResponse{
		AccessToken: accessToken,
		User:        newUserReturn(user),
	}
	ctx.JSON(http.StatusOK, res)

}
