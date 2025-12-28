<template>
  <ion-modal
    ref="modalRef"
    :is-open="isOpen"
    @did-dismiss="onClose"
  >
    <ion-header
      v-show="!isZoomMode"
      class="ion-no-border"
    >
      <ion-toolbar color="primary">
        <ion-buttons slot="start">
          <ion-button @click="media.isFavourite === false ? $emit('favourite:set') : $emit('favourite:unset')">
            <ion-icon
              slot="icon-only"
              :icon="media.isFavourite === false ? heartOutline : heart"
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

    <ion-content :color="isZoomMode ? 'dark' : 'primary'">
      <ion-button
        v-if="isZoomMode"
        slot="fixed"
        color="light"
        shape="round"
        @click="setZoomMode(false)"
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
        <ion-spinner v-if="!isMediaLoaded" />
        <video
          v-if="media.type === 'video'"
          v-show="isMediaLoaded && !isZoomMode"
          :key="'video-' + media.id"
          :src="fileUrl"
          controls
          :poster="getThumbnailUrl(media.id)"
          @loadeddata="onMediaLoaded"
          class="main-media"
        />
        <audio
          v-else-if="media.type === 'audio'"
          v-show="isMediaLoaded && !isZoomMode"
          :key="'audio-' + media.id"
          :src="fileUrl"
          controls
          @loadeddata="onMediaLoaded"
          class="main-media"
        />
        <img
          v-else-if="media.type === 'image'"
          v-show="isMediaLoaded && !isZoomMode"
          :key="'img-' + media.id"
          :src="fileUrl"
          @load="onMediaLoaded"
          @click="setZoomMode(true)"
          class="main-media"
        />
        <media-zoom
          v-if="isZoomMode"
          :is-open="isZoomMode"
          :file-url="fileUrl"
          @close="setZoomMode(false)"
        />
      </div>
      <div
        v-else
        class="media-wrapper"
      >
        <ion-spinner />
      </div>
    </ion-content>

    <ion-footer
      v-show="!isZoomMode"
      class="ion-no-border"
    >
      <ion-toolbar
        color="primary"
        class="footer"
      >
        <ion-text class="caption">
          <span
            v-for="(line, index) in media.caption.split('\n')"
            :key="index"
          >
            {{ line }}<span v-if="index < media.caption.split('\n').length - 1"><br /></span>
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
  createAnimation,
  IonSpinner,
  IonText,
} from '@ionic/vue';
import { onMounted, onUnmounted, ref, nextTick, computed, watch } from 'vue';
import { createGesture } from '@ionic/vue';
import { heart, close, ellipsisHorizontalCircleOutline, heartOutline } from 'ionicons/icons';
import mediaSelectPopover from './media-select-popover.vue';
import mediaZoom from './media-zoom.vue';
import { useThemeColor } from '@/composables/use-theme-color';
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

const { setThemeColor } = useThemeColor();
const me = useMeStore();
const { getFileUrl, getThumbnailUrl } = useThumbnail();

const modalRef = ref<InstanceType<typeof IonModal> | null>(null);
const mediaWrapperRef = ref<HTMLElement | null>(null);
const gesture = ref<ReturnType<typeof createGesture> | null>(null);
const fileUrl = ref<string | undefined>(undefined);
const isZoomMode = ref(false);
const prevFileUrl = ref<string | null>(null);
const nextFileUrl = ref<string | null>(null);

const isAdmin = me.isAdmin;

const isMediaLoaded = ref(false);

const onMediaLoaded = () => {
  isMediaLoaded.value = true;
};

/**
 * Preloads a media URL into the browser cache.
 */
const preloadMedia = (url: string | null) => {
  if (!url) return;
  const img = new Image();
  img.src = url;
};

