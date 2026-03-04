import { useState, useEffect, useCallback } from 'react'

const GATEWAY_URL = 'http://localhost:9000'

function App() {
  const [payments, setPayments] = useState([])
  const [page, setPage] = useState(1)
  const [limit] = useState(10)
  const [loading, setLoading] = useState(false)
  const [stats, setStats] = useState({ pending: 0, totalAmount: 0 })
  const [filterPending, setFilterPending] = useState(false)
 
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
      setPayments(data || [])
    } catch (err) {
      console.error("Erro ao buscar pagamentos:", err)
    } finally {
      setLoading(false)
    }
  }, [page, limit])

  const formatDate = (dateString) => {
    if (!dateString) return "-";
    return new Date(dateString).toLocaleString('pt-BR', {
      day: '2-digit',    // O correto é '2-digit' e não '2d'
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

  const handleAction = async (id, status) => {
    try {
      await fetch(`${GATEWAY_URL}/payments/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status })
      })
      fetchPayments()
      fetchGlobalStats() // Atualiza o header após aprovar
    } catch (err) {
      alert("Erro ao processar ação")
    }
  }

  // Lógica de filtragem no Frontend (para a página atual)
  const displayedPayments = filterPending
    ? payments.filter(p => p.status === 'PENDING')
    : payments

  return (
    <div className="min-h-screen bg-[#f1f5f9]">
      {/* HEADER - ESTATÍSTICAS GLOBAIS */}
      <nav className="bg-slate-900 text-white p-6 shadow-lg sticky top-0 z-50">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-indigo-500 rounded-xl flex items-center justify-center font-black">GP</div>
            <h1 className="text-xl font-bold tracking-tight">Admin Ecommerce</h1>
          </div>
          <div className="flex gap-6">
            <div className="text-right">
              <p className="text-[10px] uppercase text-slate-400 font-bold tracking-widest">Pendentes (Total)</p>
              <p className="text-xl font-black text-amber-400">{stats.pending}</p>
            </div>
            <div className="text-right border-l border-slate-700 pl-6">
              <p className="text-[10px] uppercase text-slate-400 font-bold tracking-widest">Vendas Aprovadas</p>
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

          <button onClick={() => { fetchPayments(); fetchGlobalStats(); }} className="text-slate-400 hover:text-indigo-600 transition">
            <span className="text-xs font-bold uppercase tracking-tighter">🔄 Atualizar Dados</span>
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
                        <div className="flex justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                          <button
                            onClick={() => handleAction(p.id, 'APPROVED')}
                            className="bg-emerald-500 text-white px-4 py-2 rounded-xl text-xs font-bold hover:bg-emerald-600 shadow-lg shadow-emerald-200 transition-all active:scale-95"
                          >
                            APROVAR
                          </button>
                          <button
                            onClick={() => handleAction(p.id, 'REJECTED')}
                            className="bg-white border border-slate-200 text-slate-400 px-4 py-2 rounded-xl text-xs font-bold hover:bg-red-50 hover:text-red-500 hover:border-red-100 transition-all"
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