import React, { useState, useEffect } from 'react';
import Modal from '../../components/ui/Modal';
import Input from '../../components/ui/Input';
import Button from '../../components/ui/Button';
import { Item } from '../../types';

interface ItemFormModalProps {
  isOpen: boolean;
  item: Item | null;
  onClose: () => void;
  onSave: (data: { name: string, price: number }) => void;
}

const ItemFormModal: React.FC<ItemFormModalProps> = ({ isOpen, item, onClose, onSave }) => {
  const [formData, setFormData] = useState({ name: '', price: 0 });
  const [isSaving, setIsSaving] = useState(false);

  useEffect(() => {
    if (item) {
      setFormData({ name: item.Name, price: item.Price });
    } else {
      setFormData({ name: '', price: 0 });
    }
  }, [item, isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsSaving(true);
    try {
      await onSave({ ...formData, price: Number(formData.price || 0) });
      onClose();
    } catch (err) {
      console.error(err);
    } finally {
      setIsSaving(false);
    }
  };

  return (
    <Modal 
      isOpen={isOpen} 
      onClose={onClose} 
      title={item ? 'Editar Item' : 'Novo Item'}
    >
      <form onSubmit={handleSubmit} className="space-y-6">
        <Input 
          label="Nome do Produto"
          required
          placeholder="Ex: Teclado Mecânico"
          value={formData.name}
          onChange={e => setFormData({...formData, name: e.target.value})}
        />
        <Input 
          label="Preço (R$)"
          type="number"
          step="0.01"
          required
          placeholder="0.00"
          value={formData.price}
          onChange={e => setFormData({...formData, price: Number(e.target.value)})}
        />
        <div className="flex gap-3 pt-4">
          <Button type="button" variant="secondary" className="flex-1" onClick={onClose}>
            CANCELAR
          </Button>
          <Button type="submit" className="flex-1" loading={isSaving}>
            SALVAR ITEM
          </Button>
        </div>
      </form>
    </Modal>
  );
};

export default ItemFormModal;
