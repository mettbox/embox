<template>
  <ion-modal
    class="media-upload-modal"
    ref="modalRef"
    :is-open="props.isOpen"
    @did-dismiss="onCancel"
  >
    <ion-header>
      <ion-toolbar color="primary">
        <ion-title>{{ $t('Upload') }}</ion-title>
        <ion-buttons slot="end">
          <ion-button @click="onCancel">
            <ion-icon
              slot="icon-only"
              :icon="close"
            />
          </ion-button>
        </ion-buttons>
      </ion-toolbar>
    </ion-header>

    <ion-content>
      <ion-loading
        v-if="uploadStore.isUploading"
        :is-open="true"
        :translucent="true"
        spinner="circular"
        :message="`${uploadStore.uploadProgress} / ${uploadStore.files.length}`"
      />
      <ion-list v-if="uploadStore.files.length > 0">
        <ion-item
          v-for="(file, index) in uploadStore.files"
          :key="file.fileName"
        >
          <div
            v-if="!file.date"
            class="error-indicator-wrapper"
            slot="start"
          >
            <ion-icon
              slot="icon-only"
              :icon="warning"
              color="danger"
              class="error-indicator"
            />
          </div>
          <ion-thumbnail
            slot="start"
            @click="
              uploadStore.files[index] && (uploadStore.files[index].isPublic = !uploadStore.files[index].isPublic)
            "
            :class="{
              'is-public': uploadStore.files[index].isPublic,
            }"
          >
            <ion-img
              v-if="file.type.startsWith('image/')"
              :src="file.file"
            />
            <video
              v-else-if="file.type.startsWith('video/')"
              controls
            >
              <source
                :src="file.file"
                type="video/mp4"
              />
            </video>
            <audio
              v-else-if="file.type.startsWith('audio/')"
              controls
            >
              <source
                :src="file.file"
                type="audio/mpeg"
              />
            </audio>
            <ion-icon
              v-else
              :icon="helpCircle"
              size="large"
            />
            <ion-icon
              v-if="!uploadStore.files[index].isPublic"
              class="non-public-icon"
              :icon="eyeOffOutline"
            />
          </ion-thumbnail>

          <ion-label>
            <ion-textarea
              :rows="1"
              :auto-grow="true"
              :placeholder="file.fileName"
              v-model="uploadStore.files[index].caption"
            />
            <div
              class="media-datepicker"
              :class="{ error: !file.date }"
            >
              <ion-datetime-button :datetime="`datetime-${index}`" />
              <media-datepicker
                :id="`datetime-${index}`"
                v-model="uploadStore.files[index].date"
                @update:all="setDates"
              />
            </div>
            <ion-icon
              v-if="!file.date"
              color="danger"
              :icon="warning"
            />
          </ion-label>

          <ion-button
            slot="end"
            fill="clear"
            @click="onRemove(index)"
          >
            <ion-icon
              slot="icon-only"
              :icon="closeCircleOutline"
            />
          </ion-button>
        </ion-item>
      </ion-list>

      <input
        ref="inputRef"
        type="file"
        class="input-element"
        name="imgUpload"
        :multiple="true"
        accept="image/*,video/*,audio/*"
        @change="onFileChange"
      />
    </ion-content>

    <ion-footer>
      <ion-toolbar
        color="primary"
        class="footer"
      >
        <ion-buttons slot="start">
          <ion-button @click="triggerUploadInput">
            <ion-icon
              slot="icon-only"
              :icon="addCircleOutline"
            />
          </ion-button>
          <template v-if="uploadStore.files.length >= 0">
            <ion-button @click="toggleAllPublic">
              <ion-icon
                slot="icon-only"
                :icon="hasAllPublic ? eyeOutline : eyeOffOutline"
              />
            </ion-button>
            <ion-button @click="hasOpenAlbumModal = true">
              <ion-icon
                slot="icon-only"
                :icon="albumsOutline"
              />
            </ion-button>
          </template>
        </ion-buttons>

        <ion-buttons slot="end">
          <ion-button
            v-if="uploadStore.files.length >= 0"
            @click="onSave"
            :disabled="isInvalid"
          >
            <ion-icon
              slot="icon-only"
              :icon="saveOutline"
            />
          </ion-button>
        </ion-buttons>
      </ion-toolbar>
    </ion-footer>
  </ion-modal>

  <album-select-modal
    :is-open="hasOpenAlbumModal"
    @album:set="onAlbumSelected"
    @did-dismiss="hasOpenAlbumModal = false"
  />
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import {
  IonItem,
  IonList,
  IonTextarea,
  IonThumbnail,
  IonIcon,
  IonLabel,
  IonImg,
  IonButton,
  IonModal,
  IonContent,
  IonHeader,
  IonToolbar,
  IonTitle,
  IonButtons,
  IonFooter,
  IonLoading,
  IonDatetimeButton,
} from '@ionic/vue';
import {
  helpCircle,
  warning,
  close,
  closeCircleOutline,
  saveOutline,
  addCircleOutline,
  eyeOffOutline,
  eyeOutline,
  albumsOutline,
} from 'ionicons/icons';
import mediaDatepicker from '@/components/media/media-datepicker.vue';
import * as exifr from 'exifr';
import { useUploadStore } from '@/stores/upload.store';
import albumSelectModal from '@/components/album/album-select-modal.vue';

