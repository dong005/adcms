<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Space, Modal, Form, FormItem, Input, InputNumber,
  Popconfirm, message, Card,
} from 'ant-design-vue';
import { getConfigWebList, saveConfigWebs, deleteConfigWeb } from '#/api/system/config';
import type { ConfigWebRecord } from '#/api/system/config';

const loading = ref(false);
const dataSource = ref<ConfigWebRecord[]>([]);

const modalVisible = ref(false);
const modalTitle = ref('新增网站设置');
const formState = reactive({ id: 0, name: '', code: '', value: '', sort: 0 });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '名称', dataIndex: 'name', width: 150 },
  { title: '编码', dataIndex: 'code', width: 150 },
  { title: '值', dataIndex: 'value', width: 300 },
  { title: '排序', dataIndex: 'sort', width: 70 },
  { title: '操作', dataIndex: 'action', width: 150, fixed: 'right' as const },
];

async function fetchData() {
  loading.value = true;
  try {
    dataSource.value = await getConfigWebList() || [];
  } finally {
    loading.value = false;
  }
}

function handleAdd() {
  modalTitle.value = '新增网站设置';
  Object.assign(formState, { id: 0, name: '', code: '', value: '', sort: 0 });
  modalVisible.value = true;
}

function handleEdit(record: ConfigWebRecord) {
  modalTitle.value = '编辑网站设置';
  Object.assign(formState, {
    id: record.id, name: record.name, code: record.code,
    value: record.value, sort: record.sort,
  });
  modalVisible.value = true;
}

async function handleSubmit() {
  await saveConfigWebs([{ ...formState }]);
  message.success('保存成功');
  modalVisible.value = false;
  fetchData();
}

async function handleDelete(id: number) {
  await deleteConfigWeb(id);
  message.success('删除成功');
  fetchData();
}

onMounted(fetchData);
</script>

<template>
  <div class="p-4">
    <Card title="网站设置">
      <template #extra>
        <Button type="primary" @click="handleAdd">新增设置</Button>
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
          <template v-if="column.dataIndex === 'value'">
            <span class="truncate block max-w-[300px]">{{ (record as ConfigWebRecord).value }}</span>
          </template>
          <template v-if="column.dataIndex === 'action'">
            <Space>
              <Button size="small" type="link" @click="handleEdit(record as ConfigWebRecord)">编辑</Button>
              <Popconfirm title="确认删除？" @confirm="handleDelete((record as ConfigWebRecord).id)">
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
          <Input v-model:value="formState.name" placeholder="如 网站名称" />
        </FormItem>
        <FormItem label="编码" :rules="[{ required: true }]">
          <Input v-model:value="formState.code" placeholder="如 site_name" :disabled="!!formState.id" />
        </FormItem>
        <FormItem label="值">
          <Input.TextArea v-model:value="formState.value" placeholder="请输入值" :rows="4" />
        </FormItem>
        <FormItem label="排序">
          <InputNumber v-model:value="formState.sort" :min="0" style="width: 100%" />
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>
