<script lang="ts" setup>
import { ref, reactive, computed, onMounted } from 'vue';
import {
  Table, Button, Tag, Space, Popconfirm, message, Card,
  Select, SelectOption, Badge, Modal, Input, Form, FormItem,
  Checkbox, CheckboxGroup, Divider,
} from 'ant-design-vue';
import {
  getNotificationList, markAsRead, markAllAsRead,
  deleteNotification, getUnreadCount, sendNotification,
  getNotificationDetail, replyNotification, deleteReply,
} from '#/api/system/notification';
import { getRoleList } from '#/api/system/role';
import { requestClient } from '#/api/request';
import type { NotificationRecord } from '#/api/system/notification';
import { useUserStore } from '@vben/stores';

const userStore = useUserStore();
const isSuperAdmin = computed(() => (userStore.userInfo?.roles as string[] || []).includes('super_admin'));

const loading = ref(false);
const dataSource = ref<NotificationRecord[]>([]);
const total = ref(0);
const unreadCount = ref(0);
const pagination = reactive({ current: 1, pageSize: 20 });
const filterType = ref('');
const filterRead = ref(-1);

// 发送弹窗
const sendVisible = ref(false);
const sendLoading = ref(false);
const sendForm = reactive({
  title: '',
  content: '',
  type: 'system',
  receiver_ids: [] as number[],
  role_ids: [] as number[],
  selectMode: 'role' as 'role' | 'user',
});
const roleList = ref<any[]>([]);
const userList = ref<any[]>([]);

// 详情/回复弹窗
const detailVisible = ref(false);
const detailLoading = ref(false);
const currentNotif = ref<NotificationRecord | null>(null);
const replies = ref<NotificationRecord[]>([]);
const replyContent = ref('');
const replyLoading = ref(false);

const columns = [
  { title: '发送者', dataIndex: 'sender_name', width: 100 },
  { title: '标题', dataIndex: 'title', width: 220 },
  { title: '内容', dataIndex: 'content', ellipsis: true },
  { title: '类型', dataIndex: 'type', width: 80 },
  { title: '状态', dataIndex: 'is_read', width: 70 },
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', dataIndex: 'action', width: 200, fixed: 'right' as const },
];

const typeMap: Record<string, { label: string; color: string }> = {
  system: { label: '系统', color: 'blue' },
  task: { label: '任务', color: 'orange' },
  message: { label: '消息', color: 'green' },
};

async function fetchData() {
  loading.value = true;
  try {
    const params: any = {
      page: pagination.current,
      page_size: pagination.pageSize,
    };
    if (filterType.value) params.type = filterType.value;
    if (filterRead.value >= 0) params.is_read = filterRead.value;

    const res = await getNotificationList(params);
    dataSource.value = res.list || [];
    total.value = res.total;
  } finally {
    loading.value = false;
  }
}

async function fetchUnread() {
  const res = await getUnreadCount();
  unreadCount.value = res.count;
}

async function handleRead(id: number) {
  await markAsRead(id);
  message.success('已标记为已读');
  fetchData();
  fetchUnread();
}

async function handleReadAll() {
  await markAllAsRead();
  message.success('全部标记为已读');
  fetchData();
  fetchUnread();
}

async function handleDelete(id: number) {
  await deleteNotification(id);
  message.success('删除成功');
  fetchData();
  fetchUnread();
}

function handleTableChange(pag: any) {
  pagination.current = pag.current;
  pagination.pageSize = pag.pageSize;
  fetchData();
}

function handleFilter() {
  pagination.current = 1;
  fetchData();
}

// 发送消息
async function openSendModal() {
  sendForm.title = '';
  sendForm.content = '';
  sendForm.type = 'system';
  sendForm.receiver_ids = [];
  sendForm.role_ids = [];
  sendForm.selectMode = 'role';
  sendVisible.value = true;
  if (roleList.value.length === 0) {
    roleList.value = await getRoleList();
  }
  if (userList.value.length === 0) {
    const res = await requestClient.get<any>('/users?page=1&page_size=1000');
    userList.value = (res.list || []).map((u: any) => ({ id: u.id, label: u.nickname || u.username }));
  }
}

