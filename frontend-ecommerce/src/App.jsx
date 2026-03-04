import { useState, useEffect } from 'react'

const API_URL = 'http://localhost:3000'

function App() {
  const [activeTab, setActiveTab] = useState('shop')
  const [items, setItems] = useState([])
  const [orders, setOrders] = useState([])
  const [cart, setCart] = useState([])
  const [loading, setLoading] = useState(false)

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
        body: JSON.stringify({ items: cart.map(({ item_id, quantity }) => ({ item_id, quantity })) })
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

  return (
    <div className="min-h-screen bg-[#f8fafc] font-sans text-slate-900">
      {/* Navbar Superior */}
      <nav className="bg-white/80 backdrop-blur-md border-b border-slate-200 sticky top-0 z-50">
        <div className="max-w-6xl mx-auto px-6 h-20 flex items-center justify-between">
          <div className="flex items-center gap-2">
            <div className="w-8 h-8 bg-indigo-600 rounded-lg flex items-center justify-center text-white font-bold">G</div>
            <h1 className="text-xl font-black tracking-tighter text-slate-800">GO.SHOP</h1>
          </div>
          <div className="flex bg-slate-100 p-1 rounded-xl">
            <button
              onClick={() => setActiveTab('shop')}
              className={`px-6 py-2 rounded-lg text-sm font-bold transition-all ${activeTab === 'shop' ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
            >
              Catálogo
            </button>
            <button
              onClick={() => setActiveTab('orders')}
              className={`px-6 py-2 rounded-lg text-sm font-bold transition-all ${activeTab === 'orders' ? 'bg-white text-indigo-600 shadow-sm' : 'text-slate-500 hover:text-slate-700'}`}
            >
              Pedidos
            </button>
          </div>
        </div>
      </nav>

      <main className="max-w-6xl mx-auto p-6 md:p-12">
        {activeTab === 'shop' ? (
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 items-start">

            {/* Lista de Itens Centralizada */}
            <div className="lg:col-span-8">
              <header className="mb-10 text-center md:text-left">
                <h2 className="text-3xl font-extrabold text-slate-900 tracking-tight">Produtos Selecionados</h2>
                <p className="text-slate-500 mt-2">Qualidade e tecnologia direto para sua casa.</p>
              </header>

              <div className="grid grid-cols-1 sm:grid-cols-2 gap-6">
                {items.map(item => {
                  const price = item.price ?? item.Price ?? 0;
                  const name = item.name ?? item.Name ?? "Produto";
                  const id = item.id ?? item.ID;
                  return (
                    <div key={id} className="group bg-white p-6 rounded-[2rem] border border-slate-100 shadow-xl shadow-slate-200/40 hover:shadow-2xl hover:shadow-indigo-200/40 transition-all duration-300">
                      <div className="w-full h-40 bg-slate-50 rounded-2xl mb-6 flex items-center justify-center group-hover:bg-indigo-50 transition-colors">
                        <span className="text-4xl">📦</span>
                      </div>
                      <h3 className="font-bold text-xl text-slate-800 mb-1">{name}</h3>
                      <p className="text-indigo-600 font-black text-2xl mb-6">R$ {Number(price).toFixed(2)}</p>
                      <button
                        onClick={() => addToCart({ ...item, id, name, price })}
                        className="w-full bg-slate-900 text-white py-4 rounded-2xl font-bold hover:bg-indigo-600 active:scale-95 transition-all"
                      >
                        Adicionar ao Carrinho
                      </button>
                    </div>
                  );
                })}
              </div>
            </div>

            {/* Carrinho Lateral */}
            <div className="lg:col-span-4 sticky top-28">
              <div className="bg-white p-8 rounded-[2rem] shadow-2xl shadow-slate-200 border border-slate-100">
                <h2 className="text-2xl font-bold mb-6 flex items-center gap-2 text-slate-800">
                  Resumo
                  {cart.length > 0 && <span className="text-sm bg-indigo-100 text-indigo-600 px-2 py-0.5 rounded-md">{cart.length}</span>}
                </h2>
                {cart.length === 0 ? (
                  <div className="py-12 text-center">
                    <div className="text-4xl mb-4">🛒</div>
                    <p className="text-slate-400 font-medium">Seu carrinho está <br /> pronto para ser cheio.</p>
                  </div>
                ) : (
                  <>
                    <div className="space-y-4 mb-8">
                      {cart.map(i => (
                        <div key={i.item_id} className="group flex justify-between items-center bg-slate-50 p-3 rounded-xl border border-slate-100 transition-all hover:border-red-100">
                          <div className="flex flex-col">
                            <span className="text-sm font-semibold text-slate-700">{i.quantity}x {i.name}</span>
                            <span className="text-[10px] text-slate-400 font-mono">{i.item_id.substring(0, 8)}</span>
                          </div>

                          <div className="flex items-center gap-3">
                            <span className="text-sm font-bold text-slate-900">R$ {(i.price * i.quantity).toFixed(2)}</span>

                            {/* Botão de Remover */}
                            <button
                              onClick={() => removeFromCart(i.item_id)}
                              className="text-slate-300 hover:text-red-500 transition-colors p-1"
                              title="Remover item"
                            >
                              <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                              </svg>
                            </button>
                          </div>
                        </div>
                      ))}
                    </div>
                    <div className="border-t border-slate-100 pt-6 mb-8 flex justify-between items-end">
                      <span className="text-slate-500 font-medium">Total:</span>
                      <span className="text-2xl font-black text-indigo-600 tracking-tighter">R$ {totalCart.toFixed(2)}</span>
                    </div>
                    <button
                      onClick={checkout}
                      disabled={loading}
                      className={`w-full py-5 rounded-2xl font-bold text-white shadow-lg shadow-indigo-200 transition-all active:scale-95 ${loading ? 'bg-slate-300' : 'bg-indigo-600 hover:bg-indigo-700'}`}
                    >
                      {loading ? 'Processando...' : 'Finalizar Compra'}
                    </button>
                  </>
                )}
              </div>
            </div>
          </div>
        ) : (
          /* --- ABA DE PEDIDOS (DASHBOARD) --- */
          <div className="max-w-4xl mx-auto">
            <header className="mb-12 text-center">
              <h2 className="text-4xl font-black text-slate-900 tracking-tight">Seus Pedidos</h2>
              <p className="text-slate-500 mt-2">Acompanhe o processamento em tempo real.</p>
            </header>

            <div className="space-y-8">
              {orders.slice().reverse().map((order) => (
                <div key={order.id} className="bg-white rounded-[2.5rem] border border-slate-100 shadow-xl shadow-slate-200/50 overflow-hidden">
                  <div className="p-8 flex flex-col md:flex-row justify-between gap-6">
                    <div className="flex-1">
                      <div className="flex items-center gap-3 mb-4">
                        <span className={`px-4 py-1 rounded-full text-[10px] font-black uppercase tracking-widest ${order.status === 'paid' || order.status === 'APPROVED' ? 'bg-emerald-100 text-emerald-600' :
                            order.status === 'pending' ? 'bg-amber-100 text-amber-600' : 'bg-slate-100 text-slate-500'
                          }`}>
                          {order.status}
                        </span>
                        <span className="text-xs text-slate-400 font-medium">{new Date(order.created_at).toLocaleString('pt-BR')}</span>
                      </div>
                      <h4 className="text-xs font-mono text-slate-400 mb-6 truncate max-w-xs">ID: {order.id}</h4>

                      <div className="space-y-3">
                        {order.items.map((item, idx) => (
                          <div key={idx} className="flex items-center justify-between bg-slate-50 px-4 py-2 rounded-lg border border-slate-100">
                            <span className="text-sm text-slate-600 font-bold">{item.quantity}x <span className="font-normal text-slate-400">Item {item.item_id.substring(0, 6)}</span></span>
                            <span className="text-sm font-black text-slate-800">R$ {Number(item.subtotal).toFixed(2)}</span>
                          </div>
                        ))}
                      </div>
                    </div>

                    <div className="md:w-48 bg-slate-50 rounded-3xl p-6 flex flex-col justify-center items-center text-center border border-slate-100">
                      <span className="text-xs text-slate-400 font-bold uppercase mb-1">Total Pago</span>
                      <span className="text-2xl font-black text-indigo-600">R$ {Number(order.total).toFixed(2)}</span>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </main>
    </div>
  )
}

export default App