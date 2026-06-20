<template>
  <div class="projection-wrapper">
    <zoompinch
      v-show="props.isOpen && props.fileUrl"
      ref="zoompinchRef"
      :width="props.naturalWidth || 1920"
      :height="props.naturalHeight || 1080"
      :offset="{ top: 56, right: 0, bottom: 56, left: 0 }"
      :min-scale="0.5"
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
import { ref } from 'vue';

const props = defineProps<{
  isOpen: boolean;
  fileUrl: string;
  naturalWidth?: number;
  naturalHeight?: number;
}>();

const emit = defineEmits<{
  (e: 'ready'): void;
}>();

const zoompinchRef = ref<InstanceType<typeof Zoompinch>>();
const transform = ref({ x: 0, y: 0, scale: 1, rotate: 0 });

const onImgDidLoad = () => {
  transform.value = { x: 0, y: 0, scale: 1, rotate: 0 };
  emit('ready');
};
</script>

<style scoped>
.projection-wrapper {
  width: 100%;
  height: 100%;
}
</style>
