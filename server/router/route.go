// -------------------------------------------------
// Package router
// Author: hanzhi
// Date: 2024/12/9
// -------------------------------------------------

package router

import (
	_ "gcnote/docs"
	"gcnote/server/router/apis/index_apis"
	"gcnote/server/router/apis/kb_apis"
	"gcnote/server/router/apis/user_apis"
	"gcnote/server/router/middleware"
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
	// 用户处理
	route.POST("/user/register", user_apis.Register)
	route.POST("/user/login", user_apis.Login)
	// 用户处理01
	group1 := route.Group("user").Use(middleware.VerifyJWT())
	group1.GET("/info", user_apis.Info)
	group1.POST("/update_user_name", user_apis.UpdateUserName)
	group1.POST("/update_password", user_apis.UpdatePassword)
	group1.POST("/delete", user_apis.Delete)

	// 知识库创建等
	group2 := route.Group("index").Use(middleware.VerifyJWT())
	group2.POST("/create_index", index_apis.CreateIndex)
	group2.POST("/delete_index", index_apis.DeleteIndex)
	group2.POST("/rename_index", index_apis.RenameIndex)
	group2.GET("/show_indexes", index_apis.ShowIndexes)

	// 知识库文件创建
	group2.POST("/create_file", kb_apis.CreateKBFile)
	group2.POST("/add_file", kb_apis.AddKBFile)
	group2.GET("/show_index_files", kb_apis.ShowIndexFiles)
	group2.POST("/delete_file", kb_apis.DeleteKBFile)

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
