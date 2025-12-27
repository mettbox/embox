<template>
  <ion-button
    v-if="nonPublicOnly"
    shape="round"
    color="primary"
    @click="emit('update:non-public-only', false)"
  >
    <ion-icon
      slot="icon-only"
      :icon="eyeOff"
    />
  </ion-button>
  <ion-button
    id="filter-popover"
    shape="round"
    :color="isActive ? 'primary' : 'dark'"
  >
    <ion-icon
      slot="icon-only"
      :icon="buttonIcon"
    />
  </ion-button>

  <ion-popover
    trigger="filter-popover"
    alignment="center"
    :dismiss-on-select="true"
  >
    <ion-content>
      <ion-list lines="none">
        <ion-radio-group
          :value="sort"
          @ionChange="emit('update:sort', $event.detail.value)"
        >
          <ion-item>
            <ion-radio
              justify="space-between"
              value="asc"
              >{{ $t('Date asc') }}</ion-radio
            >
          </ion-item>
          <ion-item lines="full">
            <ion-radio
              justify="space-between"
              value="desc"
              >{{ $t('Date desc') }}</ion-radio
            >
          </ion-item>
        </ion-radio-group>
      </ion-list>

      <ion-list
        v-if="isAdmin"
        lines="full"
      >
        <ion-item
          :button="true"
          :detail="false"
          @click="emit('update:non-public-only', !nonPublicOnly)"
        >
          <ion-icon
            slot="end"
            size="small"
            :color="nonPublicOnly ? 'primary' : 'medium'"
            :icon="nonPublicOnly ? eyeOff : eyeOffOutline"
          />
          <ion-label>{{ $t('Non public only') }}</ion-label>
        </ion-item>
      </ion-list>
      <ion-list lines="none">
        <ion-item
          :button="true"
          :detail="false"
          @click="emit('update:filter', '')"
        >
          <ion-icon
            slot="end"
            size="small"
            :color="filter === '' ? 'primary' : 'medium'"
            :icon="filter === '' ? apps : appsOutline"
          />
          <ion-label>{{ $t('All Media') }}</ion-label>
        </ion-item>
        <ion-item
          :button="true"
          :detail="false"
          @click="emit('update:filter', 'favourites')"
        >
          <ion-icon
            slot="end"
            size="small"
            :icon="filter === 'favourites' ? heart : heartOutline"
            :color="filter === 'favourites' ? 'primary' : 'medium'"
          />
          <ion-label>{{ $t('Favourites') }}</ion-label>
        </ion-item>
        <ion-item
          :button="true"
          :detail="false"
          @click="emit('update:filter', 'image')"
        >
          <ion-icon
            slot="end"
            size="small"
            :icon="filter === 'image' ? images : imagesOutline"
            :color="filter === 'image' ? 'primary' : 'medium'"
          />
          <ion-label>{{ $t('Photos') }}</ion-label>
        </ion-item>
        <ion-item
          :button="true"
          :detail="false"
          @click="emit('update:filter', 'video')"
        >
          <ion-icon
            slot="end"
            size="small"
            :icon="filter === 'video' ? videocam : videocamOutline"
            :color="filter === 'video' ? 'primary' : 'medium'"
          />
          <ion-label>{{ $t('Videos') }}</ion-label>
        </ion-item>
        <ion-item
          :button="true"
          :detail="false"
          @click="emit('update:filter', 'audio')"
        >
          <ion-icon
            slot="end"
            size="small"
            :color="filter === 'audio' ? 'primary' : 'medium'"
            :icon="filter === 'audio' ? musicalNotes : musicalNotesOutline"
          />
          <ion-label>{{ $t('Audios') }}</ion-label>
        </ion-item>
      </ion-list>
    </ion-content>
  </ion-popover>
</template>

<script setup lang="ts">
import {
  IonContent,
  IonIcon,
  IonPopover,
  IonList,
  IonItem,
  IonLabel,
  IonButton,
  IonRadioGroup,
  IonRadio,
} from '@ionic/vue';
import {
  heart,
  heartOutline,
  swapVertical,
  appsOutline,
  images,
  imagesOutline,
  videocam,
  videocamOutline,
  musicalNotesOutline,
  musicalNotes,
  apps,
  eyeOffOutline,
  eyeOff,
} from 'ionicons/icons';
import { computed } from 'vue';
import { useMeStore } from '@/stores/me.store';

const props = defineProps<{
  sort: 'asc' | 'desc';
  filter: '' | 'image' | 'video' | 'audio' | 'favourites';
  nonPublicOnly: boolean;
}>();

const emit = defineEmits<{
  (e: 'update:sort', value: 'asc' | 'desc'): void;
  (e: 'update:filter', value: '' | 'image' | 'video' | 'audio' | 'favourites'): void;
  (e: 'update:non-public-only', value: boolean): void;
}>();

const isActive = computed(() => {
  return props.filter !== '' || props.sort !== 'desc' || props.nonPublicOnly;
});

const buttonIcon = computed(() => {
  if (props.filter === 'favourites') return heart;
  if (props.filter === 'image') return images;
  if (props.filter === 'video') return videocam;
  if (props.filter === 'audio') return musicalNotes;
  return swapVertical;
});

const me = useMeStore();
const isAdmin = me.isAdmin;
</script>

<style scoped>
ion-radio::part(label) {
  white-space: normal;
}
</style>
