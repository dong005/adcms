<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Input, Tag, Space, Modal, Form, FormItem,
  Select, SelectOption, Popconfirm, message, Card, TreeSelect,
  InputNumber,
} from 'ant-design-vue';
import {
  getDepartmentTree, createDepartment,
  updateDepartment, deleteDepartment,
} from '#/api/system/department';
import type { DepartmentRecord } from '#/api/system/department';

const loading = ref(false);
const dataSource = ref<DepartmentRecord[]>([]);
const treeData = ref<DepartmentRecord[]>([]);

const modalVisible = ref(false);
const modalTitle = ref('新增部门');

const formState = reactive({
  id: 0,
  parent_id: 0 as number | undefined,
  name: '',
  code: '',
  leader: '',
  phone: '',
  email: '',
  sort: 0,
  status: 1,
});

const columns = [
  { title: '部门名称', dataIndex: 'name', width: 200 },
  { title: '编码', dataIndex: 'code', width: 120 },
  { title: '负责人', dataIndex: 'leader', width: 100 },
  { title: '联系电话', dataIndex: 'phone', width: 130 },
  { title: '邮箱', dataIndex: 'email', width: 180 },
  { title: '排序', dataIndex: 'sort', width: 70 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '操作', dataIndex: 'action', width: 200, fixed: 'right' as const },
];

async function fetchData() {
  loading.value = true;
  try {
    const tree = await getDepartmentTree();
    dataSource.value = tree || [];
    treeData.value = tree || [];
  } finally {
    loading.value = false;
  }
}

function handleAdd(parentId = 0) {
  modalTitle.value = '新增部门';
  Object.assign(formState, {
    id: 0, parent_id: parentId || undefined, name: '', code: '',
    leader: '', phone: '', email: '', sort: 0, status: 1,
  });
  modalVisible.value = true;
}

function handleEdit(record: DepartmentRecord) {
  modalTitle.value = '编辑部门';
  Object.assign(formState, {
    id: record.id,
    parent_id: record.parent_id || undefined,
    name: record.name,
    code: record.code,
    leader: record.leader,
    phone: record.phone,
    email: record.email,
    sort: record.sort,
    status: record.status,
  });
  modalVisible.value = true;
}

async function handleSubmit() {
  const data = { ...formState, parent_id: formState.parent_id || 0 };
  if (formState.id) {
    await updateDepartment(formState.id, data);
    message.success('更新成功');
  } else {
    await createDepartment(data);
    message.success('创建成功');
  }
  modalVisible.value = false;
  fetchData();
}

async function handleDelete(id: number) {
  await deleteDepartment(id);
  message.success('删除成功');
  fetchData();
}

function buildTreeSelectData(depts: DepartmentRecord[]): any[] {
  return depts.map((d) => ({
    value: d.id,
    title: d.name,
    children: d.children ? buildTreeSelectData(d.children) : [],
  }));
}

onMounted(() => {
  fetchData();
});
</script>

<template>
  <div class="p-4">
    <Card title="部门管理">
      <template #extra>
        <Button type="primary" @click="handleAdd(0)">新增部门</Button>
      </template>

      <Table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
        :pagination="false"
        :scroll="{ x: 1100 }"
        row-key="id"
        size="middle"
        default-expand-all-rows
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'status'">
            <Tag :color="(record as DepartmentRecord).status === 1 ? 'green' : 'red'">
              {{ (record as DepartmentRecord).status === 1 ? '启用' : '禁用' }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'action'">
            <Space>
              <Button size="small" type="link" @click="handleAdd((record as DepartmentRecord).id)">新增子部门</Button>
              <Button size="small" type="link" @click="handleEdit(record as DepartmentRecord)">编辑</Button>
              <Popconfirm title="确认删除？" @confirm="handleDelete((record as DepartmentRecord).id)">
                <Button size="small" type="link" danger>删除</Button>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <Modal v-model:open="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="560">
      <Form :model="formState" layout="vertical" class="mt-4">
        <FormItem label="上级部门">
          <TreeSelect
            v-model:value="formState.parent_id"
            :tree-data="buildTreeSelectData(treeData)"
            placeholder="请选择上级部门（留空为顶级）"
            allow-clear
            tree-default-expand-all
          />
        </FormItem>
        <FormItem label="部门名称" :rules="[{ required: true }]">
          <Input v-model:value="formState.name" placeholder="请输入部门名称" />
        </FormItem>
        <FormItem label="部门编码">
          <Input v-model:value="formState.code" placeholder="请输入部门编码" />
        </FormItem>
        <FormItem label="负责人">
          <Input v-model:value="formState.leader" placeholder="请输入负责人" />
        </FormItem>
        <FormItem label="联系电话">
          <Input v-model:value="formState.phone" placeholder="请输入联系电话" />
        </FormItem>
        <FormItem label="邮箱">
          <Input v-model:value="formState.email" placeholder="请输入邮箱" />
        </FormItem>
        <FormItem label="排序">
          <InputNumber v-model:value="formState.sort" :min="0" style="width: 100%" />
        </FormItem>
        <FormItem label="状态">
          <Select v-model:value="formState.status">
            <SelectOption :value="1">启用</SelectOption>
            <SelectOption :value="0">禁用</SelectOption>
          </Select>
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>
