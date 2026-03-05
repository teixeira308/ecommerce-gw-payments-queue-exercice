import { useState, useCallback } from 'react'
import { Toaster, toast } from 'sonner'
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { motion, AnimatePresence } from 'framer-motion'
import { Filter, RotateCw, ChevronLeft, ChevronRight, LayoutGrid, List } from 'lucide-react'

import Header from './components/Header'
import PaymentModal from './components/PaymentModal'
import PaymentTable from './components/PaymentTable'
import { PaymentService } from './services/api'
import { Payment } from './types'

function App() {
  const queryClient = useQueryClient()
  const [page, setPage] = useState(1)
  const [limit] = useState(10)
  const [filterPending, setFilterPending] = useState(false)
  const [autoRefresh, setAutoRefresh] = useState(true)
  const [newPaymentModal, setNewPaymentModal] = useState<Payment | null>(null)

  // React Query - Payments
  const { data: payments = [], isLoading, isFetching } = useQuery({
    queryKey: ['payments', page, limit],
    queryFn: () => PaymentService.getPayments(page, limit),
    refetchInterval: autoRefresh ? 30000 : false,
  })

  // React Query - Stats
  const { data: allPayments = [] } = useQuery({
    queryKey: ['stats'],
    queryFn: () => PaymentService.getAllPayments(),
    refetchInterval: autoRefresh ? 30000 : false,
  })

  const stats = {
    pending: allPayments.filter(p => p.status === 'PENDING').length,
    totalAmount: allPayments
      .filter(p => p.status === 'APPROVED')
      .reduce((acc, curr) => acc + Number(curr.amount), 0)
  }

  // Mutation for updating status
  const mutation = useMutation({
    mutationFn: ({ id, status }: { id: string, status: 'APPROVED' | 'REJECTED' }) => 
      PaymentService.updatePaymentStatus(id, status),
    onSuccess: (_, variables) => {
      toast.success(`Pagamento ${variables.status === 'APPROVED' ? 'aprovado' : 'rejeitado'}!`);
      queryClient.invalidateQueries({ queryKey: ['payments'] });
      queryClient.invalidateQueries({ queryKey: ['stats'] });
    },
    onError: () => toast.error('Erro ao processar ação')
  })

  const handleAction = (id: string, status: 'APPROVED' | 'REJECTED') => {
    mutation.mutate({ id, status })
  }

  const displayedPayments = filterPending
    ? payments.filter(p => p.status === 'PENDING')
    : payments

  return (
    <div className="min-h-screen bg-[#F8FAFC] text-slate-900 font-sans selection:bg-indigo-100 selection:text-indigo-700">
      <Toaster position="bottom-center" richColors />
      
      <PaymentModal 
        payment={newPaymentModal} 
        onClose={() => setNewPaymentModal(null)} 
      />

      <Header 
        stats={stats} 
        autoRefresh={autoRefresh} 
        onToggleAutoRefresh={() => setAutoRefresh(!autoRefresh)}
        loading={isFetching}
      />

      <main className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 sm:py-12">
        
        {/* Page Title & Actions */}
        <div className="flex flex-col md:flex-row md:items-end justify-between gap-6 mb-10">
          <div>
            <h2 className="text-3xl font-black tracking-tight text-slate-900 sm:text-4xl">Transações</h2>
            <p className="text-slate-500 mt-2 font-medium">Gerencie e monitore todos os pagamentos em tempo real.</p>
          </div>

          <div className="flex items-center gap-3">
            <div className="bg-white p-1 rounded-2xl border border-slate-200 shadow-sm flex">
              <button
                onClick={() => setFilterPending(false)}
                className={`flex items-center gap-2 px-5 py-2 rounded-xl text-xs font-bold transition-all ${!filterPending ? 'bg-slate-900 text-white shadow-lg' : 'text-slate-400 hover:text-slate-600'}`}
              >
                Todas
              </button>
              <button
                onClick={() => setFilterPending(true)}
                className={`flex items-center gap-2 px-5 py-2 rounded-xl text-xs font-bold transition-all ${filterPending ? 'bg-amber-500 text-white shadow-lg shadow-amber-100' : 'text-slate-400 hover:text-slate-600'}`}
              >
                Pendentes
                {stats.pending > 0 && (
                  <span className={`px-1.5 py-0.5 rounded-md text-[10px] ${filterPending ? 'bg-white/20' : 'bg-amber-100 text-amber-600'}`}>
                    {stats.pending}
                  </span>
                )}
              </button>
            </div>

            <button 
              onClick={() => queryClient.invalidateQueries()}
              className="p-3 bg-white border border-slate-200 rounded-2xl text-slate-400 hover:text-indigo-600 hover:border-indigo-100 transition-all active:scale-95 shadow-sm"
              title="Recarregar dados"
            >
              <RotateCw size={20} className={isFetching ? 'animate-spin' : ''} />
            </button>
          </div>
        </div>

        {/* Main Content Card */}
        <motion.div 
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          className="bg-white rounded-[2.5rem] shadow-2xl shadow-slate-200/60 border border-slate-200 overflow-hidden"
        >
          {isLoading ? (
            <div className="p-32 flex flex-col items-center gap-4">
               <div className="w-12 h-12 border-4 border-indigo-600 border-t-transparent rounded-full animate-spin" />
               <p className="text-slate-400 font-bold text-sm uppercase tracking-widest">Carregando dados...</p>
            </div>
          ) : (
            <PaymentTable 
              payments={displayedPayments} 
              processing={{ [mutation.variables?.id || '']: mutation.isPending }} 
              onAction={handleAction} 
            />
          )}

          {/* Pagination Footer */}
          <div className="bg-slate-50/50 p-6 sm:p-8 border-t border-slate-100 flex flex-col sm:flex-row items-center justify-between gap-6">
            <div className="flex items-center gap-4">
               <p className="text-xs font-black text-slate-400 uppercase tracking-[0.2em]">
                 Página <span className="text-indigo-600 ml-1">{page}</span>
               </p>
               <div className="h-4 w-px bg-slate-200" />
               <p className="text-xs font-bold text-slate-400 uppercase">
                 Mostrando {displayedPayments.length} resultados
               </p>
            </div>
            
            <div className="flex gap-3 w-full sm:w-auto">
              <button
                onClick={() => setPage(prev => Math.max(prev - 1, 1))}
                disabled={page === 1}
                className="flex-1 sm:flex-none flex items-center justify-center gap-2 px-6 py-3 bg-white border border-slate-200 rounded-2xl text-xs font-black disabled:opacity-30 hover:border-indigo-200 transition-all shadow-sm active:scale-95"
              >
                <ChevronLeft size={16} />
                ANTERIOR
              </button>
              <button
                onClick={() => setPage(prev => prev + 1)}
                disabled={payments.length < limit}
                className="flex-1 sm:flex-none flex items-center justify-center gap-2 px-6 py-3 bg-white border border-slate-200 rounded-2xl text-xs font-black hover:border-indigo-200 transition-all shadow-sm active:scale-95"
              >
                PRÓXIMA
                <ChevronRight size={16} />
              </button>
            </div>
          </div>
        </motion.div>

      </main>
    </div>
  )
}

export default App;
