import { defineStore } from 'pinia';

export const useSelectedMediaStore = defineStore('selectedMedia', {
  state: () => ({
    selectedMedias: [] as Media[],
    selectedAlbumId: null as number | null,
  }),
  getters: {
    ids: (state) => state.selectedMedias.map((media) => media.id),
    all: (state) => state.selectedMedias,
    first: (state) => (state.selectedMedias.length > 0 ? state.selectedMedias[0] : null),

    hasPublicOnly: (state) => state.selectedMedias.length > 0 && state.selectedMedias.every((media) => media.isPublic),
    hasNonPublicOnly: (state) =>
      state.selectedMedias.length > 0 && state.selectedMedias.every((media) => !media.isPublic),

    hasFavouriteOnly: (state) =>
      state.selectedMedias.length > 0 && state.selectedMedias.every((media) => media.isFavourite),
    hasNonFavouriteOnly: (state) =>
      state.selectedMedias.length > 0 && state.selectedMedias.every((media) => !media.isFavourite),

    isAlbum: (state) => state.selectedAlbumId !== null,
    albumId: (state) => state.selectedAlbumId,
  },
  actions: {
    set(media: Media) {
      this.selectedMedias = [media];
    },
    add(media: Media) {
      if (!this.selectedMedias.find((m) => m.id === media.id)) {
        this.selectedMedias.push(media);
      }
    },
    remove(mediaId: number) {
      this.selectedMedias = this.selectedMedias.filter((m) => m.id !== mediaId);
    },
    clear() {
      this.selectedMedias = [];
    },
    setAlbumId(albumId: number | null) {
      this.selectedAlbumId = albumId;
    },
    reset() {
      this.selectedMedias = [];
      this.selectedAlbumId = null;
    },
  },
});
