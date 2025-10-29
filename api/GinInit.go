// FilePath: C:/WoopsBBS/api\GinInit.go
package api

import "github.com/gin-gonic/gin"

func GinInit() {
	r := gin.Default()
	r.Run(":8080")
}
