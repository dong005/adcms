import type { RouteRecordStringComponent } from '@vben/types';

import { requestClient } from '#/api/request';

// 后端菜单类型定义（MenuTree 结构）
interface BackendMenuMeta {
  title: string;
  icon?: string;
  hideInMenu?: boolean;
  hideInTab?: boolean;
  hideInBreadcrumb?: boolean;
  keepAlive?: boolean;
  frameSrc?: string;
}

interface BackendMenu {
  id: number;
  parentId: number;
  name: string;
  path: string;
  component?: string;
  redirect?: string;
  meta: BackendMenuMeta;
  children?: BackendMenu[];
}

/**
 * 将后端菜单格式转换为前端路由格式
 */
function transformMenuToRoute(menu: BackendMenu): RouteRecordStringComponent {
  const route: RouteRecordStringComponent = {
    path: menu.path,
    name: menu.name,
    component: menu.component || 'BasicLayout',
    meta: {
      title: menu.meta?.title || menu.name,
      icon: menu.meta?.icon,
      hideInMenu: menu.meta?.hideInMenu || false,
      hideInTab: menu.meta?.hideInTab || false,
      hideInBreadcrumb: menu.meta?.hideInBreadcrumb || false,
      keepAlive: menu.meta?.keepAlive !== false,
      frameSrc: menu.meta?.frameSrc,
      order: 0,
    },
  };

  if (menu.redirect) {
    route.redirect = menu.redirect;
  }

  if (menu.children && menu.children.length > 0) {
    route.children = menu.children.map(child => transformMenuToRoute(child));
  }

  return route;
}

/**
 * 获取用户所有菜单
 */
export async function getAllMenusApi() {
  const menus: BackendMenu[] = await requestClient.get('/menus/user');
  
  // 转换菜单格式
  const routes = menus.map(menu => transformMenuToRoute(menu));
  
  return routes;
}
