<template>
  <ion-page>
    <app-header>
      {{ $t('Collections') }}
    </app-header>

    <ion-content>
      <ion-list>
        <ion-item-divider sticky>
          {{ $t('Favourites') }}
        </ion-item-divider>
        <collection-list-item
          v-for="item in favourites"
          :key="item.id"
          :item="item"
          @click="openFavourite"
        />
        <ion-item-divider sticky>
          {{ $t('Albums') }}
        </ion-item-divider>
        <collection-list-item
          v-for="item in albums"
          :key="item.id"
          :item="item"
          :can-edit="true"
          @click="openAlbum"
          @edit="onEdit"
          @remove="onRemove"
        />
      </ion-list>

      <ion-action-sheet
        :is-open="hasOpenDelete"
        :header="$t('This action is permanent and cannot be undone.')"
        :buttons="deleteActionSheetButtons"
        @did-dismiss="deleteHandler"
      />

      <album-edit-modal
        v-if="album"
        :is-open="hasOpenModal"
        :album="album"
        @did-dismiss="editHandler"
      />
    </ion-content>
  </ion-page>
</template>

<script setup lang="ts">
import {
  IonPage,
  IonContent,
  onIonViewDidEnter,
  IonList,
  IonItemDivider,
  IonActionSheet,
  ActionSheetButton,
} from '@ionic/vue';
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { useAppStore } from '@/stores/app.store';
import { FavouriteService } from '@/services/favourite.service';
import { AlbumService } from '@/services/album.service';
import appHeader from '@/components/app/app-header.vue';
import collectionListItem from '@/components/collection/collection-list-item.vue';
import albumEditModal from '@/components/album/album-edit-modal.vue';

const deleteActionSheetButtons = computed((): ActionSheetButton[] => [
  {
    text: t('Remove', { name: album.value ? album.value.name : '' }),
    role: 'destructive',
  },
  {
    text: t('Cancel'),
    role: 'cancel',
  },
]);

const { t } = useI18n();
const router = useRouter();
const app = useAppStore();

const favourites = ref<CollectionListItem[]>([]);
const albums = ref<CollectionListItem[]>([]);

const hasOpenDelete = ref(false);
const hasOpenModal = ref(false);
const album = ref<AlbumForm | null>(null);

const onEdit = (item: CollectionListItem): void => {
  album.value = {
    id: Number(item.id),
    name: item.name,
    description: item.description || '',
    isPublic: item.isPublic || false,
  };
  hasOpenModal.value = true;
};

const onRemove = (item: CollectionListItem): void => {
  album.value = {
    id: Number(item.id),
    name: item.name,
    description: item.description || '',
    isPublic: item.isPublic || false,
  };
  hasOpenDelete.value = true;
};

const deleteHandler = async (event: { detail: { role: string } }): Promise<void> => {
  hasOpenDelete.value = false;
  const action = event?.detail?.role || 'cancel';

  if (album.value === null || action !== 'destructive') {
    return;
  }

  try {
    await AlbumService.delete(Number(album.value.id));
    albums.value = albums.value.filter((a) => Number(a.id) !== album.value?.id);
  } catch (error: unknown) {
    app.setNotification(t('Failed to delete album'), error);
  }
};

const editHandler = async (event: CustomEvent): Promise<void> => {
  hasOpenModal.value = false;

  if (!event.detail || event.detail.role !== 'save') {
    return;
  }

  try {
    const album = event.detail.data as AlbumForm;
    await AlbumService.update(album.id!, album);
    const index = albums.value.findIndex((a) => Number(a.id) === album.id);
    if (index !== -1) {
      albums.value[index] = { ...albums.value[index], ...album, id: String(album.id) };
    }
  } catch (error: unknown) {
    app.setNotification(t('Failed to edit album'), error);
  }
};

const openFavourite = (id: string) => {
  router.push({ name: 'favourites', params: { id } });
};

const openAlbum = (id: string) => {
  router.push({ name: 'albums', params: { id } });
};

onIonViewDidEnter(async () => {
  favourites.value = await FavouriteService.getUsersWithLatestFavourite();
  albums.value = await AlbumService.get();
});
</script>

<style scoped></style>
