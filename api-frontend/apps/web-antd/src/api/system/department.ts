import { requestClient } from '#/api/request';

export interface DepartmentRecord {
  id: number;
  tenant_id: number;
  parent_id: number;
  name: string;
  code: string;
  leader: string;
  phone: string;
  email: string;
  sort: number;
  status: number;
  created_at: string;
  children?: DepartmentRecord[];
}

export function getDepartmentList() {
  return requestClient.get<DepartmentRecord[]>('/departments');
}

export function getDepartmentTree() {
  return requestClient.get<DepartmentRecord[]>('/departments/tree');
}

export function createDepartment(data: Partial<DepartmentRecord>) {
  return requestClient.post('/departments', data);
}

export function updateDepartment(id: number, data: Partial<DepartmentRecord>) {
  return requestClient.put(`/departments/${id}`, data);
}

export function deleteDepartment(id: number) {
  return requestClient.delete(`/departments/${id}`);
}
