import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import { CheckCircle2, XCircle, Clock, CreditCard, ChevronRight, Hash, Calendar, DollarSign } from 'lucide-react';
import { Payment } from '../types';
import { clsx, type ClassValue } from 'clsx';
import { twMerge } from 'tailwind-merge';

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

interface PaymentTableProps {
  payments: Payment[];
  processing: Record<string, boolean>;
  onAction: (id: string, status: 'APPROVED' | 'REJECTED') => void;
}

const formatDate = (dateString: string) => {
  if (!dateString) return "-";
  return new Date(dateString).toLocaleString('pt-BR', {
    day: '2-digit',
    month: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const StatusBadge = ({ status }: { status: string }) => {
  const configs = {
    APPROVED: { icon: CheckCircle2, classes: "bg-emerald-50 text-emerald-600 border-emerald-100", label: "Aprovado" },
    REJECTED: { icon: XCircle, classes: "bg-rose-50 text-rose-600 border-rose-100", label: "Recusado" },
    PENDING: { icon: Clock, classes: "bg-amber-50 text-amber-600 border-amber-100 animate-pulse", label: "Pendente" },
  };

  const config = configs[status as keyof typeof configs] || configs.PENDING;
  const Icon = config.icon;

  return (
    <span className={cn(
      "inline-flex items-center gap-1.5 px-3 py-1 rounded-full text-[10px] font-bold uppercase tracking-wider border transition-all",
      config.classes
    )}>
      <Icon size={12} strokeWidth={3} />
      {config.label}
    </span>
  );
};

const PaymentTable: React.FC<PaymentTableProps> = ({ payments, processing, onAction }) => {
  if (payments.length === 0) {
    return (
      <div className="p-32 text-center flex flex-col items-center justify-center gap-4">
        <div className="w-16 h-16 bg-slate-50 rounded-full flex items-center justify-center text-slate-300">
           <CreditCard size={32} />
        </div>
        <div>
          <p className="text-slate-900 font-bold text-lg">Nenhum dado encontrado</p>
          <p className="text-slate-400 text-sm">Tente ajustar seus filtros para ver mais resultados.</p>
        </div>
      </div>
    );
  }

  return (
    <div className="w-full">
      {/* Desktop Table View */}
      <div className="hidden lg:block overflow-x-auto">
        <table className="w-full text-left border-collapse">
          <thead>
            <tr className="border-b border-slate-100">
              <th className="px-8 py-6 text-xs font-bold text-slate-400 uppercase tracking-[0.15em]">Transação</th>
              <th className="px-8 py-6 text-xs font-bold text-slate-400 uppercase tracking-[0.15em]">Valor</th>
              <th className="px-8 py-6 text-xs font-bold text-slate-400 uppercase tracking-[0.15em]">Status</th>
              <th className="px-8 py-6 text-xs font-bold text-slate-400 uppercase tracking-[0.15em] text-right">Ações</th>
            </tr>
          </thead>
          <tbody className="divide-y divide-slate-50">
            <AnimatePresence mode='popLayout'>
              {payments.map(p => (
                <motion.tr 
                  layout
                  initial={{ opacity: 0, y: 10 }}
                  animate={{ opacity: 1, y: 0 }}
                  exit={{ opacity: 0, scale: 0.95 }}
                  key={p.id} 
                  className="hover:bg-slate-50/50 transition-all duration-300 group"
                >
                  <td className="px-8 py-6">
                    <div className="flex items-center gap-4">
                      <div className="w-10 h-10 rounded-2xl bg-slate-100 flex items-center justify-center text-slate-400 group-hover:bg-white group-hover:shadow-sm transition-all border border-transparent group-hover:border-slate-200">
                        <Hash size={16} />
                      </div>
                      <div className="flex flex-col">
                        <span className="font-bold text-slate-900 leading-tight">Pedido #{p.order_id.substring(0, 8)}</span>
                        <div className="flex items-center gap-2 mt-1">
                          <Calendar size={10} className="text-slate-300" />
                          <span className="text-xs text-slate-400 font-medium">{formatDate(p.created_at)}</span>
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="px-8 py-6">
                    <div className="flex flex-col">
                      <div className="flex items-center gap-1">
                         <span className="text-sm font-bold text-slate-400 italic">R$</span>
                         <span className="text-xl font-black text-slate-900 tracking-tight">
                            {Number(p.amount).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
                         </span>
                      </div>
                      <span className="text-[10px] font-black text-indigo-500 uppercase tracking-widest mt-0.5 ml-0.5">
                        {p.method}
                      </span>
                    </div>
                  </td>
                  <td className="px-8 py-6">
                    <StatusBadge status={p.status} />
                  </td>
                  <td className="px-8 py-6 text-right">
                    {p.status === 'PENDING' ? (
                      <div className="flex justify-end gap-3">
                        <button
                          onClick={() => onAction(p.id, 'APPROVED')}
                          disabled={processing[p.id]}
                          className="px-6 py-2.5 bg-indigo-600 text-white rounded-xl text-xs font-bold hover:bg-indigo-700 shadow-md shadow-indigo-100 transition-all active:scale-95 disabled:opacity-50"
                        >
                          {processing[p.id] ? 'Processando...' : 'Aprovar'}
                        </button>
                        <button
                          onClick={() => onAction(p.id, 'REJECTED')}
                          disabled={processing[p.id]}
                          className="px-4 py-2.5 bg-white border border-slate-200 text-slate-500 rounded-xl text-xs font-bold hover:bg-rose-50 hover:text-rose-600 hover:border-rose-100 transition-all"
                        >
                          Recusar
                        </button>
                      </div>
                    ) : (
                      <div className="flex items-center justify-end gap-1 text-slate-300">
                         <span className="text-[10px] font-black uppercase tracking-widest italic">Encerrado</span>
                         <ChevronRight size={14} />
                      </div>
                    )}
                  </td>
                </motion.tr>
              ))}
            </AnimatePresence>
          </tbody>
        </table>
      </div>

      {/* Mobile Card View */}
      <div className="lg:hidden p-4 flex flex-col gap-4 bg-slate-50/30">
        <AnimatePresence mode='popLayout'>
          {payments.map(p => (
            <motion.div 
              layout
              initial={{ opacity: 0, scale: 0.95 }}
              animate={{ opacity: 1, scale: 1 }}
              exit={{ opacity: 0, scale: 0.9 }}
              key={p.id}
              className="bg-white p-6 rounded-3xl border border-slate-200 shadow-sm shadow-slate-200/50"
            >
              <div className="flex justify-between items-start mb-6">
                <div>
                  <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">ID Transação</p>
                  <h3 className="text-lg font-black text-slate-900 leading-none">#{p.order_id.substring(0, 8)}</h3>
                  <p className="text-xs text-slate-400 mt-1 font-medium italic">{formatDate(p.created_at)}</p>
                </div>
                <StatusBadge status={p.status} />
              </div>

              <div className="flex items-end justify-between py-4 border-y border-slate-50">
                <div>
                   <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Valor Total</p>
                   <div className="flex items-baseline gap-1">
                      <span className="text-sm font-bold text-slate-400">R$</span>
                      <span className="text-2xl font-black text-slate-900 tracking-tight">
                         {Number(p.amount).toLocaleString('pt-BR', { minimumFractionDigits: 2 })}
                      </span>
                   </div>
                </div>
                <div className="text-right">
                   <p className="text-[10px] font-black text-slate-400 uppercase tracking-widest mb-1">Método</p>
                   <span className="bg-indigo-50 text-indigo-600 px-2 py-0.5 rounded text-[10px] font-black uppercase">{p.method}</span>
                </div>
              </div>

              {p.status === 'PENDING' && (
                <div className="grid grid-cols-2 gap-3 mt-6">
                   <button
                    onClick={() => onAction(p.id, 'APPROVED')}
                    disabled={processing[p.id]}
                    className="flex items-center justify-center gap-2 py-3 bg-indigo-600 text-white rounded-2xl text-xs font-bold shadow-lg shadow-indigo-100"
                  >
                    <CheckCircle2 size={14} />
                    Aprovar
                  </button>
                  <button
                    onClick={() => onAction(p.id, 'REJECTED')}
                    disabled={processing[p.id]}
                    className="flex items-center justify-center gap-2 py-3 bg-white border border-slate-200 text-slate-500 rounded-2xl text-xs font-bold"
                  >
                    <XCircle size={14} />
                    Recusar
                  </button>
                </div>
              )}
            </motion.div>
          ))}
        </AnimatePresence>
      </div>
    </div>
  );
};

export default PaymentTable;
