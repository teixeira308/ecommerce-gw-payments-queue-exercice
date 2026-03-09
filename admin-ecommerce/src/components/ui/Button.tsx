import React, { ButtonHTMLAttributes } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger' | 'success' | 'warning';
  size?: 'sm' | 'md' | 'lg' | 'icon';
  loading?: boolean;
}

const Button: React.FC<ButtonProps> = ({ 
  children, 
  variant = 'primary', 
  size = 'md', 
  className = '', 
  loading = false, 
  ...props 
}) => {
  const baseStyles = 'inline-flex items-center justify-center font-bold transition-all active:scale-95 disabled:opacity-50 disabled:pointer-events-none rounded-xl';
  
  const variants = {
    primary: 'bg-indigo-600 text-white hover:bg-indigo-700 shadow-lg shadow-indigo-200',
    secondary: 'bg-white border border-slate-200 text-slate-600 hover:bg-slate-50',
    ghost: 'bg-transparent text-slate-400 hover:text-indigo-600 hover:bg-indigo-50',
    danger: 'bg-rose-50 text-rose-600 border border-rose-100 hover:bg-rose-100',
    success: 'bg-emerald-500 text-white hover:bg-emerald-600 shadow-lg shadow-emerald-200',
    warning: 'bg-orange-500 text-white hover:bg-orange-600 shadow-lg shadow-orange-200',
  };

  const sizes = {
    sm: 'px-3 py-1.5 text-[10px] uppercase tracking-wider',
    md: 'px-6 py-2.5 text-xs',
    lg: 'px-8 py-4 text-sm',
    icon: 'p-2.5',
  };

  return (
    <button 
      className={`${baseStyles} ${variants[variant]} ${sizes[size]} ${className}`}
      disabled={loading}
      {...props}
    >
      {loading ? (
        <div className="w-4 h-4 border-2 border-current border-t-transparent rounded-full animate-spin mr-2" />
      ) : null}
      {children}
    </button>
  );
};

export default Button;
