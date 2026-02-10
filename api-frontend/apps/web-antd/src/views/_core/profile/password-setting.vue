<script setup lang="ts">
import type { VbenFormSchema } from '#/adapter/form';

import { ref } from 'vue';

import { ProfilePasswordSetting, z } from '@vben/common-ui';

import { message } from 'ant-design-vue';

import { changePasswordApi } from '#/api';

const loading = ref(false);

const formSchema = ref([
  {
    fieldName: 'old_password',
    label: '旧密码',
    component: 'VbenInputPassword',
    componentProps: {
      placeholder: '请输入旧密码',
    },
    rules: z.string().min(1, { message: '请输入旧密码' }),
  },
  {
    fieldName: 'new_password',
    label: '新密码',
    component: 'VbenInputPassword',
    componentProps: {
      passwordStrength: true,
      placeholder: '请输入新密码',
    },
    rules: z.string().min(6, { message: '新密码至少6位' }),
  },
  {
    fieldName: 'confirmPassword',
    label: '确认密码',
    component: 'VbenInputPassword',
    componentProps: {
      passwordStrength: true,
      placeholder: '请再次输入新密码',
    },
    dependencies: {
      rules(values) {
        const { new_password } = values;
        return z
          .string({ required_error: '请再次输入新密码' })
          .min(1, { message: '请再次输入新密码' })
          .refine((value) => value === new_password, {
            message: '两次输入的密码不一致',
          });
      },
      triggerFields: ['new_password'],
    },
  },
]);

async function handleSubmit(values: Record<string, any>) {
  try {
    loading.value = true;
    const { old_password, new_password } = values;
    await changePasswordApi({ old_password, new_password });
    message.success('密码修改成功');
  } catch (error) {
    message.error('密码修改失败');
  } finally {
    loading.value = false;
  }
}
</script>
<template>
  <ProfilePasswordSetting
    class="w-1/3"
    :form-schema="formSchema"
    :loading="loading"
    @submit="handleSubmit"
  />
</template>
