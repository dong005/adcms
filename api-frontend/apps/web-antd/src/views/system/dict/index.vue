<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Tag, Space, Modal, Form, FormItem, Input, InputNumber,
  Select, SelectOption, Popconfirm, message, Card, Row, Col, List, ListItem,
} from 'ant-design-vue';
import {
  getDictTypeList, createDictType, updateDictType, deleteDictType,
  getDictList, createDict, updateDict, deleteDict,
} from '#/api/system/dict';
import type { DictTypeRecord, DictRecord } from '#/api/system/dict';

// ========== 字典类型 ==========
const typeLoading = ref(false);
const typeList = ref<DictTypeRecord[]>([]);
const typeKeyword = ref('');
const selectedTypeId = ref<number>(0);
const selectedTypeName = ref('');

const typeModalVisible = ref(false);
const typeModalTitle = ref('新增字典类型');
const typeForm = reactive({
  id: 0,
  name: '',
  code: '',
  sort: 0,
  status: 1,
  remark: '',
});

async function fetchTypes() {
  typeLoading.value = true;
  try {
    typeList.value = await getDictTypeList(typeKeyword.value) || [];
    if (typeList.value.length > 0 && !selectedTypeId.value) {
      selectType(typeList.value[0]!);
    }
  } finally {
    typeLoading.value = false;
  }
}

function selectType(item: DictTypeRecord) {
  selectedTypeId.value = item.id;
  selectedTypeName.value = item.name;
  fetchDicts();
}

function handleAddType() {
  typeModalTitle.value = '新增字典类型';
  Object.assign(typeForm, { id: 0, name: '', code: '', sort: 0, status: 1, remark: '' });
  typeModalVisible.value = true;
}

function handleEditType(item: DictTypeRecord) {
  typeModalTitle.value = '编辑字典类型';
  Object.assign(typeForm, {
    id: item.id, name: item.name, code: item.code,
    sort: item.sort, status: item.status, remark: item.remark,
  });
  typeModalVisible.value = true;
}

async function handleSubmitType() {
  if (typeForm.id) {
    await updateDictType(typeForm.id, { ...typeForm });
    message.success('更新成功');
  } else {
    await createDictType({ ...typeForm });
    message.success('创建成功');
  }
  typeModalVisible.value = false;
  fetchTypes();
}

async function handleDeleteType(id: number) {
  await deleteDictType(id);
  message.success('删除成功');
  if (selectedTypeId.value === id) {
    selectedTypeId.value = 0;
    dictList.value = [];
  }
  fetchTypes();
}

// ========== 字典数据 ==========
const dictLoading = ref(false);
const dictList = ref<DictRecord[]>([]);

const dictModalVisible = ref(false);
const dictModalTitle = ref('新增字典数据');
const dictForm = reactive({
  id: 0,
  dict_type_id: 0,
  name: '',
  value: '',
  sort: 0,
  status: 1,
  remark: '',
});

const dictColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '字典标签', dataIndex: 'name', width: 150 },
  { title: '字典值', dataIndex: 'value', width: 150 },
  { title: '排序', dataIndex: 'sort', width: 70 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '备注', dataIndex: 'remark', width: 200 },
  { title: '操作', dataIndex: 'action', width: 150, fixed: 'right' as const },
];

async function fetchDicts() {
  if (!selectedTypeId.value) return;
  dictLoading.value = true;
  try {
    dictList.value = await getDictList(selectedTypeId.value) || [];
  } finally {
    dictLoading.value = false;
  }
}

function handleAddDict() {
  if (!selectedTypeId.value) {
    message.warning('请先选择字典类型');
    return;
  }
  dictModalTitle.value = '新增字典数据';
  Object.assign(dictForm, {
    id: 0, dict_type_id: selectedTypeId.value,
    name: '', value: '', sort: 0, status: 1, remark: '',
  });
  dictModalVisible.value = true;
}

function handleEditDict(record: DictRecord) {
  dictModalTitle.value = '编辑字典数据';
  Object.assign(dictForm, {
    id: record.id, dict_type_id: record.dict_type_id,
    name: record.name, value: record.value,
    sort: record.sort, status: record.status, remark: record.remark,
  });
  dictModalVisible.value = true;
}

async function handleSubmitDict() {
  if (dictForm.id) {
    await updateDict(dictForm.id, { ...dictForm });
    message.success('更新成功');
  } else {
    await createDict({ ...dictForm });
    message.success('创建成功');
  }
  dictModalVisible.value = false;
  fetchDicts();
}

async function handleDeleteDict(id: number) {
  await deleteDict(id);
  message.success('删除成功');
  fetchDicts();
}

onMounted(() => {
  fetchTypes();
});
</script>

