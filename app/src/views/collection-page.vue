<template>
  <ion-page>
    <app-header has-back-button>
      <template v-if="collectionType === 'favourite'">{{ $t("Favourite's", { name: userName }) }}</template>
      <template v-else> Album Details </template>
    </app-header>

    <ion-content color="primary"> userName: {{ userName }} | id: {{ id }} </ion-content>
  </ion-page>
</template>

<script setup lang="ts">
import { IonPage, IonContent } from '@ionic/vue';
// import { onBeforeRouteUpdate } from 'vue-router';
import appHeader from '@/components/app/app-header.vue';
import { useRoute } from 'vue-router';
import { computed } from 'vue';

const route = useRoute();

const collectionType = computed(() => (route.name == 'favourites' ? 'favourite' : 'album'));

// getFavouritesByUserID()
// const userName = route.params.userName || null;
// const albumId = route.params.id || null;

defineProps<{
  userName?: string;
  id?: string;
}>();

// onBeforeRouteUpdate((to, from) => {
//   // Reagiere auf Änderungen der Route
//   console.log('Navigated to:', to.params);
//   console.log('Previous route:', from.params);
//   // Aktualisiere lokale Zustände oder führe andere Aktionen aus
//   // userName.value = to.params.userName || null;
//   // albumId.value = to.params.id || null;
// });
</script>
