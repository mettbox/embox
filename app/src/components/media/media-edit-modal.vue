<template>
  <ion-modal
    class="edit-media-modal"
    ref="modalRef"
    :is-open="props.isOpen"
    @will-present="onWillPresent"
    @did-dismiss="onClose"
  >
    <ion-header class="ion-no-border">
      <ion-toolbar color="primary">
        <ion-title class="ion-text-center">{{ $t('Edit') }}</ion-title>
        <ion-buttons slot="end">
          <ion-button @click="onClose()">
            <ion-icon
              slot="icon-only"
              :icon="close"
            />
          </ion-button>
        </ion-buttons>
      </ion-toolbar>
    </ion-header>

    <ion-content>
      <ion-list>
        <ion-item
          v-for="(media, index) in items"
          :key="media.id"
        >
          <ion-thumbnail
            slot="start"
            @click="items[index] && (items[index].isPublic = !items[index].isPublic)"
            :class="{ 'is-public': items[index].isPublic }"
          >
            <ion-img :src="media.thumbUrl" />
            <ion-icon :icon="items[index].isPublic ? eyeOutline : eyeOffOutline" />
          </ion-thumbnail>

          <ion-label>
            <ion-textarea
              :rows="1"
              :auto-grow="true"
              :placeholder="media.caption || $t('Add a caption…')"
              v-model="items[index].caption"
            />
          </ion-label>
        </ion-item>
      </ion-list>
    </ion-content>

    <ion-footer>
      <ion-toolbar
        color="primary"
        class="footer"
      >
        <ion-buttons slot="end">
          <ion-button @click="onSave">
            <ion-icon
              slot="icon-only"
              :icon="saveOutline"
            />
          </ion-button>
        </ion-buttons>
      </ion-toolbar>
    </ion-footer>
  </ion-modal>
</template>

<script setup lang="ts">
import {
  IonModal,
  IonButton,
  IonContent,
  IonTitle,
  IonHeader,
  IonToolbar,
  IonIcon,
  IonButtons,
  IonList,
  IonFooter,
  IonItem,
  IonLabel,
  IonTextarea,
  IonThumbnail,
  IonImg,
} from '@ionic/vue';
import { ref } from 'vue';
import { close, saveOutline, eyeOffOutline, eyeOutline } from 'ionicons/icons';
import { useSelectedMediaStore } from '@/stores/selectedMedia.store';

const props = defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'did-edit', updates: Partial<Media>[]): void;
}>();

const modalRef = ref<InstanceType<typeof IonModal> | null>(null);
const selectedMedia = useSelectedMediaStore();
const items = ref<Media[]>([]);

const onClose = () => {
  if (modalRef.value) {
    modalRef.value.$el.dismiss(null, 'cancel');
  }
};

const onSave = () => {
  console.log(items.value);
  const updates = items.value.map((item) => ({
    id: item.id,
    caption: item.caption,
    isPublic: item.isPublic,
  }));
  emit('did-edit', updates);
};

const onWillPresent = async () => {
  items.value = selectedMedia.all;
};
</script>

<style>
.edit-media-modal {
  ion-item {
    --inner-padding-end: 4px;
  }

  ion-thumbnail {
    --size: 96px;
    position: relative;
    display: flex;
    justify-content: center;
    align-items: center;
    overflow: hidden;
    background-color: var(--ion-color-light);

    ion-img {
      opacity: 0.6;
    }

    &.is-public {
      ion-img {
        opacity: 1;
      }
    }

    ion-icon {
      position: absolute;
      bottom: 4px;
      right: 4px;
      font-size: 24px;
    }
  }
}
</style>
