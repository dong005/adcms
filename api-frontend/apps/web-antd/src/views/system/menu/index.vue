<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Tag, Space, Modal, Form, FormItem, Input, InputNumber,
  Select, SelectOption, Popconfirm, message, TreeSelect, Switch, Card,
  Divider,
} from 'ant-design-vue';
import { getMenuList, createMenu, updateMenu, deleteMenu, getMenuButtons, saveMenuButtons } from '#/api/system/menu';
import type { MenuRecord, ButtonRecord } from '#/api/system/menu';

const loading = ref(false);
const dataSource = ref<MenuRecord[]>([]);
const modalVisible = ref(false);
const modalTitle = ref('新增菜单');
const parentTreeData = ref<any[]>([]);

const formState = reactive({
  id: 0,
  parent_id: 0,
  name: '',
  path: '',
  component: '',
  redirect: '',
  icon: '',
  title: '',
  hide_in_menu: 0,
  keep_alive: 1,
  frame_src: '',
  is_tenant: 1,
  is_public: 0,
  type: 2,
  sort: 0,
  status: 1,
  permission_code: '',
});

const buttonList = ref<ButtonRecord[]>([]);

function addButton() {
  buttonList.value.push({ title: '', name: '', permission_code: '', sort: buttonList.value.length + 1 });
}

function removeButton(index: number) {
  buttonList.value.splice(index, 1);
}

const columns = [
  { title: '标题', dataIndex: 'title', width: 180 },
  { title: '名称', dataIndex: 'name', width: 150 },
  { title: '类型', dataIndex: 'type', width: 80 },
  { title: '路径', dataIndex: 'path', width: 200 },
  { title: '组件', dataIndex: 'component', width: 200 },
  { title: '图标', dataIndex: 'icon', width: 120 },
  { title: '排序', dataIndex: 'sort', width: 70 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '租户可见', dataIndex: 'is_tenant', width: 80 },
  { title: '公共', dataIndex: 'is_public', width: 70 },
  { title: '操作', dataIndex: 'action', width: 200, fixed: 'right' as const },
];

function buildParentTree(menus: MenuRecord[]): any[] {
  return menus.map((m) => ({
    value: m.id,
    title: m.title || m.name,
    children: m.children ? buildParentTree(m.children) : [],
  }));
}

async function fetchData() {
  loading.value = true;
  try {
    dataSource.value = await getMenuList();
    parentTreeData.value = [
      { value: 0, title: '顶级菜单', children: buildParentTree(dataSource.value) },
    ];
  } finally {
    loading.value = false;
  }
}

function handleAdd(parentId = 0) {
  modalTitle.value = '新增菜单';
  Object.assign(formState, {
    id: 0, parent_id: parentId, name: '', path: '', component: '', redirect: '',
    icon: '', title: '', hide_in_menu: 0, keep_alive: 1, frame_src: '',
    is_tenant: 1, is_public: 0, type: 2,
    sort: 0, status: 1, permission_code: '',
  });
  buttonList.value = [];
  modalVisible.value = true;
}

async function handleEdit(record: MenuRecord) {
  modalTitle.value = '编辑菜单';
  Object.assign(formState, {
    id: record.id, parent_id: record.parent_id, name: record.name,
    path: record.path, component: record.component, redirect: record.redirect,
    icon: record.icon, title: record.title, hide_in_menu: record.hide_in_menu,
    keep_alive: record.keep_alive, frame_src: record.frame_src,
    is_tenant: record.is_tenant, is_public: record.is_public, type: record.type,
    sort: record.sort, status: record.status, permission_code: record.permission_code,
  });
  // 加载按钮节点
  try {
    const buttons = await getMenuButtons(record.id);
    buttonList.value = (buttons || []).map((b: any) => ({
      title: b.title, name: b.name, permission_code: b.permission_code, sort: b.sort,
    }));
  } catch {
    buttonList.value = [];
  }
  modalVisible.value = true;
}

async function handleSubmit() {
  let menuId = formState.id;
  if (menuId) {
    await updateMenu(menuId, { ...formState });
  } else {
    const res: any = await createMenu({ ...formState });
    menuId = res?.id || res?.ID;
  }
  // 保存按钮节点
  if (menuId && buttonList.value.length > 0) {
    await saveMenuButtons(menuId, buttonList.value);
  }
  message.success(formState.id ? '更新成功' : '创建成功');
  modalVisible.value = false;
  fetchData();
}

async function handleDelete(id: number) {
  try {
    await deleteMenu(id);
    message.success('删除成功');
    fetchData();
  } catch (e: any) {
    message.error(e?.message || '删除失败，请先删除子菜单');
  }
}

onMounted(fetchData);
</script>

