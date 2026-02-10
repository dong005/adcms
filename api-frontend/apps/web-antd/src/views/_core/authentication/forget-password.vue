<script lang="ts" setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { forgotPasswordApi, resetPasswordByEmailApi } from '#/api/core/auth';

defineOptions({ name: 'ForgetPassword' });

const router = useRouter();
const step = ref(1);
const loading = ref(false);
const email = ref('');
const code = ref('');
const newPassword = ref('');
const confirmPassword = ref('');
const errorMsg = ref('');
const successMsg = ref('');
const countdown = ref(0);

let timer: ReturnType<typeof setInterval> | null = null;

function startCountdown() {
  countdown.value = 60;
  timer = setInterval(() => {
    countdown.value--;
    if (countdown.value <= 0 && timer) {
      clearInterval(timer);
      timer = null;
    }
  }, 1000);
}

async function handleSendCode() {
  if (!email.value || !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) {
    errorMsg.value = '请输入有效的邮箱地址';
    return;
  }
  loading.value = true;
  errorMsg.value = '';
  try {
    await forgotPasswordApi({ email: email.value });
    successMsg.value = '验证码已发送到您的邮箱，请查收';
    step.value = 2;
    startCountdown();
  } catch (e: any) {
    errorMsg.value = e?.response?.data?.message || '发送失败，请稍后重试';
  } finally {
    loading.value = false;
  }
}

async function handleReset() {
  errorMsg.value = '';
  if (!code.value || code.value.length !== 6) {
    errorMsg.value = '请输入6位验证码';
    return;
  }
  if (!newPassword.value || newPassword.value.length < 6) {
    errorMsg.value = '密码长度不能少于6位';
    return;
  }
  if (newPassword.value !== confirmPassword.value) {
    errorMsg.value = '两次输入的密码不一致';
    return;
  }
  loading.value = true;
  try {
    await resetPasswordByEmailApi({
      email: email.value,
      code: code.value,
      new_password: newPassword.value,
    });
    successMsg.value = '密码重置成功，即将跳转登录页...';
    errorMsg.value = '';
    setTimeout(() => router.push('/auth/login'), 2000);
  } catch (e: any) {
    errorMsg.value = e?.response?.data?.message || '重置失败，请检查验证码';
  } finally {
    loading.value = false;
  }
}

function goBack() {
  router.push('/auth/login');
}
</script>

<template>
  <div class="mx-auto w-full max-w-md px-6 py-10">
    <div class="mb-8 text-center">
      <h1 class="text-2xl font-bold">找回密码</h1>
      <p class="mt-2 text-sm text-gray-500">
        {{ step === 1 ? '输入您的注册邮箱，我们将发送验证码' : '输入验证码和新密码' }}
      </p>
    </div>

    <div v-if="errorMsg" class="mb-4 rounded bg-red-50 p-3 text-sm text-red-600">
      {{ errorMsg }}
    </div>
    <div v-if="successMsg" class="mb-4 rounded bg-green-50 p-3 text-sm text-green-600">
      {{ successMsg }}
    </div>

    <!-- Step 1: 输入邮箱 -->
    <div v-if="step === 1">
      <div class="mb-4">
        <label class="mb-1 block text-sm font-medium text-gray-700">邮箱地址</label>
        <input
          v-model="email"
          type="email"
          placeholder="请输入注册邮箱"
          class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          @keyup.enter="handleSendCode"
        />
      </div>
      <button
        :disabled="loading"
        class="w-full rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
        @click="handleSendCode"
      >
        {{ loading ? '发送中...' : '发送验证码' }}
      </button>
    </div>

    <!-- Step 2: 输入验证码 + 新密码 -->
    <div v-if="step === 2">
      <div class="mb-4">
        <label class="mb-1 block text-sm font-medium text-gray-700">邮箱</label>
        <input
          :value="email"
          disabled
          class="w-full rounded-md border border-gray-200 bg-gray-50 px-3 py-2 text-sm text-gray-500"
        />
      </div>
      <div class="mb-4">
        <label class="mb-1 block text-sm font-medium text-gray-700">验证码</label>
        <div class="flex gap-2">
          <input
            v-model="code"
            type="text"
            maxlength="6"
            placeholder="6位验证码"
            class="flex-1 rounded-md border border-gray-300 px-3 py-2 text-sm tracking-widest focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          />
          <button
            :disabled="countdown > 0"
            class="whitespace-nowrap rounded-md border border-gray-300 px-3 py-2 text-sm hover:bg-gray-50 disabled:opacity-50"
            @click="handleSendCode"
          >
            {{ countdown > 0 ? `${countdown}s` : '重新发送' }}
          </button>
        </div>
      </div>
      <div class="mb-4">
        <label class="mb-1 block text-sm font-medium text-gray-700">新密码</label>
        <input
          v-model="newPassword"
          type="password"
          placeholder="请输入新密码（至少6位）"
          class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
        />
      </div>
      <div class="mb-4">
        <label class="mb-1 block text-sm font-medium text-gray-700">确认密码</label>
        <input
          v-model="confirmPassword"
          type="password"
          placeholder="请再次输入新密码"
          class="w-full rounded-md border border-gray-300 px-3 py-2 text-sm focus:border-blue-500 focus:outline-none focus:ring-1 focus:ring-blue-500"
          @keyup.enter="handleReset"
        />
      </div>
      <button
        :disabled="loading"
        class="w-full rounded-md bg-blue-600 px-4 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50"
        @click="handleReset"
      >
        {{ loading ? '重置中...' : '重置密码' }}
      </button>
    </div>

    <div class="mt-6 text-center">
      <button class="text-sm text-blue-600 hover:underline" @click="goBack">
        返回登录
      </button>
    </div>
  </div>
</template>
