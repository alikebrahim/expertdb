import { useState } from 'react';
import { ExpertRequest } from '../types';
import { Table } from './ui';
import Button from './ui/Button';

interface AdminRequestTableProps {
  requests: ExpertRequest[];
  isLoading: boolean;
  error: string | null;
  onViewRequest: (request: ExpertRequest) => void;
  onRequestSelection: (requestId: number, selected: boolean) => void;
  selectedRequests: number[];
  pagination?: {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
  };
}

const AdminRequestTable = ({ 
  requests, 
  isLoading, 
  error, 
  onViewRequest,
  onRequestSelection,
  selectedRequests,
  pagination
}: AdminRequestTableProps) => {
  
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };
  
  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'approved':
        return 'bg-green-100 text-green-800 border-green-200';
      case 'rejected':
        return 'bg-red-100 text-red-800 border-red-200';
      case 'pending':
        return 'bg-yellow-100 text-yellow-800 border-yellow-200';
      default:
        return 'bg-gray-100 text-gray-800 border-gray-200';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'approved':
        return '✅';
      case 'rejected':
        return '❌';
      case 'pending':
        return '🟡';
      default:
        return '⚪';
    }
  };

  const handleSelectAll = (checked: boolean) => {
    if (checked) {
      requests.forEach(request => {
        if (!selectedRequests.includes(request.id)) {
          onRequestSelection(request.id, true);
        }
      });
    } else {
      requests.forEach(request => {
        if (selectedRequests.includes(request.id)) {
          onRequestSelection(request.id, false);
        }
      });
    }
  };

  const isAllSelected = requests.length > 0 && requests.every(request => selectedRequests.includes(request.id));
  const isSomeSelected = requests.some(request => selectedRequests.includes(request.id));
  
  if (error) {
    return (
      <div className="bg-red-50 text-red-800 p-4 rounded border border-red-200">
        <p>Error loading requests: {error}</p>
      </div>
    );
  }
  
  return (
    <div className="space-y-4">
      {isLoading ? (
        <div className="flex justify-center py-12">
          <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full"></div>
          <span className="ml-3 text-gray-600">Loading expert requests...</span>
        </div>
      ) : requests.length === 0 ? (
        <div className="bg-gray-50 p-8 rounded text-center">
          <div className="text-gray-400 text-4xl mb-4">📋</div>
          <p className="text-gray-600 font-medium">No expert requests found</p>
          <p className="text-sm text-gray-500 mt-1">
            No requests match your current filters. Try adjusting your search criteria.
          </p>
        </div>
      ) : (
        <div className="overflow-hidden">
          <Table>
            <Table.Header>
              <Table.Row>
                <Table.HeaderCell className="w-12">
                  <input
                    type="checkbox"
                    checked={isAllSelected}
                    ref={input => {
                      if (input) input.indeterminate = isSomeSelected && !isAllSelected;
                    }}
                    onChange={(e) => handleSelectAll(e.target.checked)}
                    className="rounded border-gray-300 text-primary focus:ring-primary"
                  />
                </Table.HeaderCell>
                <Table.HeaderCell>Expert Name</Table.HeaderCell>
                <Table.HeaderCell>Institution</Table.HeaderCell>
                <Table.HeaderCell>Specialization</Table.HeaderCell>
                <Table.HeaderCell>Status</Table.HeaderCell>
                <Table.HeaderCell>Submitted</Table.HeaderCell>
                <Table.HeaderCell>Actions</Table.HeaderCell>
              </Table.Row>
            </Table.Header>
            <Table.Body>
              {requests.map((request) => (
                <Table.Row key={request.id} className="hover:bg-gray-50">
                  <Table.Cell>
                    <input
                      type="checkbox"
                      checked={selectedRequests.includes(request.id)}
                      onChange={(e) => onRequestSelection(request.id, e.target.checked)}
                      className="rounded border-gray-300 text-primary focus:ring-primary"
                    />
                  </Table.Cell>
                  <Table.Cell>
                    <div>
                      <div className="font-medium text-gray-900">{request.name}</div>
                      <div className="text-sm text-gray-500">{request.designation}</div>
                    </div>
                  </Table.Cell>
                  <Table.Cell>
                    <div className="text-sm text-gray-900">{request.institution}</div>
                  </Table.Cell>
                  <Table.Cell>
                    <div className="text-sm text-gray-900">{request.specializedArea}</div>
                    <div className="text-xs text-gray-500">{request.role}</div>
                  </Table.Cell>
                  <Table.Cell>
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium border ${getStatusBadgeClass(request.status)}`}>
                      <span className="mr-1">{getStatusIcon(request.status)}</span>
                      {request.status.charAt(0).toUpperCase() + request.status.slice(1)}
                    </span>
                  </Table.Cell>
                  <Table.Cell>
                    <div className="text-sm text-gray-900">{formatDate(request.createdAt)}</div>
                    <div className="text-xs text-gray-500">
                      {request.createdBy ? `by ${request.createdBy}` : 'Unknown user'}
                    </div>
                  </Table.Cell>
                  <Table.Cell>
                    <div className="flex items-center space-x-2">
                      <Button 
                        variant="outline"
                        size="sm"
                        onClick={() => onViewRequest(request)}
                      >
                        View Details
                      </Button>
                      
                      {request.status === 'pending' && (
                        <div className="flex space-x-1">
                          <Button 
                            variant="primary"
                            size="sm"
                            onClick={() => onViewRequest(request)}
                            className="bg-green-600 hover:bg-green-700"
                          >
                            ✓ Quick Approve
                          </Button>
                        </div>
                      )}
                    </div>
                  </Table.Cell>
                </Table.Row>
              ))}
            </Table.Body>
          </Table>

          {/* Pagination */}
          {pagination && pagination.totalPages > 1 && (
            <div className="flex items-center justify-between mt-6 px-4 py-3 bg-gray-50 border-t border-gray-200">
              <div className="text-sm text-gray-700">
                Page {pagination.currentPage} of {pagination.totalPages}
              </div>
              <div className="flex space-x-1">
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => pagination.onPageChange(pagination.currentPage - 1)}
                  disabled={pagination.currentPage <= 1}
                >
                  Previous
                </Button>
                
                {/* Page numbers */}
                {Array.from({ length: Math.min(5, pagination.totalPages) }, (_, i) => {
                  const pageNum = i + 1;
                  return (
                    <Button
                      key={pageNum}
                      variant={pagination.currentPage === pageNum ? 'primary' : 'outline'}
                      size="sm"
                      onClick={() => pagination.onPageChange(pageNum)}
                    >
                      {pageNum}
                    </Button>
                  );
                })}
                
                <Button
                  variant="outline"
                  size="sm"
                  onClick={() => pagination.onPageChange(pagination.currentPage + 1)}
                  disabled={pagination.currentPage >= pagination.totalPages}
                >
                  Next
                </Button>
              </div>
            </div>
          )}
        </div>
      )}
    </div>
  );
};

export default AdminRequestTable;
