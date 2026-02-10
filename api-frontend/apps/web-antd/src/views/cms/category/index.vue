<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Tag, Space, Modal, Form, FormItem, Input, InputNumber,
  Popconfirm, message, Switch, Card,
} from 'ant-design-vue';
import { getCategoryList, createCategory, updateCategory, deleteCategory } from '#/api/cms';
import type { CategoryRecord } from '#/api/cms';

const loading = ref(false);
const dataSource = ref<CategoryRecord[]>([]);
const modalVisible = ref(false);
const modalTitle = ref('新增分类');
const expandedKeys = ref<number[]>([]);

const formState = reactive({
  id: 0,
  name: '',
  slug: '',
  description: '',
  parent_id: 0,
  sort: 0,
  status: 1,
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '名称', dataIndex: 'name', width: 200 },
  { title: '别名', dataIndex: 'slug', width: 150 },
  { title: '描述', dataIndex: 'description', width: 250 },
  { title: '排序', dataIndex: 'sort', width: 80 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', dataIndex: 'action', width: 180, fixed: 'right' as const },
];

function buildParentTree(categories: CategoryRecord[], excludeId?: number): any[] {
  return categories
    .filter(cat => cat.id !== excludeId)
    .map(cat => ({
      value: cat.id,
      title: cat.name,
      children: cat.children ? buildParentTree(cat.children, excludeId) : [],
    }));
}

async function fetchData() {
  loading.value = true;
  try {
    dataSource.value = await getCategoryList();
  } finally {
    loading.value = false;
  }
}

function handleAdd() {
  modalTitle.value = '新增分类';
  Object.assign(formState, { 
    id: 0, 
    name: '', 
    slug: '', 
    description: '', 
    parent_id: 0, 
    sort: 0, 
    status: 1 
  });
  modalVisible.value = true;
}

function handleEdit(record: CategoryRecord) {
  modalTitle.value = '编辑分类';
  Object.assign(formState, {
    id: record.id,
    name: record.name,
    slug: record.slug,
    description: record.description,
    parent_id: record.parent_id,
    sort: record.sort,
    status: record.status,
  });
  modalVisible.value = true;
}

async function handleSubmit() {
  try {
    const submitData = { ...formState };
    if (formState.id) {
      await updateCategory(formState.id, submitData);
      message.success('更新成功');
    } else {
      await createCategory(submitData);
      message.success('创建成功');
    }
    modalVisible.value = false;
    fetchData();
  } catch (error) {
    message.error('操作失败');
  }
}

async function handleDelete(id: number) {
  try {
    await deleteCategory(id);
    message.success('删除成功');
    fetchData();
  } catch (error) {
    message.error('删除失败');
  }
}

onMounted(() => {
  fetchData();
});
</script>

<template>
  <Card>
    <Space class="mb-4">
      <Button type="primary" @click="handleAdd">新增分类</Button>
    </Space>

    <Table
      :columns="columns"
      :data-source="dataSource"
      :loading="loading"
      :pagination="false"
      :default-expand-all="true"
      row-key="id"
      :expanded-row-keys="expandedKeys"
      @expand="(expanded, record) => {
        if (expanded) {
          expandedKeys.push(record.id);
        } else {
          const index = expandedKeys.indexOf(record.id);
          if (index > -1) {
            expandedKeys.splice(index, 1);
          }
        }
      }"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'name'">
          <span :style="{ paddingLeft: `${(record.level || 0) * 20}px` }">
            {{ record.name }}
          </span>
        </template>
        <template v-if="column.dataIndex === 'status'">
          <Tag :color="record.status === 1 ? 'green' : 'red'">
            {{ record.status === 1 ? '启用' : '禁用' }}
          </Tag>
        </template>
        <template v-if="column.dataIndex === 'created_at'">
          {{ new Date(record.created_at).toLocaleString() }}
        </template>
        <template v-if="column.dataIndex === 'action'">
          <Space>
            <Button size="small" @click="handleEdit(record)">编辑</Button>
            <Popconfirm title="确定删除这个分类吗？" @confirm="handleDelete(record.id)">
              <Button size="small" danger>删除</Button>
            </Popconfirm>
          </Space>
        </template>
      </template>
    </Table>

    <Modal
      v-model:open="modalVisible"
      :title="modalTitle"
      @ok="handleSubmit"
    >
      <Form ref="formRef" :model="formState" layout="vertical">
        <FormItem label="名称" name="name" :rules="[{ required: true, message: '请输入分类名称' }]">
          <Input v-model:value="formState.name" />
        </FormItem>
        <FormItem label="别名" name="slug">
          <Input v-model:value="formState.slug" />
        </FormItem>
        <FormItem label="父分类" name="parent_id">
          <TreeSelect
            v-model:value="formState.parent_id"
            :tree-data="buildParentTree(dataSource, formState.id || undefined)"
            placeholder="选择父分类"
            allow-clear
            tree-default-expand-all
          />
        </FormItem>
        <FormItem label="描述" name="description">
          <Input.TextArea v-model:value="formState.description" :rows="3" />
        </FormItem>
        <FormItem label="排序" name="sort">
          <InputNumber v-model:value="formState.sort" :min="0" style="width: 100%" />
        </FormItem>
        <FormItem label="状态" name="status">
          <Switch
            :checked="formState.status === 1"
            @change="(checked) => formState.status = checked ? 1 : 0"
          />
        </FormItem>
      </Form>
    </Modal>
  </Card>
</Template>
