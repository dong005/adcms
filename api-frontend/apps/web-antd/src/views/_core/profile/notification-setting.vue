<script setup lang="ts">
import { ref, onMounted } from 'vue';

import { ProfileNotificationSetting } from '@vben/common-ui';

import { getUserInfoApi, updateUserInfoApi } from '#/api';
import { message } from 'ant-design-vue';

const emailNotify = ref(true);

onMounted(async () => {
  const info = await getUserInfoApi();
  emailNotify.value = (info as any).email_notify !== 0;
  updateFormSchema();
});

function updateFormSchema() {
  const item = formSchema.value.find(i => i.fieldName === 'emailNotify');
  if (item) {
    item.value = emailNotify.value;
  }
}

async function handleToggleEmailNotify() {
  emailNotify.value = !emailNotify.value;
  try {
    const val: any = emailNotify.value ? 1 : 0;
    await updateUserInfoApi({ email_notify: val } as any);
    message.success(emailNotify.value ? '已开启邮件通知' : '已关闭邮件通知');
  } catch {
    emailNotify.value = !emailNotify.value;
    message.error('保存失败');
  }
  updateFormSchema();
}

const formSchema = ref([
  {
    value: true,
    fieldName: 'emailNotify',
    label: '邮件通知',
    description: '系统消息将同时发送到您的邮箱（需先在基本设置中配置邮箱）',
    action: handleToggleEmailNotify,
  },
  {
    value: true,
    fieldName: 'systemMessage',
    label: '系统消息',
    description: '系统消息将以站内信的形式通知',
  },
  {
    value: true,
    fieldName: 'todoTask',
    label: '待办任务',
    description: '待办任务将以站内信的形式通知',
  },
]);
</script>
<template>
  <ProfileNotificationSetting :form-schema="formSchema" />
</template>
