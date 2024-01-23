package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/jaysm12/multifinance-apps/cmd/multifinance-apps/config"
	installmentConsumer "github.com/jaysm12/multifinance-apps/internal/consumer/installment"
	"github.com/jaysm12/multifinance-apps/internal/handler/authentication"
	"github.com/jaysm12/multifinance-apps/internal/handler/installment"
	"github.com/jaysm12/multifinance-apps/internal/handler/middleware"
	"github.com/jaysm12/multifinance-apps/internal/handler/user"
	authService "github.com/jaysm12/multifinance-apps/internal/service/authentication"
	installmentService "github.com/jaysm12/multifinance-apps/internal/service/installment"
	userService "github.com/jaysm12/multifinance-apps/internal/service/user"
	creditOptionStore "github.com/jaysm12/multifinance-apps/internal/store/credit_limit"
	installmentStore "github.com/jaysm12/multifinance-apps/internal/store/installment"
	installmentPaymentHistoryStore "github.com/jaysm12/multifinance-apps/internal/store/installment_payment_history"
	userStore "github.com/jaysm12/multifinance-apps/internal/store/user"
	userKycStore "github.com/jaysm12/multifinance-apps/internal/store/user_kyc"
	"github.com/jaysm12/multifinance-apps/pkg/hash"
	"github.com/jaysm12/multifinance-apps/pkg/mysql"
	"github.com/jaysm12/multifinance-apps/pkg/rabbitmq"
	"github.com/jaysm12/multifinance-apps/pkg/token"
	"github.com/joho/godotenv"
)

// Servcer is list configuration to run Server
type Server struct {
	cfg                                  config.Config
	hashMethod                           hash.HashMethod
	tokenMethod                          token.TokenMethod
	httpServer                           *http.Server
	mysql                                mysql.MysqlMethod
	middleware                           middleware.Middleware
	authHandler                          authentication.AuthenticationHandler
	userHandler                          user.UserHandler
	installmentHandler                   installment.InstallmentHandler
	authServiceMethod                    authService.AuthenticationServiceMethod
	userStoreMethod                      userStore.UserStoreMethod
	userKycStoreMethod                   userKycStore.UserKYCStoreMethod
	installmentStoreMethod               installmentStore.InstallmentStoreMethod
	creditOptionStoreMethod              creditOptionStore.CreditOptionStoreMethod
	installmentPaymentHistoryStoreMethod installmentPaymentHistoryStore.InstallmentPaymentHistoryStoreMethod
	userServiceMethod                    userService.UserServiceMethod
	installmentServiceMethod             installmentService.InstallmentServiceMethod
	rabbitMqClient                       *rabbitmq.RabbitMqClient
}

