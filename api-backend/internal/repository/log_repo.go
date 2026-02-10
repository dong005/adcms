package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type LogRepository struct {
	db *gorm.DB
}

func NewLogRepository() *LogRepository {
	return &LogRepository{db: database.DB}
}

func (r *LogRepository) ListOperationLogs(tenantID uint, page, pageSize int, module, username string) ([]model.OperationLog, int64, error) {
	var logs []model.OperationLog
	var total int64

	query := r.db.Model(&model.OperationLog{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if module != "" {
		query = query.Where("module = ?", module)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error
	return logs, total, err
}

func (r *LogRepository) ListLoginLogs(tenantID uint, page, pageSize int, username string) ([]model.LoginLog, int64, error) {
	var logs []model.LoginLog
	var total int64

	query := r.db.Model(&model.LoginLog{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if username != "" {
		query = query.Where("username LIKE ?", "%"+username+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error
	return logs, total, err
}

func (r *LogRepository) ListEmailLogs(page, pageSize int, keyword string) ([]model.EmailLog, int64, error) {
	var logs []model.EmailLog
	var total int64

	query := r.db.Model(&model.EmailLog{})
	if keyword != "" {
		query = query.Where("`to` LIKE ? OR subject LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error
	return logs, total, err
}

func (r *LogRepository) ListSmsLogs(page, pageSize int, keyword string) ([]model.SmsLog, int64, error) {
	var logs []model.SmsLog
	var total int64

	query := r.db.Model(&model.SmsLog{})
	if keyword != "" {
		query = query.Where("phone LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&logs).Error
	return logs, total, err
}
