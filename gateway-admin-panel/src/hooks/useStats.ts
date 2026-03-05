import { useState, useCallback } from 'react';
import { Stats } from '../types';
import { PaymentService } from '../services/api';

export const useStats = () => {
  const [stats, setStats] = useState<Stats>({ pending: 0, totalAmount: 0 });

  const fetchStats = useCallback(async () => {
    try {
      const data = await PaymentService.getAllPayments();
      const all = data || [];

      const pendingCount = all.filter(p => p.status === 'PENDING').length;
      const approvedTotal = all
        .filter(p => p.status === 'APPROVED')
        .reduce((acc, curr) => acc + Number(curr.amount), 0);

      setStats({ pending: pendingCount, totalAmount: approvedTotal });
    } catch (err) {
      console.error("Error fetching global stats:", err);
    }
  }, []);

  return { stats, fetchStats };
};
