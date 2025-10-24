<template>
  <ion-item-sliding ref="item">
    <ion-item>
      <ion-label>
        <h3>{{ props.user.name }}</h3>
        <p>{{ props.user.email }}</p>
      </ion-label>
    </ion-item>
    <ion-item-options>
      <ion-item-option
        color="medium"
        @click="onAction('edit')"
        ><ion-icon
          slot="top"
          :icon="createOutline"
        ></ion-icon
        >{{ $t('Edit') }}</ion-item-option
      >
      <ion-item-option
        v-if="canDelete"
        color="primary"
        @click="onAction('remove')"
        ><ion-icon
          slot="top"
          :icon="trashOutline"
        ></ion-icon
        >{{ $t('Delete') }}</ion-item-option
      >
    </ion-item-options>
  </ion-item-sliding>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import { IonIcon, IonItem, IonItemOption, IonItemOptions, IonItemSliding, IonLabel } from '@ionic/vue';
import { trashOutline, createOutline } from 'ionicons/icons';

const props = defineProps<{
  user: User;
  canDelete: boolean;
}>();

const emit = defineEmits<{
  (e: 'edit', user: User): void;
  (e: 'remove', user: User): void;
}>();

const item = ref<VueInstanceElement>();

const onAction = (action: string): void => {
  if (item.value) {
    item.value.$el.closeOpened();
  }

  if (action === 'edit') {
    emit('edit', props.user);
  } else {
    emit('remove', props.user);
  }
};
</script>
