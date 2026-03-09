import { Product, Order, PaymentMethod, CartItem } from '../types';

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000';

export const apiService = {
  async getProducts(): Promise<Product[]> {
    const res = await fetch(`${API_URL}/items`);
    if (!res.ok) throw new Error('Failed to fetch products');
    const json = await res.json();
    const data = Array.isArray(json.data) ? json.data : [];
    // Sanitização de dados devido à inconsistência observada no review
    return data.map((item: any) => ({
      id: item.id ?? item.ID,
      name: item.name ?? item.Name ?? "Produto",
      price: Number(item.price ?? item.Price ?? 0)
    }));
  },

  async getOrders(): Promise<Order[]> {
    const res = await fetch(`${API_URL}/orders`);
    if (!res.ok) throw new Error('Failed to fetch orders');
    const json = await res.json();
    const data = Array.isArray(json.data) ? json.data : [];
    
    return data.map((order: any) => ({
      id: order.id ?? order.ID,
      items: (order.items ?? order.Items ?? []).map((item: any) => ({
        item_id: item.item_id ?? item.ItemID,
        quantity: item.quantity ?? item.Quantity,
        subtotal: Number(item.subtotal ?? item.Subtotal ?? 0)
      })),
      total: Number(order.total ?? order.Total ?? 0),
      method: (order.method ?? order.Method) as PaymentMethod,
      status: (order.status ?? order.Status) as any,
      created_at: order.created_at ?? order.CreatedAt
    }));
  },

  async createOrder(items: { item_id: string, quantity: number }[], method: PaymentMethod): Promise<void> {
    const res = await fetch(`${API_URL}/orders`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ items, method })
    });
    if (!res.ok) throw new Error('Failed to create order');
  }
};
