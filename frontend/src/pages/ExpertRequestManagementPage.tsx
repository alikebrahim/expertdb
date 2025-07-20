import React, { useState, useEffect } from 'react';
import { ExpertRequest } from '../types';
import { expertRequestsApi } from '../services/api';
import { useUI } from '../hooks/useUI';
import AdminRequestTable from '../components/AdminRequestTable';
import RequestDetailModal from '../components/RequestDetailModal';
import BatchApproveModal from '../components/BatchApproveModal';
import { Card, CardHeader, CardContent } from '../components/ui/Card';
import Button from '../components/ui/Button';
import { FormField } from '../components/ui/FormField';
import { LoadingSpinner } from '../components/ui/LoadingSpinner';

const ExpertRequestManagementPage: React.FC = () => {
  const [requests, setRequests] = useState<ExpertRequest[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedRequests, setSelectedRequests] = useState<number[]>([]);
  const [currentStatusFilter, setCurrentStatusFilter] = useState<string | null>(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedRequest, setSelectedRequest] = useState<ExpertRequest | null>(null);
  const [showBatchApproveModal, setShowBatchApproveModal] = useState(false);
  const [pagination, setPagination] = useState({
    currentPage: 1,
    totalPages: 1,
    totalCount: 0
  });

  const { showNotification } = useUI();

  const fetchRequests = async (page = 1, status: string | null = null, search = '') => {
    setLoading(true);
    setError(null);
    
    try {
      const params: Record<string, string | boolean> = {
        page: page.toString(),
        limit: '20'
      };
      
      if (status) params.status = status;
      if (search) params.search = search;
      
      const response = await expertRequestsApi.getExpertRequests(20, (page - 1) * 20, params);
      
      if (response.success) {
        setRequests(response.data.requests);
        setPagination({
          currentPage: response.data.pagination.currentPage,
          totalPages: response.data.pagination.totalPages,
          totalCount: response.data.pagination.totalCount
        });
      } else {
        setError(response.message || 'Failed to fetch requests');
      }
    } catch (err) {
      console.error('Error fetching requests:', err);
      setError('Failed to load expert requests');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchRequests(1, currentStatusFilter, searchQuery);
  }, [currentStatusFilter]);

  const handleRequestSelection = (requestId: number, selected: boolean) => {
    if (selected) {
      setSelectedRequests(prev => [...prev, requestId]);
    } else {
      setSelectedRequests(prev => prev.filter(id => id !== requestId));
    }
  };

  const handleStatusFilter = (status: string | null) => {
    setCurrentStatusFilter(status);
    setSelectedRequests([]);
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    fetchRequests(1, currentStatusFilter, searchQuery);
  };

  const handleRequestUpdate = (message: string) => {
    showNotification(message, 'success');
    setSelectedRequest(null);
    fetchRequests(pagination.currentPage, currentStatusFilter, searchQuery);
  };

  const handleBatchApprove = (requestIds: number[]) => {
    setShowBatchApproveModal(true);
  };

  const handleBatchApproveSuccess = (approvedIds: number[]) => {
    showNotification(`Successfully approved ${approvedIds.length} request(s)`, 'success');
    setShowBatchApproveModal(false);
    setSelectedRequests([]);
    fetchRequests(pagination.currentPage, currentStatusFilter, searchQuery);
  };

  const handleBatchReject = (requestIds: number[]) => {
    // For now, just show notification - could implement batch reject modal
    showNotification('Batch reject functionality coming soon', 'info');
  };

  const handlePageChange = (page: number) => {
    fetchRequests(page, currentStatusFilter, searchQuery);
  };

  const getFilteredRequests = () => {
    if (!currentStatusFilter) return requests;
    return requests.filter(request => request.status === currentStatusFilter);
  };

  const pendingCount = requests.filter(r => r.status === 'pending').length;
  const approvedCount = requests.filter(r => r.status === 'approved').length;
  const rejectedCount = requests.filter(r => r.status === 'rejected').length;

  return (
    <div className="space-y-6">
      {/* Page Header */}
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-semibold text-gray-900">Expert Request Management</h1>
          <p className="text-gray-600">Review and manage expert requests</p>
        </div>
        <div className="flex items-center space-x-4">
          <div className="text-sm text-gray-600">
            Total: {pagination.totalCount} requests
          </div>
        </div>
      </div>

      {/* Statistics Cards */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="w-8 h-8 bg-yellow-100 rounded-full flex items-center justify-center">
                  <span className="text-yellow-600 font-medium">P</span>
                </div>
              </div>
              <div className="ml-4">
                <div className="text-sm font-medium text-gray-500">Pending Review</div>
                <div className="text-2xl font-semibold text-gray-900">{pendingCount}</div>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center">
                  <span className="text-green-600 font-medium">A</span>
                </div>
              </div>
              <div className="ml-4">
                <div className="text-sm font-medium text-gray-500">Approved</div>
                <div className="text-2xl font-semibold text-gray-900">{approvedCount}</div>
              </div>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center">
              <div className="flex-shrink-0">
                <div className="w-8 h-8 bg-red-100 rounded-full flex items-center justify-center">
                  <span className="text-red-600 font-medium">R</span>
                </div>
              </div>
              <div className="ml-4">
                <div className="text-sm font-medium text-gray-500">Rejected</div>
                <div className="text-2xl font-semibold text-gray-900">{rejectedCount}</div>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Search and Filters */}
      <Card>
        <CardHeader>
          <h2 className="text-lg font-semibold">Search and Filter</h2>
        </CardHeader>
        <CardContent>
          <form onSubmit={handleSearch} className="flex items-center space-x-4">
            <div className="flex-1">
              <FormField
                name="search"
                placeholder="Search by name, institution, or specialization..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className="w-full"
              />
            </div>
            <Button type="submit" disabled={loading}>
              {loading ? <LoadingSpinner size="sm" /> : 'Search'}
            </Button>
            <Button 
              type="button" 
              variant="outline" 
              onClick={() => {
                setSearchQuery('');
                fetchRequests(1, currentStatusFilter, '');
              }}
            >
              Clear
            </Button>
          </form>
        </CardContent>
      </Card>

      {/* Requests Table */}
      <Card>
        <CardHeader>
          <h2 className="text-lg font-semibold">Expert Requests</h2>
        </CardHeader>
        <CardContent>
          <AdminRequestTable
            requests={getFilteredRequests()}
            isLoading={loading}
            error={error}
            onViewRequest={setSelectedRequest}
            onRequestSelection={handleRequestSelection}
            selectedRequests={selectedRequests}
            onStatusFilter={handleStatusFilter}
            currentStatusFilter={currentStatusFilter}
            onBatchApprove={handleBatchApprove}
            onBatchReject={handleBatchReject}
            pagination={{
              currentPage: pagination.currentPage,
              totalPages: pagination.totalPages,
              onPageChange: handlePageChange
            }}
          />
        </CardContent>
      </Card>

      {/* Request Detail Modal */}
      {selectedRequest && (
        <RequestDetailModal
          request={selectedRequest}
          onClose={() => setSelectedRequest(null)}
          onRequestUpdate={handleRequestUpdate}
        />
      )}

      {/* Batch Approve Modal */}
      {showBatchApproveModal && (
        <BatchApproveModal
          requests={requests.filter(r => selectedRequests.includes(r.id))}
          onClose={() => setShowBatchApproveModal(false)}
          onSuccess={handleBatchApproveSuccess}
        />
      )}
    </div>
  );
};

export default ExpertRequestManagementPage;