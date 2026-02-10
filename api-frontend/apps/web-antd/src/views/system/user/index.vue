<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Input, Tag, Space, Modal, Form, FormItem,
  Select, SelectOption, Popconfirm, message, Switch, Card, TreeSelect,
} from 'ant-design-vue';
import {
  getUserList, createUser, updateUser, deleteUser,
  toggleUserStatus, resetUserPassword, assignUserRoles,
} from '#/api/system/user';
import { getRoleList } from '#/api/system/role';
import { getDepartmentTree } from '#/api/system/department';
import { useAccess } from '@vben/access';
import type { UserRecord } from '#/api/system/user';
import type { RoleRecord } from '#/api/system/role';
import type { DepartmentRecord } from '#/api/system/department';

const { hasAccessByCodes } = useAccess();

const loading = ref(false);
const dataSource = ref<UserRecord[]>([]);
const total = ref(0);
const keyword = ref('');
const pagination = reactive({ current: 1, pageSize: 10 });

const modalVisible = ref(false);
const modalTitle = ref('新增用户');
const roleModalVisible = ref(false);
const currentUserId = ref(0);
const selectedRoleIds = ref<number[]>([]);
const roleList = ref<RoleRecord[]>([]);
const deptTree = ref<DepartmentRecord[]>([]);

const formState = reactive({
  id: 0,
  username: '',
  password: '',
  nickname: '',
  email: '',
  phone: '',
  department_id: undefined as number | undefined,
  status: 1,
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '用户名', dataIndex: 'username', width: 120 },
  { title: '昵称', dataIndex: 'nickname', width: 120 },
  { title: '邮箱', dataIndex: 'email', width: 180 },
  { title: '手机', dataIndex: 'phone', width: 130 },
  { title: '部门', dataIndex: 'department_id', width: 100 },
  { title: '角色', dataIndex: 'roles', width: 150 },
  { title: 'TOTP', dataIndex: 'totp_enabled', width: 70 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '最后登录', dataIndex: 'last_login_at', width: 170 },
  { title: '操作', dataIndex: 'action', width: 280, fixed: 'right' as const },
];

