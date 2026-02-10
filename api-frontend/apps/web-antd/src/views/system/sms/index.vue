<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Card, Form, FormItem, Input, Button, Space, Modal, message,
} from 'ant-design-vue';
import { getSmsConfig, updateSmsConfig, testSms } from '#/api/system/sms';

defineOptions({ name: 'SmsConfig' });

const loading = ref(false);
const testLoading = ref(false);
const testModalVisible = ref(false);
const testPhone = ref('');

const formState = reactive({
  sms_secret_id: '',
  sms_secret_key: '',
  sms_app_id: '',
  sms_sign: '',
  sms_template_id: '',
});

async function fetchConfig() {
  loading.value = true;
  try {
    const data = await getSmsConfig();
    formState.sms_secret_id = data.sms_secret_id || '';
    formState.sms_secret_key = data.sms_secret_key || '';
    formState.sms_app_id = data.sms_app_id || '';
    formState.sms_sign = data.sms_sign || '';
    formState.sms_template_id = data.sms_template_id || '';
  } catch {
    // ignore
  } finally {
    loading.value = false;
  }
}

async function handleSave() {
  loading.value = true;
  try {
    await updateSmsConfig({ ...formState });
    message.success('短信配置已保存');
  } catch {
    message.error('保存失败');
  } finally {
    loading.value = false;
  }
}

function showTestModal() {
  testPhone.value = '';
  testModalVisible.value = true;
}

async function handleTestSend() {
  if (!testPhone.value) {
    message.warning('请输入手机号');
    return;
  }
  testLoading.value = true;
  try {
    await testSms({ phone: testPhone.value });
    message.success('测试短信已发送');
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
    <Card title="短信配置（腾讯云）" :loading="loading">
      <template #extra>
        <Space>
          <Button @click="showTestModal">测试发送</Button>
          <Button type="primary" @click="handleSave">保存配置</Button>
        </Space>
      </template>

      <Form :model="formState" layout="vertical" style="max-width: 500px">
        <FormItem label="SecretId" :rules="[{ required: true }]">
          <Input v-model:value="formState.sms_secret_id" placeholder="腾讯云 SecretId" />
        </FormItem>
        <FormItem label="SecretKey" :rules="[{ required: true }]">
          <Input.Password v-model:value="formState.sms_secret_key" placeholder="腾讯云 SecretKey" />
        </FormItem>
        <FormItem label="AppId" :rules="[{ required: true }]">
          <Input v-model:value="formState.sms_app_id" placeholder="短信应用 SDK AppId" />
        </FormItem>
        <FormItem label="短信签名" :rules="[{ required: true }]">
          <Input v-model:value="formState.sms_sign" placeholder="已审核通过的短信签名" />
        </FormItem>
        <FormItem label="模板ID" :rules="[{ required: true }]">
          <Input v-model:value="formState.sms_template_id" placeholder="验证码短信模板ID" />
        </FormItem>
      </Form>
    </Card>

    <Modal
      v-model:open="testModalVisible"
      title="发送测试短信"
      :confirm-loading="testLoading"
      @ok="handleTestSend"
    >
      <div class="mt-4">
        <p class="mb-2 text-sm text-gray-500">请先保存配置，再发送测试短信。将发送验证码 123456。</p>
        <Input v-model:value="testPhone" placeholder="手机号" />
      </div>
    </Modal>
  </div>
</template>
