<script lang="ts" setup>
import { ref, reactive, onMounted } from 'vue';
import {
  Table, Button, Tag, Space, Modal, Form, FormItem, Input, Select,
  SelectOption, Popconfirm, message, Card, Upload,
} from 'ant-design-vue';
import type { UploadChangeParam } from 'ant-design-vue';
import { 
  getArticleList, createArticle, updateArticle, deleteArticle,
  publishArticle, unpublishArticle, getCategoryList, getTagList,
  createCategory, createTag, uploadMedia,
} from '#/api/cms';
import type { ArticleRecord, CategoryRecord, TagRecord } from '#/api/cms';
import dayjs from 'dayjs';
import RichTextEditor from '#/components/RichTextEditor.vue';

const loading = ref(false);
const dataSource = ref<ArticleRecord[]>([]);
const total = ref(0);
const keyword = ref('');
const statusFilter = ref<number | undefined>();
const categoryFilter = ref<number | undefined>();
const pagination = reactive({ current: 1, pageSize: 10 });

const modalVisible = ref(false);
const showCategoryInput = ref(false);
const newCategoryName = ref('');
const addingCategory = ref(false);
const showTagInput = ref(false);
const newTagName = ref('');
const addingTag = ref(false);
const uploading = ref(false);

async function handleUploadImage(info: UploadChangeParam) {
  const file = info.file.originFileObj || info.file;
  if (!file) return;
  uploading.value = true;
  try {
    const res = await uploadMedia(file as File);
    const path = (res as any).path || (res as any).url || '';
    formState.featured_image = path;
    message.success('上传成功');
  } catch {
    message.error('上传失败');
  } finally {
    uploading.value = false;
  }
}

function cancelAddCategory() {
  showCategoryInput.value = false;
  newCategoryName.value = '';
}

function cancelAddTag() {
  showTagInput.value = false;
  newTagName.value = '';
}

async function handleAddTag() {
  const name = newTagName.value.trim();
  if (!name) return;
  addingTag.value = true;
  try {
    await createTag({ name });
    message.success('标签创建成功');
    newTagName.value = '';
    showTagInput.value = false;
    await fetchTags();
  } catch {
    message.error('标签创建失败');
  } finally {
    addingTag.value = false;
  }
}

async function handleAddCategory() {
  const name = newCategoryName.value.trim();
  if (!name) return;
  addingCategory.value = true;
  try {
    await createCategory({ name, sort: 0, status: 1 });
    message.success('分类创建成功');
    newCategoryName.value = '';
    showCategoryInput.value = false;
    await fetchCategories();
  } catch {
    message.error('创建分类失败');
  } finally {
    addingCategory.value = false;
  }
}
const modalTitle = ref('新增文章');
const categoryList = ref<CategoryRecord[]>([]);
const tagList = ref<TagRecord[]>([]);
const selectedTagIds = ref<number[]>([]);

const formState = reactive({
  id: 0,
  title: '',
  slug: '',
  content: '',
  excerpt: '',
  featured_image: '',
  status: 0,
  category_id: undefined,
});

const columns = [
  { title: 'ID', dataIndex: 'id', width: 60 },
  { title: '标题', dataIndex: 'title', width: 200 },
  { title: '分类', dataIndex: 'category_id', width: 120 },
  { title: '状态', dataIndex: 'status', width: 80 },
  { title: '浏览量', dataIndex: 'view_count', width: 80 },
  { title: '作者', dataIndex: 'author_id', width: 100 },
  { title: '发布时间', dataIndex: 'published_at', width: 170 },
  { title: '创建时间', dataIndex: 'created_at', width: 170 },
  { title: '操作', dataIndex: 'action', width: 250, fixed: 'right' as const },
];

async function fetchData() {
  loading.value = true;
  try {
    const res = await getArticleList({
      page: pagination.current,
      page_size: pagination.pageSize,
      keyword: keyword.value,
      status: statusFilter.value,
      category_id: categoryFilter.value,
    });
    dataSource.value = res.list || [];
    total.value = res.total;
  } finally {
    loading.value = false;
  }
}

async function fetchCategories() {
  categoryList.value = await getCategoryList();
}

async function fetchTags() {
  const res = await getTagList({ page: 1, page_size: 100 });
  tagList.value = res.list || [];
}

