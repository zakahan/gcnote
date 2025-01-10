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
	"gcnote/server/router/apis/recycle_apis"
	"gcnote/server/router/apis/share_apis"
	"gcnote/server/router/apis/user_apis"
	"gcnote/server/router/apis/utils_apis"
	"gcnote/server/router/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitRouter() *gin.Engine {
	// 设置 Gin 模式为调试模式
	gin.SetMode(gin.DebugMode)
	route := gin.Default()
	httpCfg := cors.Config{
		AllowOrigins:     []string{"http://localhost:8080", "http://localhost:8090"}, // 允许的前端地址 "ws://localhost:8086"
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	route.Use(cors.New(httpCfg))
	// 类似 fastapi的CORSMiddleware
	//route.Use(cors.Default()) // 中间件来启用CORS支持。这将允许来自任何源的GET，POST和OPTIONS请求，并允许特定的标头和方法
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
	group2.POST("/retrieval", index_apis.RetrievalIndex)

	// 知识库文件创建
	group2.POST("/create_file", kb_apis.CreateKBFile)
	group2.POST("/add_file", kb_apis.AddKBFile)
	group2.POST("/show_files", kb_apis.ShowIndexFiles)
	group2.POST("/recycle_file", kb_apis.RecycleKBFile)
	group2.POST("/rename_file", kb_apis.RenameKBFile)
	group2.POST("/search_file", kb_apis.SearchKBFiles) // 这个有问题，没有限制用户，能搜到别的用户的文件，这个要改一下
	group2.POST("/read_file", kb_apis.ReadFile)
	group2.POST("/recent_docs", kb_apis.RecentDocs)
	group2.POST("/update_file", kb_apis.UpdateKBFile)

	// 回收站操作
	group3 := route.Group("recycle").Use(middleware.VerifyJWT())
	group3.POST("/delete_file", recycle_apis.DeleteRecycleFile)
	group3.GET("/show_files", recycle_apis.ShowRecycleFiles)
	group3.GET("/clear", recycle_apis.ClearUserRecycleBin)
	group3.POST("/clearup", recycle_apis.CleanupOldRecycleFiles)
	group3.POST("/restore", recycle_apis.RestoreRecycleFile)

	// 分享文件相关
	group4 := route.Group("share").Use(middleware.VerifyJWT())
	group4.POST("/create", share_apis.CreateShareFile)
	group4.POST("/delete", share_apis.DeleteShareFile)
	group4.GET("/exist", share_apis.CheckShareFileExist)
	group4.GET("/info", share_apis.ListShareFiles)
	group4.POST("/read", share_apis.ReadFile)
	// 实时协作
	route.GET("/share/ws/:room", share_apis.HandleWebSocket)

	route.GET("/images/:index_id/:kb_file_id/:image_name", utils_apis.GetImage)
	route.POST("/images/upload", utils_apis.UploadImage)
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
