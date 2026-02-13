<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Tag, Space, Modal, Form, FormItem, Input, InputNumber,
  Select, SelectOption, Popconfirm, message, Card, Switch,
} from 'ant-design-vue';
import { getSiteList, createSite, updateSite, deleteSite } from '#/api/system/site';
import type { SiteRecord } from '#/api/system/site';

const loading = ref(false);
const dataSource = ref<SiteRecord[]>([]);
const keyword = ref('');

const modalVisible = ref(false);
const modalTitle = ref('新增站点');
const formState = reactive({
  id: 0,
  name: '',
  type: '',
  url: '',
  image: '',
  is_domain: 0,
  status: 1,
  sort: 0,
  remark: '',
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '站点名称', dataIndex: 'name', width: 150 },
  { title: '类型', dataIndex: 'type', width: 100 },
  { title: 'URL', dataIndex: 'url', width: 250 },
  { title: '排序', dataIndex: 'sort', width: 70 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '操作', dataIndex: 'action', width: 150, fixed: 'right' as const },
];

async function fetchData() {
  loading.value = true;
  try {
    dataSource.value = await getSiteList(keyword.value) || [];
  } finally {
    loading.value = false;
  }
}

function handleAdd() {
  modalTitle.value = '新增站点';
  Object.assign(formState, {
    id: 0, name: '', type: '', url: '', image: '',
    is_domain: 0, status: 1, sort: 0, remark: '',
  });
  modalVisible.value = true;
}

function handleEdit(record: SiteRecord) {
  modalTitle.value = '编辑站点';
  Object.assign(formState, {
    id: record.id, name: record.name, type: record.type,
    url: record.url, image: record.image, is_domain: record.is_domain,
    status: record.status, sort: record.sort, remark: record.remark,
  });
  modalVisible.value = true;
}

async function handleSubmit() {
  if (formState.id) {
    await updateSite(formState.id, { ...formState });
    message.success('更新成功');
  } else {
    await createSite({ ...formState });
    message.success('创建成功');
  }
  modalVisible.value = false;
  fetchData();
}

async function handleDelete(id: number) {
  await deleteSite(id);
  message.success('删除成功');
  fetchData();
}

onMounted(fetchData);
</script>

<template>
  <div class="p-4">
    <Card title="站点管理">
      <template #extra>
        <Space>
          <Input.Search v-model:value="keyword" placeholder="搜索站点" @search="fetchData" style="width: 200px" />
          <Button type="primary" @click="handleAdd">新增站点</Button>
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
            <a :href="(record as SiteRecord).url" target="_blank">{{ (record as SiteRecord).url }}</a>
          </template>
          <template v-if="column.dataIndex === 'status'">
            <Tag :color="(record as SiteRecord).status === 1 ? 'green' : 'red'">
              {{ (record as SiteRecord).status === 1 ? '启用' : '禁用' }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'action'">
            <Space>
              <Button size="small" type="link" @click="handleEdit(record as SiteRecord)">编辑</Button>
              <Popconfirm title="确认删除？" @confirm="handleDelete((record as SiteRecord).id)">
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
          <FormItem label="站点名称" :rules="[{ required: true }]">
            <Input v-model:value="formState.name" placeholder="请输入站点名称" />
          </FormItem>
          <FormItem label="类型">
            <Input v-model:value="formState.type" placeholder="如 官网、博客" />
          </FormItem>
          <FormItem label="URL">
            <Input v-model:value="formState.url" placeholder="请输入站点URL" />
          </FormItem>
          <FormItem label="图片">
            <Input v-model:value="formState.image" placeholder="站点图片URL" />
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
          <FormItem label="绑定域名">
            <Switch v-model:checked="formState.is_domain" :checked-value="1" :un-checked-value="0" />
          </FormItem>
        </div>
        <FormItem label="备注">
          <Input.TextArea v-model:value="formState.remark" placeholder="请输入备注" :rows="3" />
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>
