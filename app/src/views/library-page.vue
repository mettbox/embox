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
            :icon="addCircleOutline"
          />
        </ion-button>
      </template>
    </app-header>

    <ion-content
      :fullscreen="true"
      color="primary"
      scroll-events
      @ionScroll="updateVisibleRange"
      ref="contentRef"
      class="ion-no-scroll-on-pinch"
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

      <section
        class="media-grid"
        :style="{ '--columns': columns }"
      >
        <template
          v-for="(media, index) in filteredMediaList"
          :key="media.id"
        >
          <div
            v-if="shouldShowHeader(index)"
            class="grid-header"
          >
            {{ getHeaderLabel(media.date) }}
          </div>
          <media-grid-item
            :is-select-mode="isSelectMode"
            :media="media"
            :is-single-column="columns === 1"
            @click="onMediaOpen(media)"
          />
        </template>
      </section>

      <media-toolbar
        :sort="sort"
        :filter="mediatype === '' && favourites ? 'favourites' : mediatype"
        :non-public-only="nonPublicOnly"
        :is-select-mode="isSelectMode"
        @toggle-select-mode="toggleSelectMode"
        @update:sort="onChangedSort"
        @update:filter="onChangedFilter"
        @update:non-public-only="onChangedNonPublicOnly"
        @favourite:set="onSetFavourites"
        @favourite:unset="onUnsetFavourites"
        @album:select="onSelectAlbum"
        @album:unselect="onAlbumUnselect"
        @album:cover="onAlbumCoverSet"
        @edit="onEditMedia"
        @delete="onDelete"
      />
    </ion-content>

    <media-single-modal
      :is-open="isMediaSingleModalOpen"
      :media="openMedia"
      :prev-media-id="prevMediaId"
      :next-media-id="nextMediaId"
      @did-dismiss="onCloseMedia"
      @prev="onShowPrevMedia"
      @next="onShowNextMedia"
      @favourite:set="onSetFavourites"
      @favourite:unset="onUnsetFavourites"
      @album:select="onSelectAlbum"
      @album:unselect="onAlbumUnselect"
      @album:cover="onAlbumCoverSet"
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
  </ion-page>
</template>

<script setup lang="ts">
import {
  IonPage,
  IonContent,
  IonIcon,
  IonButton,
  onIonViewDidEnter,
  onIonViewWillEnter,
  IonItem,
  IonLabel,
  createGesture,
} from '@ionic/vue';
import { addCircleOutline } from 'ionicons/icons';
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
import mediaGridItem from '@/components/media/media-grid-item.vue';
import mediaToolbar from '@/components/media/media-toolbar.vue';
import albumSelectModal from '@/components/album/album-select-modal.vue';
import mediaEditModal from '@/components/media/media-edit-modal.vue';
import { useSelectedMediaStore } from '@/stores/selectedMedia.store';

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
    return visibleRange.value || t('Library');
  }
});

