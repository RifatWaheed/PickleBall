import axios from 'axios';

export const http = axios.create({
  baseURL: (import.meta as any).env.VITE_API_BASE_URL,
  withCredentials: true, // if backend uses cookies for refresh tokens
});

http.interceptors.request.use((config) => {
  const token = localStorage.getItem('access_token');
  if (token) config.headers.Authorization = `Bearer ${token}`;
  return config;
});

http.interceptors.response.use(
  (res) => res,
  (err) => {
    // standardize errors
    const status = err?.response?.status;
    if (status === 401) {
      // later: refresh token flow, or force logout
      // for now: just pass through
    }
    return Promise.reject(err);
  },
);
