import { Product, Order, PaymentMethod, CartItem } from '../types';

const API_URL = 'http://localhost:3000';

export const apiService = {
  async getProducts(): Promise<Product[]> {
    const res = await fetch(`${API_URL}/items`);
    if (!res.ok) {
      throw new Error("Failed to fetch products");
    }
    const json = await res.json();
    const data = Array.isArray(json) ? json : [];
    return data.map((item: any) => ({
      id: item.id ?? item.ID,
      name: item.name ?? item.Name ?? "Produto",
      price: Number(item.price ?? item.Price ?? 0)
    }));
  },

  async getOrders(page = 1, limit = 100): Promise<Order[]> {
    const res = await fetch(`${API_URL}/orders?limit=${limit}&page=${page}`);

    if (!res.ok) {
      throw new Error('Failed to fetch orders');
    }

    const json = await res.json();
    const data = Array.isArray(json) ? json : [];

    return data.map((order: any) => ({
      id: order.id,
      items: (order.items ?? []).map((item: any) => ({
        item_id: item.item_id,
        quantity: item.quantity,
        subtotal: Number(item.subtotal ?? 0)
      })),
      total: Number(order.total ?? 0),
      method: order.method as PaymentMethod,
      status: order.status as any,
      created_at: order.created_at
    }));
  },

  async createOrder(
    items: { item_id: string; quantity: number }[],
    method: PaymentMethod
  ): Promise<void> {
    const res = await fetch(`${API_URL}/orders`, {
      method: "POST",
      headers: {
        "Content-Type": "application/json"
      },
      body: JSON.stringify({
        method: method.toLowerCase(),
        items
      })
    });

    if (!res.ok) {
      throw new Error("Failed to create order");
    }
  }
};
