<template>
  <ion-modal
    ref="modalRef"
    :is-open="isOpen"
    @did-present="onDidPresent"
    @did-dismiss="onDidDismiss"
  >
    <ion-header
      class="ion-no-border"
      v-show="canNavigate"
    >
      <ion-toolbar color="primary">
        <ion-buttons slot="start">
          <ion-button @click="!media.isFavourite ? $emit('favourite:set') : $emit('favourite:unset')">
            <ion-icon
              slot="icon-only"
              :icon="!media.isFavourite ? heartOutline : heart"
            />
          </ion-button>
        </ion-buttons>
        <ion-title class="ion-text-center">
          {{ localizedDate }}
        </ion-title>
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

    <ion-content color="primary">
      <ion-button
        slot="fixed"
        shape="round"
        color="secondary"
        v-if="!canNavigate"
        @click="onClose()"
      >
        <ion-icon
          slot="icon-only"
          :icon="close"
        />
      </ion-button>
      <div
        v-if="fileUrl"
        class="media-wrapper"
        ref="mediaWrapperRef"
      >
        <ion-spinner v-if="!isMediaLoaded && media.type !== 'video'" />
        <Transition :name="slideDirection">
          <div
            :key="media.id"
            class="media-slide"
          >
            <media-zoom
              v-if="media.type === 'image'"
              ref="mediaZoomRef"
              :is-open="true"
              :file-url="fileUrl"
              @ready="onMediaLoaded"
            />
            <video
              v-else-if="media.type === 'video'"
              :class="{ 'main-media': true, loaded: isMediaLoaded }"
              :src="fileUrl"
              controls
              playsinline
              webkit-playsinline
              preload="metadata"
              :poster="getThumbnailUrl(media.id)"
              @loadeddata="onMediaLoaded"
              @loadedmetadata="onMediaLoaded"
            />
            <audio
              v-else-if="media.type === 'audio'"
              :class="{ 'main-media': true, loaded: isMediaLoaded }"
              :src="fileUrl"
              controls
              @loadeddata="onMediaLoaded"
            />
          </div>
        </Transition>
      </div>
      <div
        v-else
        class="media-wrapper"
      >
        <ion-spinner />
      </div>
    </ion-content>

    <ion-footer
      v-show="canNavigate"
      class="ion-no-border"
    >
      <ion-toolbar
        class="footer"
        color="primary"
      >
        <ion-text class="caption">
          <span
            v-for="(line, index) in captionLines"
            :key="index"
          >
            {{ line }}<span v-if="index < captionLines.length - 1"><br /></span>
          </span>
        </ion-text>
        <ion-buttons slot="end">
          <ion-button
            v-if="isAdmin"
            id="edit-popover"
            shape="round"
          >
            <ion-icon
              slot="icon-only"
              :icon="ellipsisHorizontalCircleOutline"
            />
          </ion-button>
        </ion-buttons>
      </ion-toolbar>
    </ion-footer>

    <media-select-popover
      trigger="edit-popover"
      @favourite:set="$emit('favourite:set')"
      @favourite:unset="$emit('favourite:unset')"
      @album:set="$emit('album:select')"
      @album:unset="$emit('album:unselect')"
      @album:cover="$emit('album:cover')"
      @public:set="$emit('public:set')"
      @public:unset="$emit('public:unset')"
      @edit="$emit('edit')"
      @delete="onDelete"
    />
  </ion-modal>
</template>

<script setup lang="ts">
import {
  IonModal,
  IonButton,
  IonContent,
  IonTitle,
  IonHeader,
  IonFooter,
  IonToolbar,
  IonIcon,
  IonButtons,
  IonSpinner,
  IonText,
} from '@ionic/vue';
import { onMounted, onUnmounted, ref, computed, watch } from 'vue';
import { createGesture } from '@ionic/vue';
import { heart, close, ellipsisHorizontalCircleOutline, heartOutline } from 'ionicons/icons';
import mediaSelectPopover from './media-select-popover.vue';
import mediaZoom from './media-zoom.vue';
import { useMeStore } from '@/stores/me.store';
import { useThumbnail } from '@/composables/use-thumbnail';

