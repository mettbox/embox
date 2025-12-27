<template>
  <div
    class="media-grid-item"
    :data-id="media.id"
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
</template>

<script setup lang="ts">
import { IonImg, IonIcon, IonThumbnail } from '@ionic/vue';
import { checkmarkCircle, musicalNotes, videocam, heart } from 'ionicons/icons';
import { useThumbnail } from '@/composables/use-thumbnail';
import { useSelectedMediaStore } from '@/stores/selectedMedia.store';
import { useAppStore } from '@/stores/app.store';

const props = defineProps<{
  media: Media;
  isSelectMode: boolean;
}>();

const emit = defineEmits<{
  (e: 'media:open', media: Media): void;
  (e: 'media:select', media: Media): void;
}>();

const app = useAppStore();
const { loadThumbnail } = useThumbnail();
const selectedMedia = useSelectedMediaStore();
</script>

<style>
.media-grid-item {
  position: relative;
  aspect-ratio: 1 / 1;
  padding: 1px;

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

    pointer-events: none;
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
