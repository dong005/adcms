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
  // 新增字段
  is_tenant: number;   // 1=租户可见 0=仅超管
  is_public: number;   // 1=公共
  type: number;        // 1=目录 2=菜单 3=页面 4=按钮
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

export function getUserMenus(userId: number) {
  return requestClient.get<MenuRecord[]>(`/menus/${userId}/menus`);
}

export function assignUserMenus(userId: number, menu_ids: number[]) {
  return requestClient.put(`/menus/${userId}/menus`, { menu_ids });
}

export interface ButtonRecord {
  title: string;
  name: string;
  permission_code: string;
  sort: number;
}

export function getMenuButtons(menuId: number) {
  return requestClient.get<MenuRecord[]>(`/menus/${menuId}/buttons`);
}

export function saveMenuButtons(menuId: number, buttons: ButtonRecord[]) {
  return requestClient.put(`/menus/${menuId}/buttons`, { buttons });
}

export function getMenuTreeWithButtons() {
  return requestClient.get<MenuRecord[]>('/menus/tree-with-buttons');
}
