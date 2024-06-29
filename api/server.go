package api

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/namnv2496/go-wallet/internal/databaseaccess"
	"github.com/namnv2496/go-wallet/internal/logic"
	"github.com/namnv2496/go-wallet/internal/token"
)

type Server struct {
	router          *gin.Engine
	token           token.Maker
	accountService  logic.AccountLogic
	userService     logic.UserLogic
	transferService logic.TranserLogic
}

func NewGinServer(
	token token.Maker,
	accountService logic.AccountLogic,
	userService logic.UserLogic,
	transferService logic.TranserLogic,
) (*Server, error) {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
		v.RegisterValidation("role", validRole)
	}
	server := &Server{
		token:           token,
		accountService:  accountService,
		userService:     userService,
		transferService: transferService,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	r := gin.Default()
	r.POST("/migration", server.migration)

	router := r.Group("/api/v1/")

	router.POST("user", server.createUser)
	router.PUT("verify_user", server.verifyuser)
	router.POST("users/login", server.login)
	// router.POST("/tokens/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.token))

	authRoutes.GET("user", server.getUser)
	authRoutes.GET("users", server.getUsersByUsernameOrPhone)
	authRoutes.PUT("user", server.updateUser)

	authRoutes.POST("account", server.createAccount)
	authRoutes.GET("accounts/:id", server.getAccount)
	authRoutes.GET("accounts", server.getAccountByUserId)
	authRoutes.GET("listAccounts", server.listAccounts)

	authRoutes.POST("transfer", server.createTransfer)
	authRoutes.GET("transfer/:id", server.getTransfer)
	authRoutes.GET("transfers", server.listTransfers)

	server.router = r
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) migration(ctx *gin.Context) {
	databaseaccess.MigrateUp(ctx)
}
