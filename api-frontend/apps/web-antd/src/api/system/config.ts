import { requestClient } from '#/api/request';

export interface ConfigRecord {
  id: number;
  tenant_id: number;
  group_id: number;
  key: string;
  value: string;
  type: string;
  options: string;
  description: string;
  sort: number;
  created_at: string;
  updated_at: string;
}

export interface ConfigGroupRecord {
  id: number;
  name: string;
  sort: number;
  created_at: string;
}

export interface ConfigWebRecord {
  id: number;
  tenant_id: number;
  name: string;
  code: string;
  value: string;
  sort: number;
  created_at: string;
}

// System Configs
export function getConfigList(params?: { group?: string }) {
  return requestClient.get<ConfigRecord[]>('/configs', { params });
}

export function getConfigsByGroup(groupId: number) {
  return requestClient.get<ConfigRecord[]>('/configs/by-group', { params: { group_id: groupId } });
}

export function updateConfigs(data: Record<string, any>) {
  return requestClient.put('/configs', data);
}

// Config Groups
export function getConfigGroupList() {
  return requestClient.get<ConfigGroupRecord[]>('/config-groups');
}

export function createConfigGroup(data: Partial<ConfigGroupRecord>) {
  return requestClient.post('/config-groups', data);
}

export function updateConfigGroup(id: number, data: Partial<ConfigGroupRecord>) {
  return requestClient.put(`/config-groups/${id}`, data);
}

export function deleteConfigGroup(id: number) {
  return requestClient.delete(`/config-groups/${id}`);
}

// Config Webs
export function getConfigWebList() {
  return requestClient.get<ConfigWebRecord[]>('/config-webs');
}

export function saveConfigWebs(webs: Partial<ConfigWebRecord>[]) {
  return requestClient.put('/config-webs', { webs });
}

export function deleteConfigWeb(id: number) {
  return requestClient.delete(`/config-webs/${id}`);
}
