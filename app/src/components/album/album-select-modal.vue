<template>
  <ion-modal
    class="album-select-modal"
    ref="modalRef"
    :is-open="props.isOpen"
    @will-present="onWillPresent"
    @did-dismiss="onClose"
  >
    <ion-header class="ion-no-border">
      <ion-toolbar color="primary">
        <ion-title class="ion-text-center">{{ $t('Add to Album') }}</ion-title>
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
      <ion-segment>
        <ion-segment-button
          value="select"
          content-id="select"
        >
          <ion-label>{{ $t('Select') }}</ion-label>
        </ion-segment-button>
        <ion-segment-button
          value="add"
          content-id="add"
        >
          <ion-label>{{ $t('Create') }}</ion-label>
        </ion-segment-button>
      </ion-segment>
      <ion-segment-view>
        <ion-segment-content id="select">
          <ion-list>
            <collection-list-item
              v-for="item in albums"
              :key="item.id"
              :item="item"
              :can-edit="true"
              @click="onAlbumSelect"
            />
          </ion-list>
        </ion-segment-content>
        <ion-segment-content id="add">
          <album-form @done="onAlbumCreate" />
        </ion-segment-content>
      </ion-segment-view>
    </ion-content>
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
  IonLabel,
  IonSegment,
  IonSegmentButton,
  IonSegmentContent,
  IonSegmentView,
  IonList,
} from '@ionic/vue';
import { ref } from 'vue';
import { close } from 'ionicons/icons';
import collectionListItem from '@/components/collection/collection-list-item.vue';
import albumForm from '@/components/album/album-form.vue';
import { AlbumService } from '@/services/album.service';

const props = defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  // (e: 'close'): void;
  (e: 'album:set', id: number): void;
}>();

const modalRef = ref<InstanceType<typeof IonModal> | null>(null);
const albums = ref<CollectionListItem[]>([]);

const onAlbumSelect = (id: string) => {
  emit('album:set', Number(id));
  onClose();
};

const onAlbumCreate = async (album: AlbumForm) => {
  const createdAlbum = await AlbumService.create(album);
  emit('album:set', createdAlbum.id);
  onClose();
};

const onClose = () => {
  if (modalRef.value) {
    modalRef.value.$el.dismiss(null, 'cancel');
  }
  // emit('close');
};

const onWillPresent = async () => {
  albums.value = await AlbumService.get();
};
</script>
