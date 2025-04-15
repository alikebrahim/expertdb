import { useState, useEffect, useCallback } from 'react';
import { Engagement } from '../types';
import { engagementApi } from '../services/api';
import Button from './ui/Button';
import Table from './ui/Table';
import Modal from './Modal';
import EngagementForm from './EngagementForm';
import { useAuth } from '../hooks/useAuth';

interface EngagementListProps {
  expertId: number;
}

const EngagementList = ({ expertId }: EngagementListProps) => {
  const { user } = useAuth();
  const [engagements, setEngagements] = useState<Engagement[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 5,
    total: 0,
    totalPages: 0
  });
  const [isAddModalOpen, setIsAddModalOpen] = useState(false);
  const [isViewModalOpen, setIsViewModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedEngagement, setSelectedEngagement] = useState<Engagement | null>(null);

  const fetchEngagements = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await engagementApi.getExpertEngagements(
        expertId.toString(), 
        pagination.page, 
        pagination.limit
      );
      
      if (response.success) {
        setEngagements(response.data.data);
        setPagination({
          page: response.data.page,
          limit: response.data.limit,
          total: response.data.total,
          totalPages: response.data.totalPages
        });
      } else {
        setError(response.message || 'Failed to load engagements');
      }
    } catch (error) {
      console.error('Error fetching engagements:', error);
      setError('An error occurred while loading engagements');
    } finally {
      setIsLoading(false);
    }
  }, [expertId, pagination.page, pagination.limit]);

  useEffect(() => {
    fetchEngagements();
  }, [fetchEngagements]);

  const handlePageChange = (newPage: number) => {
    setPagination(prev => ({ ...prev, page: newPage }));
  };

  const handleEngagementCreated = (_engagement: Engagement) => {
    setIsAddModalOpen(false);
    fetchEngagements();
  };

  const handleEngagementUpdated = (_engagement: Engagement) => {
    setIsEditModalOpen(false);
    fetchEngagements();
  };

  const handleViewEngagement = (engagement: Engagement) => {
    setSelectedEngagement(engagement);
    setIsViewModalOpen(true);
  };

  const handleEditEngagement = (engagement: Engagement) => {
    setSelectedEngagement(engagement);
    setIsEditModalOpen(true);
  };

  const handleDeletePrompt = (engagement: Engagement) => {
    setSelectedEngagement(engagement);
    setIsDeleteModalOpen(true);
  };

  const handleDeleteEngagement = async () => {
    if (!selectedEngagement) return;
    
    try {
      const response = await engagementApi.deleteEngagement(selectedEngagement.id.toString());
      
      if (response.success) {
        setIsDeleteModalOpen(false);
        fetchEngagements();
      } else {
        setError(response.message || 'Failed to delete engagement');
      }
    } catch (error) {
      console.error('Error deleting engagement:', error);
      setError('An error occurred while deleting the engagement');
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  const getStatusBadgeColor = (status: string) => {
    switch (status) {
      case 'pending':
        return 'bg-yellow-100 text-yellow-800';
      case 'confirmed':
        return 'bg-blue-100 text-blue-800';
      case 'in_progress':
        return 'bg-green-100 text-green-800';
      case 'completed':
        return 'bg-gray-100 text-gray-800';
      case 'cancelled':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="bg-white rounded-md shadow-sm p-5">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-xl font-semibold">Engagements</h2>
        {user?.role === 'admin' && (
          <Button 
            variant="primary" 
            size="sm" 
            onClick={() => setIsAddModalOpen(true)}
          >
            Add Engagement
          </Button>
        )}
      </div>
      
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-600 p-4 rounded-md mb-4">
          {error}
        </div>
      )}

      {isLoading ? (
        <div className="text-center p-4">
          <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-primary mx-auto"></div>
          <p className="mt-2 text-gray-500">Loading engagements...</p>
        </div>
      ) : engagements.length === 0 ? (
        <div className="text-center p-6 bg-gray-50 rounded-md">
          <p className="text-gray-500">No engagements found for this expert.</p>
        </div>
      ) : (
        <>
          <Table>
            <Table.Header>
              <Table.Row>
                <Table.HeaderCell>Title</Table.HeaderCell>
                <Table.HeaderCell>Type</Table.HeaderCell>
                <Table.HeaderCell>Status</Table.HeaderCell>
                <Table.HeaderCell>Period</Table.HeaderCell>
                <Table.HeaderCell>Organization</Table.HeaderCell>
                <Table.HeaderCell>Actions</Table.HeaderCell>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              {engagements.map(engagement => (
                <Table.Row key={engagement.id}>
                  <Table.Cell>{engagement.title}</Table.Cell>
                  <Table.Cell className="capitalize">{engagement.engagementType}</Table.Cell>
                  <Table.Cell>
                    <span className={`px-2 py-1 rounded-full text-xs ${getStatusBadgeColor(engagement.status)}`}>
                      {engagement.status.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                    </span>
                  </Table.Cell>
                  <Table.Cell>
                    {formatDate(engagement.startDate)} - {formatDate(engagement.endDate)}
                  </Table.Cell>
                  <Table.Cell>{engagement.organizationName}</Table.Cell>
                  <Table.Cell>
                    <div className="flex space-x-2">
                      <Button 
                        variant="outline" 
                        size="sm" 
                        onClick={() => handleViewEngagement(engagement)}
                      >
                        View
                      </Button>
                      {user?.role === 'admin' && (
                        <>
                          <Button 
                            variant="secondary" 
                            size="sm" 
                            onClick={() => handleEditEngagement(engagement)}
                          >
                            Edit
                          </Button>
                          <Button 
                            variant="danger" 
                            size="sm" 
                            onClick={() => handleDeletePrompt(engagement)}
                          >
                            Delete
                          </Button>
                        </>
                      )}
                    </div>
                  </Table.Cell>
                </Table.Row>
              ))}
            </Table.Body>
          </Table>
          
          {pagination.totalPages > 1 && (
            <div className="flex justify-center mt-4">
              <div className="flex space-x-2">
                <Button
                  variant="outline"
                  size="sm"
                  disabled={pagination.page === 1}
                  onClick={() => handlePageChange(pagination.page - 1)}
                >
                  Previous
                </Button>
                <div className="flex items-center px-3">
                  Page {pagination.page} of {pagination.totalPages}
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  disabled={pagination.page === pagination.totalPages}
                  onClick={() => handlePageChange(pagination.page + 1)}
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </>
      )}
      
      {/* Add Engagement Modal */}
      <Modal
        isOpen={isAddModalOpen}
        onClose={() => setIsAddModalOpen(false)}
        title="Add New Engagement"
        size="lg"
      >
        <EngagementForm
          expertId={expertId}
          onSuccess={handleEngagementCreated}
          onCancel={() => setIsAddModalOpen(false)}
        />
      </Modal>
      
      {/* View Engagement Modal */}
      {selectedEngagement && (
        <Modal
          isOpen={isViewModalOpen}
          onClose={() => setIsViewModalOpen(false)}
          title={selectedEngagement.title}
          size="md"
        >
          <div className="space-y-4">
            <div>
              <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-1">Description</h3>
              <p>{selectedEngagement.description}</p>
            </div>
            
            <div className="grid grid-cols-2 gap-4">
              <div>
                <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-1">Type</h3>
                <p className="capitalize">{selectedEngagement.engagementType}</p>
              </div>
              <div>
                <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-1">Status</h3>
                <span className={`px-2 py-1 rounded-full text-xs ${getStatusBadgeColor(selectedEngagement.status)}`}>
                  {selectedEngagement.status.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                </span>
              </div>
            </div>
            
            <div className="grid grid-cols-2 gap-4">
              <div>
                <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-1">Start Date</h3>
                <p>{formatDate(selectedEngagement.startDate)}</p>
              </div>
              <div>
                <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-1">End Date</h3>
                <p>{formatDate(selectedEngagement.endDate)}</p>
              </div>
            </div>
            
            <div>
              <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-1">Organization</h3>
              <p>{selectedEngagement.organizationName}</p>
            </div>
            
            <div className="grid grid-cols-2 gap-4">
              <div>
                <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-1">Contact Person</h3>
                <p>{selectedEngagement.contactPerson}</p>
              </div>
              <div>
                <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-1">Contact Email</h3>
                <p>{selectedEngagement.contactEmail}</p>
              </div>
            </div>
            
            {selectedEngagement.notes && (
              <div>
                <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-1">Notes</h3>
                <p>{selectedEngagement.notes}</p>
              </div>
            )}
            
            <div className="flex justify-end pt-4">
              <Button onClick={() => setIsViewModalOpen(false)}>Close</Button>
            </div>
          </div>
        </Modal>
      )}
      
      {/* Edit Engagement Modal */}
      {selectedEngagement && (
        <Modal
          isOpen={isEditModalOpen}
          onClose={() => setIsEditModalOpen(false)}
          title="Edit Engagement"
          size="lg"
        >
          <EngagementForm
            expertId={expertId}
            engagement={selectedEngagement}
            onSuccess={handleEngagementUpdated}
            onCancel={() => setIsEditModalOpen(false)}
          />
        </Modal>
      )}
      
      {/* Delete Confirmation Modal */}
      <Modal
        isOpen={isDeleteModalOpen}
        onClose={() => setIsDeleteModalOpen(false)}
        title="Confirm Deletion"
        size="sm"
      >
        <div className="space-y-4">
          <p>Are you sure you want to delete this engagement?</p>
          <p className="text-sm text-gray-500">
            This action cannot be undone.
          </p>
          
          <div className="flex justify-end space-x-3 pt-4">
            <Button 
              variant="outline"
              onClick={() => setIsDeleteModalOpen(false)}
            >
              Cancel
            </Button>
            <Button 
              variant="danger"
              onClick={handleDeleteEngagement}
            >
              Delete
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
};

export default EngagementList;