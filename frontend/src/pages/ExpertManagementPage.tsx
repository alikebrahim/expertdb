import { useState, useEffect, useCallback } from 'react';
import { Expert } from '../types';
import { expertsApi } from '../services/api';
import ExpertTable from '../components/ExpertTable';
import Button from '../components/ui/Button';
import ExpertForm from '../components/ExpertForm';
import Modal from '../components/Modal';
import Header from '../components/layout/Header';
import Sidebar from '../components/layout/Sidebar';
import Footer from '../components/layout/Footer';

const ExpertManagementPage = () => {
  const [experts, setExperts] = useState<Expert[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [limit] = useState(10);
  
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedExpert, setSelectedExpert] = useState<Expert | null>(null);
  const [isActionLoading, setIsActionLoading] = useState(false);
  const [actionError, setActionError] = useState<string | null>(null);
  
  const fetchExperts = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await expertsApi.getExperts(page, limit);
      
      if (response.success) {
        setExperts(response.data.data);
        setTotalPages(response.data.totalPages);
      } else {
        setError(response.message || 'Failed to load experts');
      }
    } catch (error) {
      console.error('Error fetching experts:', error);
      setError('An error occurred while loading experts');
    } finally {
      setIsLoading(false);
    }
  }, [page, limit]);
  
  useEffect(() => {
    fetchExperts();
  }, [fetchExperts]);
  
  const handleCreateExpert = (expert: Expert) => {
    setExperts([expert, ...experts]);
    setIsCreateModalOpen(false);
  };
  
  const handleEditExpert = (expert: Expert) => {
    setExperts(experts.map(e => e.id === expert.id ? expert : e));
    setIsEditModalOpen(false);
    setSelectedExpert(null);
  };
  
  const handleDeleteExpert = async () => {
    if (!selectedExpert) return;
    
    setIsActionLoading(true);
    setActionError(null);
    
    try {
      const response = await expertsApi.deleteExpert(selectedExpert.id.toString());
      
      if (response.success) {
        setExperts(experts.filter(e => e.id !== selectedExpert.id));
        setIsDeleteModalOpen(false);
        setSelectedExpert(null);
      } else {
        setActionError(response.message || 'Failed to delete expert');
      }
    } catch (error) {
      console.error('Error deleting expert:', error);
      setActionError('An error occurred while deleting the expert');
    } finally {
      setIsActionLoading(false);
    }
  };
  
  const handlePageChange = (newPage: number) => {
    setPage(newPage);
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
    <div className="min-h-screen flex flex-col">
      <Header />
      
      <div className="flex flex-1">
        <Sidebar />
        
        <main className="flex-1 p-6 bg-gray-50">
          <div className="max-w-7xl mx-auto">
            <div className="mb-6 flex justify-between items-center">
              <h1 className="text-2xl font-bold text-gray-900">Expert Management</h1>
              
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
            
            <div className="bg-white shadow rounded-lg p-6">
              <ExpertTable 
                experts={experts}
                isLoading={isLoading}
                error={error}
                pagination={{
                  currentPage: page,
                  totalPages,
                  onPageChange: handlePageChange,
                }}
                onEdit={openEditModal}
                onDelete={openDeleteModal}
              />
            </div>
          </div>
        </main>
      </div>
      
      <Footer />
      
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
          
          {actionError && (
            <div className="mb-4 p-2 text-sm bg-red-50 text-red-600 rounded">
              {actionError}
            </div>
          )}
          
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
              isLoading={isActionLoading}
            >
              Delete Expert
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
};

export default ExpertManagementPage;