import { type ComponentPublicInstance } from 'vue';

declare global {
  // type of IonContent, IonHeader, etc. does have a $el property
  type VueInstanceElement = ComponentPublicInstance & { $el: HTMLElement };

  // responsive view sizes
  type ViewSize = 'xs' | 'sm' | 'md' | 'lg' | 'xl';

  type AppNotification = {
    message: string;
    error: Error | unknown | null;
  };

  type User = {
    id: string;
    name: string;
    email: string;
    isAdmin: boolean;
    hasPublicFavourites?: boolean;
  };

  type Location = {
    lat: number;
    lng: number;
  };

  type UploadMedia = {
    caption: string;
    type: string;
    size: number;
    file: string;
    fileName: string;
    location?: Location;
    date?: string;
    isPublic: boolean;
    orientation?: number;
    rawFile?: globalThis.File;
    status: 'pending' | 'uploading' | 'uploaded' | 'failed';
  };

  type UploadResponse = {
    success: Media[];
    failed: UploadMedia[];
  };

  type Media = {
    id: number;
    isPublic: boolean;
    caption: string;
    date: string; // YYYY-MM-DD
    type: string;
    // location: Location;
    createdAt: string; // datetime
    thumbUrl?: string;
    isFavourite?: boolean;
    isLoaded?: boolean; // for thumbnail loading state
  };

  type UserWithLatestFavourite = {
    user: {
      id: string;
      name: string;
      media_count: number;
    };
    media: {
      id: number;
      caption: string;
      thumbUrl?: string;
      isLoaded?: boolean;
    };
  };

  type FavouritesByUser = {
    user: {
      id: string;
      name: string;
    };
    media: Media[];
  };

  type Album = {
    id: number;
    isPublic: boolean;
    name: string;
    description: string;
    mediaCount: number;
    media: Media[];
  };

  type AlbumForm = {
    id?: number;
    name: string;
    description: string;
    isPublic: boolean;
  };

  type CollectionListItem = {
    id: string;
    name: string;
    description?: string;
    mediaCount: number;
    isPublic?: boolean;
    media: {
      id: number;
      caption: string;
      thumbUrl?: string;
      isLoaded?: boolean;
    };
  };
}

export {};
