<template>
  <div
    ref="wrapperRef"
    class="projection-wrapper"
    @touchstart="onTouchStart"
    @touchmove="onTouchMove"
    @touchend="onTouchEnd"
    @touchcancel="onTouchEnd"
  >
    <img
      v-show="props.isOpen && props.fileUrl"
      class="zoom-image"
      :src="props.fileUrl"
      :style="imageStyle"
      :draggable="false"
      @load="onImgLoad"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue';

const props = defineProps<{
  isOpen: boolean;
  fileUrl: string;
}>();

const emit = defineEmits<{
  (e: 'ready'): void;
}>();

const MIN_SCALE = 1;
const MAX_SCALE = 10;

const wrapperRef = ref<HTMLDivElement | null>(null);
const scale = ref(MIN_SCALE);
const translate = ref({ x: 0, y: 0 });
const isPinching = ref(false);

const imageStyle = computed(() => ({
  transform: `translate(${translate.value.x}px, ${translate.value.y}px) scale(${scale.value})`,
}));

const initialScale = computed(() => MIN_SCALE);
const currentScale = computed(() => scale.value);

const reset = () => {
  scale.value = MIN_SCALE;
  translate.value = { x: 0, y: 0 };
};

const onImgLoad = () => {
  reset();
  emit('ready');
};

watch(
  () => props.fileUrl,
  (newUrl) => {
    if (newUrl) reset();
  },
);

const clamp = (v: number, min: number, max: number) => Math.min(Math.max(v, min), max);
const dist = (t1: Touch, t2: Touch) => Math.hypot(t2.clientX - t1.clientX, t2.clientY - t1.clientY);

const clampTranslate = (tx: number, ty: number, s: number) => {
  const rect = wrapperRef.value?.getBoundingClientRect();
  if (!rect) return { x: tx, y: ty };
  const maxX = (rect.width * (s - 1)) / 2;
  const maxY = (rect.height * (s - 1)) / 2;
  return { x: clamp(tx, -maxX, maxX), y: clamp(ty, -maxY, maxY) };
};

let pinchStart = { distance: 0, scale: 1, cx: 0, cy: 0, tx: 0, ty: 0 };
let panStart = { cx: 0, cy: 0, tx: 0, ty: 0 };

const onTouchStart = (e: TouchEvent) => {
  const rect = wrapperRef.value?.getBoundingClientRect();
  if (!rect) return;

  if (e.touches.length === 2) {
    const [t1, t2] = [e.touches[0], e.touches[1]];
    pinchStart = {
      distance: dist(t1, t2),
      scale: scale.value,
      cx: (t1.clientX + t2.clientX) / 2 - rect.left,
      cy: (t1.clientY + t2.clientY) / 2 - rect.top,
      tx: translate.value.x,
      ty: translate.value.y,
    };
    isPinching.value = true;
  } else if (e.touches.length === 1 && scale.value > MIN_SCALE) {
    panStart = {
      cx: e.touches[0].clientX,
      cy: e.touches[0].clientY,
      tx: translate.value.x,
      ty: translate.value.y,
    };
  }
};

const onTouchMove = (e: TouchEvent) => {
  const rect = wrapperRef.value?.getBoundingClientRect();
  if (!rect) return;

  if (e.touches.length === 2) {
    const [t1, t2] = [e.touches[0], e.touches[1]];
    const currDistance = dist(t1, t2);
    const factor = currDistance / pinchStart.distance;
    const newScale = clamp(pinchStart.scale * factor, MIN_SCALE, MAX_SCALE);

    // Keep the image point under the pinch center fixed:
    //   startCenter = wrapperCenter + startTranslate + imgPoint * startScale
    //   → imgPoint  = (startCenter - wrapperCenter - startTranslate) / startScale
    //   currCenter  = wrapperCenter + newTranslate + imgPoint * newScale
    //   → newTranslate = currCenter - wrapperCenter - imgPoint * newScale
    const wcx = rect.width / 2;
    const wcy = rect.height / 2;
    const currCX = (t1.clientX + t2.clientX) / 2 - rect.left;
    const currCY = (t1.clientY + t2.clientY) / 2 - rect.top;
    const imgPtX = (pinchStart.cx - wcx - pinchStart.tx) / pinchStart.scale;
    const imgPtY = (pinchStart.cy - wcy - pinchStart.ty) / pinchStart.scale;
    const tNewX = currCX - wcx - imgPtX * newScale;
    const tNewY = currCY - wcy - imgPtY * newScale;

    scale.value = newScale;
    translate.value = clampTranslate(tNewX, tNewY, newScale);
  } else if (e.touches.length === 1 && scale.value > MIN_SCALE) {
    const dx = e.touches[0].clientX - panStart.cx;
    const dy = e.touches[0].clientY - panStart.cy;
    translate.value = clampTranslate(panStart.tx + dx, panStart.ty + dy, scale.value);
  }
};

const onTouchEnd = (e: TouchEvent) => {
  if (e.touches.length < 2) isPinching.value = false;
  if (scale.value <= MIN_SCALE) {
    scale.value = MIN_SCALE;
    translate.value = { x: 0, y: 0 };
  }
};

defineExpose({
  reset,
  initialScale,
  currentScale,
  isPinching,
});
</script>

<style scoped>
.projection-wrapper {
  width: 100%;
  height: 100%;
  overflow: hidden;
  position: relative;
  touch-action: none;
}

.zoom-image {
  display: block;
  width: 100%;
  height: 100%;
  object-fit: contain;
  transform-origin: center center;
  user-select: none;
  -webkit-user-select: none;
  -webkit-user-drag: none;
  -webkit-touch-callout: none;
  will-change: transform;
}
</style>
