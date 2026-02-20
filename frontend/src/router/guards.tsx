import { Navigate, Outlet } from 'react-router-dom';
import { useAuth } from '@/app/AuthProvider';
import type { Role } from '@/services/authApi';

export function RequireAuth() {
  const { isAuthed, isBootstrapping } = useAuth();
  if (isBootstrapping) return <div className="p-6">Loading...</div>;
  if (!isAuthed) return <Navigate to="/login" replace />;
  return <Outlet />;
}

export function RequireRole({ allowed }: { allowed: Role[] }) {
  const { user, isBootstrapping } = useAuth();
  if (isBootstrapping) return <div className="p-6">Loading...</div>;
  if (!user) return <Navigate to="/login" replace />;
  if (!allowed.includes(user.role)) return <Navigate to="/403" replace />;
  return <Outlet />;
}