const filteredMediaList = computed(() => {
  if (mediaList.value.length === 0) return [];

  const filtered = mediaList.value.filter((media) => {
    if (favourites.value && !media.isFavourite) return false;
    if (mediatype.value && media.type !== mediatype.value) return false;
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

const prevMediaId = computed(() => {
  if (!openMedia.value) return null;
  const currentIndex = filteredMediaList.value.findIndex((m) => m.id === openMedia.value.id);
  return currentIndex > 0 ? filteredMediaList.value[currentIndex - 1].id : null;
});

const nextMediaId = computed(() => {
  if (!openMedia.value) return null;
  const currentIndex = filteredMediaList.value.findIndex((m) => m.id === openMedia.value.id);
  return currentIndex < filteredMediaList.value.length - 1 ? filteredMediaList.value[currentIndex + 1].id : null;
});

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const app = useAppStore();
const me = useMeStore();
const isAdmin = me.isAdmin;

const { releaseThumbnail } = useThumbnail();
const selectedMedia = useSelectedMediaStore();

const contentRef = ref<any>(null);
const columns = ref(4);

onIonViewWillEnter(() => {
  columns.value = isCollection.value ? app.mediaGridColumns - 1 : app.mediaGridColumns;
});

let startDistance = 0;
let initialColumns = columns.value;

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
const nonPublicOnly = ref<boolean>(false);

const collectionTitle = ref<string>('');
const collectionDescription = ref<string>('');

const shouldShowHeader = (index: number) => {
  // no date header if columns are less than 10
  if (columns.value < 8) return false;
  const current = filteredMediaList.value[index];
  if (index === 0) return true;
  const prev = filteredMediaList.value[index - 1];

  const currentDate = new Date(current.date);
  const prevDate = new Date(prev.date);

  return currentDate.getMonth() !== prevDate.getMonth() || currentDate.getFullYear() !== prevDate.getFullYear();
};

const getHeaderLabel = (date: string) => {
  return new Intl.DateTimeFormat('de-DE', {
    month: 'short',
    year: 'numeric',
  }).format(new Date(date));
};

const updateVisibleRange = () => {
  const grid = document.querySelector('.media-grid');
  if (!grid) return;

  const items = Array.from(grid.querySelectorAll('.media-grid-item'));
  const visibleDates: string[] = [];

  const viewportTop = window.scrollY;
  const viewportBottom = viewportTop + window.innerHeight;

  items.forEach((item) => {
    const rect = (item as HTMLElement).getBoundingClientRect();
    const cellTop = rect.top + window.scrollY;
    const cellBottom = rect.bottom + window.scrollY;

    if (cellBottom > viewportTop && cellTop < viewportBottom) {
      const id = (item as HTMLElement).dataset.id;
      const media = filteredMediaList.value.find((m) => String(m.id) === id);
      if (media?.date) visibleDates.push(media.date);
    }
  });

  if (!visibleDates.length) {
    visibleRange.value = '';
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

  visibleRange.value =
    firstDate.getFullYear() === lastDate.getFullYear() && firstDate.getMonth() === lastDate.getMonth()
      ? firstLabel
      : `${firstLabel} – ${lastLabel}`;
};

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

const onChangedNonPublicOnly = (newPublicOnly: boolean) => {
  nonPublicOnly.value = newPublicOnly;
};

const onSetFavourites = async () => {
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

const onSelectAlbum = async () => {
  hasOpenAlbumModal.value = true;
};

const onAlbumSelected = async (id: number) => {
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

const onAlbumCoverSet = async () => {
  if (route.name !== 'albums' || !props.id || selectedMedia.ids.length === 0) {
    return;
  }

  try {
    await AlbumService.setCover(Number(props.id), selectedMedia.ids[0]);
    app.setNotification(t('Set cover for album'));
  } catch (error: unknown) {
    app.setNotification(t('Failed to set cover for album'), error);
  } finally {
    selectedMedia.clear();
    isSelectMode.value = false;
  }
};

const onEditMedia = () => {
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
  if (isSelectMode.value) {
    const index = selectedMedia.ids.indexOf(media.id);
    index === -1 ? selectedMedia.add(media) : selectedMedia.remove(media.id);
  } else {
    selectedMedia.set(media);
    openMedia.value = media;
    isMediaSingleModalOpen.value = true;
  }
};

const onCloseMedia = () => {
  isMediaSingleModalOpen.value = false;
  selectedMedia.clear();
  openMedia.value = null;
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

const initGestures = () => {
  const contentEl = contentRef.value.$el;

  const gesture = createGesture({
    el: contentEl,
    gestureName: 'pinch-zoom',
    priority: 110, // Higher priority
    threshold: 0,
    canStart: (ev) => {
      const isTwoFingers = (ev.event as TouchEvent).touches?.length === 2;
      if (isTwoFingers) {
        // Temporarily lock the element to prevent system gestures
        contentEl.style.touchAction = 'none';
      }
      return isTwoFingers;
    },
    onStart: (ev) => {
      const touches = (ev.event as TouchEvent).touches;
      if (touches.length === 2) {
        startDistance = Math.hypot(touches[0].pageX - touches[1].pageX, touches[0].pageY - touches[1].pageY);
        initialColumns = columns.value;
      }
    },
    onMove: (ev) => {
      const touches = (ev.event as TouchEvent).touches;
      if (!touches || touches.length !== 2) return;

      // Absolute prevention of native zoom
      if (ev.event.cancelable) {
        ev.event.preventDefault();
        ev.event.stopPropagation();
      }

      const currentDistance = Math.hypot(touches[0].pageX - touches[1].pageX, touches[0].pageY - touches[1].pageY);

      if (startDistance === 0) {
        startDistance = currentDistance;
        initialColumns = columns.value;
        return;
      }

      const scale = currentDistance / startDistance;
      let targetColumns = initialColumns;

      // Map scale to columns
      if (scale > 1.05) {
        const steps = Math.floor((scale - 1) / 0.15);
        targetColumns = Math.max(1, initialColumns - steps);
      } else if (scale < 0.95) {
        const steps = Math.floor((1 - scale) / 0.08);
        // Allow up to 20 columns
        targetColumns = Math.min(20, initialColumns + steps);
      }

      if (targetColumns !== columns.value) {
        columns.value = targetColumns;
        updateVisibleRange();
      }
    },
    onEnd: () => {
      startDistance = 0;
      // Restore default touch action
      contentEl.style.touchAction = '';
    },
  });

  gesture.enable(true);

  const preventNative = (e: any) => {
    if (e.touches && e.touches.length > 1) {
      if (e.cancelable) e.preventDefault();
    }
  };

  // Attach to document as well because Safari is stubborn with bubbling
  document.addEventListener('gesturestart', preventNative, { passive: false });
  document.addEventListener('gesturechange', preventNative, { passive: false });
  document.addEventListener('gestureend', preventNative, { passive: false });
  contentEl.addEventListener('touchmove', preventNative, { passive: false });
};

onIonViewDidEnter(async () => {
  await load();
  initGestures();
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

ion-content {
  .media-grid {
    display: grid;
    /* Use fixed column count for predictability */
    grid-template-columns: repeat(var(--columns), 1fr);
    gap: 1px;
    padding: 1px;

    /* Smooth transitions when column count changes */
    transition: grid-template-columns 0.3s ease-out;

    .grid-header {
      grid-column: 1 / -1;
      padding: var(--ion-padding-half, 8px) var(--ion-padding, 16px);
      font-size: 0.875rem;
    }
  }
}

/* Safari specific: prevent rubber banding and native zoom interference */
.ion-no-scroll-on-pinch {
  touch-action: pan-x pan-y;
}
</style>
