import { requestClient } from '#/api/request';

/** 获取邮箱配置 */
export async function getEmailConfig() {
  return requestClient.get<Record<string, string>>('/configs/email');
}

/** 保存邮箱配置 */
export async function updateEmailConfig(data: {
  smtp_host: string;
  smtp_port: string;
  smtp_user: string;
  smtp_password: string;
  smtp_from: string;
}) {
  return requestClient.put('/configs/email', data);
}

/** 发送测试邮件 */
export async function testEmail(data: { to: string }) {
  return requestClient.post('/configs/email/test', data);
}
