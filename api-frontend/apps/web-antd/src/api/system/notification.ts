import { requestClient } from '#/api/request';

export interface NotificationRecord {
  id: number;
  tenant_id: number;
  sender_id: number;
  sender_name: string;
  receiver_id: number;
  title: string;
  content: string;
  type: string;
  is_read: number;
  read_at: string;
  extra: string;
  reply_to_id: number;
  created_at: string;
}

export interface PageResult<T> {
  list: T[];
  total: number;
  page: number;
  page_size: number;
}

export function getNotificationList(params: {
  page: number;
  page_size: number;
  type?: string;
  is_read?: number;
}) {
  return requestClient.get<PageResult<NotificationRecord>>('/notifications', { params });
}

export function getUnreadCount() {
  return requestClient.get<{ count: number }>('/notifications/unread-count');
}

export function markAsRead(id: number) {
  return requestClient.put(`/notifications/${id}/read`);
}

export function markAllAsRead() {
  return requestClient.put('/notifications/read-all');
}

export function deleteNotification(id: number) {
  return requestClient.delete(`/notifications/${id}`);
}

export function sendNotification(data: {
  receiver_ids?: number[];
  role_ids?: number[];
  title: string;
  content?: string;
  type?: string;
}) {
  return requestClient.post('/notifications/send', data);
}

export function getNotificationDetail(id: number) {
  return requestClient.get<{
    notification: NotificationRecord;
    replies: NotificationRecord[];
  }>(`/notifications/${id}`);
}

export function replyNotification(id: number, content: string) {
  return requestClient.post(`/notifications/${id}/reply`, { content });
}

export function deleteReply(id: number) {
  return requestClient.delete(`/notifications/reply/${id}`);
}
