<template>
  <ion-menu
    content-id="main-content"
    type="overlay"
  >
    <ion-header class="ion-no-border">
      <ion-toolbar>
        <ion-title
          :color="app.isDarkMode ? 'light' : 'primary'"
          size="large"
        >
          <app-logo />
        </ion-title>
      </ion-toolbar>
    </ion-header>

    <ion-content>
      <ion-list>
        <ion-list-header
          v-if="isAuth"
          class="ion-padding-bottom"
        >
          <ion-label>{{ $t('Hello', { name: userName }) }}</ion-label>
          <ion-menu-toggle :auto-hide="false">
            <ion-button
              class="ion-margin-end"
              shape="round"
              router-direction="forward"
              :router-link="{ name: 'profile' }"
            >
              <ion-icon
                slot="icon-only"
                :icon="ellipsisHorizontalCircleOutline"
              />
            </ion-button>
          </ion-menu-toggle>
        </ion-list-header>
        <ion-menu-toggle :auto-hide="false">
          <ion-item
            router-direction="root"
            :router-link="{ name: 'home' }"
            lines="none"
            :detail="false"
          >
            <ion-icon
              aria-hidden="true"
              slot="start"
              :icon="route.name === 'home' ? time : timeOutline"
            />
            <ion-label>{{ $t('Home') }}</ion-label>
          </ion-item>

          <ion-item
            v-if="isAuth"
            router-direction="root"
            :router-link="{ name: 'library' }"
            lines="none"
            :detail="false"
          >
            <ion-icon
              aria-hidden="true"
              slot="start"
              :icon="route.name === 'library' ? images : imagesOutline"
            />
            <ion-label>{{ $t('Library') }}</ion-label>
          </ion-item>

          <ion-item
            v-if="isAuth"
            router-direction="root"
            :router-link="{ name: 'collections' }"
            lines="none"
            :detail="false"
          >
            <ion-icon
              aria-hidden="true"
              slot="start"
              :icon="route.name === 'collections' ? albums : albumsOutline"
            />
            <ion-label>{{ $t('Collections') }}</ion-label>
          </ion-item>

          <ion-item
            v-if="isAdmin"
            router-direction="root"
            :router-link="{ name: 'users' }"
            lines="none"
            :detail="false"
          >
            <ion-icon
              aria-hidden="true"
              slot="start"
              :icon="route.name === 'users' ? people : peopleOutline"
            />
            <ion-label>{{ $t('Users') }}</ion-label>
          </ion-item>

          <ion-item
            v-if="!isAuth"
            router-direction="root"
            :router-link="{ name: 'login', params: {} }"
            lines="none"
            :detail="false"
          >
            <ion-icon
              aria-hidden="true"
              slot="start"
              :icon="route.name === 'login' || 'token' ? logIn : logInOutline"
            />
            <ion-label>{{ $t('Login') }}</ion-label>
          </ion-item>

          <ion-item
            v-if="isAuth"
            router-direction="root"
            :router-link="{ name: 'home' }"
            @click="me.logout"
            lines="none"
            :detail="false"
          >
            <ion-icon
              aria-hidden="true"
              slot="start"
              :icon="logOutOutline"
            />
            <ion-label>{{ $t('Logout') }}</ion-label>
          </ion-item>
        </ion-menu-toggle>
      </ion-list>
    </ion-content>

    <ion-footer
      id="footer"
      class="ion-no-border"
    >
      <ion-list class="ion-padding-bottom">
        <ion-item
          target="_blank"
          lines="none"
          :detail="false"
          @click="reloadApp"
          :button="true"
        >
          <ion-icon
            inert
            slot="start"
            :icon="reloadOutline"
          />
          <ion-label>{{ $t('Reload') }}</ion-label>
        </ion-item>

        <ion-item
          href="https://github.com/mettbox/embox"
          target="_blank"
          lines="none"
          :detail="false"
        >
          <ion-icon
            inert
            slot="start"
            :icon="logoGithub"
          />
          <ion-label>Maik Mettenheimer</ion-label>
        </ion-item>
      </ion-list>
    </ion-footer>
  </ion-menu>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import {
  IonButton,
  IonContent,
  IonFooter,
  IonHeader,
  IonIcon,
  IonItem,
  IonLabel,
  IonList,
  IonMenu,
  IonMenuToggle,
  IonTitle,
  IonToolbar,
  IonListHeader,
} from '@ionic/vue';
import {
  logoGithub,
  time,
  timeOutline,
  people,
  peopleOutline,
  images,
  imagesOutline,
  albums,
  albumsOutline,
  logIn,
  logInOutline,
  logOutOutline,
  ellipsisHorizontalCircleOutline,
  reloadOutline,
} from 'ionicons/icons';
import appLogo from '@/components/app/app-logo.vue';
import { useAppStore } from '@/stores/app.store';
import { useMeStore } from '@/stores/me.store';

const route = useRoute();
const app = useAppStore();
const me = useMeStore();

const isAuth = computed(() => me.isLoggedIn);
const isAdmin = computed(() => me.isAdmin);
const userName = computed(() => me.name || '');

const reloadApp = () => {
  location.reload();
};
</script>

<style scoped>
ion-menu {
  ion-toolbar,
  ion-list,
  ion-item,
  ion-footer,
  ion-content {
    --background: transparent;
    background: transparent;
  }

  ion-header {
    ion-title {
      font-size: 2rem;
      font-family: 'story';
      text-transform: uppercase;
    }

    ion-buttons {
      padding: 8px;
    }
  }

  ion-item {
    ion-icon {
      font-size: 24px;
      color: var(--ion-color-secondary);
    }
  }

  ion-list-header {
    font-weight: 500;
  }
}
</style>
