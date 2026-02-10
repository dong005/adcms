<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Tag, Space, Modal, Form, FormItem, Input, InputNumber,
  Select, SelectOption, Popconfirm, message, Tree, Card,
} from 'ant-design-vue';
import { getRoleList, createRole, updateRole, deleteRole, getRoleMenus, assignRoleMenus } from '#/api/system/role';
import { getMenuList } from '#/api/system/menu';
import { getPermissionTree, getRolePermissions, assignRolePermissions } from '#/api/system/permission';
import { useAccess } from '@vben/access';
import type { RoleRecord } from '#/api/system/role';
import type { MenuRecord } from '#/api/system/menu';
import type { PermissionRecord } from '#/api/system/permission';

const { hasAccessByCodes } = useAccess();

const loading = ref(false);
const dataSource = ref<RoleRecord[]>([]);

const modalVisible = ref(false);
const modalTitle = ref('新增角色');
const menuModalVisible = ref(false);
const currentRoleId = ref(0);
const checkedMenuKeys = ref<number[]>([]);
const halfCheckedMenuKeys = ref<number[]>([]);
const menuTreeData = ref<any[]>([]);
const permModalVisible = ref(false);
const checkedPermKeys = ref<number[]>([]);
const halfCheckedPermKeys = ref<number[]>([]);
const permTreeData = ref<any[]>([]);

const formState = reactive({
  id: 0,
  name: '',
  code: '',
  description: '',
  status: 1,
  sort: 0,
  data_scope: 1,
});

const dataScopeOptions = [
  { value: 1, label: '全部数据' },
  { value: 2, label: '本部门及下级' },
  { value: 3, label: '仅本部门' },
  { value: 4, label: '仅自己' },
];

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '角色名称', dataIndex: 'name', width: 150 },
  { title: '角色编码', dataIndex: 'code', width: 150 },
  { title: '描述', dataIndex: 'description', width: 200 },
  { title: '数据权限', dataIndex: 'data_scope', width: 120 },
  { title: '排序', dataIndex: 'sort', width: 80 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', dataIndex: 'action', width: 250, fixed: 'right' as const },
];

async function fetchData() {
  loading.value = true;
  try {
    dataSource.value = await getRoleList();
  } finally {
    loading.value = false;
  }
}

function buildMenuTree(menus: MenuRecord[]): any[] {
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
  menuTreeData.value = buildMenuTree(menus);
}

const permTypeMap: Record<number, { label: string; color: string }> = {
  1: { label: '菜单', color: 'blue' },
  2: { label: '按钮', color: 'orange' },
  3: { label: 'API', color: 'green' },
};

function buildPermTree(perms: PermissionRecord[]): any[] {
  return perms.map((p) => {
    const info = permTypeMap[p.type] || { label: '未知', color: 'default' };
    return {
      key: p.id,
      title: `[${info.label}] ${p.name} (${p.code})`,
      children: p.children ? buildPermTree(p.children) : [],
    };
  });
}

async function fetchPermTree() {
  const tree = await getPermissionTree();
  permTreeData.value = buildPermTree(tree || []);
}

function handleAdd() {
  modalTitle.value = '新增角色';
  Object.assign(formState, { id: 0, name: '', code: '', description: '', status: 1, sort: 0, data_scope: 1 });
  modalVisible.value = true;
}

function handleEdit(record: RoleRecord) {
  modalTitle.value = '编辑角色';
  Object.assign(formState, {
    id: record.id, name: record.name, code: record.code,
    description: record.description, status: record.status, sort: record.sort,
    data_scope: (record as any).data_scope || 1,
  });
  modalVisible.value = true;
}

async function handleSubmit() {
  if (formState.id) {
    await updateRole(formState.id, { ...formState });
    message.success('更新成功');
  } else {
    await createRole({ ...formState });
    message.success('创建成功');
  }
  modalVisible.value = false;
  fetchData();
}

async function handleDelete(id: number) {
  await deleteRole(id);
  message.success('删除成功');
  fetchData();
}

async function handleAssignMenus(record: RoleRecord) {
  currentRoleId.value = record.id;
  await fetchMenuTree();
  try {
    const menus = await getRoleMenus(record.id);
    const allIds = (menus || []).map((m: any) => m.id);
    // 回显时只设置叶子节点为 checked，父节点会通过父子联动自动变为半选/全选
    const parentKeys = collectParentKeys(menuTreeData.value);
    checkedMenuKeys.value = allIds.filter((id: number) => !parentKeys.has(id));
    halfCheckedMenuKeys.value = [];
  } catch {
    checkedMenuKeys.value = [];
    halfCheckedMenuKeys.value = [];
  }
  menuModalVisible.value = true;
}

function handleMenuCheck(checked: any, info: any) {
  if (Array.isArray(checked)) {
    checkedMenuKeys.value = checked as number[];
    halfCheckedMenuKeys.value = (info?.halfCheckedKeys || []) as number[];
  } else {
    checkedMenuKeys.value = (checked.checked || []) as number[];
    halfCheckedMenuKeys.value = (checked.halfChecked || []) as number[];
  }
}

