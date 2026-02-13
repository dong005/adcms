package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/database"
	"adcms/pkg/email"
	"adcms/pkg/logcfg"
	"adcms/pkg/sms"
	"adcms/pkg/utils"
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userRepo *repository.UserRepository
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		userRepo: repository.NewUserRepository(),
	}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token       string `json:"token,omitempty"`
	TempToken   string `json:"temp_token,omitempty"`
	RequireTotp bool   `json:"require_totp"`
}

// Login 用户登录
// @Summary 用户登录
// @Tags 认证
// @Accept json
// @Produce json
// @Param body body LoginRequest true "登录参数"
// @Success 200 {object} LoginResponse
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	// 检查账号是否被锁定
	if locked, remainSec := middleware.IsLoginLocked(req.Username); locked {
		utils.Fail(c, 1011, fmt.Sprintf("账号已被锁定，请%d秒后再试", remainSec))
		return
	}

	user, err := h.userRepo.FindByUsernameGlobal(req.Username)
	if err != nil {
		remaining, locked := middleware.RecordLoginFail(req.Username)
		h.recordLoginLog(0, 0, req.Username, c.ClientIP(), c.Request.UserAgent(), 0, "用户不存在")
		if locked {
			utils.Fail(c, 1011, "登录失败次数过多，账号已被锁定15分钟")
		} else {
			utils.Fail(c, 1001, fmt.Sprintf("用户名或密码错误，还可尝试%d次", remaining))
		}
		return
	}

	if !utils.ComparePassword(user.Password, req.Password) {
		remaining, locked := middleware.RecordLoginFail(req.Username)
		h.recordLoginLog(user.TenantID, user.ID, req.Username, c.ClientIP(), c.Request.UserAgent(), 0, "密码错误")
		if locked {
			utils.Fail(c, 1011, "登录失败次数过多，账号已被锁定15分钟")
		} else {
			utils.Fail(c, 1001, fmt.Sprintf("用户名或密码错误，还可尝试%d次", remaining))
		}
		return
	}

	if user.Status == 0 {
		h.recordLoginLog(user.TenantID, user.ID, req.Username, c.ClientIP(), c.Request.UserAgent(), 0, "用户已禁用")
		utils.Fail(c, 1002, "用户已被禁用")
		return
	}
	if user.Status == 2 {
		h.recordLoginLog(user.TenantID, user.ID, req.Username, c.ClientIP(), c.Request.UserAgent(), 0, "用户已锁定")
		utils.Fail(c, 1012, "账号已被锁定，请联系管理员解锁")
		return
	}

	// 动态检查：用户长期未登录自动锁定（从 config_webs 读取天数，0=不锁定）
	if user.IsAdmin != 2 && user.LastLoginAt != nil {
		lockDays := h.getInactiveLockDays()
		if lockDays > 0 {
			threshold := time.Now().AddDate(0, 0, -lockDays)
			if user.LastLoginAt.Before(threshold) {
				// 自动将状态设为锁定
				h.userRepo.UpdateStatus(user.ID, 2)
				h.recordLoginLog(user.TenantID, user.ID, req.Username, c.ClientIP(), c.Request.UserAgent(), 0, fmt.Sprintf("超过%d天未登录自动锁定", lockDays))
				utils.Fail(c, 1012, fmt.Sprintf("账号因超过%d天未登录已被锁定，请联系管理员解锁", lockDays))
				return
			}
		}
	}

	if user.TOTPEnabled == 1 {
		tempToken, err := utils.GenerateTempToken(user.ID, user.TenantID, user.Username)
		if err != nil {
			utils.ServerError(c, "生成token失败")
			return
		}
		utils.Success(c, LoginResponse{
			TempToken:   tempToken,
			RequireTotp: true,
		})
		return
	}

	token, err := utils.GenerateToken(user.ID, user.TenantID, user.Username, user.IsAdmin)
	if err != nil {
		utils.ServerError(c, "生成token失败")
		return
	}

	// 登录成功，清除失败记录
	middleware.ClearLoginFail(req.Username)
	h.userRepo.UpdateLoginInfo(user.ID, c.ClientIP())
	h.recordLoginLog(user.TenantID, user.ID, req.Username, c.ClientIP(), c.Request.UserAgent(), 1, "登录成功")

	utils.Success(c, LoginResponse{
		Token:       token,
		RequireTotp: false,
	})
}