const props = defineProps<{
  isOpen: boolean;
  media: Media;
  prevMediaId: number | null;
  nextMediaId: number | null;
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'prev', currentMediaId: number): void;
  (e: 'next', currentMediaId: number): void;
  (e: 'delete'): void;
  (e: 'favourite:set'): void;
  (e: 'favourite:unset'): void;
  (e: 'public:set'): void;
  (e: 'public:unset'): void;
  (e: 'album:select'): void;
  (e: 'album:unselect'): void;
  (e: 'album:cover'): void;
  (e: 'edit'): void;
}>();

const localizedDate = computed(() => {
  if (!props.media?.date) return '';
  const date = new Date(props.media.date);
  return date.toLocaleDateString('de-DE', { year: 'numeric', month: 'long', day: 'numeric' });
});

const captionLines = computed(() => (props.media?.caption ?? '').split('\n'));

const me = useMeStore();
const { getFileUrl, getThumbnailUrl } = useThumbnail();

const modalRef = ref<InstanceType<typeof IonModal> | null>(null);
const mediaWrapperRef = ref<HTMLElement | null>(null);
const mediaZoomRef = ref<InstanceType<typeof mediaZoom> | null>(null);
const gesture = ref<ReturnType<typeof createGesture> | null>(null);
const fileUrl = ref<string | undefined>(undefined);
const slideDirection = ref<'slide-left' | 'slide-right'>('slide-left');

const isAdmin = computed(() => me.isAdmin);

const canNavigate = computed(() => {
  if (!mediaZoomRef.value) return true;
  const initialScale = mediaZoomRef.value.initialScale;
  const currentScale = mediaZoomRef.value.currentScale;
  const isPinching = mediaZoomRef.value.isPinching;
  return currentScale <= initialScale * 1.01 && !isPinching;
});

const isMediaLoaded = ref(false);
let originalPreloader: HTMLImageElement | null = null;
let thumbPreloaders: HTMLImageElement[] = [];

const onMediaLoaded = () => {
  isMediaLoaded.value = true;
};

const cancelOriginalPreload = () => {
  if (originalPreloader) {
    originalPreloader.onload = null;
    originalPreloader.onerror = null;
    originalPreloader.src = '';
    originalPreloader = null;
  }
};

const loadOriginalForImage = (mediaId: number) => {
  cancelOriginalPreload();
  const img = new Image();
  originalPreloader = img;
  img.onload = () => {
    if (props.media?.id === mediaId) {
      fileUrl.value = getFileUrl(mediaId);
    }
    originalPreloader = null;
  };
  img.onerror = () => {
    originalPreloader = null;
  };
  img.src = getFileUrl(mediaId);
};

const cancelThumbPreloads = () => {
  for (const img of thumbPreloaders) img.src = '';
  thumbPreloaders = [];
};

const preloadAdjacentThumbs = (...ids: (number | null)[]) => {
  cancelThumbPreloads();
  for (const id of ids) {
    if (!id) continue;
    const img = new Image();
    img.src = getThumbnailUrl(id);
    thumbPreloaders.push(img);
  }
};

watch(
  () => props.media,
  (newMedia, oldMedia) => {
    if (newMedia?.id === oldMedia?.id && fileUrl.value) return;

    cancelOriginalPreload();
    isMediaLoaded.value = false;

    if (!newMedia) {
      fileUrl.value = undefined;
      return;
    }

    if (newMedia.type === 'image') {
      // Show thumbnail immediately, swap to original once preloaded into HTTP cache.
      fileUrl.value = getThumbnailUrl(newMedia.id);
      loadOriginalForImage(newMedia.id);
    } else {
      fileUrl.value = getFileUrl(newMedia.id);
    }

    preloadAdjacentThumbs(props.prevMediaId, props.nextMediaId);
  },
  { immediate: true },
);

const onClose = () => {
  cancelOriginalPreload();
  fileUrl.value = undefined;
  if (modalRef.value) {
    modalRef.value.$el.dismiss(null, 'cancel');
  }
};

