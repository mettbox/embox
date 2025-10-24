<template>
  <ion-progress-bar
    v-if="app.isLoading"
    class="loading-bar"
    type="indeterminate"
    :color="app.isMdMode ? 'light' : 'primary'"
  />
</template>

<script setup lang="ts">
import { IonProgressBar } from '@ionic/vue';
import { useAppStore } from '@/stores/app.store';

const app = useAppStore();
</script>
