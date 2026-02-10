<script setup lang="ts">
import type { BasicOption } from '@vben/types';

import type { VbenFormSchema } from '#/adapter/form';

import { onMounted, ref } from 'vue';

import { ProfileBaseSetting } from '@vben/common-ui';
import { message } from 'ant-design-vue';

import { getUserInfoApi, updateUserInfoApi } from '#/api';

const profileBaseSettingRef = ref();
const loading = ref(false);

const MOCK_ROLES_OPTIONS: BasicOption[] = [
  {
    label: '管理员',
    value: 'super',
  },
  {
    label: '用户',
    value: 'user',
  },
  {
    label: '测试',
    value: 'test',
  },
];

const formSchema = ref([
  {
    fieldName: 'realName',
    component: 'Input',
    label: '姓名',
  },
  {
    fieldName: 'username',
    component: 'Input',
    label: '用户名',
    disabled: true,
  },
  {
    fieldName: 'email',
    component: 'Input',
    label: '邮箱',
  },
  {
    fieldName: 'roles',
    component: 'Select',
    componentProps: {
      mode: 'tags',
      options: MOCK_ROLES_OPTIONS,
      disabled: true,
    },
    label: '角色',
  },
  {
    fieldName: 'introduction',
    component: 'Textarea',
    label: '个人简介',
  },
]);

// 加载用户信息
onMounted(async () => {
  const data = await getUserInfoApi();
  profileBaseSettingRef.value.getFormApi().setValues({
    realName: data.realName,
    username: data.username,
    email: (data as any).email || '',
    roles: data.roles || [],
    introduction: '',
  });
});

// 提交更新
const handleSubmit = async (values: Record<string, any>) => {
  try {
    loading.value = true;
    await updateUserInfoApi({
      nickname: values.realName,
      email: values.email,
    });
    message.success('保存成功');
  } catch (error: any) {
    message.error(error?.message || '保存失败');
  } finally {
    loading.value = false;
  }
};
</script>
<template>
  <ProfileBaseSetting 
    ref="profileBaseSettingRef" 
    :form-schema="formSchema"
    :loading="loading"
    @submit="handleSubmit"
  />
</template>
