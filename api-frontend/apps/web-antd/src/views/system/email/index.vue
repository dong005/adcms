<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Card, Form, FormItem, Input, Button, Space, Modal, message,
} from 'ant-design-vue';
import { getEmailConfig, updateEmailConfig, testEmail } from '#/api/system/email';

defineOptions({ name: 'EmailConfig' });

const loading = ref(false);
const testLoading = ref(false);
const testModalVisible = ref(false);
const testTo = ref('');

const formState = reactive({
  smtp_host: '',
  smtp_port: '587',
  smtp_user: '',
  smtp_password: '',
  smtp_from: '',
});

async function fetchConfig() {
  loading.value = true;
  try {
    const data = await getEmailConfig();
    formState.smtp_host = data.smtp_host || '';
    formState.smtp_port = data.smtp_port || '587';
    formState.smtp_user = data.smtp_user || '';
    formState.smtp_password = data.smtp_password || '';
    formState.smtp_from = data.smtp_from || '';
  } catch {
    // ignore
  } finally {
    loading.value = false;
  }
}

async function handleSave() {
  loading.value = true;
  try {
    await updateEmailConfig({ ...formState });
    message.success('邮箱配置已保存');
  } catch {
    message.error('保存失败');
  } finally {
    loading.value = false;
  }
}

function showTestModal() {
  testTo.value = '';
  testModalVisible.value = true;
}

async function handleTestSend() {
  if (!testTo.value) {
    message.warning('请输入收件邮箱');
    return;
  }
  testLoading.value = true;
  try {
    await testEmail({ to: testTo.value });
    message.success('测试邮件已发送');
    testModalVisible.value = false;
  } catch (e: any) {
    message.error(e?.message || '发送失败');
  } finally {
    testLoading.value = false;
  }
}

onMounted(fetchConfig);
</script>

<template>
  <div class="p-4">
    <Card title="邮箱配置" :loading="loading">
      <template #extra>
        <Space>
          <Button @click="showTestModal">测试发送</Button>
          <Button type="primary" @click="handleSave">保存配置</Button>
        </Space>
      </template>

      <Form :model="formState" layout="vertical" style="max-width: 500px">
        <FormItem label="SMTP 服务器" :rules="[{ required: true }]">
          <Input v-model:value="formState.smtp_host" placeholder="如 smtp.qq.com / smtp.gmail.com" />
        </FormItem>
        <FormItem label="SMTP 端口">
          <Input v-model:value="formState.smtp_port" placeholder="587（TLS）或 465（SSL）" />
        </FormItem>
        <FormItem label="SMTP 账号" :rules="[{ required: true }]">
          <Input v-model:value="formState.smtp_user" placeholder="发件邮箱账号" />
        </FormItem>
        <FormItem label="SMTP 密码/授权码" :rules="[{ required: true }]">
          <Input.Password v-model:value="formState.smtp_password" placeholder="邮箱密码或授权码" />
        </FormItem>
        <FormItem label="发件人地址">
          <Input v-model:value="formState.smtp_from" placeholder="留空则使用 SMTP 账号" />
        </FormItem>
      </Form>
    </Card>

    <Modal
      v-model:open="testModalVisible"
      title="发送测试邮件"
      :confirm-loading="testLoading"
      @ok="handleTestSend"
    >
      <div class="mt-4">
        <p class="mb-2 text-sm text-gray-500">请先保存配置，再发送测试邮件。</p>
        <Input v-model:value="testTo" placeholder="收件人邮箱" />
      </div>
    </Modal>
  </div>
</template>
