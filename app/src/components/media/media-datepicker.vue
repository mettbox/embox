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
      presentation="date-time"
    >
      <ion-buttons slot="buttons">
        <ion-button
          color="secondary"
          fill="solid"
          @click="cancel"
        >
          <ion-icon
            slot="icon-only"
            :icon="closeCircleOutline"
          />
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
          <ion-icon
            slot="icon-only"
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
  const value = datetimeElem.value?.$el.value ?? '';

  emit('update:modelValue', value);
  if (modalElem.value) {
    modalElem.value.$el.dismiss();
  }
};

const setAll = () => {
  if (datetimeElem.value) {
    datetimeElem.value.$el.confirm();
  }
  const value = datetimeElem.value?.$el.value ?? '';

  emit('update:all', value);
  if (modalElem.value) {
    modalElem.value.$el.dismiss();
  }
};
</script>
