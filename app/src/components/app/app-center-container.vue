<template>
  <ion-grid :class="{ fullwidth: props.fullwidth }">
    <ion-row class="ion-align-items-center">
      <ion-col>
        <div :class="{ 'ion-text-center': !props.fullwidth }">
          <slot />
        </div>
      </ion-col>
    </ion-row>
  </ion-grid>
</template>

<script setup lang="ts">
import { IonCol, IonGrid, IonRow } from '@ionic/vue';

const props = withDefaults(defineProps<{ fullwidth?: boolean }>(), {
  fullwidth: false,
});
</script>

<style scoped>
ion-grid {
  height: 100%;
  max-width: 640px;
  display: flex;
  flex-direction: column;
}

ion-grid.fullwidth {
  --ion-grid-padding: 0;
  --ion-grid-column-padding: 0;
  max-width: 100%;
}

ion-row {
  flex: 1;
}
</style>
