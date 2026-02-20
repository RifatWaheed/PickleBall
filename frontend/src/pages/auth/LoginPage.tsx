import { useForm } from 'react-hook-form';
import { z } from 'zod';
import { zodResolver } from '@hookform/resolvers/zod';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '@/app/AuthProvider';
import Input from '@/components/ui/Input';
import Button from '@/components/ui/Button';

const schema = z.object({
  email: z.string().email('Invalid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
});

type FormData = z.infer<typeof schema>;

export default function LoginPage() {
  const navigate = useNavigate();
  const { login, user } = useAuth();

  const {
    register,
    handleSubmit,
    formState: { errors, isSubmitting },
    setError,
  } = useForm<FormData>({
    resolver: zodResolver(schema),
  });

  const onSubmit = async (data: FormData) => {
    try {
      await login(data.email, data.password);

      if (user?.role === 'ADMIN') navigate('/admin');
      else if (user?.role === 'SUPER_ADMIN') navigate('/super-admin');
      else navigate('/booking');

    } catch (err: any) {
      setError('root', {
        message: err?.response?.data?.message || 'Invalid credentials',
      });
    }
  };

  return (
    <div className="max-w-md mx-auto mt-16 space-y-6">
      <h1 className="text-2xl font-semibold">Login</h1>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
        <Input
          label="Email"
          type="email"
          {...register('email')}
          error={errors.email?.message}
        />

        <Input
          label="Password"
          type="password"
          {...register('password')}
          error={errors.password?.message}
        />

        {errors.root && (
          <p className="text-sm text-red-600">
            {errors.root.message}
          </p>
        )}

        <Button type="submit" isLoading={isSubmitting} className="w-full">
          Sign In
        </Button>
      </form>
    </div>
  );
}
