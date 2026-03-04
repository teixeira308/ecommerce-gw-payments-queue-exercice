import { useState, useEffect, useCallback, useRef } from 'react'

const GATEWAY_URL = 'http://localhost:9000'
const ECOMMERCE_URL = 'http://localhost:3000'

// Componente de Modal de Notificação
const Modal = ({ payment, onClose }) => {
  if (!payment) return null;
  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4 animate-in fade-in duration-300">
      <div className="bg-white rounded-[2.5rem] p-8 max-w-md w-full shadow-2xl border border-slate-100 animate-in zoom-in-95 duration-300">
        <div className="flex flex-col items-center text-center">
          <div className="w-16 h-16 bg-amber-100 text-amber-600 rounded-2xl flex items-center justify-center text-2xl mb-4">🔔</div>
          <h2 className="text-2xl font-black text-slate-900 mb-2">Novo Pagamento!</h2>
          <p className="text-slate-500 mb-6">Um novo pagamento de <span className="font-bold text-slate-900">R$ {Number(payment.amount || 0).toFixed(2)}</span> chegou e aguarda processamento.</p>
          <button onClick={onClose} className="w-full bg-slate-900 text-white py-4 rounded-2xl font-black text-sm hover:bg-slate-800 transition-all active:scale-95 shadow-xl shadow-slate-200">ENTENDIDO</button>
        </div>
      </div>
    </div>
  );
};

// Componente de Formulário de Item
const ItemFormModal = ({ item, onClose, onSave }) => {
  const [formData, setFormData] = useState(item ? { name: item.Name, price: item.Price } : { name: '', price: '' });
  const [isSaving, setIsSaving] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    setIsSaving(true);
    await onSave({ ...formData, price: Number(formData.price || 0) });
    setIsSaving(false);
    onClose();
  };

  return (
    <div className="fixed inset-0 bg-black/60 backdrop-blur-sm z-[100] flex items-center justify-center p-4">
      <div className="bg-white rounded-[2.5rem] p-8 max-w-md w-full shadow-2xl border border-slate-100 animate-in zoom-in-95 duration-300">
        <h2 className="text-2xl font-black text-slate-900 mb-6">{item ? 'Editar Item' : 'Novo Item'}</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="text-[10px] font-black text-slate-400 uppercase mb-1 block">Nome do Produto</label>
            <input required placeholder="Ex: Produto Top" className="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-indigo-500/20" value={formData.name} onChange={e => setFormData({...formData, name: e.target.value})} />
          </div>
          <div>
            <label className="text-[10px] font-black text-slate-400 uppercase mb-1 block">Preço (R$)</label>
            <input type="number" step="0.01" required placeholder="0.00" className="w-full bg-slate-50 border border-slate-200 rounded-xl px-4 py-3 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-indigo-500/20" value={formData.price} onChange={e => setFormData({...formData, price: e.target.value})} />
          </div>
          <div className="flex gap-3 pt-4">
            <button type="button" onClick={onClose} className="flex-1 px-6 py-4 rounded-2xl font-black text-xs text-slate-400 hover:bg-slate-50 transition-all">CANCELAR</button>
            <button type="submit" disabled={isSaving} className="flex-1 bg-indigo-600 text-white px-6 py-4 rounded-2xl font-black text-xs hover:bg-indigo-700 transition-all active:scale-95 disabled:opacity-50">{isSaving ? 'SALVANDO...' : 'SALVAR ITEM'}</button>
          </div>
        </form>
      </div>
    </div>
  );
};

