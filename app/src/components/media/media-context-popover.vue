<template>
  <ion-popover
    ref="popoverRef"
    class="media-context-popover"
    :is-open="props.isOpen"
    alignment="center"
    :dismiss-on-select="false"
    @did-dismiss="emit('close')"
  >
    <ion-content>
      <ion-list lines="none">
        <ion-item
          v-if="!selectedMedia.hasFavouriteOnly"
          :button="true"
          :detail="false"
          lines="full"
          @click="$emit('favourite:set')"
        >
          <ion-icon
            slot="end"
            size="small"
            :icon="heartOutline"
            color="medium"
          />
          <ion-label>{{ $t('Favourite') }}</ion-label>
        </ion-item>

        <ion-item
          v-if="!selectedMedia.hasNonFavouriteOnly"
          :button="true"
          :detail="false"
          lines="full"
          @click="$emit('favourite:unset')"
        >
          <ion-icon
            slot="end"
            size="small"
            :icon="heartDislikeOutline"
            color="medium"
          />
          <ion-label>{{ $t('No Favourite') }}</ion-label>
        </ion-item>

        <template v-if="isAdmin">
          <ion-item
            v-if="!selectedMedia.isAlbum"
            :button="true"
            :detail="false"
            lines="full"
            @click="$emit('album:set')"
          >
            <ion-icon
              slot="end"
              size="small"
              :icon="albumsOutline"
              color="medium"
            />
            <ion-label>{{ $t('Add to Album') }}</ion-label>
          </ion-item>

          <ion-item
            v-if="selectedMedia.isAlbum"
            :button="true"
            :detail="false"
            lines="full"
            @click="$emit('album:unset')"
          >
            <ion-icon
              slot="end"
              size="small"
              :icon="trashBinOutline"
              color="medium"
            />
            <ion-label>{{ $t('Remove from Album') }}</ion-label>
          </ion-item>

          <ion-item
            v-if="!selectedMedia.hasPublicOnly"
            :button="true"
            :detail="false"
            lines="full"
            @click="$emit('public:set')"
          >
            <ion-icon
              slot="end"
              size="small"
              :icon="eyeOutline"
              color="medium"
            />
            <ion-label>{{ $t('Public') }}</ion-label>
          </ion-item>

          <ion-item
            v-if="!selectedMedia.hasNonPublicOnly"
            :button="true"
            :detail="false"
            lines="full"
            @click="$emit('public:unset')"
          >
            <ion-icon
              slot="end"
              size="small"
              :icon="eyeOffOutline"
              color="medium"
            />
            <ion-label>{{ $t('Non Public') }}</ion-label>
          </ion-item>

          <ion-item
            :button="true"
            :detail="false"
            lines="full"
            @click="$emit('edit')"
          >
            <ion-icon
              slot="end"
              size="small"
              :icon="createOutline"
              color="medium"
            />
            <ion-label>{{ $t('Edit') }}</ion-label>
          </ion-item>

          <ion-item
            v-if="!selectedMedia.isAlbum"
            :button="true"
            :detail="false"
            lines="full"
            @click="hasOpenDelete = true"
          >
            <ion-icon
              slot="end"
              size="small"
              :icon="trashBinOutline"
              color="medium"
            />
            <ion-label>{{ $t('Delete') }}</ion-label>
          </ion-item>
        </template>
      </ion-list>
    </ion-content>
  </ion-popover>

  <ion-action-sheet
    :is-open="hasOpenDelete"
    :header="$t('This action is permanent and cannot be undone.')"
    :buttons="deleteActionSheetButtons"
    @did-dismiss="deleteHandler"
  />
</template>

<script setup lang="ts">
import { computed, ref } from 'vue';
import { IonPopover, IonContent, IonList, IonItem, IonIcon, IonLabel } from '@ionic/vue';
import {
  heartOutline,
  heartDislikeOutline,
  albumsOutline,
  trashBinOutline,
  createOutline,
  eyeOffOutline,
  eyeOutline,
} from 'ionicons/icons';
import { useI18n } from 'vue-i18n';
import { ActionSheetButton, IonActionSheet } from '@ionic/vue';
import { useSelectedMediaStore } from '@/stores/selectedMedia.store';
import { useMeStore } from '@/stores/me.store';

const props = defineProps<{
  isOpen: boolean;
}>();

const emit = defineEmits<{
  (e: 'favourite:set'): void;
  (e: 'favourite:unset'): void;
  (e: 'album:set'): void;
  (e: 'album:unset'): void;
  (e: 'public:set'): void;
  (e: 'public:unset'): void;
  (e: 'edit'): void;
  (e: 'delete'): void;
  (e: 'close'): void;
}>();

const deleteActionSheetButtons = computed((): ActionSheetButton[] => [
  {
    text: selectedMedia.ids.length === 1 ? t('Delete') : t('Remove Selection', { count: selectedMedia.ids.length }),
    role: 'destructive',
  },
  {
    text: t('Cancel'),
    role: 'cancel',
  },
]);

const { t } = useI18n();
const me = useMeStore();
const isAdmin = me.isAdmin;
const selectedMedia = useSelectedMediaStore();
const hasOpenDelete = ref(false);

const popoverRef = ref<InstanceType<typeof IonPopover> | null>(null);

const deleteHandler = async (event: { detail: { role: string } }): Promise<void> => {
  hasOpenDelete.value = false;
  const action = event?.detail?.role || 'cancel';

  if (action !== 'destructive' || selectedMedia.ids.length === 0) {
    return;
  }

  emit('delete');
};
</script>
