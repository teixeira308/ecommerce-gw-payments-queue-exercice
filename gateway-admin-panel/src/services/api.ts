import { Payment } from '../types';

const GATEWAY_URL = import.meta.env.VITE_GATEWAY_URL || 'http://localhost:4000';

export const PaymentService = {
  async getPayments(page = 1, limit = 10): Promise<Payment[]> {
    const res = await fetch(`${GATEWAY_URL}/payments?limit=${limit}&page=${page}`);
    if (!res.ok) throw new Error("Failed to fetch payments");
    const json = await res.json();
    return json.data || [];
  },

  async getAllPayments(): Promise<Payment[]> {
    const res = await fetch(`${GATEWAY_URL}/payments?limit=9999`);
    if (!res.ok) throw new Error("Failed to fetch all payments for stats");
    const json = await res.json();
    return json.data || [];
  },

  async updatePaymentStatus(id: string, status: 'APPROVED' | 'REJECTED'): Promise<Payment> {
    const res = await fetch(`${GATEWAY_URL}/payments/${id}/status`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ status })
    });
    if (!res.ok) throw new Error("Failed to update payment status");
    const json = await res.json();
    return json.data || json;
  }
};
