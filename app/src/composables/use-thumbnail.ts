const baseURL = import.meta.env.VITE_API_URL || '';

export function useThumbnail() {
  const getThumbnailUrl = (id: number) => {
    return `${baseURL}/media/${id}/thumbnail`;
  };

  const getFileUrl = (id: number) => {
    return `${baseURL}/media/${id}/file`;
  };

  return {
    getThumbnailUrl,
    getFileUrl,
  };
}
