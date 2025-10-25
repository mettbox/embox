<template>
  <ion-page>
    <app-header :has-back-button="isCollection">
      {{ pageTitle }}
      <template #buttons>
        <ion-button
          v-if="isAdmin && !isCollection"
          @click="openUploadModal"
          :selected="hasOpenUploadModal === true"
        >
          <ion-icon
            slot="icon-only"
            :icon="cloudUploadOutline"
          />
        </ion-button>
      </template>
    </app-header>

    <ion-content
      :fullscreen="true"
      color="primary"
    >
      <div
        v-if="hasMediaListLoaded && filteredMediaList.length === 0"
        class="ion-padding"
      >
        <p>{{ $t('No media found.') }}</p>
      </div>

      <ion-item
        v-if="collectionDescription"
        lines="none"
        color="primary"
      >
        <ion-label class="ion-text-wrap ion-padding-top ion-padding-bottom ion-text-center">
          {{ collectionDescription }}
        </ion-label>
      </ion-item>

      <media-grid
        :mediaList="filteredMediaList"
        :columns="isCollection ? app.mediaGridColumns - 1 : app.mediaGridColumns"
        :isSelectMode="isSelectMode"
        @update:visible-range="(range: string) => (visibleRange = range)"
        @media:open="onMediaOpen"
        @media:select="onMediaSingleSelect"
      />

      <media-toolbar
        :sort="sort"
        :filter="mediatype === '' && favourites ? 'favourites' : mediatype"
        :public-only="publicOnly"
        :is-select-mode="isSelectMode"
        @toggle-select-mode="toggleSelectMode"
        @update:sort="onChangedSort"
        @update:filter="onChangedFilter"
        @update:public-only="onChangedPublicOnly"
        @favourite:set="onSetFavourites"
        @favourite:unset="onUnsetFavourites"
        @album:select="onSelectAlbum"
        @album:unselect="onAlbumUnselect"
        @public:set="onSetPublic"
        @public:unset="onUnsetPublic"
        @edit="onEditMedia"
        @delete="onDelete"
      />
    </ion-content>

    <media-single-modal
      :is-open="isMediaSingleModalOpen"
      :media="openMedia"
      :has-prev="hasPrev"
      :has-next="hasNext"
      @did-dismiss="onCloseMedia"
      @prev="onShowPrevMedia"
      @next="onShowNextMedia"
      @favourite:set="onSetFavourites"
      @favourite:unset="onUnsetFavourites"
      @album:select="onSelectAlbum"
      @album:unselect="onAlbumUnselect"
      @public:set="onSetPublic"
      @public:unset="onUnsetPublic"
      @edit="onEditMedia"
      @delete="onDelete"
    />

    <media-upload-modal
      ref="uploadModalRef"
      :is-open="hasOpenUploadModal"
      @did-dismiss="hasOpenUploadModal = false"
      @did-upload="onDidUpload"
    />

    <album-select-modal
      :is-open="hasOpenAlbumModal"
      @album:set="onAlbumSelected"
      @did-dismiss="hasOpenAlbumModal = false"
    />

    <media-edit-modal
      ref="editMediaModalRef"
      :is-open="hasOpenEditMediaModal"
      @did-dismiss="hasOpenEditMediaModal = false"
      @did-edit="onDidEditMedia"
    />

    <media-context-popover
      :is-open="hasOpenContextPopover"
      @close="onCloseMediaContextPopover"
      @update:sort="onChangedSort"
      @update:filter="onChangedFilter"
      @update:public-only="onChangedPublicOnly"
      @favourite:set="onSetFavourites"
      @favourite:unset="onUnsetFavourites"
      @public:set="onSetPublic"
      @public:unset="onUnsetPublic"
      @album:set="onSelectAlbum"
      @album:unset="onAlbumUnselect"
      @edit="onEditMedia"
      @delete="onDelete"
    />
  </ion-page>
</template>

