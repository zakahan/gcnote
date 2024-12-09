package router

import (
	"gcnote/server/router/apis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	// 设置 Gin 模式为调试模式
	gin.SetMode(gin.DebugMode)
	route := gin.Default()
	// 类似 fastapi的CORSMiddleware
	route.Use(cors.Default()) // 中间件来启用CORS支持。这将允许来自任何源的GET，POST和OPTIONS请求，并允许特定的标头和方法
	route.POST("/user/register", apis.Register)
	route.POST("/user/login", apis.Login)

	// swagger
	route.GET(
		"/swagger/*any",
		ginSwagger.WrapHandler(
			swaggerFiles.Handler,
			ginSwagger.DocExpansion("none"),
		),
	)

	return route
}
