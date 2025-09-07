package handler

import (
	"english-word-card-api/internal/middleware"
	"english-word-card-api/internal/model"
	"errors"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


type AuthenticateHandler struct {
	db *gorm.DB
}

type CreateUserRequest struct {
	Username      string `json:"username" binding:"required"`
	Email     string `json:"email"`
	Password  string `json:"password" binding:"required"`
}

func NewAuthenticateHandler(db *gorm.DB) *AuthenticateHandler {
	return &AuthenticateHandler{db: db}
}


func (h *AuthenticateHandler) Login(c *gin.Context) {
	var req CreateUserRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// 檢查必需欄位
	if req.Username == "" || req.Password == "" {
		c.JSON(400, gin.H{"error": "Username and password are required"})
		return
	}

	// 從資料庫查找用戶
	var user model.User
	isNewUser := false
	if err := h.db.Where("username = ?", req.Username).First(&user).Error; err != nil {
		newUser,err := h.createUser(req)

		if err != nil{
			c.JSON(400,gin.H{"error": err.Error()})
			return
		}
		
		user = *newUser
		isNewUser = true
	}

	// 只有現有用戶才需要比對密碼
	if !isNewUser {
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
			c.JSON(401, gin.H{"error": "Invalid username or password"})
			return
		}
	}	// 生成 access token (短期的)
	accessToken := middleware.GenerateToken(req.Username)

	// refresh token (24h)
	refreshToken := middleware.GenerateRefreshToken(req.Username)

	// 更新用戶的 refresh token
	user.RefreshToken = refreshToken
	if err := h.db.Save(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to update user"})
		return
	}

    c.JSON(200, gin.H{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
        "token_type":    "Bearer",
        "expires_in":    900, // 15分鐘
    })
}

func (h *AuthenticateHandler) Logout(c *gin.Context) {
    userID := c.GetString("userID") // 從中間件獲取
    
    // 清除資料庫中的 refresh token
    h.db.Model(&model.User{}).Where("username = ?", userID).Update("refresh_token", "")
    
    c.JSON(200, gin.H{
        "message": "Logout successful",
    })
}

func (h *AuthenticateHandler) Register(c *gin.Context) {
	var req = CreateUserRequest{}
	// 獲取請求參數
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	// 檢查參數是否完整
	if req.Username == "" || req.Password == "" {
		c.JSON(400, gin.H{"error": "Missing required fields!!"})
		return
	}

	// 檢查用戶是否已存在
	var existingUser model.User
	// 檢查 username 是否重複
	if err := h.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		c.JSON(400, gin.H{"error": "Username already exists"})
		return
	}
	
	// 只有當 email 不為空時才檢查 email 重複
	if req.Email != "" {
		if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			c.JSON(400, gin.H{"error": "Email already exists"})
			return
		}
	}

	// 加密密碼
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to hash password"})
		return
	}

	// 創建新用戶
	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// 保存到資料庫
	if err := h.db.Create(&user).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create user"})
		return
	}

	// 生成 tokens
	accessToken := middleware.GenerateToken(req.Username)
	refreshToken := middleware.GenerateRefreshToken(req.Username)

	// 更新用戶的 refresh token
	user.RefreshToken = refreshToken
	h.db.Save(&user)

	c.JSON(201, gin.H{
		"message":       "User registered successfully",
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"token_type":    "Bearer",
		"expires_in":    900,
	})
}

func (h *AuthenticateHandler) RefreshToken(c *gin.Context) {
    refreshToken := c.PostForm("refresh_token")
    
    // 驗證 refresh token 並從資料庫查找
    var user model.User
    if err := h.db.Where("refresh_token = ?", refreshToken).First(&user).Error; err != nil {
        c.JSON(401, gin.H{"error": "Invalid refresh token"})
        return
    }
    
    // 生成新的 access token
    newAccessToken := middleware.GenerateToken(user.Username)
    
    c.JSON(200, gin.H{
        "access_token": newAccessToken,
        "token_type":   "Bearer",
        "expires_in":   900,
    })
}

func (h *AuthenticateHandler) createUser(req CreateUserRequest) (*model.User, error) {
	var existingUser model.User
	println("Creating user:", req.Username, req.Email)
	
	// 檢查 username 是否重複
	if err := h.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}
	
	// 只有當 email 不為空時才檢查 email 重複
	if req.Email != "" {
		if err := h.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
			return nil, errors.New("email already exists")
		}
	}

	//加密密碼
	hashedPassword,err := bcrypt.GenerateFromPassword([]byte(req.Password),bcrypt.DefaultCost)
	if err != nil {
		return nil,err
	}

	user := model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

   if err := h.db.Create(&user).Error; err != nil {
	   return nil,err
   }

	return &user, nil
}