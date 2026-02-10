<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { Button, Modal, Input, Table, Tag, message } from 'ant-design-vue';
import {
  getUserInfoApi, updateUserInfoApi,
  sendSmsCodeApi, bindPhoneApi, getLoginHistoryApi,
} from '#/api';
import TotpModal from './totp-modal.vue';

// 用户信息
const userEmail = ref('');
const userPhone = ref('');
const isTotpBound = ref(false);
const lastLogin = ref('');

// TOTP 弹窗
const totpVisible = ref(false);

// 邮箱弹窗
const emailVisible = ref(false);
const emailInput = ref('');
const emailLoading = ref(false);

// 手机弹窗
const phoneVisible = ref(false);
const phoneInput = ref('');
const smsCode = ref('');
const phoneLoading = ref(false);
const smsSending = ref(false);
const smsCountdown = ref(0);
let smsTimer: any = null;

// 登录历史弹窗
const historyVisible = ref(false);
const historyLoading = ref(false);
const historyData = ref<any[]>([]);

const historyColumns = [
  { title: 'IP', dataIndex: 'ip', width: 140 },
  { title: '时间', dataIndex: 'created_at', width: 180 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '信息', dataIndex: 'message', ellipsis: true },
];

function maskEmail(email: string) {
  if (!email) return '';
  const parts = email.split('@');
  if (parts.length < 2 || !parts[0] || !parts[1]) return email;
  return parts[0].charAt(0) + '***@' + parts[1];
}

function maskPhone(phone: string) {
  if (!phone || phone.length < 7) return phone;
  return phone.substring(0, 3) + '****' + phone.substring(phone.length - 4);
}

async function fetchUserInfo() {
  try {
    const info = await getUserInfoApi();
    userEmail.value = (info as any).email || '';
    userPhone.value = (info as any).phone || '';
    isTotpBound.value = (info as any).totp_enabled || false;
  } catch {
    // ignore
  }
}

// 邮箱绑定
function openEmailModal() {
  emailInput.value = userEmail.value;
  emailVisible.value = true;
}

async function handleSaveEmail() {
  if (!emailInput.value.trim()) {
    message.warning('请输入邮箱');
    return;
  }
  emailLoading.value = true;
  try {
    await updateUserInfoApi({ email: emailInput.value } as any);
    message.success('邮箱已更新');
    userEmail.value = emailInput.value;
    emailVisible.value = false;
  } catch (e: any) {
    message.error(e?.message || '保存失败');
  } finally {
    emailLoading.value = false;
  }
}

// 手机绑定
function openPhoneModal() {
  phoneInput.value = '';
  smsCode.value = '';
  smsCountdown.value = 0;
  phoneVisible.value = true;
}

async function handleSendSms() {
  if (!phoneInput.value.trim()) {
    message.warning('请输入手机号');
    return;
  }
  smsSending.value = true;
  try {
    await sendSmsCodeApi({ phone: phoneInput.value });
    message.success('验证码已发送');
    smsCountdown.value = 60;
    smsTimer = setInterval(() => {
      smsCountdown.value--;
      if (smsCountdown.value <= 0) {
        clearInterval(smsTimer);
      }
    }, 1000);
  } catch (e: any) {
    message.error(e?.message || '发送失败');
  } finally {
    smsSending.value = false;
  }
}

async function handleBindPhone() {
  if (!phoneInput.value.trim() || !smsCode.value.trim()) {
    message.warning('请输入手机号和验证码');
    return;
  }
  phoneLoading.value = true;
  try {
    await bindPhoneApi({ phone: phoneInput.value, code: smsCode.value });
    message.success('手机绑定成功');
    userPhone.value = phoneInput.value;
    phoneVisible.value = false;
  } catch (e: any) {
    message.error(e?.message || '绑定失败');
  } finally {
    phoneLoading.value = false;
  }
}

// 登录历史
async function openHistory() {
  historyVisible.value = true;
  historyLoading.value = true;
  try {
    historyData.value = await getLoginHistoryApi();
    if (historyData.value.length > 0) {
      lastLogin.value = historyData.value[0].created_at;
    }
  } catch {
    historyData.value = [];
  } finally {
    historyLoading.value = false;
  }
}

onMounted(() => {
  fetchUserInfo();
});
</script>

