import { useState } from 'react';
import { useEcommerce } from './hooks/useEcommerce';
import DashboardLayout from './layouts/DashboardLayout';
import PaymentTable from './features/payments/PaymentTable';
import CatalogGrid from './features/catalog/CatalogGrid';
import ItemFormModal from './features/catalog/ItemFormModal';
import NotificationModal from './features/payments/NotificationModal';
import Button from './components/ui/Button';
import { Item } from './types';

function App() {
  const ecommerce = useEcommerce();
  const [filterPending, setFilterPending] = useState(false);
  const [itemModal, setItemModal] = useState<{ open: boolean, item: Item | null }>({ open: false, item: null });

  return (
    <DashboardLayout
      activeTab={ecommerce.activeTab}
      setActiveTab={ecommerce.setActiveTab}
      stats={ecommerce.stats}
      autoRefresh={ecommerce.autoRefresh}
      setAutoRefresh={ecommerce.setAutoRefresh}
      onRefresh={ecommerce.fetchData}
      loading={ecommerce.loading}
    >
      <NotificationModal 
        payment={ecommerce.newPaymentModal} 
        onClose={() => ecommerce.setNewPaymentModal(null)} 
      />

      <ItemFormModal 
        isOpen={itemModal.open}
        item={itemModal.item}
        onClose={() => setItemModal({ open: false, item: null })}
        onSave={(data) => ecommerce.saveItem(itemModal.item, data)}
      />

      <div className="mb-8 flex justify-end">
        {ecommerce.activeTab === 'catalog' && (
          <Button onClick={() => setItemModal({ open: true, item: null })}>
            + NOVO ITEM
          </Button>
        )}
      </div>

      {ecommerce.activeTab === 'transactions' ? (
        <PaymentTable 
          payments={ecommerce.payments}
          filterPending={filterPending}
          setFilterPending={setFilterPending}
          onProcess={ecommerce.handleProcess}
          onReject={ecommerce.handleReject}
          processing={ecommerce.processing}
        />
      ) : (
        <CatalogGrid 
          items={ecommerce.items}
          onEdit={(item) => setItemModal({ open: true, item })}
          onDelete={ecommerce.deleteItem}
        />
      )}
    </DashboardLayout>
  );
}

export default App;
