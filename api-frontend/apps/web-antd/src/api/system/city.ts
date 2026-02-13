import { requestClient } from '#/api/request';

export interface CityRecord {
  id: number;
  pid: number;
  level: number;
  name: string;
  citycode: string;
  adcode: string;
  p_adcode: string;
  lng: number;
  lat: number;
  sort: number;
  created_at: string;
  children?: CityRecord[];
}

export function getCityList(pid: number = 0) {
  return requestClient.get<CityRecord[]>('/cities', { params: { pid } });
}

export function getCityTree(maxLevel: number = 2) {
  return requestClient.get<CityRecord[]>('/cities/tree', { params: { max_level: maxLevel } });
}

export function createCity(data: Partial<CityRecord>) {
  return requestClient.post('/cities', data);
}

export function updateCity(id: number, data: Partial<CityRecord>) {
  return requestClient.put(`/cities/${id}`, data);
}

export function deleteCity(id: number) {
  return requestClient.delete(`/cities/${id}`);
}
