import React from 'react';
import { motion } from 'framer-motion';
import { Zap, ShieldCheck, TrendingUp, RefreshCcw, Bell } from 'lucide-react';
import { Stats } from '../types';

interface HeaderProps {
  stats: Stats;
  autoRefresh: boolean;
  onToggleAutoRefresh: () => void;
  loading?: boolean;
}

const Header: React.FC<HeaderProps> = ({ stats, autoRefresh, onToggleAutoRefresh, loading }) => {
  return (
    <header className="bg-white border-b border-slate-200 sticky top-0 z-40 backdrop-blur-md bg-white/80">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex justify-between items-center h-20">
          
          {/* Logo Section */}
          <div className="flex items-center gap-4">
            <div className="w-12 h-12 bg-indigo-600 rounded-2xl flex items-center justify-center text-white shadow-lg shadow-indigo-200 overflow-hidden relative group">
              <Zap size={24} fill="currentColor" />
              <motion.div 
                className="absolute inset-0 bg-white/20"
                animate={{ x: ['-100%', '100%'] }}
                transition={{ duration: 1.5, repeat: Infinity, ease: "linear" }}
              />
            </div>
            <div className="hidden sm:block">
              <h1 className="text-lg font-black text-slate-900 tracking-tight leading-none">GoPay Admin</h1>
              <p className="text-[10px] text-slate-400 font-bold uppercase tracking-widest mt-1">Financeiro & Gateway</p>
            </div>
          </div>

          {/* Stats & Actions */}
          <div className="flex items-center gap-3 sm:gap-8">
            <div className="hidden md:flex items-center gap-8">
              <div className="text-right">
                <div className="flex items-center justify-end gap-1.5 text-amber-500 mb-0.5">
                  <Bell size={12} strokeWidth={3} />
                  <span className="text-[10px] uppercase font-black tracking-widest">Pendentes</span>
                </div>
                <p className="text-xl font-black text-slate-900">{stats.pending}</p>
              </div>

              <div className="text-right border-l border-slate-100 pl-8">
                <div className="flex items-center justify-end gap-1.5 text-emerald-500 mb-0.5">
                  <TrendingUp size={12} strokeWidth={3} />
                  <span className="text-[10px] uppercase font-black tracking-widest">Faturamento</span>
                </div>
                <p className="text-xl font-black text-slate-900">
                  <span className="text-xs font-bold text-slate-300 mr-1">R$</span>
                  {stats.totalAmount.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
                </p>
              </div>
            </div>

            <div className="flex items-center gap-4">
               <button 
                onClick={onToggleAutoRefresh}
                className={`relative flex items-center gap-2 px-4 py-2 rounded-2xl transition-all font-bold text-xs ${
                  autoRefresh 
                  ? 'bg-indigo-50 text-indigo-600 ring-1 ring-indigo-100' 
                  : 'bg-slate-50 text-slate-400'
                }`}
              >
                <RefreshCcw size={14} className={autoRefresh && loading ? 'animate-spin' : ''} />
                <span className="hidden lg:inline">{autoRefresh ? 'Auto-update ON' : 'Auto-update OFF'}</span>
              </button>
            </div>
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;
