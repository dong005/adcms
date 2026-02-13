import { requestClient } from '#/api/request';
import type { UserRecord } from './user';

export interface AdminRecord extends UserRecord {
  user_count: number;
  article_count: number;
  category_count: number;
  media_count: number;
}

export interface AdminStatistics {
  user_count: number;
  article_count: number;
  category_count: number;
  media_count: number;
  login_count: number;
  last_login_at: string | null;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

export function getAdminList(params: { page: number; page_size: number; keyword?: string }) {
  return requestClient.get<PageResult<AdminRecord>>('/admins', { params });
}

export function getAdminDetail(id: number) {
  return requestClient.get<AdminRecord>(`/admins/${id}`);
}

export function createAdmin(data: Partial<UserRecord> & { password: string }) {
  return requestClient.post('/admins', data);
}

export function updateAdmin(id: number, data: Partial<UserRecord>) {
  return requestClient.put(`/admins/${id}`, data);
}

export function deleteAdmin(id: number) {
  return requestClient.delete(`/admins/${id}`);
}

export function toggleAdminStatus(id: number, status: number) {
  return requestClient.put(`/admins/${id}/status`, { status });
}

export function resetAdminPassword(id: number, password: string) {
  return requestClient.put(`/admins/${id}/reset-password`, { password });
}

export function getAdminStatistics(id: number) {
  return requestClient.get<AdminStatistics>(`/admins/${id}/statistics`);
}
