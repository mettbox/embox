<template>
  <ion-button
    v-show="isSelectMode"
    id="select-popover"
    shape="round"
    size="small"
    :color="selectedMedia.ids.length === 0 ? 'dark' : 'primary'"
    :disabled="selectedMedia.ids.length === 0"
  >
    <ion-icon
      slot="end"
      :icon="ellipsisHorizontalCircleOutline"
    />
    {{ $t('Selected', { count: selectedMedia.ids.length }) }}
  </ion-button>

  <ion-button
    icon-only
    shape="round"
    :color="isSelectMode ? 'primary' : 'dark'"
    @click="$emit('toggle-select-mode')"
  >
    <ion-icon
      slot="icon-only"
      :icon="isSelectMode ? close : checkmark"
    />
  </ion-button>

  <media-select-popover
    trigger="select-popover"
    @favourite:set="$emit('favourite:set')"
    @favourite:unset="$emit('favourite:unset')"
    @album:set="$emit('album:select')"
    @album:unset="$emit('album:unselect')"
    @album:cover="$emit('album:cover')"
    @public:set="$emit('public:set')"
    @public:unset="$emit('public:unset')"
    @edit="$emit('edit')"
    @delete="$emit('delete')"
  />
</template>

<script setup lang="ts">
import { IonButton, IonIcon } from '@ionic/vue';
import { checkmark, close, ellipsisHorizontalCircleOutline } from 'ionicons/icons';
import mediaSelectPopover from './media-select-popover.vue';
import { useSelectedMediaStore } from '@/stores/selectedMedia.store';

defineProps<{
  isSelectMode: boolean;
}>();

defineEmits<{
  (e: 'toggle-select-mode'): void;
  (e: 'favourite:set'): void;
  (e: 'favourite:unset'): void;
  (e: 'album:select'): void;
  (e: 'album:unselect'): void;
  (e: 'album:cover'): void;
  (e: 'public:set'): void;
  (e: 'public:unset'): void;
  (e: 'edit'): void;
  (e: 'delete'): void;
}>();

const selectedMedia = useSelectedMediaStore();
</script>

<style scoped>
ion-button[disabled] {
  opacity: 1 !important;
}
</style>
