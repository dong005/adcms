package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository() *UserRepository {
	return &UserRepository{db: database.DB}
}

func (r *UserRepository) Create(user *model.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) Update(user *model.User) error {
	return r.db.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.db.Delete(&model.User{}, id).Error
}

func (r *UserRepository) FindByID(id uint) (*model.User, error) {
	var user model.User
	err := r.db.Preload("Roles").First(&user, id).Error
	return &user, err
}

func (r *UserRepository) FindByUsername(tenantID uint, username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("tenant_id = ? AND username = ?", tenantID, username).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByUsernameGlobal(username string) (*model.User, error) {
	var user model.User
	err := r.db.Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByEmail(tenantID uint, email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("tenant_id = ? AND email = ?", tenantID, email).First(&user).Error
	return &user, err
}

func (r *UserRepository) FindByEmailGlobal(email string) (*model.User, error) {
	var user model.User
	err := r.db.Where("email = ? AND email != ''", email).First(&user).Error
	return &user, err
}

func (r *UserRepository) List(tenantID uint, page, pageSize int, keyword string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := r.db.Model(&model.User{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if keyword != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ? OR email LIKE ?",
			"%"+keyword+"%", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("Roles").Offset(offset).Limit(pageSize).Order("id DESC").Find(&users).Error
	return users, total, err
}

func (r *UserRepository) UpdatePassword(id uint, password string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("password", password).Error
}

func (r *UserRepository) UpdateTOTP(id uint, enabled int8, secret string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"totp_enabled": enabled,
		"totp_secret":  secret,
	}).Error
}

func (r *UserRepository) UpdateStatus(id uint, status int8) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Update("status", status).Error
}

func (r *UserRepository) UpdateLoginInfo(id uint, ip string) error {
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(map[string]interface{}{
		"last_login_at": gorm.Expr("NOW()"),
		"last_login_ip": ip,
	}).Error
}

func (r *UserRepository) AssignRoles(userID uint, roleIDs []uint) error {
	if err := r.db.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
		return err
	}

	for _, roleID := range roleIDs {
		userRole := model.UserRole{UserID: userID, RoleID: roleID}
		if err := r.db.Create(&userRole).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *UserRepository) GetUserRoles(userID uint) ([]model.Role, error) {
	var user model.User
	err := r.db.Preload("Roles").First(&user, userID).Error
	return user.Roles, err
}

func (r *UserRepository) SetTenantID(userID, tenantID uint) error {
	return r.db.Model(&model.User{}).Where("id = ?", userID).Update("tenant_id", tenantID).Error
}

func (r *UserRepository) AssignMenus(userID uint, menuIDs []uint) error {
	// 先清除旧记录
	if err := r.db.Where("user_id = ?", userID).Delete(&model.UserMenu{}).Error; err != nil {
		return err
	}

	// 批量插入新记录
	for _, menuID := range menuIDs {
		userMenu := model.UserMenu{
			UserID: userID,
			MenuID: menuID,
		}
		if err := r.db.Create(&userMenu).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *UserRepository) GetUserMenus(userID uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Joins("JOIN user_menus ON user_menus.menu_id = menus.id").
		Where("user_menus.user_id = ?", userID).
		Find(&menus).Error
	return menus, err
}