<template>
  <div class="p-4">
    <Card title="菜单管理">
      <template #extra>
        <Button type="primary" @click="handleAdd(0)">新增菜单</Button>
      </template>

      <Table
        :columns="columns"
        :data-source="dataSource"
        :loading="loading"
        :pagination="false"
        :scroll="{ x: 1200 }"
        row-key="id"
        size="middle"
        default-expand-all-rows
      >
        <template #bodyCell="{ column, record }">
          <template v-if="column.dataIndex === 'type'">
            <Tag :color="(record as MenuRecord).type === 1 ? 'blue' : (record as MenuRecord).type === 2 ? 'green' : (record as MenuRecord).type === 3 ? 'orange' : 'purple'">
              {{ (record as MenuRecord).type === 1 ? '目录' : (record as MenuRecord).type === 2 ? '菜单' : (record as MenuRecord).type === 3 ? '页面' : '按钮' }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'status'">
            <Tag :color="(record as MenuRecord).status === 1 ? 'green' : 'red'">
              {{ (record as MenuRecord).status === 1 ? '启用' : '禁用' }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'is_tenant'">
            <Tag :color="(record as MenuRecord).is_tenant === 1 ? 'green' : 'red'">
              {{ (record as MenuRecord).is_tenant === 1 ? '是' : '否' }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'is_public'">
            <Tag :color="(record as MenuRecord).is_public === 1 ? 'blue' : 'default'">
              {{ (record as MenuRecord).is_public === 1 ? '是' : '否' }}
            </Tag>
          </template>
          <template v-if="column.dataIndex === 'action'">
            <Space>
              <Button size="small" type="link" @click="handleAdd((record as MenuRecord).id)">添加子菜单</Button>
              <Button size="small" type="link" @click="handleEdit(record as MenuRecord)">编辑</Button>
              <Popconfirm title="确认删除？" @confirm="handleDelete((record as MenuRecord).id)">
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
          <FormItem label="父级菜单">
            <TreeSelect
              v-model:value="formState.parent_id"
              :tree-data="parentTreeData"
              placeholder="请选择父级菜单"
              tree-default-expand-all
              style="width: 100%"
            />
          </FormItem>
          <FormItem label="菜单标题" :rules="[{ required: true }]">
            <Input v-model:value="formState.title" placeholder="显示的标题" />
          </FormItem>
          <FormItem label="路由名称" :rules="[{ required: true }]">
            <Input v-model:value="formState.name" placeholder="唯一路由名称" />
          </FormItem>
          <FormItem label="路由路径">
            <Input v-model:value="formState.path" placeholder="如 /system/user" />
          </FormItem>
          <FormItem label="组件路径">
            <Input v-model:value="formState.component" placeholder="如 BasicLayout 或 /system/user/index" />
          </FormItem>
          <FormItem label="重定向">
            <Input v-model:value="formState.redirect" placeholder="重定向路径" />
          </FormItem>
          <FormItem label="图标">
            <Input v-model:value="formState.icon" placeholder="如 lucide:users" />
          </FormItem>
          <FormItem label="类型">
            <Select v-model:value="formState.type">
              <SelectOption :value="1">目录</SelectOption>
              <SelectOption :value="2">菜单</SelectOption>
              <SelectOption :value="3">页面</SelectOption>
              <SelectOption :value="4">按钮</SelectOption>
            </Select>
          </FormItem>
          <FormItem label="权限标识">
            <Input v-model:value="formState.permission_code" placeholder="如 user:list" />
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
          <FormItem label="租户可见">
            <Switch v-model:checked="formState.is_tenant" :checked-value="1" :un-checked-value="0" />
          </FormItem>
          <FormItem label="公共菜单">
            <Switch v-model:checked="formState.is_public" :checked-value="1" :un-checked-value="0" />
          </FormItem>
          <FormItem label="隐藏菜单">
            <Switch v-model:checked="formState.hide_in_menu" :checked-value="1" :un-checked-value="0" />
          </FormItem>
          <FormItem label="缓存">
            <Switch v-model:checked="formState.keep_alive" :checked-value="1" :un-checked-value="0" />
          </FormItem>
        </div>
        <FormItem label="iframe地址">
          <Input v-model:value="formState.frame_src" placeholder="外部iframe地址" />
        </FormItem>
      </Form>

      <Divider v-if="formState.type === 2 || formState.type === 3" orientation="left" style="margin-top: 0">
        权限按钮
      </Divider>
      <div v-if="formState.type === 2 || formState.type === 3">
        <div v-for="(btn, idx) in buttonList" :key="idx" style="display: flex; gap: 8px; margin-bottom: 8px; align-items: center;">
          <Input v-model:value="btn.title" placeholder="按钮名称" style="width: 120px" />
          <Input v-model:value="btn.name" placeholder="标识名" style="width: 140px" />
          <Input v-model:value="btn.permission_code" placeholder="权限码 如 user:create" style="flex: 1" />
          <InputNumber v-model:value="btn.sort" :min="0" placeholder="排序" style="width: 70px" />
          <Button type="text" danger size="small" @click="removeButton(idx)">×</Button>
        </div>
        <Button type="dashed" block @click="addButton">+ 添加按钮</Button>
      </div>
    </Modal>
  </div>
</template>
