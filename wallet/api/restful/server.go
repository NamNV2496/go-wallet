package restful

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/namnv2496/go-wallet/config"
	"github.com/namnv2496/go-wallet/internal/databaseaccess"
	"github.com/namnv2496/go-wallet/internal/logic"
	"github.com/namnv2496/go-wallet/internal/mq/consumer"
	"github.com/namnv2496/go-wallet/internal/mq/producer"
	"github.com/namnv2496/go-wallet/internal/token"
	"github.com/namnv2496/go-wallet/internal/worker"
)

type Server struct {
	config          config.Config
	router          *gin.Engine
	token           token.Maker
	accountService  logic.AccountLogic
	userService     logic.UserLogic
	transferService logic.TransferLogic
	sessionService  logic.SessionLogic
	producer        *producer.Producer
	Consumer        *consumer.Consumer
	Queue           worker.RedisTaskProcessor
}

func NewGinServer(
	config config.Config,
	token token.Maker,
	accountService logic.AccountLogic,
	userService logic.UserLogic,
	transferService logic.TransferLogic,
	sessionService logic.SessionLogic,
	producer *producer.Producer,
	consumer *consumer.Consumer,
	queue worker.RedisTaskProcessor,
) (*Server, error) {

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
		v.RegisterValidation("role", validRole)
	}
	server := &Server{
		config:          config,
		token:           token,
		accountService:  accountService,
		userService:     userService,
		transferService: transferService,
		sessionService:  sessionService,
		producer:        producer,
		Consumer:        consumer,
		Queue:           queue,
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "*"
		},
		MaxAge: 12 * time.Hour,
	}))

	r.POST("/migration", server.migration)
	router := r.Group("/api/v1/")

	router.POST("user", server.createUser)
	router.GET("verify_email", server.verifyuser)
	router.POST("users/login", server.login)
	router.POST("users/renew_access", server.renewAccessToken)

	authRoutes := router.Group("/").Use(authMiddleware(server.token))

	authRoutes.GET("user", server.getUser)
	authRoutes.POST("users", server.getUsersByUsernameOrPhone)
	authRoutes.PUT("user", server.updateUser)

	authRoutes.POST("account", server.createAccount)
	authRoutes.GET("accounts/:id", server.getAccount)
	authRoutes.GET("accounts", server.getAccountByUserId)
	authRoutes.GET("listAccounts", server.listAccounts)

	authRoutes.POST("transfer", server.createTransfer)
	authRoutes.GET("transfer/:id", server.getTransfer)
	authRoutes.GET("transfers", server.listTransfers)

	// TCC
	router.POST("TransIn", server.TccTransInTry)
	router.POST("TransInConfirm", server.TccTransInConfirm)
	router.POST("TransInRevert", server.TccTransInCancel)

	router.POST("TransOut", server.TccTransOutTry)
	router.POST("TransOutConfirm", server.TccTransOutConfirm)
	router.POST("TransOutRevert", server.TccTransOutCancel)

	// authRoutes.POST("transfers/fakeResult", server.fakeResult)

	server.router = r
}

func (server *Server) StartHttpServer(address string) error {
	server.router.Run(address)
	return nil
}

func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

func (server *Server) migration(ctx *gin.Context) {
	databaseaccess.MigrateUp(ctx)
}
