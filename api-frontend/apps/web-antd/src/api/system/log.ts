import { requestClient } from '#/api/request';

export interface OperationLogRecord {
  id: number;
  tenant_id: number;
  user_id: number;
  username: string;
  module: string;
  action: string;
  method: string;
  path: string;
  ip: string;
  user_agent: string;
  params: string;
  result: string;
  status: number; // 1-成功 0-失败
  error: string;
  duration: number;
  created_at: string;
}

export interface LoginLogRecord {
  id: number;
  tenant_id: number;
  user_id: number;
  username: string;
  ip: string;
  user_agent: string;
  login_type: string; // password, totp
  status: number; // 1-成功 0-失败
  error: string;
  created_at: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

// 操作日志
export function getOperationLogList(params: {
  page: number;
  page_size: number;
  keyword?: string;
  module?: string;
  status?: number;
  start_time?: string;
  end_time?: string;
}) {
  return requestClient.get<PageResult<OperationLogRecord>>('/logs/operation', { params });
}

// 登录日志
export function getLoginLogList(params: {
  page: number;
  page_size: number;
  keyword?: string;
  status?: number;
  start_time?: string;
  end_time?: string;
}) {
  return requestClient.get<PageResult<LoginLogRecord>>('/logs/login', { params });
}

export interface EmailLogRecord {
  id: number;
  tenant_id: number;
  to: string;
  subject: string;
  status: number;
  error: string;
  created_at: string;
}

export interface SmsLogRecord {
  id: number;
  tenant_id: number;
  phone: string;
  template_id: string;
  params: string;
  status: number;
  error: string;
  created_at: string;
}

// 邮件日志
export function getEmailLogList(params: {
  page: number;
  page_size: number;
  keyword?: string;
}) {
  return requestClient.get<PageResult<EmailLogRecord>>('/logs/email', { params });
}

// 短信日志
export function getSmsLogList(params: {
  page: number;
  page_size: number;
  keyword?: string;
}) {
  return requestClient.get<PageResult<SmsLogRecord>>('/logs/sms', { params });
}

export interface LogConfig {
  log_operation_enabled: string;
  log_login_enabled: string;
  log_email_enabled: string;
  log_sms_enabled: string;
}

// 日志配置
export function getLogConfig() {
  return requestClient.get<LogConfig>('/configs/log');
}

export function updateLogConfig(data: LogConfig) {
  return requestClient.put('/configs/log', data);
}
