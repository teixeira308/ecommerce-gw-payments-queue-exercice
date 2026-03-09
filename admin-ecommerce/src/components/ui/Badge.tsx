import React from 'react';

interface BadgeProps {
  children: React.ReactNode;
  variant?: 'APPROVED' | 'REJECTED' | 'PENDING' | string;
  className?: string;
}

const Badge: React.FC<BadgeProps> = ({ children, variant = 'PENDING', className = '' }) => {
  const variants: Record<string, string> = {
    APPROVED: 'bg-emerald-50 text-emerald-600 border-emerald-100',
    REJECTED: 'bg-rose-50 text-rose-600 border-rose-100',
    PENDING: 'bg-orange-50 text-orange-600 border-orange-100 animate-pulse',
    default: 'bg-slate-50 text-slate-500 border-slate-100'
  };

  return (
    <span className={`px-3 py-1 rounded-full text-[10px] font-black uppercase border tracking-wider ${variants[variant] || variants.default} ${className}`}>
      {children}
    </span>
  );
};

export default Badge;
