package application

import (
	"database/sql"
	"fmt"
	"github.com/amryamanah/go-boilerplate/internal/auth"
	store "github.com/amryamanah/go-boilerplate/internal/store/sqlc"
	"github.com/amryamanah/go-boilerplate/pkg/util"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"net/http"
	"time"
)

type createUserRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
	Password string `json:"password" binding:"required,min=6"`
}

type createUserResponse struct {
	Id                int64     `json:"id"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	FullName          string    `json:"full_name"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func (a *Application) CreateUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse(err))
		return
	}

	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	fmt.Println(req)


	arg := store.CreateUserParams{
		Email: req.Email,
		HashedPassword: hashedPassword,
	}
	if req.FullName != "" {
		arg.FullName = sql.NullString{
			String: req.FullName,
			Valid:  true,
		}
	}
	if req.Phone != "" {
		arg.Phone = sql.NullString{
			String: req.Phone,
			Valid:  true,
		}
	}

	user, err := a.Store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "unique_violation":
				ctx.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	rsp := createUserResponse{
		Id:                user.ID,
		FullName:          user.FullName.String,
		Email:             user.Email,
		Phone:             user.Phone.String,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}

	ctx.JSON(http.StatusOK, rsp)
}

func (a *Application) GetMe(ctx *gin.Context) {
	tokenAuth, err := auth.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}

	userId, err := auth.FetchAuth(ctx, tokenAuth)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(err))
		return
	}
	fmt.Printf("tokenAuth: %+v , userId: %v \n", tokenAuth, userId)
	user, err := a.Store.GetUserByID(ctx, userId)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, ErrorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, user)
}

func (a *Application) Logout(ctx *gin.Context) {
	au, err := auth.ExtractTokenMetadata(ctx.Request)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	deleted, delErr := auth.DeleteAuth(ctx, au.AccessUuid)
	if delErr != nil || deleted == 0 { //if any goes wrong
		ctx.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}
	ctx.JSON(http.StatusOK, "Successfully logged out")
}
