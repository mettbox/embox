import { createRouter, createWebHistory } from '@ionic/vue-router';
import { routeGuard } from './guard';

// static import for better performance for frequently used components
import LibraryPage from '@/views/library-page.vue';
import HomePage from '@/views/home-page.vue';

const routes = [
  {
    name: 'home',
    path: '/',
    component: HomePage,
  },
  {
    name: 'login',
    path: '/login/:email?',
    component: () => import('../views/login-request-page.vue'),
  },
  {
    name: 'token',
    path: '/verify-token',
    component: () => import('../views/login-token-page.vue'),
  },
  {
    name: 'library',
    path: '/library',
    component: LibraryPage,
  },
  {
    name: 'collections',
    path: '/collections',
    component: () => import('@/views/collections-page.vue'),
  },
  {
    path: '/favourites/:id',
    name: 'favourites',
    component: LibraryPage,
    props: true,
  },
  {
    path: '/albums/:id',
    name: 'albums',
    component: LibraryPage,
    props: true,
  },
  {
    name: 'users',
    path: '/users',
    component: () => import('@/views/users-page.vue'),
  },
  {
    name: 'profile',
    path: '/profile',
    component: () => import('@/views/profile-page.vue'),
  },
  {
    name: '404',
    path: '/:catchAll(.*)',
    redirect: '/',
  },
];

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes,
});

router.beforeEach(routeGuard);

export default router;
