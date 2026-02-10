<script lang="ts" setup>
import { ref, reactive, onMounted, computed } from 'vue';
import {
  Button, Form, FormItem, Input, InputNumber,
  Switch, message, Card, Tabs, TabPane,
} from 'ant-design-vue';
import { getConfigList, updateConfigs } from '#/api/system/config';
import type { ConfigRecord } from '#/api/system/config';

const loading = ref(false);
const saving = ref(false);
const allConfigs = ref<ConfigRecord[]>([]);
const activeTab = ref('basic');

// 按组分类配置
const configGroups = computed(() => {
  const groups: Record<string, ConfigRecord[]> = {};
  allConfigs.value.forEach(config => {
    const groupKey = config.group || 'default';
    if (!groups[groupKey]) {
      groups[groupKey] = [];
    }
    groups[groupKey].push(config);
  });
  // 按排序字段排序
  Object.keys(groups).forEach(group => {
    const groupArray = groups[group];
    if (groupArray) {
      groupArray.sort((a, b) => a.sort - b.sort);
    }
  });
  return groups;
});

// 表单数据
const formData = reactive<Record<string, any>>({});

// 配置组件映射
const configComponents: Record<string, any> = {
  string: Input,
  number: InputNumber,
  boolean: Switch,
  json: Input.TextArea,
};

async function fetchData() {
  loading.value = true;
  try {
    allConfigs.value = await getConfigList();
    // 初始化表单数据
    allConfigs.value.forEach(config => {
      if (config.type === 'boolean') {
        formData[config.key] = config.value === 'true';
      } else if (config.type === 'number') {
        formData[config.key] = Number(config.value);
      } else {
        formData[config.key] = config.value;
      }
    });
  } finally {
    loading.value = false;
  }
}

async function handleSave() {
  try {
    saving.value = true;
    const updateData: Record<string, any> = {};
    Object.keys(formData).forEach(key => {
      const config = allConfigs.value.find(c => c.key === key);
      if (config) {
        if (config.type === 'boolean') {
          updateData[key] = formData[key] ? 'true' : 'false';
        } else {
          updateData[key] = String(formData[key]);
        }
      }
    });
    await updateConfigs(updateData);
    message.success('保存成功');
  } catch (error) {
    message.error('保存失败');
  } finally {
    saving.value = false;
  }
}

onMounted(() => {
  fetchData();
});
</script>

<template>
  <Card>
    <div class="flex justify-between items-center mb-4">
      <h2>系统配置</h2>
      <Button type="primary" :loading="saving" @click="handleSave">
        保存配置
      </Button>
    </div>

    <Tabs v-model:activeKey="activeTab">
      <TabPane
        v-for="(configs, group) in configGroups"
        :key="group"
        :tab="group === 'basic' ? '基础设置' : group === 'site' ? '站点设置' : group"
      >
        <Form :model="formData" layout="vertical">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormItem
              v-for="config in configs"
              :key="config.key"
              :label="config.description || config.key"
              :name="config.key"
            >
              <component :is="configComponents[config.type] || configComponents.string" v-model="formData[config.key]" />
            </FormItem>
          </div>
        </Form>
      </TabPane>
    </Tabs>
  </Card>
</template>

<style scoped>
.grid {
  display: grid;
  gap: 1rem;
}

.grid-cols-1 {
  grid-template-columns: repeat(1, minmax(0, 1fr));
}

.md\:grid-cols-2 {
  @media (min-width: 768px) {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
