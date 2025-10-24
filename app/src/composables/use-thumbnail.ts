import { MediaService } from '@/services/media.service';

type ThumbnailMedia = Pick<Media, 'id' | 'isLoaded' | 'thumbUrl'>;

export function useThumbnail() {
  const loadThumbnail = async (item: ThumbnailMedia) => {
    if (!item.thumbUrl && !item.isLoaded) {
      item.isLoaded = true;
      try {
        const blob = await MediaService.getThumbnail(item.id);

        if (item.thumbUrl) {
          URL.revokeObjectURL(item.thumbUrl);
        }

        item.thumbUrl = URL.createObjectURL(blob);
      } catch (err) {
        item.thumbUrl = '';
      } finally {
        item.isLoaded = false;
      }
    }
    return item.thumbUrl;
  };

  const releaseThumbnail = (item: any) => {
    if (item.thumbUrl) {
      URL.revokeObjectURL(item.thumbUrl);
      item.thumbUrl = null;
    }
  };

  return {
    loadThumbnail,
    releaseThumbnail,
  };
}
