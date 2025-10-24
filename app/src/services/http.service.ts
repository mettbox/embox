import { useAppStore } from '@/stores/app.store';
import { useMeStore } from '@/stores/me.store';

enum RequestCredentials {
  INCLUDE = 'include',
  OMIT = 'omit',
  SAME_ORIGIN = 'same-origin',
}

type RequestMethod = 'get' | 'post' | 'put' | 'delete' | 'patch' | 'options' | 'head' | 'connect' | 'trace';

interface RequestOptions {
  headers?: Record<string, string>;
  credentials?: RequestCredentials;
  [key: string]: any;
}

interface HttpHeaders {
  [key: string]: string;
}

/**
 * HttpError
 *
 * @class HttpError
 * @extends {Error}
 */
export class HttpError extends Error {
  code: number;

  /**
   * Constructor
   *
   * @param {string} message
   * @param {number} code
   */
  constructor(message: string, code: number = 0) {
    super(message);
    this.code = code;
    this.name = 'HttpError';
  }
}

const baseURL = import.meta.env.VITE_API_URL || '';
const authEndpoint = 'auth/refresh';

let csrfToken: string | null = null;

/**
 * Get cookie by name
 *
 * @param {string} name
 * @returns {string | null}
 */
const getCookie = (name: string): string | null => {
  const value = `; ${document.cookie}`;
  const parts = value.split(`; ${name}=`);

  if (parts.length === 2) {
    return parts.pop()?.split(';').shift() || null;
  }

  return null;
};

/**
 * Fetch CSRF token from the server and store it in csrfToken variable
 *
 * @returns {Promise<void>}
 * @throws {HttpError}
 */
const fetchCsrfToken = async (): Promise<void> => {
  try {
    const url = `${baseURL}/csrf-token`;

    const response = await fetch(url, {
      method: 'get',
      credentials: 'include' as RequestCredentials,
    });

    if (!response.ok) {
      throw new Error('Failed to fetch CSRF token');
    }

    const cookie = getCookie('XSRF-TOKEN');
    if (cookie) {
      csrfToken = decodeURIComponent(cookie);
    } else {
      throw new Error('Failed to fetch CSRF token');
    }
  } catch (error) {
    handleHttpError(error, 'Failed to fetch CSRF token', 500);
  }
};

/**
 * Get headers for the request
 *
 * @param {string} csrfToken
 * @param {RequestOptions} options
 * @returns {HttpHeaders}
 */
const getHeaders = (csrfToken: string | null, options: RequestOptions = {}): HttpHeaders => {
  return {
    'Content-Type': 'application/json',
    accept: 'application/json',
    ...(csrfToken ? { 'X-XSRF-TOKEN': csrfToken } : {}),
    ...options.headers,
  };
};

export const handleHttpError = (error: unknown, defaultMessage: string, defaultCode: number = 500): never => {
  const app = useAppStore();

  let err: HttpError;

  if (error instanceof HttpError) {
    err = error;
  } else if (error instanceof Error) {
    err = new HttpError(error.message, defaultCode);
  } else {
    err = new HttpError(defaultMessage, defaultCode);
  }

  app.setNotification(err.message, err);
  throw err;
};

/**
 * Http request
 *
 * @param {string} endpoint - api endpoint w/o leading slash
 * @param {RequestMethod} [method=get]
 * @param {Record<string, unknown> | FormData} params - request parameters as object
 * @param {RequestOptions} options - additional header options
 * @returns {Promise<any>}
 * @throws {HttpError}
 */
