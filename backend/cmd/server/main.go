package main

import (
	"log"
	"mypage-backend/internal/config"
	"mypage-backend/internal/handler"
	"mypage-backend/internal/middleware"
	"mypage-backend/internal/repo"
	"mypage-backend/internal/service"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/zap"
)

func main() {
	// 加载配置
	cfg := config.Load()
	
	// 初始化日志
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	// 初始化数据库
	db, err := repo.InitDB(cfg)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}

	// 自动迁移
	if err := repo.AutoMigrate(db); err != nil {
		logger.Fatal("Failed to migrate database", zap.Error(err))
	}

	// 初始化仓库层
	userRepo := repo.NewUserRepository(db)
	commentRepo := repo.NewCommentRepository(db)

	// 初始化服务层
	userService := service.NewUserService(userRepo, cfg.JWTSecret)
	commentService := service.NewCommentService(commentRepo)

	// 初始化处理器
	authHandler := handler.NewAuthHandler(userService)
	commentHandler := handler.NewCommentHandler(commentService)

	// 创建Fiber应用
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler,
	})

	// 中间件
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000,http://127.0.0.1:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization",
		AllowCredentials: true,
	}))

	app.Use(middleware.LoggerMiddleware(logger))

	// API路由
	api := app.Group("/api")

	// 认证路由
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)
	auth.Get("/verify", middleware.JWTMiddleware(cfg.JWTSecret), authHandler.VerifyToken)

	// 评论路由
	comments := api.Group("/comments")
	comments.Get("/", commentHandler.GetComments)
	comments.Post("/", middleware.JWTMiddleware(cfg.JWTSecret), commentHandler.CreateComment)
	comments.Delete("/:id", middleware.JWTMiddleware(cfg.JWTSecret), commentHandler.DeleteComment)

	// 启动服务器
	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(app.Listen(":" + cfg.Port))
}