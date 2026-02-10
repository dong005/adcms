<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Tag, Space, Modal, Form, FormItem, Input,
  Popconfirm, message, Card,
} from 'ant-design-vue';
import { getTagList, createTag, updateTag, deleteTag } from '#/api/cms';
import type { TagRecord } from '#/api/cms';

const loading = ref(false);
const dataSource = ref<TagRecord[]>([]);
const total = ref(0);
const keyword = ref('');
const pagination = reactive({ current: 1, pageSize: 10 });

const modalVisible = ref(false);
const modalTitle = ref('新增标签');

const formState = reactive({
  id: 0,
  name: '',
  slug: '',
  description: '',
  color: '#1890ff',
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '名称', dataIndex: 'name', width: 150 },
  { title: '别名', dataIndex: 'slug', width: 150 },
  { title: '颜色', dataIndex: 'color', width: 100 },
  { title: '描述', dataIndex: 'description', width: 250 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', dataIndex: 'action', width: 150, fixed: 'right' as const },
];

async function fetchData() {
  loading.value = true;
  try {
    const res = await getTagList({
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
  modalTitle.value = '新增标签';
  Object.assign(formState, { 
    id: 0, 
    name: '', 
    slug: '', 
    description: '', 
    color: '#1890ff' 
  });
  modalVisible.value = true;
}

function handleEdit(record: TagRecord) {
  modalTitle.value = '编辑标签';
  Object.assign(formState, {
    id: record.id,
    name: record.name,
    slug: record.slug,
    description: record.description,
    color: record.color,
  });
  modalVisible.value = true;
}

async function handleSubmit() {
  try {
    const submitData = { ...formState };
    if (formState.id) {
      await updateTag(formState.id, submitData);
      message.success('更新成功');
    } else {
      await createTag(submitData);
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
    await deleteTag(id);
    message.success('删除成功');
    fetchData();
  } catch (error) {
    message.error('删除失败');
  }
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
  <Card>
    <Space class="mb-4">
      <Input v-model:value="keyword" placeholder="搜索标签" style="width: 200px" />
      <Button type="primary" @click="fetchData">搜索</Button>
      <Button @click="handleAdd">新增标签</Button>
    </Space>

    <Table
      :columns="columns"
      :data-source="dataSource"
      :loading="loading"
      :pagination="{ ...pagination, total }"
      row-key="id"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'color'">
          <Tag :color="record.color">{{ record.color }}</Tag>
        </template>
        <template v-if="column.dataIndex === 'created_at'">
          {{ new Date(record.created_at).toLocaleString() }}
        </template>
        <template v-if="column.dataIndex === 'action'">
          <Space>
            <Button size="small" @click="handleEdit(record as TagRecord)">编辑</Button>
            <Popconfirm title="确定删除这个标签吗？" @confirm="handleDelete(record.id)">
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
        <FormItem label="名称" name="name" :rules="[{ required: true, message: '请输入标签名称' }]">
          <Input v-model:value="formState.name" />
        </FormItem>
        <FormItem label="别名" name="slug">
          <Input v-model:value="formState.slug" />
        </FormItem>
        <FormItem label="颜色" name="color">
          <Input v-model:value="formState.color" type="color" />
        </FormItem>
        <FormItem label="描述" name="description">
          <Input.TextArea v-model:value="formState.description" :rows="3" />
        </FormItem>
      </Form>
    </Modal>
  </Card>
</template>
