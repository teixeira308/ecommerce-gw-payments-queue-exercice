import { useState, useEffect } from 'react'

const API_URL = 'http://localhost:3000'

function App() {
  const [activeTab, setActiveTab] = useState('shop')
  const [items, setItems] = useState([])
  const [orders, setOrders] = useState([])
  const [cart, setCart] = useState([])
  const [loading, setLoading] = useState(false)
  const [paymentMethod, setPaymentMethod] = useState('CREDIT_CARD')

  const fetchData = async () => {
    try {
      const [itemsRes, ordersRes] = await Promise.all([
        fetch(`${API_URL}/items`),
        fetch(`${API_URL}/orders`)
      ])
      const itemsData = await itemsRes.json()
      const ordersData = await ordersRes.json()
      setItems(itemsData || [])
      setOrders(ordersData || [])
    } catch (err) {
      console.error("Erro ao sincronizar:", err)
    }
  }

  useEffect(() => {
    fetchData()
    const interval = setInterval(fetchData, 5000)
    return () => clearInterval(interval)
  }, [])

  const addToCart = (product) => {
    setCart(prev => {
      const exists = prev.find(i => i.item_id === product.id)
      if (exists) return prev.map(i => i.item_id === product.id ? { ...i, quantity: i.quantity + 1 } : i)
      return [...prev, { item_id: product.id, name: product.name, price: product.price, quantity: 1 }]
    })
  }

  const removeFromCart = (productId) => {
    setCart(prev => prev.filter(item => item.item_id !== productId))
  }

  const checkout = async () => {
    if (cart.length === 0) return
    setLoading(true)
    try {
      const res = await fetch(`${API_URL}/orders`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          items: cart.map(({ item_id, quantity }) => ({ item_id, quantity })),
          method: paymentMethod 
        })
      })
      if (res.ok) {
        setCart([])
        setActiveTab('orders')
        fetchData()
      }
    } catch (err) {
      alert("Erro na conexão.")
    } finally {
      setLoading(false)
    }
  }

  const totalCart = cart.reduce((acc, curr) => acc + (curr.price * curr.quantity), 0)

  const paymentMethods = [
    { id: 'CREDIT_CARD', label: 'Cartão de Crédito', icon: '💳' },
    { id: 'PIX', label: 'Pix', icon: '📱' },
    { id: 'BOLETO', label: 'Boleto', icon: '📄' },
  ]

  return (
    <div className="min-h-screen bg-slate-950 text-slate-200 selection:bg-indigo-500/30 font-sans">
      {/* Background Decorativo - Mesh Gradients */}
      <div className="fixed inset-0 overflow-hidden pointer-events-none">
        <div className="absolute top-[-10%] right-[-10%] w-[50%] h-[50%] bg-indigo-900/20 rounded-full blur-[120px]"></div>
        <div className="absolute bottom-[-10%] left-[-10%] w-[50%] h-[50%] bg-blue-900/20 rounded-full blur-[120px]"></div>
        <div className="absolute top-[30%] left-[20%] w-[30%] h-[30%] bg-purple-900/10 rounded-full blur-[100px]"></div>
      </div>

      {/* Navbar Superior */}
      <nav className="bg-slate-950/50 backdrop-blur-2xl border-b border-white/5 sticky top-0 z-50">
        <div className="max-w-6xl mx-auto px-6 h-20 flex items-center justify-between">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-2xl flex items-center justify-center text-white font-black shadow-xl shadow-indigo-500/20">G</div>
            <h1 className="text-xl font-black tracking-tight text-white">GO.SHOP</h1>
          </div>
          <div className="flex bg-white/5 p-1 rounded-2xl border border-white/10">
            <button onClick={() => setActiveTab('shop')} className={`px-6 py-2 rounded-xl text-xs font-black transition-all ${activeTab === 'shop' ? 'bg-indigo-600 text-white shadow-lg' : 'text-slate-400 hover:text-slate-200'}`}>LOJA</button>
            <button onClick={() => setActiveTab('orders')} className={`px-6 py-2 rounded-xl text-xs font-black transition-all ${activeTab === 'orders' ? 'bg-indigo-600 text-white shadow-lg' : 'text-slate-400 hover:text-slate-200'}`}>PEDIDOS</button>
          </div>
        </div>
      </nav>

      <main className="max-w-6xl mx-auto p-6 md:p-12 relative z-10">
        {activeTab === 'shop' ? (
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 items-start">
            
            {/* Seção de Produtos */}
            <div className="lg:col-span-8">
              <header className="mb-12">
                <span className="text-[10px] font-black text-indigo-400 uppercase tracking-[0.3em] mb-3 block">Digital Storefront</span>
                <h2 className="text-5xl font-black text-white tracking-tight">Tech Essentials</h2>
                <p className="text-slate-500 mt-4 text-lg">Curadoria exclusiva dos melhores produtos para seu setup.</p>
              </header>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-8">
                {items.map(item => {
                  const price = item.price ?? item.Price ?? 0;
                  const name = item.name ?? item.Name ?? "Produto";
                  const id = item.id ?? item.ID;
                  return (
                    <div key={id} className="group bg-white/5 backdrop-blur-sm p-3 rounded-[3rem] border border-white/10 hover:border-indigo-500/50 transition-all duration-500 hover:-translate-y-2">
                      <div className="bg-slate-900/50 rounded-[2.5rem] p-10 flex items-center justify-center relative overflow-hidden">
                        <span className="text-7xl group-hover:scale-110 transition-transform duration-700 z-10 drop-shadow-2xl">📦</span>
                        <div className="absolute inset-0 bg-gradient-to-br from-indigo-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                      </div>
                      <div className="p-8">
                        <h3 className="font-bold text-xl text-white mb-2">{name}</h3>
                        <div className="flex items-center gap-2 mb-6">
                           <span className="w-2 h-2 bg-emerald-500 rounded-full animate-pulse"></span>
                           <span className="text-[10px] text-slate-400 uppercase font-black tracking-widest">Em Estoque</span>
                        </div>
                        <div className="flex items-center justify-between">
                          <span className="text-2xl font-black text-white">R$ {Number(price).toFixed(2)}</span>
                          <button onClick={() => addToCart({ ...item, id, name, price })} className="w-14 h-14 bg-indigo-600 text-white rounded-3xl flex items-center justify-center hover:bg-indigo-500 transition-all active:scale-90 shadow-2xl shadow-indigo-500/20 text-xl font-bold">＋</button>
                        </div>
                      </div>
                    </div>
                  );
                })}
              </div>
            </div>

            {/* Carrinho Lateral - Glassmorphism */}
            <div className="lg:col-span-4 sticky top-28">
              <div className="bg-white/[0.03] backdrop-blur-3xl p-8 rounded-[3rem] shadow-2xl border border-white/10">
                <div className="flex items-center justify-between mb-10">
                  <h2 className="text-2xl font-black text-white">Sacola</h2>
                  <span className="bg-indigo-600/20 text-indigo-400 text-xs font-black px-3 py-1 rounded-full border border-indigo-500/30">{cart.length} itens</span>
                </div>
                
                {cart.length === 0 ? (
                  <div className="py-20 text-center">
                    <div className="text-5xl mb-6 grayscale opacity-20">🛒</div>
                    <p className="text-[10px] font-black text-slate-500 uppercase tracking-widest leading-relaxed">Adicione produtos para <br /> começar sua jornada.</p>
                  </div>
                ) : (
                  <>
                    <div className="space-y-6 mb-10 max-h-[350px] overflow-y-auto pr-2 custom-scrollbar">
                      {cart.map(i => (
                        <div key={i.item_id} className="flex justify-between items-center group animate-in fade-in slide-in-from-right-4">
                          <div className="flex flex-col">
                            <span className="text-sm font-bold text-white mb-1">{i.quantity}x {i.name}</span>
                            <span className="text-[9px] text-slate-500 font-mono tracking-widest uppercase">REF: {i.item_id.substring(0, 8)}</span>
                          </div>
                          <div className="flex items-center gap-4">
                            <span className="text-sm font-black text-indigo-400">R$ {(i.price * i.quantity).toFixed(2)}</span>
                            <button onClick={() => removeFromCart(i.item_id)} className="w-8 h-8 rounded-full flex items-center justify-center bg-white/5 text-slate-500 hover:bg-red-500/20 hover:text-red-400 transition-all">✕</button>
                          </div>
                        </div>
                      ))}
                    </div>

                    <div className="bg-white/5 rounded-[2rem] p-6 mb-10 border border-white/5">
                      <div className="flex justify-between items-center mb-2">
                        <span className="text-[10px] font-black text-slate-500 uppercase tracking-widest">Subtotal</span>
                        <span className="text-sm font-bold text-slate-300">R$ {totalCart.toFixed(2)}</span>
                      </div>
                      <div className="flex justify-between items-center">
                        <span className="text-[10px] font-black text-slate-500 uppercase tracking-widest font-black">Total</span>
                        <span className="text-3xl font-black text-white tracking-tighter">R$ {totalCart.toFixed(2)}</span>
                      </div>
                    </div>
                    
                    <div className="mb-10">
                      <p className="text-[10px] font-black text-slate-500 uppercase mb-4 tracking-[0.2em]">Pagamento</p>
                      <div className="grid grid-cols-1 gap-3">
                        {paymentMethods.map(m => (
                          <button key={m.id} onClick={() => setPaymentMethod(m.id)} className={`flex items-center justify-between p-5 rounded-2xl border transition-all text-sm font-bold ${paymentMethod === m.id ? 'bg-indigo-600 border-indigo-400 text-white shadow-lg shadow-indigo-500/20' : 'bg-white/5 border-white/5 text-slate-400 hover:bg-white/10'}`}>
                            <div className="flex items-center gap-4">
                              <span className="text-xl">{m.icon}</span>
                              {m.label}
                            </div>
                            {paymentMethod === m.id && <div className="w-2 h-2 bg-white rounded-full shadow-[0_0_8px_white]"></div>}
                          </button>
                        ))}
                      </div>
                    </div>

                    <button onClick={checkout} disabled={loading} className={`w-full py-6 rounded-3xl font-black text-xs uppercase tracking-[0.3em] shadow-2xl transition-all active:scale-95 border ${loading ? 'bg-slate-800 border-white/5 text-slate-600' : 'bg-white border-white text-slate-900 hover:bg-indigo-50 hover:shadow-indigo-500/20'}`}>{loading ? 'Sincronizando...' : 'Fechar Pedido'}</button>
                  </>
                )}
              </div>
            </div>
          </div>
        ) : (
          /* Aba de Pedidos - Dark Style */
          <div className="max-w-4xl mx-auto">
            <header className="mb-20 text-center">
              <span className="text-[10px] font-black text-indigo-400 uppercase tracking-[0.3em] mb-3 block">Orders Dashboard</span>
              <h2 className="text-5xl font-black text-white tracking-tight">Meus Pedidos</h2>
              <div className="w-16 h-1.5 bg-gradient-to-r from-indigo-600 to-purple-600 mt-6 mx-auto rounded-full"></div>
            </header>

            <div className="space-y-10">
              {orders.slice().reverse().map((order) => (
                <div key={order.id} className="bg-white/[0.02] backdrop-blur-md rounded-[3.5rem] border border-white/5 shadow-2xl overflow-hidden group hover:bg-white/[0.04] hover:border-white/10 transition-all duration-700">
                  <div className="p-12 flex flex-col md:flex-row justify-between gap-10">
                    <div className="flex-1">
                      <div className="flex items-center gap-5 mb-8">
                        <span className={`px-6 py-2 rounded-full text-[10px] font-black uppercase tracking-widest border ${order.status === 'paid' || order.status === 'APPROVED' ? 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20' : order.status === 'pending' ? 'bg-amber-500/10 text-amber-400 border-amber-500/20' : 'bg-white/5 text-slate-400 border-white/10'}`}>{order.status}</span>
                        <span className="text-xs text-slate-500 font-bold uppercase tracking-tighter">{new Date(order.created_at).toLocaleString('pt-BR')}</span>
                      </div>
                      <h4 className="text-[10px] font-mono text-slate-600 mb-10 uppercase tracking-[0.3em]">REF: {order.id}</h4>
                      <div className="space-y-4">
                        {order.items.map((item, idx) => (
                          <div key={idx} className="flex items-center justify-between bg-white/[0.03] px-6 py-4 rounded-[1.5rem] border border-white/5 group-hover:bg-white/[0.01] transition-all">
                            <span className="text-sm text-slate-300 font-bold">{item.quantity}x <span className="font-medium text-slate-500 ml-3">Produto SKU-{item.item_id.substring(0, 6)}</span></span>
                            <span className="text-sm font-black text-white">R$ {Number(item.subtotal).toFixed(2)}</span>
                          </div>
                        ))}
                      </div>
                    </div>
                    <div className="md:w-64 bg-white/[0.03] rounded-[2.5rem] p-10 flex flex-col justify-center items-center text-center border border-white/5 group-hover:border-indigo-500/20 transition-all">
                      <span className="text-[10px] text-slate-500 font-black uppercase mb-3 tracking-[0.2em]">Investimento</span>
                      <span className="text-4xl font-black text-white tracking-tighter">R$ {Number(order.total).toFixed(2)}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
            {orders.length === 0 && <div className="p-32 text-center bg-white/[0.02] rounded-[4rem] border border-dashed border-white/10 text-slate-600 font-black uppercase text-xs tracking-[0.4em]">Nenhuma atividade registrada.</div>}
          </div>
        )}
      </main>
    </div>
  )
}

export default App;