<template>
  <ion-page class="page page-login-request">
    <app-header>
      <app-logo />
    </app-header>

    <ion-content class="ion-padding">
      <v-center-container>
        <h1>{{ $t('Login') }}</h1>
        <p v-html="$t('registration_intro')" />

        <ion-input
          ref="emailInput"
          v-model="form.email"
          class="ion-text-left"
          fill="outline"
          :placeholder="$t('Enter your email address…')"
          type="email"
          aria-label="Email"
          :error-text="isTouched.email ? fieldErrors.email : ''"
          @ion-blur="() => validateField('email')"
          @keyup.enter="request"
        />

        <ion-button
          class="ion-margin-top"
          expand="block"
          @click="request"
          >{{ $t('Submit') }}</ion-button
        >

        <ion-button
          class="ion-margin-top"
          fill="clear"
          size="small"
          @click="router.replace({ name: 'token' })"
          >{{ $t('I already have a login code') }}</ion-button
        >
      </v-center-container>
    </ion-content>
  </ion-page>
</template>

<script setup lang="ts">
import { IonButton, IonContent, IonInput, IonPage } from '@ionic/vue';
import { useRoute } from 'vue-router';
import { useI18n } from 'vue-i18n';
import { ref, onMounted, reactive } from 'vue';
import { useRouter } from 'vue-router';
import { useAppStore } from '@/stores/app.store';
import { useMeStore } from '@/stores/me.store';
import { useValidation, type FormField } from '@/composables/use-validation';
import vCenterContainer from '@/components/app/app-center-container.vue';
import appLogo from '@/components/app/app-logo.vue';
import appHeader from '@/components/app/app-header.vue';

const { t } = useI18n();
const route = useRoute();
const router = useRouter();
const app = useAppStore();
const me = useMeStore();

const emailInput = ref<InstanceType<typeof IonInput> | null>(null);

const form = reactive({
  email: '',
});

const fields: Record<string, FormField> = {
  email: {
    ref: emailInput,
    rules: [
      {
        validate: (val) => val.trim() !== '',
        message: t('Please enter an email address'),
      },
      {
        validate: (val) => /^[^\s@]+@[^\s@]+\.[^\s@]{2,}$/.test(String(val ?? '')),
        message: t('Please enter a valid email address'),
      },
    ],
  },
};

const { validateField, validateAllFields, isFormValid, fieldErrors, isTouched } = useValidation(fields, form);

const request = async () => {
  validateAllFields();
  if (!isFormValid.value) {
    return;
  }

  try {
    await me.requestToken(form.email);
    router.replace('/verify-token');
  } catch (error: unknown) {
    console.log('Login request error:', error);

    app.setNotification(t('Failed to login'), error);
  }
};

onMounted(() => {
  const user = route.params.email ? String(route.params.email) : null;
  if (user) {
    form.email = user;
  }
});
</script>
