import { useState, useCallback } from 'react';
import { Expert } from '../types';
import { expertsApi } from '../services/api';
import Layout from '../components/layout/Layout';
import ExpertTable from '../components/tables/ExpertTable';
import ExpertFilters from '../components/ExpertFilters';
import Button from '../components/ui/Button';
import ExpertForm from '../components/ExpertForm';
import Modal from '../components/Modal';
import { useFetch, useOptimisticCollection } from '../hooks';
import { LoadingOverlay } from '../components/ui/LoadingSpinner';

interface ExpertFiltersType {
  name?: string;
  role?: string;
  type?: string;
  affiliation?: string;
  expertAreaId?: string;
  nationality?: string;
  rating?: string;
  isAvailable?: boolean;
  isBahraini?: boolean;
}

const ExpertManagementPage = () => {
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [filters, setFilters] = useState<ExpertFiltersType>({
    isAvailable: true // Default to show only available experts
  });
  
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [isAreaManagementOpen, setIsAreaManagementOpen] = useState(false);
  const [selectedExpert, setSelectedExpert] = useState<Expert | null>(null);
  
  // Create fetch function based on current filters and pagination
  const fetchExperts = useCallback(async () => {
    // Convert filters to API-friendly params
    const params: Record<string, string | boolean | number> = {};
    
    if (filters.name) params.name = filters.name;
    if (filters.role) params.role = filters.role;
    if (filters.type) params.employmentType = filters.type;
    if (filters.affiliation) params.affiliation = filters.affiliation;
    if (filters.expertAreaId) params.generalArea = parseInt(filters.expertAreaId);
    if (filters.nationality) params.nationality = filters.nationality;
    if (filters.rating) params.minRating = parseInt(filters.rating);
    if (filters.isAvailable !== undefined) params.isAvailable = filters.isAvailable;
    if (filters.isBahraini !== undefined) params.isBahraini = filters.isBahraini;
    
    const offset = (page - 1) * limit;
    const response = await expertsApi.getExperts(limit, offset, params);
    
    if (response.success && response.data) {
      return {
        experts: response.data.experts,
        totalPages: response.data.pagination.totalPages
      };
    } else {
      throw new Error(response.message || 'Failed to fetch experts');
    }
  }, [page, limit, filters]);
  
  // Use optimistic collection for expert management
  const {
    items: experts,
    isLoading,
    error,
    addItem,
    updateItem,
    deleteItem
  } = useOptimisticCollection<Expert>(
    () => fetchExperts().then(data => data.experts)
  );
  
  // Use separate fetch for pagination info
  const { data: paginationData } = useFetch(fetchExperts, {
    deps: [page, limit, filters],
    showErrorNotifications: false,
  });
  
  const totalPages = paginationData?.totalPages || 1;
  
  const handleCreateExpert = async (expert: Expert) => {
    try {
      // Since the API expects FormData, this would normally be handled by the form
      // For now, just close the modal and refresh
      setIsCreateModalOpen(false);
      // Refresh the list
      window.location.reload();
    } catch (error) {
      console.error('Error creating expert:', error);
    }
  };
  
  const handleEditExpert = async (expert: Expert) => {
    try {
      const response = await expertsApi.updateExpert(expert.id.toString(), expert as any);
      if (response.success) {
        await updateItem(expert, () => Promise.resolve(expert), {
          successMessage: 'Expert updated successfully',
          errorMessage: 'Failed to update expert'
        });
        setIsEditModalOpen(false);
        setSelectedExpert(null);
      }
    } catch (error) {
      // Error is handled by the hook
      console.error('Error updating expert:', error);
    }
  };
  
  const handleDeleteExpert = async () => {
    if (!selectedExpert) return;
    
    try {
      await deleteItem(
        selectedExpert.id,
        (id) => expertsApi.deleteExpert(id.toString()).then(() => {}),
        {
          successMessage: 'Expert deleted successfully',
          errorMessage: 'Failed to delete expert'
        }
      );
      setIsDeleteModalOpen(false);
      setSelectedExpert(null);
    } catch (error) {
      // Error is handled by the hook
      console.error('Error deleting expert:', error);
    }
  };
  
  const handleBatchExport = (expertIds: number[]) => {
    // Implement batch export functionality
    console.log('Exporting experts:', expertIds);
  };
  
  const handleBatchPublish = (expertIds: number[]) => {
    // Implement batch publish functionality
    console.log('Publishing experts:', expertIds);
  };
  
  const handleBatchUnpublish = (expertIds: number[]) => {
    // Implement batch unpublish functionality
    console.log('Unpublishing experts:', expertIds);
  };
  
  const handlePageChange = (newPage: number) => {
    setPage(newPage);
  };
  
  const handleFilterChange = (newFilters: ExpertFiltersType) => {
    setFilters(newFilters);
    setPage(1); // Reset to first page when filters change
  };
  
  const openEditModal = (expert: Expert) => {
    setSelectedExpert(expert);
    setIsEditModalOpen(true);
  };
  
  const openDeleteModal = (expert: Expert) => {
    setSelectedExpert(expert);
    setIsDeleteModalOpen(true);
  };
  
  return (
    <Layout>
      <div className="mb-6 flex flex-wrap justify-between items-center gap-4">
        <h1 className="text-2xl font-bold text-primary">Expert Management</h1>
        
        <div className="flex flex-wrap gap-2">
          <Button 
            onClick={() => setIsAreaManagementOpen(true)}
            variant="secondary"
          >
            Manage Expert Areas
          </Button>
          <Button 
            onClick={() => setIsCreateModalOpen(true)}
            icon={
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 5a1 1 0 011 1v3h3a1 1 0 110 2h-3v3a1 1 0 11-2 0v-3H6a1 1 0 110-2h3V6a1 1 0 011-1z" clipRule="evenodd" />
              </svg>
            }
          >
            Add Expert
          </Button>
        </div>
      </div>
      
      <ExpertFilters onFilterChange={handleFilterChange} initialFilters={filters} />
      
      <LoadingOverlay 
        isLoading={isLoading && experts.length === 0}
        className="bg-white shadow rounded-lg overflow-hidden"
        label="Loading experts..."
      >
        <ExpertTable 
          experts={experts}
          isLoading={isLoading}
          error={error ? error.message : null}
          pagination={{
            currentPage: page,
            totalPages,
            onPageChange: handlePageChange,
          }}
          onEdit={openEditModal}
          onDelete={openDeleteModal}
          enableBatchActions={true}
          onBatchExport={handleBatchExport}
          onBatchPublish={handleBatchPublish}
          onBatchUnpublish={handleBatchUnpublish}
        />
      </LoadingOverlay>
      
      {/* Create Expert Modal */}
      <Modal
        isOpen={isCreateModalOpen}
        onClose={() => setIsCreateModalOpen(false)}
        title="Create New Expert"
        size="lg"
      >
        <ExpertForm
          onSuccess={handleCreateExpert}
          onCancel={() => setIsCreateModalOpen(false)}
        />
      </Modal>
      
      {/* Edit Expert Modal */}
      <Modal
        isOpen={isEditModalOpen}
        onClose={() => setIsEditModalOpen(false)}
        title="Edit Expert"
        size="lg"
      >
        {selectedExpert && (
          <ExpertForm
            expert={selectedExpert}
            onSuccess={handleEditExpert}
            onCancel={() => setIsEditModalOpen(false)}
          />
        )}
      </Modal>
      
      {/* Delete Expert Modal */}
      <Modal
        isOpen={isDeleteModalOpen}
        onClose={() => setIsDeleteModalOpen(false)}
        title="Delete Expert"
        size="sm"
      >
        <div className="text-center">
          <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100 mb-4">
            <svg className="h-6 w-6 text-red-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor" aria-hidden="true">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
            </svg>
          </div>
          
          <h3 className="text-lg font-medium text-gray-900 mb-2">Delete Expert</h3>
          
          <p className="text-gray-500 mb-4">
            Are you sure you want to delete <span className="font-semibold">{selectedExpert?.name}</span>? This action cannot be undone.
          </p>
          
          <div className="flex justify-center space-x-3 mt-4">
            <Button 
              variant="outline"
              onClick={() => setIsDeleteModalOpen(false)}
            >
              Cancel
            </Button>
            <Button 
              variant="danger"
              onClick={handleDeleteExpert}
            >
              Delete Expert
            </Button>
          </div>
        </div>
      </Modal>
      
      {/* Area Management Modal */}
      <Modal
        isOpen={isAreaManagementOpen}
        onClose={() => setIsAreaManagementOpen(false)}
        title="Expert Areas Management"
        size="md"
      >
        {/* TODO: Replace with ExpertAreaManagement component */}
        <div className="p-4">
          <p>Expert area management will be implemented in a future update.</p>
        </div>
        <div className="px-6 py-4 flex justify-end border-t">
          <Button 
            onClick={() => setIsAreaManagementOpen(false)}
          >
            Close
          </Button>
        </div>
      </Modal>
    </Layout>
  );
};

export default ExpertManagementPage;