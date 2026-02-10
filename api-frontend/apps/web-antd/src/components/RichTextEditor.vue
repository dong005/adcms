<script lang="ts" setup>
import { watch, ref } from 'vue';
import { QuillEditor } from '@vueup/vue-quill';
import '@vueup/vue-quill/dist/vue-quill.snow.css';

const props = defineProps<{
  value?: string;
}>();

const emit = defineEmits<{
  (e: 'update:value', val: string): void;
}>();

const content = ref(props.value || '');

watch(() => props.value, (val) => {
  if (val !== content.value) {
    content.value = val || '';
  }
});

function handleUpdate(val: string) {
  content.value = val;
  emit('update:value', val);
}

const toolbarOptions = [
  [{ header: [1, 2, 3, 4, 5, 6, false] }],
  ['bold', 'italic', 'underline', 'strike'],
  [{ color: [] }, { background: [] }],
  [{ align: [] }],
  [{ list: 'ordered' }, { list: 'bullet' }],
  [{ indent: '-1' }, { indent: '+1' }],
  ['blockquote', 'code-block'],
  ['link', 'image', 'video'],
  ['clean'],
];
</script>

<template>
  <QuillEditor
    theme="snow"
    :content="content"
    content-type="html"
    :toolbar="toolbarOptions"
    style="min-height: 300px"
    @update:content="handleUpdate"
  />
</template>
