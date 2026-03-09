import React from 'react';
import { motion } from 'framer-motion';
import Card from '../../components/ui/Card';
import Button from '../../components/ui/Button';
import { Item } from '../../types';

interface CatalogGridProps {
  items: Item[];
  onEdit: (item: Item) => void;
  onDelete: (id: string) => void;
}

const CatalogGrid: React.FC<CatalogGridProps> = ({ items, onEdit, onDelete }) => {
  const formatDate = (dateString: string) => {
    if (!dateString) return "-";
    return new Date(dateString).toLocaleDateString('pt-BR');
  };

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  };

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: {
      y: 0,
      opacity: 1
    }
  };

  return (
    <motion.div 
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
    >
      {items.map(item => (
        <motion.div 
          key={item.ID} 
          variants={itemVariants}
          whileHover={{ y: -5 }}
          className="h-full"
        >
          <Card className="relative group overflow-hidden h-full flex flex-col">
            <motion.div 
              initial={{ scale: 1 }}
              whileHover={{ scale: 1.1 }}
              className="absolute top-0 right-0 w-32 h-32 bg-indigo-50 rounded-full -mr-16 -mt-16 group-hover:bg-indigo-100 transition-colors"
            />
            
            <div className="relative z-10 flex-1">
              <div className="flex justify-between items-start mb-4">
                <div>
                  <span className="text-[10px] font-mono text-indigo-400 uppercase block mb-1 font-bold">
                    ID: {(item.ID || '').substring(0, 8)}
                  </span>
                  <h3 className="text-xl font-black text-slate-900 leading-tight">{item.Name}</h3>
                </div>
                <span className="bg-indigo-600 text-white px-4 py-1.5 rounded-xl text-sm font-black shadow-lg shadow-indigo-200">
                  R$ {Number(item.Price || 0).toFixed(2)}
                </span>
              </div>

              <div className="flex items-center justify-between pt-6 border-t border-slate-50 mt-auto">
                <div className="flex flex-col">
                  <span className="text-[9px] font-black text-slate-400 uppercase tracking-tighter">Criado em</span>
                  <span className="text-sm font-bold text-slate-700">{formatDate(item.CreatedAt)}</span>
                </div>
                <div className="flex gap-2">
                  <Button variant="secondary" size="icon" onClick={() => onEdit(item)}>
                    ✏️
                  </Button>
                  <Button variant="danger" size="icon" onClick={() => onDelete(item.ID)}>
                    🗑️
                  </Button>
                </div>
              </div>
            </div>
          </Card>
        </motion.div>
      ))}
      {items.length === 0 && (
        <motion.div 
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          className="col-span-full p-20 text-center bg-white rounded-[2.5rem] border border-dashed border-slate-200 text-slate-400 font-medium"
        >
          O catálogo está vazio. Comece adicionando um item.
        </motion.div>
      )}
    </motion.div>
  );
};

export default CatalogGrid;
