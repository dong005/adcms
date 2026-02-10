package handler

import (
	"adcms/internal/middleware"
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/database"
	"adcms/pkg/email"
	"adcms/pkg/utils"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type NotificationHandler struct {
	notifRepo *repository.NotificationRepository
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{notifRepo: repository.NewNotificationRepository()}
}

func (h *NotificationHandler) List(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	nType := c.Query("type")
	isRead := int8(-1)
	if r := c.Query("is_read"); r != "" {
		v, _ := strconv.Atoi(r)
		isRead = int8(v)
	}

	notifications, total, err := h.notifRepo.List(tenantID, userID, page, pageSize, nType, isRead)
	if err != nil {
		utils.ServerError(c, "查询失败")
		return
	}
	utils.SuccessWithPage(c, notifications, total, page, pageSize)
}

func (h *NotificationHandler) UnreadCount(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)
	count := h.notifRepo.UnreadCount(tenantID, userID)
	utils.Success(c, map[string]int64{"count": count})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetUserID(c)

	if err := h.notifRepo.MarkAsRead(uint(id), userID); err != nil {
		utils.ServerError(c, "操作失败")
		return
	}
	utils.SuccessWithMessage(c, "已读", nil)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	tenantID := middleware.GetTenantID(c)
	userID := middleware.GetUserID(c)

	if err := h.notifRepo.MarkAllAsRead(tenantID, userID); err != nil {
		utils.ServerError(c, "操作失败")
		return
	}
	utils.SuccessWithMessage(c, "全部已读", nil)
}

func (h *NotificationHandler) Delete(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)
	userID := middleware.GetUserID(c)

	if err := h.notifRepo.Delete(uint(id), userID); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}

type SendNotificationRequest struct {
	ReceiverIDs []uint `json:"receiver_ids"`
	RoleIDs     []uint `json:"role_ids"`
	Title       string `json:"title" binding:"required"`
	Content     string `json:"content"`
	Type        string `json:"type"`
}

func (h *NotificationHandler) Send(c *gin.Context) {
	// 仅超管可发送消息
	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可发送消息")
		return
	}

	var req SendNotificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "参数错误")
		return
	}

	tenantID := middleware.GetTenantID(c)
	senderID := middleware.GetUserID(c)
	nType := req.Type
	if nType == "" {
		nType = "system"
	}

	// 合并收件人：直接指定的用户 + 按角色查找的用户
	receiverSet := make(map[uint]bool)
	for _, id := range req.ReceiverIDs {
		receiverSet[id] = true
	}
	if len(req.RoleIDs) > 0 {
		roleUserIDs, err := h.notifRepo.FindUsersByRoleIDs(req.RoleIDs)
		if err == nil {
			for _, id := range roleUserIDs {
				receiverSet[id] = true
			}
		}
	}

	if len(receiverSet) == 0 {
		utils.BadRequest(c, "请选择收件人")
		return
	}

	receiverIDs := make([]uint, 0, len(receiverSet))
	for id := range receiverSet {
		receiverIDs = append(receiverIDs, id)
	}

	if err := h.notifRepo.SendToUsers(tenantID, senderID, receiverIDs, req.Title, req.Content, nType); err != nil {
		utils.ServerError(c, "发送失败")
		return
	}

	// 异步抄送邮件给开启了邮件通知的用户
	go func() {
		users, err := h.notifRepo.FindUsersWithEmailNotify(receiverIDs)
		if err != nil {
			return
		}
		for _, u := range users {
			subject := "[ADCMS] " + req.Title
			body := fmt.Sprintf(`<div style="max-width:600px;margin:0 auto;padding:20px;font-family:Arial,sans-serif;">
				<h2 style="color:#1890ff;">%s</h2>
				<div style="padding:15px;background:#f5f5f5;border-radius:4px;">%s</div>
				<p style="color:#999;font-size:12px;margin-top:15px;">此邮件由 ADCMS 系统自动发送，您可以在个人设置中关闭邮件通知。</p>
			</div>`, req.Title, req.Content)
			if err := email.SendMail(u.Email, subject, body); err != nil {
				fmt.Printf("[Email] 发送通知邮件失败 to=%s err=%v\n", u.Email, err)
			}
		}
	}()

	utils.SuccessWithMessage(c, fmt.Sprintf("已发送给 %d 位用户", len(receiverIDs)), nil)
}

func (h *NotificationHandler) Detail(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	notif, err := h.notifRepo.GetByID(uint(id))
	if err != nil {
		utils.Fail(c, 4001, "消息不存在")
		return
	}

	replies, _ := h.notifRepo.GetReplies(uint(id))
	if replies == nil {
		replies = []model.Notification{}
	}

	utils.Success(c, map[string]interface{}{
		"notification": notif,
		"replies":      replies,
	})
}

type ReplyRequest struct {
	Content string `json:"content" binding:"required"`
}

func (h *NotificationHandler) Reply(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	var req ReplyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "请输入回复内容")
		return
	}

	tenantID := middleware.GetTenantID(c)
	senderID := middleware.GetUserID(c)

	receiverID, originalTitle, err := h.notifRepo.Reply(tenantID, senderID, uint(id), req.Content)
	if err != nil {
		utils.ServerError(c, "回复失败")
		return
	}

	// 异步发送邮件通知给原消息发送者
	go func() {
		if receiverID == 0 {
			return
		}
		users, err := h.notifRepo.FindUsersWithEmailNotify([]uint{receiverID})
		if err != nil || len(users) == 0 {
			return
		}
		// 获取回复者名称
		var sender model.User
		database.DB.Select("nickname, username").First(&sender, senderID)
		senderName := sender.Nickname
		if senderName == "" {
			senderName = sender.Username
		}

		for _, u := range users {
			subject := fmt.Sprintf("[ADCMS] %s 回复了「%s」", senderName, originalTitle)
			body := fmt.Sprintf(`<div style="max-width:600px;margin:0 auto;padding:20px;font-family:Arial,sans-serif;">
				<h2 style="color:#1890ff;">%s 回复了您的消息</h2>
				<p style="color:#666;">原消息：%s</p>
				<div style="padding:15px;background:#f5f5f5;border-radius:4px;border-left:3px solid #1890ff;">%s</div>
				<p style="color:#999;font-size:12px;margin-top:15px;">此邮件由 ADCMS 系统自动发送，您可以在个人设置中关闭邮件通知。</p>
			</div>`, senderName, originalTitle, req.Content)
			if err := email.SendMail(u.Email, subject, body); err != nil {
				fmt.Printf("[Email] 发送回复通知邮件失败 to=%s err=%v\n", u.Email, err)
			}
		}
	}()

	utils.SuccessWithMessage(c, "回复成功", nil)
}

// DeleteReply 删除回复（仅超管）
func (h *NotificationHandler) DeleteReply(c *gin.Context) {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 64)

	if !middleware.IsSuperAdmin(middleware.GetUserID(c)) {
		utils.Fail(c, 4003, "仅超级管理员可删除回复")
		return
	}

	if err := h.notifRepo.DeleteByID(uint(id)); err != nil {
		utils.ServerError(c, "删除失败")
		return
	}
	utils.SuccessWithMessage(c, "删除成功", nil)
}
