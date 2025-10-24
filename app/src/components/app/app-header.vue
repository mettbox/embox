<template>
  <ion-header class="ion-no-border">
    <ion-toolbar color="primary">
      <ion-menu-button
        v-if="hasBackButton !== true"
        slot="start"
        color="light"
      />
      <ion-back-button
        v-if="hasBackButton === true"
        slot="start"
        :text="$t('Back')"
        color="light"
      />
      <ion-title>
        <slot />
      </ion-title>
      <ion-buttons slot="end">
        <slot name="buttons" />
      </ion-buttons>
      <app-loading />
    </ion-toolbar>
  </ion-header>
</template>

<script setup lang="ts">
import { IonHeader, IonToolbar, IonMenuButton, IonTitle, IonButtons, IonBackButton } from '@ionic/vue';
import appLoading from '@/components/app/app-loading.vue';

defineProps<{
  hasBackButton?: boolean;
}>();
</script>
