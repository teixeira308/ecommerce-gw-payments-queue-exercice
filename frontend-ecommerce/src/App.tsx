import { useState } from 'react';
import { useQuery } from '@tanstack/react-query';
import { motion, AnimatePresence } from 'framer-motion';
import { ShoppingBag, LayoutDashboard, Store } from 'lucide-react';
import { apiService } from './services/api';
import { ProductCard } from './features/shop/ProductCard';
import { Cart } from './features/shop/Cart';
import { OrderList } from './features/orders/OrderList';

export default function App() {
  const [activeTab, setActiveTab] = useState<'shop' | 'orders'>('shop');
  
  const { data: products, isLoading } = useQuery({
    queryKey: ['products'],
    queryFn: apiService.getProducts,
  });

  return (
    <div className="min-h-screen bg-slate-950 text-slate-200 selection:bg-indigo-500/30 font-sans antialiased">
      {/* Background Decorativo - SaaS Premium */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-[-10%] right-[-10%] w-[60%] h-[60%] bg-indigo-900/10 rounded-full blur-[140px]"></div>
        <div className="absolute bottom-[-10%] left-[-10%] w-[60%] h-[60%] bg-blue-900/10 rounded-full blur-[140px]"></div>
        <div className="absolute inset-0 bg-[url('https://grainy-gradients.vercel.app/noise.svg')] opacity-20 mix-blend-overlay"></div>
      </div>

      <nav className="bg-slate-950/40 backdrop-blur-2xl border-b border-white/5 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-6 h-20 flex items-center justify-between">
          <div className="flex items-center gap-4">
            <div className="w-11 h-11 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-2xl flex items-center justify-center text-white font-black shadow-2xl shadow-indigo-500/30">G</div>
            <div className="flex flex-col">
              <h1 className="text-lg font-black tracking-tighter text-white leading-none">GO.SHOP</h1>
              <span className="text-[9px] font-bold text-indigo-400 tracking-[0.2em] uppercase">Enterprise</span>
            </div>
          </div>
          
          <div className="flex bg-white/5 p-1.5 rounded-[1.25rem] border border-white/10">
            <button 
              onClick={() => setActiveTab('shop')} 
              className={`flex items-center gap-2 px-6 py-2.5 rounded-xl text-xs font-black transition-all relative ${activeTab === 'shop' ? 'text-white' : 'text-slate-500 hover:text-slate-300'}`}
            >
              {activeTab === 'shop' && (
                <motion.div layoutId="nav-pill" className="absolute inset-0 bg-indigo-600 rounded-xl shadow-lg shadow-indigo-500/20" />
              )}
              <Store className="w-4 h-4 relative z-10" />
              <span className="relative z-10">LOJA</span>
            </button>
            <button 
              onClick={() => setActiveTab('orders')} 
              className={`flex items-center gap-2 px-6 py-2.5 rounded-xl text-xs font-black transition-all relative ${activeTab === 'orders' ? 'text-white' : 'text-slate-500 hover:text-slate-300'}`}
            >
              {activeTab === 'orders' && (
                <motion.div layoutId="nav-pill" className="absolute inset-0 bg-indigo-600 rounded-xl shadow-lg shadow-indigo-500/20" />
              )}
              <LayoutDashboard className="w-4 h-4 relative z-10" />
              <span className="relative z-10">PEDIDOS</span>
            </button>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto p-6 md:p-12 relative z-10">
        <AnimatePresence mode="wait">
          {activeTab === 'shop' ? (
            <motion.div 
              key="shop"
              initial={{ opacity: 0, x: -20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: 20 }}
              className="grid grid-cols-1 lg:grid-cols-12 gap-12 items-start"
            >
              <div className="lg:col-span-8">
                {isLoading ? (
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-8 animate-pulse">
                    {[1, 2, 3, 4].map(i => (
                      <div key={i} className="h-[400px] bg-white/5 rounded-[3rem] border border-white/10"></div>
                    ))}
                  </div>
                ) : (
                  <div className="grid grid-cols-1 sm:grid-cols-2 gap-8">
                    {products?.map(product => (
                      <ProductCard key={product.id} product={product} />
                    ))}
                  </div>
                )}
              </div>

              <div className="lg:col-span-4">
                <Cart onCheckoutSuccess={() => setActiveTab('orders')} />
              </div>
            </motion.div>
          ) : (
            <motion.div
              key="orders"
              initial={{ opacity: 0, x: 20 }}
              animate={{ opacity: 1, x: 0 }}
              exit={{ opacity: 0, x: -20 }}
            >
              <OrderList />
            </motion.div>
          )}
        </AnimatePresence>
      </main>
    </div>
  );
}
