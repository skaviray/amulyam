package api

import (
	db "amulyam/db/sqlc"
	"amulyam/token"
	"amulyam/utils"
	"context"
	"fmt"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"github.com/razorpay/razorpay-go"
	"golang.org/x/crypto/bcrypt"
)

type Server struct {
	config         utils.Config
	store          db.Store
	tokenMaker     token.Maker
	razorPayClient *razorpay.Client
	router         *gin.Engine
}

func New(store db.Store, config utils.Config) (*Server, error) {
	maker, err := token.NewPasetoMaker(config.SecretKey)
	razorPayClient := razorpay.NewClient(config.RazorpayKey, config.RazorpaySecret)
	if err != nil {
		return nil, fmt.Errorf("unable to create token maker: %c", err)
	}
	server := &Server{
		config:         config,
		store:          store,
		tokenMaker:     maker,
		razorPayClient: razorPayClient,
	}
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}
func (server *Server) setupRouter() {
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://127.0.0.1:3000", "http://localhost:3001", "http://localhost:3002"}, // Adjust based on your frontend URL
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "PUT", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	authRoutes := router.Group("/api")
	{
		// User Public Routes
		authRoutes.POST("/register", server.CreateUser)
		authRoutes.POST("/login", server.LoginUser)
		authRoutes.POST("/users", server.CreateUser)
		authRoutes.GET("/userinfo", server.GetUserInfo)

		// Payment routes
		authRoutes.POST("/create-order", server.CreateOrder)
		authRoutes.POST("/verify-payment", server.VerifySignature)
	}

	protected := router.Group("/api").Use(server.authMiddleware(), server.isAdmin())
	// protected.GET("/", server.GetUser)
	// protected.POST("/users", server.CreateUser)
	protected.GET("/users/:id", server.GetUser)
	protected.GET("/users", server.ListUsers)
	protected.PATCH("/users/password", server.UpdateUserPassword)
	protected.DELETE("/users/:id", server.DeleteUser)

	// customers Routes
	protected.POST("/customers", server.CreateCustomer)
	protected.GET("/customers", server.ListCustomers)
	protected.GET("/customers/:id", server.GetCustomer)
	protected.PATCH("/customers", server.UpdateCustomer)
	protected.DELETE("/customers/:id", server.DeleteCustomer)

	server.router = router
}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func (server *Server) CreateDefaultAdminUser() error {
	ctx := context.Background()
	email := "admin@example.com"
	password := "Admin123@" // You can hash this

	// Check if admin exists
	_, err := server.store.GetUserByEmail(ctx, email)
	if err == nil {
		return nil // Admin already exists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}
	arg := db.CreateUserParams{
		Username:          "admin",
		HashedPassword:    string(hashedPassword),
		Email:             email,
		FullName:          "System Admin",
		PasswordChangedAt: time.Now(),
		IsAdmin:           true,
	}

	_, err = server.store.CreateUser(ctx, arg)
	if err != nil {
		return fmt.Errorf("failed to create default admin user: %w", err)
	}

	return nil
}
