<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Tag, Space, Modal, Form, FormItem, Input, InputNumber,
  Select, SelectOption, Popconfirm, message, Switch, Card, DatePicker, Descriptions,
} from 'ant-design-vue';
import dayjs from 'dayjs';
import {
  getAdminList, createAdmin, updateAdmin, deleteAdmin,
  toggleAdminStatus, resetAdminPassword, getAdminStatistics,
} from '#/api/system/admin';
import { useAccess } from '@vben/access';
import type { AdminRecord, AdminStatistics } from '#/api/system/admin';

const { hasAccessByCodes } = useAccess();

const loading = ref(false);
const dataSource = ref<AdminRecord[]>([]);
const total = ref(0);
const keyword = ref('');
const pagination = reactive({ current: 1, pageSize: 10 });

const modalVisible = ref(false);
const modalTitle = ref('新增管理员');
const statisticsModalVisible = ref(false);
const currentAdmin = ref<AdminRecord | null>(null);
const statistics = ref<AdminStatistics | null>(null);

const formState = reactive({
  id: 0,
  username: '',
  password: '',
  nickname: '',
  email: '',
  phone: '',
  company: '',
  domain: '',
  expire_time: null as any,
  max_users: 0,
  status: 1,
  remark: '',
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '用户名', dataIndex: 'username', width: 120 },
  { title: '昵称', dataIndex: 'nickname', width: 120 },
  { title: '公司名称', dataIndex: 'company', width: 150 },
  { title: '域名', dataIndex: 'domain', width: 150 },
  { title: '用户数', dataIndex: 'user_count', width: 80 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '到期时间', dataIndex: 'expire_time', width: 120 },
  { title: '登录次数', dataIndex: 'login_count', width: 100 },
  { title: '最后登录', dataIndex: 'last_login_at', width: 170 },
  { title: '操作', dataIndex: 'action', width: 280, fixed: 'right' as const },
];

async function fetchData() {
  loading.value = true;
  try {
    const res = await getAdminList({
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: keyword.value,
    });
    dataSource.value = res.list || [];
    total.value = res.total;
  } finally {
    loading.value = false;
  }
}

function handleAdd() {
  modalTitle.value = '新增管理员';
  Object.assign(formState, {
    id: 0, username: '', password: '', nickname: '', email: '', phone: '',
    company: '', domain: '', expire_time: null, max_users: 0, status: 1, remark: '',
  });
  modalVisible.value = true;
}

function handleEdit(record: AdminRecord) {
  modalTitle.value = '编辑管理员';
  Object.assign(formState, {
    id: record.id,
    username: record.username,
    password: '',
    nickname: record.nickname,
    email: record.email,
    phone: record.phone,
    company: record.company,
    domain: record.domain,
    expire_time: record.expire_time ? dayjs(record.expire_time) : null,
    max_users: record.max_users,
    status: record.status,
    remark: record.remark,
  });
  modalVisible.value = true;
}

async function handleSubmit() {
  if (formState.id) {
    await updateAdmin(formState.id, {
      nickname: formState.nickname,
      email: formState.email,
      phone: formState.phone,
      company: formState.company,
      domain: formState.domain,
      expire_time: formState.expire_time?.format('YYYY-MM-DD HH:mm:ss') || null,
      max_users: formState.max_users,
      status: formState.status,
      remark: formState.remark,
    });
    message.success('更新成功');
  } else {
    await createAdmin({
      username: formState.username,
      password: formState.password,
      nickname: formState.nickname,
      email: formState.email,
      phone: formState.phone,
      company: formState.company,
      domain: formState.domain,
      expire_time: formState.expire_time?.format('YYYY-MM-DD HH:mm:ss') || null,
      max_users: formState.max_users,
      status: formState.status,
      remark: formState.remark,
    });
    message.success('创建成功');
  }
  modalVisible.value = false;
  fetchData();
}

async function handleDelete(id: number) {
  await deleteAdmin(id);
  message.success('删除成功');
  fetchData();
}

async function handleToggleStatus(record: AdminRecord) {
  const newStatus = record.status === 1 ? 0 : 1;
  await toggleAdminStatus(record.id, newStatus);
  message.success('状态更新成功');
  fetchData();
}

async function handleResetPassword(id: number) {
  try {
    await resetAdminPassword(id, '123456');
    message.success('密码已重置为 123456');
  } catch (e: any) {
    message.error(e?.message || '重置失败');
  }
}

async function handleStatistics(record: AdminRecord) {
  currentAdmin.value = record;
  statistics.value = await getAdminStatistics(record.id);
  statisticsModalVisible.value = true;
}

function handleTableChange(pag: any) {
  pagination.current = pag.current;
  pagination.pageSize = pag.pageSize;
  fetchData();
}

onMounted(() => {
  fetchData();
});
</script>

