import { computed } from 'vue';
import { useAccessStore } from '@vben/stores';
import { $t } from '#/locales';

/**
 * 菜单国际化混合器
 */
export function useMenuI18n() {
  const accessStore = useAccessStore();

  // 带国际化的菜单
  const menusWithI18n = computed(() => {
    function translateMenu(menu: any): any {
      const translated = { ...menu };
      
      // 翻译菜单标题
      if (menu.meta?.title) {
        const translatedTitle = $t(`menu.${menu.name}`);
        if (translatedTitle !== `menu.${menu.name}`) {
          translated.meta = { ...menu.meta, title: translatedTitle };
        }
      }
      
      // 递归翻译子菜单
      if (menu.children && menu.children.length > 0) {
        translated.children = menu.children.map(translateMenu);
      }
      
      return translated;
    }
    
    return accessStore.accessMenus.map(translateMenu);
  });

  return {
    menusWithI18n,
  };
}
