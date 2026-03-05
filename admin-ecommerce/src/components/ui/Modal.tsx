import React from 'react';

interface ModalProps {
  isOpen: boolean;
  onClose: () => void;
  title: string;
  children: React.ReactNode;
  maxWidth?: string;
}

const Modal: React.FC<ModalProps> = ({ isOpen, onClose, title, children, maxWidth = 'max-w-md' }) => {
  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-slate-900/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4 animate-in fade-in duration-300">
      <div className={`bg-white rounded-[2.5rem] p-8 ${maxWidth} w-full shadow-2xl border border-slate-100 animate-in zoom-in-95 duration-300`}>
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-black text-slate-900 leading-tight">{title}</h2>
          <button onClick={onClose} className="p-2 hover:bg-slate-50 rounded-xl transition-colors text-slate-400">
            ✕
          </button>
        </div>
        {children}
      </div>
    </div>
  );
};

export default Modal;
