<template>
  <ion-header class="ion-no-border">
    <ion-toolbar color="primary">
      <ion-menu-button
        v-if="hasBackButton !== true"
        slot="start"
        color="light"
      />
      <ion-buttons
        slot="start"
        v-if="hasBackButton"
      >
        <ion-button
          @click="$router.back()"
          color="light"
        >
          <ion-icon
            :icon="chevronBackOutline"
            slot="icon-only"
          />
        </ion-button>
      </ion-buttons>
      <ion-text>
        <slot />
      </ion-text>
      <ion-buttons slot="end">
        <slot name="buttons" />
      </ion-buttons>
      <app-loading />
    </ion-toolbar>
  </ion-header>
</template>

<script setup lang="ts">
import { IonHeader, IonToolbar, IonMenuButton, IonText, IonButtons, IonButton, IonIcon } from '@ionic/vue';
import { chevronBackOutline } from 'ionicons/icons';

import appLoading from '@/components/app/app-loading.vue';

defineProps<{
  hasBackButton?: boolean;
}>();
</script>
