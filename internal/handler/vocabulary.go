package handler

import (
	"english-word-card-api/internal/helper"
	"english-word-card-api/internal/model"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type VocabularyHandler struct {
	db *gorm.DB
}

// 創建單字的請求結構
type CreateVocabularyRequest struct {
	Vocabulary string `json:"vocabulary" binding:"required"`
	Mean       string `json:"mean" binding:"required"`
	Category   string `json:"category"`
	Difficulty int    `json:"difficulty"`
	UserId	int   `json:"user_id"`
}

func NewVocabularyHandler(db *gorm.DB) *VocabularyHandler {
	return &VocabularyHandler{db: db}
}

// 取得所有單字
// (h * VocabularyHandler) 是 Method Receiver，告訴 Go「這個函式屬於 WordHandler」
// (c *gin.Context) 函式的輸入參數，代表這次 HTTP 請求的上下文。
func (h *VocabularyHandler) GetAllVocabulary(c *gin.Context) {
	var allVocabulary []model.Vocabulary

	userId, err := helper.GetCurrentUserID(c, h.db)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	result := h.db.Find(&allVocabulary, "user_id = ?", userId)
	
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": allVocabulary,
		"total": len(allVocabulary),
	})
}

// 取得單一單字
func (h *VocabularyHandler) GetVocabulary(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid ID",
		})
		return
	}

	var vocabulary model.Vocabulary
	result := h.db.First(&vocabulary, id)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Word not found",
			})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": result.Error.Error(),
			})
		}
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": vocabulary,
	})
}

//獲得難度5的單字列表
func (h*VocabularyHandler) GetHardVocabulary(c *gin.Context){
	var hardVocabulary []model.Vocabulary

	result := h.db.Where("difficulty = ?", 5).Limit(20).Find(&hardVocabulary)

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"data": hardVocabulary,
		"total": len(hardVocabulary),
	})
}

//創建單字
// (h * VocabularyHandler) 是 Method Receiver，告訴 Go「這個函式屬於 WordHandler」
// (c *gin.Context) 函式的輸入參數，代表這次 HTTP 請求的上下文。
func (h *VocabularyHandler) CreateVocabulary(c *gin.Context) {
	var req CreateVocabularyRequest

	// 使用 helper 取得當前用戶
	user, err := helper.GetCurrentUser(c, h.db)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 綁定 JSON 到 struct
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 創建 Word model
	vocabulary := model.Vocabulary{
		Vocabulary:    req.Vocabulary, //單字參數 (string)
		Mean:    req.Mean,  //解釋參數 (string)
		Category:   req.Category, //種類 (string)
		Difficulty: req.Difficulty, //難度 (int)
		UserId:  user.ID,// 使用 helper 取得的用戶 ID
	}

	fmt.Printf("Creating vocabulary for user %s: %+v\n", user.Username, vocabulary)

	// 存入資料庫
	result := h.db.Create(&vocabulary)
	// 如果Result.Error不等於空，代表存入失敗
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": result.Error.Error(),
		})
		return
	}
	// 回傳成功訊息
	c.JSON(http.StatusCreated, gin.H{
		"data": vocabulary,
		"message":"單字建立成功",
	})
}

 