async function handleSend() {
  if (!sendForm.title.trim()) {
    message.warning('请输入标题');
    return;
  }
  if (sendForm.selectMode === 'role' && sendForm.role_ids.length === 0) {
    message.warning('请选择角色');
    return;
  }
  if (sendForm.selectMode === 'user' && sendForm.receiver_ids.length === 0) {
    message.warning('请选择用户');
    return;
  }
  sendLoading.value = true;
  try {
    await sendNotification({
      title: sendForm.title,
      content: sendForm.content,
      type: sendForm.type,
      receiver_ids: sendForm.selectMode === 'user' ? sendForm.receiver_ids : undefined,
      role_ids: sendForm.selectMode === 'role' ? sendForm.role_ids : undefined,
    });
    message.success('发送成功');
    sendVisible.value = false;
    fetchData();
  } catch (e: any) {
    message.error(e?.message || '发送失败');
  } finally {
    sendLoading.value = false;
  }
}

// 详情/回复
async function openDetail(record: NotificationRecord) {
  detailLoading.value = true;
  detailVisible.value = true;
  replyContent.value = '';
  try {
    const res = await getNotificationDetail(record.id);
    currentNotif.value = res.notification;
    replies.value = res.replies || [];
    // 自动标记已读
    if (record.is_read === 0) {
      await markAsRead(record.id);
      fetchData();
      fetchUnread();
    }
  } finally {
    detailLoading.value = false;
  }
}

async function handleReply() {
  if (!replyContent.value.trim()) {
    message.warning('请输入回复内容');
    return;
  }
  if (!currentNotif.value) return;
  replyLoading.value = true;
  try {
    await replyNotification(currentNotif.value.id, replyContent.value);
    message.success('回复成功');
    replyContent.value = '';
    // 刷新回复列表
    const res = await getNotificationDetail(currentNotif.value.id);
    replies.value = res.replies || [];
  } catch (e: any) {
    message.error(e?.message || '回复失败');
  } finally {
    replyLoading.value = false;
  }
}

async function handleDeleteReply(replyId: number) {
  try {
    await deleteReply(replyId);
    message.success('回复已删除');
    if (currentNotif.value) {
      const res = await getNotificationDetail(currentNotif.value.id);
      replies.value = res.replies || [];
    }
  } catch (e: any) {
    message.error(e?.message || '删除失败');
  }
}

onMounted(() => {
  fetchData();
  fetchUnread();
});
</script>

