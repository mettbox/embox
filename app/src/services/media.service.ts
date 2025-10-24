import { httpService } from '@/services/http.service';

export const MediaService = {
  async getMediaList(): Promise<Media[]> {
    const { data } = await httpService('media/', 'get');
    return data;
  },

  async getMediaFile(id: number): Promise<Blob> {
    return await httpService(`media/${id}/file`, 'get', {}, { responseType: 'blob' });
  },

  async getThumbnail(id: number): Promise<Blob> {
    return await httpService(`media/${id}/thumbnail`, 'get', {}, { responseType: 'blob' });
  },

  async upload(files: UploadMedia[]): Promise<any> {
    const formData = new FormData();

    const meta = files.map((file) => ({
      caption: file.caption,
      type: file.type, // e.g "image/jpeg"
      fileName: file.fileName, // original file name
      isPublic: file.isPublic,
      date: file.date, // yyyy-mm-dd
      orientation: file.orientation ?? '',
      locationLat: file.location?.lat ?? 0,
      locationLng: file.location?.lng ?? 0,
    }));
    formData.append('meta', JSON.stringify(meta));

    files.forEach((file) => {
      if (file.rawFile) {
        formData.append('files', file.rawFile, file.fileName);
      }
    });

    return await httpService('media/', 'post', formData, {}, 600000); // 10 minutes timeout
  },

  async setPublic(ids: number[], isPublic: boolean): Promise<void> {
    const updates = ids.map((id) => ({
      id,
      isPublic,
    }));
    return this.update(updates);
  },

  async update(updates: Partial<Media>[]): Promise<void> {
    const { data } = await httpService('media/', 'put', { updates });
    return data;
  },

  async delete(ids: number[]): Promise<void> {
    await httpService('media/', 'delete', { ids });
  },
};
