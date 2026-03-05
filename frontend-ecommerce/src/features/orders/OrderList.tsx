import { useQuery } from '@tanstack/react-query';
import { apiService } from '../../services/api';

export const OrderList = () => {
  const { data: orders, isLoading, isError } = useQuery({
    queryKey: ['orders'],
    queryFn: apiService.getOrders,
    refetchInterval: 10000, // Revalidação a cada 10s em vez de 5s para performance
  });

  if (isLoading) {
    return <div className="text-center p-20 text-slate-500 font-black uppercase tracking-[0.4em]">Carregando pedidos...</div>;
  }

  if (isError) {
    return <div className="text-center p-20 text-red-400 font-black uppercase tracking-[0.4em]">Erro ao carregar histórico.</div>;
  }

  return (
    <div className="max-w-4xl mx-auto">
      <header className="mb-20 text-center">
        <span className="text-[10px] font-black text-indigo-400 uppercase tracking-[0.3em] mb-3 block">Orders Dashboard</span>
        <h2 className="text-5xl font-black text-white tracking-tight">Meus Pedidos</h2>
        <div className="w-16 h-1.5 bg-gradient-to-r from-indigo-600 to-purple-600 mt-6 mx-auto rounded-full"></div>
      </header>

      <div className="space-y-10">
        {orders?.slice().reverse().map((order) => (
          <div key={order.id} className="bg-white/[0.02] backdrop-blur-md rounded-[3.5rem] border border-white/5 shadow-2xl overflow-hidden group hover:bg-white/[0.04] hover:border-white/10 transition-all duration-700">
            <div className="p-12 flex flex-col md:flex-row justify-between gap-10">
              <div className="flex-1">
                <div className="flex items-center gap-5 mb-8">
                  <span className={`px-6 py-2 rounded-full text-[10px] font-black uppercase tracking-widest border ${order.status === 'APPROVED' ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : 'bg-amber-500/10 text-amber-400 border-amber-500/20'}`}>
                    {order.status}
                  </span>
                  <span className="text-xs text-slate-500 font-bold uppercase tracking-tighter">
                    {new Date(order.created_at).toLocaleString('pt-BR')}
                  </span>
                </div>
                <h4 className="text-[10px] font-mono text-slate-600 mb-10 uppercase tracking-[0.3em]">REF: {order.id}</h4>
                <div className="space-y-4">
                  {order.items.map((item, idx) => (
                    <div key={idx} className="flex items-center justify-between bg-white/[0.03] px-6 py-4 rounded-[1.5rem] border border-white/5">
                      <span className="text-sm text-slate-300 font-bold">
                        {item.quantity}x <span className="font-medium text-slate-500 ml-3">Produto SKU-{item.item_id.substring(0, 6)}</span>
                      </span>
                      <span className="text-sm font-black text-white">R$ {Number(item.subtotal).toFixed(2)}</span>
                    </div>
                  ))}
                </div>
              </div>
              <div className="md:w-64 bg-white/[0.03] rounded-[2.5rem] p-10 flex flex-col justify-center items-center text-center border border-white/5 group-hover:border-indigo-500/20 transition-all">
                <span className="text-[10px] text-slate-500 font-black uppercase mb-3 tracking-[0.2em]">Total</span>
                <span className="text-4xl font-black text-white tracking-tighter">R$ {Number(order.total).toFixed(2)}</span>
                <span className="text-[10px] text-slate-500 font-black uppercase mt-6 mb-3 tracking-[0.2em]">Pagamento</span>
                <span className="text-2xl font-black text-white tracking-tighter">{order.method}</span>
              </div>
            </div>
          </div>
        ))}
      </div>
      {orders?.length === 0 && (
        <div className="p-32 text-center bg-white/[0.02] rounded-[4rem] border border-dashed border-white/10 text-slate-600 font-black uppercase text-xs tracking-[0.4em]">
          Nenhuma atividade registrada.
        </div>
      )}
    </div>
  );
};
