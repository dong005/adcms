package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type MenuHandler struct {
	menuRepo *repository.MenuRepository
	roleRepo *repository.RoleRepository
	userRepo *repository.UserRepository
}

func NewMenuHandler() *MenuHandler {
	return &MenuHandler{
		menuRepo: repository.NewMenuRepository(),
		roleRepo: repository.NewRoleRepository(),
		userRepo: repository.NewUserRepository(),
	}
}

type CreateMenuRequest struct {
	ParentID         uint   `json:"parent_id"`
	Name             string `json:"name" binding:"required"`
	Path             string `json:"path"`
	Component        string `json:"component"`
	Redirect         string `json:"redirect"`
	Icon             string `json:"icon"`
	Title            string `json:"title" binding:"required"`
	HideInMenu       int8   `json:"hide_in_menu"`
	HideInTab        int8   `json:"hide_in_tab"`
	HideInBreadcrumb int8   `json:"hide_in_breadcrumb"`
	KeepAlive        int8   `json:"keep_alive"`
	FrameSrc         string `json:"frame_src"`
	Sort             int    `json:"sort"`
	Status           int8   `json:"status"`
	PermissionCode   string `json:"permission_code"`
	// 新增字段
	IsTenant         int8   `json:"is_tenant"`      // 1=租户可见 0=仅超管
	IsPublic         int8   `json:"is_public"`      // 1=公共(不需权限)
	Type             int8   `json:"type"`           // 1=目录 2=菜单 3=页面 4=按钮/权限
}

func (h *MenuHandler) Create(c *gin.Context) {
	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tenantID := middleware.GetTenantID(c)
	menu := model.Menu{
		TenantID:         tenantID,
		ParentID:         req.ParentID,
		Name:             req.Name,
		Path:             req.Path,
		Component:        req.Component,
		Redirect:         req.Redirect,
		Icon:             req.Icon,
		Title:            req.Title,
		HideInMenu:       req.HideInMenu,
		HideInTab:        req.HideInTab,
		HideInBreadcrumb: req.HideInBreadcrumb,
		KeepAlive:        req.KeepAlive,
		FrameSrc:         req.FrameSrc,
		Sort:             req.Sort,
		Status:           req.Status,
		PermissionCode:   req.PermissionCode,
		IsTenant:         req.IsTenant,
		IsPublic:         req.IsPublic,
		Type:             req.Type,
	}

	if err := h.menuRepo.Create(&menu); err != nil {
		utils.ServerError(c, "创建菜单失败")
		return
	}

	utils.Success(c, menu)
}

func (h *MenuHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	menu, err := h.menuRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 2001, "菜单不存在")
		return
	}

	var req CreateMenuRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	menu.ParentID = req.ParentID
	menu.Name = req.Name
	menu.Path = req.Path
	menu.Component = req.Component
	menu.Redirect = req.Redirect
	menu.Icon = req.Icon
	menu.Title = req.Title
	menu.HideInMenu = req.HideInMenu
	menu.HideInTab = req.HideInTab
	menu.HideInBreadcrumb = req.HideInBreadcrumb
	menu.KeepAlive = req.KeepAlive
	menu.FrameSrc = req.FrameSrc
	menu.Sort = req.Sort
	menu.Status = req.Status
	menu.PermissionCode = req.PermissionCode
	menu.IsTenant = req.IsTenant
	menu.IsPublic = req.IsPublic
	menu.Type = req.Type

	if err := h.menuRepo.Update(menu); err != nil {
		utils.ServerError(c, "更新菜单失败")
		return
	}

	utils.Success(c, menu)
}

func (h *MenuHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	hasChildren, err := h.menuRepo.HasChildren(uint(id))
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	if hasChildren {
		utils.Fail(c, 2002, "请先删除子菜单")
		return
	}

	if err := h.menuRepo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除菜单失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *MenuHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	isAdmin := middleware.IsAdmin(middleware.GetUserID(c))
	menus, err := h.menuRepo.FindAll(tenantID, isAdmin)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}

	tree := repository.BuildMenuTree(menus, 0)
	utils.Success(c, tree)
}

func (h *MenuHandler) Tree(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	isAdmin := middleware.IsAdmin(middleware.GetUserID(c))
	menus, err := h.menuRepo.FindAll(tenantID, isAdmin)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}

	tree := repository.ConvertToMenuTree(menus)
	utils.Success(c, tree)
}

func (h *MenuHandler) UserMenus(c *gin.Context) {
	userID := middleware.GetUserID(c)

	roles, err := h.userRepo.GetUserRoles(userID)
	if err != nil {
		utils.ServerError(c, "查询角色失败")
		return
	}

	var roleIDs []uint
	isSuperAdmin := false
	for _, role := range roles {
		roleIDs = append(roleIDs, role.ID)
		if role.Code == "super_admin" {
			isSuperAdmin = true
		}
	}

	var menus []model.Menu
	if isSuperAdmin {
		tenantID := middleware.GetTenantID(c)
		menus, err = h.menuRepo.FindAll(tenantID, true)
	} else {
		menus, err = h.menuRepo.FindByRoleIDs(roleIDs)
	}

	if err != nil {
		utils.ServerError(c, "查询菜单失败")
		return
	}

	tree := repository.ConvertToMenuTree(menus)
	utils.Success(c, tree)
}

type AssignMenusRequest struct {
	MenuIDs []uint `json:"menu_ids" binding:"required"`
}

func (h *MenuHandler) AssignRoleMenus(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	role, err := h.roleRepo.FindByID(uint(roleID))
	if err != nil {
		utils.Fail(c, 4001, "角色不存在")
		return
	}

	operatorID := middleware.GetUserID(c)
	if !middleware.CanOperateRole(operatorID, role.Code) {
		utils.Fail(c, 4003, "无权修改该角色的菜单")
		return
	}

	var req AssignMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.roleRepo.AssignMenus(uint(roleID), req.MenuIDs); err != nil {
		utils.ServerError(c, "分配菜单失败")
		return
	}

	utils.SuccessWithMessage(c, "分配成功", nil)
}

// GetUserMenus 获取用户的独立菜单权限
func (h *MenuHandler) GetUserMenus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	menus, err := h.userRepo.GetUserMenus(uint(id))
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}

	utils.Success(c, menus)
}

// AssignUserMenus 分配用户的独立菜单权限
func (h *MenuHandler) AssignUserMenus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req AssignMenusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.userRepo.AssignMenus(uint(id), req.MenuIDs); err != nil {
		utils.ServerError(c, "分配失败")
		return
	}

	utils.SuccessWithMessage(c, "分配成功", nil)
}

func (h *MenuHandler) GetRoleMenus(c *gin.Context) {
	roleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	menus, err := h.roleRepo.GetRoleMenus(uint(roleID))
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}

	utils.Success(c, menus)
}
