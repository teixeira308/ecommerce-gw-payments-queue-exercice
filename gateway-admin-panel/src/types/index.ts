export interface Payment {
  id: string;
  order_id: string;
  amount: number | string;
  status: 'PENDING' | 'APPROVED' | 'REJECTED';
  created_at: string;
  method: string;
}

export interface Stats {
  pending: number;
  totalAmount: number;
}

export interface PaymentFilters {
  page: number;
  limit: number;
  status?: string;
}

export interface ApiResponse<T> {
  data: T;
  meta?: {
    total: number;
    page: number;
    limit: number;
  };
}
