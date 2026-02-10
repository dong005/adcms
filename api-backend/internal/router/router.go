package router

import (
	"adcms/internal/handler"
	"adcms/internal/middleware"
	"time"

	_ "adcms/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func SetupRouter(mode string) *gin.Engine {
	gin.SetMode(mode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.Cors())

	r.Static("/uploads", "./uploads")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	authHandler := handler.NewAuthHandler()
	menuHandler := handler.NewMenuHandler()
	userHandler := handler.NewUserHandler()
	roleHandler := handler.NewRoleHandler()
	categoryHandler := handler.NewCategoryHandler()
	tagHandler := handler.NewTagHandler()
	articleHandler := handler.NewArticleHandler()
	mediaHandler := handler.NewMediaHandler()
	configHandler := handler.NewConfigHandler()
	logHandler := handler.NewLogHandler()
	permissionHandler := handler.NewPermissionHandler()
	departmentHandler := handler.NewDepartmentHandler()
	notificationHandler := handler.NewNotificationHandler()

	api := r.Group("/api")
	api.Use(middleware.GlobalRateLimit(300)) // 每个IP每分钟最多300次请求
	{
		auth := api.Group("/auth")
		{
			auth.POST("/login", middleware.RateLimit(10, time.Minute), authHandler.Login)
			auth.POST("/verify-totp", middleware.RateLimit(10, time.Minute), authHandler.VerifyTOTP)
			auth.POST("/forgot-password", middleware.RateLimit(5, time.Minute), authHandler.ForgotPassword)
			auth.POST("/reset-password", middleware.RateLimit(10, time.Minute), authHandler.ResetPasswordByEmail)
		}

		protected := api.Group("")
		protected.Use(middleware.JWTAuth())
		protected.Use(middleware.APIPermissionCheck())
		protected.Use(middleware.DataScopeFilter())
		protected.Use(middleware.OperationLogger())
		{
			// Auth
			protectedAuth := protected.Group("/auth")
			{
				protectedAuth.GET("/user-info", authHandler.GetUserInfo)
				protectedAuth.PUT("/user-info", authHandler.UpdateUserInfo)
				protectedAuth.PUT("/password", authHandler.ChangePassword)
				protectedAuth.POST("/logout", authHandler.Logout)
				protectedAuth.POST("/totp/generate", authHandler.GenerateTOTP)
				protectedAuth.POST("/totp/bind", authHandler.BindTOTP)
				protectedAuth.POST("/totp/disable", authHandler.DisableTOTP)
				protectedAuth.GET("/codes", authHandler.GetPermissionCodes)
				protectedAuth.GET("/login-history", authHandler.LoginHistory)
				protectedAuth.POST("/send-sms-code", authHandler.SendSmsCode)
				protectedAuth.POST("/bind-phone", authHandler.BindPhone)
			}

			// Menus
			menus := protected.Group("/menus")
			{
				menus.GET("", menuHandler.List)
				menus.GET("/tree", menuHandler.Tree)
				menus.GET("/user", menuHandler.UserMenus)
				menus.POST("", menuHandler.Create)
				menus.PUT("/:id", menuHandler.Update)
				menus.DELETE("/:id", menuHandler.Delete)
			}

			// Users
			users := protected.Group("/users")
			{
				users.GET("", userHandler.List)
				users.GET("/:id", userHandler.Detail)
				users.POST("", userHandler.Create)
				users.PUT("/:id", userHandler.Update)
				users.DELETE("/:id", userHandler.Delete)
				users.PUT("/:id/status", userHandler.ToggleStatus)
				users.PUT("/:id/reset-password", userHandler.ResetPassword)
				users.PUT("/:id/roles", userHandler.AssignRoles)
				users.GET("/export", userHandler.Export)
				users.GET("/import-template", userHandler.ImportTemplate)
				users.POST("/import", userHandler.Import)
			}

			// Roles
			roles := protected.Group("/roles")
			{
				roles.GET("", roleHandler.List)
				roles.GET("/:id", roleHandler.Detail)
				roles.POST("", roleHandler.Create)
				roles.PUT("/:id", roleHandler.Update)
				roles.DELETE("/:id", roleHandler.Delete)
				roles.PUT("/:id/menus", menuHandler.AssignRoleMenus)
				roles.GET("/:id/menus", menuHandler.GetRoleMenus)
				roles.PUT("/:id/permissions", roleHandler.AssignPermissions)
				roles.GET("/:id/permissions", roleHandler.GetPermissions)
			}

			// Departments
			departments := protected.Group("/departments")
			{
				departments.GET("", departmentHandler.List)
				departments.GET("/tree", departmentHandler.Tree)
				departments.POST("", departmentHandler.Create)
				departments.PUT("/:id", departmentHandler.Update)
				departments.DELETE("/:id", departmentHandler.Delete)
			}

			// Categories
			categories := protected.Group("/categories")
			{
				categories.GET("", categoryHandler.List)
				categories.POST("", categoryHandler.Create)
				categories.PUT("/:id", categoryHandler.Update)
				categories.DELETE("/:id", categoryHandler.Delete)
			}

			// Tags
			tags := protected.Group("/tags")
			{
				tags.GET("", tagHandler.List)
				tags.POST("", tagHandler.Create)
				tags.PUT("/:id", tagHandler.Update)
				tags.DELETE("/:id", tagHandler.Delete)
			}

			// Articles
			articles := protected.Group("/articles")
			{
				articles.GET("", articleHandler.List)
				articles.GET("/:id", articleHandler.Detail)
				articles.POST("", articleHandler.Create)
				articles.PUT("/:id", articleHandler.Update)
				articles.DELETE("/:id", articleHandler.Delete)
				articles.PUT("/:id/publish", articleHandler.Publish)
				articles.PUT("/:id/draft", articleHandler.Draft)
			}

			// Media
			media := protected.Group("/media")
			{
				media.GET("", mediaHandler.List)
				media.POST("/upload", mediaHandler.Upload)
				media.DELETE("/:id", mediaHandler.Delete)
			}

			// System Configs
			configs := protected.Group("/configs")
			{
				configs.GET("", configHandler.List)
				configs.PUT("", configHandler.Update)
				configs.GET("/email", configHandler.GetEmailConfig)
				configs.PUT("/email", configHandler.UpdateEmailConfig)
				configs.POST("/email/test", configHandler.TestEmail)
				configs.GET("/sms", configHandler.GetSmsConfig)
				configs.PUT("/sms", configHandler.UpdateSmsConfig)
				configs.POST("/sms/test", configHandler.TestSms)
				configs.GET("/log", configHandler.GetLogConfig)
				configs.PUT("/log", configHandler.UpdateLogConfig)
			}

			// Logs
			logs := protected.Group("/logs")
			{
				logs.GET("/operation", logHandler.OperationLogs)
				logs.GET("/login", logHandler.LoginLogs)
				logs.GET("/email", logHandler.EmailLogs)
				logs.GET("/sms", logHandler.SmsLogs)
			}

			// Permissions
			permissions := protected.Group("/permissions")
			{
				permissions.GET("", permissionHandler.List)
				permissions.GET("/tree", permissionHandler.Tree)
				permissions.POST("", permissionHandler.Create)
				permissions.PUT("/:id", permissionHandler.Update)
				permissions.DELETE("/:id", permissionHandler.Delete)
			}

			// Notifications
			notifications := protected.Group("/notifications")
			{
				notifications.GET("", notificationHandler.List)
				notifications.GET("/unread-count", notificationHandler.UnreadCount)
				notifications.GET("/:id", notificationHandler.Detail)
				notifications.PUT("/:id/read", notificationHandler.MarkAsRead)
				notifications.PUT("/read-all", notificationHandler.MarkAllAsRead)
				notifications.DELETE("/:id", notificationHandler.Delete)
				notifications.POST("/send", notificationHandler.Send)
				notifications.POST("/:id/reply", notificationHandler.Reply)
				notifications.DELETE("/reply/:id", notificationHandler.DeleteReply)
			}
		}
	}

	return r
}
