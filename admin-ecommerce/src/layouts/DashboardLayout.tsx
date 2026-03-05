import React, { useState } from 'react';
import { Menu, X, DollarSign, Package, RefreshCw } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';
import Button from '../components/ui/Button';
import { Stats } from '../types';

interface DashboardLayoutProps {
  children: React.ReactNode;
  activeTab: 'transactions' | 'catalog';
  setActiveTab: (tab: 'transactions' | 'catalog') => void;
  stats: Stats;
  autoRefresh: boolean;
  setAutoRefresh: (val: boolean) => void;
  onRefresh: () => void;
  loading: boolean;
}

const DashboardLayout: React.FC<DashboardLayoutProps> = ({ 
  children, activeTab, setActiveTab, stats, autoRefresh, setAutoRefresh, onRefresh, loading 
}) => {
  const [isSidebarOpen, setIsSidebarOpen] = useState(false);

  const menuItems = [
    { id: 'transactions' as const, label: 'Transações', icon: <DollarSign size={20} /> },
    { id: 'catalog' as const, label: 'Catálogo', icon: <Package size={20} /> },
  ];

  return (
    <div className="min-h-screen bg-slate-50 flex flex-col md:flex-row font-sans">
      {/* Mobile Header */}
      <motion.div 
        initial={{ y: -100 }}
        animate={{ y: 0 }}
        className="md:hidden bg-slate-900 text-white p-4 flex justify-between items-center fixed top-0 left-0 right-0 z-[70] shadow-lg h-16"
      >
        <div className="flex items-center gap-2">
          <div className="w-8 h-8 bg-indigo-500 rounded-lg flex items-center justify-center font-black text-[10px]">ECO</div>
          <h1 className="text-sm font-bold tracking-tight uppercase">Admin</h1>
        </div>
        <button 
          onClick={() => setIsSidebarOpen(!isSidebarOpen)} 
          className="p-2 hover:bg-slate-800 rounded-lg transition-colors active:scale-95"
        >
          {isSidebarOpen ? <X size={24} /> : <Menu size={24} />}
        </button>
      </motion.div>

      {/* Sidebar Overlay */}
      <AnimatePresence>
        {isSidebarOpen && (
          <motion.div 
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-[50] md:hidden" 
            onClick={() => setIsSidebarOpen(false)}
          />
        )}
      </AnimatePresence>

      {/* Sidebar */}
      <aside className={`
        fixed inset-y-0 left-0 w-72 bg-slate-900 text-white flex flex-col z-[60] md:relative md:translate-x-0 md:w-64
        transition-transform duration-300 ease-in-out
        ${isSidebarOpen ? 'translate-x-0' : '-translate-x-full md:translate-x-0'}
      `}>
        <div className="hidden md:flex p-8 items-center gap-3">
          <motion.div 
            whileHover={{ rotate: 15 }}
            className="w-10 h-10 bg-indigo-500 rounded-xl flex items-center justify-center font-black text-xs shadow-lg shadow-indigo-500/20"
          >
            ECO
          </motion.div>
          <h1 className="text-lg font-bold tracking-tight">Admin</h1>
        </div>

        <div className="md:hidden h-20" /> 

        <nav className="flex-1 px-4 space-y-2 overflow-y-auto">
          {menuItems.map((item) => (
            <motion.button
              key={item.id}
              whileHover={{ x: 4 }}
              whileTap={{ scale: 0.98 }}
              onClick={() => {
                setActiveTab(item.id);
                setIsSidebarOpen(false);
              }}
              className={`w-full flex items-center gap-3 px-4 py-4 md:py-3 rounded-xl text-sm font-bold transition-all ${
                activeTab === item.id 
                ? 'bg-indigo-600 text-white shadow-lg shadow-indigo-600/20' 
                : 'text-slate-400 hover:text-white hover:bg-slate-800'
              }`}
            >
              {item.icon}
              {item.label}
            </motion.button>
          ))}
        </nav>

        <div className="p-6 border-t border-slate-800 space-y-4">
          <div className="flex items-center justify-between text-[10px] font-black text-slate-500 uppercase tracking-widest">
            <span>Auto-Update</span>
            <button 
              onClick={() => setAutoRefresh(!autoRefresh)}
              className={`w-8 h-4 rounded-full relative transition-colors ${autoRefresh ? 'bg-indigo-500' : 'bg-slate-700'}`}
            >
              <motion.div 
                animate={{ left: autoRefresh ? '1.125rem' : '0.125rem' }}
                className="absolute top-0.5 w-3 h-3 bg-white rounded-full" 
              />
            </button>
          </div>
          <Button 
            variant="ghost" 
            size="sm" 
            className="w-full justify-start px-0 hover:bg-transparent" 
            onClick={() => { onRefresh(); setIsSidebarOpen(false); }}
            loading={loading}
          >
            <RefreshCw size={14} className="mr-2" /> Sincronizar Agora
          </Button>
        </div>
      </aside>

      {/* Main Content */}
      <main className="flex-1 p-6 md:p-10 overflow-x-hidden mt-16 md:mt-0">
        <header className="flex flex-col lg:flex-row justify-between lg:items-center gap-6 mb-8 md:mb-10">
          <motion.div
            initial={{ opacity: 0, x: -20 }}
            animate={{ opacity: 1, x: 0 }}
            transition={{ duration: 0.5 }}
          >
            <h2 className="text-2xl md:text-3xl font-black text-slate-900 tracking-tight capitalize">
              {activeTab === 'transactions' ? 'Transações' : 'Catálogo'}
            </h2>
            <p className="text-slate-500 text-sm mt-1">Gerencie seu ecommerce com facilidade.</p>
          </motion.div>

          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 md:gap-4">
            {[
              { label: 'Pendentes', value: stats.pending, color: 'text-orange-500' },
              { label: 'Vendas Totais', value: `R$ ${stats.totalAmount.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}`, color: 'text-indigo-600' }
            ].map((stat, idx) => (
              <motion.div 
                key={stat.label}
                initial={{ opacity: 0, y: 20 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ delay: 0.1 * idx }}
                whileHover={{ y: -4 }}
                className="bg-white p-5 rounded-3xl border border-slate-100 shadow-sm flex flex-col justify-center min-w-[180px]"
              >
                <p className="text-[10px] uppercase text-slate-400 font-bold mb-1">{stat.label}</p>
                <p className={`text-xl font-black ${stat.color}`}>{stat.value}</p>
              </motion.div>
            ))}
          </div>
        </header>

        {/* Animate tab content change */}
        <AnimatePresence mode="wait">
          <motion.div
            key={activeTab}
            initial={{ opacity: 0, y: 10 }}
            animate={{ opacity: 1, y: 0 }}
            exit={{ opacity: 0, y: -10 }}
            transition={{ duration: 0.2 }}
          >
            {children}
          </motion.div>
        </AnimatePresence>
      </main>
    </div>
  );
};

export default DashboardLayout;
