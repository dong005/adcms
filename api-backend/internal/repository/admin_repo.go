package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"
	"time"
	"gorm.io/gorm"
)

type AdminRepository struct {
	db *gorm.DB
}

func NewAdminRepository() *AdminRepository {
	return &AdminRepository{db: database.DB}
}

type AdminDetail struct {
	model.User
	UserCount     uint      `json:"user_count"`     // 租户用户数
	ArticleCount  uint      `json:"article_count"`  // 文章数
	CategoryCount uint      `json:"category_count"` // 分类数
	MediaCount    uint      `json:"media_count"`    // 媒体数
	CreatedAt     string    `json:"created_at"`
	UpdatedAt     string    `json:"updated_at"`
}

type AdminStatistics struct {
	UserCount     uint      `json:"user_count"`     // 租户用户数
	ArticleCount  uint      `json:"article_count"`  // 文章数
	CategoryCount uint      `json:"category_count"` // 分类数
	MediaCount    uint      `json:"media_count"`    // 媒体数
	LoginCount    uint      `json:"login_count"`    // 登录次数
	LastLoginAt   *string   `json:"last_login_at"`  // 最后登录时间
}

func (r *AdminRepository) List(page, pageSize int, keyword string) ([]AdminDetail, int64, error) {
	var admins []AdminDetail
	var total int64

	query := r.db.Model(&model.User{}).Where("is_admin = 1")
	
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR company LIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	
	// 查询管理员列表
	rows, err := r.db.Raw(`
		SELECT 
			u.*,
			(SELECT COUNT(*) FROM users WHERE tenant_id = u.id AND deleted_at IS NULL) as user_count,
			(SELECT COUNT(*) FROM articles WHERE tenant_id = u.id AND deleted_at IS NULL) as article_count,
			(SELECT COUNT(*) FROM categories WHERE tenant_id = u.id AND deleted_at IS NULL) as category_count,
			(SELECT COUNT(*) FROM media WHERE tenant_id = u.id AND deleted_at IS NULL) as media_count
		FROM users u
		WHERE u.is_admin = 1 AND u.deleted_at IS NULL
		ORDER BY u.created_at DESC
		LIMIT ? OFFSET ?
	`, pageSize, offset).Rows()
	
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var admin AdminDetail
		err := rows.Scan(
			&admin.ID, &admin.CreatedAt, &admin.UpdatedAt, &admin.DeletedAt,
			&admin.TenantID, &admin.Username, &admin.Password, &admin.Email,
			&admin.Phone, &admin.Nickname, &admin.Avatar, &admin.DepartmentID,
			&admin.Status, &admin.TOTPEnabled, &admin.TOTPSecret, &admin.EmailNotify,
			&admin.LastLoginAt, &admin.LastLoginIP, &admin.IsAdmin, &admin.Company,
			&admin.Domain, &admin.ExpireTime, &admin.MaxUsers, &admin.LoginCount,
			&admin.Remark, &admin.UserCount, &admin.ArticleCount,
			&admin.CategoryCount, &admin.MediaCount,
		)
		if err != nil {
			continue
		}
		admins = append(admins, admin)
	}

	return admins, total, nil
}

func (r *AdminRepository) Detail(id uint) (*AdminDetail, error) {
	var admin AdminDetail
	
	err := r.db.Raw(`
		SELECT 
			u.*,
			(SELECT COUNT(*) FROM users WHERE tenant_id = u.id AND deleted_at IS NULL) as user_count,
			(SELECT COUNT(*) FROM articles WHERE tenant_id = u.id AND deleted_at IS NULL) as article_count,
			(SELECT COUNT(*) FROM categories WHERE tenant_id = u.id AND deleted_at IS NULL) as category_count,
			(SELECT COUNT(*) FROM media WHERE tenant_id = u.id AND deleted_at IS NULL) as media_count
		FROM users u
		WHERE u.id = ? AND u.is_admin = 1 AND u.deleted_at IS NULL
	`, id).Scan(&admin).Error
	
	if err != nil {
		return nil, err
	}
	
	return &admin, nil
}

func (r *AdminRepository) Statistics(id uint) (*AdminStatistics, error) {
	var stats AdminStatistics
	
	err := r.db.Raw(`
		SELECT 
			(SELECT COUNT(*) FROM users WHERE tenant_id = ? AND deleted_at IS NULL) as user_count,
			(SELECT COUNT(*) FROM articles WHERE tenant_id = ? AND deleted_at IS NULL) as article_count,
			(SELECT COUNT(*) FROM categories WHERE tenant_id = ? AND deleted_at IS NULL) as category_count,
			(SELECT COUNT(*) FROM media WHERE tenant_id = ? AND deleted_at IS NULL) as media_count,
			login_count,
			last_login_at
		FROM users 
		WHERE id = ? AND is_admin = 1
	`, id, id, id, id, id).Scan(&stats).Error
	
	if err != nil {
		return nil, err
	}
	
	return &stats, nil
}

// GetTenantUsers 获取租户下的所有用户
func (r *AdminRepository) GetTenantUsers(tenantID uint, page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{}).Where("tenant_id = ? AND is_admin != 2", tenantID)
	
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?", 
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 查询总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error
	
	return users, total, err
}

// DisableTenant 禁用租户（禁用租户下的所有用户）
func (r *AdminRepository) DisableTenant(tenantID uint) error {
	return r.db.Model(&model.User{}).Where("tenant_id = ?", tenantID).Update("status", 0).Error
}

// EnableTenant 启用租户
func (r *AdminRepository) EnableTenant(tenantID uint) error {
	return r.db.Model(&model.User{}).Where("tenant_id = ? AND is_admin != 2", tenantID).Update("status", 1).Error
}

// CheckTenantExpired 检查租户是否过期
func (r *AdminRepository) CheckTenantExpired() ([]model.User, error) {
	var expiredTenants []model.User
	now := time.Now()
	
	err := r.db.Where("is_admin = 1 AND expire_time IS NOT NULL AND expire_time < ?", now).
		Find(&expiredTenants).Error
	
	return expiredTenants, err
}
