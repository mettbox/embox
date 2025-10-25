import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { MediaService } from '@/services/media.service';
import { AlbumService } from '@/services/album.service';

export const useUploadStore = defineStore('upload', () => {
  const files = ref<UploadMedia[]>([]);
  const albumId = ref<number | null>(null);
  const isUploading = ref(false);
  const uploadProgress = ref(0);
  const batchSize = 6;
  const responses = ref<Record<string, any>>({});

  /**
   * Check if all files are uploaded
   *
   * @return {boolean}
   */
  const isComplete = computed(() => files.value.every((file) => file.status === 'uploaded'));

  /**
   * Add new files to the upload queue
   *
   * @param {UploadMedia[]} newFiles
   * @return {void}
   */
  const addFiles = (newFiles: UploadMedia[]): void => {
    newFiles.forEach((file) => {
      files.value.push({ ...file, status: 'pending' }); // Status: 'pending', 'uploading', 'uploaded', 'failed'
    });
  };

  const setAlbum = (id: number): void => {
    albumId.value = id;
  };

  /**
   * Upload files in batches
   *
   * @returns {Promise<{ success: Media[]; failed: UploadMedia[] }>}
   */
  const uploadFiles = async (): Promise<{ success: Media[]; failed: UploadMedia[] }> => {
    // Avoid multiple simultaneous uploads
    if (isUploading.value) {
      return { success: [], failed: [] };
    }

    isUploading.value = true;
    uploadProgress.value = 0;

    const uploadedFiles: Media[] = [];
    const failedFiles: UploadMedia[] = [];

    try {
      while (files.value.some((file) => file.status === 'pending')) {
        // Take the next batch of files
        const batch = files.value.filter((file) => file.status === 'pending').slice(0, batchSize);

        batch.forEach((file) => (file.status = 'uploading'));

        const uploadPromises = batch.map((file) =>
          MediaService.upload([file])
            .then((response) => {
              file.status = 'uploaded';
              responses.value[file.fileName] = response.data[0];
              uploadedFiles.push(response.data[0]);
            })
            .catch(() => {
              file.status = 'failed';
              failedFiles.push(file);
            })
            .finally(() => {
              uploadProgress.value += 1;
            }),
        );

        await Promise.all(uploadPromises);
      }
      if (albumId.value) {
        const mediaIds = Object.values(responses.value).map((media: Media) => media.id);
        await AlbumService.addMedia(albumId.value, mediaIds);
      }
    } finally {
      isUploading.value = false;
      albumId.value = null;
    }

    return { success: uploadedFiles, failed: failedFiles };
  };

  /**
   * Clear all files from the upload queue
   *
   * @return {void}
   */
  const clearFiles = (): void => {
    files.value = [];
    uploadProgress.value = 0;
  };

  return {
    files,
    isUploading,
    uploadProgress,
    isComplete,
    addFiles,
    uploadFiles,
    clearFiles,
    setAlbum,
  };
});
