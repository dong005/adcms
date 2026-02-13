package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/database"
	"adcms/pkg/excel"
	"adcms/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo *repository.UserRepository
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		userRepo: repository.NewUserRepository(),
	}
}

type CreateUserRequest struct {
	Username    string     `json:"username" binding:"required"`
	Password    string     `json:"password" binding:"required"`
	Email       string     `json:"email" binding:"email"`
	Phone       string     `json:"phone"`
	Nickname    string     `json:"nickname"`
	Avatar      string     `json:"avatar"`
	DepartmentID uint       `json:"department_id"`
	Status      int8       `json:"status"`
	RoleIDs     []uint     `json:"role_ids" binding:"required"`
	IsAdmin     int8       `json:"is_admin"`                    // 0=普通用户 1=管理员(租户) 2=超级管理员
	Company     string     `json:"company"`                     // 租户公司名称
	Domain      string     `json:"domain"`                      // 租户绑定域名
	ExpireTime  *time.Time `json:"expire_time"`                 // 租户到期时间
	MaxUsers    uint       `json:"max_users"`                   // 最大用户数，0=不限
	Remark      string     `json:"remark"`                      // 备注
}

func (h *UserHandler) Create(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)

	// 检查是否有权分配这些角色
	if len(req.RoleIDs) > 0 && !middleware.CanAssignRoles(operatorID, req.RoleIDs) {
		utils.Fail(c, 4003, "无权分配该角色，不能分配与自己同级或更高级别的角色")
		return
	}

	tenantID := middleware.GetTenantID(c)

	_, err := h.userRepo.FindByUsername(tenantID, req.Username)
	if err == nil {
		utils.Fail(c, 3001, "用户名已存在")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	user := model.User{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		Username:        req.Username,
		Password:        hashedPassword,
		Email:           req.Email,
		Phone:           req.Phone,
		Nickname:        req.Nickname,
		Avatar:          req.Avatar,
		DepartmentID:    req.DepartmentID,
		Status:          req.Status,
		IsAdmin:         req.IsAdmin,
		Company:         req.Company,
		Domain:          req.Domain,
		ExpireTime:      req.ExpireTime,
		MaxUsers:        req.MaxUsers,
		Remark:          req.Remark,
	}

	if err := h.userRepo.Create(&user); err != nil {
		utils.ServerError(c, "创建用户失败")
		return
	}

	if len(req.RoleIDs) > 0 {
		h.userRepo.AssignRoles(user.ID, req.RoleIDs)
	}

	// 如果超管创建的是管理员（is_admin=1），tenant_id 设为新用户自己的 ID
	if tenantID == 0 && req.IsAdmin == 1 {
		user.TenantID = user.ID
		h.userRepo.Update(&user)
	}

	utils.Success(c, user)
}

type UpdateUserRequest struct {
	Email       string     `json:"email" binding:"email"`
	Phone       string     `json:"phone"`
	Nickname    string     `json:"nickname"`
	Avatar      string     `json:"avatar"`
	DepartmentID uint       `json:"department_id"`
	Status      int8       `json:"status"`
	RoleIDs     []uint     `json:"role_ids"`
	// 新增字段
	Company     string     `json:"company"`                     // 租户公司名称
	Domain      string     `json:"domain"`                      // 租户绑定域名
	ExpireTime  *time.Time `json:"expire_time"`                 // 租户到期时间
	MaxUsers    uint       `json:"max_users"`                   // 最大用户数，0=不限
	Remark      string     `json:"remark"`                      // 备注
}

func (h *UserHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	targetID := uint(id)

	// 不能修改权限等级 >= 自己的用户（除非是自己）
	if operatorID != targetID && !middleware.HasHigherLevel(operatorID, targetID) {
		utils.Fail(c, 4003, "无权修改该用户，目标用户权限等级不低于您")
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 3002, "用户不存在")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	user.Email = req.Email
	user.Phone = req.Phone
	user.Nickname = req.Nickname
	user.Avatar = req.Avatar
	user.DepartmentID = req.DepartmentID
	user.Status = req.Status
	user.Company = req.Company
	user.Domain = req.Domain
	user.ExpireTime = req.ExpireTime
	user.MaxUsers = req.MaxUsers
	user.Remark = req.Remark

	if err := h.userRepo.Update(user); err != nil {
		utils.ServerError(c, "更新用户失败")
		return
	}

	// 更新角色
	if req.RoleIDs != nil {
		h.userRepo.AssignRoles(user.ID, req.RoleIDs)
	}

	utils.Success(c, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	targetID := uint(id)

	// 不能删除自己
	if operatorID == targetID {
		utils.Fail(c, 4003, "不能删除自己")
		return
	}

	// 不能删除权限等级 >= 自己的用户
	if !middleware.HasHigherLevel(operatorID, targetID) {
		utils.Fail(c, 4003, "无权删除该用户，目标用户权限等级不低于您")
		return
	}

	if err := h.userRepo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除用户失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *UserHandler) List(c *gin.Context) {
	isAdmin := middleware.GetIsAdmin(c)
	
	// 超级管理员可以看到所有用户，其他用户只能看到自己租户的用户
	var tenantID uint
	if isAdmin == 2 { // 超级管理员
		tenantID = 0 // 查询所有租户
	} else {
		tenantID = middleware.GetTenantID(c)
	}
	
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	users, total, err := h.userRepo.List(tenantID, page, pageSize, keyword)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}

	utils.SuccessWithPage(c, users, total, page, pageSize)
}

