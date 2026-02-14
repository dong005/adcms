import type { UserInfo } from '@vben/types';

import { requestClient } from '#/api/request';

/**
 * 获取用户信息
 */
export async function getUserInfoApi() {
  const res = await requestClient.get<{
    id: number;
    username: string;
    nickname: string;
    email: string;
    phone: string;
    avatar: string;
    roles: string[];
    totp_enabled: boolean;
    email_notify: number;
    is_admin: number;
  }>('/auth/user-info');
  return {
    userId: `${res.id}`,
    username: res.username,
    realName: res.nickname || res.username,
    avatar: res.avatar || '',
    roles: res.roles || [],
    desc: '',
    token: '',
    homePath: '/dashboard/analytics',
    email: res.email || '',
    phone: res.phone || '',
    totp_enabled: res.totp_enabled || false,
    email_notify: res.email_notify ?? 1,
    is_admin: res.is_admin ?? 0,
  } as UserInfo;
}

/**
 * 更新用户信息
 */
export async function updateUserInfoApi(data: {
  nickname?: string;
  email?: string;
  avatar?: string;
}) {
  return requestClient.put('/auth/user-info', data);
}

/**
 * 修改密码
 */
export async function changePasswordApi(data: {
  old_password: string;
  new_password: string;
}) {
  return requestClient.put('/auth/password', data);
}

/**
 * 生成 TOTP 密钥
 */
export async function generateTOTPApi() {
  return requestClient.post<{
    secret: string;
    qr_code: string;
  }>('/auth/totp/generate');
}

/**
 * 绑定 TOTP
 */
export async function bindTOTPApi(data: { code: string }) {
  return requestClient.post('/auth/totp/bind', data);
}

/**
 * 解绑 TOTP
 */
export async function unbindTOTPApi(data: { code: string }) {
  return requestClient.post('/auth/totp/disable', data);
}

/**
 * 发送手机验证码
 */
export async function sendSmsCodeApi(data: { phone: string }) {
  return requestClient.post('/auth/send-sms-code', data);
}

/**
 * 绑定手机
 */
export async function bindPhoneApi(data: { phone: string; code: string }) {
  return requestClient.post('/auth/bind-phone', data);
}

/**
 * 获取登录历史
 */
export async function getLoginHistoryApi() {
  return requestClient.get<any[]>('/auth/login-history');
}
