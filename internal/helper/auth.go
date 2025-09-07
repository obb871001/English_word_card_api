package helper

import (
	"english-word-card-api/internal/model"
	"errors"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCurrentUser 從 JWT middleware 取得當前用戶
func GetCurrentUser(c *gin.Context, db *gorm.DB) (*model.User, error) {
	// 從 context 取得 user_id (username)
	userID, exists := c.Get("user_id")
	if !exists {
		return nil, errors.New("user not authenticated")
	}

	// 轉換為字串
	username, ok := userID.(string)
	if !ok {
		return nil, errors.New("invalid user ID format")
	}

	// 從資料庫查詢用戶
	var user model.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return &user, nil
}

// GetCurrentUserID 只取得用戶 ID (uint)
func GetCurrentUserID(c *gin.Context, db *gorm.DB) (uint, error) {
	user, err := GetCurrentUser(c, db)
	if err != nil {
		return 0, err
	}
	return user.ID, nil
}

// GetCurrentUsername 只取得用戶名
func GetCurrentUsername(c *gin.Context) (string, error) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", errors.New("user not authenticated")
	}

	username, ok := userID.(string)
	if !ok {
		return "", errors.New("invalid user ID format")
	}

	return username, nil
}
