<template>
  <div class="projection-wrapper">
    <zoompinch
      v-show="props.isOpen && props.fileUrl"
      ref="zoompinchRef"
      :width="1536"
      :height="2048"
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
}>();

const zoompinchRef = ref<InstanceType<typeof Zoompinch>>();
const transform = ref({ x: 0, y: 0, scale: 1, rotate: 0 });

const onImgDidLoad = () => {
  if (zoompinchRef.value) {
    transform.value = { x: 0, y: 0, scale: 1, rotate: 0 };
  }
};
</script>

<style scoped>
.projection-wrapper {
  width: 100%;
  height: 100%;
}

ion-content ion-button[slot='fixed'] {
  top: var(--ion-padding);
  right: var(--ion-padding);
}
</style>
