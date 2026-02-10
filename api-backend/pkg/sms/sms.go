package sms

import (
	"adcms/internal/model"
	"adcms/pkg/database"
	"adcms/pkg/logcfg"
	"fmt"
	"strings"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	smsClient "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

// SMSConfig 腾讯云短信配置
type SMSConfig struct {
	SecretId   string
	SecretKey  string
	AppId      string
	Sign       string
	TemplateId string
}

// GetSMSConfig 从 system_configs 表读取短信配置
func GetSMSConfig() (*SMSConfig, error) {
	var configs []model.SystemConfig
	err := database.DB.Where("`key` IN ?", []string{
		"sms_secret_id", "sms_secret_key", "sms_app_id", "sms_sign", "sms_template_id",
	}).Find(&configs).Error
	if err != nil {
		return nil, fmt.Errorf("读取短信配置失败: %w", err)
	}

	configMap := make(map[string]string)
	for _, c := range configs {
		configMap[c.Key] = c.Value
	}

	secretId := configMap["sms_secret_id"]
	secretKey := configMap["sms_secret_key"]
	if secretId == "" || secretKey == "" {
		return nil, fmt.Errorf("未配置腾讯云短信密钥")
	}

	appId := configMap["sms_app_id"]
	if appId == "" {
		return nil, fmt.Errorf("未配置短信 AppId")
	}

	return &SMSConfig{
		SecretId:   secretId,
		SecretKey:  secretKey,
		AppId:      appId,
		Sign:       configMap["sms_sign"],
		TemplateId: configMap["sms_template_id"],
	}, nil
}

// SendSMS 使用腾讯云发送短信
func SendSMS(phone string, templateParams []string) error {
	cfg, err := GetSMSConfig()
	if err != nil {
		return err
	}
	return SendSMSWithConfig(cfg, phone, templateParams)
}

// SendSMSWithConfig 使用指定配置发送短信
func SendSMSWithConfig(cfg *SMSConfig, phone string, templateParams []string) error {
	credential := common.NewCredential(cfg.SecretId, cfg.SecretKey)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com"

	client, err := smsClient.NewClient(credential, "ap-guangzhou", cpf)
	if err != nil {
		return fmt.Errorf("创建短信客户端失败: %w", err)
	}

	request := smsClient.NewSendSmsRequest()
	request.SmsSdkAppId = common.StringPtr(cfg.AppId)
	request.SignName = common.StringPtr(cfg.Sign)
	request.TemplateId = common.StringPtr(cfg.TemplateId)
	request.PhoneNumberSet = common.StringPtrs([]string{"+86" + phone})

	// 模板参数
	request.TemplateParamSet = common.StringPtrs(templateParams)

	response, err := client.SendSms(request)
	if err != nil {
		logSms(phone, cfg.TemplateId, templateParams, 0, err.Error())
		return fmt.Errorf("发送短信失败: %w", err)
	}

	// 检查发送状态
	if response.Response != nil && len(response.Response.SendStatusSet) > 0 {
		status := response.Response.SendStatusSet[0]
		if *status.Code != "Ok" {
			errMsg := fmt.Sprintf("%s - %s", *status.Code, *status.Message)
			logSms(phone, cfg.TemplateId, templateParams, 0, errMsg)
			return fmt.Errorf("短信发送失败: %s", errMsg)
		}
	}

	logSms(phone, cfg.TemplateId, templateParams, 1, "")
	return nil
}

func logSms(phone, templateID string, params []string, status int8, errMsg string) {
	go func() {
		if !logcfg.IsLogEnabled("log_sms_enabled") {
			return
		}
		database.DB.Create(&model.SmsLog{
			Phone:      phone,
			TemplateID: templateID,
			Params:     strings.Join(params, ","),
			Status:     status,
			Error:      errMsg,
		})
	}()
}

// SendVerifyCode 发送验证码短信
func SendVerifyCode(phone, code string) error {
	return SendSMS(phone, []string{code, "5"})
}
