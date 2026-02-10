package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/database"
	"adcms/pkg/email"
	"adcms/pkg/sms"
	"adcms/pkg/logcfg"
	"adcms/pkg/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ConfigHandler struct {
	configRepo *repository.ConfigRepository
}

func NewConfigHandler() *ConfigHandler {
	return &ConfigHandler{configRepo: repository.NewConfigRepository()}
}

func (h *ConfigHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	configs, err := h.configRepo.FindAll(tenantID)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.Success(c, configs)
}

type UpdateConfigRequest struct {
	Configs []struct {
		Key         string `json:"key"`
		Value       string `json:"value"`
		Description string `json:"description"`
	} `json:"configs"`
}

func (h *ConfigHandler) Update(c *gin.Context) {
	var req UpdateConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tenantID := middleware.GetTenantID(c)
	for _, cfg := range req.Configs {
		config := model.SystemConfig{
			TenantID:    tenantID,
			Key:         cfg.Key,
			Value:       cfg.Value,
			Description: cfg.Description,
		}
		h.configRepo.Upsert(&config)
	}

	utils.SuccessWithMessage(c, "更新成功", nil)
}

type LogHandler struct {
	logRepo *repository.LogRepository
}

func NewLogHandler() *LogHandler {
	return &LogHandler{logRepo: repository.NewLogRepository()}
}

func (h *LogHandler) OperationLogs(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	module := c.Query("module")
	username := c.Query("username")

	logs, total, err := h.logRepo.ListOperationLogs(tenantID, page, pageSize, module, username)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.SuccessWithPage(c, logs, total, page, pageSize)
}

func (h *LogHandler) LoginLogs(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	username := c.Query("username")

	logs, total, err := h.logRepo.ListLoginLogs(tenantID, page, pageSize, username)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.SuccessWithPage(c, logs, total, page, pageSize)
}

func (h *LogHandler) EmailLogs(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可查看")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	logs, total, err := h.logRepo.ListEmailLogs(page, pageSize, keyword)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.SuccessWithPage(c, logs, total, page, pageSize)
}

func (h *LogHandler) SmsLogs(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可查看")
		return
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	logs, total, err := h.logRepo.ListSmsLogs(page, pageSize, keyword)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.SuccessWithPage(c, logs, total, page, pageSize)
}

type PermissionHandler struct {
	permRepo *repository.PermissionRepository
}

func NewPermissionHandler() *PermissionHandler {
	return &PermissionHandler{permRepo: repository.NewPermissionRepository()}
}

func (h *PermissionHandler) List(c *gin.Context) {
	permissions, err := h.permRepo.FindAll()
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.Success(c, permissions)
}

func (h *PermissionHandler) Tree(c *gin.Context) {
	permissions, err := h.permRepo.FindAll()
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	tree := repository.BuildPermissionTree(permissions, 0)
	utils.Success(c, tree)
}

func (h *PermissionHandler) Create(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可操作")
		return
	}
	var req struct {
		Name        string `json:"name" binding:"required"`
		Code        string `json:"code" binding:"required"`
		Type        int8   `json:"type"`
		ParentID    uint   `json:"parent_id"`
		Path        string `json:"path"`
		Method      string `json:"method"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}
	perm := &model.Permission{
		Name:        req.Name,
		Code:        req.Code,
		Type:        req.Type,
		ParentID:    req.ParentID,
		Path:        req.Path,
		Method:      req.Method,
		Description: req.Description,
	}
	if err := h.permRepo.Create(perm); err != nil {
		utils.ServerError(c, "创建失败")
		return
	}
	utils.Success(c, perm)
}

func (h *PermissionHandler) Update(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可操作")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}
	perm, err := h.permRepo.FindByID(uint(id))
	if err != nil {
		utils.Fail(c, 4001, "权限不存在")
		return
	}
	var req struct {
		Name        string `json:"name"`
		Code        string `json:"code"`
		Type        int8   `json:"type"`
		ParentID    uint   `json:"parent_id"`
		Path        string `json:"path"`
		Method      string `json:"method"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}
	if req.Name != "" {
		perm.Name = req.Name
	}
	if req.Code != "" {
		perm.Code = req.Code
	}
	perm.Type = req.Type
	perm.ParentID = req.ParentID
	perm.Path = req.Path
	perm.Method = req.Method
	perm.Description = req.Description

	if err := h.permRepo.Update(perm); err != nil {
		utils.ServerError(c, "更新失败")
		return
	}
	utils.Success(c, perm)
}

func (h *PermissionHandler) Delete(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可操作")
		return
	}
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}
	if err := h.permRepo.Delete(uint(id)); err != nil {
		utils.Fail(c, 4002, err.Error())
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

// 邮箱配置相关接口

var emailConfigKeys = []string{"smtp_host", "smtp_port", "smtp_user", "smtp_password", "smtp_from"}

func (h *ConfigHandler) GetEmailConfig(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可操作")
		return
	}

	var configs []model.SystemConfig
	database.DB.Where("`key` IN ?", emailConfigKeys).Find(&configs)

	result := make(map[string]string)
	for _, cfg := range configs {
		result[cfg.Key] = cfg.Value
	}
	utils.Success(c, result)
}

