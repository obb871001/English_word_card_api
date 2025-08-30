package main

import (
	"english-word-card-api/internal/config"
	"english-word-card-api/internal/handler"
	"english-word-card-api/internal/middleware"
	"english-word-card-api/internal/model"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化資料庫
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	
	// 自動遷移
	db.AutoMigrate(&model.Vocabulary{}, &model.User{})
	
	// 初始化處理器
	vocabularyHandler := handler.NewVocabularyHandler(db)
	auth := handler.NewAuthenticateHandler(db)
	
	router := gin.Default()
	
	// 設置 CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
	
	// API 路由
	api := router.Group("/api")
	
	// 公開路由（不需要驗證）
	public := api.Group("/")
	{
		// 認證相關
		public.POST("auth/login", auth.Login)
		public.POST("auth/register", auth.Register)
		public.POST("auth/refresh", auth.RefreshToken)
		
		// 公開的單字查詢（如果你想要公開的話）
		public.GET("vocabulary", vocabularyHandler.GetAllVocabulary)
		public.GET("vocabulary/:id", vocabularyHandler.GetVocabulary)
	}
	
	// 受保護路由（需要驗證）
	protected := api.Group("/")
	protected.Use(middleware.Authenticate) // 使用 JWT 中間件
	{
		// 需要登入才能創建單字
		protected.POST("vocabulary", vocabularyHandler.CreateVocabulary)
		
		// 登出（需要驗證才能登出）
		protected.POST("auth/logout", auth.Logout)
		
		// 其他需要驗證的 API...
	}
	
	router.Run(":8080")
}

