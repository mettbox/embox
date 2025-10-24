<template>
  <ion-page class="page page-login">
    <app-header>
      <app-logo />
    </app-header>

    <ion-content class="ion-padding">
      <v-center-container>
        <h1>{{ $t('Login') }}</h1>
        <p>
          {{ $t('Please enter your login code.') }}
          <br />{{ $t('Please note that it is only valid once within the next 15 minutes.') }}
        </p>

        <ion-input-otp
          shape="soft"
          type="number"
          :length="6"
          @ion-complete="onPinEntered"
        >
          <ion-button
            class="ion-margin-top"
            fill="clear"
            size="small"
            @click="router.replace({ name: 'login' })"
            >{{ $t('Request a new code') }}</ion-button
          >
        </ion-input-otp>
      </v-center-container>
    </ion-content>
  </ion-page>
</template>

<script setup lang="ts">
import { IonButton, IonContent, IonPage, IonInputOtp } from '@ionic/vue';
import vCenterContainer from '@/components/app/app-center-container.vue';
import appLogo from '@/components/app/app-logo.vue';
import { HttpError } from '@/services/http.service';
import { useRouter } from 'vue-router';
import { useAppStore } from '@/stores/app.store';
import { useMeStore } from '@/stores/me.store';
import { useI18n } from 'vue-i18n';
import appHeader from '@/components/app/app-header.vue';

const { t } = useI18n();
const router = useRouter();
const app = useAppStore();
const me = useMeStore();

const onPinEntered = async (e: CustomEvent): Promise<void> => {
  try {
    await me.login(e.detail.value);
    router.replace({ name: 'library' });
  } catch (error: unknown) {
    if (error instanceof HttpError && (error.code === 410 || error.code === 429)) {
      router.replace('/login');
    }
    app.setNotification(t('Invalid login token'), error);
  }
};
</script>

<style scoped>
.page-login {
  ion-grid {
    height: 100%;
    display: flex;
    flex-direction: column;
  }

  ion-row {
    flex: 1;
  }
}
</style>
