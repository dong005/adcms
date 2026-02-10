import type { App, Directive, DirectiveBinding } from 'vue';
import { useAccessStore } from '@vben/stores';

export interface PermissionBinding extends DirectiveBinding {
  value: string | string[];
}

function checkPermission(el: HTMLElement, binding: PermissionBinding) {
  const { value } = binding;
  const accessStore = useAccessStore();
  const { accessCodes } = accessStore;

  if (!value) return;

  const hasPermission = Array.isArray(value)
    ? value.some(code => accessCodes.includes(code))
    : accessCodes.includes(value);

  if (!hasPermission) {
    el.style.display = 'none';
    // 或者移除元素
    // el.parentNode?.removeChild(el);
  }
}

const permission: Directive = {
  mounted(el: HTMLElement, binding: PermissionBinding) {
    checkPermission(el, binding);
  },
  updated(el: HTMLElement, binding: PermissionBinding) {
    checkPermission(el, binding);
  },
};

export function setupPermissionDirective(app: App) {
  app.directive('permission', permission);
}

export default permission;