// NewServer is func to create server with all configuration
func NewServer() (*Server, error) {
	s := &Server{}

	// ======== Init Dependencies Related ========
	// Load Env File
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Error loading .env file: %v", err)
		return s, err
	}

	// Get Config from yaml and replace by secret
	{
		cfg, err := config.GetConfig()
		if err != nil {
			fmt.Print("[Got Error]-Load Config :", err)
			return s, err
		}
		s.cfg = cfg

		log.Println("LOAD-Config")
	}

	// Init Mysql
	{
		mysqlMethod, err := mysql.NewMysqlClient(s.cfg.Mysql.Config)
		if err != nil {
			fmt.Print("[Got Error]-Mysql :", err)
			return s, err
		}

		s.mysql = mysqlMethod

		// run migration
		RunMigration(s.mysql.GetDB())

		log.Println("Init-Mysql")
	}

	// Init RabbitMQ
	{
		rabbitmqClient, err := rabbitmq.NewRabbitMQClient(s.cfg.RabbitMQ.Config)
		if err != nil {
			fmt.Print("[Got Error]-RabbitMQ :", err)
			return s, err
		}

		s.rabbitMqClient = rabbitmqClient

		log.Println("Init-RabbitMQ")
	}

	// Init Hash Package
	{
		hashMethod := hash.NewHashMethod(s.cfg.Hash.Cost)
		s.hashMethod = hashMethod
		log.Println("Init-Hash Package")
	}

	// Init Token Package
	{
		tokenMethod := token.NewTokenMethod(s.cfg.Token.Secret, s.cfg.Token.ExpInHour)
		s.tokenMethod = tokenMethod
		log.Println("Init-Token Package")
	}

	// ======== Init Dependencies Store ========
	// Init User Store
	{
		userStoreMethod := userStore.NewUserStore(s.mysql)
		s.userStoreMethod = userStoreMethod
		log.Println("Init-User Store")
	}

	{
		userKycStoreMethod := userKycStore.NewUserKYCStore(s.mysql)
		s.userKycStoreMethod = userKycStoreMethod
		log.Println("Init-User KYC Store")
	}

	{
		installmentStoreMethod := installmentStore.NewInstallmentStore(s.mysql)
		s.installmentStoreMethod = installmentStoreMethod
		log.Println("Init-Installment Store")
	}

	{
		creditOptionStoreMethod := creditOptionStore.NewCreditOptionStore(s.mysql)
		s.creditOptionStoreMethod = creditOptionStoreMethod
		log.Println("Init-Credit Limit Store")
	}

	{
		installmentPaymentHistoryStoreMethod := installmentPaymentHistoryStore.NewInstallmentPaymentHistoryStore(s.mysql)
		s.installmentPaymentHistoryStoreMethod = installmentPaymentHistoryStoreMethod
		log.Println("Init-Installment Payment History Store")
	}

	// // ======== Init Dependencies Service ========
	// // Init User Service
	{
		userServiceMethod := userService.NewUserService(s.userStoreMethod, s.userKycStoreMethod, s.creditOptionStoreMethod)
		s.userServiceMethod = userServiceMethod
		log.Println("Init-User Service")
	}

	{
		authServiceMethod := authService.NewAuthenticationService(s.userStoreMethod, s.tokenMethod, s.hashMethod)
		s.authServiceMethod = authServiceMethod
		log.Println("Init-Auth Service")
	}

	{
		installmentServiceMethod := installmentService.NewInstallmentService(s.installmentStoreMethod, s.userStoreMethod, s.creditOptionStoreMethod, s.installmentPaymentHistoryStoreMethod)
		s.installmentServiceMethod = installmentServiceMethod
		log.Println("Init-Installment Service")
	}

	// init ConsumerInstallment
	{
		installmentConsumerMethod := installmentConsumer.NewInstallmentConsumer(s.installmentServiceMethod, s.rabbitMqClient)
		errChan := make(chan error, 2)
		go installmentConsumerMethod.CreateInstallmentConsumer(errChan)
		go installmentConsumerMethod.PayInstallmentConsumer(errChan)

		go func() {
			for i := 0; i < 2; i++ {
				errFromConsumer := <-errChan
				if errFromConsumer != nil {
					log.Printf("Error from consumer: %v", errFromConsumer)
				}
			}
		}()

		log.Println("Init-Installment Consumer")
	}

	// // ======== Init Dependencies Handler ========
	// Init Middleware
	{
		midlewareService := middleware.NewMiddleware(s.tokenMethod, s.userStoreMethod)
		s.middleware = midlewareService
		log.Println("Init-Middleware")
	}

	// Init User Handler
	{
		var opts []user.Option
		opts = append(opts, user.WithTimeoutOptions(s.cfg.UserHandler.TimeoutInSec))
		userHandler := user.NewUserHandler(s.userServiceMethod, opts...)
		s.userHandler = *userHandler
		log.Println("Init-User Handler")
	}

	// Init Auth Handler
	{
		var opts []authentication.Option
		opts = append(opts, authentication.WithTimeoutOptions(s.cfg.AuthHandler.TimeoutInSec))
		authHandler := authentication.NewAuthenticationHandler(s.authServiceMethod, opts...)
		s.authHandler = *authHandler
		log.Println("Init-Auth Handler")
	}

	// Init installment handler
	{
		var opts []installment.Option
		opts = append(opts, installment.WithTimeoutOptions(s.cfg.AuthHandler.TimeoutInSec))
		installmentHandler := installment.NewInstallmentHandler(s.installmentServiceMethod, s.rabbitMqClient, opts...)
		s.installmentHandler = *installmentHandler
		log.Println("Init-Installment Handler")
	}

	// Init Router
	{
		r := mux.NewRouter()
		// Init Guest Path
		r.HandleFunc("/v1/login", s.authHandler.LoginUserHandler).Methods("POST")
		r.HandleFunc("/v1/register", s.authHandler.RegisterUserHandler).Methods("POST")

		r.HandleFunc("/v1/user/kyc", s.middleware.MiddlewareVerifyToken(s.userHandler.CreateUserKyc)).Methods("POST")

		r.HandleFunc("/v1/installment", s.middleware.MiddlewareVerifyToken(s.installmentHandler.CreateInstallment)).Methods("POST")
		r.HandleFunc("/v1/pay-installment/{contract_id}", s.middleware.MiddlewareVerifyToken(s.installmentHandler.PayInstallment)).Methods("POST")

		port := ":" + s.cfg.Port
		log.Println("running on port ", port)

		server := &http.Server{
			Addr:    port,
			Handler: r,
		}

		s.httpServer = server
	}
	return s, nil
}

func (s *Server) Start() int {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			fmt.Println(err)
		}
	}()

	// Wait for a signal to shut down the application
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a context with a timeout to allow the server to cleanly shut down
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	defer s.rabbitMqClient.Connection.Close()
	defer s.rabbitMqClient.Channel.Close()
	s.httpServer.Shutdown(ctx)
	log.Println("complete, shutting down.")
	return 0
}

// Run is func to create server and invoke Start()
func Run() int {
	s, err := NewServer()
	if err != nil {
		return 1
	}

	return s.Start()
}
