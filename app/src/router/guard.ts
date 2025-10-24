import { NavigationGuardNext, RouteLocationNormalized } from 'vue-router';
import { useMeStore } from '@/stores/me.store';

export const routeGuard = async (
  to: RouteLocationNormalized,
  _: RouteLocationNormalized,
  next: NavigationGuardNext,
) => {
  const me = useMeStore();

  const publicPaths = ['/', '/login', '/verify-token'];
  const isPublicPath = publicPaths.includes(to.path);

  const loginPaths = ['/login', '/verify-token'];
  const isLoginPath = loginPaths.includes(to.path);

  const adminPaths = ['/users'];
  const isAdminPath = adminPaths.some((path) => to.path.startsWith(path));

  try {
    await me.init();
  } catch (error) {
    return next();
  }

  const isLoggedIn = me.isLoggedIn;
  const isAdmin = me.isAdmin;

  if (!isLoggedIn && !isPublicPath) {
    return next('/login');
  }

  if (isLoggedIn && isLoginPath) {
    return next('/');
  }

  if (isLoggedIn && !isAdmin && isAdminPath) {
    return next('/');
  }

  next();
};
