import { requestClient } from '#/api/request';

export interface PermissionRecord {
  id: number;
  name: string;
  code: string;
  type: number;
  parent_id: number;
  path: string;
  method: string;
  description: string;
  children?: PermissionRecord[];
}

export function getPermissionList() {
  return requestClient.get<PermissionRecord[]>('/permissions');
}

export function getPermissionTree() {
  return requestClient.get<PermissionRecord[]>('/permissions/tree');
}

export function getRolePermissions(roleId: number) {
  return requestClient.get<PermissionRecord[]>(`/roles/${roleId}/permissions`);
}

export function assignRolePermissions(roleId: number, permission_ids: number[]) {
  return requestClient.put(`/roles/${roleId}/permissions`, { permission_ids });
}

export function createPermission(data: Partial<PermissionRecord>) {
  return requestClient.post<PermissionRecord>('/permissions', data);
}

export function updatePermission(id: number, data: Partial<PermissionRecord>) {
  return requestClient.put<PermissionRecord>(`/permissions/${id}`, data);
}

export function deletePermission(id: number) {
  return requestClient.delete(`/permissions/${id}`);
}
