<template>
  <div class="projection-wrapper">
    <zoompinch
      v-show="props.isOpen && props.fileUrl"
      ref="zoompinchRef"
      :width="props.naturalWidth || 1920"
      :height="props.naturalHeight || 1080"
      :offset="{ top: 0, right: 0, bottom: 0, left: 0 }"
      :min-scale="initialScale"
      :max-scale="10"
      :rotation="false"
      :bounds="true"
      :mouse="true"
      :touch="true"
      :wheel="true"
      :gesture="true"
      v-model:transform="transform"
    >
      <template #canvas>
        <ion-img
          :src="props.fileUrl"
          @ion-img-did-load="onImgDidLoad"
        />
      </template>
    </zoompinch>
  </div>
</template>

<script setup lang="ts">
import { Zoompinch } from 'zoompinch';
import '@/theme/zoompinch.css';
import { IonImg } from '@ionic/vue';
import { ref, computed } from 'vue';

const props = defineProps<{
  isOpen: boolean;
  fileUrl: string;
  naturalWidth?: number;
  naturalHeight?: number;
  containerWidth?: number;
}>();

const emit = defineEmits<{
  (e: 'ready'): void;
}>();

const zoompinchRef = ref<InstanceType<typeof Zoompinch>>();
const transform = ref({ x: 0, y: 0, scale: 1, rotate: 0 });

const initialScale = computed(() => {
  if (!props.containerWidth || !props.naturalWidth) return 1;
  return props.containerWidth / props.naturalWidth;
});

const currentScale = computed(() => transform.value.scale);

const reset = () => {
  transform.value = { x: 0, y: 0, scale: initialScale.value, rotate: 0 };
};

const onImgDidLoad = () => {
  reset();
  emit('ready');
};

defineExpose({
  reset,
  initialScale,
  currentScale,
});
</script>

<style scoped>
.projection-wrapper {
  width: 100%;
  height: 100%;
}
</style>
