<template>
  <div
    class="button-row"
    slot="fixed"
  >
    <media-filter-button
      :sort="sort"
      :filter="filter"
      :non-public-only="nonPublicOnly"
      @update:sort="(newSort: 'asc' | 'desc') => emit('update:sort', newSort)"
      @update:filter="(type: '' | 'image' | 'video' | 'audio' | 'favourites') => emit('update:filter', type)"
      @update:non-public-only="(nonPublicOnly: boolean) => emit('update:non-public-only', nonPublicOnly)"
    />

    <media-select-button
      :is-select-mode="isSelectMode"
      @toggle-select-mode="emit('toggle-select-mode')"
      @favourite:set="emit('favourite:set')"
      @favourite:unset="emit('favourite:unset')"
      @album:select="emit('album:select')"
      @album:unselect="emit('album:unselect')"
      @album:cover="emit('album:cover')"
      @public:set="$emit('public:set')"
      @public:unset="$emit('public:unset')"
      @edit="emit('edit')"
      @delete="emit('delete')"
    />
  </div>
</template>

<script setup lang="ts">
import mediaFilterButton from '@/components/media/media-filter-button.vue';
import mediaSelectButton from '@/components/media/media-select-button.vue';

defineProps<{
  isSelectMode: boolean;
  sort: 'asc' | 'desc';
  filter: '' | 'image' | 'video' | 'audio' | 'favourites';
  nonPublicOnly: boolean;
}>();

const emit = defineEmits<{
  (e: 'toggle-select-mode'): void;
  (e: 'update:sort', value: 'asc' | 'desc'): void;
  (e: 'update:filter', value: '' | 'image' | 'video' | 'audio' | 'favourites'): void;
  (e: 'update:non-public-only', value: boolean): void;
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
</script>

<style scoped>
.button-row {
  bottom: 32px;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  width: 100%;
  pointer-events: none;
}

.button-row > * {
  pointer-events: auto;
}
</style>
