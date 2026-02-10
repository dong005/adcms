package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"
	"time"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{db: database.DB}
}

func (r *NotificationRepository) Create(n *model.Notification) error {
	return r.db.Create(n).Error
}

// SendToUser 发送通知给指定用户
func (r *NotificationRepository) SendToUser(tenantID, senderID, receiverID uint, title, content, nType string) error {
	n := model.Notification{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		SenderID:        senderID,
		ReceiverID:      receiverID,
		Title:           title,
		Content:         content,
		Type:            nType,
	}
	return r.db.Create(&n).Error
}

// SendToUsers 批量发送通知
func (r *NotificationRepository) SendToUsers(tenantID, senderID uint, receiverIDs []uint, title, content, nType string) error {
	for _, rid := range receiverIDs {
		n := model.Notification{
			TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
			SenderID:        senderID,
			ReceiverID:      rid,
			Title:           title,
			Content:         content,
			Type:            nType,
		}
		if err := r.db.Create(&n).Error; err != nil {
			return err
		}
	}
	return nil
}

// List 获取用户的通知列表（分页），返回 sender_name
func (r *NotificationRepository) List(tenantID, userID uint, page, pageSize int, nType string, isRead int8) ([]model.Notification, int64, error) {
	var notifications []model.Notification
	var total int64

	query := r.db.Where("tenant_id = ? AND receiver_id = ? AND reply_to_id = 0", tenantID, userID)
	if nType != "" {
		query = query.Where("type = ?", nType)
	}
	if isRead >= 0 {
		query = query.Where("is_read = ?", isRead)
	}

	query.Model(&model.Notification{}).Count(&total)

	err := query.Order("created_at DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&notifications).Error

	// 填充 sender_name
	r.fillSenderNames(notifications)

	return notifications, total, err
}

// fillSenderNames 批量填充发送者名称
func (r *NotificationRepository) fillSenderNames(notifications []model.Notification) {
	idSet := make(map[uint]bool)
	for _, n := range notifications {
		if n.SenderID > 0 {
			idSet[n.SenderID] = true
		}
	}
	if len(idSet) == 0 {
		for i := range notifications {
			if notifications[i].SenderID == 0 {
				notifications[i].SenderName = "系统"
			}
		}
		return
	}
	ids := make([]uint, 0, len(idSet))
	for id := range idSet {
		ids = append(ids, id)
	}
	var users []model.User
	r.db.Select("id, username, nickname").Where("id IN ?", ids).Find(&users)
	nameMap := make(map[uint]string)
	for _, u := range users {
		name := u.Nickname
		if name == "" {
			name = u.Username
		}
		nameMap[u.ID] = name
	}
	for i := range notifications {
		if notifications[i].SenderID == 0 {
			notifications[i].SenderName = "系统"
		} else {
			notifications[i].SenderName = nameMap[notifications[i].SenderID]
		}
	}
}

// GetByID 获取消息详情
func (r *NotificationRepository) GetByID(id uint) (*model.Notification, error) {
	var n model.Notification
	err := r.db.First(&n, id).Error
	if err != nil {
		return nil, err
	}
	// 填充 sender_name
	list := []model.Notification{n}
	r.fillSenderNames(list)
	n.SenderName = list[0].SenderName
	return &n, nil
}

// GetReplies 获取某条消息的回复列表
func (r *NotificationRepository) GetReplies(replyToID uint) ([]model.Notification, error) {
	var replies []model.Notification
	err := r.db.Where("reply_to_id = ?", replyToID).Order("created_at ASC").Find(&replies).Error
	if err != nil {
		return nil, err
	}
	r.fillSenderNames(replies)
	return replies, nil
}

// FindUsersByRoleIDs 根据角色ID列表查找用户ID列表
func (r *NotificationRepository) FindUsersByRoleIDs(roleIDs []uint) ([]uint, error) {
	var userIDs []uint
	err := r.db.Model(&model.UserRole{}).Where("role_id IN ?", roleIDs).Distinct("user_id").Pluck("user_id", &userIDs).Error
	return userIDs, err
}

// FindUsersWithEmailNotify 查找开启了邮件通知的用户
func (r *NotificationRepository) FindUsersWithEmailNotify(userIDs []uint) ([]model.User, error) {
	var users []model.User
	err := r.db.Select("id, email, email_notify").Where("id IN ? AND email_notify = 1 AND email != ''", userIDs).Find(&users).Error
	return users, err
}

// MarkAsRead 标记为已读
func (r *NotificationRepository) MarkAsRead(id, userID uint) error {
	now := time.Now()
	return r.db.Model(&model.Notification{}).
		Where("id = ? AND receiver_id = ?", id, userID).
		Updates(map[string]interface{}{"is_read": 1, "read_at": &now}).Error
}

// MarkAllAsRead 全部标记为已读
func (r *NotificationRepository) MarkAllAsRead(tenantID, userID uint) error {
	now := time.Now()
	return r.db.Model(&model.Notification{}).
		Where("tenant_id = ? AND receiver_id = ? AND is_read = 0", tenantID, userID).
		Updates(map[string]interface{}{"is_read": 1, "read_at": &now}).Error
}

// UnreadCount 获取未读数量
func (r *NotificationRepository) UnreadCount(tenantID, userID uint) int64 {
	var count int64
	r.db.Model(&model.Notification{}).
		Where("tenant_id = ? AND receiver_id = ? AND is_read = 0", tenantID, userID).
		Count(&count)
	return count
}

// Delete 删除通知
func (r *NotificationRepository) Delete(id, userID uint) error {
	return r.db.Where("id = ? AND receiver_id = ?", id, userID).Delete(&model.Notification{}).Error
}

// DeleteByID 按ID删除通知（不限制receiver_id，超管用）
func (r *NotificationRepository) DeleteByID(id uint) error {
	return r.db.Delete(&model.Notification{}, id).Error
}

// Reply 回复消息，返回 (receiverID, originalTitle, error)
func (r *NotificationRepository) Reply(tenantID, senderID, replyToID uint, content string) (uint, string, error) {
	// 获取原消息
	var original model.Notification
	if err := r.db.First(&original, replyToID).Error; err != nil {
		return 0, "", err
	}

	// 回复发给原消息的发送者
	receiverID := original.SenderID

	n := model.Notification{
		TenantBaseModel: model.TenantBaseModel{TenantID: tenantID},
		SenderID:        senderID,
		ReceiverID:      receiverID,
		Title:           "回复: " + original.Title,
		Content:         content,
		Type:            "message",
		ReplyToID:       replyToID,
	}
	return receiverID, original.Title, r.db.Create(&n).Error
}
