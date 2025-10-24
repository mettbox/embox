<template>
  <ion-backdrop
    v-if="isOpen && notification.error"
    :tappable="true"
    @ion-backdrop-tap="onDidDismiss"
  />

  <ion-toast
    :is-open="isOpen"
    :message="message"
    :icon="icon"
    :duration="duration"
    :buttons="[
      {
        icon: closeOutline,
        role: 'cancel',
      },
    ]"
    :translucent="true"
    :color="color"
    @did-dismiss="onDidDismiss"
  />
</template>

<script lang="ts" setup>
import { IonBackdrop, IonToast } from '@ionic/vue';
import { warning, informationCircle, closeOutline } from 'ionicons/icons';
import { computed } from 'vue';
import { HttpError } from '@/services/http.service';
import { useAppStore } from '@/stores/app.store';

const app = useAppStore();

const notification = computed(() => app.getNotification);
const isOpen = computed(() => app.hasNotification);

const message = computed(() => {
  let error = '';

  if (notification.value.error) {
    if (notification.value.error instanceof HttpError) {
      error = `${notification.value.error.code}: ${notification.value.error.message}`;
    } else if (notification.value.error instanceof Error) {
      error = notification.value.error.message;
    } else {
      error = String(notification.value.error);
    }
  }

  return error ? `<b>${notification.value.message}</b><br>${error}` : notification.value.message;
});

const isError = computed(
  () => notification.value.error instanceof HttpError || notification.value.error instanceof Error,
);

const icon = computed(() => (isError.value ? warning : informationCircle));

const color = computed(() => (isError.value ? 'danger' : 'light'));

const duration = computed(() => (isError.value ? -1 : 1500));

const onDidDismiss = () => app.unsetNotification();
</script>

<style scoped>
ion-backdrop {
  background: var(--ion-color-dark);
  opacity: 0.3;
}
</style>
