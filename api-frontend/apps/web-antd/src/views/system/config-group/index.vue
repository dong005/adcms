<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Space, Modal, Form, FormItem, Input, InputNumber,
  Popconfirm, message, Card,
} from 'ant-design-vue';
import {
  getConfigGroupList, createConfigGroup, updateConfigGroup, deleteConfigGroup,
} from '#/api/system/config';
import type { ConfigGroupRecord } from '#/api/system/config';

const loading = ref(false);
const dataSource = ref<ConfigGroupRecord[]>([]);

const modalVisible = ref(false);
const modalTitle = ref('新增配置分组');
const formState = reactive({ id: 0, name: '', sort: 0 });

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '分组名称', dataIndex: 'name', width: 200 },
  { title: '排序', dataIndex: 'sort', width: 80 },
  { title: '创建时间', dataIndex: 'created_at', width: 180 },
  { title: '操作', dataIndex: 'action', width: 150, fixed: 'right' as const },
];

async function fetchData() {
  loading.value = true;
  try {
    dataSource.value = await getConfigGroupList() || [];
  } finally {
    loading.value = false;
  }
}

function handleAdd() {
  modalTitle.value = '新增配置分组';
  Object.assign(formState, { id: 0, name: '', sort: 0 });
  modalVisible.value = true;
}

function handleEdit(record: ConfigGroupRecord) {
  modalTitle.value = '编辑配置分组';
  Object.assign(formState, { id: record.id, name: record.name, sort: record.sort });
  modalVisible.value = true;
}

async function handleSubmit() {
  if (formState.id) {
    await updateConfigGroup(formState.id, { ...formState });
    message.success('更新成功');
  } else {
    await createConfigGroup({ ...formState });
    message.success('创建成功');
  }
  modalVisible.value = false;
  fetchData();
}

async function handleDelete(id: number) {
  await deleteConfigGroup(id);
  message.success('删除成功');
  fetchData();
}

onMounted(fetchData);
</script>

<template>
  <div class="p-4">
    <Card title="配置分组">
      <template #extra>
        <Button type="primary" @click="handleAdd">新增分组</Button>
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
          <template v-if="column.dataIndex === 'action'">
            <Space>
              <Button size="small" type="link" @click="handleEdit(record as ConfigGroupRecord)">编辑</Button>
              <Popconfirm title="确认删除？" @confirm="handleDelete((record as ConfigGroupRecord).id)">
                <Button size="small" type="link" danger>删除</Button>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <Modal v-model:open="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="420">
      <Form :model="formState" layout="vertical" class="mt-4">
        <FormItem label="分组名称" :rules="[{ required: true }]">
          <Input v-model:value="formState.name" placeholder="请输入分组名称" />
        </FormItem>
        <FormItem label="排序">
          <InputNumber v-model:value="formState.sort" :min="0" style="width: 100%" />
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>
