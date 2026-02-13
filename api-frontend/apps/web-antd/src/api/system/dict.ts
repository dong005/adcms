import { requestClient } from '#/api/request';

export interface DictTypeRecord {
  id: number;
  tenant_id: number;
  name: string;
  code: string;
  sort: number;
  status: number;
  remark: string;
  created_at: string;
}

export interface DictRecord {
  id: number;
  tenant_id: number;
  dict_type_id: number;
  name: string;
  value: string;
  sort: number;
  status: number;
  remark: string;
  created_at: string;
}

export function getDictTypeList(keyword?: string) {
  return requestClient.get<DictTypeRecord[]>('/dict-types', { params: { keyword } });
}

export function createDictType(data: Partial<DictTypeRecord>) {
  return requestClient.post('/dict-types', data);
}

export function updateDictType(id: number, data: Partial<DictTypeRecord>) {
  return requestClient.put(`/dict-types/${id}`, data);
}

export function deleteDictType(id: number) {
  return requestClient.delete(`/dict-types/${id}`);
}

export function getDictList(typeId: number) {
  return requestClient.get<DictRecord[]>('/dicts', { params: { type_id: typeId } });
}

export function createDict(data: Partial<DictRecord>) {
  return requestClient.post('/dicts', data);
}

export function updateDict(id: number, data: Partial<DictRecord>) {
  return requestClient.put(`/dicts/${id}`, data);
}

export function deleteDict(id: number) {
  return requestClient.delete(`/dicts/${id}`);
}

export function getDictsByCode(code: string) {
  return requestClient.get<DictRecord[]>(`/dicts/code/${code}`);
}
