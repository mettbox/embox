<template>
  <ion-item-sliding ref="itemRef">
    <ion-item
      :detail="props.hasDetails"
      :button="true"
      @click="() => onClick(item.id)"
    >
      <ion-thumbnail slot="start">
        <ion-img
          v-if="item.media"
          :src="getThumbnailUrl(item.media.id)"
        />
      </ion-thumbnail>
      <ion-label>
        {{ item.name }}
        <ion-note v-if="item.description">{{ item.description }}</ion-note>
      </ion-label>
      <ion-note slot="end">{{ item.mediaCount }}</ion-note>
    </ion-item>
    <ion-item-options v-if="canEdit">
      <ion-item-option
        color="medium"
        @click="() => onAction('edit', item)"
        ><ion-icon
          slot="top"
          :icon="createOutline"
        ></ion-icon
        >{{ $t('Edit') }}</ion-item-option
      >
      <ion-item-option
        color="primary"
        @click="() => onAction('remove', item)"
        ><ion-icon
          slot="top"
          :icon="trashOutline"
        ></ion-icon
        >{{ $t('Delete') }}</ion-item-option
      >
    </ion-item-options>
  </ion-item-sliding>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import {
  IonImg,
  IonThumbnail,
  IonItem,
  IonLabel,
  IonNote,
  IonItemSliding,
  IonItemOptions,
  IonItemOption,
  IonIcon,
} from '@ionic/vue';
import { trashOutline, createOutline } from 'ionicons/icons';
import { useThumbnail } from '@/composables/use-thumbnail';

const props = withDefaults(
  defineProps<{
    item: CollectionListItem;
    hasDetails?: boolean;
    canEdit?: boolean;
  }>(),
  {
    hasDetails: true,
    canEdit: false,
  },
);

const emit = defineEmits<{
  (e: 'click', id: string): void;
  (e: 'edit', item: CollectionListItem): void;
  (e: 'remove', item: CollectionListItem): void;
}>();

const { getThumbnailUrl } = useThumbnail();

const itemRef = ref<VueInstanceElement>();

const onClick = (id: string) => {
  emit('click', id);
};

const onAction = (action: string, item: CollectionListItem): void => {
  if (itemRef.value) {
    itemRef.value.$el.closeOpened();
  }

  if (action === 'edit') {
    emit('edit', item);
  } else {
    emit('remove', item);
  }
};
</script>

<style scoped>
ion-thumbnail {
  --size: 96px;
}

ion-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

ion-label {
  ion-note {
    display: block;
  }
}
</style>
