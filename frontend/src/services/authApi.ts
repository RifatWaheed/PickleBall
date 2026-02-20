import { http } from './http';

export type Role = 'USER' | 'ADMIN' | 'SUPER_ADMIN';

export type MeResponse = {
  id: string;
  email: string;
  role: Role;
};

export async function login(req: { email: string; password: string }) {
  const { data } = await http.post('/auth/login', req);
  // expects { accessToken }
  return data as { accessToken: string };
}

export async function refresh() {
  const { data } = await http.post('/auth/refresh');
  return data as { accessToken: string };
}

export async function me() {
  const { data } = await http.get('/auth/me');
  return data as MeResponse;
}

export async function logout() {
  await http.post('/auth/logout');
}
