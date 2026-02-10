<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';
import type { Recordable } from '@vben/types';

import { computed, ref } from 'vue';

import { AuthenticationLogin, z } from '@vben/common-ui';

import { useAuthStore } from '#/store';

defineOptions({ name: 'TotpLogin' });

const authStore = useAuthStore();
const loading = ref(false);
const CODE_LENGTH = 6;

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      component: 'VbenPinInput',
      componentProps: {
        codeLength: CODE_LENGTH,
        placeholder: '请输入6位验证码',
      },
      fieldName: 'code',
      label: 'TOTP 验证码',
      rules: z.string().length(CODE_LENGTH, {
        message: '请输入6位验证码',
      }),
    },
  ];
});

/**
 * 异步处理 TOTP 验证
 * @param values TOTP 表单数据
 */
async function handleVerify(values: Recordable<any>) {
  const { code } = values;
  await authStore.verifyTotp(code);
}
</script>

<template>
  <AuthenticationLogin
    :form-schema="formSchema"
    :loading="loading"
    submit-text="验证"
    @submit="handleVerify"
  >
    <template #description>
      <div class="mb-4 text-center text-sm text-gray-600">
        请打开谷歌验证器，输入6位验证码
      </div>
    </template>
  </AuthenticationLogin>
</template>
