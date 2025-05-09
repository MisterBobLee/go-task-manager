package controllers

import (
	"go-task-api/models"
	"go-task-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginInput struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

func Register(c *gin.Context) {
    var input RegisterInput
    if err := c.ShouldBindJSON(&input); err != nil {
        if input.Username == "" && input.Password == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "請輸入使用者名稱和密碼"})
        } else if input.Username == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "請輸入使用者名稱"})
        } else if input.Password == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "請輸入密碼"})
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": "格式錯誤"})
        }
        return
    }

    if len(input.Password) < 8 || !containsLetterAndNumber(input.Password) {
        c.JSON(http.StatusBadRequest, gin.H{"error": "密碼需至少 8 碼，且包含數字與英文字母"})
        return
    }

    hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
    user := models.User{Username: input.Username, Password: string(hashedPassword)}

    if err := models.DB.Create(&user).Error; err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "使用者名稱已被註冊"})
        return
    }

    token, _ := utils.GenerateJWT(user.ID)

    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "username": user.Username,
    })
}

func Login(c *gin.Context) {
    var input LoginInput
    if err := c.ShouldBindJSON(&input); err != nil {
        if input.Username == "" && input.Password == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "請輸入使用者名稱和密碼"})
        } else if input.Username == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "請輸入使用者名稱"})
        } else if input.Password == "" {
            c.JSON(http.StatusBadRequest, gin.H{"error": "請輸入密碼"})
        } else {
            c.JSON(http.StatusBadRequest, gin.H{"error": "格式錯誤"})
        }
        return
    }

    var user models.User
    if err := models.DB.Where("username = ?", input.Username).First(&user).Error; err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "使用者不存在"})
        return
    }    

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "密碼錯誤"})
        return
    }

    token, _ := utils.GenerateJWT(user.ID)

    c.JSON(http.StatusOK, gin.H{
        "token": token,
        "username": user.Username,
    })
}

func containsLetterAndNumber(s string) bool {
	hasLetter := false
	hasDigit := false
	for _, c := range s {
		switch {
		case 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z':
			hasLetter = true
		case '0' <= c && c <= '9':
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}