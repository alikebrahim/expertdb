import { useState, useEffect, useCallback } from 'react';
import { engagementApi, expertsApi } from '../services/api';
import { Engagement, Expert } from '../types';
import { Button, Table, Input } from '../components/ui';
import Modal from '../components/Modal';
import EngagementForm from '../components/EngagementForm';

const EngagementManagementPage = () => {
  const [engagements, setEngagements] = useState<Engagement[]>([]);
  const [experts, setExperts] = useState<Expert[]>([]);
  const [expertMap, setExpertMap] = useState<Record<number, Expert>>({});
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
    totalPages: 0
  });
  const [filters, setFilters] = useState({
    status: '',
    type: '',
    search: ''
  });
  const [selectedExpertId, setSelectedExpertId] = useState<number | null>(null);
  const [isCreateModalOpen, setIsCreateModalOpen] = useState(false);
  const [isViewModalOpen, setIsViewModalOpen] = useState(false);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [selectedEngagement, setSelectedEngagement] = useState<Engagement | null>(null);

  const fetchEngagements = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const params: Record<string, string | boolean> = {};
      if (filters.status) params.status = filters.status;
      if (filters.type) params.type = filters.type;
      if (filters.search) params.search = filters.search;
      
      const response = await engagementApi.getEngagements(pagination.page, pagination.limit, params);
      
      if (response.success) {
        setEngagements(response.data);
        setPagination({
          page: pagination.page,
          limit: pagination.limit,
          total: response.data.length,
          totalPages: Math.ceil(response.data.length / pagination.limit)
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
  }, [pagination.page, pagination.limit, filters]);

  const fetchExperts = useCallback(async () => {
    try {
      const response = await expertsApi.getExperts(1, 100);
      if (response.success) {
        const expertsList = response.data.experts;
        setExperts(expertsList);
        
        // Create a map of expert IDs to expert objects for quick lookup
        const expertMapObj: Record<number, Expert> = {};
        expertsList.forEach(expert => {
          expertMapObj[expert.id] = expert;
        });
        setExpertMap(expertMapObj);
      }
    } catch (error) {
      console.error('Error fetching experts:', error);
    }
  }, []);

  useEffect(() => {
    fetchEngagements();
    fetchExperts();
  }, [fetchEngagements, fetchExperts]);

  const handlePageChange = (newPage: number) => {
    setPagination(prev => ({ ...prev, page: newPage }));
  };

  const handleFilterChange = (e: React.ChangeEvent<HTMLSelectElement | HTMLInputElement>) => {
    const { name, value } = e.target;
    setFilters(prev => ({ ...prev, [name]: value }));
    setPagination(prev => ({ ...prev, page: 1 })); // Reset to first page on filter change
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    fetchEngagements();
  };

  const resetFilters = () => {
    setFilters({
      status: '',
      type: '',
      search: ''
    });
    setPagination(prev => ({ ...prev, page: 1 }));
  };

  const handleEngagementCreated = (_engagement: Engagement) => {
    setIsCreateModalOpen(false);
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

  const handleCreateEngagement = () => {
    // If no expert is selected, show expert selection modal
    if (experts.length > 0) {
      setSelectedExpertId(null);
      setIsCreateModalOpen(true);
    } else {
      setError('No experts available. Please add experts first.');
    }
  };

  const getExpertName = (expertId: number) => {
    return expertMap[expertId]?.name || 'Unknown Expert';
  };

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="mb-6">
        <h1 className="text-2xl font-bold mb-2">Engagement Management</h1>
        <p className="text-gray-600">
          Manage expert engagements, track status, and create new engagements.
        </p>
      </div>
      
      <div className="bg-white rounded-md shadow-sm p-6 mb-6">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-xl font-semibold">Engagements</h2>
          <Button
            variant="primary"
            onClick={handleCreateEngagement}
          >
            Create Engagement
          </Button>
        </div>
        
        {/* Filters */}
        <form onSubmit={handleSearch} className="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select
              name="status"
              value={filters.status}
              onChange={handleFilterChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            >
              <option value="">All Statuses</option>
              <option value="pending">Pending</option>
              <option value="confirmed">Confirmed</option>
              <option value="in_progress">In Progress</option>
              <option value="completed">Completed</option>
              <option value="cancelled">Cancelled</option>
            </select>
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Type</label>
            <select
              name="type"
              value={filters.type}
              onChange={handleFilterChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            >
              <option value="">All Types</option>
              <option value="validator">Validator</option>
              <option value="evaluator">Evaluator</option>
            </select>
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Search</label>
            <Input
              name="search"
              value={filters.search}
              onChange={handleFilterChange}
              placeholder="Search by title, org, or expert..."
            />
          </div>
          
          <div className="flex items-end gap-2">
            <Button
              type="submit"
              variant="secondary"
            >
              Apply Filters
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={resetFilters}
            >
              Reset
            </Button>
          </div>
        </form>
        
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-600 p-4 rounded-md mb-4">
            {error}
          </div>
        )}
        
        {/* Engagements Table */}
        {isLoading ? (
          <div className="text-center p-8">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary mx-auto"></div>
            <p className="mt-4 text-gray-600">Loading engagements...</p>
          </div>
        ) : engagements.length === 0 ? (
          <div className="text-center p-8 bg-gray-50 rounded-md">
            <p className="text-gray-500">No engagements found.</p>
            <p className="text-gray-500 mt-2">Try adjusting your filters or create a new engagement.</p>
          </div>
        ) : (
          <>
            <Table>
              <Table.Header>
                <Table.Row>
                  <Table.HeaderCell>Title</Table.HeaderCell>
                  <Table.HeaderCell>Expert</Table.HeaderCell>
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
                    <Table.Cell>{getExpertName(engagement.expertId)}</Table.Cell>
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
                      </div>
                    </Table.Cell>
                  </Table.Row>
                ))}
              </Table.Body>
            </Table>
            
            {/* Pagination */}
            {pagination.totalPages > 1 && (
              <div className="flex justify-center mt-6">
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
      </div>
      
      {/* Create Engagement Modal - First select an expert */}
      <Modal
        isOpen={isCreateModalOpen && selectedExpertId === null}
        onClose={() => setIsCreateModalOpen(false)}
        title="Select Expert for Engagement"
        size="md"
      >
        <div className="space-y-4">
          <p className="text-gray-600">Select an expert to create an engagement for:</p>
          
          <div className="max-h-96 overflow-y-auto">
            <Table>
              <Table.Header>
                <Table.Row>
                  <Table.HeaderCell>Name</Table.HeaderCell>
                  <Table.HeaderCell>Affiliation</Table.HeaderCell>
                  <Table.HeaderCell>Action</Table.HeaderCell>
                </Table.Row>
              </Table.Header>
              <Table.Body>
                {experts.map(expert => (
                  <Table.Row key={expert.id}>
                    <Table.Cell>{expert.name}</Table.Cell>
                    <Table.Cell>{expert.affiliation}</Table.Cell>
                    <Table.Cell>
                      <Button
                        variant="secondary"
                        size="sm"
                        onClick={() => setSelectedExpertId(expert.id)}
                      >
                        Select
                      </Button>
                    </Table.Cell>
                  </Table.Row>
                ))}
              </Table.Body>
            </Table>
          </div>
          
          <div className="flex justify-end pt-4">
            <Button
              variant="outline"
              onClick={() => setIsCreateModalOpen(false)}
            >
              Cancel
            </Button>
          </div>
        </div>
      </Modal>
      
      {/* Create Engagement Form Modal */}
      <Modal
        isOpen={isCreateModalOpen && selectedExpertId !== null}
        onClose={() => setIsCreateModalOpen(false)}
        title="Create New Engagement"
        size="lg"
      >
        {selectedExpertId && (
          <div>
            <div className="mb-4 pb-4 border-b">
              <p className="text-gray-600">
                Creating engagement for expert: <strong>{expertMap[selectedExpertId]?.name}</strong>
              </p>
            </div>
            <EngagementForm
              expertId={selectedExpertId}
              onSuccess={handleEngagementCreated}
              onCancel={() => setIsCreateModalOpen(false)}
            />
          </div>
        )}
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
              <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-1">Expert</h3>
              <p>{getExpertName(selectedEngagement.expertId)}</p>
            </div>
            
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
            
            <div className="flex justify-end space-x-3 pt-4">
              <Button 
                variant="outline"
                onClick={() => setIsViewModalOpen(false)}
              >
                Close
              </Button>
              <Button 
                variant="secondary"
                onClick={() => {
                  setIsViewModalOpen(false);
                  handleEditEngagement(selectedEngagement);
                }}
              >
                Edit
              </Button>
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
            expertId={selectedEngagement.expertId}
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

export default EngagementManagementPage;