import { createApp } from 'vue';
import App from './app.vue';
import router from './router';
import { createPinia } from 'pinia';
import { i18n } from './i18n/i18n';
import { IonicVue } from '@ionic/vue';

/* Core CSS required for Ionic components to work properly */
import '@ionic/vue/css/core.css';

/* Basic CSS for apps built with Ionic */
import '@ionic/vue/css/normalize.css';
import '@ionic/vue/css/structure.css';
import '@ionic/vue/css/typography.css';

/* Optional CSS utils that can be commented out */
import '@ionic/vue/css/padding.css';
import '@ionic/vue/css/float-elements.css';
import '@ionic/vue/css/text-alignment.css';
import '@ionic/vue/css/text-transformation.css';
import '@ionic/vue/css/flex-utils.css';
import '@ionic/vue/css/display.css';

/**
 * Ionic Dark Mode
 * -----------------------------------------------------
 * For more info, please see:
 * https://ionicframework.com/docs/theming/dark-mode
 */

// import '@ionic/vue/css/palettes/dark.always.css';
// import '@ionic/vue/css/palettes/dark.system.css';
import '@ionic/vue/css/palettes/dark.class.css';

/* Theme variables */
import './theme/variables.css';

/* App custom styles */
import './theme/app.css';
import { useAppStore } from './stores/app.store';

const app = createApp(App)
  .use(createPinia())
  .use(IonicVue, {
    // For more info, please see: https://ionicframework.com/docs/developing/config
    // mode: 'ios',
    // rippleEffect: false,
    // swipeBackEnabled: false,
    // IonToast content is parsed as plaintext by default. innerHTMLTemplatesEnabled enables custom HTML
    innerHTMLTemplatesEnabled: true,
  })
  .use(router)
  .use(i18n);

router.isReady().then(() => {
  useAppStore().initDarkMode();
  app.mount('#app');
});