type VerifyTOTPRequest struct {
	TempToken string `json:"temp_token" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

func (h *AuthHandler) VerifyTOTP(c *gin.Context) {
	var req VerifyTOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	claims, err := utils.ParseTempToken(req.TempToken)
	if err != nil || !claims.RequireOTP {
		utils.Fail(c, 1005, "临时token无效")
		return
	}

	user, err := h.userRepo.FindByID(claims.UserID)
	if err != nil {
		utils.Fail(c, 1006, "用户不存在")
		return
	}

	if !utils.ValidateTOTPCode(user.TOTPSecret, req.Code) {
		utils.Fail(c, 1007, "验证码错误")
		return
	}

	token, err := utils.GenerateToken(user.ID, user.TenantID, user.Username, user.IsAdmin)
	if err != nil {
		utils.ServerError(c, "生成token失败")
		return
	}

	h.userRepo.UpdateLoginInfo(user.ID, c.ClientIP())
	h.recordLoginLog(user.TenantID, user.ID, user.Username, c.ClientIP(), c.Request.UserAgent(), 1, "登录成功(TOTP)")

	utils.Success(c, LoginResponse{
		Token:       token,
		RequireTotp: false,
	})
}

type GenerateTOTPResponse struct {
	Secret string `json:"secret"`
	QRCode string `json:"qr_code"`
}

func (h *AuthHandler) GenerateTOTP(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		utils.Fail(c, 1006, "用户不存在")
		return
	}

	key, err := utils.GenerateTOTPSecret("ADCMS", user.Username)
	if err != nil {
		utils.ServerError(c, "生成密钥失败")
		return
	}

	qrCode, err := utils.GenerateTOTPQRCode(key)
	if err != nil {
		utils.ServerError(c, "生成二维码失败")
		return
	}

	ctx := context.Background()
	database.RDB.Set(ctx, "totp:temp:"+user.Username, key.Secret(), 5*time.Minute)

	utils.Success(c, GenerateTOTPResponse{
		Secret: key.Secret(),
		QRCode: "data:image/png;base64," + qrCode,
	})
}

type BindTOTPRequest struct {
	Code string `json:"code" binding:"required"`
}

func (h *AuthHandler) BindTOTP(c *gin.Context) {
	var req BindTOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		utils.Fail(c, 1006, "用户不存在")
		return
	}

	ctx := context.Background()
	secret, err := database.RDB.Get(ctx, "totp:temp:"+user.Username).Result()
	if err != nil {
		utils.Fail(c, 1008, "请先生成二维码")
		return
	}

	if !utils.ValidateTOTPCode(secret, req.Code) {
		utils.Fail(c, 1007, "验证码错误")
		return
	}

	if err := h.userRepo.UpdateTOTP(userID, 1, secret); err != nil {
		utils.ServerError(c, "绑定失败")
		return
	}

	database.RDB.Del(ctx, "totp:temp:"+user.Username)

	utils.SuccessWithMessage(c, "绑定成功", nil)
}

func (h *AuthHandler) DisableTOTP(c *gin.Context) {
	var req BindTOTPRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		utils.Fail(c, 1006, "用户不存在")
		return
	}

	if user.TOTPEnabled != 1 {
		utils.Fail(c, 1009, "未启用TOTP")
		return
	}

	if !utils.ValidateTOTPCode(user.TOTPSecret, req.Code) {
		utils.Fail(c, 1007, "验证码错误")
		return
	}

	if err := h.userRepo.UpdateTOTP(userID, 0, ""); err != nil {
		utils.ServerError(c, "解绑失败")
		return
	}

	utils.SuccessWithMessage(c, "解绑成功", nil)
}

type UserInfoResponse struct {
	ID          uint     `json:"id"`
	Username    string   `json:"username"`
	Nickname    string   `json:"nickname"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	Avatar      string   `json:"avatar"`
	TenantID    uint     `json:"tenant_id"`
	TOTPEnabled bool     `json:"totp_enabled"`
	EmailNotify int8     `json:"email_notify"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`
}

