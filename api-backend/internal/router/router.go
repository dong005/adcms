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
	adminHandler := handler.NewAdminHandler()
	categoryHandler := handler.NewCategoryHandler()
	tagHandler := handler.NewTagHandler()
	articleHandler := handler.NewArticleHandler()
	mediaHandler := handler.NewMediaHandler()
	configHandler := handler.NewConfigHandler()
	logHandler := handler.NewLogHandler()
	permissionHandler := handler.NewPermissionHandler()
	departmentHandler := handler.NewDepartmentHandler()
	notificationHandler := handler.NewNotificationHandler()
	dictHandler := handler.NewDictHandler()
	siteHandler := handler.NewSiteHandler()
	linkHandler := handler.NewLinkHandler()
	crontabHandler := handler.NewCrontabHandler()
	databaseHandler := handler.NewDatabaseHandler()
	cityHandler := handler.NewCityHandler()

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
				menus.GET("/:id/menus", menuHandler.GetUserMenus) // 新增
				menus.PUT("/:id/menus", menuHandler.AssignUserMenus) // 新增
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
				users.PUT("/:id/menus", userHandler.AssignMenus) // 新增
				users.PUT("/:id/unlock", userHandler.UnlockUser)
				users.GET("/export", userHandler.Export)
				users.GET("/import-template", userHandler.ImportTemplate)
				users.POST("/import", userHandler.Import)
			}

			// Admins - 管理员管理（仅超管）
			admins := protected.Group("/admins")
			{
				admins.GET("", adminHandler.List)
				admins.GET("/:id", adminHandler.Detail)
				admins.POST("", adminHandler.Create)
				admins.PUT("/:id", adminHandler.Update)
				admins.DELETE("/:id", adminHandler.Delete)
				admins.PUT("/:id/status", adminHandler.ToggleStatus)
				admins.PUT("/:id/reset-password", adminHandler.ResetPassword)
				admins.GET("/:id/statistics", adminHandler.Statistics)
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
				configs.GET("/by-group", configHandler.ListByGroup)
				configs.GET("/email", configHandler.GetEmailConfig)
				configs.PUT("/email", configHandler.UpdateEmailConfig)
				configs.POST("/email/test", configHandler.TestEmail)
				configs.GET("/sms", configHandler.GetSmsConfig)
				configs.PUT("/sms", configHandler.UpdateSmsConfig)
				configs.POST("/sms/test", configHandler.TestSms)
				configs.GET("/log", configHandler.GetLogConfig)
				configs.PUT("/log", configHandler.UpdateLogConfig)
			}

			// Config Groups
			configGroups := protected.Group("/config-groups")
			{
				configGroups.GET("", configHandler.ListGroups)
				configGroups.POST("", configHandler.CreateGroup)
				configGroups.PUT("/:id", configHandler.UpdateGroup)
				configGroups.DELETE("/:id", configHandler.DeleteGroup)
			}

			// Config Webs (网站设置)
			configWebs := protected.Group("/config-webs")
			{
				configWebs.GET("", configHandler.ListWebs)
				configWebs.PUT("", configHandler.SaveWebs)
				configWebs.DELETE("/:id", configHandler.DeleteWeb)
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

			// Dict Types
			dictTypes := protected.Group("/dict-types")
			{
				dictTypes.GET("", dictHandler.ListTypes)
				dictTypes.POST("", dictHandler.CreateType)
				dictTypes.PUT("/:id", dictHandler.UpdateType)
				dictTypes.DELETE("/:id", dictHandler.DeleteType)
			}

			// Dicts
			dicts := protected.Group("/dicts")
			{
				dicts.GET("", dictHandler.ListDicts)
				dicts.POST("", dictHandler.CreateDict)
				dicts.PUT("/:id", dictHandler.UpdateDict)
				dicts.DELETE("/:id", dictHandler.DeleteDict)
				dicts.GET("/code/:code", dictHandler.GetDictsByCode)
			}

			// Sites
			sites := protected.Group("/sites")
			{
				sites.GET("", siteHandler.List)
				sites.POST("", siteHandler.Create)
				sites.PUT("/:id", siteHandler.Update)
				sites.DELETE("/:id", siteHandler.Delete)
			}

			// Links
			links := protected.Group("/links")
			{
				links.GET("", linkHandler.List)
				links.POST("", linkHandler.Create)
				links.PUT("/:id", linkHandler.Update)
				links.DELETE("/:id", linkHandler.Delete)
			}

			// Crontabs
			crontabs := protected.Group("/crontabs")
			{
				crontabs.GET("", crontabHandler.List)
				crontabs.POST("", crontabHandler.Create)
				crontabs.PUT("/:id", crontabHandler.Update)
				crontabs.DELETE("/:id", crontabHandler.Delete)
			}

			// Database
			db := protected.Group("/database")
			{
				db.GET("/tables", databaseHandler.Tables)
				db.GET("/tables/:table/columns", databaseHandler.Columns)
			}

			// Cities
			cities := protected.Group("/cities")
			{
				cities.GET("", cityHandler.List)
				cities.GET("/tree", cityHandler.Tree)
				cities.POST("", cityHandler.Create)
				cities.PUT("/:id", cityHandler.Update)
				cities.DELETE("/:id", cityHandler.Delete)
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
