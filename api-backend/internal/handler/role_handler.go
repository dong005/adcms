package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleHandler struct {
	roleRepo *repository.RoleRepository
}

func NewRoleHandler() *RoleHandler {
	return &RoleHandler{
		roleRepo: repository.NewRoleRepository(),
	}
}

type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Code        string `json:"code" binding:"required"`
	Description string `json:"description"`
	Status      int8   `json:"status"`
	Sort        int    `json:"sort"`
}

func (h *RoleHandler) Create(c *gin.Context) {
	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tenantID := middleware.GetTenantID(c)
	role := model.Role{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		Name:            req.Name,
		Code:            req.Code,
		Description:     req.Description,
		Status:          req.Status,
		Sort:            req.Sort,
	}

	if err := h.roleRepo.Create(&role); err != nil {
		utils.ServerError(c, "创建角色失败")
		return
	}

	utils.Success(c, role)
}

func (h *RoleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	role, err := h.roleRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 4001, "角色不存在")
		return
	}

	operatorID := middleware.GetUserID(c)

	// 不能修改比自己权限高或同级的角色
	if !middleware.CanOperateRole(operatorID, role.Code) {
		utils.Fail(c, 4003, "无权修改该角色")
		return
	}

	var req CreateRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 不允许通过修改code来提升角色等级
	if req.Code != role.Code {
		newLevel := middleware.RoleCodeToLevel(req.Code)
		operatorLevel := middleware.GetUserRoleLevel(operatorID)
		if newLevel <= operatorLevel {
			utils.Fail(c, 4003, "不能将角色编码修改为同级或更高级别")
			return
		}
	}

	role.Name = req.Name
	role.Code = req.Code
	role.Description = req.Description
	role.Status = req.Status
	role.Sort = req.Sort

	if err := h.roleRepo.Update(role); err != nil {
		utils.ServerError(c, "更新角色失败")
		return
	}

	utils.Success(c, role)
}

func (h *RoleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	role, err := h.roleRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 4001, "角色不存在")
		return
	}

	operatorID := middleware.GetUserID(c)

	// 不能删除比自己权限高或同级的角色
	if !middleware.CanOperateRole(operatorID, role.Code) {
		utils.Fail(c, 4003, "无权删除该角色")
		return
	}

	// 不能删除内置角色
	if role.Code == "super_admin" || role.Code == "admin" || role.Code == "user" {
		utils.Fail(c, 4003, "不能删除系统内置角色")
		return
	}

	if err := h.roleRepo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除角色失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *RoleHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	roles, err := h.roleRepo.List(tenantID)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}

	utils.Success(c, roles)
}

func (h *RoleHandler) Detail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	role, err := h.roleRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 4001, "角色不存在")
		return
	}

	utils.Success(c, role)
}

type AssignPermissionsRequest struct {
	PermissionIDs []uint `json:"permission_ids" binding:"required"`
}

func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	role, err := h.roleRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 4001, "角色不存在")
		return
	}

	operatorID := middleware.GetUserID(c)
	if !middleware.CanOperateRole(operatorID, role.Code) {
		utils.Fail(c, 4003, "无权修改该角色的权限")
		return
	}

	var req AssignPermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	if err := h.roleRepo.AssignPermissions(uint(id), req.PermissionIDs); err != nil {
		utils.ServerError(c, "分配权限失败")
		return
	}

	utils.SuccessWithMessage(c, "分配成功", nil)
}

func (h *RoleHandler) GetPermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	permissions, err := h.roleRepo.GetRolePermissions(uint(id))
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}

	utils.Success(c, permissions)
}
