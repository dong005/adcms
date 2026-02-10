import { requestClient } from '#/api/request';

export interface ConfigRecord {
  id: number;
  tenant_id: number;
  key: string;
  value: string;
  type: string; // string, number, boolean, json
  description: string;
  group: string;
  sort: number;
  created_at: string;
  updated_at: string;
}

export function getConfigList(params?: { group?: string }) {
  return requestClient.get<ConfigRecord[]>('/configs', { params });
}

export function getConfig(key: string) {
  return requestClient.get<string>(`/configs/${key}`);
}

export function updateConfigs(data: Record<string, any>) {
  return requestClient.put('/configs', data);
}

export function updateConfig(key: string, value: any) {
  return requestClient.put(`/configs/${key}`, { value });
}