<template>
  <div class="space-y-4" style="max-width: 700px;">
    <!-- 账户密码 -->
    <div class="flex items-center justify-between rounded-lg border p-4">
      <div class="flex items-center gap-3">
        <span style="font-size: 20px;">&#128274;</span>
        <div>
          <div class="text-base font-medium">账户密码</div>
          <div class="text-sm text-gray-500">当前密码强度：已设置</div>
        </div>
      </div>
      <Button size="small" @click="$router.push({ path: '/profile', query: { tab: 'password' } })">修改密码</Button>
    </div>

    <!-- 备用邮箱 -->
    <div class="flex items-center justify-between rounded-lg border p-4">
      <div class="flex items-center gap-3">
        <span style="font-size: 20px;">&#9993;</span>
        <div>
          <div class="text-base font-medium">备用邮箱</div>
          <div class="text-sm text-gray-500">
            {{ userEmail ? '已绑定 ' + maskEmail(userEmail) : '未绑定邮箱' }}
          </div>
        </div>
      </div>
      <Button size="small" @click="openEmailModal">{{ userEmail ? '修改邮箱' : '绑定邮箱' }}</Button>
    </div>

    <!-- 密保手机 -->
    <div class="flex items-center justify-between rounded-lg border p-4">
      <div class="flex items-center gap-3">
        <span style="font-size: 20px;">&#128241;</span>
        <div>
          <div class="text-base font-medium">密保手机</div>
          <div class="text-sm text-gray-500">
            {{ userPhone ? '已绑定 ' + maskPhone(userPhone) : '未绑定手机' }}
          </div>
        </div>
      </div>
      <Button size="small" @click="openPhoneModal">{{ userPhone ? '修改手机' : '绑定手机' }}</Button>
    </div>

    <!-- MFA 设备 -->
    <div class="flex items-center justify-between rounded-lg border p-4">
      <div class="flex items-center gap-3">
        <span style="font-size: 20px;">&#128272;</span>
        <div>
          <div class="text-base font-medium">MFA 设备</div>
          <div class="text-sm text-gray-500">
            {{ isTotpBound ? '已绑定 TOTP 验证器' : '未绑定 MFA 设备，绑定后可进行二次确认' }}
          </div>
        </div>
      </div>
      <Button size="small" @click="totpVisible = true">{{ isTotpBound ? '解绑' : '绑定' }}</Button>
    </div>

    <!-- 登录历史 -->
    <div class="flex items-center justify-between rounded-lg border p-4">
      <div class="flex items-center gap-3">
        <span style="font-size: 20px;">&#128337;</span>
        <div>
          <div class="text-base font-medium">登录历史</div>
          <div class="text-sm text-gray-500">
            {{ lastLogin ? '上次登录：' + lastLogin : '查看最近登录记录' }}
          </div>
        </div>
      </div>
      <Button size="small" @click="openHistory">查看</Button>
    </div>
  </div>

  <!-- TOTP 弹窗 -->
  <TotpModal
    v-model:visible="totpVisible"
    :is-bound="isTotpBound"
    @success="fetchUserInfo"
  />

  <!-- 邮箱弹窗 -->
  <Modal v-model:open="emailVisible" title="绑定/修改邮箱" :confirm-loading="emailLoading" @ok="handleSaveEmail" :width="400">
    <div class="mt-4">
      <Input v-model:value="emailInput" placeholder="请输入邮箱地址" />
    </div>
  </Modal>

  <!-- 手机绑定弹窗 -->
  <Modal v-model:open="phoneVisible" title="绑定手机" :footer="null" :width="400">
    <div class="mt-4 space-y-3">
      <Input v-model:value="phoneInput" placeholder="请输入手机号" />
      <div class="flex gap-2">
        <Input v-model:value="smsCode" placeholder="6位验证码" style="flex: 1;" />
        <Button
          :loading="smsSending"
          :disabled="smsCountdown > 0"
          @click="handleSendSms"
        >
          {{ smsCountdown > 0 ? `${smsCountdown}s` : '发送验证码' }}
        </Button>
      </div>
      <Button type="primary" block :loading="phoneLoading" @click="handleBindPhone">确认绑定</Button>
    </div>
  </Modal>

  <!-- 登录历史弹窗 -->
  <Modal v-model:open="historyVisible" title="登录历史" :footer="null" :width="650">
    <Table
      :columns="historyColumns"
      :data-source="historyData"
      :loading="historyLoading"
      :pagination="false"
      row-key="id"
      size="small"
      class="mt-4"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'status'">
          <Tag :color="record.status === 1 ? 'green' : 'red'">
            {{ record.status === 1 ? '成功' : '失败' }}
          </Tag>
        </template>
      </template>
    </Table>
  </Modal>
</template>
