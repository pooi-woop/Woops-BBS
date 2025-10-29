// FilePath: C:/WoopsBBS/global/model\User.go
package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	
	UserID    int64     `gorm:"primaryKey" json:"user_id"`
	Username  string    `gorm:"size:50;not null;uniqueIndex" json:"username"`
	Password  string    `gorm:"size:255;not null" json:"-"` // 密码哈希值，不返回给前端
	Salt      string    `gorm:"size:64;not null" json:"-"`  // 盐值，不返回给前端
	Avatar    string    `gorm:"size:255" json:"avatar"`     // 头像URL
	Email     string    `gorm:"size:100;uniqueIndex" json:"email"`
	Homepage  string    `gorm:"size:255" json:"homepage"` // 个人主页
	CreatedAt time.Time `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
