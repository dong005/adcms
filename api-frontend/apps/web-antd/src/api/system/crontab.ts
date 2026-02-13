import { requestClient } from '#/api/request';

export interface CrontabRecord {
  id: number;
  tenant_id: number;
  name: string;
  expression: string;
  command: string;
  status: number;
  remark: string;
  last_run_at: string | null;
  next_run_at: string | null;
  created_at: string;
}

export function getCrontabList() {
  return requestClient.get<CrontabRecord[]>('/crontabs');
}

export function createCrontab(data: Partial<CrontabRecord>) {
  return requestClient.post('/crontabs', data);
}

export function updateCrontab(id: number, data: Partial<CrontabRecord>) {
  return requestClient.put(`/crontabs/${id}`, data);
}

export function deleteCrontab(id: number) {
  return requestClient.delete(`/crontabs/${id}`);
}
