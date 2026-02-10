<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Tag, Space, Card, Input, Select, SelectOption,
  RangePicker, Descriptions, Switch, Form, FormItem, message,
  Tabs, Modal,
} from 'ant-design-vue';
import {
  getOperationLogList, getLoginLogList,
  getEmailLogList, getSmsLogList,
  getLogConfig, updateLogConfig,
} from '#/api/system/log';
import type {
  OperationLogRecord, LoginLogRecord,
  EmailLogRecord, SmsLogRecord,
  LogConfig,
} from '#/api/system/log';
import dayjs from 'dayjs';

const activeTab = ref('operation');
const loading = ref(false);
const operationData = ref<OperationLogRecord[]>([]);
const loginData = ref<LoginLogRecord[]>([]);
const emailData = ref<EmailLogRecord[]>([]);
const smsData = ref<SmsLogRecord[]>([]);
const total = ref(0);
const keyword = ref('');
const moduleFilter = ref<string>('');
const statusFilter = ref<number | undefined>();
const dateRange = ref<[dayjs.Dayjs, dayjs.Dayjs] | undefined>(undefined);
const pagination = reactive({ current: 1, pageSize: 20 });

const detailVisible = ref(false);
const detailRecord = ref<OperationLogRecord | LoginLogRecord | null>(null);

const operationColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '用户', dataIndex: 'username', width: 120 },
  { title: '模块', dataIndex: 'module', width: 100 },
  { title: '操作', dataIndex: 'action', width: 120 },
  { title: '方法', dataIndex: 'method', width: 80 },
  { title: '路径', dataIndex: 'path', width: 200 },
  { title: 'IP', dataIndex: 'ip', width: 120 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '耗时', dataIndex: 'duration', width: 80 },
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', dataIndex: 'action', width: 80, fixed: 'right' as const },
];

const loginColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '用户', dataIndex: 'username', width: 120 },
  { title: '登录类型', dataIndex: 'login_type', width: 100 },
  { title: 'IP', dataIndex: 'ip', width: 120 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', dataIndex: 'action', width: 80, fixed: 'right' as const },
];

const emailColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '收件人', dataIndex: 'to', width: 200 },
  { title: '主题', dataIndex: 'subject', ellipsis: true },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '错误', dataIndex: 'error', ellipsis: true, width: 200 },
  { title: '时间', dataIndex: 'created_at', width: 170 },
];

const smsColumns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '手机号', dataIndex: 'phone', width: 140 },
  { title: '模板ID', dataIndex: 'template_id', width: 120 },
  { title: '参数', dataIndex: 'params', width: 150 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '错误', dataIndex: 'error', ellipsis: true, width: 200 },
  { title: '时间', dataIndex: 'created_at', width: 170 },
];

async function fetchOperationLogs() {
  loading.value = true;
  try {
    const params: any = {
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: keyword.value,
    };
    if (moduleFilter.value) params.module = moduleFilter.value;
    if (statusFilter.value !== undefined) params.status = statusFilter.value;
    if (dateRange.value) {
      params.start_time = dateRange.value[0].toISOString();
      params.end_time = dateRange.value[1].toISOString();
    }
    
    const res = await getOperationLogList(params);
    operationData.value = res.list || [];
    total.value = res.total;
  } finally {
    loading.value = false;
  }
}

async function fetchLoginLogs() {
  loading.value = true;
  try {
    const params: any = {
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: keyword.value,
    };
    if (statusFilter.value !== undefined) params.status = statusFilter.value;
    if (dateRange.value) {
      params.start_time = dateRange.value[0].toISOString();
      params.end_time = dateRange.value[1].toISOString();
    }
    
    const res = await getLoginLogList(params);
    loginData.value = res.list || [];
    total.value = res.total;
  } finally {
    loading.value = false;
  }
}

async function fetchEmailLogs() {
  loading.value = true;
  try {
    const res = await getEmailLogList({
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: keyword.value || undefined,
    });
    emailData.value = res.list || [];
    total.value = res.total;
  } finally {
    loading.value = false;
  }
}