<template>
  <div class="p-4">
    <Card title="管理员管理">
      <template #extra>
        <Space>
          <Input.Search
            v-model:value="keyword"
            placeholder="搜索用户名/昵称/公司"
            style="width: 250px"
            @search="() => { pagination.current = 1; fetchData(); }"
          />
          <Button v-if="hasAccessByCodes(['admin:create'])" type="primary" @click="handleAdd">新增管理员</Button>
        </Space>
      </template>

      <Table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
        :pagination="{ current: pagination.current, pageSize: pagination.pageSize, total, showSizeChanger: true, showTotal: (t: number) => `共 ${t} 条` }"
        :scroll="{ x: 1400 }"
        row-key="id"
        size="middle"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'status'">
            <Switch
              :checked="(record as AdminRecord).status === 1"
              checked-children="启用"
              un-checked-children="禁用"
              @change="() => handleToggleStatus(record as AdminRecord)"
            />
          </template>
          <template v-if="column.dataIndex === 'expire_time'">
            <span>{{ (record as AdminRecord).expire_time || '不限' }}</span>
          </template>
          <template v-if="column.dataIndex === 'action'">
            <Space>
              <Button v-if="hasAccessByCodes(['admin:update'])" size="small" type="link" @click="handleEdit(record as AdminRecord)">编辑</Button>
              <Button v-if="hasAccessByCodes(['admin:reset-password'])" size="small" type="link" @click="handleResetPassword((record as AdminRecord).id)">重置密码</Button>
              <Button v-if="hasAccessByCodes(['admin:statistics'])" size="small" type="link" @click="handleStatistics(record as AdminRecord)">统计</Button>
              <Popconfirm v-if="hasAccessByCodes(['admin:delete'])" title="确认删除？" @confirm="handleDelete((record as AdminRecord).id)">
                <Button size="small" type="link" danger>删除</Button>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <Modal v-model:open="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="640">
      <Form :model="formState" layout="vertical" class="mt-4">
        <div class="grid grid-cols-2 gap-x-4">
          <FormItem label="用户名" :rules="[{ required: true }]">
            <Input v-model:value="formState.username" :disabled="!!formState.id" placeholder="请输入用户名" />
          </FormItem>
          <FormItem v-if="!formState.id" label="密码" :rules="[{ required: true }]">
            <Input.Password v-model:value="formState.password" placeholder="请输入密码" />
          </FormItem>
          <FormItem label="昵称">
            <Input v-model:value="formState.nickname" placeholder="请输入昵称" />
          </FormItem>
          <FormItem label="公司名称" :rules="[{ required: true }]">
            <Input v-model:value="formState.company" placeholder="请输入公司名称" />
          </FormItem>
          <FormItem label="邮箱">
            <Input v-model:value="formState.email" placeholder="请输入邮箱" />
          </FormItem>
          <FormItem label="手机">
            <Input v-model:value="formState.phone" placeholder="请输入手机号" />
          </FormItem>
          <FormItem label="域名">
            <Input v-model:value="formState.domain" placeholder="请输入绑定域名" />
          </FormItem>
          <FormItem label="到期时间">
            <DatePicker
              v-model:value="formState.expire_time"
              show-time
              format="YYYY-MM-DD HH:mm:ss"
              placeholder="不限制"
              style="width: 100%"
            />
          </FormItem>
          <FormItem label="最大用户数">
            <InputNumber v-model:value="formState.max_users" :min="0" placeholder="0=不限制" style="width: 100%" />
          </FormItem>
          <FormItem label="状态">
            <Select v-model:value="formState.status">
              <SelectOption :value="1">启用</SelectOption>
              <SelectOption :value="0">禁用</SelectOption>
            </Select>
          </FormItem>
        </div>
        <FormItem label="备注">
          <Input.TextArea v-model:value="formState.remark" placeholder="请输入备注" :rows="3" />
        </FormItem>
      </Form>
    </Modal>

    <Modal v-model:open="statisticsModalVisible" title="租户统计" :footer="null" :width="480">
      <Descriptions v-if="statistics" :column="1" bordered class="mt-4">
        <Descriptions.Item label="用户总数">{{ statistics.user_count }}</Descriptions.Item>
        <Descriptions.Item label="文章总数">{{ statistics.article_count }}</Descriptions.Item>
        <Descriptions.Item label="分类总数">{{ statistics.category_count }}</Descriptions.Item>
        <Descriptions.Item label="媒体总数">{{ statistics.media_count }}</Descriptions.Item>
        <Descriptions.Item label="登录次数">{{ statistics.login_count }}</Descriptions.Item>
        <Descriptions.Item label="最后登录">{{ statistics.last_login_at || '从未登录' }}</Descriptions.Item>
      </Descriptions>
    </Modal>
  </div>
</template>
