<template>
  <ion-page class="home-page">
    <app-header />

    <ion-content color="primary">
      <v-center-container fullwidth>
        <v-fit-to-width
          class="text-fit-header"
          text="<b>E</b>.<small>M</small>.<i>♥︎</i>"
        />
        <v-fit-to-width
          v-for="(time, index) in times"
          :text="time"
          :key="`fit-${index}`"
        />
      </v-center-container>
    </ion-content>
  </ion-page>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { useI18n } from 'vue-i18n';
import { IonContent, IonPage } from '@ionic/vue';
import vFitToWidth from '@/components/v-fit-to-width.vue';
import vCenterContainer from '@/components/app/app-center-container.vue';
import appHeader from '@/components/app/app-header.vue';

const { t } = useI18n();
const startDate = new Date('2024-05-27T12:23:00');
const times = ref<string[]>([]);

const updateTimes = (): void => {
  const now = new Date();
  const diffInMs = now.getTime() - startDate.getTime();
  const years = diffInMs / 1000 / 60 / 60 / 24 / 365;

  times.value = [
    `${(years < 1 ? years.toFixed(3) : Math.floor(years)).toLocaleString()} <i>${t('Years')}</i>`,
    `${Math.floor(diffInMs / 1000 / 60 / 60 / 24 / 30).toLocaleString()} ${t('Months')}`,
    `${Math.floor(diffInMs / 1000 / 60 / 60 / 24 / 7).toLocaleString()} ${t('Weeks')}`,
    `${Math.floor(diffInMs / 1000 / 60 / 60 / 24).toLocaleString()} ${t('Days')}`,
    `${Math.floor(diffInMs / 1000 / 60 / 60).toLocaleString()} ${t('Hours')}`,
    `${Math.floor(diffInMs / 1000 / 60).toLocaleString()} ${t('Minutes')}`,
    `${Math.floor(diffInMs / 1000).toLocaleString()} ${t('Seconds')}`,
  ];
};

updateTimes();
setInterval(updateTimes, 1000);
</script>

<style>
.home-page {
  ion-content {
    font-family: 'story';

    .text-fit:nth-child(even) {
      background-color: #a73356;
    }

    .text-fit:last-child {
      margin-bottom: var(--ion-padding-double);
    }

    .text-fit-header {
      i {
        font-size: 0.65em;
      }
    }
  }
}
</style>
