<template>
  <DynamicScroller
    class="media-grid"
    ref="dynamicScrollerRef"
    :items="rows"
    :min-item-size="96"
    :buffer="200"
    key-field="rowIndex"
    :emitUpdate="true"
    @update="onUpdate"
    v-slot="{ item, index, active }"
  >
    <DynamicScrollerItem
      :item="item"
      :active="active"
      :data-index="index"
      :size-dependencies="[item.rowIndex]"
    >
      <div class="row">
        <div
          v-for="media in item.items"
          :key="media.id"
          class="cell"
          :style="{ flex: `1 1 calc(100% / ${columns})` }"
          :data-id="media.id"
          @click="onClickMedia(media)"
        >
          <ion-thumbnail v-if="media.type === 'audio'">
            <ion-icon
              :icon="musicalNotes"
              size="large"
            />
          </ion-thumbnail>
          <ion-img
            v-else
            :src="media.thumbUrl || (loadThumbnail(media), '')"
          />
          <ion-icon
            v-if="media.isFavourite"
            :icon="heart"
            color="light"
            size="default"
            style="position: absolute; top: 4px; left: 6px; z-index: 2"
          />
          <ion-icon
            v-if="media.type === 'video'"
            :icon="videocam"
            color="light"
            size="default"
            style="position: absolute; bottom: 4px; left: 6px; z-index: 2"
          />
          <div
            v-if="isSelectMode && selectedMedia.ids.includes(media.id)"
            class="select-overlay"
          >
            <ion-icon
              :icon="checkmarkCircle"
              color="light"
              size="medium"
            />
          </div>
        </div>
      </div>
    </DynamicScrollerItem>
  </DynamicScroller>
</template>

<script setup lang="ts">
import { DynamicScroller, DynamicScrollerItem } from 'vue-virtual-scroller';
import { IonImg, IonIcon, IonThumbnail } from '@ionic/vue';
import { checkmarkCircle, musicalNotes, videocam, heart } from 'ionicons/icons';
import { computed, ref } from 'vue';
import { useThumbnail } from '@/composables/use-thumbnail';
import { useSelectedMediaStore } from '@/stores/selectedMedia.store';

const props = defineProps<{
  mediaList: Media[];
  columns: number;
  isSelectMode: boolean;
}>();

const emit = defineEmits<{
  (e: 'media:open', media: Media): void;
  (e: 'media:select', media: Media): void;
  (e: 'update:visible-range', range: string): void;
}>();

const rows = computed<
  {
    rowIndex: number;
    items: Media[];
  }[]
>(() => {
  const out: any[] = [];
  for (let i = 0; i < props.mediaList.length; i += props.columns) {
    out.push({
      rowIndex: i / props.columns,
      items: props.mediaList.slice(i, i + props.columns),
    });
  }
  return out;
});

const { loadThumbnail } = useThumbnail();
const selectedMedia = useSelectedMediaStore();

const onUpdate = (start: number, end: number) => calculateVisibleRange(start, end);

const calculateVisibleRange = (startIndex: number, endIndex: number) => {
  const visibleRows = rows.value.slice(startIndex, endIndex + 1);
  const visibleMedia: any[] = [];
  visibleRows.forEach((row) => visibleMedia.push(...row.items));
  const visibleDates = visibleMedia.map((media) => media.date).filter(Boolean);

  if (!visibleDates.length) {
    emit('update:visible-range', '');
    return;
  }

  visibleDates.sort();
  const firstDate = new Date(visibleDates[0]);
  const lastDate = new Date(visibleDates[visibleDates.length - 1]);

  const formatter = new Intl.DateTimeFormat('de-DE', {
    month: 'short',
    year: 'numeric',
  });

  const firstLabel = formatter.format(firstDate);
  const lastLabel = formatter.format(lastDate);

  const visibleRange =
    firstDate.getFullYear() === lastDate.getFullYear() && firstDate.getMonth() === lastDate.getMonth()
      ? firstLabel
      : `${firstLabel} – ${lastLabel}`;

  emit('update:visible-range', visibleRange);
};