<script setup lang="ts">
import { IonPage, IonContent, IonIcon, IonButton, onIonViewDidEnter, IonItem, IonLabel } from '@ionic/vue';
import { cloudUploadOutline } from 'ionicons/icons';
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { useAppStore } from '@/stores/app.store';
import { useMeStore } from '@/stores/me.store';
import { useThumbnail } from '@/composables/use-thumbnail';
import { useRoute, useRouter } from 'vue-router';
import { MediaService } from '@/services/media.service';
import { FavouriteService } from '@/services/favourite.service';
import { AlbumService } from '@/services/album.service';
import mediaSingleModal from '@/components/media/media-single-modal.vue';
import mediaUploadModal from '@/components/media/media-upload-modal.vue';
import appHeader from '@/components/app/app-header.vue';
import mediaGrid from '@/components/media/media-grid.vue';
import mediaToolbar from '@/components/media/media-toolbar.vue';
import albumSelectModal from '@/components/album/album-select-modal.vue';
import mediaEditModal from '@/components/media/media-edit-modal.vue';
import { useSelectedMediaStore } from '@/stores/selectedMedia.store';
import 'vue-virtual-scroller/dist/vue-virtual-scroller.css';
import mediaContextPopover from '@/components/media/media-context-popover.vue';

const props = defineProps<{
  id?: string;
}>();

const isCollection = computed(() => route.name === 'favourites' || route.name === 'albums');

const pageTitle = computed(() => {
  if (route.name === 'favourites') {
    return t("Favourite's", { name: collectionTitle.value });
  } else if (route.name === 'albums') {
    return collectionTitle.value;
  } else {
    return visibleRange.value || 'Library';
  }
});

const filteredMediaList = computed(() => {
  if (mediaList.value.length === 0) return [];

  const filtered = mediaList.value.filter((media) => {
    if (favourites.value && !media.isFavourite) return false;
    if (mediatype.value && media.type !== mediatype.value) return false;
    if (isAdmin && publicOnly.value && !media.isPublic) return false;
    return true;
  });

  return filtered.sort((a, b) => {
    if (sort.value === 'asc') {
      return new Date(a.date).getTime() - new Date(b.date).getTime();
    } else {
      return new Date(b.date).getTime() - new Date(a.date).getTime();
    }
  });
});

const hasPrev = computed(() => {
  if (!openMedia.value) return false;
  const currentIndex = filteredMediaList.value.findIndex((m) => m.id === openMedia.value.id);
  return currentIndex > 0;
});

const hasNext = computed(() => {
  if (!openMedia.value) return false;
  const currentIndex = filteredMediaList.value.findIndex((m) => m.id === openMedia.value.id);
  return currentIndex < filteredMediaList.value.length - 1;
});

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const app = useAppStore();
const me = useMeStore();
const isAdmin = me.isAdmin;

const { releaseThumbnail } = useThumbnail();
const selectedMedia = useSelectedMediaStore();

const mediaList = ref<Media[]>([]);
const visibleRange = ref('');
const hasMediaListLoaded = ref(false);

const uploadModalRef = ref<InstanceType<typeof mediaUploadModal> | null>(null);
const hasOpenUploadModal = ref(false);

const editMediaModalRef = ref<InstanceType<typeof mediaEditModal> | null>(null);
const hasOpenEditMediaModal = ref(false);

const hasOpenAlbumModal = ref(false);

const openMedia = ref<any>(null);
const isMediaSingleModalOpen = ref(false);

const isSelectMode = ref(false);

const sort = ref<'asc' | 'desc'>('desc');
const mediatype = ref<'image' | 'video' | 'audio' | ''>('');
const favourites = ref<boolean>(false);
const publicOnly = ref<boolean>(false);

const collectionTitle = ref<string>('');
const collectionDescription = ref<string>('');

const openUploadModal = () => {
  hasOpenUploadModal.value = true;
  setTimeout(() => {
    if (uploadModalRef.value) {
      uploadModalRef.value.triggerUploadInput();
    }
  }, 100);
};

const onDidUpload = async (response: UploadResponse) => {
  hasOpenUploadModal.value = false;
  if (response.failed.length > 0) {
    app.setNotification(
      t('Some files failed to upload', { success: response.success.length, failed: response.failed }),
    );
  } else {
    app.setNotification(t('Uploaded media', { count: response.success.length }));
  }

  const newMedia = response.success.map((item: any) => ({
    ...item,
    thumbUrl: '',
    isLoaded: false,
  }));
  mediaList.value = [...newMedia, ...mediaList.value];
};

const load = async () => {
  try {
    let list = [];
    selectedMedia.reset();
    if (route.name === 'favourites' && props.id) {
      const favourites = (await FavouriteService.getFavouritesByUserID(props.id)) || [];
      collectionTitle.value = favourites.user.name;
      list = favourites.media;
    } else if (route.name === 'albums' && props.id) {
      const album = (await AlbumService.getByID(Number(props.id))) || [];
      selectedMedia.setAlbumId(Number(props.id));
      collectionTitle.value = album.name;
      collectionDescription.value = album.description;
      list = album.media;
    } else {
      list = (await MediaService.getMediaList()) || [];
    }

    mediaList.value = list.map((item: any) => ({
      ...item,
      thumbUrl: '',
      isLoaded: false,
    }));
  } catch (err) {
    console.error('Error fetching media', err);
  } finally {
    hasMediaListLoaded.value = true;
  }
};