async function fetchSmsLogs() {
  loading.value = true;
  try {
    const res = await getSmsLogList({
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: keyword.value || undefined,
    });
    smsData.value = res.list || [];
    total.value = res.total;
  } finally {
    loading.value = false;
  }
}

function fetchCurrentTab() {
  const map: Record<string, () => void> = {
    operation: fetchOperationLogs,
    login: fetchLoginLogs,
    email: fetchEmailLogs,
    sms: fetchSmsLogs,
    settings: fetchLogConfig,
  };
  (map[activeTab.value] || fetchOperationLogs)();
}

function handleSearch() {
  pagination.current = 1;
  fetchCurrentTab();
}

function handleViewDetail(record: any) {
  detailRecord.value = record;
  detailVisible.value = true;
}

function handleTableChange(pag: any) {
  pagination.current = pag.current;
  pagination.pageSize = pag.pageSize;
  fetchCurrentTab();
}

function onTabChange() {
  pagination.current = 1;
  keyword.value = '';
  moduleFilter.value = '';
  statusFilter.value = undefined;
  dateRange.value = undefined;
  fetchCurrentTab();
}

// 日志设置
const logConfig = reactive<LogConfig>({
  log_operation_enabled: '1',
  log_login_enabled: '1',
  log_email_enabled: '1',
  log_sms_enabled: '1',
});
const configLoading = ref(false);

async function fetchLogConfig() {
  configLoading.value = true;
  try {
    const res = await getLogConfig();
    Object.assign(logConfig, res);
  } finally {
    configLoading.value = false;
  }
}

async function handleSaveLogConfig() {
  configLoading.value = true;
  try {
    await updateLogConfig({ ...logConfig });
    message.success('日志配置已保存');
  } catch (e: any) {
    message.error(e?.message || '保存失败');
  } finally {
    configLoading.value = false;
  }
}

onMounted(() => {
  fetchOperationLogs();
});
</script>

