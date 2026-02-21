/// <reference types="vite/client" />
import axios from 'axios';

export const http = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL,
  withCredentials: true,
});

let accessToken: string | null = null;
export function setAccessTokenMemory(token: string | null) {
  accessToken = token;
}

http.interceptors.request.use((config) => {
  if (accessToken) config.headers.Authorization = `Bearer ${accessToken}`;
  return config;
});

let isRefreshing = false;
let pending: Array<(token: string | null) => void> = [];

function notify(token: string | null) {
  pending.forEach((cb) => cb(token));
  pending = [];
}

http.interceptors.response.use(
  (res) => res,
  async (err) => {
    const original = err.config;

    // Prevent infinite loops
    if (err.response?.status !== 401 || original._retry) {
      return Promise.reject(err);
    }
    original._retry = true;

    // Queue requests while refresh happens
    if (isRefreshing) {
      return new Promise((resolve, reject) => {
        pending.push((token) => {
          if (!token) return reject(err);
          original.headers.Authorization = `Bearer ${token}`;
          resolve(http(original));
        });
      });
    }

    isRefreshing = true;
    try {
      const { data } = await http.post('/auth/refresh');
      setAccessTokenMemory(data.accessToken);
      notify(data.accessToken);
      original.headers.Authorization = `Bearer ${data.accessToken}`;
      return http(original);
    } catch (e) {
      setAccessTokenMemory(null);
      notify(null);
      return Promise.reject(e);
    } finally {
      isRefreshing = false;
    }
  },
);
