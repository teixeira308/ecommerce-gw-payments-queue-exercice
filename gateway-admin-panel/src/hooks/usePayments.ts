import { useState, useCallback, useRef } from 'react';
import { Payment } from '../types';
import { PaymentService } from '../services/api';

export const usePayments = (page: number, limit: number, onNewPayment?: (payment: Payment) => void) => {
  const [payments, setPayments] = useState<Payment[]>([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const seenIds = useRef<Set<string>>(new Set());
  const isFirstLoad = useRef(true);

  const fetchPayments = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await PaymentService.getPayments(page, limit);
      const paymentsData = data || [];

      // Check for new payments
      const newOnes = paymentsData.filter(p => p.status === 'PENDING' && !seenIds.current.has(p.id));
      if (!isFirstLoad.current && newOnes.length > 0 && onNewPayment) {
        onNewPayment(newOnes[0]);
      }

      // Update seen IDs
      paymentsData.forEach(p => seenIds.current.add(p.id));
      isFirstLoad.current = false;
      setPayments(paymentsData);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }, [page, limit, onNewPayment]);

  return { payments, loading, error, fetchPayments };
};