const dynamicScrollerRef = ref<VueInstanceElement | null>(null);
// const longPressTimer = ref<number | null>(null);
// const longPressTriggered = ref(false);
// const touchStartY = ref(0);

const onClickMedia = (media: Media) => {
  if (props.isSelectMode) {
    const index = selectedMedia.ids.indexOf(media.id);
    index === -1 ? selectedMedia.add(media) : selectedMedia.remove(media.id);
  } else {
    selectedMedia.set(media);
    emit('media:open', media);
  }
};

// const onLongPress = (media: Media) => {
//   emit('media:select', media);
// };

// const onTouchStart = (event: TouchEvent | MouseEvent) => {
//   const target = (event.target as HTMLElement).closest('.cell');
//   if (!target) return;

//   const id = target.getAttribute('data-id');
//   const media = props.mediaList.find((m) => String(m.id) === id);
//   if (!media) return;

//   // Get Y coordinate depending on event type
//   let clientY = 0;
//   if ('touches' in event && event.touches.length > 0) {
//     clientY = event.touches[0].clientY;
//   } else if ('clientY' in event) {
//     clientY = event.clientY;
//   } else {
//     console.warn('Unknown event type, no Y coordinate found');
//     return;
//   }

//   touchStartY.value = clientY;
//   longPressTriggered.value = false;
//   longPressTimer.value = window.setTimeout(() => {
//     longPressTriggered.value = true;
//     onLongPress(media);
//   }, 900); // 1 second press for long press
// };

// const onTouchEnd = (event: TouchEvent | MouseEvent) => {
//   if (longPressTimer.value) clearTimeout(longPressTimer.value);

//   let clientY = 0;
//   if ('changedTouches' in event && event.changedTouches.length > 0) {
//     clientY = event.changedTouches[0].clientY;
//   } else if ('clientY' in event) {
//     clientY = event.clientY;
//   } else {
//     return;
//   }

//   if (Math.abs(clientY - touchStartY.value) > 10) return;
//   if (longPressTriggered.value) return;

//   const target = (event.target as HTMLElement).closest('.cell');
//   if (!target) return;
//   const id = target.getAttribute('data-id');
//   const media = props.mediaList.find((m) => String(m.id) === id);
//   if (media) onClickMedia(media);
// };

// onMounted(() => {
//   const el = dynamicScrollerRef.value?.$el;
//   if (!el) return;

//   el.addEventListener('touchstart', onTouchStart);
//   el.addEventListener('mousedown', onTouchStart);
//   el.addEventListener('touchend', onTouchEnd);
//   el.addEventListener('mouseup', onTouchEnd);
//   el.addEventListener('mouseleave', () => {
//     if (longPressTimer.value) clearTimeout(longPressTimer.value);
//   });
// });

// onBeforeUnmount(() => {
//   const el = dynamicScrollerRef.value?.$el;
//   if (!el) return;

//   el.removeEventListener('touchstart', onTouchStart);
//   el.removeEventListener('mousedown', onTouchStart);
//   el.removeEventListener('touchend', onTouchEnd);
//   el.removeEventListener('mouseup', onTouchEnd);
// });
</script>

<style>
.media-grid {
  height: 100%;

  .row {
    display: flex;
    width: 100%;
    max-height: calc(100vh - 56px);
  }

  .cell {
    position: relative;
    aspect-ratio: 1 / 1;
    padding: 1px;
    max-width: 384px;
    max-height: 384px;

    ion-img {
      position: absolute;
      width: 100%;
      height: 100%;
      object-fit: cover;
    }

    ion-thumbnail {
      width: 100%;
      height: 100%;
      display: flex;
      align-items: center;
      justify-content: center;
      background-color: var(--ion-color-medium-shade);
      color: var(--ion-color-light);
    }
  }

  .select-overlay {
    position: absolute;
    top: 0;
    left: 0;
    z-index: 2;
    background: rgba(0, 0, 0, 0.4);
    width: 100%;
    height: 100%;

    ion-icon {
      position: absolute;
      bottom: 4px;
      right: 4px;
      z-index: 2;
    }
  }
}
</style>
