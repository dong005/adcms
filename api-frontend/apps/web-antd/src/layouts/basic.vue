<script lang="ts" setup>
import type { NotificationItem } from '@vben/layouts';

import { computed, onMounted, onUnmounted, ref, watch } from 'vue';
import { useRouter } from 'vue-router';

import { AuthenticationLoginExpiredModal } from '@vben/common-ui';
import { useWatermark } from '@vben/hooks';
import {
  BasicLayout,
  LockScreen,
  Notification,
  UserDropdown,
} from '@vben/layouts';
import { preferences } from '@vben/preferences';
import { useAccessStore, useUserStore } from '@vben/stores';

import { $t } from '#/locales';
import { useAuthStore } from '#/store';
import LoginForm from '#/views/_core/authentication/login.vue';
import {
  getNotificationList,
  getUnreadCount,
  markAsRead,
  markAllAsRead,
  deleteNotification,
  getNotificationDetail,
  replyNotification,
} from '#/api/system/notification';
import type { NotificationRecord } from '#/api/system/notification';
import { Modal, Input, Button, Divider, Tag, message } from 'ant-design-vue';

const notifications = ref<NotificationItem[]>([]);
const unreadCount = ref(0);
let pollTimer: ReturnType<typeof setInterval> | null = null;

async function fetchNotifications() {
  try {
    const res = await getNotificationList({ page: 1, page_size: 10 });
    notifications.value = (res.list || []).map((n: any) => ({
      id: n.id,
      avatar: `https://avatar.vercel.sh/${encodeURIComponent(n.sender_name || 'system')}.svg?text=${(n.sender_name || '系')[0]}`,
      date: n.created_at,
      isRead: n.is_read === 1,
      message: n.content || '',
      title: n.title,
    }));
  } catch {
    // ignore
  }
}

async function fetchUnread() {
  try {
    const res = await getUnreadCount();
    unreadCount.value = res.count;
  } catch {
    // ignore
  }
}

async function loadNotifications() {
  await Promise.all([fetchNotifications(), fetchUnread()]);
}

const router = useRouter();
const userStore = useUserStore();
const authStore = useAuthStore();
const accessStore = useAccessStore();
const { destroyWatermark, updateWatermark } = useWatermark();
const showDot = computed(() => unreadCount.value > 0);

onMounted(() => {
  loadNotifications();
  pollTimer = setInterval(fetchUnread, 60000);
});

onUnmounted(() => {
  if (pollTimer) clearInterval(pollTimer);
});

const menus = computed(() => [
  {
    handler: () => {
      router.push({ name: 'Profile' });
    },
    icon: 'lucide:user',
    text: $t('page.auth.profile'),
  },
]);

const avatar = computed(() => {
  return userStore.userInfo?.avatar ?? preferences.app.defaultAvatar;
});

async function handleLogout() {
  await authStore.logout(false);
}

async function handleNoticeClear() {
  await markAllAsRead();
  notifications.value = [];
  unreadCount.value = 0;
}

async function markRead(id: number | string) {
  await markAsRead(Number(id));
  const item = notifications.value.find((item) => item.id === id);
  if (item) {
    item.isRead = true;
  }
  unreadCount.value = Math.max(0, unreadCount.value - 1);
}

async function remove(id: number | string) {
  await deleteNotification(Number(id));
  notifications.value = notifications.value.filter((item) => item.id !== id);
  fetchUnread();
}

async function handleMakeAll() {
  await markAllAsRead();
  notifications.value.forEach((item) => (item.isRead = true));
  unreadCount.value = 0;
}

function handleViewAll() {
  router.push('/system/notification');
}

// 详情/回复弹窗
const detailVisible = ref(false);
const currentNotif = ref<NotificationRecord | null>(null);
const replies = ref<NotificationRecord[]>([]);
const replyContent = ref('');
const replyLoading = ref(false);
const detailLoading = ref(false);

const typeMap: Record<string, { label: string; color: string }> = {
  system: { label: '系统', color: 'blue' },
  task: { label: '任务', color: 'orange' },
  message: { label: '消息', color: 'green' },
};