export const httpService = async (
  endpoint: string,
  method: RequestMethod = 'get',
  params: Record<string, unknown> | FormData = {},
  options: RequestOptions = {},
  timeout: number = 60000, // 1 minutes in ms
): Promise<any> => {
  const app = useAppStore();
  const me = useMeStore();

  let url = `${baseURL}/${endpoint}`;

  if (!csrfToken) {
    try {
      app.setLoading(true);
      await fetchCsrfToken();
    } catch (error) {
      handleHttpError(error, 'Failed to fetch CSRF token', 500);
    } finally {
      app.setLoading(false);
    }
  }

  app.setLoading(true);

  const isFormData = typeof FormData !== 'undefined' && params instanceof FormData;

  const headers = isFormData
    ? { ...options.headers, ...(csrfToken ? { 'X-XSRF-TOKEN': csrfToken } : {}) }
    : getHeaders(csrfToken, options);

  if (method === 'get' && params) {
    const queryParams = new URLSearchParams(params as Record<string, string>);
    url += `?${queryParams.toString()}`;
  }

  let body: any = null;
  if (method !== 'get' && params) {
    body = isFormData ? params : JSON.stringify(params);
  }

  const controller = new AbortController();
  const signal = controller.signal;

  const timeoutId = setTimeout(() => {
    controller.abort();
  }, timeout);

  const config = {
    method: method.toString(),
    credentials: 'include' as RequestCredentials,
    headers,
    signal,
    ...(body ? { body } : {}),
    ...options,
  };

  try {
    const response = await fetch(url, config);

    // Handle blob responses
    if (options.responseType === 'blob') {
      if (!response.ok) {
        // TODO
        throw new HttpError(response.statusText, response.status);
      }
      return await response.blob();
    }

    // Handle successful responses w/o content
    // const results = response.status !== 204 ? await response.json() : null;
    let results = null;
    if (response.status !== 204) {
      const text = await response.text();
      try {
        results = text ? JSON.parse(text) : null;
      } catch {
        results = text;
      }
    }

    if (!response.ok) {
      if (response.status === 419) {
        console.warn('CSRF token mismatch, refreshing token and retrying request...');
        await fetchCsrfToken();
        headers['X-XSRF-TOKEN'] = csrfToken!;
        const retryResponse = await fetch(url, { ...config, headers });
        if (!retryResponse.ok) {
          throw new HttpError(retryResponse.statusText, retryResponse.status);
        }

        return results;
      }

      if (response.status === 401 && endpoint !== authEndpoint) {
        try {
          const refreshResponse = await fetch(`${baseURL}/${authEndpoint}`, {
            method: 'POST',
            credentials: 'include',
            headers: getHeaders(csrfToken),
          });

          if (refreshResponse.ok) {
            console.log('Session refreshed successfully, retrying request...');
            const retryResponse = await fetch(url, config);
            if (!retryResponse.ok) {
              throw new HttpError(retryResponse.statusText, retryResponse.status);
            }
            return options.responseType === 'blob' ? await retryResponse.blob() : await retryResponse.json();
          }
          throw new Error('Unauthorized');
        } catch (error: unknown) {
          console.error('Failed to refresh session:', error);
          const err = new HttpError('Unauthorized', 401);
          await me.logout();
          app.setNotification(err.message, err);
          throw err;
        }
      }

      // silent refresh
      if (response.status === 401 && endpoint === authEndpoint) {
        return null;
      }

      const err = new HttpError(results?.message || response.statusText, response.status);
      app.setNotification(err.message, err);

      throw err;
    }

    if (options.responseType === 'blob') {
      return response;
    }

    return results;
  } catch (error: unknown) {
    console.log('HTTP Error:', error);

    if (error instanceof Error) {
      if (error.name === 'AbortError') {
        const err = new HttpError('Request timed out', 408);
        app.setNotification(err.message, err);
        throw err;
      }

      const err = new HttpError(error.message, 500); // Default for all other errors
      app.setNotification(err.message, err);
      throw err;
    }

    // Netzwerkfehler explizit behandeln
    if (error instanceof TypeError && error.message === 'Failed to fetch') {
      const err = new HttpError('Network error: Failed to connect to the server', 503);
      app.setNotification(err.message, err);
      throw err;
    }

    // Fallback for non-Error objects
    app.setNotification('An unknown error occurred', 500);
    throw new HttpError('An unknown error occurred', 500);
  } finally {
    clearTimeout(timeoutId);
    app.setLoading(false);
  }
};