type EmailConfigRequest struct {
	SmtpHost     string `json:"smtp_host"`
	SmtpPort     string `json:"smtp_port"`
	SmtpUser     string `json:"smtp_user"`
	SmtpPassword string `json:"smtp_password"`
	SmtpFrom     string `json:"smtp_from"`
}

func (h *ConfigHandler) UpdateEmailConfig(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可操作")
		return
	}

	var req EmailConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	items := map[string]string{
		"smtp_host":     req.SmtpHost,
		"smtp_port":     req.SmtpPort,
		"smtp_user":     req.SmtpUser,
		"smtp_password": req.SmtpPassword,
		"smtp_from":     req.SmtpFrom,
	}

	for key, value := range items {
		cfg := model.SystemConfig{
			TenantID:    0,
			Key:         key,
			Value:       value,
			Description: "邮箱配置",
		}
		h.configRepo.Upsert(&cfg)
	}

	utils.SuccessWithMessage(c, "邮箱配置已保存", nil)
}

type TestEmailRequest struct {
	To string `json:"to" binding:"required,email"`
}

func (h *ConfigHandler) TestEmail(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可操作")
		return
	}

	var req TestEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请输入有效的收件邮箱")
		return
	}

	cfg, err := email.GetSMTPConfig()
	if err != nil {
		utils.Fail(c, 1030, "邮箱配置不完整: "+err.Error())
		return
	}

	err = email.SendMailWithConfig(cfg, req.To, "ADCMS 邮件测试", "<h3>恭喜！ADCMS 邮件配置测试成功。</h3><p>此邮件由系统自动发送，用于验证SMTP配置是否正确。</p>")
	if err != nil {
		utils.Fail(c, 1031, "发送失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "测试邮件已发送", nil)
}

// 短信配置相关接口

var smsConfigKeys = []string{"sms_secret_id", "sms_secret_key", "sms_app_id", "sms_sign", "sms_template_id"}

func (h *ConfigHandler) GetSmsConfig(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可操作")
		return
	}

	var configs []model.SystemConfig
	database.DB.Where("`key` IN ?", smsConfigKeys).Find(&configs)

	result := make(map[string]string)
	for _, cfg := range configs {
		result[cfg.Key] = cfg.Value
	}
	utils.Success(c, result)
}

type SmsConfigRequest struct {
	SecretId   string `json:"sms_secret_id"`
	SecretKey  string `json:"sms_secret_key"`
	AppId      string `json:"sms_app_id"`
	Sign       string `json:"sms_sign"`
	TemplateId string `json:"sms_template_id"`
}

func (h *ConfigHandler) UpdateSmsConfig(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可操作")
		return
	}

	var req SmsConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	items := map[string]string{
		"sms_secret_id":   req.SecretId,
		"sms_secret_key":  req.SecretKey,
		"sms_app_id":      req.AppId,
		"sms_sign":        req.Sign,
		"sms_template_id": req.TemplateId,
	}

	for key, value := range items {
		cfg := model.SystemConfig{
			TenantID:    0,
			Key:         key,
			Value:       value,
			Description: "短信配置",
		}
		h.configRepo.Upsert(&cfg)
	}

	utils.SuccessWithMessage(c, "短信配置已保存", nil)
}

type TestSmsRequest struct {
	Phone string `json:"phone" binding:"required"`
}

func (h *ConfigHandler) TestSms(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可操作")
		return
	}

	var req TestSmsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请输入手机号")
		return
	}

	cfg, err := sms.GetSMSConfig()
	if err != nil {
		utils.Fail(c, 1030, "短信配置不完整: "+err.Error())
		return
	}

	err = sms.SendSMSWithConfig(cfg, req.Phone, []string{"123456", "5"})
	if err != nil {
		utils.Fail(c, 1031, "发送失败: "+err.Error())
		return
	}

	utils.SuccessWithMessage(c, "测试短信已发送", nil)
}

// 日志配置相关接口
var logConfigKeys = []string{"log_operation_enabled", "log_login_enabled", "log_email_enabled", "log_sms_enabled"}

func (h *ConfigHandler) GetLogConfig(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可操作")
		return
	}

	var configs []model.SystemConfig
	database.DB.Where("`key` IN ?", logConfigKeys).Find(&configs)

	result := make(map[string]string)
	for _, key := range logConfigKeys {
		result[key] = "1" // 默认启用
	}
	for _, cfg := range configs {
		result[cfg.Key] = cfg.Value
	}
	utils.Success(c, result)
}

type LogConfigRequest struct {
	OperationEnabled string `json:"log_operation_enabled"`
	LoginEnabled     string `json:"log_login_enabled"`
	EmailEnabled     string `json:"log_email_enabled"`
	SmsEnabled       string `json:"log_sms_enabled"`
}

func (h *ConfigHandler) UpdateLogConfig(c *gin.Context) {
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可操作")
		return
	}

	var req LogConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	items := map[string]string{
		"log_operation_enabled": req.OperationEnabled,
		"log_login_enabled":    req.LoginEnabled,
		"log_email_enabled":    req.EmailEnabled,
		"log_sms_enabled":      req.SmsEnabled,
	}

	for key, value := range items {
		cfg := model.SystemConfig{
			TenantID:    0,
			Key:         key,
			Value:       value,
			Description: "日志配置",
		}
		h.configRepo.Upsert(&cfg)
	}

	// 清除缓存
	logcfg.ClearCache()

	utils.SuccessWithMessage(c, "日志配置已保存", nil)
}