function handleAdd() {
  modalTitle.value = '新增文章';
  Object.assign(formState, { 
    id: 0, 
    title: '', 
    slug: '', 
    content: '', 
    excerpt: '', 
    featured_image: '', 
    status: 0,
    category_id: undefined,
  });
  selectedTagIds.value = [];
  modalVisible.value = true;
}

function handleEdit(record: ArticleRecord) {
  modalTitle.value = '编辑文章';
  Object.assign(formState, {
    id: record.id,
    title: record.title,
    slug: record.slug,
    content: record.content,
    excerpt: record.excerpt,
    featured_image: record.featured_image,
    status: record.status,
    category_id: record.category_id,
  });
  selectedTagIds.value = record.tags?.map(tag => tag.id) || [];
  modalVisible.value = true;
}

async function handleSubmit() {
  try {
    const submitData = {
      ...formState,
      tag_ids: selectedTagIds.value,
    };
    if (formState.id) {
      await updateArticle(formState.id, submitData);
      message.success('更新成功');
    } else {
      await createArticle(submitData);
      message.success('创建成功');
    }
    modalVisible.value = false;
    fetchData();
  } catch (error) {
    message.error('操作失败');
  }
}

async function handleDelete(id: number) {
  try {
    await deleteArticle(id);
    message.success('删除成功');
    fetchData();
  } catch (error) {
    message.error('删除失败');
  }
}

async function handlePublish(id: number) {
  try {
    await publishArticle(id);
    message.success('发布成功');
    fetchData();
  } catch (error) {
    message.error('发布失败');
  }
}

async function handleUnpublish(id: number) {
  try {
    await unpublishArticle(id);
    message.success('下线成功');
    fetchData();
  } catch (error) {
    message.error('下线失败');
  }
}

function handleTableChange(pag: any) {
  pagination.current = pag.current;
  pagination.pageSize = pag.pageSize;
  fetchData();
}

onMounted(() => {
  fetchData();
  fetchCategories();
  fetchTags();
});
</script>

