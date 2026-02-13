<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Tag, Space, Modal, Form, FormItem, Input, InputNumber,
  Select, SelectOption, Popconfirm, message, Switch, Card, TreeSelect,
  DatePicker, Tree,
} from 'ant-design-vue';
import dayjs from 'dayjs';
import {
  getUserList, createUser, updateUser, deleteUser,
  toggleUserStatus, resetUserPassword, assignUserRoles,
  assignUserMenus, unlockUser,
} from '#/api/system/user';
import { getRoleList } from '#/api/system/role';
import { getDepartmentTree } from '#/api/system/department';
import { getMenuList } from '#/api/system/menu';
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
const menuModalVisible = ref(false);
const currentUserId = ref(0);
const selectedMenuKeys = ref<number[]>([]);
const halfCheckedMenuKeys = ref<number[]>([]);
const menuTreeData = ref<any[]>([]);
const menuList = ref<any[]>([]);
const roleModalVisible = ref(false);
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
  company: '',
  domain: '',
  expire_time: null as any,
  max_users: 0,
  remark: '',
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '用户名', dataIndex: 'username', width: 120 },
  { title: '昵称', dataIndex: 'nickname', width: 120 },
  { title: '邮箱', dataIndex: 'email', width: 180 },
  { title: '手机', dataIndex: 'phone', width: 130 },
  { title: '部门', dataIndex: 'department_id', width: 100 },
  { title: '角色', dataIndex: 'roles', width: 150 },
  { title: '管理员', dataIndex: 'is_admin', width: 80 },
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
    department_id: record.department_id,
    status: record.status,
    company: record.company,
    domain: record.domain,
    expire_time: record.expire_time ? dayjs(record.expire_time) : null,
    max_users: record.max_users,
    remark: record.remark,
  });
  modalVisible.value = true;
}

async function handleSubmit() {
  if (formState.id) {
    await updateUser(formState.id, {
      nickname: formState.nickname,
      email: formState.email,
      phone: formState.phone,
      department_id: formState.department_id,
      status: formState.status,
      company: formState.company,
      domain: formState.domain,
      expire_time: formState.expire_time?.format('YYYY-MM-DD HH:mm:ss') || null,
      max_users: formState.max_users,
      remark: formState.remark,
    });
    message.success('更新成功');
  } else {
    await createUser({
      username: formState.username,
      password: formState.password,
      nickname: formState.nickname,
      email: formState.email,
      phone: formState.phone,
      department_id: formState.department_id,
      status: formState.status,
      company: formState.company,
      domain: formState.domain,
      expire_time: formState.expire_time?.format('YYYY-MM-DD HH:mm:ss') || null,
      max_users: formState.max_users,
      remark: formState.remark,
    });
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

async function handleUnlock(id: number) {
  await unlockUser(id);
  message.success('用户已解锁');
  fetchData();
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

function buildMenuTree(menus: any[]): any[] {
  return menus.map((m) => ({
    key: m.id,
    title: m.title || m.name,
    children: m.children ? buildMenuTree(m.children) : [],
  }));
}

// 收集树中所有非叶子节点的 key
function collectParentKeys(tree: any[]): Set<number> {
  const parentKeys = new Set<number>();
  function walk(nodes: any[]) {
    for (const node of nodes) {
      if (node.children && node.children.length > 0) {
        parentKeys.add(node.key);
        walk(node.children);
      }
    }
  }
  walk(tree);
  return parentKeys;
}

async function fetchMenuTree() {
  const menus = await getMenuList();
  menuList.value = menus;
  menuTreeData.value = buildMenuTree(menus);
}

function handleAssignMenus(record: UserRecord) {
  currentUserId.value = record.id;
  fetchMenuTree();
  // TODO: 加载用户已分配的菜单
  selectedMenuKeys.value = [];
  halfCheckedMenuKeys.value = [];
  menuModalVisible.value = true;
}

function handleMenuCheck(checked: any, info: any) {
  if (Array.isArray(checked)) {
    selectedMenuKeys.value = checked as number[];
    halfCheckedMenuKeys.value = (info?.halfCheckedKeys || []) as number[];
  } else {
    selectedMenuKeys.value = (checked.checked || []) as number[];
    halfCheckedMenuKeys.value = (checked.halfChecked || []) as number[];
  }
}

async function handleSaveMenus() {
  const allMenuIds = [...new Set([...selectedMenuKeys.value, ...halfCheckedMenuKeys.value])];
  await assignUserMenus(currentUserId.value, allMenuIds);
  message.success('菜单分配成功');
  menuModalVisible.value = false;
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
          <template v-if="column.dataIndex === 'is_admin'">
            <Tag :color="(record as UserRecord).is_admin === 2 ? 'red' : (record as UserRecord).is_admin === 1 ? 'blue' : 'default'">
              {{ (record as UserRecord).is_admin === 2 ? '超管' : (record as UserRecord).is_admin === 1 ? '管理员' : '普通' }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'totp_enabled'">
            <Tag :color="(record as UserRecord).totp_enabled === 1 ? 'green' : 'default'">
              {{ (record as UserRecord).totp_enabled === 1 ? '已启用' : '未启用' }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'status'">
            <Tag v-if="(record as UserRecord).status === 2" color="orange">已锁定</Tag>
            <Switch
              v-else
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
              <Button v-if="hasAccessByCodes(['user:assign-menus'])" size="small" type="link" @click="handleAssignMenus(record as UserRecord)">分配菜单</Button>
              <Popconfirm v-if="hasAccessByCodes(['user:reset-password'])" title="确认重置密码？" @confirm="handleResetPassword((record as UserRecord).id)">
                <Button size="small" type="link">重置密码</Button>
              </Popconfirm>
              <Popconfirm v-if="(record as UserRecord).status === 2 && hasAccessByCodes(['user:unlock'])" title="确认解锁该用户？" @confirm="handleUnlock((record as UserRecord).id)">
                <Button size="small" type="link" style="color: orange">解锁</Button>
              </Popconfirm>
              <Popconfirm v-if="hasAccessByCodes(['user:delete'])" title="确认删除？" @confirm="handleDelete((record as UserRecord).id)">
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
          <FormItem label="公司名称">
            <Input v-model:value="formState.company" placeholder="请输入公司名称" />
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

    <Modal v-model:open="menuModalVisible" title="分配菜单" @ok="handleSaveMenus" :width="480">
      <div class="mt-4" style="max-height: 400px; overflow-y: auto;">
        <Tree
          v-model:checkedKeys="selectedMenuKeys"
          :tree-data="menuTreeData"
          checkable
          default-expand-all
          @check="handleMenuCheck"
        />
      </div>
    </Modal>
  </div>
</template>
