<template>
  <ion-page>
    <app-header>
      {{ $t('Users') }}
      <template #buttons>
        <ion-button @click="onAdd">
          <ion-icon
            :icon="addCircleOutline"
            slot="icon-only"
          />
        </ion-button>
      </template>
    </app-header>

    <ion-content :fullscreen="true">
      <ion-list>
        <user-list-item
          v-for="user in users"
          :key="user.id"
          :user="user"
          :can-delete="isDeleteable(user)"
          @edit="onEdit"
          @remove="onRemove"
        />
      </ion-list>

      <ion-action-sheet
        :is-open="hasOpenDelete"
        :header="$t('This action is permanent and cannot be undone.')"
        :buttons="deleteActionSheetButtons"
        @did-dismiss="deleteHandler"
      />

      <user-modal
        :is-open="hasOpenModal"
        :id="user?.id"
        :name="user?.name"
        :email="user?.email"
        :is-admin="user?.isAdmin ?? false"
        :is-current-user="user?.id === me.id"
        :existing-emails="users.map((u: User) => u.email)"
        @did-dismiss="editHandler"
      />
    </ion-content>
  </ion-page>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useI18n } from 'vue-i18n';
import { IonActionSheet, IonContent, IonIcon, IonList, IonButton, IonPage, onIonViewDidEnter } from '@ionic/vue';
import { ActionSheetButton } from '@ionic/core';
import { addCircleOutline } from 'ionicons/icons';
import UserListItem from '@/components/user/user-list-item.vue';
import UserModal from '@/components/user/user-modal.vue';
import { useAppStore } from '@/stores/app.store';
import { useMeStore } from '@/stores/me.store';
import { UserService } from '@/services/user.service';
import appHeader from '@/components/app/app-header.vue';

const lastAdminId = computed(() => {
  const admins = users.value.filter((a) => a.isAdmin);

  return admins.length === 1 ? admins[0].id : null;
});

const deleteActionSheetButtons = computed((): ActionSheetButton[] => [
  {
    text: t('Remove', { name: user.value ? user.value.name : '' }),
    role: 'destructive',
  },
  {
    text: t('Cancel'),
    role: 'cancel',
  },
]);

const { t } = useI18n();
const app = useAppStore();
const me = useMeStore();

const hasOpenDelete = ref(false);
const hasOpenModal = ref(false);

const users = ref([] as User[]);
const user = ref<User | null>(null);

const isDeleteable = (u: User): boolean => {
  return u.id !== lastAdminId.value && u.id !== me.id;
};

const onAdd = (): void => {
  user.value = null;
  hasOpenModal.value = true;
};

const onEdit = (u: User): void => {
  user.value = u;
  hasOpenModal.value = true;
};

const onRemove = (u: User): void => {
  user.value = u;
  hasOpenDelete.value = true;
};

const deleteHandler = async (event: { detail: { role: string } }): Promise<void> => {
  hasOpenDelete.value = false;

  const action = event?.detail?.role || 'cancel';

  if (user.value === null || action !== 'destructive') {
    return;
  }

  try {
    await UserService.delete(user.value.id);
  } catch (error: unknown) {
    app.setNotification(t('Failed to delete user'), error);
  } finally {
    await load();
  }
};

const editHandler = async (event: CustomEvent): Promise<void> => {
  hasOpenModal.value = false;

  if (!event.detail || event.detail.role !== 'save') {
    return;
  }

  try {
    const u = event.detail.data as User;
    await UserService.save(u);
  } catch (error: unknown) {
    app.setNotification(t('Failed to edit user'), error);
  } finally {
    await load();
  }
};

const load = async (): Promise<void> => {
  try {
    users.value = await UserService.list();
    user.value = null;
  } catch (error: unknown) {
    app.setNotification(t('Failed to load users'), error);
  }
};

onIonViewDidEnter(async () => {
  await load();
});
</script>
