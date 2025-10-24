import { defineStore } from 'pinia';
import { Mode } from '@ionic/core';

export const useAppStore = defineStore('app', {
  state: () => ({
    loading: false,
    notification: '',
    error: null as Error | unknown | null,
    screenWidth: window.innerWidth,
    breakpoint: 'xs' as ViewSize,
    appMode: 'md' as Mode,
    isDarkMode: window.matchMedia('(prefers-color-scheme: dark)').matches,
  }),

  actions: {
    getBreakpoint(width: number): ViewSize {
      if (width >= 1200) return 'xl';
      if (width >= 992) return 'lg';
      if (width >= 768) return 'md';
      if (width >= 576) return 'sm';
      return 'xs';
    },
    updateBreakpoint() {
      this.screenWidth = window.innerWidth;
      this.breakpoint = this.getBreakpoint(this.screenWidth);
    },

    /**
     * Set the app mode (md or ios)
     *
     * @param {Mode} mode - 'md' or 'ios'
     */
    setAppMode(mode: Mode) {
      this.appMode = mode;
    },

    initDarkMode() {
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)');
      this.isDarkMode = prefersDark.matches;
      document.documentElement.classList.toggle('ion-palette-dark', this.isDarkMode);
    },

    setAppearance(mode: 'dark' | 'light' | 'system') {
      localStorage.setItem('appearance', mode);
      if (mode === 'system') {
        this.initDarkMode();
      } else {
        this.isDarkMode = mode === 'dark';
        document.documentElement.classList.toggle('ion-palette-dark', this.isDarkMode);
      }
    },

    loadAppearance() {
      const mode = localStorage.getItem('appearance') as 'dark' | 'light' | 'system' | null;
      if (mode) {
        this.setAppearance(mode);
      } else {
        this.initDarkMode();
      }
    },

    /**
     * Set the loading state. Used by http service to indicate loading state
     *
     * @param {boolean} value - true to set loading, false to unset
     */
    setLoading(value: boolean): void {
      this.loading = value;
    },

    /**
     * Global notification handler.
     * Sets the notification message and error state.
     *
     * @param {string} notification
     * @param {Error|unknown|null} [error=null] - Optional error object to set
     */
    setNotification(notification: string, error: unknown | Error | null = null): void {
      this.notification = notification;
      this.error = error;
    },

    /**
     * Clear the current notification and error state
     */
    unsetNotification(): void {
      this.notification = '';
      this.error = null;
    },
  },

  getters: {
    isIosMode: (state) => state.appMode === 'ios',
    isMdMode: (state) => state.appMode === 'md',
    isLoading: (state) => state.loading,
    hasNotification: (state) => !!state.notification || !!state.error,
    hasError: (state) => state.error !== null,
    getNotification: (state): AppNotification => ({
      message: state.notification,
      error: state.error,
    }),
    isXs: (state) => state.breakpoint === 'xs',
    isSm: (state) => ['xs', 'sm'].includes(state.breakpoint),
    isMd: (state) => ['xs', 'sm', 'md'].includes(state.breakpoint),
    isLg: (state) => ['xs', 'sm', 'md', 'lg'].includes(state.breakpoint),
    isXl: () => true, // always true, xl is the largest breakpoint
    mediaGridColumns: (state) => {
      switch (state.breakpoint) {
        case 'xs':
          return 4;
        case 'sm':
          return 5;
        case 'md':
          return 6;
        case 'lg':
          return 7;
        default:
          return 8;
      }
    },
  },
});
