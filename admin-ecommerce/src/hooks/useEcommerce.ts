import { useState, useRef, useEffect } from 'react';
import { useQuery, useMutation, useQueryClient } from '@tanstack/react-query';
import { toast } from 'sonner';
import api from '../services/api';
import { Payment, Item, Stats } from '../types';

export const useEcommerce = () => {
  const queryClient = useQueryClient();
  const [activeTab, setActiveTab] = useState<'transactions' | 'catalog'>('transactions');
  const [autoRefresh, setAutoRefresh] = useState(true);
  const [newPaymentModal, setNewPaymentModal] = useState<Payment | null>(null);
  const seenIds = useRef<Set<string>>(new Set());
  const isFirstLoad = useRef(true);

  // Queries
  const { data: payments = [], isLoading: loadingPayments, refetch: refetchPayments } = useQuery<Payment[]>({
    queryKey: ['payments'],
    queryFn: () => api.payments.getAll(1, 100),
    refetchInterval: autoRefresh ? 30000 : false,
  });

  const { data: items = [], isLoading: loadingItems, refetch: refetchItems } = useQuery<Item[]>({
    queryKey: ['items'],
    queryFn: api.items.getAll,
    refetchInterval: autoRefresh ? 30000 : false,
  });

  // Stats computation
  const stats: Stats = {
    pending: payments.filter(p => p.status === 'PENDING').length,
    totalAmount: payments
      .filter(p => p.status === 'APPROVED')
      .reduce((acc, curr) => acc + Number(curr.amount || 0), 0)
  };

  // Notification logic
  useEffect(() => {
    if (!payments.length) return;
    
    const pendingOnes = payments.filter(p => p.status === 'PENDING');
    const newOnes = pendingOnes.filter(p => !seenIds.current.has(p.id));

    if (!isFirstLoad.current && newOnes.length > 0) {
      setNewPaymentModal(newOnes[0]);
      setActiveTab('transactions');
      toast.info(`Novo pagamento de R$ ${Number(newOnes[0].amount).toFixed(2)}`);
    }

    payments.forEach(p => seenIds.current.add(p.id));
    isFirstLoad.current = false;
  }, [payments]);

  // Mutations
  const processMutation = useMutation({
    mutationFn: (id: string) => api.payments.process(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payments'] });
      toast.success("Pagamento enviado com sucesso!");
    },
    onError: (err: any) => toast.error(err.message)
  });

  const statusMutation = useMutation({
    mutationFn: ({ id, status }: { id: string, status: string }) => api.payments.updateStatus(id, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['payments'] });
      toast.success("Status atualizado!");
    },
    onError: () => toast.error("Erro ao atualizar status")
  });

  const saveItemMutation = useMutation({
    mutationFn: ({ item, data }: { item: Item | null, data: { name: string, price: number } }) => api.items.save(item, data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['items'] });
      toast.success("Item salvo com sucesso!");
    },
    onError: () => toast.error("Erro ao salvar item")
  });

  const deleteItemMutation = useMutation({
    mutationFn: (id: string) => api.items.delete(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['items'] });
      toast.success("Item excluído!");
    },
    onError: () => toast.error("Erro ao excluir item")
  });

  const refreshAll = () => {
    toast.promise(Promise.all([refetchPayments(), refetchItems()]), {
      loading: 'Sincronizando...',
      success: 'Dados atualizados!',
      error: 'Erro na sincronização',
    });
  };

  return {
    activeTab, setActiveTab,
    payments, items,
    loading: loadingPayments || loadingItems,
    stats,
    autoRefresh, setAutoRefresh,
    processing: {
      ...processMutation.isPending && { [processMutation.variables as string]: true },
      ...statusMutation.isPending && { [statusMutation.variables?.id as string]: true }
    } as Record<string, boolean>,
    newPaymentModal, setNewPaymentModal,
    fetchData: refreshAll,
    handleProcess: (id: string) => processMutation.mutate(id),
    handleReject: (id: string) => statusMutation.mutate({ id, status: 'REJECTED' }),
    saveItem: (item: Item | null, data: { name: string, price: number }) => saveItemMutation.mutate({ item, data }),
    deleteItem: (id: string) => {
      if (window.confirm("Deseja realmente excluir este item?")) {
        deleteItemMutation.mutate(id);
      }
    }
  };
};
