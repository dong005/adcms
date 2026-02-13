<script lang="ts" setup>
import { ref, onMounted } from 'vue';
import {
  Table, Button, Card, Modal, Tag,
} from 'ant-design-vue';
import { getTableList, getTableColumns } from '#/api/system/database';
import type { TableInfo, ColumnInfo } from '#/api/system/database';

const loading = ref(false);
const tableList = ref<TableInfo[]>([]);

const columnModalVisible = ref(false);
const columnModalTitle = ref('');
const columnList = ref<ColumnInfo[]>([]);
const columnLoading = ref(false);

const tableColumns = [
  { title: '表名', dataIndex: 'name', width: 200 },
  { title: '引擎', dataIndex: 'engine', width: 100 },
  { title: '行数', dataIndex: 'rows', width: 100 },
  { title: '数据大小', dataIndex: 'data_size', width: 120 },
  { title: '索引大小', dataIndex: 'index_size', width: 120 },
  { title: '备注', dataIndex: 'comment', width: 200 },
  { title: '操作', dataIndex: 'action', width: 100, fixed: 'right' as const },
];

const colColumns = [
  { title: '字段名', dataIndex: 'name', width: 150 },
  { title: '类型', dataIndex: 'type', width: 180 },
  { title: '可空', dataIndex: 'nullable', width: 70 },
  { title: '键', dataIndex: 'key', width: 70 },
  { title: '默认值', dataIndex: 'default', width: 120 },
  { title: '备注', dataIndex: 'comment', width: 200 },
];

async function fetchTables() {
  loading.value = true;
  try {
    tableList.value = await getTableList() || [];
  } finally {
    loading.value = false;
  }
}

async function handleViewColumns(tableName: string) {
  columnModalTitle.value = `表结构 - ${tableName}`;
  columnModalVisible.value = true;
  columnLoading.value = true;
  try {
    columnList.value = await getTableColumns(tableName) || [];
  } finally {
    columnLoading.value = false;
  }
}

onMounted(fetchTables);
</script>

<template>
  <div class="p-4">
    <Card title="数据库管理">
      <template #extra>
        <Button @click="fetchTables">刷新</Button>
      </template>
      <Table
        :columns="tableColumns"
        :data-source="tableList"
        :loading="loading"
        :pagination="false"
        row-key="name"
        size="middle"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'action'">
            <Button size="small" type="link" @click="handleViewColumns((record as TableInfo).name)">查看结构</Button>
          </template>
        </template>
      </Table>
    </Card>

    <Modal v-model:open="columnModalVisible" :title="columnModalTitle" :footer="null" :width="800">
      <Table
        :columns="colColumns"
        :data-source="columnList"
        :loading="columnLoading"
        :pagination="false"
        row-key="name"
        size="small"
        class="mt-4"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'key'">
            <Tag v-if="(record as ColumnInfo).key === 'PRI'" color="blue">PRI</Tag>
            <Tag v-else-if="(record as ColumnInfo).key === 'UNI'" color="green">UNI</Tag>
            <Tag v-else-if="(record as ColumnInfo).key === 'MUL'" color="orange">MUL</Tag>
            <span v-else>-</span>
          </template>
          <template v-if="column.dataIndex === 'nullable'">
            <Tag :color="(record as ColumnInfo).nullable === 'YES' ? 'green' : 'red'">
              {{ (record as ColumnInfo).nullable }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'default'">
            {{ (record as ColumnInfo).default ?? '-' }}
          </template>
        </template>
      </Table>
    </Modal>
  </div>
</template>
