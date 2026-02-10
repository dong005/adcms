package database

import (
	"adcms/internal/model"
	"adcms/pkg/utils"
	"time"
)

func AutoMigrate() error {
	return DB.AutoMigrate(
		&model.User{},
		&model.Role{},
		&model.Permission{},
		&model.UserRole{},
		&model.RolePermission{},
		&model.RoleMenu{},
		&model.Menu{},
		&model.Category{},
		&model.Article{},
		&model.Tag{},
		&model.ArticleTag{},
		&model.Media{},
		&model.SystemConfig{},
		&model.OperationLog{},
		&model.LoginLog{},
		&model.Department{},
		&model.UserDepartment{},
		&model.Notification{},
		&model.EmailLog{},
		&model.SmsLog{},
	)
}

func InitData() error {
	var count int64
	DB.Model(&model.User{}).Count(&count)
	if count > 0 {
		return fixMenuRelations()
	}

	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		return err
	}

	roles := []model.Role{
		{
			TenantBaseModel: model.TenantBaseModel{TenantID: 0},
			Name: "超级管理员", Code: "super_admin", Description: "拥有所有权限", Status: 1, Sort: 0,
		},
		{
			TenantBaseModel: model.TenantBaseModel{TenantID: 0},
			Name: "管理员", Code: "admin", Description: "代理商管理员", Status: 1, Sort: 1,
		},
		{
			TenantBaseModel: model.TenantBaseModel{TenantID: 0},
			Name: "普通用户", Code: "user", Description: "仅可查看和编辑自己的内容", Status: 1, Sort: 2,
		},
	}
	for i := range roles {
		if err := DB.Create(&roles[i]).Error; err != nil {
			return err
		}
	}
	superAdminRole := roles[0]

	now := time.Now()
	admin := model.User{
		TenantBaseModel: model.TenantBaseModel{
			TenantID: 0,
		},
		Username:    "admin",
		Password:    hashedPassword,
		Email:       "admin@adcms.com",
		Nickname:    "管理员",
		Status:      1,
		LastLoginAt: &now,
	}
	if err := DB.Create(&admin).Error; err != nil {
		return err
	}

	userRole := model.UserRole{
		UserID:    admin.ID,
		RoleID:    superAdminRole.ID,
		CreatedAt: now,
	}
	if err := DB.Create(&userRole).Error; err != nil {
		return err
	}

	menus := []model.Menu{
		{TenantID: 0, ParentID: 0, Name: "Dashboard", Path: "/dashboard", Component: "BasicLayout", Icon: "lucide:layout-dashboard", Title: "仪表盘", Sort: 0, Status: 1},
		{TenantID: 0, ParentID: 1, Name: "Analytics", Path: "/dashboard/analytics", Component: "/dashboard/analytics/index", Icon: "lucide:area-chart", Title: "分析页", Sort: 0, Status: 1},
		{TenantID: 0, ParentID: 1, Name: "Workspace", Path: "/dashboard/workspace", Component: "/dashboard/workspace/index", Icon: "carbon:workspace", Title: "工作台", Sort: 1, Status: 1},
		{TenantID: 0, ParentID: 0, Name: "System", Path: "/system", Component: "BasicLayout", Icon: "lucide:settings", Title: "系统管理", Sort: 100, Status: 1},
		{TenantID: 0, ParentID: 4, Name: "UserManagement", Path: "/system/user", Component: "/system/user/index", Icon: "lucide:users", Title: "用户管理", Sort: 0, Status: 1},
		{TenantID: 0, ParentID: 4, Name: "RoleManagement", Path: "/system/role", Component: "/system/role/index", Icon: "lucide:shield", Title: "角色管理", Sort: 1, Status: 1},
		{TenantID: 0, ParentID: 4, Name: "MenuManagement", Path: "/system/menu", Component: "/system/menu/index", Icon: "lucide:menu", Title: "菜单管理", Sort: 2, Status: 1},
		{TenantID: 0, ParentID: 0, Name: "Content", Path: "/cms", Component: "BasicLayout", Icon: "lucide:file-text", Title: "内容管理", Sort: 50, Status: 1},
		{TenantID: 0, ParentID: 10, Name: "ArticleManagement", Path: "/cms/article", Component: "/cms/article/index", Icon: "lucide:newspaper", Title: "文章管理", Sort: 0, Status: 1},
		{TenantID: 0, ParentID: 10, Name: "CategoryManagement", Path: "/cms/category", Component: "/cms/category/index", Icon: "lucide:folder-tree", Title: "分类管理", Sort: 1, Status: 1},
		{TenantID: 0, ParentID: 10, Name: "TagManagement", Path: "/cms/tag", Component: "/cms/tag/index", Icon: "lucide:tags", Title: "标签管理", Sort: 2, Status: 1},
		{TenantID: 0, ParentID: 10, Name: "MediaManagement", Path: "/cms/media", Component: "/cms/media/index", Icon: "lucide:image", Title: "媒体管理", Sort: 3, Status: 1},
	}

	for _, menu := range menus {
		if err := DB.Create(&menu).Error; err != nil {
			return err
		}
	}

	// 初始化数据完成后修复一次菜单关系
	if err := fixMenuRelations(); err != nil {
		return err
	}

	// super_admin 拥有所有菜单
	var allMenus []model.Menu
	DB.Find(&allMenus)
	for _, menu := range allMenus {
		DB.Create(&model.RoleMenu{RoleID: superAdminRole.ID, MenuID: menu.ID, CreatedAt: now})
	}

	// admin 角色拥有所有菜单
	adminRole := roles[1]
	for _, menu := range allMenus {
		DB.Create(&model.RoleMenu{RoleID: adminRole.ID, MenuID: menu.ID, CreatedAt: now})
	}

	// user 仅拥有仪表盘和内容管理
	userRole2 := roles[2]
	for _, menu := range allMenus {
		if menu.Name == "Dashboard" || menu.Name == "Analytics" || menu.Name == "Workspace" ||
			menu.Name == "Content" || menu.Name == "ArticleManagement" ||
			menu.Name == "CategoryManagement" || menu.Name == "TagManagement" ||
			menu.Name == "MediaManagement" {
			DB.Create(&model.RoleMenu{RoleID: userRole2.ID, MenuID: menu.ID, CreatedAt: now})
		}
	}

	// 初始化 API 权限数据
	if err := initPermissions(superAdminRole.ID, adminRole.ID, userRole2.ID, now); err != nil {
		return err
	}

	return nil
}