const toggleSelectMode = () => {
  isSelectMode.value = !isSelectMode.value;
  if (!isSelectMode.value) {
    selectedMedia.clear();
  }
};

const onChangedSort = (newSort: 'asc' | 'desc') => {
  sort.value = newSort;
};

const onChangedFilter = (newFilter: '' | 'image' | 'video' | 'audio' | 'favourites') => {
  if (newFilter === 'favourites') {
    favourites.value = !favourites.value;
    mediatype.value = '';
  } else {
    mediatype.value = newFilter;
    favourites.value = false;
  }
};

const onChangedPublicOnly = (newPublicOnly: boolean) => {
  publicOnly.value = newPublicOnly;
};

const onSetFavourites = async () => {
  hasOpenContextPopover.value = false;
  if (selectedMedia.ids.length === 0) return;

  try {
    const ids = selectedMedia.ids;
    await FavouriteService.add(ids);
    app.setNotification(t('Added favourites', { count: ids.length }));

    ids.forEach((id) => {
      const media = mediaList.value.find((m) => m.id === id);
      if (media) {
        media.isFavourite = true;
      }
    });
  } catch (error: unknown) {
    app.setNotification(t('Failed to add favourites'), error);
  } finally {
    if (openMedia.value && selectedMedia.ids.includes(openMedia.value.id)) {
      openMedia.value.isFavourite = true;
    } else {
      selectedMedia.clear();
      isSelectMode.value = false;
    }
  }
};

const onUnsetFavourites = async () => {
  hasOpenContextPopover.value = false;
  if (selectedMedia.ids.length === 0) return;

  try {
    const ids = selectedMedia.ids;
    await FavouriteService.remove(ids);
    app.setNotification(t('Removed favourites', { count: ids.length }));

    ids.forEach((id) => {
      const media = mediaList.value.find((m) => m.id === id);
      if (media) {
        media.isFavourite = false;
      }
    });
  } catch (error: unknown) {
    app.setNotification(t('Failed to remove favourites'), error);
  } finally {
    if (openMedia.value && selectedMedia.ids.includes(openMedia.value.id)) {
      openMedia.value.isFavourite = false;
    } else {
      selectedMedia.clear();
      isSelectMode.value = false;
    }
  }
};

const onSetPublic = async () => {
  hasOpenContextPopover.value = false;
  if (selectedMedia.ids.length === 0) return;

  try {
    const ids = selectedMedia.ids;
    await MediaService.setPublic(ids, true);
    app.setNotification(t('Set public', { count: ids.length }));

    ids.forEach((id) => {
      const media = mediaList.value.find((m) => m.id === id);
      if (media) {
        media.isPublic = true;
      }
    });
  } catch (error: unknown) {
    app.setNotification(t('Failed to set public'), error);
  } finally {
    if (openMedia.value && selectedMedia.ids.includes(openMedia.value.id)) {
      openMedia.value.isPublic = true;
    } else {
      selectedMedia.clear();
      isSelectMode.value = false;
    }
  }
};

const onUnsetPublic = async () => {
  hasOpenContextPopover.value = false;
  if (selectedMedia.ids.length === 0) return;

  try {
    const ids = selectedMedia.ids;
    await MediaService.setPublic(ids, false);
    app.setNotification(t('Set non public', { count: ids.length }));

    ids.forEach((id) => {
      const media = mediaList.value.find((m) => m.id === id);
      if (media) {
        media.isPublic = false;
      }
    });
  } catch (error: unknown) {
    app.setNotification(t('Failed to set non public'), error);
  } finally {
    if (openMedia.value && selectedMedia.ids.includes(openMedia.value.id)) {
      openMedia.value.isPublic = false;
    } else {
      selectedMedia.clear();
      isSelectMode.value = false;
    }
  }
};

const onSelectAlbum = async () => {
  hasOpenContextPopover.value = false;
  hasOpenAlbumModal.value = true;
};