<template>
  <div class="p-4">
    <Card>
      <template #title>
        <Space>
          <span>消息管理</span>
          <Badge :count="unreadCount" :offset="[6, -2]" />
        </Space>
      </template>
      <template #extra>
        <Space>
          <Select v-model:value="filterType" placeholder="类型" style="width: 100px" allow-clear @change="handleFilter">
            <SelectOption value="">全部</SelectOption>
            <SelectOption value="system">系统</SelectOption>
            <SelectOption value="task">任务</SelectOption>
            <SelectOption value="message">消息</SelectOption>
          </Select>
          <Select v-model:value="filterRead" style="width: 100px" @change="handleFilter">
            <SelectOption :value="-1">全部</SelectOption>
            <SelectOption :value="0">未读</SelectOption>
            <SelectOption :value="1">已读</SelectOption>
          </Select>
          <Button @click="handleReadAll" :disabled="unreadCount === 0">全部已读</Button>
          <Button v-if="isSuperAdmin" type="primary" @click="openSendModal">发送消息</Button>
        </Space>
      </template>

      <Table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
        :pagination="{ current: pagination.current, pageSize: pagination.pageSize, total, showSizeChanger: true, showTotal: (t: number) => `共 ${t} 条` }"
        row-key="id"
        size="middle"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'title'">
            <a :style="{ fontWeight: (record as NotificationRecord).is_read === 0 ? 'bold' : 'normal' }" @click="openDetail(record as NotificationRecord)">
              {{ (record as NotificationRecord).title }}
            </a>
          </template>
          <template v-if="column.dataIndex === 'type'">
            <Tag :color="typeMap[(record as NotificationRecord).type]?.color || 'default'">
              {{ typeMap[(record as NotificationRecord).type]?.label || (record as NotificationRecord).type }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'is_read'">
            <Tag :color="(record as NotificationRecord).is_read === 1 ? 'default' : 'blue'">
              {{ (record as NotificationRecord).is_read === 1 ? '已读' : '未读' }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'action'">
            <Space>
              <Button size="small" type="link" @click="openDetail(record as NotificationRecord)">查看</Button>
              <Button
                v-if="(record as NotificationRecord).is_read === 0"
                size="small" type="link"
                @click="handleRead((record as NotificationRecord).id)"
              >已读</Button>
              <Popconfirm title="确认删除？" @confirm="handleDelete((record as NotificationRecord).id)">
                <Button size="small" type="link" danger>删除</Button>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <!-- 发送消息弹窗 -->
    <Modal v-model:open="sendVisible" title="发送消息" :confirm-loading="sendLoading" @ok="handleSend" :width="560">
      <Form layout="vertical" class="mt-4">
        <FormItem label="收件方式">
          <Select v-model:value="sendForm.selectMode" style="width: 200px">
            <SelectOption value="role">按角色</SelectOption>
            <SelectOption value="user">按用户</SelectOption>
          </Select>
        </FormItem>
        <FormItem v-if="sendForm.selectMode === 'role'" label="选择角色">
          <CheckboxGroup v-model:value="sendForm.role_ids">
            <div v-for="role in roleList" :key="role.id" style="padding: 2px 0;">
              <Checkbox :value="role.id">{{ role.name }} ({{ role.code }})</Checkbox>
            </div>
          </CheckboxGroup>
        </FormItem>
        <FormItem v-if="sendForm.selectMode === 'user'" label="选择用户">
          <Select
            v-model:value="sendForm.receiver_ids"
            mode="multiple"
            placeholder="搜索选择用户"
            style="width: 100%"
            :options="userList.map(u => ({ value: u.id, label: u.label }))"
            show-search
            :filter-option="(input: string, option: any) => option.label.toLowerCase().includes(input.toLowerCase())"
          />
        </FormItem>
        <FormItem label="消息类型">
          <Select v-model:value="sendForm.type" style="width: 200px">
            <SelectOption value="system">系统通知</SelectOption>
            <SelectOption value="task">任务通知</SelectOption>
            <SelectOption value="message">普通消息</SelectOption>
          </Select>
        </FormItem>
        <FormItem label="标题">
          <Input v-model:value="sendForm.title" placeholder="请输入消息标题" />
        </FormItem>
        <FormItem label="内容">
          <Input.TextArea v-model:value="sendForm.content" placeholder="请输入消息内容" :rows="4" />
        </FormItem>
      </Form>
    </Modal>

    <!-- 详情/回复弹窗 -->
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
        <div v-for="reply in replies" :key="reply.id" style="padding: 8px 12px; margin-bottom: 8px; border: 1px solid var(--border); border-radius: 4px; border-left: 3px solid #1890ff; position: relative;">
          <div style="font-size: 12px; opacity: 0.6; margin-bottom: 4px;">
            {{ reply.sender_name || '系统' }} · {{ reply.created_at }}
            <Popconfirm v-if="isSuperAdmin" title="确认删除这条回复？" @confirm="handleDeleteReply(reply.id)">
              <Button type="link" size="small" danger style="float: right; padding: 0; height: auto;">删除</Button>
            </Popconfirm>
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
  </div>
</template>
