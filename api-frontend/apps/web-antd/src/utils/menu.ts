import { $t } from '#/locales';

/**
 * 翻译菜单标题
 * @param title 原始标题
 * @param name 菜单名称（用于翻译键）
 * @returns 翻译后的标题
 */
export function translateMenuTitle(title: string, name: string): string {
  try {
    // 尝试从语言包中获取翻译
    const translated = $t(`menu.${name}`);
    // 如果翻译键不存在，返回原始标题
    return translated !== `menu.${name}` ? translated : title;
  } catch {
    // 如果翻译失败，返回原始标题
    return title;
  }
}
