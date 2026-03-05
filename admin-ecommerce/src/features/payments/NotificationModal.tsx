import React from 'react';
import Modal from '../../components/ui/Modal';
import Button from '../../components/ui/Button';
import { Payment } from '../../types';

interface NotificationModalProps {
  payment: Payment | null;
  onClose: () => void;
}

const NotificationModal: React.FC<NotificationModalProps> = ({ payment, onClose }) => {
  return (
    <Modal
      isOpen={!!payment}
      onClose={onClose}
      title="Novo Pagamento!"
    >
      <div className="flex flex-col items-center text-center">
        <div className="w-20 h-20 bg-indigo-100 text-indigo-600 rounded-3xl flex items-center justify-center text-3xl mb-6 shadow-inner">
          🔔
        </div>
        <p className="text-slate-500 mb-8 leading-relaxed">
          Um novo pagamento de <span className="font-black text-slate-900 underline decoration-indigo-200">R$ {Number(payment?.amount || 0).toFixed(2)}</span> chegou e aguarda processamento.
        </p>
        <Button onClick={onClose} className="w-full" size="lg">
          ENTENDIDO
        </Button>
      </div>
    </Modal>
  );
};

export default NotificationModal;