<template>
  <div class="p-4">
    <Row :gutter="16">
      <Col :span="6">
        <Card title="字典类型" :loading="typeLoading">
          <template #extra>
            <Button type="primary" size="small" @click="handleAddType">新增</Button>
          </template>
          <Input.Search
            v-model:value="typeKeyword"
            placeholder="搜索类型"
            class="mb-3"
            @search="fetchTypes"
          />
          <List size="small" :data-source="typeList" :split="true">
            <template #renderItem="{ item }">
              <ListItem
                class="cursor-pointer"
                :class="{ 'dict-type-active': selectedTypeId === (item as DictTypeRecord).id }"
                @click="selectType(item as DictTypeRecord)"
              >
                <div class="flex items-center justify-between w-full">
                  <span>{{ (item as DictTypeRecord).name }}</span>
                  <Space size="small">
                    <Button type="link" size="small" @click.stop="handleEditType(item as DictTypeRecord)">编辑</Button>
                    <Popconfirm title="删除将同时删除该类型下所有数据，确认？" @confirm="handleDeleteType((item as DictTypeRecord).id)">
                      <Button type="link" size="small" danger @click.stop>删除</Button>
                    </Popconfirm>
                  </Space>
                </div>
              </ListItem>
            </template>
          </List>
        </Card>
      </Col>

      <Col :span="18">
        <Card :title="selectedTypeName ? `字典数据 - ${selectedTypeName}` : '字典数据'">
          <template #extra>
            <Button type="primary" @click="handleAddDict" :disabled="!selectedTypeId">新增数据</Button>
          </template>
          <Table
            :columns="dictColumns"
            :data-source="dictList"
            :loading="dictLoading"
            :pagination="false"
            row-key="id"
            size="middle"
          >
            <template #bodyCell="{ column, record }">
              <template v-if="column.dataIndex === 'status'">
                <Tag :color="(record as DictRecord).status === 1 ? 'green' : 'red'">
                  {{ (record as DictRecord).status === 1 ? '启用' : '禁用' }}
                </Tag>
              </template>
              <template v-if="column.dataIndex === 'action'">
                <Space>
                  <Button size="small" type="link" @click="handleEditDict(record as DictRecord)">编辑</Button>
                  <Popconfirm title="确认删除？" @confirm="handleDeleteDict((record as DictRecord).id)">
                    <Button size="small" type="link" danger>删除</Button>
                  </Popconfirm>
                </Space>
              </template>
            </template>
          </Table>
        </Card>
      </Col>
    </Row>

    <!-- 字典类型弹窗 -->
    <Modal v-model:open="typeModalVisible" :title="typeModalTitle" @ok="handleSubmitType" :width="520">
      <Form :model="typeForm" layout="vertical" class="mt-4">
        <FormItem label="类型名称" :rules="[{ required: true }]">
          <Input v-model:value="typeForm.name" placeholder="请输入类型名称" />
        </FormItem>
        <FormItem label="类型编码" :rules="[{ required: true }]">
          <Input v-model:value="typeForm.code" placeholder="请输入类型编码" :disabled="!!typeForm.id" />
        </FormItem>
        <div class="grid grid-cols-2 gap-x-4">
          <FormItem label="排序">
            <InputNumber v-model:value="typeForm.sort" :min="0" style="width: 100%" />
          </FormItem>
          <FormItem label="状态">
            <Select v-model:value="typeForm.status">
              <SelectOption :value="1">启用</SelectOption>
              <SelectOption :value="0">禁用</SelectOption>
            </Select>
          </FormItem>
        </div>
        <FormItem label="备注">
          <Input.TextArea v-model:value="typeForm.remark" placeholder="请输入备注" :rows="3" />
        </FormItem>
      </Form>
    </Modal>

    <!-- 字典数据弹窗 -->
    <Modal v-model:open="dictModalVisible" :title="dictModalTitle" @ok="handleSubmitDict" :width="520">
      <Form :model="dictForm" layout="vertical" class="mt-4">
        <FormItem label="字典标签" :rules="[{ required: true }]">
          <Input v-model:value="dictForm.name" placeholder="请输入字典标签" />
        </FormItem>
        <FormItem label="字典值" :rules="[{ required: true }]">
          <Input v-model:value="dictForm.value" placeholder="请输入字典值" />
        </FormItem>
        <div class="grid grid-cols-2 gap-x-4">
          <FormItem label="排序">
            <InputNumber v-model:value="dictForm.sort" :min="0" style="width: 100%" />
          </FormItem>
          <FormItem label="状态">
            <Select v-model:value="dictForm.status">
              <SelectOption :value="1">启用</SelectOption>
              <SelectOption :value="0">禁用</SelectOption>
            </Select>
          </FormItem>
        </div>
        <FormItem label="备注">
          <Input.TextArea v-model:value="dictForm.remark" placeholder="请输入备注" :rows="3" />
        </FormItem>
      </Form>
    </Modal>
  </div>
</template>

<style scoped>
.dict-type-active {
  background-color: var(--ant-primary-1, rgba(22, 119, 255, 0.1));
  border-radius: 4px;
}
</style>