const uploadStore = useUploadStore();

const props = defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: 'did-upload', response: UploadResponse): void;
  (e: 'did-dismiss'): void;
  (e: 'cancel'): void;
}>();

const isInvalid = computed(() => {
  return uploadStore.files.length === 0 || uploadStore.files.some((f) => !f.date);
});

const modalRef = ref<InstanceType<typeof IonModal> | null>(null);
const inputRef = ref<HTMLInputElement | null>(null);
const isUploading = ref(false);
const hasAllPublic = ref(false);
const hasOpenAlbumModal = ref(false);

const onFileChange = async (e: Event) => {
  const fileList = (e.target as HTMLInputElement).files;
  if (!fileList) return;

  const newFiles: UploadMedia[] = [];

  for (const file of Array.from(fileList)) {
    const alreadyExists = uploadStore.files.some((m) => m.size === file.size && m.fileName === file.name);
    if (alreadyExists) continue;

    let location;
    let date;
    let orientation;

    if (file.type.startsWith('image/')) {
      try {
        const exif = await exifr.parse(file, {
          gps: true,
          tiff: true,
          exif: true,
        });
        if (exif) {
          if (exif.latitude && exif.longitude) {
            location = { lat: exif.latitude, lng: exif.longitude };
          }
          if (exif.DateTimeOriginal) {
            const d = new Date(exif.DateTimeOriginal);
            date = d.toISOString().slice(0, 10); // yyyy-mm-dd
          }
          if (exif.Orientation) {
            orientation = exif.Orientation;
          }
        }
      } catch (err) {
        console.warn('EXIF extraction failed', err);
      }
    }

    newFiles.push({
      caption: '',
      type: file.type,
      size: file.size,
      file: window.URL.createObjectURL(file),
      fileName: file.name,
      date,
      isPublic: true,
      location,
      orientation,
      rawFile: file, // Store the original File object for upload
      status: 'pending',
    });
  }

  uploadStore.addFiles(newFiles);

  // Reset input to allow re-selection of the same file
  (e.target as HTMLInputElement).value = '';
};

const triggerUploadInput = () => {
  if (inputRef.value) {
    inputRef.value.click();
  }
};

const onCancel = async () => {
  uploadStore.clearFiles();
  isUploading.value = false;
  if (inputRef.value) {
    inputRef.value.value = '';
  }

  if (modalRef.value) {
    modalRef.value.$el.dismiss(null, 'cancel');
  }
};

const onRemove = (index: number) => {
  uploadStore.files.splice(index, 1);
};

const toggleAllPublic = () => {
  hasAllPublic.value = !hasAllPublic.value;
  uploadStore.files.forEach((f) => {
    f.isPublic = hasAllPublic.value;
  });
};

// set same date for all files w/o date
const setDates = (date: string) => {
  uploadStore.files.forEach((f) => {
    if (!f.date) {
      f.date = date;
    }
  });
};

const onAlbumSelected = (albumId: number) => {
  uploadStore.setAlbum(albumId);
  hasOpenAlbumModal.value = false;
};

const onSave = async () => {
  const { success, failed } = await uploadStore.uploadFiles();
  emit('did-upload', { success, failed });
  onCancel();
};

defineExpose({
  triggerUploadInput,
});
</script>

<style>
.media-upload-modal {
  .input-element {
    position: fixed;
    top: -500px;
    left: -500;
    z-index: -100;
    opacity: 0;
  }

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

    ion-img,
    video,
    audio {
      opacity: 0.6;
    }

    &.is-public {
      ion-img,
      video,
      audio {
        opacity: 1;
      }
    }

    .non-public-icon {
      position: absolute;
      bottom: 4px;
      right: 4px;
      font-size: 24px;
    }
  }

  .error-indicator {
    position: absolute;
    top: 12px;
    width: 10px;
    height: 10px;
    inset-inline-start: 12px;
  }

  .media-datepicker {
    display: inline-block;
    margin-top: 8px;
  }

  .media-datepicker.error {
    opacity: 0.5;
  }
}
</style>