func initPermissions(superAdminRoleID, adminRoleID, userRoleID uint, now time.Time) error {
	// 权限类型: 1=菜单 2=按钮 3=API
	permissions := []model.Permission{
		// ---- 用户管理 ----
		{Name: "用户列表", Code: "user:list", Type: 3, ParentID: 0, Path: "/api/users", Method: "GET"},
		{Name: "用户详情", Code: "user:detail", Type: 3, ParentID: 0, Path: "/api/users/:id", Method: "GET"},
		{Name: "创建用户", Code: "user:create", Type: 3, ParentID: 0, Path: "/api/users", Method: "POST"},
		{Name: "更新用户", Code: "user:update", Type: 3, ParentID: 0, Path: "/api/users/:id", Method: "PUT"},
		{Name: "删除用户", Code: "user:delete", Type: 3, ParentID: 0, Path: "/api/users/:id", Method: "DELETE"},
		{Name: "切换用户状态", Code: "user:status", Type: 3, ParentID: 0, Path: "/api/users/:id/status", Method: "PUT"},
		{Name: "重置密码", Code: "user:reset-password", Type: 3, ParentID: 0, Path: "/api/users/:id/reset-password", Method: "PUT"},
		{Name: "分配角色", Code: "user:assign-roles", Type: 3, ParentID: 0, Path: "/api/users/:id/roles", Method: "PUT"},

		// ---- 角色管理 ----
		{Name: "角色列表", Code: "role:list", Type: 3, ParentID: 0, Path: "/api/roles", Method: "GET"},
		{Name: "角色详情", Code: "role:detail", Type: 3, ParentID: 0, Path: "/api/roles/:id", Method: "GET"},
		{Name: "创建角色", Code: "role:create", Type: 3, ParentID: 0, Path: "/api/roles", Method: "POST"},
		{Name: "更新角色", Code: "role:update", Type: 3, ParentID: 0, Path: "/api/roles/:id", Method: "PUT"},
		{Name: "删除角色", Code: "role:delete", Type: 3, ParentID: 0, Path: "/api/roles/:id", Method: "DELETE"},
		{Name: "分配角色权限", Code: "role:assign-permissions", Type: 3, ParentID: 0, Path: "/api/roles/:id/permissions", Method: "PUT"},
		{Name: "查看角色权限", Code: "role:get-permissions", Type: 3, ParentID: 0, Path: "/api/roles/:id/permissions", Method: "GET"},
		{Name: "分配角色菜单", Code: "role:assign-menus", Type: 3, ParentID: 0, Path: "/api/roles/:id/menus", Method: "PUT"},
		{Name: "查看角色菜单", Code: "role:get-menus", Type: 3, ParentID: 0, Path: "/api/roles/:id/menus", Method: "GET"},

		// ---- 菜单管理 ----
		{Name: "菜单列表", Code: "menu:list", Type: 3, ParentID: 0, Path: "/api/menus", Method: "GET"},
		{Name: "菜单树", Code: "menu:tree", Type: 3, ParentID: 0, Path: "/api/menus/tree", Method: "GET"},
		{Name: "创建菜单", Code: "menu:create", Type: 3, ParentID: 0, Path: "/api/menus", Method: "POST"},
		{Name: "更新菜单", Code: "menu:update", Type: 3, ParentID: 0, Path: "/api/menus/:id", Method: "PUT"},
		{Name: "删除菜单", Code: "menu:delete", Type: 3, ParentID: 0, Path: "/api/menus/:id", Method: "DELETE"},


		// ---- 内容管理 ----
		{Name: "分类列表", Code: "category:list", Type: 3, ParentID: 0, Path: "/api/categories", Method: "GET"},
		{Name: "创建分类", Code: "category:create", Type: 3, ParentID: 0, Path: "/api/categories", Method: "POST"},
		{Name: "更新分类", Code: "category:update", Type: 3, ParentID: 0, Path: "/api/categories/:id", Method: "PUT"},
		{Name: "删除分类", Code: "category:delete", Type: 3, ParentID: 0, Path: "/api/categories/:id", Method: "DELETE"},

		{Name: "标签列表", Code: "tag:list", Type: 3, ParentID: 0, Path: "/api/tags", Method: "GET"},
		{Name: "创建标签", Code: "tag:create", Type: 3, ParentID: 0, Path: "/api/tags", Method: "POST"},
		{Name: "更新标签", Code: "tag:update", Type: 3, ParentID: 0, Path: "/api/tags/:id", Method: "PUT"},
		{Name: "删除标签", Code: "tag:delete", Type: 3, ParentID: 0, Path: "/api/tags/:id", Method: "DELETE"},

		{Name: "文章列表", Code: "article:list", Type: 3, ParentID: 0, Path: "/api/articles", Method: "GET"},
		{Name: "文章详情", Code: "article:detail", Type: 3, ParentID: 0, Path: "/api/articles/:id", Method: "GET"},
		{Name: "创建文章", Code: "article:create", Type: 3, ParentID: 0, Path: "/api/articles", Method: "POST"},
		{Name: "更新文章", Code: "article:update", Type: 3, ParentID: 0, Path: "/api/articles/:id", Method: "PUT"},
		{Name: "删除文章", Code: "article:delete", Type: 3, ParentID: 0, Path: "/api/articles/:id", Method: "DELETE"},
		{Name: "发布文章", Code: "article:publish", Type: 3, ParentID: 0, Path: "/api/articles/:id/publish", Method: "PUT"},
		{Name: "文章存草稿", Code: "article:draft", Type: 3, ParentID: 0, Path: "/api/articles/:id/draft", Method: "PUT"},

		{Name: "媒体列表", Code: "media:list", Type: 3, ParentID: 0, Path: "/api/media", Method: "GET"},
		{Name: "上传媒体", Code: "media:upload", Type: 3, ParentID: 0, Path: "/api/media/upload", Method: "POST"},
		{Name: "删除媒体", Code: "media:delete", Type: 3, ParentID: 0, Path: "/api/media/:id", Method: "DELETE"},

		// ---- 系统管理 ----
		{Name: "系统配置列表", Code: "config:list", Type: 3, ParentID: 0, Path: "/api/configs", Method: "GET"},
		{Name: "更新系统配置", Code: "config:update", Type: 3, ParentID: 0, Path: "/api/configs", Method: "PUT"},
		{Name: "操作日志", Code: "log:operation", Type: 3, ParentID: 0, Path: "/api/logs/operation", Method: "GET"},
		{Name: "登录日志", Code: "log:login", Type: 3, ParentID: 0, Path: "/api/logs/login", Method: "GET"},
		{Name: "权限列表", Code: "permission:list", Type: 3, ParentID: 0, Path: "/api/permissions", Method: "GET"},
		{Name: "权限树", Code: "permission:tree", Type: 3, ParentID: 0, Path: "/api/permissions/tree", Method: "GET"},
	}

	for i := range permissions {
		if err := DB.Create(&permissions[i]).Error; err != nil {
			return err
		}
	}

	// admin 角色拥有所有权限
	// user 角色仅拥有内容管理相关权限
	for _, perm := range permissions {
		DB.Create(&model.RolePermission{RoleID: adminRoleID, PermissionID: perm.ID, CreatedAt: now})

		// user: 仅内容管理 + 列表查看
		isContentPerm := (len(perm.Code) >= 9 && perm.Code[:9] == "category:") ||
			(len(perm.Code) >= 4 && perm.Code[:4] == "tag:") ||
			(len(perm.Code) >= 8 && perm.Code[:8] == "article:") ||
			(len(perm.Code) >= 6 && perm.Code[:6] == "media:")
		if isContentPerm {
			DB.Create(&model.RolePermission{RoleID: userRoleID, PermissionID: perm.ID, CreatedAt: now})
		}
	}

	return nil
}

func fixMenuRelations() error {
	var dashboard model.Menu
	var system model.Menu
	var content model.Menu

	DB.Where("name = ?", "Dashboard").First(&dashboard)
	DB.Where("name = ?", "System").First(&system)
	DB.Where("name = ?", "Content").First(&content)

	// 如果父菜单不存在则不处理
	if dashboard.ID == 0 || system.ID == 0 || content.ID == 0 {
		return nil
	}

	// Dashboard children
	DB.Model(&model.Menu{}).Where("name IN ?", []string{"Analytics", "Workspace"}).Updates(map[string]any{"parent_id": dashboard.ID})
	// System children
	DB.Model(&model.Menu{}).Where("name IN ?", []string{"UserManagement", "RoleManagement", "MenuManagement"}).Updates(map[string]any{"parent_id": system.ID})
	// CMS children
	DB.Model(&model.Menu{}).Where("name IN ?", []string{"ArticleManagement", "CategoryManagement", "TagManagement", "MediaManagement"}).Updates(map[string]any{"parent_id": content.ID})

	return nil
}