func (h *UserHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 3002, "用户不存在")
		return
	}

	utils.Success(c, user)
}

func (h *UserHandler) ToggleStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	targetID := uint(id)

	if operatorID == targetID {
		utils.Fail(c, 4003, "不能禁用自己")
		return
	}

	if !middleware.HasHigherLevel(operatorID, targetID) {
		utils.Fail(c, 4003, "无权修改该用户状态")
		return
	}

	var req struct {
		Status int8 `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.userRepo.UpdateStatus(uint(id), req.Status); err != nil {
		utils.ServerError(c, "更新状态失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

func (h *UserHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	targetID := uint(id)

	// 可以重置自己的密码，或者权限更高的用户可以重置
	if operatorID != targetID && !middleware.HasHigherLevel(operatorID, targetID) {
		utils.Fail(c, 4003, "无权重置该用户密码")
		return
	}

	hashedPassword, err := utils.HashPassword("123456")
	if err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	if err := h.userRepo.UpdatePassword(uint(id), hashedPassword); err != nil {
		utils.ServerError(c, "重置密码失败")
		return
	}

	utils.SuccessWithMessage(c, "密码已重置为123456", nil)
}

type AssignRolesRequest struct {
	RoleIDs []uint `json:"role_ids" binding:"required"`
}

func (h *UserHandler) AssignRoles(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	operatorID := middleware.GetUserID(c)
	targetID := uint(id)

	// 不能修改权限等级 >= 自己的用户的角色（除非是超级管理员）
	if operatorID != targetID && !middleware.HasHigherLevel(operatorID, targetID) {
		utils.Fail(c, 4003, "无权修改该用户的角色")
		return
	}

	var req AssignRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 检查是否有权分配这些角色
	if !middleware.CanAssignRoles(operatorID, req.RoleIDs) {
		utils.Fail(c, 4003, "无权分配该角色，不能分配与自己同级或更高级别的角色")
		return
	}

	if err := h.userRepo.AssignRoles(uint(id), req.RoleIDs); err != nil {
		utils.ServerError(c, "分配角色失败")
		return
	}

	utils.SuccessWithMessage(c, "分配成功", nil)
}

// isAdminRole 检查角色ID列表中是否包含 admin 角色
func isAdminRole(roleIDs []uint) bool {
	for _, id := range roleIDs {
		var role model.Role
		if err := database.DB.First(&role, id).Error; err == nil {
			if role.Code == "admin" {
				return true
			}
		}
	}
	return false
}

var userExcelColumns = []excel.ColumnDef{
	{Header: "用户名", Field: "Username", Width: 15},
	{Header: "昵称", Field: "Nickname", Width: 15},
	{Header: "邮箱", Field: "Email", Width: 25},
	{Header: "手机", Field: "Phone", Width: 15},
	{Header: "状态", Field: "Status", Width: 8},
	{Header: "创建时间", Field: "CreatedAt", Width: 20},
}

func (h *UserHandler) Export(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	keyword := c.Query("keyword")
	users, _, err := h.userRepo.List(tenantID, 1, 10000, keyword)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	excel.Export(c, "用户列表.xlsx", "用户", userExcelColumns, users)
}

func (h *UserHandler) ImportTemplate(c *gin.Context) {
	excel.ExportTemplate(c, "用户导入模板.xlsx", "用户", userExcelColumns)
}

func (h *UserHandler) Import(c *gin.Context) {
	records, err := excel.Import(c, "file", userExcelColumns)
	if err != nil {
		utils.BadRequest(c, err.Error())
		return
	}

	tenantID := middleware.GetTenantID(c)
	imported := 0
	for _, rec := range records {
		username := rec["Username"]
		if username == "" {
			continue
		}
		password, _ := utils.HashPassword("123456")
		user := model.User{
			TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
			Username:        username,
			Password:        password,
			Nickname:        rec["Nickname"],
			Email:           rec["Email"],
			Phone:           rec["Phone"],
			Status:          1,
		}
		if err := h.userRepo.Create(&user); err == nil {
			imported++
		}
	}
	utils.Success(c, map[string]int{"imported": imported, "total": len(records)})
}

func (h *UserHandler) AssignMenus(c *gin.Context) {
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
		utils.ServerError(c, "分配菜单失败")
		return
	}

	utils.SuccessWithMessage(c, "分配成功", nil)
}

// UnlockUser 解锁因长期未登录被锁定的用户（仅超管）
func (h *UserHandler) UnlockUser(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Forbidden(c, "仅超级管理员可解锁用户")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 3002, "用户不存在")
		return
	}

	if user.Status != 2 {
		utils.Fail(c, 4004, "该用户未被锁定")
		return
	}

	if err := h.userRepo.UpdateStatus(uint(id), 1); err != nil {
		utils.ServerError(c, "解锁失败")
		return
	}

	utils.SuccessWithMessage(c, "用户已解锁", nil)
}
