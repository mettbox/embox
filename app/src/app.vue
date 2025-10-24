<template>
  <ion-app>
    <app-menu />

    <ion-page id="main-content">
      <ion-router-outlet />
    </ion-page>

    <app-notification />
  </ion-app>
</template>

<script setup lang="ts">
import { getMode } from '@ionic/core';
import { HttpError } from '@/services/http.service';
import { IonApp, IonPage, IonRouterOutlet } from '@ionic/vue';
import { onMounted, onUnmounted } from 'vue';
import { useAppStore } from '@/stores/app.store';
import { useI18n } from 'vue-i18n';

import appMenu from '@/components/app/app-menu.vue';
import appNotification from '@/components/app/app-notification.vue';

const app = useAppStore();
const { t } = useI18n();

const onlineHandler = () => {
  app.unsetNotification();
};

const offlineHandler = () => {
  app.setNotification(t('No Internet Connection'), new HttpError('Offline', 0));
};

const disablePinchZoom = (e: Event) => {
  e.preventDefault();
};

onMounted(async () => {
  document.addEventListener('gesturestart', disablePinchZoom);

  app.setAppMode(getMode());
  app.loadAppearance();
  app.updateBreakpoint();

  window.addEventListener('resize', app.updateBreakpoint);

  window.addEventListener('online', () => {
    app.unsetNotification();
  });

  window.addEventListener('offline', () => {
    app.setNotification(t('No Internet Connection'), new HttpError('Offline', 0));
  });
});

// Clean up event listeners when the component is hot-reloaded during development
onUnmounted(() => {
  document.removeEventListener('gesturestart', disablePinchZoom);
  window.removeEventListener('resize', app.updateBreakpoint);
  window.removeEventListener('online', onlineHandler);
  window.removeEventListener('offline', offlineHandler);
});
</script>
