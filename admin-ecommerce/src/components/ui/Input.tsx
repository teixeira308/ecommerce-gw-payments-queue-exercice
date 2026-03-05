import React, { InputHTMLAttributes } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
}

const Input: React.FC<InputProps> = ({ label, error, className = '', ...props }) => {
  return (
    <div className={`w-full ${className}`}>
      {label && (
        <label className="text-[10px] font-black text-slate-400 uppercase mb-1.5 block tracking-wider">
          {label}
        </label>
      )}
      <input
        className={`w-full bg-slate-50 border ${error ? 'border-rose-300 ring-rose-500/10' : 'border-slate-200 focus:ring-indigo-500/10'} rounded-xl px-4 py-3 text-sm font-medium transition-all focus:outline-none focus:ring-4 focus:bg-white`}
        {...props}
      />
      {error && <span className="text-[10px] text-rose-500 font-bold mt-1 block">{error}</span>}
    </div>
  );
};

export default Input;
