import { requestClient } from '#/api/request';

export interface LinkRecord {
  id: number;
  tenant_id: number;
  name: string;
  url: string;
  logo: string;
  desc: string;
  sort: number;
  status: number;
  created_at: string;
}

export function getLinkList(keyword?: string) {
  return requestClient.get<LinkRecord[]>('/links', { params: { keyword } });
}

export function createLink(data: Partial<LinkRecord>) {
  return requestClient.post('/links', data);
}

export function updateLink(id: number, data: Partial<LinkRecord>) {
  return requestClient.put(`/links/${id}`, data);
}

export function deleteLink(id: number) {
  return requestClient.delete(`/links/${id}`);
}
