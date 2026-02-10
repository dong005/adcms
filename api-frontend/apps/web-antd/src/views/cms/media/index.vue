<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Space, Card, Upload, Image, Input, Select, SelectOption,
  Popconfirm, message,
} from 'ant-design-vue';
import { getMediaList, uploadMedia, deleteMedia } from '#/api/cms';
import type { MediaRecord } from '#/api/cms';
import type { UploadProps } from 'ant-design-vue';

const loading = ref(false);
const uploading = ref(false);
const dataSource = ref<MediaRecord[]>([]);
const total = ref(0);
const keyword = ref('');
const typeFilter = ref<string>('');
const pagination = reactive({ current: 1, pageSize: 20 });

const previewVisible = ref(false);
const previewImage = ref('');

const columns = [
  { title: '预览', dataIndex: 'preview', width: 100 },
  { title: '文件名', dataIndex: 'original_name', width: 200 },
  { title: '类型', dataIndex: 'mime_type', width: 150 },
  { title: '大小', dataIndex: 'size', width: 100 },
  { title: '上传时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', dataIndex: 'action', width: 150, fixed: 'right' as const },
];

async function fetchData() {
  loading.value = true;
  try {
    const res = await getMediaList({
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: keyword.value,
      type: typeFilter.value,
    });
    dataSource.value = res.list || [];
    total.value = res.total;
  } finally {
    loading.value = false;
  }
}

const uploadProps: UploadProps = {
  name: 'file',
  multiple: true,
  showUploadList: false,
  customRequest: async ({ file, onSuccess, onError }) => {
    try {
      uploading.value = true;
      await uploadMedia(file as File);
      message.success('上传成功');
      onSuccess?.(file);
      fetchData();
    } catch (error) {
      message.error('上传失败');
      onError?.(error as Error);
    } finally {
      uploading.value = false;
    }
  },
};

async function handleDelete(id: number) {
  try {
    await deleteMedia(id);
    message.success('删除成功');
    fetchData();
  } catch (error) {
    message.error('删除失败');
  }
}

function handlePreview(record: MediaRecord) {
  if (record.mime_type.startsWith('image/')) {
    previewImage.value = record.path;
    previewVisible.value = true;
  } else {
    window.open(record.path, '_blank');
  }
}

function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B';
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(2) + ' KB';
  return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
}

function handleTableChange(pag: any) {
  pagination.current = pag.current;
  pagination.pageSize = pag.pageSize;
  fetchData();
}

onMounted(() => {
  fetchData();
});
</script>

<template>
  <Card>
    <Space class="mb-4">
      <Input v-model:value="keyword" placeholder="搜索文件名" style="width: 200px" />
      <Select v-model:value="typeFilter" placeholder="文件类型" allow-clear style="width: 150px">
        <SelectOption value="image">图片</SelectOption>
        <SelectOption value="video">视频</SelectOption>
        <SelectOption value="audio">音频</SelectOption>
        <SelectOption value="document">文档</SelectOption>
      </Select>
      <Button type="primary" @click="fetchData">搜索</Button>
      <Upload v-bind="uploadProps">
        <Button type="primary" :loading="uploading">
          上传文件
        </Button>
      </Upload>
    </Space>

    <Table
      :columns="columns"
      :data-source="dataSource"
      :loading="loading"
      :pagination="{ ...pagination, total, showSizeChanger: true, showQuickJumper: true }"
      row-key="id"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'preview'">
          <div class="preview-container">
            <Image
              v-if="record.mime_type.startsWith('image/')"
              :width="60"
              :height="60"
              :src="record.path"
              fit="cover"
              style="object-fit: cover; border-radius: 4px;"
            />
            <div v-else class="file-icon">
              {{ record.mime_type.split('/')[0]?.toUpperCase() || 'FILE' }}
            </div>
          </div>
        </template>
        <template v-if="column.dataIndex === 'size'">
          {{ formatSize(record.size) }}
        </template>
        <template v-if="column.dataIndex === 'created_at'">
          {{ new Date(record.created_at).toLocaleString() }}
        </template>
        <template v-if="column.dataIndex === 'action'">
          <Space>
            <Button size="small" type="link" @click="handlePreview(record as MediaRecord)">
              预览
            </Button>
            <Popconfirm title="确定删除这个文件吗？" @confirm="handleDelete(record.id)">
              <Button size="small" type="link" danger>
                删除
              </Button>
            </Popconfirm>
          </Space>
        </template>
      </template>
    </Table>

    <Image
      :style="{ display: 'none' }"
      :preview="{ visible: previewVisible, onVisibleChange: (vis) => previewVisible = vis }"
      :src="previewImage"
    />
  </Card>
</template>

<style scoped>
.preview-container {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 60px;
  height: 60px;
  background: #f5f5f5;
  border-radius: 4px;
}

.file-icon {
  font-size: 12px;
  color: #666;
  text-align: center;
  line-height: 60px;
}
</style>
