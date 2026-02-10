import { requestClient } from '#/api/request';

export function getSmsConfig() {
  return requestClient.get<Record<string, string>>('/configs/sms');
}

export function updateSmsConfig(data: {
  sms_secret_id: string;
  sms_secret_key: string;
  sms_app_id: string;
  sms_sign: string;
  sms_template_id: string;
}) {
  return requestClient.put('/configs/sms', data);
}

export function testSms(data: { phone: string }) {
  return requestClient.post('/configs/sms/test', data);
}