<template>
  <Card>
    <Tabs v-model:activeKey="activeTab" @change="onTabChange">
      <Tabs.TabPane key="operation" tab="操作日志">
        <Space class="mb-4">
          <Input v-model:value="keyword" placeholder="搜索用户/路径" style="width: 200px" />
          <Select v-model:value="moduleFilter" placeholder="模块" allow-clear style="width: 120px">
            <SelectOption value="user">用户管理</SelectOption>
            <SelectOption value="role">角色管理</SelectOption>
            <SelectOption value="menu">菜单管理</SelectOption>
            <SelectOption value="tenant">租户管理</SelectOption>
            <SelectOption value="article">文章管理</SelectOption>
            <SelectOption value="category">分类管理</SelectOption>
            <SelectOption value="tag">标签管理</SelectOption>
            <SelectOption value="media">媒体管理</SelectOption>
          </Select>
          <Select v-model:value="statusFilter" placeholder="状态" allow-clear style="width: 100px">
            <SelectOption :value="1">成功</SelectOption>
            <SelectOption :value="0">失败</SelectOption>
          </Select>
          <RangePicker v-model:value="dateRange" />
          <Button type="primary" @click="handleSearch">搜索</Button>
        </Space>

        <Table
          :columns="operationColumns"
          :data-source="operationData"
          :loading="loading"
          :pagination="{ ...pagination, total }"
          row-key="id"
          @change="handleTableChange"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'status'">
              <Tag :color="record.status === 1 ? 'green' : 'red'">
                {{ record.status === 1 ? '成功' : '失败' }}
              </Tag>
            </template>
            <template v-if="column.dataIndex === 'duration'">
              {{ record.duration }}ms
            </template>
            <template v-if="column.dataIndex === 'created_at'">
              {{ dayjs(record.created_at).format('YYYY-MM-DD HH:mm:ss') }}
            </template>
            <template v-if="column.dataIndex === 'action'">
              <Button size="small" type="link" @click="handleViewDetail(record)">
                详情
              </Button>
            </template>
          </template>
        </Table>
      </Tabs.TabPane>

      <Tabs.TabPane key="login" tab="登录日志">
        <Space class="mb-4">
          <Input v-model:value="keyword" placeholder="搜索用户" style="width: 200px" />
          <Select v-model:value="statusFilter" placeholder="状态" allow-clear style="width: 100px">
            <SelectOption :value="1">成功</SelectOption>
            <SelectOption :value="0">失败</SelectOption>
          </Select>
          <RangePicker v-model:value="dateRange" />
          <Button type="primary" @click="handleSearch">搜索</Button>
        </Space>

        <Table
          :columns="loginColumns"
          :data-source="loginData"
          :loading="loading"
          :pagination="{ ...pagination, total }"
          row-key="id"
          @change="handleTableChange"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'login_type'">
              <Tag>{{ record.login_type === 'password' ? '密码登录' : 'TOTP登录' }}</Tag>
            </template>
            <template v-if="column.dataIndex === 'status'">
              <Tag :color="record.status === 1 ? 'green' : 'red'">
                {{ record.status === 1 ? '成功' : '失败' }}
              </Tag>
            </template>
            <template v-if="column.dataIndex === 'created_at'">
              {{ dayjs(record.created_at).format('YYYY-MM-DD HH:mm:ss') }}
            </template>
            <template v-if="column.dataIndex === 'action'">
              <Button size="small" type="link" @click="handleViewDetail(record)">
                详情
              </Button>
            </template>
          </template>
        </Table>
      </Tabs.TabPane>

      <Tabs.TabPane key="email" tab="邮件日志">
        <Space class="mb-4">
          <Input v-model:value="keyword" placeholder="搜索收件人/主题" style="width: 200px" />
          <Button type="primary" @click="handleSearch">搜索</Button>
        </Space>

        <Table
          :columns="emailColumns"
          :data-source="emailData"
          :loading="loading"
          :pagination="{ ...pagination, total }"
          row-key="id"
          @change="handleTableChange"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'status'">
              <Tag :color="record.status === 1 ? 'green' : 'red'">
                {{ record.status === 1 ? '成功' : '失败' }}
              </Tag>
            </template>
            <template v-if="column.dataIndex === 'created_at'">
              {{ dayjs(record.created_at).format('YYYY-MM-DD HH:mm:ss') }}
            </template>
          </template>
        </Table>
      </Tabs.TabPane>

      <Tabs.TabPane key="sms" tab="短信日志">
        <Space class="mb-4">
          <Input v-model:value="keyword" placeholder="搜索手机号" style="width: 200px" />
          <Button type="primary" @click="handleSearch">搜索</Button>
        </Space>

        <Table
          :columns="smsColumns"
          :data-source="smsData"
          :loading="loading"
          :pagination="{ ...pagination, total }"
          row-key="id"
          @change="handleTableChange"
        >
          <template #bodyCell="{ column, record }">
            <template v-if="column.dataIndex === 'status'">
              <Tag :color="record.status === 1 ? 'green' : 'red'">
                {{ record.status === 1 ? '成功' : '失败' }}
              </Tag>
            </template>
            <template v-if="column.dataIndex === 'created_at'">
              {{ dayjs(record.created_at).format('YYYY-MM-DD HH:mm:ss') }}
            </template>
          </template>
        </Table>
      </Tabs.TabPane>

      <Tabs.TabPane key="settings" tab="日志设置">
        <Form layout="vertical" style="max-width: 500px;">
          <FormItem label="操作日志">
            <Switch
              :checked="logConfig.log_operation_enabled === '1'"
              checked-children="开启"
              un-checked-children="关闭"
              @change="(v: any) => logConfig.log_operation_enabled = v ? '1' : '0'"
            />
            <span style="margin-left: 12px; opacity: 0.6;">记录用户的增删改查操作</span>
          </FormItem>
          <FormItem label="登录日志">
            <Switch
              :checked="logConfig.log_login_enabled === '1'"
              checked-children="开启"
              un-checked-children="关闭"
              @change="(v: any) => logConfig.log_login_enabled = v ? '1' : '0'"
            />
            <span style="margin-left: 12px; opacity: 0.6;">记录用户登录成功/失败</span>
          </FormItem>
          <FormItem label="邮件日志">
            <Switch
              :checked="logConfig.log_email_enabled === '1'"
              checked-children="开启"
              un-checked-children="关闭"
              @change="(v: any) => logConfig.log_email_enabled = v ? '1' : '0'"
            />
            <span style="margin-left: 12px; opacity: 0.6;">记录邮件发送记录</span>
          </FormItem>
          <FormItem label="短信日志">
            <Switch
              :checked="logConfig.log_sms_enabled === '1'"
              checked-children="开启"
              un-checked-children="关闭"
              @change="(v: any) => logConfig.log_sms_enabled = v ? '1' : '0'"
            />
            <span style="margin-left: 12px; opacity: 0.6;">记录短信发送记录</span>
          </FormItem>
          <FormItem>
            <Button type="primary" :loading="configLoading" @click="handleSaveLogConfig">保存设置</Button>
          </FormItem>
        </Form>
      </Tabs.TabPane>
    </Tabs>

    <!-- 详情弹窗 -->
    <Modal
      v-model:open="detailVisible"
      title="日志详情"
      width="800px"
      :footer="null"
    >
      <Descriptions v-if="detailRecord" :column="2" bordered>
        <Descriptions.Item label="ID">{{ detailRecord.id }}</Descriptions.Item>
        <Descriptions.Item label="用户">{{ detailRecord.username }}</Descriptions.Item>
        <Descriptions.Item label="IP">{{ detailRecord.ip }}</Descriptions.Item>
        <Descriptions.Item label="时间">{{ dayjs(detailRecord.created_at).format('YYYY-MM-DD HH:mm:ss') }}</Descriptions.Item>
        
        <!-- 操作日志特有字段 -->
        <template v-if="activeTab === 'operation'">
          <Descriptions.Item label="模块">{{ (detailRecord as OperationLogRecord).module }}</Descriptions.Item>
          <Descriptions.Item label="操作">{{ (detailRecord as OperationLogRecord).action }}</Descriptions.Item>
          <Descriptions.Item label="方法">{{ (detailRecord as OperationLogRecord).method }}</Descriptions.Item>
          <Descriptions.Item label="路径">{{ (detailRecord as OperationLogRecord).path }}</Descriptions.Item>
          <Descriptions.Item label="耗时">{{ (detailRecord as OperationLogRecord).duration }}ms</Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag :color="(detailRecord as OperationLogRecord).status === 1 ? 'green' : 'red'">
              {{ (detailRecord as OperationLogRecord).status === 1 ? '成功' : '失败' }}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item label="请求参数" :span="2">
            <pre>{{ (detailRecord as OperationLogRecord).params || '-' }}</pre>
          </Descriptions.Item>
          <Descriptions.Item label="响应结果" :span="2">
            <pre>{{ (detailRecord as OperationLogRecord).result || '-' }}</pre>
          </Descriptions.Item>
          <Descriptions.Item v-if="(detailRecord as OperationLogRecord).error" label="错误信息" :span="2">
            <pre style="color: red">{{ (detailRecord as OperationLogRecord).error }}</pre>
          </Descriptions.Item>
        </template>
        
        <!-- 登录日志特有字段 -->
        <template v-else>
          <Descriptions.Item label="登录类型">
            <Tag>{{ (detailRecord as LoginLogRecord).login_type === 'password' ? '密码登录' : 'TOTP登录' }}</Tag>
          </Descriptions.Item>
          <Descriptions.Item label="状态">
            <Tag :color="(detailRecord as LoginLogRecord).status === 1 ? 'green' : 'red'">
              {{ (detailRecord as LoginLogRecord).status === 1 ? '成功' : '失败' }}
            </Tag>
          </Descriptions.Item>
          <Descriptions.Item v-if="(detailRecord as LoginLogRecord).error" label="错误信息" :span="2">
            <pre style="color: red">{{ (detailRecord as LoginLogRecord).error }}</pre>
          </Descriptions.Item>
        </template>
        
        <Descriptions.Item label="User Agent" :span="2">
          <div style="word-break: break-all">{{ detailRecord.user_agent }}</div>
        </Descriptions.Item>
      </Descriptions>
    </Modal>
  </Card>
</template>