const onAlbumSelected = async (id: number) => {
  hasOpenContextPopover.value = false;
  if (selectedMedia.ids.length === 0) {
    hasOpenAlbumModal.value = false;
    return;
  }

  try {
    await AlbumService.addMedia(id, selectedMedia.ids);
    app.setNotification(t('Added media to album', { count: selectedMedia.ids.length }));
    router.push({ name: 'albums', params: { id } });
  } catch (error: unknown) {
    app.setNotification(t('Failed to add media to album'), error);
  } finally {
    selectedMedia.clear();
    isSelectMode.value = false;
    hasOpenAlbumModal.value = false;
  }
};

const onAlbumUnselect = async () => {
  hasOpenContextPopover.value = false;
  if (route.name !== 'albums' || !props.id || selectedMedia.ids.length === 0) {
    return;
  }

  try {
    await AlbumService.removeMedia(Number(props.id), selectedMedia.ids);
    app.setNotification(t('Removed media from album', { count: selectedMedia.ids.length }));
    mediaList.value = mediaList.value.filter((media) => !selectedMedia.ids.includes(media.id));
  } catch (error: unknown) {
    app.setNotification(t('Failed to remove media from album'), error);
  } finally {
    selectedMedia.clear();
    isSelectMode.value = false;
  }
};

const onEditMedia = () => {
  hasOpenContextPopover.value = false;
  if (selectedMedia.ids.length === 0) return;
  hasOpenEditMediaModal.value = true;
};

const onDidEditMedia = async (items: Partial<Media>[]) => {
  hasOpenEditMediaModal.value = false;

  if (openMedia.value && selectedMedia.ids.includes(openMedia.value.id)) {
    const updatedItem = items.find((item) => item.id === openMedia.value.id);
    if (updatedItem) {
      openMedia.value = { ...openMedia.value, ...updatedItem };
    }
  } else {
    selectedMedia.clear();
    isSelectMode.value = false;
  }

  try {
    await MediaService.update(items);

    // TODO: is this teh best way to update the mediaList?
    const mediaMap = new Map(mediaList.value.map((media) => [media.id, media]));
    items.forEach((item) => {
      if (!item.id) return;
      const media = mediaMap.get(item.id);
      if (media) {
        Object.assign(media, item);
      }
    });

    app.setNotification(t('Media updated'));
  } catch (error: unknown) {
    app.setNotification(t('Failed to update caption'), error);
  }
};

const onDelete = async () => {
  hasOpenContextPopover.value = false;
  if (selectedMedia.ids.length === 0) return;

  try {
    const ids = selectedMedia.ids;
    const mediaToDelete = mediaList.value.filter((media) => ids.includes(media.id));
    mediaList.value = mediaList.value.filter((media) => !ids.includes(media.id));
    mediaToDelete.forEach((media) => releaseThumbnail(media));

    // Important: small delay to allow UI to update before deletion
    await new Promise((resolve) => setTimeout(resolve, 100));

    await MediaService.delete(ids);
    app.setNotification(t('Deleted media', { count: ids.length }));
  } catch (error: unknown) {
    app.setNotification(t('Failed to delete media'), error);
  } finally {
    selectedMedia.clear();
    isSelectMode.value = false;
  }
};

const onMediaOpen = (media: Media) => {
  openMedia.value = media;
  isMediaSingleModalOpen.value = true;
};

const onCloseMedia = () => {
  isMediaSingleModalOpen.value = false;
  selectedMedia.clear();
  openMedia.value = null;
};

const hasOpenContextPopover = ref(false);

const onMediaSingleSelect = async (media: Media) => {
  selectedMedia.clear();
  selectedMedia.add(media);
  hasOpenContextPopover.value = true;
};

const onCloseMediaContextPopover = () => {
  hasOpenContextPopover.value = false;
  selectedMedia.clear();
};

const onShowPrevMedia = (id: number) => {
  const currentIndex = filteredMediaList.value.findIndex((m) => m.id === id);
  if (currentIndex > 0) {
    openMedia.value = filteredMediaList.value[currentIndex - 1];
    selectedMedia.set(openMedia.value);
  }
};

const onShowNextMedia = (id: number) => {
  const currentIndex = filteredMediaList.value.findIndex((m) => m.id === id);
  if (currentIndex < filteredMediaList.value.length - 1) {
    openMedia.value = filteredMediaList.value[currentIndex + 1];
    selectedMedia.set(openMedia.value);
  }
};

onIonViewDidEnter(async () => {
  await load();
});
</script>

<style scoped>
div[slot='fixed'] {
  bottom: 32px;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 8px;
  width: 100%;
}
</style>
