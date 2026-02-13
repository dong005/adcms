import { requestClient } from '#/api/request';

export interface SiteRecord {
  id: number;
  tenant_id: number;
  name: string;
  type: string;
  url: string;
  image: string;
  is_domain: number;
  status: number;
  sort: number;
  remark: string;
  created_at: string;
}

export function getSiteList(keyword?: string) {
  return requestClient.get<SiteRecord[]>('/sites', { params: { keyword } });
}

export function createSite(data: Partial<SiteRecord>) {
  return requestClient.post('/sites', data);
}

export function updateSite(id: number, data: Partial<SiteRecord>) {
  return requestClient.put(`/sites/${id}`, data);
}

export function deleteSite(id: number) {
  return requestClient.delete(`/sites/${id}`);
}
