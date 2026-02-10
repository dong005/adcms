<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';

import { computed, markRaw, ref } from 'vue';

import { AuthenticationLogin, SliderCaptcha, z } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { useAuthStore } from '#/store';
import TotpLogin from './totp-login.vue';

defineOptions({ name: 'Login' });

const authStore = useAuthStore();
const showTotp = ref(false);

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: $t('authentication.usernameTip'),
      },
      fieldName: 'username',
      label: $t('authentication.username'),
      rules: z.string().min(1, { message: $t('authentication.usernameTip') }),
    },
    {
      component: 'VbenInputPassword',
      componentProps: {
        placeholder: $t('authentication.password'),
      },
      fieldName: 'password',
      label: $t('authentication.password'),
      rules: z.string().min(1, { message: $t('authentication.passwordTip') }),
    },
    {
      component: markRaw(SliderCaptcha),
      fieldName: 'captcha',
      rules: z.boolean().refine((value) => value, {
        message: $t('authentication.verifyRequiredTip'),
      }),
    },
  ];
});

/**
 * 处理登录提交
 */
async function handleLogin(values: any) {
  const result = await authStore.authLogin({
    username: values.username,
    password: values.password,
  });
  
  // 如果需要 TOTP 验证
  if (result.require_totp) {
    showTotp.value = true;
  }
}

/**
 * TOTP 验证成功回调
 */
function onTotpSuccess() {
  showTotp.value = false;
}
</script>

<template>
  <AuthenticationLogin
    v-if="!showTotp"
    :form-schema="formSchema"
    :loading="authStore.loginLoading"
    @submit="handleLogin"
  />
  <TotpLogin
    v-else
    :loading="authStore.loginLoading"
    @success="onTotpSuccess"
  />
</template>