func (h *AuthHandler) GetUserInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		utils.Fail(c, 1006, "用户不存在")
		return
	}

	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Code)
	}

	// 获取用户权限码列表
	permissions := middleware.GetUserPermissionCodes(userID)
	if permissions == nil {
		permissions = []string{}
	}

	utils.Success(c, UserInfoResponse{
		ID:          user.ID,
		Username:    user.Username,
		Nickname:    user.Nickname,
		Email:       user.Email,
		Phone:       user.Phone,
		Avatar:      user.Avatar,
		TenantID:    user.TenantID,
		TOTPEnabled: user.TOTPEnabled == 1,
		EmailNotify: user.EmailNotify,
		Roles:       roleNames,
		Permissions: permissions,
	})
}

type UpdateUserInfoRequest struct {
	Nickname    string `json:"nickname"`
	Email       string `json:"email"`
	Avatar      string `json:"avatar"`
	EmailNotify *int8  `json:"email_notify"`
}

func (h *AuthHandler) UpdateUserInfo(c *gin.Context) {
	var req UpdateUserInfoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		utils.Fail(c, 1006, "用户不存在")
		return
	}

	user.Nickname = req.Nickname
	user.Email = req.Email
	user.Avatar = req.Avatar
	if req.EmailNotify != nil {
		user.EmailNotify = *req.EmailNotify
	}

	if err := h.userRepo.Update(user); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		utils.Fail(c, 1006, "用户不存在")
		return
	}

	if !utils.ComparePassword(user.Password, req.OldPassword) {
		utils.Fail(c, 1010, "原密码错误")
		return
	}

	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	if err := h.userRepo.UpdatePassword(userID, hashedPassword); err != nil {
		utils.ServerError(c, "修改密码失败")
		return
	}

	utils.SuccessWithMessage(c, "密码修改成功", nil)
}

func (h *AuthHandler) Logout(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader != "" {
		token := authHeader[7:]
		middleware.BlacklistToken(token, 24*time.Hour)
	}

	utils.SuccessWithMessage(c, "登出成功", nil)
}

// GetPermissionCodes 获取当前用户的权限码列表（供前端按钮级权限使用）
func (h *AuthHandler) GetPermissionCodes(c *gin.Context) {
	userID := middleware.GetUserID(c)
	codes := middleware.GetUserPermissionCodes(userID)
	if codes == nil {
		codes = []string{}
	}
	utils.Success(c, codes)
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请输入有效的邮箱地址")
		return
	}

	ctx := context.Background()

	// 限流：同一邮箱1分钟内只能发1次
	rateLimitKey := "pwd_reset_limit:" + req.Email
	if database.RDB.Exists(ctx, rateLimitKey).Val() > 0 {
		utils.Fail(c, 1020, "发送过于频繁，请1分钟后再试")
		return
	}

	// 查找用户
	user, err := h.userRepo.FindByEmailGlobal(req.Email)
	if err != nil {
		// 为防止邮箱枚举攻击，即使用户不存在也返回成功
		utils.SuccessWithMessage(c, "如果该邮箱已注册，验证码将发送到您的邮箱", nil)
		return
	}

	// 生成6位验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 存入Redis，5分钟过期
	codeKey := "pwd_reset:" + req.Email
	database.RDB.Set(ctx, codeKey, code, 5*time.Minute)
	database.RDB.Set(ctx, rateLimitKey, "1", time.Minute)

	// 异步发送邮件
	go func() {
		if err := email.SendResetCode(req.Email, code); err != nil {
			fmt.Printf("[Email] 发送重置验证码失败: %v\n", err)
		}
	}()

	_ = user
	utils.SuccessWithMessage(c, "如果该邮箱已注册，验证码将发送到您的邮箱", nil)
}

