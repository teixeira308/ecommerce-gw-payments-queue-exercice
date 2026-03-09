import { motion } from 'framer-motion';
import { Plus, Package } from 'lucide-react';
import { Product } from '../../types';
import { useCartStore } from '../../store/useCartStore';
import { Button } from '../../components/ui/Button';
import { toast } from 'sonner';

interface ProductCardProps {
  product: Product;
}

export const ProductCard = ({ product }: ProductCardProps) => {
  const addToCart = useCartStore((state) => state.addToCart);

  const handleAdd = () => {
    addToCart(product);
    toast.success(`${product.name} adicionado ao carrinho!`, {
      icon: <Plus className="w-4 h-4" />,
    });
  };

  return (
    <motion.div 
      initial={{ opacity: 0, y: 20 }}
      animate={{ opacity: 1, y: 0 }}
      whileHover={{ y: -8 }}
      className="group bg-white/5 backdrop-blur-sm p-3 rounded-[3rem] border border-white/10 hover:border-indigo-500/50 transition-all duration-500"
    >
      <div className="bg-slate-900/50 rounded-[2.5rem] p-10 flex items-center justify-center relative overflow-hidden">
        <Package className="w-20 h-20 text-indigo-400 group-hover:scale-110 transition-transform duration-700 z-10 opacity-80" />
        <div className="absolute inset-0 bg-gradient-to-br from-indigo-500/10 to-transparent opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
      </div>
      <div className="p-8">
        <h3 className="font-bold text-xl text-white mb-2">{product.name}</h3>
        <div className="flex items-center justify-between mt-6">
          <div className="flex flex-col">
            <span className="text-[10px] text-slate-500 font-black uppercase tracking-widest mb-1">Preço</span>
            <span className="text-2xl font-black text-white">R$ {product.price.toFixed(2)}</span>
          </div>
          <Button 
            size="icon" 
            onClick={handleAdd}
            aria-label={`Adicionar ${product.name}`}
          >
            <Plus className="w-6 h-6" />
          </Button>
        </div>
      </div>
    </motion.div>
  );
};
