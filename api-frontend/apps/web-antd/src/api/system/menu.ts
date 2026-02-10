import { requestClient } from '#/api/request';

export interface MenuRecord {
  id: number;
  tenant_id: number;
  parent_id: number;
  name: string;
  path: string;
  component: string;
  redirect: string;
  icon: string;
  title: string;
  hide_in_menu: number;
  hide_in_tab: number;
  hide_in_breadcrumb: number;
  keep_alive: number;
  frame_src: string;
  sort: number;
  status: number;
  permission_code: string;
  created_at: string;
  children?: MenuRecord[];
}

export function getMenuList() {
  return requestClient.get<MenuRecord[]>('/menus');
}

export function getMenuTree() {
  return requestClient.get<MenuRecord[]>('/menus/tree');
}

export function createMenu(data: Partial<MenuRecord>) {
  return requestClient.post('/menus', data);
}

export function updateMenu(id: number, data: Partial<MenuRecord>) {
  return requestClient.put(`/menus/${id}`, data);
}

export function deleteMenu(id: number) {
  return requestClient.delete(`/menus/${id}`);
}
