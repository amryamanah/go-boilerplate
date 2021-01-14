package application

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/amryamanah/go-boilerplate/internal/auth"
	"github.com/amryamanah/go-boilerplate/pkg/config"
	"github.com/amryamanah/go-boilerplate/pkg/util"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type loginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// Login return access token and response token
func (a *Application) Login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, ErrorResponse(errors.New("invalid_json")))
		return
	}

	user, err := a.Store.GetUserByEmail(ctx, req.Email)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, ErrorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, ErrorResponse(err))
		return
	}

	if err := util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		ctx.JSON(http.StatusUnauthorized, ErrorResponse(errors.New("invalid_login_details")))
		return
	}

	ts, err := auth.CreateToken(user.ID)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	saveErr := auth.CreateAuth(ctx, user.ID, ts)
	if saveErr != nil {
		ctx.JSON(http.StatusUnprocessableEntity, saveErr.Error())
		return
	}

	resp := loginResponse{
		AccessToken:  ts.AccessToken,
		RefreshToken: ts.RefreshToken,
	}

	ctx.JSON(http.StatusOK, resp)
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// Refresh return new access and refresh token based on existing refresh token
func (a *Application) Refresh(ctx *gin.Context) {
	var req refreshTokenRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, ErrorResponse(errors.New("invalid_json")))
		return
	}

	refreshToken := req.RefreshToken

	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Config.AppRefreshSecret), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, "Refresh token expired")
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		ctx.JSON(http.StatusUnauthorized, err)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string) //convert the interface to string
		if !ok {
			ctx.JSON(http.StatusUnprocessableEntity, err)
			return
		}
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["user_id"]), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusUnprocessableEntity, "Error occured")
			return
		}
		//Delete the previous Refresh Token
		deleted, delErr := auth.DeleteAuth(ctx, refreshUUID)
		if delErr != nil || deleted == 0 { //if any goes wrong
			ctx.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		//Create new pairs of refresh and access tokens
		ts, createErr := auth.CreateToken(userID)
		if createErr != nil {
			ctx.JSON(http.StatusForbidden, createErr.Error())
			return
		}
		//save the tokens metadata to redis
		saveErr := auth.CreateAuth(ctx, userID, ts)
		if saveErr != nil {
			ctx.JSON(http.StatusForbidden, saveErr.Error())
			return
		}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		ctx.JSON(http.StatusCreated, tokens)
	} else {
		ctx.JSON(http.StatusUnauthorized, "refresh expired")
	}
}

// SignUp endpoint to register new user
func (a *Application) SignUp(ctx *gin.Context) {
	ctx.JSON(http.StatusNotImplemented, ErrorResponse(errors.New("not_implemented")))
}
