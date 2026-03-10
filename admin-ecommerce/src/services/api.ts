import { Payment, Item } from '../types';

const GATEWAY_URL = 'http://localhost:9000';
const ECOMMERCE_URL = 'http://localhost:3000';

const api = {
  payments: {
    getAll: async (page = 1, limit = 10): Promise<Payment[]> => {
      const res = await fetch(`${GATEWAY_URL}/payments?limit=${limit}&page=${page}`);
      const json = await res.json();
      return json.data || [];
    },
    updateStatus: async (id: string, status: string): Promise<Payment> => {
      const res = await fetch(`${GATEWAY_URL}/payments/${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ status })
      });
      const json = await res.json();
      return json.data || json;
    },
    process: async (id: string): Promise<any> => {
      const res = await fetch(`${GATEWAY_URL}/payments/${id}/process`, { method: 'POST' });
      if (!res.ok) throw new Error("Erro no gateway");
      const json = await res.json();
      return json.data || json;
    }
  },
  items: {
    getAll: async (): Promise<Item[]> => {
      const res = await fetch(`${ECOMMERCE_URL}/items`);
      const json = await res.json();
      return Array.isArray(json) ? json : (json.data || []);
    },
    save: async (item: Item | null, formData: { name: string, price: number }): Promise<Item> => {
      const isEdit = !!item;
      const url = isEdit ? `${ECOMMERCE_URL}/items/${item.ID}` : `${ECOMMERCE_URL}/items`;
      
      const body = {
        Name: formData.name,
        Price: formData.price
      };

      const res = await fetch(url, {
        method: isEdit ? 'PUT' : 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
      });
      if (!res.ok) throw new Error("Erro ao salvar");
      const json = await res.json();
      return json.data || json;
    },
    delete: async (id: string): Promise<boolean> => {
      const res = await fetch(`${ECOMMERCE_URL}/items/${id}`, { method: 'DELETE' });
      if (!res.ok) throw new Error("Erro ao excluir");
      return true;
    }
  }
};

export default api;
