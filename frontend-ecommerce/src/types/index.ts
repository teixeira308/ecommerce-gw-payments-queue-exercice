export type PaymentMethod = 'CREDIT_CARD' | 'PIX' | 'BOLETO';

export type OrderStatus = 'PENDING' | 'APPROVED' | 'CANCELLED';

export interface Product {
  id: string;
  name: string;
  price: number;
}

export interface CartItem extends Product {
  quantity: number;
}

export interface OrderItem {
  item_id: string;
  quantity: number;
  subtotal: number;
}

export interface Order {
  id: string;
  items: OrderItem[];
  total: number;
  method: PaymentMethod;
  status: OrderStatus;
  created_at: string;
}
