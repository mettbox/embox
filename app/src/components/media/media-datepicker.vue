<template>
  <ion-modal
    ref="modalElem"
    :keep-contents-mounted="true"
  >
    <ion-datetime
      ref="datetimeElem"
      :id="datetimeId"
      :value="props.modelValue"
      :first-day-of-week="1"
      :prefer-wheel="true"
      presentation="date"
    >
      <ion-buttons slot="buttons">
        <ion-button
          color="secondary"
          fill="solid"
          @click="cancel"
        >
          <ion-icon
            slot="start"
            :icon="closeCircleOutline"
          />
          {{ $t('Cancel') }}
        </ion-button>

        <ion-button
          color="secondary"
          fill="solid"
          @click="setAll()"
        >
          <ion-icon
            slot="start"
            :icon="checkmarkDoneCircleOutline"
          />
          {{ $t('Apply to all w/o date') }}
        </ion-button>

        <ion-button
          color="secondary"
          fill="solid"
          :strong="true"
          @click="confirm()"
        >
          {{ $t('Apply') }}
          <ion-icon
            slot="start"
            :icon="checkmarkCircleOutline"
          />
        </ion-button>
      </ion-buttons>
    </ion-datetime>
  </ion-modal>
</template>

<script lang="ts" setup>
import { IonDatetime, IonModal, IonButtons, IonButton, IonIcon } from '@ionic/vue';
import { checkmarkCircleOutline, checkmarkDoneCircleOutline, closeCircleOutline } from 'ionicons/icons';

import { ref, computed } from 'vue';

const props = defineProps<{
  modelValue: string | undefined;
  id?: string;
}>();

const emit = defineEmits<{
  (e: 'update:modelValue', value: string): void;
  (e: 'update:all', value: string): void;
}>();

const modalElem = ref<InstanceType<typeof IonModal> | null>(null);
const datetimeElem = ref<InstanceType<typeof IonDatetime> | null>(null);

const datetimeId = computed(() => props.id ?? 'datetime-' + Math.random().toString(36).slice(2, 8));

const cancel = () => {
  if (modalElem.value) {
    modalElem.value.$el.dismiss();
  }
};

const confirm = () => {
  if (datetimeElem.value) {
    datetimeElem.value.$el.confirm();
  }
  const value = datetimeElem.value?.$el.value?.slice(0, 10) ?? '';

  emit('update:modelValue', value);
  if (modalElem.value) {
    modalElem.value.$el.dismiss();
  }
};

const setAll = () => {
  if (datetimeElem.value) {
    datetimeElem.value.$el.confirm();
  }
  const value = datetimeElem.value?.$el.value?.slice(0, 10) ?? '';

  emit('update:all', value);
  if (modalElem.value) {
    modalElem.value.$el.dismiss();
  }
};
</script>

<style scoped>
ion-datetime {
  --background: var(--ion-color-light);
  --background-rgb: var(--ion-color-light-rgb);
  box-shadow: rgba(var(--ion-color-light-rgb), 0.3) 0px 10px 15px -3px;
}
</style>