const handleKey = (e: KeyboardEvent) => {
  if (e.key === 'ArrowLeft') onPrev();
  if (e.key === 'ArrowRight') onNext();
};

const onPrev = () => {
  if (!props.prevMediaId || !canNavigate.value) return;
  slideDirection.value = 'slide-right';
  mediaZoomRef.value?.reset();
  isMediaLoaded.value = false;
  emit('prev', props.media.id);
};

const onNext = () => {
  if (!props.nextMediaId || !canNavigate.value) return;
  slideDirection.value = 'slide-left';
  mediaZoomRef.value?.reset();
  isMediaLoaded.value = false;
  emit('next', props.media.id);
};

const initSwipe = () => {
  gesture.value?.destroy();
  gesture.value = null;

  const el =
    (modalRef.value?.$el?.querySelector('.modal-wrapper') as HTMLElement | null) ??
    (mediaWrapperRef.value as HTMLElement | null);
  if (!el) return;

  let firedThisGesture = false;

  gesture.value = createGesture({
    el,
    gestureName: 'swipe',
    direction: 'x',
    threshold: 15,
    onStart() {
      firedThisGesture = false;
    },
    onMove(ev) {
      if (firedThisGesture) return;
      if (!canNavigate.value) return;
      if (ev.deltaX > 60 && props.prevMediaId) {
        firedThisGesture = true;
        onPrev();
      } else if (ev.deltaX < -60 && props.nextMediaId) {
        firedThisGesture = true;
        onNext();
      }
    },
  });
  gesture.value.enable(true);
};

const onDidPresent = () => {
  initSwipe();
};

const onDidDismiss = () => {
  gesture.value?.destroy();
  gesture.value = null;
  cancelOriginalPreload();
  fileUrl.value = undefined;
};

const onDelete = () => {
  emit('delete');
  onClose();
};

onMounted(() => {
  window.addEventListener('keydown', handleKey);
});

onUnmounted(() => {
  window.removeEventListener('keydown', handleKey);
  gesture.value?.destroy();
  cancelOriginalPreload();
  cancelThumbPreloads();
});
</script>

<style scoped>
ion-modal {
  --height: 100%;
  --width: 100%;
  --border-radius: 0;
}

ion-button[slot='fixed'] {
  top: 0;
  right: 0;
}

.main-media {
  position: absolute;
  display: block;
  margin-left: auto;
  margin-right: auto;
  width: auto;
  height: 100%;
  max-width: 100%;
  object-fit: contain;
  opacity: 0;
  transition: opacity 0.2s ease-in-out;
  user-select: none;
  -webkit-user-drag: none;
  pointer-events: none;
}

.main-media.loaded {
  opacity: 1;
  pointer-events: auto;
}

/* Video shows poster immediately, doesn't need to wait for loaded state */
video.main-media {
  opacity: 1;
  pointer-events: auto;
}

/* Audio shouldn't stretch to full height */
audio.main-media {
  height: 50px;
  width: 300px;
}

.media-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  width: 100%;
  position: relative;
  overflow: hidden;
}

.media-slide {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 10;
}

/* Ensure zoom component takes full space */
.media-wrapper :deep(.projection-wrapper) {
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0;
  left: 0;
}

.slide-left-enter-active,
.slide-left-leave-active,
.slide-right-enter-active,
.slide-right-leave-active {
  transition: transform 250ms ease-out, opacity 250ms ease-out;
}

.slide-left-enter-from {
  transform: translateX(100%);
  opacity: 0;
}
.slide-left-leave-to {
  transform: translateX(-100%);
  opacity: 0;
}
.slide-right-enter-from {
  transform: translateX(-100%);
  opacity: 0;
}
.slide-right-leave-to {
  transform: translateX(100%);
  opacity: 0;
}

audio {
  display: block;
  margin-left: auto;
  margin-right: auto;
}

.caption {
  display: block;
  font-size: 1rem;
  white-space: normal;
  text-align: center;
  line-height: 1.4;
  padding-left: 48px;
}

ion-content ion-button[slot='fixed'] {
  top: var(--ion-padding);
  right: var(--ion-padding);
}
</style>
