import { useState, useEffect, useCallback, useRef } from 'react'

const GATEWAY_URL = 'http://localhost:4000/gopay'

// Componente de Modal Simples
const Modal = ({ payment, onClose }) => {
  if (!payment) return null;
  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4 animate-in fade-in duration-300">
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

function App() {
  const [payments, setPayments] = useState([])
  const [page, setPage] = useState(1)
  const [limit] = useState(10)
  const [loading, setLoading] = useState(false)
  const [stats, setStats] = useState({ pending: 0, totalAmount: 0 })
  const [filterPending, setFilterPending] = useState(false)

  // Novos estados
  const [autoRefresh, setAutoRefresh] = useState(true)
  const [processing, setProcessing] = useState({})
  const [newPaymentModal, setNewPaymentModal] = useState(null)
  const [seenIds, setSeenIds] = useState(new Set())
  const isFirstLoad = useRef(true)

  // 1. Busca estatísticas de TODOS os pagamentos (sem paginação)
  const fetchGlobalStats = useCallback(async () => {
    try {
      const res = await fetch(`${GATEWAY_URL}/payments?limit=9999`)
      const data = await res.json()
      const all = data || []

      const pendingCount = all.filter(p => p.status === 'PENDING').length
      const approvedTotal = all
        .filter(p => p.status === 'APPROVED')
        .reduce((acc, curr) => acc + Number(curr.amount), 0)

      setStats({ pending: pendingCount, totalAmount: approvedTotal })
    } catch (err) {
      console.error("Erro ao buscar estatísticas globais:", err)
    }
  }, [])

  // 2. Busca pagamentos para a tabela (com paginação)
  const fetchPayments = useCallback(async () => {
    setLoading(true)
    try {
      const res = await fetch(`${GATEWAY_URL}/payments?limit=${limit}&page=${page}`)
      const data = await res.json()
      const paymentsData = data || []

      setSeenIds(prev => {
        const newOnes = paymentsData.filter(p => p.status === 'PENDING' && !prev.has(p.id))
        
        if (!isFirstLoad.current && newOnes.length > 0) {
          setNewPaymentModal(newOnes[0])
          setFilterPending(true) // Ativa o filtro de pendentes automaticamente
        }

        const next = new Set(prev)
        paymentsData.forEach(p => next.add(p.id))
        return next
      })
      isFirstLoad.current = false

      setPayments(paymentsData)
    } catch (err) {
      console.error("Erro ao buscar pagamentos:", err)
    } finally {
      setTimeout(() => setLoading(false), 300)
    }
  }, [page, limit])

  const formatDate = (dateString) => {
    if (!dateString) return "-";
    return new Date(dateString).toLocaleString('pt-BR', {
      day: '2-digit',
      month: '2-digit',
      year: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  useEffect(() => {
    fetchPayments()
    fetchGlobalStats()
  }, [fetchPayments, fetchGlobalStats])

  // Effect para atualização automática
  useEffect(() => {
    if (!autoRefresh) return;
    const interval = setInterval(() => {
      fetchPayments();
      fetchGlobalStats();
    }, 30000);
    return () => clearInterval(interval);
  }, [autoRefresh, fetchPayments, fetchGlobalStats]);

  const handleAction = async (id, status) => {
    setProcessing(prev => ({ ...prev, [id]: true }))
    try {
      const res = await fetch(`${GATEWAY_URL}/payments/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status })
      })
      if (!res.ok) throw new Error("Erro na comunicação")
      
      fetchPayments()
      fetchGlobalStats()
    } catch (err) {
      alert("Erro ao processar ação")
      setProcessing(prev => ({ ...prev, [id]: false }))
    }
  }

  // Lógica de filtragem no Frontend (para a página atual)
  const displayedPayments = filterPending
    ? payments.filter(p => p.status === 'PENDING')
    : payments

  return (
    <div className="min-h-screen bg-[#f1f5f9]">
      <Modal payment={newPaymentModal} onClose={() => setNewPaymentModal(null)} />

      {/* HEADER - ESTATÍSTICAS GLOBAIS */}
      <nav className="bg-slate-900 text-white p-6 shadow-lg sticky top-0 z-50">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-indigo-500 rounded-xl flex items-center justify-center font-black">GP</div>
            <h1 className="text-xl font-bold tracking-tight">Gateway Admin Panel</h1>
          </div>
          <div className="flex gap-6 items-center">
            {/* Toggle Auto-Refresh */}
            <div className="flex flex-col items-end mr-4">
              <span className="text-[9px] uppercase text-slate-500 font-black mb-1">Auto-refresh (30s)</span>
              <button 
                onClick={() => setAutoRefresh(!autoRefresh)}
                className={`w-10 h-5 rounded-full transition-colors relative ${autoRefresh ? 'bg-indigo-500' : 'bg-slate-700'}`}
              >
                <div className={`absolute top-1 w-3 h-3 bg-white rounded-full transition-all ${autoRefresh ? 'left-6' : 'left-1'}`} />
              </button>
            </div>

            <div className="text-right">
              <p className="text-[10px] uppercase text-slate-400 font-bold tracking-widest">Pendentes (Total)</p>
              <p className="text-xl font-black text-amber-400">{stats.pending}</p>
            </div>
            <div className="text-right border-l border-slate-700 pl-6">
              <p className="text-[10px] uppercase text-slate-400 font-bold tracking-widest">Faturamento Aprovado</p>
              <p className="text-xl font-black text-emerald-400">R$ {stats.totalAmount.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}</p>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto p-8">
        {/* FILTROS E AÇÕES */}
        <div className="flex justify-between items-center mb-6">
          <div className="flex bg-white p-1 rounded-2xl border border-slate-200 shadow-sm">
            <button
              onClick={() => setFilterPending(false)}
              className={`px-6 py-2 rounded-xl text-xs font-bold transition-all ${!filterPending ? 'bg-slate-900 text-white shadow-md' : 'text-slate-500 hover:bg-slate-50'}`}
            >
              Todos
            </button>
            <button
              onClick={() => setFilterPending(true)}
              className={`px-6 py-2 rounded-xl text-xs font-bold transition-all ${filterPending ? 'bg-amber-500 text-white shadow-md' : 'text-slate-500 hover:bg-slate-50'}`}
            >
              Pendentes
            </button>
          </div>

          <button 
            onClick={() => { fetchPayments(); fetchGlobalStats(); }} 
            disabled={loading}
            className="flex items-center gap-3 px-4 py-2 text-slate-400 hover:text-indigo-600 transition-all active:scale-95 disabled:opacity-50"
          >
            <div className="relative w-4 h-4 flex items-center justify-center">
               <span className={`absolute transition-all duration-300 ${loading ? 'opacity-0 scale-50' : 'opacity-100 scale-100'}`}>🔄</span>
               <div className={`absolute w-4 h-4 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin transition-all duration-300 ${loading ? 'opacity-100 scale-100' : 'opacity-0 scale-50'}`}></div>
            </div>
            <span className="text-xs font-bold uppercase tracking-tighter">Atualizar Dados</span>
          </button>
        </div>

        {/* TABELA */}
        <div className="bg-white rounded-[2.5rem] shadow-xl shadow-slate-200/50 border border-slate-200 overflow-hidden">
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead>
                <tr className="bg-slate-50/50 border-b border-slate-100">
                  <th className="p-6 text-xs font-black text-slate-400 uppercase tracking-widest">Identificação</th>
                  <th className="p-6 text-xs font-black text-slate-400 uppercase tracking-widest">Valor</th>
                  <th className="p-6 text-xs font-black text-slate-400 uppercase tracking-widest">Status</th>
                  <th className="p-6 text-xs font-black text-slate-400 uppercase tracking-widest text-right">Ações</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-50">
                {displayedPayments.map(p => (
                  <tr key={p.id} className="hover:bg-slate-50/80 transition-colors group">
                    <td className="p-6">
                      <div className="flex flex-col">
                        <span className="text-[10px] font-mono text-slate-400 bg-slate-100 w-fit px-2 py-0.5 rounded mb-1 uppercase">ID: {p.id.substring(0, 8)}...</span>
                        <span className="font-bold text-slate-700">Pedido #{p.order_id.substring(0, 5)}</span>
                        <span className="text-[11px] text-slate-400 mt-1">{formatDate(p.created_at)}</span>
                      </div>
                    </td>
                    <td className="p-6">
                      <div className="flex flex-col">
                        <span className="text-lg font-black text-slate-900">R$ {Number(p.amount).toFixed(2)}</span>
                        <span className="text-[10px] font-bold text-indigo-500 uppercase tracking-tighter">
                          {p.method}
                        </span>
                      </div>
                    </td>
                    <td className="p-6">
                      <span className={`px-4 py-1.5 rounded-full text-[10px] font-black tracking-tighter uppercase border ${p.status === 'APPROVED' ? 'bg-emerald-50 text-emerald-600 border-emerald-100' :
                        p.status === 'REJECTED' ? 'bg-red-50 text-red-600 border-red-100' : 'bg-amber-50 text-amber-600 border-amber-100 animate-pulse'
                        }`}>
                        {p.status}
                      </span>
                    </td>
                    <td className="p-6 text-right">
                      {p.status === 'PENDING' ? (
                        <div className="flex justify-end gap-2 opacity-0 group-hover:opacity-100 focus-within:opacity-100 transition-opacity">
                          <button
                            onClick={() => handleAction(p.id, 'APPROVED')}
                            disabled={processing[p.id]}
                            className="bg-emerald-500 text-white px-4 py-2 rounded-xl text-xs font-bold hover:bg-emerald-600 shadow-lg shadow-emerald-200 transition-all active:scale-95 disabled:opacity-50 disabled:cursor-not-allowed"
                          >
                            {processing[p.id] ? '...' : 'APROVAR'}
                          </button>
                          <button
                            onClick={() => handleAction(p.id, 'REJECTED')}
                            disabled={processing[p.id]}
                            className="bg-white border border-slate-200 text-slate-400 px-4 py-2 rounded-xl text-xs font-bold hover:bg-red-50 hover:text-red-500 hover:border-red-100 transition-all disabled:opacity-50 disabled:cursor-not-allowed"
                          >
                            RECUSAR
                          </button>
                        </div>
                      ) : (
                        <span className="text-[10px] font-bold text-slate-300 uppercase italic">Concluído</span>
                      )}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>

          {displayedPayments.length === 0 && (
            <div className="p-20 text-center">
              <p className="text-slate-400 font-medium">Nenhuma transação encontrada com este filtro.</p>
            </div>
          )}

          {/* PAGINAÇÃO */}
          <div className="bg-slate-50/50 p-6 border-t border-slate-100 flex items-center justify-between">
            <p className="text-xs font-bold text-slate-400 uppercase tracking-widest">
              Página <span className="text-indigo-600">{page}</span>
            </p>
            <div className="flex gap-3">
              <button
                onClick={() => setPage(prev => Math.max(prev - 1, 1))}
                disabled={page === 1}
                className="px-6 py-2 bg-white border border-slate-200 rounded-xl text-xs font-black disabled:opacity-30 hover:border-indigo-200 transition-all shadow-sm"
              >
                ← ANTERIOR
              </button>
              <button
                onClick={() => setPage(prev => prev + 1)}
                disabled={payments.length < limit}
                className="px-6 py-2 bg-white border border-slate-200 rounded-xl text-xs font-black hover:border-indigo-200 transition-all shadow-sm"
              >
                PRÓXIMA →
              </button>
            </div>
          </div>
        </div>
      </main>
    </div>
  )
}

export default App;