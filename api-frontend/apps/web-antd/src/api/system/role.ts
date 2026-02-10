import { requestClient } from '#/api/request';

export interface RoleRecord {
  id: number;
  tenant_id: number;
  name: string;
  code: string;
  description: string;
  status: number;
  sort: number;
  created_at: string;
}

export function getRoleList() {
  return requestClient.get<RoleRecord[]>('/roles');
}

export function getRoleDetail(id: number) {
  return requestClient.get<RoleRecord>(`/roles/${id}`);
}

export function createRole(data: Partial<RoleRecord>) {
  return requestClient.post('/roles', data);
}

export function updateRole(id: number, data: Partial<RoleRecord>) {
  return requestClient.put(`/roles/${id}`, data);
}

export function deleteRole(id: number) {
  return requestClient.delete(`/roles/${id}`);
}

export function getRoleMenus(id: number) {
  return requestClient.get<any[]>(`/roles/${id}/menus`);
}

export function assignRoleMenus(id: number, menu_ids: number[]) {
  return requestClient.put(`/roles/${id}/menus`, { menu_ids });
}
