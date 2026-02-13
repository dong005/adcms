package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct {
	adminRepo *repository.AdminRepository
	userRepo  *repository.UserRepository
}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{
		adminRepo: repository.NewAdminRepository(),
		userRepo:  repository.NewUserRepository(),
	}
}

type CreateAdminRequest struct {
	Username    string     `json:"username" binding:"required"`
	Password    string     `json:"password" binding:"required"`
	Email       string     `json:"email" binding:"email"`
	Phone       string     `json:"phone"`
	Nickname    string     `json:"nickname"`
	Company     string     `json:"company" binding:"required"`
	Domain      string     `json:"domain"`
	ExpireTime  *time.Time `json:"expire_time"`
	MaxUsers    uint       `json:"max_users"`
	Remark      string     `json:"remark"`
	Status      int8       `json:"status"`
}

func (h *AdminHandler) Create(c *gin.Context) {
	// 只有超管可以创建管理员
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Forbidden(c, "无权操作")
		return
	}

	var req CreateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 检查用户名是否已存在
	_, err := h.userRepo.FindByUsername(0, req.Username)
	if err == nil {
		utils.Fail(c, 3001, "用户名已存在")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	now := time.Now()
	user := model.User{
		TenantBaseModel: model.TenantBaseModel{TenantID: 0}, // 先设为0，创建后更新为自身ID
		Username:        req.Username,
		Password:        hashedPassword,
		Email:           req.Email,
		Phone:           req.Phone,
		Nickname:        req.Nickname,
		Status:          req.Status,
		IsAdmin:         1, // 管理员
		Company:         req.Company,
		Domain:          req.Domain,
		ExpireTime:      req.ExpireTime,
		MaxUsers:        req.MaxUsers,
		Remark:          req.Remark,
		LoginCount:      0,
		LastLoginAt:     &now,
	}

	if err := h.userRepo.Create(&user); err != nil {
		utils.ServerError(c, "创建管理员失败")
		return
	}

	// 设置租户ID为用户自身ID
	user.TenantID = user.ID
	if err := h.userRepo.Update(&user); err != nil {
		utils.ServerError(c, "更新租户信息失败")
		return
	}

	// 分配默认角色
	// 这里可以分配一个默认的管理员角色

	utils.Success(c, user)
}

type UpdateAdminRequest struct {
	Email      string     `json:"email" binding:"email"`
	Phone      string     `json:"phone"`
	Nickname   string     `json:"nickname"`
	Company    string     `json:"company"`
	Domain     string     `json:"domain"`
	ExpireTime *time.Time `json:"expire_time"`
	MaxUsers   uint       `json:"max_users"`
	Remark     string     `json:"remark"`
	Status     int8       `json:"status"`
}

func (h *AdminHandler) Update(c *gin.Context) {
	// 只有超管可以修改管理员信息
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Forbidden(c, "无权操作")
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

	// 只能修改管理员（is_admin=1）
	if user.IsAdmin != 1 {
		utils.Fail(c, 4003, "该用户不是管理员")
		return
	}

	var req UpdateAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	user.Email = req.Email
	user.Phone = req.Phone
	user.Nickname = req.Nickname
	user.Company = req.Company
	user.Domain = req.Domain
	user.ExpireTime = req.ExpireTime
	user.MaxUsers = req.MaxUsers
	user.Remark = req.Remark
	user.Status = req.Status

	if err := h.userRepo.Update(user); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}

	utils.Success(c, user)
}

func (h *AdminHandler) Delete(c *gin.Context) {
	// 只有超管可以删除管理员
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Forbidden(c, "无权操作")
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

	// 只能删除管理员（is_admin=1）
	if user.IsAdmin != 1 {
		utils.Fail(c, 4003, "该用户不是管理员")
		return
	}

	// TODO: 检查该租户下是否还有其他用户

	if err := h.userRepo.Delete(uint(id)); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}

	utils.SuccessWithMessage(c, "删除成功", nil)
}

func (h *AdminHandler) List(c *gin.Context) {
	// 只有超管可以查看管理员列表
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Forbidden(c, "无权操作")
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	keyword := c.Query("keyword")

	admins, total, err := h.adminRepo.List(page, pageSize, keyword)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}

	utils.SuccessWithPage(c, admins, total, page, pageSize)
}

func (h *AdminHandler) Detail(c *gin.Context) {
	// 只有超管可以查看管理员详情
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Forbidden(c, "无权操作")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	admin, err := h.adminRepo.Detail(uint(id))
	if err != nil {
		utils.Fail(c, 3002, "管理员不存在")
		return
	}

	utils.Success(c, admin)
}

func (h *AdminHandler) ToggleStatus(c *gin.Context) {
	// 只有超管可以禁用/启用管理员
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Forbidden(c, "无权操作")
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

	if user.IsAdmin != 1 {
		utils.Fail(c, 4003, "该用户不是管理员")
		return
	}

	// 切换状态
	if user.Status == 1 {
		user.Status = 0
	} else {
		user.Status = 1
	}

	if err := h.userRepo.Update(user); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "状态更新成功", user)
}

func (h *AdminHandler) ResetPassword(c *gin.Context) {
	// 只有超管可以重置管理员密码
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Forbidden(c, "无权操作")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	var req struct {
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	user, err := h.userRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 3002, "用户不存在")
		return
	}

	if user.IsAdmin != 1 {
		utils.Fail(c, 4003, "该用户不是管理员")
		return
	}

	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	user.Password = hashedPassword
	if err := h.userRepo.Update(user); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "密码重置成功", nil)
}

func (h *AdminHandler) Statistics(c *gin.Context) {
	// 只有超管可以查看统计信息
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Forbidden(c, "无权操作")
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	stats, err := h.adminRepo.Statistics(uint(id))
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}

	utils.Success(c, stats)
}
