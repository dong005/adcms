<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Tag, Space, Modal, Form, FormItem, Input, InputNumber,
  Select, SelectOption, Popconfirm, message, Tree, Card,
} from 'ant-design-vue';
import { getRoleList, createRole, updateRole, deleteRole, getRoleMenus, assignRoleMenus } from '#/api/system/role';
import { getMenuTreeWithButtons } from '#/api/system/menu';
import { useAccess } from '@vben/access';
import type { RoleRecord } from '#/api/system/role';
import type { MenuRecord } from '#/api/system/menu';

const { hasAccessByCodes } = useAccess();

const loading = ref(false);
const dataSource = ref<RoleRecord[]>([]);

const modalVisible = ref(false);
const modalTitle = ref('新增角色');
const assignModalVisible = ref(false);
const assignSaving = ref(false);
const currentRoleId = ref(0);
const currentRoleName = ref('');
const checkedMenuKeys = ref<number[]>([]);
const halfCheckedMenuKeys = ref<number[]>([]);
const menuTreeData = ref<any[]>([]);

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

const menuTypeMap: Record<number, string> = {
  1: '目录',
  2: '菜单',
  3: '页面',
  4: '按钮',
};

function buildMenuTree(menus: MenuRecord[]): any[] {
  return menus.map((m) => {
    const typeLabel = menuTypeMap[m.type] || '';
    const codeLabel = m.permission_code ? ` (${m.permission_code})` : '';
    const suffix = m.type === 4 ? ` [${typeLabel}]${codeLabel}` : '';
    return {
      key: m.id,
      title: `${m.title || m.name}${suffix}`,
      children: m.children ? buildMenuTree(m.children) : [],
    };
  });
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
  const menus = await getMenuTreeWithButtons();
  menuTreeData.value = buildMenuTree(menus);
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

async function handleAssignPermissions(record: RoleRecord) {
  currentRoleId.value = record.id;
  currentRoleName.value = record.name;
  await fetchMenuTree();
  try {
    const menus = await getRoleMenus(record.id);
    const allMenuIds = (menus || []).map((m: any) => m.id);
    const menuParentKeys = collectParentKeys(menuTreeData.value);
    checkedMenuKeys.value = allMenuIds.filter((id: number) => !menuParentKeys.has(id));
    halfCheckedMenuKeys.value = [];
  } catch {
    checkedMenuKeys.value = [];
    halfCheckedMenuKeys.value = [];
  }
  assignModalVisible.value = true;
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

async function handleSaveAll() {
  assignSaving.value = true;
  try {
    const allMenuIds = [...new Set([...checkedMenuKeys.value, ...halfCheckedMenuKeys.value])];
    await assignRoleMenus(currentRoleId.value, allMenuIds);
    message.success('权限分配成功');
    assignModalVisible.value = false;
  } finally {
    assignSaving.value = false;
  }
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
              <Button v-if="hasAccessByCodes(['role:assign-menus']) || hasAccessByCodes(['role:assign-permissions'])" size="small" type="link" @click="handleAssignPermissions(record as RoleRecord)">分配权限</Button>
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

    <Modal v-model:open="assignModalVisible" :title="`分配权限 - ${currentRoleName}`" @ok="handleSaveAll" :confirm-loading="assignSaving" :width="520">
      <div class="mt-4" style="max-height: 450px; overflow-y: auto;">
        <Tree
          v-model:checkedKeys="checkedMenuKeys"
          :tree-data="menuTreeData"
          checkable
          default-expand-all
          @check="handleMenuCheck"
        />
      </div>
    </Modal>
  </div>
</template>
