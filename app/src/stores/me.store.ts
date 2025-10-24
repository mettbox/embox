import { defineStore } from 'pinia';
import { httpService } from '@/services/http.service';
import router from '@/router';

export const useMeStore = defineStore('me', {
  state: () => ({
    isInitialized: false,
    isAuthenticated: false,
    user: null as User | null,
    error: null as unknown | null,
  }),

  getters: {
    isLoggedIn: (state) => state.isAuthenticated,
    isAdmin: (state) => (state.user ? state.user.isAdmin : false),

    // Fail Fast Getter (=Exception-Throwing Getter) Pattern
    id: (state) => {
      if (state.user) return state.user.id;
      throw new Error('user.id is not callable');
    },
    name: (state) => {
      if (state.user) return state.user.name;
      throw new Error('user.name is not callable');
    },
    email: (state) => {
      if (state.user) return state.user.email;
      throw new Error('user.email is not callable');
    },
    hasPublicFavourites: (state) => {
      if (state.user) return state.user.hasPublicFavourites;
      throw new Error('user.hasPublicFavourites is not callable');
    },
  },

  actions: {
    /**
     * Initialize the user store
     *
     * This method is called from the router guard to ensure the user store is initialized
     * Check if the user is authenticated and fetch the user data
     *
     * @returns {Promise<void>}
     */
    async init(): Promise<void> {
      if (this.isInitialized && this.isAuthenticated) return;

      const { data } = await httpService('auth/refresh', 'post');
      this.isInitialized = true;

      if (data && data.id) {
        this.user = data;
        this.isAuthenticated = true;
      }
    },

    /**
     * Request login token for the given email
     *
     * @param {string} email
     * @returns {Promise<void>}
     */
    async requestToken(email: string): Promise<void> {
      await httpService('auth/token', 'post', { email });
    },

    /**
     * Login with the auth token from the email
     *
     * @param {string} auth_token
     * @returns {Promise<void>}
     */
    async login(auth_token: string): Promise<void> {
      const { data } = await httpService('auth/login', 'post', { token: `${auth_token}` });
      if (data) {
        this.user = data;
        this.isAuthenticated = true;
      }
    },

    /**
     * Logout and redirect to the login page
     *
     * @returns {Promise<void>}
     */
    async logout(): Promise<void> {
      await httpService('auth/logout', 'post');
      this.user = null;
      this.isAuthenticated = false;
      router.replace({ name: 'login' });
    },

    /**
     * Update the current user with the given updates
     *
     * @param {Partial<User>} updates - Partial user object with fields to update
     * @returns {Promise<void>}
     */
    async update(updates: Partial<User>): Promise<void> {
      if (!this.user) throw new Error('Not authenticated');
      const needsReAuth = updates.email && updates.email !== this.user.email;
      const { data } = await httpService(`user/${this.user.id}`, 'put', updates);
      this.user = data;
      if (needsReAuth) await this.logout();
    },

    /**
     * Delete the current user's account
     *
     * @returns {Promise<void>}
     */
    async deleteAccount(): Promise<void> {
      if (!this.user) throw new Error('Not authenticated');
      await httpService(`user/${this.user.id}`, 'delete');
      this.user = null;
      this.isAuthenticated = false;
      router.replace({ name: 'login' });
    },
  },
});
