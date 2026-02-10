import { requestClient } from '#/api/request';

export interface ArticleRecord {
  id: number;
  tenant_id: number;
  title: string;
  slug: string;
  content: string;
  excerpt: string;
  featured_image: string;
  status: number; // 0-草稿 1-发布 2-下线
  author_id: number;
  category_id: number;
  view_count: number;
  created_at: string;
  updated_at: string;
  published_at?: string;
  categories?: CategoryRecord[];
  tags?: TagRecord[];
}

export interface CategoryRecord {
  id: number;
  tenant_id: number;
  name: string;
  slug: string;
  description: string;
  parent_id: number;
  sort: number;
  status: number;
  created_at: string;
  children?: CategoryRecord[];
}

export interface TagRecord {
  id: number;
  tenant_id: number;
  name: string;
  slug: string;
  description: string;
  color: string;
  created_at: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

// 文章管理
export function getArticleList(params: { 
  page: number; 
  page_size: number; 
  keyword?: string;
  status?: number;
  category_id?: number;
}) {
  return requestClient.get<PageResult<ArticleRecord>>('/articles', { params });
}

export function getArticleDetail(id: number) {
  return requestClient.get<ArticleRecord>(`/articles/${id}`);
}

export function createArticle(data: Partial<ArticleRecord> & { tag_ids?: number[] }) {
  return requestClient.post('/articles', data);
}

export function updateArticle(id: number, data: Partial<ArticleRecord> & { tag_ids?: number[] }) {
  return requestClient.put(`/articles/${id}`, data);
}

export function deleteArticle(id: number) {
  return requestClient.delete(`/articles/${id}`);
}

export function publishArticle(id: number) {
  return requestClient.put(`/articles/${id}/publish`);
}

export function unpublishArticle(id: number) {
  return requestClient.put(`/articles/${id}/unpublish`);
}

// 分类管理
export function getCategoryList() {
  return requestClient.get<CategoryRecord[]>('/categories');
}

export function getCategoryTree() {
  return requestClient.get<CategoryRecord[]>('/categories/tree');
}

export function createCategory(data: Partial<CategoryRecord>) {
  return requestClient.post('/categories', data);
}

export function updateCategory(id: number, data: Partial<CategoryRecord>) {
  return requestClient.put(`/categories/${id}`, data);
}

export function deleteCategory(id: number) {
  return requestClient.delete(`/categories/${id}`);
}

// 标签管理
export function getTagList(params: { page?: number; page_size?: number; keyword?: string }) {
  return requestClient.get<PageResult<TagRecord>>('/tags', { params });
}

export function createTag(data: Partial<TagRecord>) {
  return requestClient.post('/tags', data);
}

export function updateTag(id: number, data: Partial<TagRecord>) {
  return requestClient.put(`/tags/${id}`, data);
}

export function deleteTag(id: number) {
  return requestClient.delete(`/tags/${id}`);
}

// 媒体管理
export interface MediaRecord {
  id: number;
  tenant_id: number;
  filename: string;
  original_name: string;
  mime_type: string;
  size: number;
  path: string;
  url: string;
  created_at: string;
}

export function getMediaList(params: { 
  page: number; 
  page_size: number; 
  keyword?: string;
  type?: string;
}) {
  return requestClient.get<PageResult<MediaRecord>>('/media', { params });
}

export function uploadMedia(file: File) {
  const formData = new FormData();
  formData.append('file', file);
  return requestClient.post<MediaRecord>('/media/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
}

export function deleteMedia(id: number) {
  return requestClient.delete(`/media/${id}`);
}
