<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Space, Modal, Form, FormItem, Input, InputNumber,
  Popconfirm, message, Card,
} from 'ant-design-vue';
import { getCityList, createCity, updateCity, deleteCity } from '#/api/system/city';
import type { CityRecord } from '#/api/system/city';

const loading = ref(false);
const dataSource = ref<CityRecord[]>([]);
const currentPid = ref(0);
const breadcrumbs = ref<{ id: number; name: string }[]>([{ id: 0, name: '全部' }]);

const modalVisible = ref(false);
const modalTitle = ref('新增区域');
const formState = reactive({
  id: 0,
  pid: 0,
  level: 0,
  name: '',
  citycode: '',
  adcode: '',
  p_adcode: '',
  lng: 0,
  lat: 0,
  sort: 0,
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '名称', dataIndex: 'name', width: 200 },
  { title: '城市编码', dataIndex: 'citycode', width: 100 },
  { title: '区域编码', dataIndex: 'adcode', width: 100 },
  { title: '经度', dataIndex: 'lng', width: 120 },
  { title: '纬度', dataIndex: 'lat', width: 120 },
  { title: '排序', dataIndex: 'sort', width: 70 },
  { title: '操作', dataIndex: 'action', width: 200, fixed: 'right' as const },
];

async function fetchData() {
  loading.value = true;
  try {
    dataSource.value = await getCityList(currentPid.value) || [];
  } finally {
    loading.value = false;
  }
}

function drillDown(record: CityRecord) {
  currentPid.value = record.id;
  breadcrumbs.value.push({ id: record.id, name: record.name });
  fetchData();
}

function goTo(index: number) {
  const item = breadcrumbs.value[index]!;
  currentPid.value = item.id;
  breadcrumbs.value = breadcrumbs.value.slice(0, index + 1);
  fetchData();
}

function handleAdd() {
  modalTitle.value = '新增区域';
  const level = breadcrumbs.value.length - 1;
  Object.assign(formState, {
    id: 0, pid: currentPid.value, level,
    name: '', citycode: '', adcode: '', p_adcode: '',
    lng: 0, lat: 0, sort: 0,
  });
  modalVisible.value = true;
}

function handleEdit(record: CityRecord) {
  modalTitle.value = '编辑区域';
  Object.assign(formState, {
    id: record.id, pid: record.pid, level: record.level,
    name: record.name, citycode: record.citycode, adcode: record.adcode,
    p_adcode: record.p_adcode, lng: record.lng, lat: record.lat, sort: record.sort,
  });
  modalVisible.value = true;
}

async function handleSubmit() {
  if (formState.id) {
    await updateCity(formState.id, { ...formState });
    message.success('更新成功');
  } else {
    await createCity({ ...formState });
    message.success('创建成功');
  }
  modalVisible.value = false;
  fetchData();
}

async function handleDelete(id: number) {
  await deleteCity(id);
  message.success('删除成功');
  fetchData();
}

onMounted(fetchData);
</script>

<template>
  <div class="p-4">
    <Card title="行政区域管理">
      <template #extra>
        <Button type="primary" @click="handleAdd">新增区域</Button>
      </template>

      <div class="mb-3">
        <span
          v-for="(item, index) in breadcrumbs"
          :key="item.id"
        >
          <a v-if="index < breadcrumbs.length - 1" class="text-blue-500 cursor-pointer" @click="goTo(index)">{{ item.name }}</a>
          <span v-else class="text-gray-500">{{ item.name }}</span>
          <span v-if="index < breadcrumbs.length - 1" class="mx-1">/</span>
        </span>
      </div>

      <Table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
        :pagination="false"
        row-key="id"
        size="middle"
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'name'">
            <a class="text-blue-500 cursor-pointer" @click="drillDown(record as CityRecord)">
              {{ (record as CityRecord).name }}
            </a>
          </template>
          <template v-if="column.dataIndex === 'action'">
            <Space>
              <Button size="small" type="link" @click="drillDown(record as CityRecord)">下级</Button>
              <Button size="small" type="link" @click="handleEdit(record as CityRecord)">编辑</Button>
              <Popconfirm title="删除将同时删除所有下级区域，确认？" @confirm="handleDelete((record as CityRecord).id)">
                <Button size="small" type="link" danger>删除</Button>
              </Popconfirm>
            </Space>
          </template>
        </template>
      </Table>
    </Card>

    <Modal v-model:open="modalVisible" :title="modalTitle" @ok="handleSubmit" :width="640">
      <Form :model="formState" layout="vertical" class="mt-4">
        <FormItem label="名称" :rules="[{ required: true }]">
          <Input v-model:value="formState.name" placeholder="请输入区域名称" />
        </FormItem>
        <div class="grid grid-cols-2 gap-x-4">
          <FormItem label="城市编码">
            <Input v-model:value="formState.citycode" placeholder="citycode" />
          </FormItem>
          <FormItem label="区域编码">
            <Input v-model:value="formState.adcode" placeholder="adcode" />
          </FormItem>
          <FormItem label="经度">
            <InputNumber v-model:value="formState.lng" :step="0.000001" style="width: 100%" />
          </FormItem>
          <FormItem label="纬度">
            <InputNumber v-model:value="formState.lat" :step="0.000001" style="width: 100%" />
          </FormItem>
          <FormItem label="排序">
            <InputNumber v-model:value="formState.sort" :min="0" style="width: 100%" />
          </FormItem>
        </div>
      </Form>
    </Modal>
  </div>
</template>
