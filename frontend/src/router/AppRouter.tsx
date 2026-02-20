import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import RootLayout from '@/layouts/RootLayout';
import HomePage from '@/pages/home/HomePage';
import LoginPage from '@/pages/auth/LoginPage';
import NotFoundPage from '@/pages/misc/NotFoundPage';
import ForbiddenPage from '@/pages/misc/ForbiddenPage';

const router = createBrowserRouter([
  {
    path: '/',
    element: <RootLayout />,
    errorElement: <NotFoundPage />,
    children: [
      { index: true, element: <HomePage /> },
      { path: 'login', element: <LoginPage /> }
    ],
  },
  { 
    path: '/403',
    element: <ForbiddenPage /> 
  }
]);

export default function AppRouter() {
  return <RouterProvider router={router} />;
}
