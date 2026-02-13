<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Tag, Space, Modal, Form, FormItem, Input, InputNumber,
  Select, SelectOption, Popconfirm, message, Card,
} from 'ant-design-vue';
import { getLinkList, createLink, updateLink, deleteLink } from '#/api/system/link';
import type { LinkRecord } from '#/api/system/link';

const loading = ref(false);
const dataSource = ref<LinkRecord[]>([]);
const keyword = ref('');

const modalVisible = ref(false);
const modalTitle = ref('新增友链');
const formState = reactive({
  id: 0,
  name: '',
  url: '',
  logo: '',
  desc: '',
  sort: 0,
  status: 1,
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '名称', dataIndex: 'name', width: 150 },
  { title: 'URL', dataIndex: 'url', width: 250 },
  { title: 'Logo', dataIndex: 'logo', width: 150 },
  { title: '描述', dataIndex: 'desc', width: 200 },
  { title: '排序', dataIndex: 'sort', width: 70 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '操作', dataIndex: 'action', width: 150, fixed: 'right' as const },
];

async function fetchData() {
  loading.value = true;
  try {
    dataSource.value = await getLinkList(keyword.value) || [];
  } finally {
    loading.value = false;
  }
}

function handleAdd() {
  modalTitle.value = '新增友链';
  Object.assign(formState, { id: 0, name: '', url: '', logo: '', desc: '', sort: 0, status: 1 });
  modalVisible.value = true;
}

function handleEdit(record: LinkRecord) {
  modalTitle.value = '编辑友链';
  Object.assign(formState, {
    id: record.id, name: record.name, url: record.url,
    logo: record.logo, desc: record.desc, sort: record.sort, status: record.status,
  });
  modalVisible.value = true;
}

async function handleSubmit() {
  if (formState.id) {
    await updateLink(formState.id, { ...formState });
    message.success('更新成功');
  } else {
    await createLink({ ...formState });
    message.success('创建成功');
  }
  modalVisible.value = false;
  fetchData();
}

async function handleDelete(id: number) {
  await deleteLink(id);
  message.success('删除成功');
  fetchData();
}

onMounted(fetchData);
</script>

<template>
  <div class="p-4">
    <Card title="友链管理">
      <template #extra>
        <Space>
          <Input.Search v-model:value="keyword" placeholder="搜索友链" @search="fetchData" style="width: 200px" />
          <Button type="primary" @click="handleAdd">新增友链</Button>
        </Space>
      </template>
      <Table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
        :pagination="false"
        row-key="id"
        size="middle"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'url'">
            <a :href="(record as LinkRecord).url" target="_blank">{{ (record as LinkRecord).url }}</a>
          </template>
          <template v-if="column.dataIndex === 'status'">
            <Tag :color="(record as LinkRecord).status === 1 ? 'green' : 'red'">
              {{ (record as LinkRecord).status === 1 ? '启用' : '禁用' }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'action'">
            <Space>
              <Button size="small" type="link" @click="handleEdit(record as LinkRecord)">编辑</Button>
              <Popconfirm title="确认删除？" @confirm="handleDelete((record as LinkRecord).id)">
                <Button size="small" type="link" danger>删除</Button>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <Modal v-model:open="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="520">
      <Form :model="formState" layout="vertical" class="mt-4">
        <FormItem label="名称" :rules="[{ required: true }]">
          <Input v-model:value="formState.name" placeholder="请输入友链名称" />
        </FormItem>
        <FormItem label="URL" :rules="[{ required: true }]">
          <Input v-model:value="formState.url" placeholder="请输入友链URL" />
        </FormItem>
        <FormItem label="Logo">
          <Input v-model:value="formState.logo" placeholder="Logo图片URL" />
        </FormItem>
        <FormItem label="描述">
          <Input.TextArea v-model:value="formState.desc" placeholder="请输入描述" :rows="2" />
        </FormItem>
        <div class="grid grid-cols-2 gap-x-4">
          <FormItem label="排序">
            <InputNumber v-model:value="formState.sort" :min="0" style="width: 100%" />
          </FormItem>
          <FormItem label="状态">
            <Select v-model:value="formState.status">
              <SelectOption :value="1">启用</SelectOption>
              <SelectOption :value="0">禁用</SelectOption>
            </Select>
          </FormItem>
        </div>
      </Form>
    </Modal>
  </div>
</template>