async function fetchData() {
  loading.value = true;
  try {
    const res = await getUserList({
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

async function fetchRoles() {
  roleList.value = await getRoleList();
}

async function fetchDepts() {
  deptTree.value = (await getDepartmentTree()) || [];
}

function buildDeptTreeSelect(depts: DepartmentRecord[]): any[] {
  return depts.map((d) => ({
    value: d.id,
    title: d.name,
    children: d.children ? buildDeptTreeSelect(d.children) : [],
  }));
}

function handleAdd() {
  modalTitle.value = '新增用户';
  Object.assign(formState, { id: 0, username: '', password: '', nickname: '', email: '', phone: '', department_id: undefined, status: 1 });
  modalVisible.value = true;
}

function handleEdit(record: UserRecord) {
  modalTitle.value = '编辑用户';
  Object.assign(formState, {
    id: record.id,
    username: record.username,
    password: '',
    nickname: record.nickname,
    email: record.email,
    phone: record.phone,
    department_id: (record as any).department_id || undefined,
    status: record.status,
  });
  modalVisible.value = true;
}

async function handleSubmit() {
  if (formState.id) {
    await updateUser(formState.id, {
      nickname: formState.nickname,
      email: formState.email,
      phone: formState.phone,
      department_id: formState.department_id || 0,
      status: formState.status,
    } as any);
    message.success('更新成功');
  } else {
    await createUser({
      username: formState.username,
      password: formState.password,
      nickname: formState.nickname,
      email: formState.email,
      phone: formState.phone,
      department_id: formState.department_id || 0,
      status: formState.status,
    } as any);
    message.success('创建成功');
  }
  modalVisible.value = false;
  fetchData();
}

async function handleDelete(id: number) {
  await deleteUser(id);
  message.success('删除成功');
  fetchData();
}

async function handleToggleStatus(record: UserRecord) {
  const newStatus = record.status === 1 ? 0 : 1;
  await toggleUserStatus(record.id, newStatus);
  message.success('状态更新成功');
  fetchData();
}

async function handleResetPassword(id: number) {
  await resetUserPassword(id);
  message.success('密码已重置为 123456');
}

function handleAssignRoles(record: UserRecord) {
  currentUserId.value = record.id;
  selectedRoleIds.value = (record.roles || []).map((r) => r.id);
  roleModalVisible.value = true;
}

async function handleSaveRoles() {
  await assignUserRoles(currentUserId.value, selectedRoleIds.value);
  message.success('角色分配成功');
  roleModalVisible.value = false;
  fetchData();
}

function handleTableChange(pag: any) {
  pagination.current = pag.current;
  pagination.pageSize = pag.pageSize;
  fetchData();
}

onMounted(() => {
  fetchData();
  fetchRoles();
  fetchDepts();
});
</script>

<template>
  <div class="p-4">
    <Card title="用户管理">
      <template #extra>
        <Space>
          <Input.Search
            v-model:value="keyword"
            placeholder="搜索用户名/昵称/邮箱"
            style="width: 250px"
            @search="() => { pagination.current = 1; fetchData(); }"
          />
          <Button v-if="hasAccessByCodes(['user:create'])" type="primary" @click="handleAdd">新增用户</Button>
        </Space>
      </template>

      <Table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
        :pagination="{ current: pagination.current, pageSize: pagination.pageSize, total, showSizeChanger: true, showTotal: (t: number) => `共 ${t} 条` }"
        :scroll="{ x: 1200 }"
        row-key="id"
        size="middle"
        @change="handleTableChange"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'roles'">
            <Tag v-for="role in (record as UserRecord).roles" :key="role.id" color="blue">{{ role.name }}</Tag>
          </template>
          <template v-if="column.dataIndex === 'totp_enabled'">
            <Tag :color="(record as UserRecord).totp_enabled === 1 ? 'green' : 'default'">
              {{ (record as UserRecord).totp_enabled === 1 ? '已启用' : '未启用' }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'status'">
            <Switch
              :checked="(record as UserRecord).status === 1"
              checked-children="启用"
              un-checked-children="禁用"
              @change="() => handleToggleStatus(record as UserRecord)"
            />
          </template>
          <template v-if="column.dataIndex === 'action'">
            <Space>
              <Button v-if="hasAccessByCodes(['user:update'])" size="small" type="link" @click="handleEdit(record as UserRecord)">编辑</Button>
              <Button v-if="hasAccessByCodes(['user:assign-roles'])" size="small" type="link" @click="handleAssignRoles(record as UserRecord)">分配角色</Button>
              <Popconfirm v-if="hasAccessByCodes(['user:reset-password'])" title="确认重置密码？" @confirm="handleResetPassword((record as UserRecord).id)">
                <Button size="small" type="link">重置密码</Button>
              </Popconfirm>
              <Popconfirm v-if="hasAccessByCodes(['user:delete'])" title="确认删除？" @confirm="handleDelete((record as UserRecord).id)">
                <Button size="small" type="link" danger>删除</Button>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <Modal v-model:open="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="520">
      <Form :model="formState" layout="vertical" class="mt-4">
        <FormItem label="用户名" :rules="[{ required: true }]">
          <Input v-model:value="formState.username" :disabled="!!formState.id" placeholder="请输入用户名" />
        </FormItem>
        <FormItem v-if="!formState.id" label="密码" :rules="[{ required: true }]">
          <Input.Password v-model:value="formState.password" placeholder="请输入密码" />
        </FormItem>
        <FormItem label="昵称">
          <Input v-model:value="formState.nickname" placeholder="请输入昵称" />
        </FormItem>
        <FormItem label="邮箱">
          <Input v-model:value="formState.email" placeholder="请输入邮箱" />
        </FormItem>
        <FormItem label="手机">
          <Input v-model:value="formState.phone" placeholder="请输入手机号" />
        </FormItem>
        <FormItem label="部门">
          <TreeSelect
            v-model:value="formState.department_id"
            :tree-data="buildDeptTreeSelect(deptTree)"
            placeholder="请选择部门"
            allow-clear
            tree-default-expand-all
          />
        </FormItem>
        <FormItem label="状态">
          <Select v-model:value="formState.status">
            <SelectOption :value="1">启用</SelectOption>
            <SelectOption :value="0">禁用</SelectOption>
          </Select>
        </FormItem>
      </Form>
    </Modal>

    <Modal v-model:open="roleModalVisible" title="分配角色" @ok="handleSaveRoles" :width="420">
      <Select
        v-model:value="selectedRoleIds"
        mode="multiple"
        placeholder="请选择角色"
        style="width: 100%"
        class="mt-4"
      >
        <SelectOption v-for="role in roleList" :key="role.id" :value="role.id">{{ role.name }}</SelectOption>
      </Select>
    </Modal>
  </div>
</template>
