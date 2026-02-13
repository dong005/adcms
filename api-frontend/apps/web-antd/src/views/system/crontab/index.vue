<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Tag, Space, Modal, Form, FormItem, Input,
  Select, SelectOption, Popconfirm, message, Card,
} from 'ant-design-vue';
import { getCrontabList, createCrontab, updateCrontab, deleteCrontab } from '#/api/system/crontab';
import type { CrontabRecord } from '#/api/system/crontab';

const loading = ref(false);
const dataSource = ref<CrontabRecord[]>([]);

const modalVisible = ref(false);
const modalTitle = ref('新增定时任务');
const formState = reactive({
  id: 0,
  name: '',
  expression: '',
  command: '',
  status: 0,
  remark: '',
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '任务名称', dataIndex: 'name', width: 150 },
  { title: 'Cron表达式', dataIndex: 'expression', width: 150 },
  { title: '执行命令', dataIndex: 'command', width: 250 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '上次执行', dataIndex: 'last_run_at', width: 170 },
  { title: '下次执行', dataIndex: 'next_run_at', width: 170 },
  { title: '操作', dataIndex: 'action', width: 150, fixed: 'right' as const },
];

async function fetchData() {
  loading.value = true;
  try {
    dataSource.value = await getCrontabList() || [];
  } finally {
    loading.value = false;
  }
}

function handleAdd() {
  modalTitle.value = '新增定时任务';
  Object.assign(formState, { id: 0, name: '', expression: '', command: '', status: 0, remark: '' });
  modalVisible.value = true;
}

function handleEdit(record: CrontabRecord) {
  modalTitle.value = '编辑定时任务';
  Object.assign(formState, {
    id: record.id, name: record.name, expression: record.expression,
    command: record.command, status: record.status, remark: record.remark,
  });
  modalVisible.value = true;
}

async function handleSubmit() {
  if (formState.id) {
    await updateCrontab(formState.id, { ...formState });
    message.success('更新成功');
  } else {
    await createCrontab({ ...formState });
    message.success('创建成功');
  }
  modalVisible.value = false;
  fetchData();
}

async function handleDelete(id: number) {
  await deleteCrontab(id);
  message.success('删除成功');
  fetchData();
}

onMounted(fetchData);
</script>

<template>
  <div class="p-4">
    <Card title="定时任务管理">
      <template #extra>
        <Button type="primary" @click="handleAdd">新增任务</Button>
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
          <template v-if="column.dataIndex === 'status'">
            <Tag :color="(record as CrontabRecord).status === 1 ? 'green' : 'red'">
              {{ (record as CrontabRecord).status === 1 ? '运行中' : '已停止' }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'last_run_at'">
            {{ (record as CrontabRecord).last_run_at || '-' }}
          </template>
          <template v-if="column.dataIndex === 'next_run_at'">
            {{ (record as CrontabRecord).next_run_at || '-' }}
          </template>
          <template v-if="column.dataIndex === 'action'">
            <Space>
              <Button size="small" type="link" @click="handleEdit(record as CrontabRecord)">编辑</Button>
              <Popconfirm title="确认删除？" @confirm="handleDelete((record as CrontabRecord).id)">
                <Button size="small" type="link" danger>删除</Button>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <Modal v-model:open="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="520">
      <Form :model="formState" layout="vertical" class="mt-4">
        <FormItem label="任务名称" :rules="[{ required: true }]">
          <Input v-model:value="formState.name" placeholder="请输入任务名称" />
        </FormItem>
        <FormItem label="Cron表达式" :rules="[{ required: true }]">
          <Input v-model:value="formState.expression" placeholder="如 0 0 * * * (每小时)" />
        </FormItem>
        <FormItem label="执行命令" :rules="[{ required: true }]">
          <Input.TextArea v-model:value="formState.command" placeholder="请输入执行命令或URL" :rows="3" />
        </FormItem>
        <FormItem label="状态">
          <Select v-model:value="formState.status">
            <SelectOption :value="1">启用</SelectOption>
            <SelectOption :value="0">停止</SelectOption>
          </Select>
        </FormItem>
        <FormItem label="备注">
          <Input.TextArea v-model:value="formState.remark" placeholder="请输入备注" :rows="2" />
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>