function App() {
  const [activeTab, setActiveTab] = useState('transactions')
  const [payments, setPayments] = useState([])
  const [items, setItems] = useState([])
  const [page, setPage] = useState(1)
  const [limit] = useState(10)
  const [loading, setLoading] = useState(false)
  const [stats, setStats] = useState({ pending: 0, totalAmount: 0 })
  const [filterPending, setFilterPending] = useState(false)
  
  const [autoRefresh, setAutoRefresh] = useState(true)
  const [processing, setProcessing] = useState({})
  const [newPaymentModal, setNewPaymentModal] = useState(null)
  const [itemFormModal, setItemFormModal] = useState({ open: false, item: null })
  const [seenIds, setSeenIds] = useState(new Set())
  const isFirstLoad = useRef(true)

  const fetchGlobalStats = useCallback(async () => {
    try {
      const res = await fetch(`${GATEWAY_URL}/payments?limit=9999`)
      const data = await res.json()
      const all = data || []
      const pendingCount = all.filter(p => p.status === 'PENDING').length
      const approvedTotal = all.filter(p => p.status === 'APPROVED').reduce((acc, curr) => acc + Number(curr.amount || 0), 0)
      setStats({ pending: pendingCount, totalAmount: approvedTotal })
    } catch (err) { console.error(err) }
  }, [])

  const fetchPayments = useCallback(async () => {
    setLoading(true)
    try {
      const res = await fetch(`${GATEWAY_URL}/payments?limit=${limit}&page=${page}`)
      const data = await res.json()
      const paymentsData = data || []
      
      setSeenIds(prev => {
        if (!isFirstLoad.current && activeTab === 'transactions') {
          const newOnes = paymentsData.filter(p => p.status === 'PENDING' && !prev.has(p.id))
          if (newOnes.length > 0) setNewPaymentModal(newOnes[0])
        }
        const next = new Set(prev)
        paymentsData.forEach(p => next.add(p.id))
        return next
      })
      
      setPayments(paymentsData)
    } catch (err) { console.error(err) } finally { setTimeout(() => setLoading(false), 400) }
  }, [page, limit, activeTab])

  const fetchItems = useCallback(async () => {
    setLoading(true)
    try {
      const res = await fetch(`${ECOMMERCE_URL}/items`)
      const data = await res.json()
      setItems(data || [])
    } catch (err) { console.error(err) } finally { setTimeout(() => setLoading(false), 400) }
  }, [])

  useEffect(() => {
    if (activeTab === 'transactions') fetchPayments()
    if (activeTab === 'catalog') fetchItems()
    fetchGlobalStats()
    isFirstLoad.current = false
  }, [activeTab, fetchPayments, fetchItems, fetchGlobalStats])

  useEffect(() => {
    if (!autoRefresh) return;
    const interval = setInterval(() => {
      if (activeTab === 'transactions') fetchPayments()
      if (activeTab === 'catalog') fetchItems()
      fetchGlobalStats()
    }, 30000);
    return () => clearInterval(interval);
  }, [autoRefresh, activeTab, fetchPayments, fetchItems, fetchGlobalStats]);

  const handleAction = async (id, status) => {
    setProcessing(prev => ({ ...prev, [id]: true }))
    try {
      await fetch(`${GATEWAY_URL}/payments/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status })
      })
      fetchPayments(); fetchGlobalStats();
    } catch (err) { alert("Erro na ação"); setProcessing(prev => ({ ...prev, [id]: false })) }
  }

  const handleProcess = async (id) => {
    setProcessing(prev => ({ ...prev, [id]: true }))
    try {
      const res = await fetch(`${GATEWAY_URL}/payments/${id}/process`, { method: 'POST' })
      if (!res.ok) throw new Error("Erro no gateway")
      fetchPayments(); fetchGlobalStats();
    } catch (err) { alert(err.message); setProcessing(prev => ({ ...prev, [id]: false })) }
  }

  const saveItem = async (formData) => {
    try {
      const isEdit = !!itemFormModal.item;
      const url = isEdit ? `${ECOMMERCE_URL}/items/${itemFormModal.item.ID}` : `${ECOMMERCE_URL}/items`;
      const res = await fetch(url, {
        method: isEdit ? 'PUT' : 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(formData)
      });
      if (!res.ok) throw new Error("Erro ao salvar");
      fetchItems();
    } catch (err) { alert(err.message) }
  }

  const deleteItem = async (id) => {
    if (!confirm("Excluir item?")) return;
    try {
      await fetch(`${ECOMMERCE_URL}/items/${id}`, { method: 'DELETE' });
      fetchItems();
    } catch (err) { alert("Erro ao excluir") }
  }

  const formatDate = (dateString) => {
    if (!dateString) return "-";
    return new Date(dateString).toLocaleString('pt-BR', { day: '2-digit', month: '2-digit', year: '2-digit', hour: '2-digit', minute: '2-digit' });
  };

  return (
    <div className="min-h-screen bg-[#f1f5f9]">
      <Modal payment={newPaymentModal} onClose={() => setNewPaymentModal(null)} />
      {itemFormModal.open && (
        <ItemFormModal item={itemFormModal.item} onClose={() => setItemFormModal({ open: false, item: null })} onSave={saveItem} />
      )}
      
      <nav className="bg-slate-900 text-white p-6 shadow-lg sticky top-0 z-50">
        <div className="max-w-7xl mx-auto flex justify-between items-center">
          <div className="flex items-center gap-8">
            <div className="flex items-center gap-3">
              <div className="w-10 h-10 bg-indigo-500 rounded-xl flex items-center justify-center font-black text-xs">HUB</div>
              <h1 className="text-xl font-bold tracking-tight">Admin Hub</h1>
            </div>
            <div className="flex bg-slate-800 p-1 rounded-xl">
              <button onClick={() => setActiveTab('transactions')} className={`px-4 py-2 rounded-lg text-xs font-bold transition-all ${activeTab === 'transactions' ? 'bg-slate-700 text-white shadow-lg' : 'text-slate-400'}`}>Transações</button>
              <button onClick={() => setActiveTab('catalog')} className={`px-4 py-2 rounded-lg text-xs font-bold transition-all ${activeTab === 'catalog' ? 'bg-slate-700 text-white shadow-lg' : 'text-slate-400'}`}>Catálogo</button>
            </div>
          </div>
          <div className="flex gap-6 items-center">
            <div className="flex flex-col items-end mr-2">
              <span className="text-[8px] uppercase text-slate-500 font-black mb-1">AUTO-UPDATE</span>
              <button onClick={() => setAutoRefresh(!autoRefresh)} className={`w-8 h-4 rounded-full relative transition-colors ${autoRefresh ? 'bg-emerald-500' : 'bg-slate-700'}`}>
                <div className={`absolute top-0.5 w-3 h-3 bg-white rounded-full transition-all ${autoRefresh ? 'left-4.5' : 'left-0.5'}`} />
              </button>
            </div>
            <div className="text-right">
              <p className="text-[10px] uppercase text-slate-400 font-bold">Pendentes</p>
              <p className="text-xl font-black text-amber-400">{stats.pending}</p>
            </div>
            <div className="text-right border-l border-slate-700 pl-6">
              <p className="text-[10px] uppercase text-slate-400 font-bold">Total Vendas</p>
              <p className="text-xl font-black text-emerald-400">R$ {stats.totalAmount.toLocaleString('pt-BR', { minimumFractionDigits: 2 })}</p>
            </div>
          </div>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto p-8">
        <div className="flex justify-between items-center mb-6">
          {activeTab === 'transactions' ? (
            <div className="flex bg-white p-1 rounded-2xl border border-slate-200 shadow-sm">
              <button onClick={() => setFilterPending(false)} className={`px-6 py-2 rounded-xl text-xs font-bold transition-all ${!filterPending ? 'bg-slate-900 text-white shadow-md' : 'text-slate-500'}`}>Todos</button>
              <button onClick={() => setFilterPending(true)} className={`px-6 py-2 rounded-xl text-xs font-bold transition-all ${filterPending ? 'bg-amber-500 text-white shadow-md' : 'text-slate-500'}`}>Pendentes</button>
            </div>
          ) : (
            <h2 className="text-xl font-black text-slate-900">Catálogo de Produtos</h2>
          )}
          
          <div className="flex items-center gap-4">
            {activeTab === 'catalog' && (
              <button onClick={() => setItemFormModal({ open: true, item: null })} className="bg-indigo-600 text-white px-6 py-2 rounded-xl text-[10px] font-black hover:bg-indigo-700 shadow-lg shadow-indigo-200 transition-all active:scale-95">+ NOVO ITEM</button>
            )}
            <button 
              onClick={() => { if(activeTab === 'transactions') fetchPayments(); else fetchItems(); fetchGlobalStats(); }} 
              disabled={loading} 
              className="flex items-center justify-center gap-3 px-4 py-2 min-w-[140px] text-slate-400 hover:text-indigo-600 transition-all active:scale-95 disabled:opacity-50"
            >
              <div className="relative w-4 h-4 flex items-center justify-center">
                <span className={`absolute transition-all duration-300 ${loading ? 'opacity-0 scale-50 rotate-180' : 'opacity-100 scale-100 rotate-0'}`}>🔄</span>
                <div className={`absolute w-4 h-4 border-2 border-indigo-500 border-t-transparent rounded-full animate-spin transition-all duration-300 ${loading ? 'opacity-100 scale-100' : 'opacity-0 scale-50'}`}></div>
              </div>
              <span className="text-[10px] font-black uppercase tracking-tighter">Atualizar</span>
            </button>
          </div>
        </div>

        {activeTab === 'transactions' ? (
          <div className="bg-white rounded-[2.5rem] shadow-xl border border-slate-200 overflow-hidden">
            <table className="w-full text-left">
              <thead>
                <tr className="bg-slate-50/50 border-b border-slate-100">
                  <th className="p-6 text-xs font-black text-slate-400 uppercase">Identificação</th>
                  <th className="p-6 text-xs font-black text-slate-400 uppercase">Valor</th>
                  <th className="p-6 text-xs font-black text-slate-400 uppercase">Status</th>
                  <th className="p-6 text-xs font-black text-slate-400 uppercase text-right">Ações</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-slate-50">
                {payments.filter(p => !filterPending || p.status === 'PENDING').map(p => (
                  <tr key={p.id} className="hover:bg-slate-50/80 group">
                    <td className="p-6">
                      <div className="flex flex-col">
                        <span className="text-[10px] font-mono text-slate-400 bg-slate-100 px-2 py-0.5 rounded w-fit uppercase">ID: {(p.id || '').substring(0, 8)}</span>
                        <span className="font-bold text-slate-700">Pedido #{(p.order_id || '').substring(0, 5)}</span>
                        <span className="text-[11px] text-slate-400 mt-1">{formatDate(p.created_at)}</span>
                      </div>
                    </td>
                    <td className="p-6">
                      <span className="text-lg font-black text-slate-900 block">R$ {Number(p.amount || 0).toFixed(2)}</span>
                      <span className="text-[10px] font-bold text-indigo-500 uppercase">{p.method}</span>
                    </td>
                    <td className="p-6">
                      <span className={`px-4 py-1.5 rounded-full text-[10px] font-black uppercase border ${p.status === 'APPROVED' ? 'bg-emerald-50 text-emerald-600 border-emerald-100' : p.status === 'REJECTED' ? 'bg-red-50 text-red-600 border-red-100' : 'bg-amber-50 text-amber-600 animate-pulse'}`}>{p.status}</span>
                    </td>
                    <td className="p-6 text-right">
                      {p.status === 'PENDING' ? (
                        <div className="flex justify-end gap-2 opacity-0 group-hover:opacity-100 transition-opacity">
                          <button onClick={() => handleProcess(p.id)} disabled={processing[p.id]} className="bg-indigo-500 text-white px-4 py-2 rounded-xl text-xs font-bold hover:bg-indigo-600">ENVIAR</button>
                          <button onClick={() => handleAction(p.id, 'REJECTED')} disabled={processing[p.id]} className="bg-white border border-slate-200 text-slate-400 px-4 py-2 rounded-xl text-xs font-bold hover:text-red-500">RECUSAR</button>
                        </div>
                      ) : <span className="text-[10px] font-bold text-slate-300 uppercase italic">Concluído</span>}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {items.map(item => (
              <div key={item.ID} className="bg-white rounded-[2rem] p-6 shadow-xl shadow-slate-200/50 border border-slate-100 hover:border-indigo-200 transition-all group">
                <div className="flex justify-between items-start mb-4">
                  <div>
                    <span className="text-[10px] font-mono text-slate-400 uppercase block mb-1">ID: {(item.ID || '').substring(0, 8)}</span>
                    <h3 className="text-lg font-black text-slate-900 leading-tight">{item.Name}</h3>
                  </div>
                  <span className="bg-indigo-50 text-indigo-600 px-3 py-1 rounded-full text-xs font-black">R$ {Number(item.Price || 0).toFixed(2)}</span>
                </div>
                <div className="flex items-center justify-between pt-6 border-t border-slate-50">
                  <div className="flex flex-col">
                    <span className="text-[9px] font-black text-slate-400 uppercase">Cadastro</span>
                    <span className="text-sm font-bold text-slate-700">{formatDate(item.CreatedAt)}</span>
                  </div>
                  <div className="flex gap-2">
                    <button onClick={() => setItemFormModal({ open: true, item })} className="p-3 bg-slate-50 text-slate-400 rounded-xl hover:bg-indigo-50 hover:text-indigo-600 transition-all">✏️</button>
                    <button onClick={() => deleteItem(item.ID)} className="p-3 bg-slate-50 text-slate-400 rounded-xl hover:bg-red-50 hover:text-red-500 transition-all">🗑️</button>
                  </div>
                </div>
              </div>
            ))}
            {items.length === 0 && <div className="col-span-full p-20 text-center bg-white rounded-[2.5rem] border border-dashed border-slate-300 text-slate-400 font-medium">Catálogo vazio.</div>}
          </div>
        )}
      </main>
    </div>
  )
}

export default App;