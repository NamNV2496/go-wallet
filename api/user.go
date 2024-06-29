package api

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	db "github.com/namnv2496/go-wallet/internal/databaseaccess/sqlc"
	"github.com/namnv2496/go-wallet/internal/token"
	"github.com/namnv2496/go-wallet/internal/util"
)

type getUsersByUsernameOrPhoneRequest struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Limit    int32  `json:"limit" binding:"required"`
}

func (s *Server) getUsersByUsernameOrPhone(ctx *gin.Context) {

	var req getUsersByUsernameOrPhoneRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	result, err := s.userService.GetUsersByUsernameOrPhone(ctx, req.Username, req.Phone, req.Limit)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, result)
}

type createUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	Phone             string    `json:"phone"`
	Role              string    `json:"role"`
	IsEmailVerified   string    `json:"is_email_verified"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username.String,
		FullName:          user.FullName,
		Email:             user.Email,
		Phone:             user.Phone,
		Role:              user.Role,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}

func (s *Server) createUser(ctx *gin.Context) {

	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := s.userService.CreateUser(ctx, req.Username, req.Password, req.FullName, req.Email, req.Phone, req.Role)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusCreated, newUserResponse(user))
}

func (s *Server) getUser(ctx *gin.Context) {

	username := ctx.Query("username")
	authPayload := ctx.MustGet(authorPayloadKey).(*token.Payload)
	if authPayload.Username != username {
		err := errors.New("Cannot get user information of another person")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := s.userService.GetUser(ctx, username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, newUserResponse(user))
}

type updateUserRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Phone    string `json:"phone" binding:"required"`
}

func (s *Server) updateUser(ctx *gin.Context) {

	var req updateUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	authPayload := ctx.MustGet(authorPayloadKey).(*token.Payload)
	if authPayload.Username != req.Username {
		err := errors.New("Cannot update user information of another person")
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	user, err := s.userService.UpdateUser(ctx, req.Password, req.FullName, req.Email, req.Phone, req.Username)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, newUserResponse(user))

}
func (s *Server) verifyuser(ctx *gin.Context) {

	userName := ctx.Query("username")
	err := s.userService.VerifyEmail(ctx, userName)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"status": "thành công"})
}

type loginRequest struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

type loginResponse struct {
	Token     string    `json:"token"`
	Username  string    `json:"username"`
	ExpiredAt time.Time `json:"expired_at"`
}

func (s *Server) login(ctx *gin.Context) {
	var req loginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		log.Println("Invalid request")
		return
	}

	user, err := s.userService.GetUser(ctx, req.Username)
	if err != nil {
		log.Println("Failed to get user")
		return
	}
	if !util.IsCorrectPassword(req.Password, user.HashedPassword) {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Wrong password"})
		return
	}

	token, accessPayload, err := s.token.CreateToken(user.ID, req.Username, user.Role, time.Duration(time.Hour*24))
	if err != nil {
		log.Println("Failed to create token")
		return
	}

	rsp := loginResponse{
		Token:     token,
		Username:  req.Username,
		ExpiredAt: accessPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, rsp)
}
