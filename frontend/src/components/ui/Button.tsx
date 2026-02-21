import { ButtonHTMLAttributes } from 'react';

type Variant = 'primary' | 'secondary' | 'ghost';

type Props = ButtonHTMLAttributes<HTMLButtonElement> & {
  variant?: Variant;
  isLoading?: boolean;
};

const styles: Record<Variant, string> = {
  primary: 'bg-black text-white hover:opacity-90',
  secondary: 'bg-gray-100 text-gray-900 hover:bg-gray-200',
  ghost: 'bg-transparent hover:bg-gray-100',
};

export default function Button({ variant = 'primary', isLoading, className = '', disabled, ...props }: Props) {
  return (
    <button
      {...props}
      disabled={disabled || isLoading}
      className={`inline-flex items-center justify-center rounded-lg px-4 py-2 text-sm font-medium transition
      disabled:opacity-50 disabled:cursor-not-allowed ${styles[variant]} ${className}`}
    >
      {isLoading ? 'Loading...' : props.children}
    </button>
  );
}
