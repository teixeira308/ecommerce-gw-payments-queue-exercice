import React from 'react';
import { Payment } from '../types';

interface PaymentModalProps {
  payment: Payment | null;
  onClose: () => void;
}

const PaymentModal: React.FC<PaymentModalProps> = ({ payment, onClose }) => {
  if (!payment) return null;
  return (
    <div 
      className="fixed inset-0 bg-black/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4 animate-in fade-in duration-300"
      role="dialog"
      aria-modal="true"
    >
      <div className="bg-white rounded-[2.5rem] p-8 max-w-md w-full shadow-2xl border border-slate-100 animate-in zoom-in-95 duration-300">
        <div className="flex flex-col items-center text-center">
          <div className="w-16 h-16 bg-emerald-100 text-emerald-600 rounded-2xl flex items-center justify-center text-2xl mb-4">
            💸
          </div>
          <h2 className="text-2xl font-black text-slate-900 mb-2">Novo Pagamento!</h2>
          <p className="text-slate-500 mb-6">Um novo pagamento de <span className="font-bold text-slate-900">R$ {Number(payment.amount).toFixed(2)}</span> foi recebido pelo gateway.</p>
          <div className="w-full bg-slate-50 rounded-2xl p-4 mb-6 text-left">
            <div className="flex justify-between mb-1">
              <span className="text-[10px] font-black text-slate-400 uppercase">ID Interno</span>
              <span className="text-[10px] font-mono text-slate-600">#{payment.id.substring(0, 8)}</span>
            </div>
            <div className="flex justify-between">
              <span className="text-[10px] font-black text-slate-400 uppercase">Pedido</span>
              <span className="text-[10px] font-bold text-slate-700 uppercase">#{payment.order_id.substring(0, 8)}</span>
            </div>
          </div>
          <button
            onClick={onClose}
            className="w-full bg-indigo-600 text-white py-4 rounded-2xl font-black text-sm hover:bg-indigo-700 transition-all active:scale-95 shadow-xl shadow-indigo-200"
          >
            ENTENDIDO
          </button>
        </div>
      </div>
    </div>
  );
};

export default PaymentModal;
