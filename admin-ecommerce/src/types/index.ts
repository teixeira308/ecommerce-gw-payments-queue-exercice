export interface Payment {
  id: string;
  order_id: string;
  amount: number;
  method: string;
  status: 'PENDING' | 'APPROVED' | 'REJECTED';
  created_at: string;
}

export interface Item {
  ID: string;
  Name: string;
  Price: number;
  CreatedAt: string;
}

export interface Stats {
  pending: number;
  totalAmount: number;
}
