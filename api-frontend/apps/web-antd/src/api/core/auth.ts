import { baseRequestClient, requestClient } from '#/api/request';

export namespace AuthApi {
  /** 登录接口参数 */
  export interface LoginParams {
    password?: string;
    username?: string;
    tenant?: string;
  }

  /** 登录接口返回值 */
  export interface LoginResult {
    accessToken: string;
    require_totp?: boolean;
    temp_token?: string;
  }

  export interface RefreshTokenResult {
    data: string;
    status: number;
  }
}

/**
 * 登录
 */
export async function loginApi(data: AuthApi.LoginParams) {
  const res = await requestClient.post<{ token: string; require_totp: boolean; temp_token?: string }>('/auth/login', data);
  return { 
    accessToken: res.token, 
    require_totp: res.require_totp,
    temp_token: res.temp_token,
  } as AuthApi.LoginResult;
}

/**
 * 刷新accessToken
 */
export async function refreshTokenApi() {
  return baseRequestClient.post<AuthApi.RefreshTokenResult>('/auth/refresh', {
    withCredentials: true,
  });
}

/**
 * 退出登录
 */
export async function logoutApi() {
  return requestClient.post('/auth/logout');
}

/**
 * 获取用户权限码
 */
export async function getAccessCodesApi() {
  // 从用户信息接口获取权限码
  const res = await requestClient.get<{ permissions: string[] }>('/auth/user-info');
  return res.permissions || [];
}

/**
 * TOTP 验证（登录时使用）
 */
export async function verifyTotpApi(data: { code: string; temp_token: string }) {
  return requestClient.post<{ token: string }>('/auth/verify-totp', data);
}

/**
 * 忘记密码 - 发送验证码
 */
export async function forgotPasswordApi(data: { email: string }) {
  return baseRequestClient.post('/api/auth/forgot-password', data);
}

/**
 * 忘记密码 - 重置密码
 */
export async function resetPasswordByEmailApi(data: {
  email: string;
  code: string;
  new_password: string;
}) {
  return baseRequestClient.post('/api/auth/reset-password', data);
}
