// FilePath: C:/WoopsBBS/api/handler\auth.go
package handler

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"

	"WoopsBBS/global/Database/mysql"
	"WoopsBBS/global/model"
	"github.com/bwmarrin/snowflake"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest 注册请求参数
type RegisterRequest struct {
	Name     string `form:"name" binding:"required,min=3,max=50"`
	Password string `form:"password" binding:"required,min=6"`
	Email    string `form:"email" binding:"required,email"`
}

// 雪花算法节点
var node *snowflake.Node

func init() {
	// 初始化雪花算法节点（使用默认节点0）
	n, err := snowflake.NewNode(0)
	if err != nil {
		panic("初始化雪花算法节点失败: " + err.Error())
	}
	node = n
}

// 生成随机盐值
func generateSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

// 密码加盐哈希
func hashPassword(password, salt string) (string, error) {
	// 将盐值添加到密码中
	saltedPassword := password + salt
	// 使用bcrypt进行哈希，成本因子设为10
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), 10)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// Register 用户注册函数
func Register(c *gin.Context) {
	// 1. 绑定并验证请求参数
	var req RegisterRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"code":    400,
			"message": "参数错误: " + err.Error(),
		})
		return
	}

	// 2. 检查用户名是否已存在
	var existingUser model.User
	result := mysql.DB.Where("username = ?", req.Name).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "用户名已存在",
		})
		return
	}

	// 3. 检查邮箱是否已存在
	result = mysql.DB.Where("email = ?", req.Email).First(&existingUser)
	if result.Error == nil {
		c.JSON(http.StatusConflict, gin.H{
			"code":    409,
			"message": "邮箱已被注册",
		})
		return
	}

	// 4. 生成雪花ID
	userID := node.Generate().Int64()

	// 5. 生成盐值
	salt, err := generateSalt()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "生成盐值失败",
		})
		return
	}

	// 6. 对密码进行加盐哈希
	hashedPassword, err := hashPassword(req.Password, salt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "密码加密失败",
		})
		return
	}

	// 7. 创建用户模型实例
	user := model.User{
		UserID:   userID,
		Username: req.Name,
		Password: hashedPassword,
		Salt:     salt,
		Email:    req.Email,
		// Avatar和Homepage使用默认值
	}

	// 8. 保存到数据库
	if err := mysql.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code":    500,
			"message": "创建用户失败: " + err.Error(),
		})
		return
	}

	// 9. 返回成功响应
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "注册成功",
		"data": gin.H{
			"user_id":  user.UserID,
			"username": user.Username,
			"email":    user.Email,
		},
	})
}

func Login(c *gin.Context) {
	// 登录功能待实现
	c.JSON(http.StatusOK, gin.H{
		"code":    200,
		"message": "登录接口待实现",
	})
}
