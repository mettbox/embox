import { httpService } from '@/services/http.service';

export const FavouriteService = {
  async add(ids: number[]): Promise<void> {
    await httpService('favourite/', 'post', { ids });
  },

  async remove(ids: number[]): Promise<void> {
    await httpService('favourite/', 'delete', { ids });
  },

  async getUsersWithLatestFavourite(): Promise<CollectionListItem[]> {
    const { data } = await httpService('favourite/user', 'get');

    const collection = data.users.map((item: UserWithLatestFavourite) => ({
      id: item.user.id,
      name: item.user.name,
      mediaCount: item.user.media_count,
      media: {
        ...item.media,
        thumbUrl: '',
        isLoaded: false,
      },
    }));

    return collection;
  },

  async getFavouritesByUserID(id: string): Promise<FavouritesByUser> {
    const { data } = await httpService(`favourite/user/${id}`, 'get');

    return {
      user: data.user,
      media: data.media.map((media: Media) => ({
        ...media,
        thumbUrl: '',
        isLoaded: false,
      })),
    };
  },
};
