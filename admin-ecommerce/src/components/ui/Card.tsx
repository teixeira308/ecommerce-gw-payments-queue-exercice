import React from 'react';

interface CardProps {
  children: React.ReactNode;
  className?: string;
  padding?: string;
}

const Card: React.FC<CardProps> = ({ children, className = '', padding = 'p-6' }) => {
  return (
    <div className={`bg-white rounded-[2rem] shadow-sm border border-slate-100 hover:border-indigo-100 transition-all ${padding} ${className}`}>
      {children}
    </div>
  );
};

export default Card;
