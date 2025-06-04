import { useState, useEffect } from 'react';
import { ExpertRequest } from '../types';
import { expertRequestsApi, expertAreasApi } from '../services/api';
import { useAuth } from '../hooks/useAuth';
import Button from '../components/ui/Button';
import { Card } from '../components/ui/Card';
import { Alert } from '../components/ui/Alert';
import AdminRequestTable from '../components/AdminRequestTable';
import RequestDetailModal from '../components/RequestDetailModal';

type StatusFilter = 'all' | 'pending' | 'approved' | 'rejected';

interface ExpertArea {
  id: number;
  name: string;
}

const AdminPanelPage = () => {
  const { user } = useAuth();
  const [requests, setRequests] = useState<ExpertRequest[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [statusFilter, setStatusFilter] = useState<StatusFilter>('pending');
  const [searchTerm, setSearchTerm] = useState('');
  const [institutionFilter, setInstitutionFilter] = useState('');
  const [areaFilter, setAreaFilter] = useState('');
  const [selectedRequest, setSelectedRequest] = useState<ExpertRequest | null>(null);
  const [showModal, setShowModal] = useState(false);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [limit] = useState(10);
  const [expertAreas, setExpertAreas] = useState<ExpertArea[]>([]);
  const [selectedRequests, setSelectedRequests] = useState<number[]>([]);
  const [showBatchModal, setShowBatchModal] = useState(false);
  const [batchAction, setBatchAction] = useState<'approve' | 'reject' | null>(null);
  const [batchComment, setBatchComment] = useState('');
  const [batchFile, setBatchFile] = useState<File | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  // Load expert areas for filtering
  useEffect(() => {
    const fetchExpertAreas = async () => {
      try {
        const response = await expertAreasApi.getExpertAreas();
        if (response.success && response.data) {
          setExpertAreas(response.data);
        }
      } catch (error) {
        console.error('Failed to fetch expert areas:', error);
      }
    };
    
    fetchExpertAreas();
  }, []);

  // Fetch expert requests with filters
  useEffect(() => {
    const fetchRequests = async () => {
      if (!user || (user.role !== 'admin' && user.role !== 'super_user')) {
        setError('Access denied. Admin privileges required.');
        setIsLoading(false);
        return;
      }
      
      setIsLoading(true);
      setError(null);
      
      try {
        const filters: any = {};
        if (statusFilter !== 'all') {
          filters.status = statusFilter;
        }
        if (searchTerm) {
          filters.search = searchTerm;
        }
        if (institutionFilter) {
          filters.institution = institutionFilter;
        }
        if (areaFilter) {
          filters.generalArea = areaFilter;
        }

        const response = await expertRequestsApi.getExpertRequests(limit, (page - 1) * limit, filters);
        
        if (response.success && response.data) {
          const data = response.data as any;
          setRequests(data.requests || []);
          setTotalPages(data.pagination?.totalPages || 1);
        } else {
          setError(response.message || 'Failed to fetch expert requests');
          setRequests([]);
          setTotalPages(1);
        }
      } catch (error) {
        console.error('Error fetching expert requests:', error);
        setError('An error occurred while fetching requests');
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchRequests();
  }, [user, page, limit, statusFilter, searchTerm, institutionFilter, areaFilter]);

  const handleStatusFilterChange = (status: StatusFilter) => {
    setStatusFilter(status);
    setPage(1); // Reset to first page when changing filters
  };

  const handleSearch = (term: string) => {
    setSearchTerm(term);
    setPage(1);
  };

  const handleViewRequest = (request: ExpertRequest) => {
    setSelectedRequest(request);
    setShowModal(true);
  };

  const handleCloseModal = () => {
    setShowModal(false);
    setSelectedRequest(null);
  };

  const handleRequestUpdate = (message: string) => {
    setSuccessMessage(message);
    setShowModal(false);
    setSelectedRequest(null);
    
    // Refresh the requests list
    setIsLoading(true);
    const refreshRequests = async () => {
      try {
        const filters: any = {};
        if (statusFilter !== 'all') {
          filters.status = statusFilter;
        }
        if (searchTerm) {
          filters.search = searchTerm;
        }
        if (institutionFilter) {
          filters.institution = institutionFilter;
        }
        if (areaFilter) {
          filters.generalArea = areaFilter;
        }

        const response = await expertRequestsApi.getExpertRequests(limit, (page - 1) * limit, filters);
        
        if (response.success && response.data) {
          const data = response.data as any;
          setRequests(data.requests || []);
          setTotalPages(data.pagination?.totalPages || 1);
        }
      } catch (error) {
        console.error('Error refreshing requests:', error);
      } finally {
        setIsLoading(false);
      }
    };
    
    refreshRequests();
  };

  const handleBatchApprove = async () => {
    if (selectedRequests.length === 0) {
      setError('Please select requests to approve');
      return;
    }

    setBatchAction('approve');
    setShowBatchModal(true);
  };

  const handleBatchReject = async () => {
    if (selectedRequests.length === 0) {
      setError('Please select requests to reject');
      return;
    }

    setBatchAction('reject');
    setShowBatchModal(true);
  };

  const handleBatchSubmit = async () => {
    if (!batchAction) return;

    if (batchAction === 'approve' && !batchFile) {
      setError('Approval document is required for batch approval');
      return;
    }

    setIsLoading(true);
    try {
      const results = await Promise.allSettled(
        selectedRequests.map(requestId => {
          if (batchAction === 'approve') {
            return expertRequestsApi.approveExpertRequest(requestId, batchFile!, batchComment);
          } else {
            return expertRequestsApi.rejectExpertRequest(requestId, batchComment);
          }
        })
      );

      const successful = results.filter(result => result.status === 'fulfilled').length;
      const failed = results.filter(result => result.status === 'rejected').length;

      if (successful > 0) {
        setSuccessMessage(`Successfully ${batchAction === 'approve' ? 'approved' : 'rejected'} ${successful} request(s)`);
      }
      if (failed > 0) {
        setError(`Failed to process ${failed} request(s)`);
      }

      // Reset batch state
      setSelectedRequests([]);
      setBatchAction(null);
      setBatchComment('');
      setBatchFile(null);
      setShowBatchModal(false);

      // Refresh requests
      handleRequestUpdate('Batch operation completed');
    } catch (error) {
      console.error('Error in batch operation:', error);
      setError('An error occurred during batch operation');
    } finally {
      setIsLoading(false);
    }
  };

  const handlePageChange = (newPage: number) => {
    setPage(newPage);
  };

  const handleRequestSelection = (requestId: number, selected: boolean) => {
    if (selected) {
      setSelectedRequests(prev => [...prev, requestId]);
    } else {
      setSelectedRequests(prev => prev.filter(id => id !== requestId));
    }
  };

  const clearFilters = () => {
    setSearchTerm('');
    setInstitutionFilter('');
    setAreaFilter('');
    setStatusFilter('all');
    setPage(1);
  };

  // Check if user has admin access
  if (!user || (user.role !== 'admin' && user.role !== 'super_user')) {
    return (
      <div className="p-6">
        <Alert variant="error">
          Access denied. Administrator privileges required to view this page.
        </Alert>
      </div>
    );
  }

  return (
    <div className="p-6">
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-primary">Expert Request Management</h1>
        <p className="text-neutral-600">
          Review and manage expert profile submissions from staff members
        </p>
      </div>

      {successMessage && (
        <Alert variant="success" className="mb-6" onClose={() => setSuccessMessage(null)}>
          {successMessage}
        </Alert>
      )}

      {error && (
        <Alert variant="error" className="mb-6" onClose={() => setError(null)}>
          {error}
        </Alert>
      )}

      {/* Status Filter Tabs */}
      <Card className="mb-6">
        <div className="flex flex-wrap gap-2 mb-4">
          {(['all', 'pending', 'approved', 'rejected'] as StatusFilter[]).map((status) => (
            <Button
              key={status}
              variant={statusFilter === status ? 'primary' : 'outline'}
              size="sm"
              onClick={() => handleStatusFilterChange(status)}
            >
              {status.charAt(0).toUpperCase() + status.slice(1)}
              {status === 'pending' && requests.filter(r => r.status === 'pending').length > 0 && (
                <span className="ml-2 bg-red-500 text-white text-xs rounded-full px-2 py-1">
                  {requests.filter(r => r.status === 'pending').length}
                </span>
              )}
            </Button>
          ))}
        </div>

        {/* Search and Filter Controls */}
        <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <input
              type="text"
              placeholder="Search by name or email..."
              value={searchTerm}
              onChange={(e) => handleSearch(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            />
          </div>
          
          <div>
            <input
              type="text"
              placeholder="Filter by institution..."
              value={institutionFilter}
              onChange={(e) => setInstitutionFilter(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            />
          </div>
          
          <div>
            <select
              value={areaFilter}
              onChange={(e) => setAreaFilter(e.target.value)}
              className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            >
              <option value="">All areas</option>
              {expertAreas.map(area => (
                <option key={area.id} value={area.id.toString()}>
                  {area.name}
                </option>
              ))}
            </select>
          </div>
          
          <div>
            <Button variant="outline" onClick={clearFilters} size="sm" className="w-full">
              Clear Filters
            </Button>
          </div>
        </div>
      </Card>

      {/* Batch Actions */}
      {selectedRequests.length > 0 && (
        <Card className="mb-6">
          <div className="flex items-center justify-between">
            <span className="text-sm text-gray-600">
              {selectedRequests.length} request(s) selected
            </span>
            <div className="flex gap-2">
              <Button 
                variant="primary" 
                size="sm" 
                onClick={handleBatchApprove}
                disabled={selectedRequests.length === 0}
              >
                Batch Approve
              </Button>
              <Button 
                variant="outline" 
                size="sm" 
                onClick={handleBatchReject}
                disabled={selectedRequests.length === 0}
              >
                Batch Reject
              </Button>
              <Button 
                variant="outline" 
                size="sm" 
                onClick={() => setSelectedRequests([])}
              >
                Clear Selection
              </Button>
            </div>
          </div>
        </Card>
      )}

      {/* Main Table */}
      <Card>
        <AdminRequestTable
          requests={requests}
          isLoading={isLoading}
          error={error}
          onViewRequest={handleViewRequest}
          onRequestSelection={handleRequestSelection}
          selectedRequests={selectedRequests}
          pagination={{
            currentPage: page,
            totalPages: totalPages,
            onPageChange: handlePageChange
          }}
        />
      </Card>

      {/* Request Detail Modal */}
      {showModal && selectedRequest && (
        <RequestDetailModal
          request={selectedRequest}
          onClose={handleCloseModal}
          onRequestUpdate={handleRequestUpdate}
        />
      )}

      {/* Batch Action Modal */}
      {showBatchModal && batchAction && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-lg shadow-xl max-w-md w-full mx-4">
            <div className="p-6">
              <h3 className="text-lg font-semibold text-gray-900 mb-4">
                Batch {batchAction === 'approve' ? 'Approve' : 'Reject'} Requests
              </h3>
              <p className="text-gray-600 mb-4">
                You are about to {batchAction} {selectedRequests.length} request(s).
              </p>
              
              {batchAction === 'approve' && (
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-2">
                    Approval Document <span className="text-red-500">*</span>
                  </label>
                  <input
                    type="file"
                    accept=".pdf"
                    onChange={(e) => setBatchFile(e.target.files?.[0] || null)}
                    className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
                    required
                  />
                </div>
              )}
              
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  {batchAction === 'approve' ? 'Comment (optional)' : 'Rejection Reason *'}
                </label>
                <textarea
                  value={batchComment}
                  onChange={(e) => setBatchComment(e.target.value)}
                  placeholder={batchAction === 'approve' ? 'Add a comment...' : 'Please provide a reason for rejection...'}
                  className="w-full px-3 py-2 border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
                  rows={3}
                  required={batchAction === 'reject'}
                />
              </div>
              
              <div className="flex gap-3 justify-end">
                <Button
                  variant="outline"
                  onClick={() => {
                    setShowBatchModal(false);
                    setBatchAction(null);
                    setBatchComment('');
                    setBatchFile(null);
                  }}
                >
                  Cancel
                </Button>
                <Button
                  variant={batchAction === 'approve' ? 'primary' : 'outline'}
                  onClick={handleBatchSubmit}
                  disabled={isLoading || (batchAction === 'approve' && !batchFile) || (batchAction === 'reject' && !batchComment.trim())}
                >
                  {isLoading ? 'Processing...' : `${batchAction === 'approve' ? 'Approve' : 'Reject'} Requests`}
                </Button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminPanelPage;
