<template>
  <ion-modal
    class="album-edit-modal"
    ref="modalRef"
    :is-open="props.isOpen"
    @did-dismiss="onClose"
  >
    <ion-header class="ion-no-border">
      <ion-toolbar color="primary">
        <ion-title class="ion-text-center">{{ $t('Edit Album') }}</ion-title>
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
      <album-form
        :album="album"
        @done="onUpdateAlbum"
      />
    </ion-content>
  </ion-modal>
</template>

<script setup lang="ts">
import { IonModal, IonButton, IonContent, IonTitle, IonHeader, IonToolbar, IonIcon, IonButtons } from '@ionic/vue';
import { ref } from 'vue';
import { close } from 'ionicons/icons';
import albumForm from '@/components/album/album-form.vue';

const props = defineProps<{
  isOpen: boolean;
  album: AlbumForm;
}>();

const modalRef = ref();

const onUpdateAlbum = async (album: AlbumForm) => {
  modalRef.value?.$el.dismiss(album, 'save');
};

const onClose = () => {
  if (modalRef.value) {
    modalRef.value.$el.dismiss(null, 'cancel');
  }
};
</script>
