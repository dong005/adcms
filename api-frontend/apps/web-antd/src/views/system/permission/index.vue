<script lang="ts" setup>
import { ref, reactive, onMounted, computed } from 'vue';
import {
  Table, Tag, Card, Button, Space, Modal, Form, FormItem,
  Input, Select, SelectOption, TreeSelect, Popconfirm, message,
} from 'ant-design-vue';
import {
  getPermissionTree, createPermission, updatePermission, deletePermission,
} from '#/api/system/permission';
import { useAccess } from '@vben/access';
import type { PermissionRecord } from '#/api/system/permission';

const { hasAccessByCodes } = useAccess();

const loading = ref(false);
const dataSource = ref<PermissionRecord[]>([]);
const modalVisible = ref(false);
const modalTitle = ref('新增权限');

const typeMap: Record<number, { label: string; color: string }> = {
  1: { label: '模块', color: 'blue' },
  2: { label: '按钮', color: 'orange' },
  3: { label: 'API', color: 'green' },
};

const methodColors: Record<string, string> = {
  GET: 'blue',
  POST: 'green',
  PUT: 'orange',
  DELETE: 'red',
};

const formState = reactive({
  id: 0,
  name: '',
  code: '',
  type: 3 as number,
  parent_id: 0 as number,
  path: '',
  method: '',
  description: '',
});

const columns = [
  { title: '权限名称', dataIndex: 'name', width: 200 },
  { title: '权限编码', dataIndex: 'code', width: 200 },
  { title: '类型', dataIndex: 'type', width: 80 },
  { title: '请求方法', dataIndex: 'method', width: 100 },
  { title: '请求路径', dataIndex: 'path', width: 250 },
  { title: '描述', dataIndex: 'description', width: 200 },
  { title: '操作', dataIndex: 'action', width: 200, fixed: 'right' as const },
];

function buildTreeSelectData(tree: PermissionRecord[]): any[] {
  return tree.map((item) => ({
    value: item.id,
    title: item.name,
    children: item.children ? buildTreeSelectData(item.children) : [],
  }));
}

const parentTreeData = computed(() => {
  return [{ value: 0, title: '顶级（无父级）', children: buildTreeSelectData(dataSource.value) }];
});

async function fetchData() {
  loading.value = true;
  try {
    const tree = await getPermissionTree();
    dataSource.value = tree || [];
  } finally {
    loading.value = false;
  }
}

function handleAdd(parentId = 0) {
  modalTitle.value = '新增权限';
  Object.assign(formState, {
    id: 0, name: '', code: '', type: 3, parent_id: parentId,
    path: '', method: '', description: '',
  });
  modalVisible.value = true;
}

function handleEdit(record: PermissionRecord) {
  modalTitle.value = '编辑权限';
  Object.assign(formState, {
    id: record.id, name: record.name, code: record.code,
    type: record.type, parent_id: record.parent_id,
    path: record.path, method: record.method,
    description: record.description,
  });
  modalVisible.value = true;
}

async function handleSubmit() {
  try {
    if (formState.id) {
      await updatePermission(formState.id, { ...formState });
      message.success('更新成功');
    } else {
      await createPermission({ ...formState });
      message.success('创建成功');
    }
    modalVisible.value = false;
    fetchData();
  } catch {
    message.error('操作失败');
  }
}

async function handleDelete(id: number) {
  try {
    await deletePermission(id);
    message.success('删除成功');
    fetchData();
  } catch {
    message.error('删除失败，可能存在子权限');
  }
}

onMounted(() => {
  fetchData();
});
</script>

<template>
  <div class="p-4">
    <Card title="权限管理">
      <template #extra>
        <Button v-if="hasAccessByCodes(['permission:list'])" type="primary" @click="handleAdd(0)">新增权限</Button>
      </template>

      <Table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
        :pagination="false"
        :scroll="{ x: 1200 }"
        row-key="id"
        size="middle"
        default-expand-all-rows
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'type'">
            <Tag :color="typeMap[(record as PermissionRecord).type]?.color || 'default'">
              {{ typeMap[(record as PermissionRecord).type]?.label || (record as PermissionRecord).type }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'method'">
            <Tag v-if="(record as PermissionRecord).method" :color="methodColors[(record as PermissionRecord).method] || 'default'">
              {{ (record as PermissionRecord).method }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'path'">
            <code v-if="(record as PermissionRecord).path" style="font-size: 12px; opacity: 0.7;">
              {{ (record as PermissionRecord).path }}
            </code>
          </template>
          <template v-if="column.dataIndex === 'action'">
            <Space>
              <Button size="small" type="link" @click="handleAdd((record as PermissionRecord).id)">添加子级</Button>
              <Button size="small" type="link" @click="handleEdit(record as PermissionRecord)">编辑</Button>
              <Popconfirm title="确认删除？" @confirm="handleDelete((record as PermissionRecord).id)">
                <Button size="small" type="link" danger>删除</Button>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <Modal v-model:open="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="560">
      <Form :model="formState" layout="vertical" class="mt-4">
        <FormItem label="父级权限">
          <TreeSelect
            v-model:value="formState.parent_id"
            :tree-data="parentTreeData"
            placeholder="选择父级"
            tree-default-expand-all
            style="width: 100%"
          />
        </FormItem>
        <FormItem label="权限名称">
          <Input v-model:value="formState.name" placeholder="如：用户列表" />
        </FormItem>
        <FormItem label="权限编码">
          <Input v-model:value="formState.code" placeholder="如：user:list" />
        </FormItem>
        <FormItem label="类型">
          <Select v-model:value="formState.type">
            <SelectOption :value="1">模块</SelectOption>
            <SelectOption :value="2">按钮</SelectOption>
            <SelectOption :value="3">API</SelectOption>
          </Select>
        </FormItem>
        <FormItem label="请求方法" v-if="formState.type === 3">
          <Select v-model:value="formState.method" allow-clear placeholder="选择请求方法">
            <SelectOption value="GET">GET</SelectOption>
            <SelectOption value="POST">POST</SelectOption>
            <SelectOption value="PUT">PUT</SelectOption>
            <SelectOption value="DELETE">DELETE</SelectOption>
          </Select>
        </FormItem>
        <FormItem label="请求路径" v-if="formState.type === 3">
          <Input v-model:value="formState.path" placeholder="如：/api/users/:id" />
        </FormItem>
        <FormItem label="描述">
          <Input.TextArea v-model:value="formState.description" placeholder="权限描述" :rows="2" />
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>
