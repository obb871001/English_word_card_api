package main

import (
	"english-word-card-api/internal/config"
	"english-word-card-api/internal/handler"
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
	db.AutoMigrate(&model.Vocabulary{})
	
	// 初始化處理器
	vocabularyHandler := handler.NewVocabularyHandler(db)
	
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
	{
		api.GET("/vocabulary", vocabularyHandler.GetAllVocabulary)
		api.GET("/vocabulary/:id", vocabularyHandler.GetVocabulary)
		api.POST("/vocabulary", vocabularyHandler.CreateVocabulary)
	}
	
	router.Run(":8080")
}

