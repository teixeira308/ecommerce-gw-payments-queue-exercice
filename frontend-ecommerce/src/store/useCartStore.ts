import { create } from 'zustand';
import { CartItem, Product, PaymentMethod } from '../types';

interface CartState {
  cart: CartItem[];
  paymentMethod: PaymentMethod;
  addToCart: (product: Product) => void;
  removeFromCart: (productId: string) => void;
  clearCart: () => void;
  setPaymentMethod: (method: PaymentMethod) => void;
  total: number;
}

export const useCartStore = create<CartState>((set, get) => ({
  cart: [],
  paymentMethod: 'CREDIT_CARD',
  total: 0,

  addToCart: (product) => {
    set((state) => {
      const exists = state.cart.find((i) => i.id === product.id);
      let newCart;
      
      if (exists) {
        newCart = state.cart.map((i) => 
          i.id === product.id ? { ...i, quantity: i.quantity + 1 } : i
        );
      } else {
        newCart = [...state.cart, { ...product, quantity: 1 }];
      }

      const total = newCart.reduce((acc, curr) => acc + (curr.price * curr.quantity), 0);
      return { cart: newCart, total };
    });
  },

  removeFromCart: (productId) => {
    set((state) => {
      const newCart = state.cart.filter((item) => item.id !== productId);
      const total = newCart.reduce((acc, curr) => acc + (curr.price * curr.quantity), 0);
      return { cart: newCart, total };
    });
  },

  clearCart: () => set({ cart: [], total: 0 }),

  setPaymentMethod: (method) => set({ paymentMethod: method }),
}));