watch(
  () => props.media,
  (newMedia, oldMedia) => {
    // Only reset if it's actually a different media file
    if (newMedia?.id === oldMedia?.id && fileUrl.value) return;

    isMediaLoaded.value = false;
    if (!newMedia) return;
    fileUrl.value = getFileUrl(newMedia.id);

    // Set URLs for transitions
    prevFileUrl.value = props.prevMediaId ? getFileUrl(props.prevMediaId) : null;
    nextFileUrl.value = props.nextMediaId ? getFileUrl(props.nextMediaId) : null;

    // Actively preload adjacent items
    preloadMedia(prevFileUrl.value);
    preloadMedia(nextFileUrl.value);
  },
  { immediate: true },
);

const load = async () => {
  // Direct URL is set in watch, no further action needed
};

const onClose = () => {
  fileUrl.value = undefined;
  setZoomMode(false);
  if (modalRef.value) {
    modalRef.value.$el.dismiss(null, 'cancel');
  }
};

const handleKey = (e: KeyboardEvent) => {
  if (e.key === 'ArrowLeft') onPrev();
  if (e.key === 'ArrowRight') onNext();
};

const slideOut = async (direction: 'left' | 'right') => {
  if (!mediaWrapperRef.value) return;

  const anim = createAnimation()
    .addElement(mediaWrapperRef.value)
    .duration(300)
    .fromTo('transform', 'translateX(0)', `translateX(${direction === 'left' ? '-100%' : '100%'})`)
    .fromTo('opacity', '1', '0');
  await anim.play();
  anim.destroy();
};

const onPrev = async () => {
  if (!props.prevMediaId || isZoomMode.value) return;
  await slideOut('right');
  isMediaLoaded.value = false;
  fileUrl.value = undefined; // Trigger loading spinner
  emit('prev', props.media.id);
};

const onNext = async () => {
  if (!props.nextMediaId || isZoomMode.value) return;
  await slideOut('left');
  isMediaLoaded.value = false;
  fileUrl.value = undefined; // Trigger loading spinner
  emit('next', props.media.id);
};

const initSwipe = async () => {
  const gestureEl =
    modalRef.value?.$el?.querySelector('.modal-wrapper, ion-content') || modalRef.value?.$el || modalRef.value;

  gesture.value = createGesture({
    el: gestureEl,
    gestureName: 'swipe',
    direction: 'x',
    threshold: 10,
    onMove(ev) {
      if (ev.deltaX > 50) {
        // stop further detection until next media is shown
        gesture.value?.enable(false);
        onPrev();
      } else if (ev.deltaX < -50) {
        // stop further detection until next media is shown
        gesture.value?.enable(false);
        onNext();
      }
      // avoid immediate re-triggering and skip files
      setTimeout(() => {
        gesture.value?.enable(true);
      }, 200);
    },
  });
  gesture.value.enable(true);
};

const onDelete = () => {
  emit('delete');
  onClose();
};

const setZoomMode = (zoom: boolean) => {
  isZoomMode.value = zoom;
  const theme = zoom ? '--ion-color-dark' : '--ion-color-primary';
  setThemeColor(theme);
};

onMounted(async () => {
  window.addEventListener('keydown', handleKey);
  await nextTick();
  initSwipe();
});

onUnmounted(() => {
  window.removeEventListener('keydown', handleKey);
  if (gesture.value) {
    gesture.value.destroy();
  }
});
</script>

<style scoped>
ion-modal {
  --height: 100%;
  --width: 100%;
  --border-radius: 0;
}

.media-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 100%;
  width: 100%;
}

ion-img,
video,
.main-media {
  display: block;
  margin-left: auto;
  margin-right: auto;
  width: auto;
  height: 100%;
  max-width: 100%;
  object-fit: contain;
}

audio {
  display: block;
  margin-left: auto;
  margin-right: auto;
}

.slide-left-leave-active,
.slide-right-leave-active {
  transition:
    transform 0.3s,
    opacity 0.3s;
  position: absolute;
  width: 100%;
}

.slide-left-leave-to {
  transform: translateX(-100%);
  opacity: 0;
}

.slide-right-leave-to {
  transform: translateX(100%);
  opacity: 0;
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