type ResetPasswordByEmailRequest struct {
	Email       string `json:"email" binding:"required,email"`
	Code        string `json:"code" binding:"required,len=6"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *AuthHandler) ResetPasswordByEmail(c *gin.Context) {
	var req ResetPasswordByEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	ctx := context.Background()
	codeKey := "pwd_reset:" + req.Email

	// 验证验证码
	storedCode, err := database.RDB.Get(ctx, codeKey).Result()
	if err != nil || storedCode != req.Code {
		utils.Fail(c, 1021, "验证码错误或已过期")
		return
	}

	// 查找用户
	user, err := h.userRepo.FindByEmailGlobal(req.Email)
	if err != nil {
		utils.Fail(c, 1006, "该邮箱未注册")
		return
	}

	// 重置密码
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		utils.ServerError(c, "密码加密失败")
		return
	}

	if err := h.userRepo.UpdatePassword(user.ID, hashedPassword); err != nil {
		utils.ServerError(c, "重置密码失败")
		return
	}

	// 删除验证码
	database.RDB.Del(ctx, codeKey)

	utils.SuccessWithMessage(c, "密码重置成功", nil)
}

// LoginHistory 获取当前用户最近登录记录
func (h *AuthHandler) LoginHistory(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var logs []model.LoginLog
	database.DB.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(20).
		Find(&logs)

	utils.Success(c, logs)
}

// SendSmsCode 发送手机验证码
type SendSmsCodeRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (h *AuthHandler) SendSmsCode(c *gin.Context) {
	var req SendSmsCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请输入手机号")
		return
	}

	userID := middleware.GetUserID(c)
	ctx := context.Background()

	// 限流：同一用户1分钟内只能发1次
	rateLimitKey := fmt.Sprintf("sms_limit:%d", userID)
	if database.RDB.Exists(ctx, rateLimitKey).Val() > 0 {
		utils.Fail(c, 1020, "发送过于频繁，请1分钟后再试")
		return
	}

	// 生成6位验证码
	code := fmt.Sprintf("%06d", rand.Intn(1000000))

	// 存入Redis，5分钟过期
	codeKey := fmt.Sprintf("sms_code:%d", userID)
	database.RDB.Set(ctx, codeKey, code+":"+req.Phone, 5*time.Minute)
	database.RDB.Set(ctx, rateLimitKey, "1", time.Minute)

	// 异步发送短信
	go func() {
		if err := sms.SendVerifyCode(req.Phone, code); err != nil {
			fmt.Printf("[SMS] 发送验证码失败: %v\n", err)
		}
	}()

	utils.SuccessWithMessage(c, "验证码已发送", nil)
}

// BindPhone 绑定手机号
type BindPhoneRequest struct {
	Phone string `json:"phone" binding:"required"`
	Code  string `json:"code" binding:"required,len=6"`
}

func (h *AuthHandler) BindPhone(c *gin.Context) {
	var req BindPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	userID := middleware.GetUserID(c)
	ctx := context.Background()
	codeKey := fmt.Sprintf("sms_code:%d", userID)

	// 验证验证码
	stored, err := database.RDB.Get(ctx, codeKey).Result()
	if err != nil {
		utils.Fail(c, 1021, "验证码已过期，请重新发送")
		return
	}

	// stored 格式: "code:phone"
	parts := fmt.Sprintf("%s:%s", req.Code, req.Phone)
	if stored != parts {
		utils.Fail(c, 1021, "验证码错误")
		return
	}

	// 更新手机号
	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		utils.Fail(c, 1006, "用户不存在")
		return
	}

	user.Phone = req.Phone
	if err := h.userRepo.Update(user); err != nil {
		utils.ServerError(c, "绑定失败")
		return
	}

	// 删除验证码
	database.RDB.Del(ctx, codeKey)

	utils.SuccessWithMessage(c, "手机绑定成功", nil)
}

// getInactiveLockDays 从 config_webs 表读取用户未登录锁定天数，默认30，0=不锁定
func (h *AuthHandler) getInactiveLockDays() int {
	var web model.ConfigWeb
	if err := database.DB.Where("code = ? AND tenant_id = 0", "user_lock_inactive_days").First(&web).Error; err == nil {
		var v int
		if n, _ := fmt.Sscanf(web.Value, "%d", &v); n == 1 && v >= 0 {
			return v
		}
	}
	return 30
}

func (h *AuthHandler) recordLoginLog(tenantID, userID uint, username, ip, userAgent string, status int8, message string) {
	if !logcfg.IsLogEnabled("log_login_enabled") {
		return
	}
	log := model.LoginLog{
		TenantID:  tenantID,
		UserID:    userID,
		Username:  username,
		IP:        ip,
		UserAgent: userAgent,
		Status:    status,
		Message:   message,
		CreatedAt: time.Now(),
	}
	go database.DB.Create(&log)
}
