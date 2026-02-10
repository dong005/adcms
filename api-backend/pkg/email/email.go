package email

import (
	"adcms/internal/model"
	"adcms/pkg/database"
	"adcms/pkg/logcfg"
	"fmt"
	"strconv"

	"gopkg.in/gomail.v2"
)

// SMTPConfig SMTP配置
type SMTPConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
}

// GetSMTPConfig 从 system_configs 表读取 SMTP 配置
func GetSMTPConfig() (*SMTPConfig, error) {
	var configs []model.SystemConfig
	err := database.DB.Where("`key` IN ?", []string{
		"smtp_host", "smtp_port", "smtp_user", "smtp_password", "smtp_from",
	}).Find(&configs).Error
	if err != nil {
		return nil, fmt.Errorf("读取邮箱配置失败: %w", err)
	}

	configMap := make(map[string]string)
	for _, c := range configs {
		configMap[c.Key] = c.Value
	}

	host := configMap["smtp_host"]
	if host == "" {
		return nil, fmt.Errorf("未配置SMTP服务器地址")
	}

	port, _ := strconv.Atoi(configMap["smtp_port"])
	if port == 0 {
		port = 587
	}

	user := configMap["smtp_user"]
	password := configMap["smtp_password"]
	from := configMap["smtp_from"]
	if from == "" {
		from = user
	}

	return &SMTPConfig{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		From:     from,
	}, nil
}

// SendMail 发送邮件
func SendMail(to, subject, body string) error {
	cfg, err := GetSMTPConfig()
	if err != nil {
		return err
	}
	return SendMailWithConfig(cfg, to, subject, body)
}

// SendMailWithConfig 使用指定配置发送邮件
func SendMailWithConfig(cfg *SMTPConfig, to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", cfg.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(cfg.Host, cfg.Port, cfg.User, cfg.Password)

	if err := d.DialAndSend(m); err != nil {
		logEmail(to, subject, 0, err.Error())
		return fmt.Errorf("发送邮件失败: %w", err)
	}
	logEmail(to, subject, 1, "")
	return nil
}

func logEmail(to, subject string, status int8, errMsg string) {
	go func() {
		if !logcfg.IsLogEnabled("log_email_enabled") {
			return
		}
		database.DB.Create(&model.EmailLog{
			To:      to,
			Subject: subject,
			Status:  status,
			Error:   errMsg,
		})
	}()
}

// SendResetCode 发送密码重置验证码邮件
func SendResetCode(to, code string) error {
	subject := "ADCMS 密码重置验证码"
	body := fmt.Sprintf(`
		<div style="max-width:500px;margin:0 auto;padding:20px;font-family:Arial,sans-serif;">
			<h2 style="color:#1890ff;">ADCMS 密码重置</h2>
			<p>您正在进行密码重置操作，验证码为：</p>
			<div style="background:#f5f5f5;padding:15px;text-align:center;font-size:28px;font-weight:bold;letter-spacing:8px;color:#333;border-radius:4px;">
				%s
			</div>
			<p style="color:#999;font-size:12px;margin-top:15px;">
				验证码有效期为5分钟，请勿将验证码泄露给他人。<br/>
				如非本人操作，请忽略此邮件。
			</p>
		</div>
	`, code)
	return SendMail(to, subject, body)
}
