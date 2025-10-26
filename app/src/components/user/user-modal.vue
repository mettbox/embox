<template>
  <ion-modal
    ref="modalRef"
    :is-open="props.isOpen"
  >
    <ion-header>
      <ion-toolbar color="primary">
        <ion-title>{{ title }}</ion-title>
        <ion-buttons slot="end">
          <ion-button @click="onCancel">
            <ion-icon
              slot="icon-only"
              :icon="close"
            />
          </ion-button>
        </ion-buttons>
      </ion-toolbar>
    </ion-header>

    <ion-content>
      <ion-list>
        <ion-item>
          <ion-input
            ref="nameInput"
            v-model="form.name"
            :placeholder="$t('Name')"
            :error-text="isTouched.name ? fieldErrors.name : ''"
            @ion-blur="() => validateField('name')"
          />
        </ion-item>
        <ion-item>
          <ion-input
            ref="emailInput"
            v-model="form.email"
            :disabled="props.isCurrentUser"
            :placeholder="$t('Email address')"
            type="email"
            :error-text="isTouched.email ? fieldErrors.email : ''"
            :helper-text="props.isCurrentUser ? $t('You cannot change your own email address.') : undefined"
            @ion-blur="() => validateField('email')"
          />
        </ion-item>
        <ion-item v-if="!props.isCurrentUser">
          <ion-toggle
            v-model="form.isAdmin"
            justify="space-between"
            :disabled="props.isCurrentUser"
          >
            {{ $t('Administrator?') }}
          </ion-toggle>
        </ion-item>
      </ion-list>

      <ion-list v-if="!props.isCurrentUser">
        <ion-item>
          <ion-toggle
            v-model="sendMail"
            justify="space-between"
            :disabled="props.isCurrentUser"
          >
            {{ $t('Send email?') }}
          </ion-toggle>
        </ion-item>
        <ion-item
          v-if="sendMail"
          lines="none"
        >
          <ion-input
            ref="subjectInput"
            v-model="form.subject"
            :placeholder="$t('Subject')"
            :error-text="isTouched.subject ? fieldErrors.subject : ''"
            @ion-blur="() => validateField('subject')"
          />
        </ion-item>
        <ion-item
          v-if="sendMail"
          lines="none"
        >
          <ion-textarea
            ref="messageInput"
            v-model="form.message"
            fill="solid"
            :auto-grow="true"
            :placeholder="$t('Message')"
            :error-text="isTouched.message ? fieldErrors.message : ''"
            @ion-blur="() => validateField('message')"
          />
        </ion-item>
      </ion-list>
    </ion-content>

    <ion-footer>
      <ion-toolbar
        color="primary"
        class="footer"
      >
        <ion-buttons slot="end">
          <ion-button @click="onSave">
            <ion-icon
              slot="icon-only"
              :icon="saveOutline"
            />
          </ion-button>
        </ion-buttons>
      </ion-toolbar>
    </ion-footer>
  </ion-modal>
</template>

<script setup lang="ts">
import { computed, reactive, ref, watch } from 'vue';
import { useI18n } from 'vue-i18n';
import {
  IonButton,
  IonButtons,
  IonContent,
  IonHeader,
  IonInput,
  IonItem,
  IonList,
  IonModal,
  IonTextarea,
  IonTitle,
  IonToggle,
  IonToolbar,
  IonFooter,
} from '@ionic/vue';
import { IonIcon } from '@ionic/vue';
import { close, saveOutline } from 'ionicons/icons';
import { useValidation, type FormField } from '@/composables/use-validation';

const props = defineProps<{
  isOpen: boolean;
  id?: string;
  name?: string;
  email?: string;
  isAdmin: boolean;
  existingEmails: string[];
  isCurrentUser: boolean;
}>();

const loginUrl = computed(() => {
  const url = new URL(window.location.href);
  url.pathname = '/login';
  return url.toString();
});

const helloMessage = computed(() => {
  return t('Hello message', { name: props.name ? props.name : '' }) + loginUrl.value;
});

const { t } = useI18n();
const modalRef = ref<InstanceType<typeof IonModal> | null>(null);
const title = ref('');
const nameInput = ref<InstanceType<typeof IonInput> | null>(null);
const emailInput = ref<InstanceType<typeof IonInput> | null>(null);
const sendMail = ref(false);
const subjectInput = ref<InstanceType<typeof IonInput> | null>(null);
const messageInput = ref<InstanceType<typeof IonTextarea> | null>(null);

const getDefaults = (): Record<string, any> => ({
  id: props.id || null,
  name: props.name || '',
  email: props.email || '',
  isAdmin: props.isAdmin || false,
  subject: props.isCurrentUser === false ? t('Welcome!') : '',
  message: helloMessage.value,
});

const form = reactive(getDefaults());

const fields: Record<string, FormField> = {
  name: {
    ref: nameInput,
    rules: [
      {
        validate: (val) => val.trim() !== '',
        message: t('Please enter a name'),
      },
    ],
  },
  email: {
    ref: emailInput,
    shouldValidate: () => props.isCurrentUser === false,
    rules: [
      {
        validate: (val) => val.trim() !== '',
        message: t('Please enter an email address'),
      },
      {
        validate: (val) => /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(String(val ?? '')),
        message: t('Please enter a valid email address'),
      },
    ],
  },
  subject: {
    ref: subjectInput,
    shouldValidate: () => sendMail.value,
    rules: [
      {
        validate: (val) => val.trim() !== '',
        message: t('Subject is required'),
      },
    ],
  },
  message: {
    ref: messageInput,
    shouldValidate: () => sendMail.value,
    rules: [
      {
        validate: (val) => val.trim() !== '',
        message: t('Message is required'),
      },
    ],
  },
};

const { validateField, validateAllFields, isFormValid, fieldErrors, isTouched } = useValidation(fields, form);

const resetForm = () => {
  title.value = props.id ? props.name || '' : t('New User');
  sendMail.value = !props.id;
  const defaults = getDefaults();
  Object.keys(defaults).forEach((key) => {
    form[key] = defaults[key];
  });
};

const onCancel = () => {
  if (modalRef.value) {
    modalRef.value.$el.dismiss(null, 'cancel');
  }
};

const onSave = () => {
  validateAllFields();
  if (!isFormValid.value) {
    return;
  }

  const currentUrl = new URL(window.location.href);
  currentUrl.pathname = `/login`;
  const message = sendMail.value
    ? form.message
        .replace(loginUrl.value, `<a href="${currentUrl.toString()}" target="_blank">${loginUrl.value}</a>`)
        .replace(/\n/g, '<br>')
    : '';

  const result = {
    ...form,
    subject: sendMail.value ? form.subject : '',
    message,
  };

  modalRef.value?.$el.dismiss(result, 'save');
};

watch(
  () => props.isOpen,
  (open) => {
    if (open) {
      resetForm();
    }
  },
);
</script>
