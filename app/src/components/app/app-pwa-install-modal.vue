<template>
  <ion-modal
    ref="modalRef"
    :is-open="isOpen"
    @did-dismiss="isOpen = false"
  >
    <ion-header>
      <ion-toolbar color="primary">
        <ion-title>{{ $t('pwa_install_title') }}</ion-title>
        <ion-buttons slot="end">
          <ion-button @click="onClose()">
            <ion-icon
              slot="icon-only"
              :icon="close"
            />
          </ion-button>
        </ion-buttons>
      </ion-toolbar>
    </ion-header>

    <ion-content class="ion-padding">
      <p>{{ $t('pwa_install_intro') }}</p>

      <template v-if="isIos">
        <ion-list lines="none">
          <ion-item>
            <ion-icon
              :icon="shareOutline"
              slot="start"
              color="primary"
            />
            <ion-label class="ion-text-wrap">{{ $t('pwa_install_ios_step1') }}</ion-label>
          </ion-item>
          <ion-item>
            <ion-icon
              :icon="addCircleOutline"
              slot="start"
              color="primary"
            />
            <ion-label class="ion-text-wrap">{{ $t('pwa_install_ios_step2') }}</ion-label>
          </ion-item>
        </ion-list>
      </template>

      <template v-else>
        <ion-list lines="none">
          <ion-item>
            <ion-icon
              :icon="ellipsisVertical"
              slot="start"
              color="primary"
            />
            <ion-label class="ion-text-wrap">{{ $t('pwa_install_android_step1') }}</ion-label>
          </ion-item>
          <ion-item>
            <ion-icon
              :icon="addCircleOutline"
              slot="start"
              color="primary"
            />
            <ion-label class="ion-text-wrap">{{ $t('pwa_install_android_step2') }}</ion-label>
          </ion-item>
        </ion-list>
      </template>
    </ion-content>

    <ion-footer>
      <ion-toolbar>
        <ion-button
          expand="block"
          fill="clear"
          @click="dismissPermanently"
        >
          {{ $t('pwa_install_dismiss') }}
        </ion-button>
      </ion-toolbar>
    </ion-footer>
  </ion-modal>
</template>

<script setup lang="ts">
import {
  IonButton,
  IonContent,
  IonFooter,
  IonHeader,
  IonIcon,
  IonItem,
  IonLabel,
  IonList,
  IonModal,
  IonTitle,
  IonToolbar,
  IonButtons,
} from '@ionic/vue';
import { addCircleOutline, ellipsisVertical, shareOutline, close } from 'ionicons/icons';
import { onMounted, ref } from 'vue';

const modalRef = ref<InstanceType<typeof IonModal> | null>(null);
const isOpen = ref(false);
const isIos = ref(false);

onMounted(() => {
  const ua = navigator.userAgent;
  const iosDevice = /iphone|ipad|ipod/i.test(ua) || (/macintosh/i.test(ua) && navigator.maxTouchPoints > 1);
  const androidDevice = /android/i.test(ua);
  const standalone =
    window.matchMedia('(display-mode: standalone)').matches || (window.navigator as any).standalone === true;
  const dismissed = localStorage.getItem('pwa-install-dismissed') === 'true';

  isIos.value = iosDevice;
  isOpen.value = (iosDevice || androidDevice) && !standalone && !dismissed;
});

const onClose = () => {
  if (modalRef.value) {
    modalRef.value.$el.dismiss(null, 'cancel');
  }
};

const dismissPermanently = () => {
  localStorage.setItem('pwa-install-dismissed', 'true');
  isOpen.value = false;
};
</script>