<template>
  <Card>
    <Space class="mb-4">
      <Input v-model:value="keyword" placeholder="搜索标题" style="width: 200px" />
      <Select v-model:value="statusFilter" placeholder="状态" allow-clear style="width: 120px">
        <SelectOption :value="0">草稿</SelectOption>
        <SelectOption :value="1">已发布</SelectOption>
        <SelectOption :value="2">已下线</SelectOption>
      </Select>
      <Select v-model:value="categoryFilter" placeholder="分类" allow-clear style="width: 150px">
        <SelectOption v-for="cat in categoryList" :key="cat.id" :value="cat.id">
          {{ cat.name }}
        </SelectOption>
      </Select>
      <Button type="primary" @click="fetchData">搜索</Button>
      <Button @click="handleAdd">新增文章</Button>
    </Space>

    <Table
      :columns="columns"
      :data-source="dataSource"
      :loading="loading"
      :pagination="{ ...pagination, total }"
      :scroll="{ x: 1200 }"
      row-key="id"
      @change="handleTableChange"
    >
      <template #bodyCell="{ column, record }">
        <template v-if="column.dataIndex === 'status'">
          <Tag :color="record.status === 1 ? 'green' : record.status === 2 ? 'red' : 'orange'">
            {{ record.status === 1 ? '已发布' : record.status === 2 ? '已下线' : '草稿' }}
          </Tag>
        </template>
        <template v-if="column.dataIndex === 'category_id'">
          {{ record.categories?.[0]?.name || '-' }}
        </template>
        <template v-if="column.dataIndex === 'published_at'">
          {{ record.published_at ? dayjs(record.published_at).format('YYYY-MM-DD HH:mm') : '-' }}
        </template>
        <template v-if="column.dataIndex === 'created_at'">
          {{ dayjs(record.created_at).format('YYYY-MM-DD HH:mm') }}
        </template>
        <template v-if="column.dataIndex === 'action'">
          <Space>
            <Button size="small" @click="handleEdit(record as ArticleRecord)">编辑</Button>
            <template v-if="record.status === 0">
              <Button size="small" type="primary" @click="handlePublish(record.id)">发布</Button>
            </template>
            <template v-else-if="record.status === 1">
              <Button size="small" type="default" @click="handleUnpublish(record.id)">下线</Button>
            </template>
            <Popconfirm title="确定删除这篇文章吗？" @confirm="handleDelete(record.id)">
              <Button size="small" danger>删除</Button>
            </Popconfirm>
          </Space>
        </template>
      </template>
    </Table>

    <Modal
      v-model:open="modalVisible"
      :title="modalTitle"
      width="800px"
      @ok="handleSubmit"
    >
      <Form ref="formRef" :model="formState" layout="vertical">
        <FormItem label="标题" name="title" :rules="[{ required: true, message: '请输入标题' }]">
          <Input v-model:value="formState.title" />
        </FormItem>
        <FormItem label="别名" name="slug">
          <Input v-model:value="formState.slug" />
        </FormItem>
        <FormItem label="分类">
          <div style="display: flex; gap: 8px; align-items: center">
            <FormItem name="category_id" no-style>
              <Select v-model:value="formState.category_id" placeholder="选择分类" style="flex: 1">
                <SelectOption v-for="cat in categoryList" :key="cat.id" :value="cat.id">
                  {{ cat.name }}
                </SelectOption>
              </Select>
            </FormItem>
            <Button v-if="!showCategoryInput" shape="circle" size="small" @click="showCategoryInput = true">+</Button>
            <template v-else>
              <Input
                v-model:value="newCategoryName"
                placeholder="分类名称"
                size="small"
                style="width: 120px"
                @keydown.enter.prevent="handleAddCategory"
              />
              <Button type="primary" size="small" :loading="addingCategory" @click="handleAddCategory">确定</Button>
              <Button size="small" @click="cancelAddCategory">取消</Button>
            </template>
          </div>
        </FormItem>
        <FormItem label="标签">
          <div style="display: flex; gap: 8px; align-items: center">
            <FormItem name="tags" no-style>
              <Select v-model:value="selectedTagIds" mode="multiple" placeholder="选择标签" style="flex: 1">
                <SelectOption v-for="tag in tagList" :key="tag.id" :value="tag.id">
                  {{ tag.name }}
                </SelectOption>
              </Select>
            </FormItem>
            <Button v-if="!showTagInput" shape="circle" size="small" @click="showTagInput = true">+</Button>
            <template v-else>
              <Input
                v-model:value="newTagName"
                placeholder="标签名称"
                size="small"
                style="width: 120px"
                @keydown.enter.prevent="handleAddTag"
              />
              <Button type="primary" size="small" :loading="addingTag" @click="handleAddTag">确定</Button>
              <Button size="small" @click="cancelAddTag">取消</Button>
            </template>
          </div>
        </FormItem>
        <FormItem label="摘要" name="excerpt">
          <Input.TextArea v-model:value="formState.excerpt" :rows="3" />
        </FormItem>
        <FormItem label="特色图片" name="featured_image">
          <Upload
            :show-upload-list="false"
            :before-upload="() => false"
            accept="image/*"
            @change="handleUploadImage"
          >
            <div class="featured-image-box">
              <template v-if="formState.featured_image">
                <img :src="formState.featured_image" alt="特色图片" />
                <span class="featured-image-delete" @click.stop="formState.featured_image = ''">×</span>
              </template>
              <template v-else>
                <span style="font-size: 24px; opacity: 0.4">+</span>
                <span style="font-size: 12px; opacity: 0.4">上传图片</span>
              </template>
            </div>
          </Upload>
        </FormItem>
        <FormItem label="状态" name="status">
          <Select v-model:value="formState.status">
            <SelectOption :value="0">草稿</SelectOption>
            <SelectOption :value="1">发布</SelectOption>
            <SelectOption :value="2">下线</SelectOption>
          </Select>
        </FormItem>
        <FormItem label="内容" name="content">
          <RichTextEditor v-model:value="formState.content" />
        </FormItem>
      </Form>
    </Modal>
  </Card>
</template>

<style scoped>
.featured-image-box {
  width: 104px;
  height: 104px;
  border: 1px dashed var(--border, #d9d9d9);
  border-radius: 6px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  position: relative;
  overflow: hidden;
}

.featured-image-box:hover {
  border-color: var(--ant-primary-color, #1677ff);
}

.featured-image-box img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.featured-image-delete {
  position: absolute;
  top: 2px;
  right: 2px;
  width: 20px;
  height: 20px;
  background: rgba(0, 0, 0, 0.45);
  color: #fff;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
  opacity: 0;
  transition: opacity 0.2s;
}

.featured-image-box:hover .featured-image-delete {
  opacity: 1;
}
</style>
