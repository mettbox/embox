<template>
  <div class="projection-wrapper">
    <zoompinch
      v-show="props.isOpen && props.fileUrl"
      ref="zoompinchRef"
      :width="naturalWidth || 1920"
      :height="naturalHeight || 1080"
      :offset="{ top: 0, right: 0, bottom: 0, left: 0 }"
      :min-scale="1"
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
        <img
          :src="props.fileUrl"
          :width="naturalWidth || undefined"
          :height="naturalHeight || undefined"
          :draggable="false"
          @load="onImgLoad"
        />
      </template>
    </zoompinch>
  </div>
</template>

<script setup lang="ts">
import { Zoompinch } from 'zoompinch';
import '@/theme/zoompinch.css';
import { ref, computed, watch } from 'vue';

const props = defineProps<{
  isOpen: boolean;
  fileUrl: string;
}>();

const emit = defineEmits<{
  (e: 'ready'): void;
}>();

// zoompinch's scale=1 means "fit canvas into view box, centered" as long as
// the canvas width/height match the image's aspect ratio.
const FIT_SCALE = 1;

const zoompinchRef = ref<InstanceType<typeof Zoompinch>>();
const transform = ref({ x: 0, y: 0, scale: FIT_SCALE, rotate: 0 });
const naturalWidth = ref(0);
const naturalHeight = ref(0);

const initialScale = computed(() => FIT_SCALE);
const currentScale = computed(() => transform.value.scale);

const reset = () => {
  transform.value = { x: 0, y: 0, scale: FIT_SCALE, rotate: 0 };
};

const applyDimensions = (w: number, h: number) => {
  if (!w || !h) return;
  naturalWidth.value = w;
  naturalHeight.value = h;
  reset();
  emit('ready');
};

const onImgLoad = (ev: Event) => {
  const img = ev.target as HTMLImageElement;
  applyDimensions(img.naturalWidth, img.naturalHeight);
};

// Re-fit when natural dimensions change.
watch([naturalWidth, naturalHeight], () => reset());

// When the source URL changes, try to resolve dimensions BEFORE the visible
// <img> re-renders. For cached images (e.g. preloaded thumbnails), this
// happens synchronously, so the canvas aspect ratio is correct on the very
// first paint. For uncached images, falls back to the @load event.
watch(
  () => props.fileUrl,
  (newUrl) => {
    if (!newUrl) return;
    const probe = new Image();
    probe.src = newUrl;
    if (probe.complete && probe.naturalWidth) {
      applyDimensions(probe.naturalWidth, probe.naturalHeight);
    } else {
      probe.onload = () => applyDimensions(probe.naturalWidth, probe.naturalHeight);
    }
  },
  { immediate: true },
);

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

.projection-wrapper :deep(.zoompinch .canvas img) {
  display: block;
  width: 100%;
  height: 100%;
}
</style>
