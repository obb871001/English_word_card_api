package middleware

import (
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var jwtSecret = []byte("secret")

func GenerateToken(userID string) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,jwt.MapClaims{
		"user_id": userID,
		"expired":time.Now().Add(15 * time.Minute).Unix(),
	})

	tokenString, _ := token.SignedString(jwtSecret)
	print("Generated token: %s\n", tokenString)
	return tokenString
}

func GenerateRefreshToken(userID string) string {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": userID,
        "token_type": "refresh",  // 標記為 refresh token
        "expired": time.Now().Add(7 * 24 * time.Hour).Unix(), // 7天過期
    })

    tokenString, _ := token.SignedString(jwtSecret)
    print("Generated refresh token: %s\n", tokenString)
    return tokenString
}

func Authenticate(c *gin.Context) {
	// 從 Authorization header 獲取 token
	authHeader := c.GetHeader("Authorization")
	
	// 檢查 header 格式：Bearer <token>
	if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Missing or invalid authorization header",
		})
		c.Abort()
		return
	}
	
	// 提取 token（去掉 "Bearer " 前綴）
	tokenString := authHeader[7:]

	// 解析 token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 確保使用正確的簽名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.NewValidationError("Invalid signing method", jwt.ValidationErrorSignatureInvalid)
		}
		return jwtSecret, nil
	})

	// 檢查 token 是否有效
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid or expired token",
		})
		c.Abort()
		return
	}

	// 獲取 claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid token claims",
		})
		c.Abort()
		return
	}

	// 檢查過期時間
	if exp, ok := claims["expired"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Token expired",
			})
			c.Abort()
			return
		}
	}

	// 獲取用戶 ID
	userID, ok := claims["user_id"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Invalid user ID in token",
		})
		c.Abort()
		return
	}

	// 將用戶 ID 存到 context 中，供後續處理器使用
	c.Set("user_id", userID)
	c.Next() // 繼續執行下一個處理器
}

