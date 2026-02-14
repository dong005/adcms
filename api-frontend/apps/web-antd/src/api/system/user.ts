import { requestClient } from '#/api/request';

export interface UserRecord {
  id: number;
  tenant_id: number;
  username: string;
  email: string;
  phone: string;
  nickname: string;
  avatar: string;
  status: number;
  totp_enabled: number;
  last_login_at: string;
  last_login_ip: string;
  created_at: string;
  roles?: RoleRecord[];
  department_id?: number;
  // 新增字段
  is_admin: number;      // 0=普通 1=管理员 2=超管
  company: string;
  domain: string;
  expire_time: string | null;
  max_users: number;
  login_count: number;
  remark: string;
}

export interface RoleRecord {
  id: number;
  name: string;
  code: string;
  status: number;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

export function getUserList(params: { page: number; page_size: number; keyword?: string }) {
  return requestClient.get<PageResult<UserRecord>>('/users', { params });
}

export function getUserDetail(id: number) {
  return requestClient.get<UserRecord>(`/users/${id}`);
}

export function createUser(data: Partial<UserRecord> & { password: string; role_ids?: number[] }) {
  return requestClient.post('/users', data);
}

export function updateUser(id: number, data: Partial<UserRecord>) {
  return requestClient.put(`/users/${id}`, data);
}

export function deleteUser(id: number) {
  return requestClient.delete(`/users/${id}`);
}

export function toggleUserStatus(id: number, status: number) {
  return requestClient.put(`/users/${id}/status`, { status });
}

export function resetUserPassword(id: number) {
  return requestClient.put(`/users/${id}/reset-password`);
}

export function assignUserRoles(id: number, role_ids: number[]) {
  return requestClient.put(`/users/${id}/roles`, { role_ids });
}

export function assignUserMenus(id: number, menu_ids: number[]) {
  return requestClient.put(`/users/${id}/menus`, { menu_ids });
}

export function unlockUser(id: number) {
  return requestClient.put(`/users/${id}/unlock`);
}

export function loginAsUser(id: number) {
  return requestClient.post<{ token: string }>(`/users/${id}/login-as`);
}
