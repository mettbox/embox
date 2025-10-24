import { createI18n } from 'vue-i18n';
import de from './de';

const userLocale = navigator.language.split('-')[0];
const availableLocales = ['de'];
const locale = availableLocales.includes(userLocale) ? userLocale : 'de';

export const i18n = createI18n({
  legacy: false, // important for composition API
  locale: locale,
  fallbackLocale: 'de',
  messages: {
    de,
  },
});
