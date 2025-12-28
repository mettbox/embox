import { httpService } from '@/services/http.service';

export const AlbumService = {
  async create(album: AlbumForm): Promise<Album> {
    const { data } = await httpService('album/', 'post', album);
    return data;
  },

  async update(id: number, album: AlbumForm): Promise<Album> {
    const { data } = await httpService(`album/${id}`, 'put', album);
    return data;
  },

  async get(): Promise<CollectionListItem[]> {
    const { data } = await httpService('album/', 'get');

    if (!data || data.length === 0) {
      return [];
    }

    const collection = data.map((item: Album) => {
      const media =
        item.media && item.media.length > 0
          ? {
              ...item.media[0],
              thumbUrl: '',
              isLoaded: false,
            }
          : null;

      return {
        id: `${item.id}`,
        name: item.name,
        description: item.description,
        mediaCount: item.mediaCount,
        media,
      };
    });

    return collection;
  },

  async getByID(id: number): Promise<Album> {
    const { data } = await httpService(`album/${id}`, 'get');

    return data;
  },

  async addMedia(albumId: number, mediaIds: number[]): Promise<void> {
    await httpService(`album/${albumId}/media`, 'post', { mediaIds });
  },

  async removeMedia(albumId: number, mediaIds: number[]): Promise<void> {
    await httpService(`album/${albumId}/media`, 'delete', { mediaIds });
  },

  async delete(albumId: number): Promise<void> {
    await httpService(`album/${albumId}`, 'delete');
  },
};
