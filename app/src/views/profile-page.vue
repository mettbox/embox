<template>
  <ion-page>
    <app-header>
      {{ $t('Profile') }}
    </app-header>

    <ion-content>
      <ion-list class="content-list">
        <ion-item-divider>
          <ion-label>{{ $t('Appearance') }}</ion-label>
        </ion-item-divider>
        <ion-radio-group
          :value="appearance"
          @ionChange="onAppearanceChange"
        >
          <ion-item>
            <ion-radio
              value="dark"
              justify="space-between"
            >
              <ion-icon
                :icon="moon"
                class="ion-padding-end"
              />
              {{ $t('Dark Mode') }}
            </ion-radio>
          </ion-item>
          <ion-item>
            <ion-radio
              value="light"
              justify="space-between"
            >
              <ion-icon
                :icon="sunny"
                class="ion-padding-end"
              />
              {{ $t('Light Mode') }}
            </ion-radio>
          </ion-item>
          <ion-item>
            <ion-radio
              value="system"
              justify="space-between"
            >
              <ion-icon
                :icon="contrastOutline"
                class="ion-padding-end"
              />
              {{ $t('System') }}
            </ion-radio>
          </ion-item>
        </ion-radio-group>

        <ion-item-divider>
          <ion-label>{{ $t('Favourites') }}</ion-label>
        </ion-item-divider>
        <ion-item>
          <ion-toggle
            justify="space-between"
            :checked="me.hasPublicFavourites"
            @ionChange="onTogglePublicFavourites"
          >
            {{ $t('My favourites are public') }}
          </ion-toggle>
        </ion-item>

        <ion-item-divider>
          <ion-label>{{ $t('Profile') }}</ion-label>
        </ion-item-divider>
        <ion-item>
          <ion-input
            ref="nameInput"
            type="text"
            v-model="form.name"
            label-placement="stacked"
            :label="$t('Name')"
            :placeholder="$t('Enter your name…')"
            :error-text="isTouched.name ? fieldErrors.name : ''"
            @ion-blur="() => validateField('name')"
            @keyup.enter="updateProfile"
          />
        </ion-item>
        <ion-item>
          <ion-input
            ref="emailInput"
            type="email"
            v-model="form.email"
            label-placement="stacked"
            :label="$t('Email address')"
            :placeholder="$t('Enter your email address…')"
            :error-text="isTouched.email ? fieldErrors.email : ''"
            :helper-text="$t('You will be logged out and need to login again after changing your email address.')"
            @ion-blur="() => validateField('email')"
            @keyup.enter="updateProfile"
          />
        </ion-item>
        <ion-item class="ion-margin-top">
          <ion-button
            v-if="!isAdmin"
            size="default"
            slot="end"
            color="danger"
            @click="hasOpenDelete = true"
            >{{ $t('Delete Account') }}</ion-button
          >
          <ion-button
            size="default"
            slot="end"
            @click="updateProfile"
            >{{ $t('Update Profil') }}</ion-button
          >
        </ion-item>
      </ion-list>

      <ion-action-sheet
        :is-open="hasOpenDelete"
        :header="$t('Are you sure you want to delete your profile?')"
        :buttons="deleteActionSheetButtons"
        @did-dismiss="deleteHandler"
      />
    </ion-content>
  </ion-page>
</template>

<script setup lang="ts">
import { ref, computed, reactive } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  IonActionSheet,
  IonContent,
  IonInput,
  IonPage,
  onIonViewDidEnter,
  IonList,
  IonToggle,
  IonItem,
  IonRadioGroup,
  IonRadio,
  IonLabel,
  IonIcon,
  IonButton,
  IonItemDivider,
} from '@ionic/vue';
import { ActionSheetButton } from '@ionic/core';
import { useValidation, type FormField } from '@/composables/use-validation';
import { useAppStore } from '@/stores/app.store';
import { useMeStore } from '@/stores/me.store';
import appHeader from '@/components/app/app-header.vue';
import { moon, sunny, contrastOutline } from 'ionicons/icons';

const deleteActionSheetButtons = computed((): ActionSheetButton[] => [
  {
    text: t('Yes, delete my all my data'),
    role: 'destructive',
  },
  {
    text: t('Cancel'),
    role: 'cancel',
  },
]);

const { t } = useI18n();
const app = useAppStore();
const me = useMeStore();
const isAdmin = me.isAdmin;

const nameInput = ref<InstanceType<typeof IonInput> | null>(null);
const emailInput = ref<InstanceType<typeof IonInput> | null>(null);
const hasOpenDelete = ref(false);
const appearance = ref(localStorage.getItem('appearance') || 'system');

const onAppearanceChange = (e: CustomEvent) => {
  appearance.value = e.detail.value;
  app.setAppearance(appearance.value as 'dark' | 'light' | 'system');
};

const form = reactive({
  name: me.name || '',
  email: me.email || '',
});

const fields: Record<string, FormField> = {
  name: {
    ref: nameInput,
    rules: [
      {
        validate: (val) => val.trim() !== '',
        message: t('Please enter your name'),
      },
    ],
  },
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

const onTogglePublicFavourites = async (event: CustomEvent) => {
  const newValue = event.detail.checked;

  try {
    await me.update({ hasPublicFavourites: newValue });
    app.setNotification(t('Updated public favourites setting'));
  } catch (error: unknown) {
    app.setNotification(t('Failed to update public favourites setting'), error);
  }
};

const updateProfile = async () => {
  validateAllFields();
  if (!isFormValid.value) {
    return;
  }

  try {
    await me.update(form);
    app.setNotification(t('Profile updated'));
  } catch (error: unknown) {
    app.setNotification(t('Failed to update profile'), error);
  }
};

const deleteHandler = async (event: { detail: { role: string } }): Promise<void> => {
  hasOpenDelete.value = false;
  const action = event?.detail?.role || 'cancel';
  if (action !== 'destructive') return;

  try {
    await me.deleteAccount();
  } catch (error: unknown) {
    app.setNotification(t('Failed to delete user'), error);
  }
};

onIonViewDidEnter(async () => {
  appearance.value = localStorage.getItem('appearance') || 'system';
});
</script>

<style scoped>
.content-list {
  max-width: 640px;
  margin: 0 auto;
}

ion-list,
ion-item {
  --background: transparent;
  background: transparent;
}

ion-list-header {
  font-weight: 500;
}
</style>
