<script setup lang="ts">
import { ref } from 'vue';
import { Modal, Form, FormItem, Input, Button, QRCode, message } from 'ant-design-vue';
import { generateTOTPApi, bindTOTPApi, unbindTOTPApi } from '#/api';

const props = defineProps<{
  visible: boolean;
  isBound: boolean;
}>();

const emit = defineEmits<{
  'update:visible': [value: boolean];
  'success': [];
}>();

const loading = ref(false);
const qrCode = ref('');
const secret = ref('');
const bindCode = ref('');
const disableCode = ref('');

// 生成二维码
const generateQR = async () => {
  try {
    loading.value = true;
    const res = await generateTOTPApi();
    qrCode.value = res.qr_code;
    secret.value = res.secret;
  } catch (error: any) {
    message.error(error?.message || '生成二维码失败');
  } finally {
    loading.value = false;
  }
};

// 绑定 TOTP
const handleBind = async () => {
  if (!bindCode.value) {
    message.error('请输入验证码');
    return;
  }
  try {
    loading.value = true;
    await bindTOTPApi({ code: bindCode.value });
    message.success('TOTP 绑定成功');
    emit('success');
    handleClose();
  } catch (error: any) {
    message.error(error?.message || '绑定失败，请检查验证码');
  } finally {
    loading.value = false;
  }
};

// 解绑 TOTP
const handleDisable = async () => {
  if (!disableCode.value) {
    message.error('请输入验证码');
    return;
  }
  try {
    loading.value = true;
    await unbindTOTPApi({ code: disableCode.value });
    message.success('TOTP 解绑成功');
    emit('success');
    handleClose();
  } catch (error: any) {
    message.error(error?.message || '解绑失败，请检查验证码');
  } finally {
    loading.value = false;
  }
};

const handleClose = () => {
  emit('update:visible', false);
  qrCode.value = '';
  secret.value = '';
  bindCode.value = '';
  disableCode.value = '';
};
</script>

<template>
  <Modal
    :open="props.visible"
    :title="props.isBound ? '解绑 TOTP' : '绑定 TOTP'"
    :footer="null"
    width="400px"
    @cancel="handleClose"
  >
    <!-- 绑定 TOTP -->
    <div v-if="!isBound">
      <div v-if="!qrCode" class="text-center">
        <p>请点击下方按钮生成二维码</p>
        <Button type="primary" :loading="loading" @click="generateQR">
          生成二维码
        </Button>
      </div>
      <div v-else>
        <div class="text-center mb-4">
          <p class="mb-2">请使用谷歌验证器扫描下方二维码</p>
          <QRCode :value="qrCode" :size="200" />
          <p class="mt-2 text-gray-500">或手动输入密钥：{{ secret }}</p>
        </div>
        <Form @submit="handleBind">
          <FormItem label="验证码">
            <Input
              v-model:value="bindCode"
              placeholder="请输入6位验证码"
              max-length="6"
            />
          </FormItem>
          <FormItem class="mb-0">
            <Button type="primary" html-type="submit" :loading="loading" block>
              确认绑定
            </Button>
          </FormItem>
        </Form>
      </div>
    </div>

    <!-- 解绑 TOTP -->
    <div v-else>
      <p class="mb-4">解绑 TOTP 将降低账户安全性，请确认操作</p>
      <Form @submit="handleDisable">
        <FormItem label="验证码">
          <Input
            v-model:value="disableCode"
            placeholder="请输入6位验证码"
            max-length="6"
          />
        </FormItem>
        <FormItem class="mb-0">
          <Button type="primary" html-type="submit" :loading="loading" block>
            确认解绑
          </Button>
        </FormItem>
      </Form>
    </div>
  </Modal>
</Template>
