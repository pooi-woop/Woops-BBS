// FilePath: C:/WoopsBBS/api/router\auth.go
package router

import (
	"WoopsBBS/api/handler"
	"github.com/gin-gonic/gin"
)

func AuthRegister(r *gin.Engine) {
	r.GET("/auth/register", handler.Register)
	r.GET("/auth/login", handler.Login)
}
