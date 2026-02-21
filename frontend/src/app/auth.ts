export type AuthUser = {
  id: string;
  email: string;
  role: 'USER' | 'ADMIN' | 'SUPER_ADMIN';
};

export function setAccessToken(token: string) {
  localStorage.setItem('access_token', token);
}

export function clearAccessToken() {
  localStorage.removeItem('access_token');
}
