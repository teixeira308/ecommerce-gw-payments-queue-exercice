import { motion, AnimatePresence } from 'framer-motion';
import { ShoppingBag, X, CreditCard, Smartphone, FileText, CheckCircle2 } from 'lucide-react';
import { useCartStore } from '../../store/useCartStore';
import { apiService } from '../../services/api';
import { useMutation, useQueryClient } from '@tanstack/react-query';
import { PaymentMethod } from '../../types';
import { Button } from '../../components/ui/Button';
import { toast } from 'sonner';

export const Cart = ({ onCheckoutSuccess }: { onCheckoutSuccess: () => void }) => {
  const { cart, total, paymentMethod, setPaymentMethod, removeFromCart, clearCart } = useCartStore();
  const queryClient = useQueryClient();

  const mutation = useMutation({
    mutationFn: () => apiService.createOrder(
      cart.map(i => ({ item_id: i.id, quantity: i.quantity })),
      paymentMethod
    ),
    onSuccess: () => {
      toast.success("Pedido realizado com sucesso!", {
        description: "Seu pedido foi processado e já está na fila.",
        icon: <CheckCircle2 className="w-5 h-5 text-emerald-500" />,
      });
      clearCart();
      queryClient.invalidateQueries({ queryKey: ['orders'] });
      onCheckoutSuccess();
    },
    onError: () => toast.error("Erro na conexão", {
      description: "Não foi possível sincronizar seu pedido. Tente novamente."
    })
  });

  const paymentMethods: { id: PaymentMethod, label: string, icon: React.ReactNode }[] = [
    { id: 'CREDIT_CARD', label: 'Cartão de Crédito', icon: <CreditCard className="w-5 h-5" /> },
    { id: 'PIX', label: 'Pix', icon: <Smartphone className="w-5 h-5" /> },
    { id: 'BOLETO', label: 'Boleto', icon: <FileText className="w-5 h-5" /> },
  ];

  return (
    <div className="bg-white/[0.03] backdrop-blur-3xl p-8 rounded-[3rem] shadow-2xl border border-white/10 sticky top-28">
      <div className="flex items-center justify-between mb-10">
        <div className="flex items-center gap-3">
          <ShoppingBag className="w-6 h-6 text-indigo-400" />
          <h2 className="text-2xl font-black text-white">Sacola</h2>
        </div>
        <span className="bg-indigo-600/20 text-indigo-400 text-[10px] font-black px-4 py-1.5 rounded-full border border-indigo-500/30 uppercase tracking-widest">
          {cart.length} itens
        </span>
      </div>

      {cart.length === 0 ? (
        <div className="py-20 text-center">
          <div className="w-20 h-20 bg-white/5 rounded-full flex items-center justify-center mx-auto mb-6 text-slate-700">
            <ShoppingBag className="w-10 h-10" />
          </div>
          <p className="text-[10px] font-black text-slate-500 uppercase tracking-widest leading-relaxed">
            Adicione produtos para <br /> começar sua jornada.
          </p>
        </div>
      ) : (
        <>
          <div className="space-y-6 mb-10 max-h-[350px] overflow-y-auto pr-2 custom-scrollbar">
            <AnimatePresence mode="popLayout">
              {cart.map(i => (
                <motion.div 
                  layout
                  initial={{ opacity: 0, x: 20 }}
                  animate={{ opacity: 1, x: 0 }}
                  exit={{ opacity: 0, scale: 0.8 }}
                  key={i.id} 
                  className="flex justify-between items-center group"
                >
                  <div className="flex flex-col">
                    <span className="text-sm font-bold text-white mb-1">{i.quantity}x {i.name}</span>
                    <span className="text-[9px] text-slate-600 font-mono tracking-widest uppercase">ID: {i.id.substring(0, 8)}</span>
                  </div>
                  <div className="flex items-center gap-4">
                    <span className="text-sm font-black text-indigo-400">R$ {(i.price * i.quantity).toFixed(2)}</span>
                    <Button 
                      variant="ghost" 
                      size="icon" 
                      className="w-8 h-8"
                      onClick={() => removeFromCart(i.id)} 
                    >
                      <X className="w-4 h-4" />
                    </Button>
                  </div>
                </motion.div>
              ))}
            </AnimatePresence>
          </div>

          <div className="bg-white/5 rounded-[2rem] p-6 mb-10 border border-white/5">
            <div className="flex justify-between items-center mb-2">
              <span className="text-[10px] font-black text-slate-600 uppercase tracking-widest">Subtotal</span>
              <span className="text-sm font-bold text-slate-400">R$ {total.toFixed(2)}</span>
            </div>
            <div className="flex justify-between items-center">
              <span className="text-[10px] font-black text-slate-500 uppercase tracking-widest">Total Final</span>
              <span className="text-3xl font-black text-white tracking-tighter">R$ {total.toFixed(2)}</span>
            </div>
          </div>

          <div className="mb-10">
            <p className="text-[10px] font-black text-slate-500 uppercase mb-4 tracking-[0.2em]">Método de Pagamento</p>
            <div className="grid grid-cols-1 gap-2.5">
              {paymentMethods.map(m => (
                <button 
                  key={m.id} 
                  onClick={() => setPaymentMethod(m.id)} 
                  className={`flex items-center justify-between p-4 rounded-2xl border transition-all text-xs font-bold ${paymentMethod === m.id ? 'bg-indigo-600 border-indigo-400 text-white shadow-lg shadow-indigo-500/20' : 'bg-white/5 border-white/5 text-slate-400 hover:bg-white/10'}`}
                >
                  <div className="flex items-center gap-4">
                    {m.icon}
                    {m.label}
                  </div>
                  {paymentMethod === m.id && <div className="w-1.5 h-1.5 bg-white rounded-full shadow-[0_0_8px_white]"></div>}
                </button>
              ))}
            </div>
          </div>

          <Button 
            className="w-full"
            size="lg"
            onClick={() => mutation.mutate()} 
            disabled={mutation.isPending} 
          >
            {mutation.isPending ? 'Sincronizando...' : 'Finalizar Pedido'}
          </Button>
        </>
      )}
    </div>
  );
};