async function handleNotifClick(item: NotificationItem) {
  if (!item.id) return;
  detailVisible.value = true;
  detailLoading.value = true;
  replyContent.value = '';
  try {
    const res = await getNotificationDetail(Number(item.id));
    currentNotif.value = res.notification;
    replies.value = res.replies || [];
    // 自动标记已读
    if (!item.isRead) {
      await markAsRead(Number(item.id));
      item.isRead = true;
      unreadCount.value = Math.max(0, unreadCount.value - 1);
    }
  } catch {
    currentNotif.value = null;
    replies.value = [];
  } finally {
    detailLoading.value = false;
  }
}

async function handleReply() {
  if (!replyContent.value.trim() || !currentNotif.value) return;
  replyLoading.value = true;
  try {
    await replyNotification(currentNotif.value.id, replyContent.value);
    message.success('回复成功');
    replyContent.value = '';
    const res = await getNotificationDetail(currentNotif.value.id);
    replies.value = res.replies || [];
  } catch (e: any) {
    message.error(e?.message || '回复失败');
  } finally {
    replyLoading.value = false;
  }
}
watch(
  () => ({
    enable: preferences.app.watermark,
    content: preferences.app.watermarkContent,
  }),
  async ({ enable, content }) => {
    if (enable) {
      await updateWatermark({
        content:
          content ||
          `${userStore.userInfo?.username} - ${userStore.userInfo?.realName}`,
      });
    } else {
      destroyWatermark();
    }
  },
  {
    immediate: true,
  },
);
</script>

<template>
  <BasicLayout @clear-preferences-and-logout="handleLogout">
    <template #user-dropdown>
      <UserDropdown
        :avatar
        :menus
        :text="userStore.userInfo?.realName"
        :description="userStore.userInfo?.email || ''"
        tag-text="Pro"
        @logout="handleLogout"
      />
    </template>
    <template #notification>
      <Notification
        :dot="showDot"
        :notifications="notifications"
        @clear="handleNoticeClear"
        @click="handleNotifClick"
        @read="(item) => item.id && markRead(item.id)"
        @remove="(item) => item.id && remove(item.id)"
        @make-all="handleMakeAll"
        @view-all="handleViewAll"
      />
    </template>
    <template #extra>
      <AuthenticationLoginExpiredModal
        v-model:open="accessStore.loginExpired"
        :avatar
      >
        <LoginForm />
      </AuthenticationLoginExpiredModal>
    </template>
    <template #lock-screen>
      <LockScreen :avatar @to-login="handleLogout" />
    </template>
  </BasicLayout>

  <!-- 消息详情/回复弹窗 -->
  <Modal v-model:open="detailVisible" title="消息详情" :footer="null" :width="600">
    <div v-if="currentNotif" class="mt-4">
      <div style="margin-bottom: 12px;">
        <Tag :color="typeMap[currentNotif.type]?.color || 'default'">{{ typeMap[currentNotif.type]?.label || currentNotif.type }}</Tag>
        <span style="font-weight: bold; font-size: 16px; margin-left: 8px;">{{ currentNotif.title }}</span>
      </div>
      <div style="color: #999; font-size: 12px; margin-bottom: 12px;">
        发送者: {{ currentNotif.sender_name || '系统' }} · {{ currentNotif.created_at }}
      </div>
      <div style="padding: 12px; border: 1px solid var(--border); border-radius: 4px; margin-bottom: 16px; white-space: pre-wrap;">
        {{ currentNotif.content || '（无内容）' }}
      </div>

      <Divider v-if="replies.length > 0" orientation="left" style="font-size: 13px;">回复 ({{ replies.length }})</Divider>
      <div v-for="reply in replies" :key="reply.id" style="padding: 8px 12px; margin-bottom: 8px; border: 1px solid var(--border); border-radius: 4px; border-left: 3px solid #1890ff;">
        <div style="font-size: 12px; opacity: 0.6; margin-bottom: 4px;">
          {{ reply.sender_name || '系统' }} · {{ reply.created_at }}
        </div>
        <div style="white-space: pre-wrap;">{{ reply.content }}</div>
      </div>

      <Divider orientation="left" style="font-size: 13px;">回复消息</Divider>
      <Input.TextArea v-model:value="replyContent" placeholder="输入回复内容..." :rows="3" />
      <div style="margin-top: 8px; text-align: right;">
        <Button type="primary" :loading="replyLoading" @click="handleReply" :disabled="!replyContent.trim()">发送回复</Button>
      </div>
    </div>
  </Modal>
</template>