async function handleSaveMenus() {
  // 合并全选节点 + 半选父节点，确保父菜单不丢失
  const allMenuIds = [...new Set([...checkedMenuKeys.value, ...halfCheckedMenuKeys.value])];
  await assignRoleMenus(currentRoleId.value, allMenuIds);
  message.success('菜单分配成功');
  menuModalVisible.value = false;
}

async function handleAssignPerms(record: RoleRecord) {
  currentRoleId.value = record.id;
  await fetchPermTree();
  try {
    const perms = await getRolePermissions(record.id);
    const allIds = (perms || []).map((p: any) => p.id);
    const parentKeys = collectParentKeys(permTreeData.value);
    checkedPermKeys.value = allIds.filter((id: number) => !parentKeys.has(id));
    halfCheckedPermKeys.value = [];
  } catch {
    checkedPermKeys.value = [];
    halfCheckedPermKeys.value = [];
  }
  permModalVisible.value = true;
}

function handlePermCheck(checked: any, info: any) {
  if (Array.isArray(checked)) {
    checkedPermKeys.value = checked as number[];
    halfCheckedPermKeys.value = (info?.halfCheckedKeys || []) as number[];
  } else {
    checkedPermKeys.value = (checked.checked || []) as number[];
    halfCheckedPermKeys.value = (checked.halfChecked || []) as number[];
  }
}

async function handleSavePerms() {
  const allPermIds = [...new Set([...checkedPermKeys.value, ...halfCheckedPermKeys.value])];
  await assignRolePermissions(currentRoleId.value, allPermIds);
  message.success('权限分配成功');
  permModalVisible.value = false;
}

function getDataScopeLabel(scope: number) {
  return dataScopeOptions.find(o => o.value === scope)?.label || '未知';
}

onMounted(fetchData);
</script>

<template>
  <div class="p-4">
    <Card title="角色管理">
      <template #extra>
        <Button v-if="hasAccessByCodes(['role:create'])" type="primary" @click="handleAdd">新增角色</Button>
      </template>

      <Table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
        :pagination="false"
        :scroll="{ x: 1000 }"
        row-key="id"
        size="middle"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'status'">
            <Tag :color="(record as RoleRecord).status === 1 ? 'green' : 'red'">
              {{ (record as RoleRecord).status === 1 ? '启用' : '禁用' }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'data_scope'">
            <Tag color="purple">{{ getDataScopeLabel((record as any).data_scope || 1) }}</Tag>
          </template>
          <template v-if="column.dataIndex === 'action'">
            <Space>
              <Button v-if="hasAccessByCodes(['role:update'])" size="small" type="link" @click="handleEdit(record as RoleRecord)">编辑</Button>
              <Button v-if="hasAccessByCodes(['role:assign-menus'])" size="small" type="link" @click="handleAssignMenus(record as RoleRecord)">分配菜单</Button>
              <Button v-if="hasAccessByCodes(['role:assign-permissions'])" size="small" type="link" @click="handleAssignPerms(record as RoleRecord)">分配权限</Button>
              <Popconfirm v-if="hasAccessByCodes(['role:delete'])" title="确认删除？" @confirm="handleDelete((record as RoleRecord).id)">
                <Button size="small" type="link" danger>删除</Button>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <Modal v-model:open="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="520">
      <Form :model="formState" layout="vertical" class="mt-4">
        <FormItem label="角色名称" :rules="[{ required: true }]">
          <Input v-model:value="formState.name" placeholder="请输入角色名称" />
        </FormItem>
        <FormItem label="角色编码" :rules="[{ required: true }]">
          <Input v-model:value="formState.code" :disabled="!!formState.id" placeholder="如 admin, editor" />
        </FormItem>
        <FormItem label="描述">
          <Input.TextArea v-model:value="formState.description" placeholder="请输入描述" :rows="3" />
        </FormItem>
        <FormItem label="排序">
          <InputNumber v-model:value="formState.sort" :min="0" style="width: 100%" />
        </FormItem>
        <FormItem label="数据权限">
          <Select v-model:value="formState.data_scope">
            <SelectOption v-for="opt in dataScopeOptions" :key="opt.value" :value="opt.value">{{ opt.label }}</SelectOption>
          </Select>
        </FormItem>
        <FormItem label="状态">
          <Select v-model:value="formState.status">
            <SelectOption :value="1">启用</SelectOption>
            <SelectOption :value="0">禁用</SelectOption>
          </Select>
        </FormItem>
      </Form>
    </Modal>

    <Modal v-model:open="menuModalVisible" title="分配菜单" @ok="handleSaveMenus" :width="480">
      <div class="mt-4" style="max-height: 400px; overflow-y: auto;">
        <Tree
          v-model:checkedKeys="checkedMenuKeys"
          :tree-data="menuTreeData"
          checkable
          default-expand-all
          @check="handleMenuCheck"
        />
      </div>
    </Modal>

    <Modal v-model:open="permModalVisible" title="分配权限" @ok="handleSavePerms" :width="520">
      <div class="mt-4" style="max-height: 400px; overflow-y: auto;">
        <Tree
          v-model:checkedKeys="checkedPermKeys"
          :tree-data="permTreeData"
          checkable
          default-expand-all
          @check="handlePermCheck"
        />
      </div>
    </Modal>
  </div>
</template>
