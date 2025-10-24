<template>
  <ion-list class="content-list">
    <ion-item>
      <ion-input
        ref="nameInput"
        type="text"
        v-model="form.name"
        label-placement="stacked"
        :label="$t('Name')"
        :placeholder="$t('Enter an Album name…')"
        :error-text="isTouched.name ? fieldErrors.name : ''"
        @ion-blur="() => validateField('name')"
        @keyup.enter="update"
      />
    </ion-item>
    <ion-item>
      <ion-textarea
        ref="descriptionInput"
        v-model="form.description"
        :label="$t('Description')"
        label-placement="stacked"
        :placeholder="$t('Enter an album description…')"
        :error-text="isTouched.description ? fieldErrors.description : ''"
        @ion-blur="() => validateField('description')"
      />
    </ion-item>
    <ion-item>
      <ion-toggle
        justify="space-between"
        v-model="form.isPublic"
      >
        {{ $t('Is public') }}
      </ion-toggle>
    </ion-item>
    <ion-item class="ion-margin-top">
      <ion-button
        size="default"
        slot="end"
        @click="update"
        >{{ $t('Done') }}</ion-button
      >
    </ion-item>
  </ion-list>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useI18n } from 'vue-i18n';
import { IonInput, IonList, IonToggle, IonItem, IonButton, IonTextarea } from '@ionic/vue';
import { useValidation, type FormField } from '@/composables/use-validation';

const props = defineProps<{
  album?: AlbumForm;
}>();

const emit = defineEmits<{
  (e: 'done', album: AlbumForm): void;
}>();

const form = reactive<AlbumForm>({
  id: props.album?.id,
  name: props.album?.name || '',
  description: props.album?.description || '',
  isPublic: props.album?.isPublic || false,
});

const { t } = useI18n();

const nameInput = ref<InstanceType<typeof IonInput> | null>(null);
const descriptionInput = ref<InstanceType<typeof IonInput> | null>(null);

const fields: Record<string, FormField> = {
  name: {
    ref: nameInput,
    rules: [
      {
        validate: (val) => val.trim() !== '',
        message: t('Please enter an album name'),
      },
    ],
  },
  description: {
    ref: descriptionInput,
    rules: [],
  },
};

const { validateField, validateAllFields, isFormValid, fieldErrors, isTouched } = useValidation(fields, form);

const update = async () => {
  validateAllFields();
  if (!isFormValid.value) {
    return;
  }
  emit('done', {
    id: form.id,
    name: form.name.trim(),
    description: form.description.trim(),
    isPublic: form.isPublic,
  });
};
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
