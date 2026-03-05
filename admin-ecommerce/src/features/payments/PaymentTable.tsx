import React from 'react';
import { motion, AnimatePresence } from 'framer-motion';
import Badge from '../../components/ui/Badge';
import Button from '../../components/ui/Button';
import Card from '../../components/ui/Card';
import { CreditCard, Calendar, Hash, ArrowRight, XCircle, Package } from 'lucide-react';
import { Payment } from '../../types';

interface PaymentTableProps {
  payments: Payment[];
  filterPending: boolean;
  setFilterPending: (val: boolean) => void;
  onProcess: (id: string) => void;
  onReject: (id: string) => void;
  processing: Record<string, boolean>;
}

const PaymentTable: React.FC<PaymentTableProps> = ({ 
  payments, filterPending, setFilterPending, onProcess, onReject, processing 
}) => {
  const formatDate = (dateString: string) => {
    if (!dateString) return "-";
    return new Date(dateString).toLocaleString('pt-BR', { 
      day: '2-digit', month: '2-digit', year: '2-digit', hour: '2-digit', minute: '2-digit' 
    });
  };

  const tableVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.05
      }
    }
  };

  const rowVariants = {
    hidden: { x: -10, opacity: 0 },
    visible: { x: 0, opacity: 1 }
  };

  return (
    <Card padding="p-0 overflow-hidden">
      <div className="flex p-4 md:p-6 border-b border-slate-50 gap-2 overflow-x-auto no-scrollbar">
        <Button 
          variant={!filterPending ? 'primary' : 'secondary'} 
          size="sm" 
          onClick={() => setFilterPending(false)}
          className="whitespace-nowrap"
        >
          Todos
        </Button>
        <Button 
          variant={filterPending ? 'warning' : 'secondary'} 
          size="sm" 
          onClick={() => setFilterPending(true)}
          className="whitespace-nowrap"
        >
          Pendentes
        </Button>
      </div>

      <div className="overflow-x-auto">
        <table className="w-full text-left min-w-[600px]">
          <thead>
            <tr className="bg-slate-50/50 border-b border-slate-100">
              <th className="p-4 md:p-6 text-[10px] font-black text-slate-400 uppercase tracking-widest">Identificação</th>
              <th className="p-4 md:p-6 text-[10px] font-black text-slate-400 uppercase tracking-widest">Valor</th>
              <th className="p-4 md:p-6 text-[10px] font-black text-slate-400 uppercase tracking-widest">Status</th>
              <th className="p-4 md:p-6 text-[10px] font-black text-slate-400 uppercase tracking-widest text-right">Ações</th>
            </tr>
          </thead>
          <motion.tbody 
            variants={tableVariants}
            initial="hidden"
            animate="visible"
            className="divide-y divide-slate-50"
          >
            <AnimatePresence mode="popLayout">
              {payments.filter(p => !filterPending || p.status === 'PENDING').map(p => (
                <motion.tr 
                  key={p.id} 
                  variants={rowVariants}
                  layout
                  initial="hidden"
                  animate="visible"
                  exit={{ opacity: 0, x: 20 }}
                  className="hover:bg-slate-50/50 group transition-colors"
                >
                  <td className="p-4 md:p-6">
                    <div className="flex flex-col">
                      <div className="flex items-center gap-2">
                        <Hash size={12} className="text-indigo-400" />
                        <span className="text-[10px] font-mono text-indigo-500 bg-indigo-50 px-2 py-0.5 rounded font-bold">
                          {(p.id || '').substring(0, 8)}
                        </span>
                      </div>
                      <span className="font-bold text-slate-700 mt-2 flex items-center gap-1.5">
                        <Package size={14} className="text-slate-400" /> 
                        Pedido #{(p.order_id || '').substring(0, 5)}
                      </span>
                      <span className="text-[11px] text-slate-400 mt-1 flex items-center gap-1.5">
                        <Calendar size={12} /> {formatDate(p.created_at)}
                      </span>
                    </div>
                  </td>
                  <td className="p-4 md:p-6">
                    <span className="text-lg font-black text-slate-900 block">
                      R$ {Number(p.amount || 0).toFixed(2)}
                    </span>
                    <span className="text-[10px] font-bold text-indigo-500 uppercase tracking-tight flex items-center gap-1.5 mt-0.5">
                      <CreditCard size={12} /> {p.method}
                    </span>
                  </td>
                  <td className="p-4 md:p-6">
                    <Badge variant={p.status}>{p.status}</Badge>
                  </td>
                  <td className="p-4 md:p-6 text-right">
                    {p.status === 'PENDING' ? (
                      <div className="flex justify-end gap-2">
                        <Button 
                          size="sm" 
                          onClick={() => onProcess(p.id)} 
                          loading={processing[p.id]}
                          className="px-3"
                        >
                          <ArrowRight size={14} className="md:mr-2" />
                          <span className="hidden md:inline">ENVIAR</span>
                        </Button>
                        <Button 
                          variant="danger" 
                          size="sm" 
                          onClick={() => onReject(p.id)} 
                          loading={processing[p.id]}
                          className="px-3"
                        >
                          <XCircle size={14} className="md:mr-2" />
                          <span className="hidden md:inline">RECUSAR</span>
                        </Button>
                      </div>
                    ) : (
                      <span className="text-[10px] font-bold text-slate-300 uppercase italic tracking-wider">Concluído</span>
                    )}
                  </td>
                </motion.tr>
              ))}
            </AnimatePresence>
          </motion.tbody>
        </table>
      </div>
      {payments.length === 0 && (
        <div className="p-20 text-center text-slate-400 font-medium">Nenhuma transação encontrada.</div>
      )}
    </Card>
  );
};

export default PaymentTable;
