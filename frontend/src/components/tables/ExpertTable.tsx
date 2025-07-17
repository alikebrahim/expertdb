import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Expert } from '../../types';
import { useUI } from '../../hooks/useUI';
import { Table, Button, ColumnSelector, DEFAULT_COLUMNS } from '../ui';
import Modal from '../Modal';
import ExpertProfileModal from '../ExpertProfileModal';
import { formatDate } from '../../utils/formatters';
import type { ColumnConfig } from '../ui/ColumnSelector';

export interface SortConfig {
  field: string;
  direction: 'asc' | 'desc';
}

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

interface ExpertTableProps {
  experts: Expert[];
  isLoading: boolean;
  error: string | null;
  onEdit?: (expert: Expert) => void;
  onDelete?: (expert: Expert) => void;
  enableBatchActions?: boolean;
  onBatchExport?: (expertIds: number[]) => void;
  onBatchPublish?: (expertIds: number[]) => void;
  onBatchUnpublish?: (expertIds: number[]) => void;
  pagination?: PaginationProps;
  sortConfig?: SortConfig;
  onSort?: (field: string) => void;
  showColumnSelector?: boolean;
}

const ExpertTable = ({ 
  experts, 
  isLoading = false,
  error = null,
  onEdit, 
  onDelete,
  enableBatchActions = false,
  onBatchExport,
  onBatchPublish,
  onBatchUnpublish,
  pagination,
  sortConfig,
  onSort,
  showColumnSelector = false
}: ExpertTableProps) => {
  const [selectedExpert, setSelectedExpert] = useState<Expert | null>(null);
  const [selectedExperts, setSelectedExperts] = useState<number[]>([]);
  const [isExportModalOpen, setIsExportModalOpen] = useState(false);
  const [columns, setColumns] = useState<ColumnConfig[]>(DEFAULT_COLUMNS);
  const navigate = useNavigate();
  const { addNotification } = useUI();
  
  const handleSelectExpert = (expert: Expert) => {
    setSelectedExpert(expert);
  };
  
  const handleViewExpertDetails = (expert: Expert, e: React.MouseEvent) => {
    e.stopPropagation();
    navigate(`/experts/${expert.id}`);
  };

  const handleSelectAll = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.checked) {
      setSelectedExperts(experts.map((expert: Expert) => expert.id));
    } else {
      setSelectedExperts([]);
    }
  };

  const handleSelectSingle = (expertId: number, isChecked: boolean) => {
    if (isChecked) {
      setSelectedExperts((prev: number[]) => [...prev, expertId]);
    } else {
      setSelectedExperts((prev: number[]) => prev.filter((id: number) => id !== expertId));
    }
  };

  const handleBatchExport = () => {
    if (selectedExperts.length === 0) {
      addNotification({
        type: 'warning',
        message: 'Please select at least one expert to export',
        duration: 3000,
      });
      return;
    }

    if (onBatchExport) {
      onBatchExport(selectedExperts);
    } else {
      setIsExportModalOpen(true);
    }
  };

  const handleBatchPublish = () => {
    if (selectedExperts.length === 0) {
      addNotification({
        type: 'warning',
        message: 'Please select at least one expert to publish',
        duration: 3000,
      });
      return;
    }

    if (onBatchPublish) {
      onBatchPublish(selectedExperts);
    } else {
      addNotification({
        type: 'info',
        message: 'Publish functionality not implemented yet',
        duration: 3000,
      });
    }
  };

  const handleBatchUnpublish = () => {
    if (selectedExperts.length === 0) {
      addNotification({
        type: 'warning',
        message: 'Please select at least one expert to unpublish',
        duration: 3000,
      });
      return;
    }

    if (onBatchUnpublish) {
      onBatchUnpublish(selectedExperts);
    } else {
      addNotification({
        type: 'info',
        message: 'Unpublish functionality not implemented yet',
        duration: 3000,
      });
    }
  };
  
  if (error) {
    return (
      <div className="bg-secondary bg-opacity-10 text-secondary p-4 rounded">
        <p>Error loading experts: {error}</p>
      </div>
    );
  }

  const visibleColumns = columns.filter(col => col.visible);
  
  const getCellValue = (expert: Expert, columnKey: string) => {
    switch (columnKey) {
      case 'id':
        return expert.id;
      case 'name':
        return expert.name;
      case 'affiliation':
        return expert.affiliation;
      case 'specializedArea':
        return expert.specializedArea;
      case 'rating':
        return expert.rating;
      case 'role':
        return expert.role;
      case 'employmentType':
        return expert.employmentType;
      case 'generalArea':
        return expert.generalAreaName || expert.generalArea;
      case 'phone':
        return expert.phone || 'N/A';
      case 'email':
        return expert.email || 'N/A';
      case 'cvPath':
        return expert.cvPath;
      case 'isAvailable':
        return expert.isAvailable;
      case 'nationality':
        return 'N/A'; // Nationality not available in current Expert type
      case 'isTrained':
        return expert.isTrained;
      case 'created_at':
        return formatDate(expert.created_at);
      case 'actions':
        return null; // Actions column is handled separately
      default:
        return expert[columnKey as keyof Expert];
    }
  };

  const renderCellContent = (expert: Expert, column: ColumnConfig) => {
    if (column.key === 'actions') {
      return (
        <div className="flex flex-wrap gap-2">
          <Button 
            variant="primary"
            size="sm"
            onClick={(e: React.MouseEvent) => handleViewExpertDetails(expert, e)}
          >
            View
          </Button>
          {onEdit && (
            <Button 
              variant="secondary"
              size="sm"
              onClick={(e: React.MouseEvent) => {
                e.stopPropagation();
                onEdit(expert);
              }}
            >
              Edit
            </Button>
          )}
          {onDelete && (
            <Button 
              variant="danger"
              size="sm"
              onClick={(e: React.MouseEvent) => {
                e.stopPropagation();
                onDelete(expert);
              }}
            >
              Delete
            </Button>
          )}
        </div>
      );
    }
    
    if (column.key === 'rating') {
      const rating = expert.rating;
      const displayRating = rating === 0 ? 'N/A' : rating.toString();
      return (
        <div className="flex items-center">
          <span className="mr-1">{displayRating}</span>
          {rating > 0 && (
            <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-yellow-500" viewBox="0 0 20 20" fill="currentColor">
              <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
            </svg>
          )}
        </div>
      );
    }
    
    if (column.key === 'cvPath' && column.type === 'file') {
      const cvPath = getCellValue(expert, column.key);
      if (cvPath) {
        return (
          <button
            className="text-blue-600 hover:text-blue-800 text-sm underline"
            onClick={(e) => {
              e.stopPropagation();
              const link = document.createElement('a');
              link.href = `/api/documents/download/${cvPath}`;
              link.download = `${expert.name}_CV.pdf`;
              link.click();
            }}
          >
            Download
          </button>
        );
      }
      return <span className="text-gray-400 text-sm">No CV</span>;
    }
    
    if (column.type === 'status') {
      const value = getCellValue(expert, column.key);
      return (
        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
          value ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
        }`}>
          {value ? 'Yes' : 'No'}
        </span>
      );
    }
    
    return getCellValue(expert, column.key);
  };

  return (
    <div className="space-y-4">
      {/* Header with controls */}
      <div className="flex justify-between items-center">
        <div className="flex items-center gap-4">
          {showColumnSelector && (
            <ColumnSelector 
              columns={columns}
              onColumnChange={setColumns}
            />
          )}
          {enableBatchActions && (
            <div className="flex flex-wrap gap-2">
              <Button 
                variant="outline" 
                size="sm"
                onClick={handleBatchExport}
                disabled={selectedExperts.length === 0}
              >
                Export Selected ({selectedExperts.length})
              </Button>
              <Button 
                variant="primary" 
                size="sm"
                onClick={handleBatchPublish}
                disabled={selectedExperts.length === 0}
              >
                Publish Selected
              </Button>
              <Button 
                variant="secondary" 
                size="sm"
                onClick={handleBatchUnpublish}
                disabled={selectedExperts.length === 0}
              >
                Unpublish Selected
              </Button>
            </div>
          )}
        </div>
        
        <div className="flex items-center gap-4">
          {/* Right side content can go here if needed */}
        </div>
      </div>
      
      {isLoading ? (
        <div className="flex justify-center py-8">
          <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full"></div>
          <span className="sr-only">Loading...</span>
        </div>
      ) : experts.length === 0 ? (
        <div className="bg-accent p-6 rounded text-center">
          <p className="text-neutral-600">No experts found matching your filters.</p>
          <p className="text-sm text-neutral-500 mt-1">Try adjusting your search criteria.</p>
        </div>
      ) : (
        <Table className="w-auto">
          <Table.Header>
            <Table.Row className="bg-neutral-50">
              {enableBatchActions && (
                <Table.HeaderCell className="whitespace-nowrap">
                  <input
                    type="checkbox"
                    onChange={handleSelectAll}
                    checked={selectedExperts.length === experts.length && experts.length > 0}
                    className="h-4 w-4 text-primary focus:ring-primary border-neutral-300 rounded"
                  />
                </Table.HeaderCell>
              )}
              {visibleColumns.map((column) => (
                <Table.HeaderCell key={column.key} className="whitespace-nowrap">
                  {column.sortable && onSort ? (
                    <button
                      onClick={() => onSort(column.key)}
                      className="flex items-center space-x-1 hover:text-primary font-medium whitespace-nowrap"
                    >
                      <span>{column.label}</span>
                      {sortConfig?.field === column.key ? (
                        <span className="ml-1">
                          {sortConfig.direction === 'asc' ? (
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
                            </svg>
                          ) : (
                            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                            </svg>
                          )}
                        </span>
                      ) : (
                        <svg className="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 9l4-4 4 4m0 6l-4 4-4-4" />
                        </svg>
                      )}
                    </button>
                  ) : (
                    column.label
                  )}
                </Table.HeaderCell>
              ))}
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {experts.map((expert) => (
              <Table.Row 
                key={expert.id} 
                className="hover:bg-gray-50 cursor-pointer"
                onClick={() => handleSelectExpert(expert)}
              >
                {enableBatchActions && (
                  <Table.Cell onClick={(e) => e.stopPropagation()} className="whitespace-nowrap">
                    <input
                      type="checkbox"
                      checked={selectedExperts.includes(expert.id)}
                      onChange={(e) => handleSelectSingle(expert.id, e.target.checked)}
                      className="h-4 w-4 text-primary focus:ring-primary border-neutral-300 rounded"
                    />
                  </Table.Cell>
                )}
                {visibleColumns.map((column) => (
                  <Table.Cell key={column.key} className="whitespace-nowrap">
                    {renderCellContent(expert, column)}
                  </Table.Cell>
                ))}
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      )}
      
      {/* Pagination controls */}
      {pagination && pagination.totalPages > 1 && (
        <div className="px-6 py-3 border-t border-gray-200 bg-gray-50 flex items-center justify-between">
          <div className="text-sm text-gray-700">
            Page {pagination.currentPage} of {pagination.totalPages}
          </div>
          <div className="flex space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => pagination.onPageChange(pagination.currentPage - 1)}
              disabled={pagination.currentPage <= 1}
            >
              Previous
            </Button>
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
      
      {/* Expert details modal */}
      <ExpertProfileModal 
        expert={selectedExpert}
        isOpen={!!selectedExpert}
        onClose={() => setSelectedExpert(null)}
      />

      {/* Export Modal */}
      <Modal
        isOpen={isExportModalOpen}
        onClose={() => setIsExportModalOpen(false)}
        title="Export Experts"
      >
        <div className="p-4">
          <p className="mb-4">Select the export format:</p>
          <div className="space-y-3">
            <button
              className="w-full bg-white border border-gray-300 rounded-md p-3 flex items-center justify-between hover:bg-gray-50"
              onClick={() => {
                addNotification({
                  type: 'success',
                  message: `Exported ${selectedExperts.length} experts as Excel`,
                  duration: 3000,
                });
                setIsExportModalOpen(false);
              }}
            >
              <span className="font-medium">Excel (.xlsx)</span>
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-green-600" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
            </button>
            
            <button
              className="w-full bg-white border border-gray-300 rounded-md p-3 flex items-center justify-between hover:bg-gray-50"
              onClick={() => {
                addNotification({
                  type: 'success',
                  message: `Exported ${selectedExperts.length} experts as CSV`,
                  duration: 3000,
                });
                setIsExportModalOpen(false);
              }}
            >
              <span className="font-medium">CSV (.csv)</span>
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
            </button>
            
            <button
              className="w-full bg-white border border-gray-300 rounded-md p-3 flex items-center justify-between hover:bg-gray-50"
              onClick={() => {
                addNotification({
                  type: 'success',
                  message: `Exported ${selectedExperts.length} experts as PDF`,
                  duration: 3000,
                });
                setIsExportModalOpen(false);
              }}
            >
              <span className="font-medium">PDF (.pdf)</span>
              <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
                <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd" />
              </svg>
            </button>
          </div>
          
          <div className="mt-6 flex justify-end space-x-3">
            <Button
              variant="outline"
              onClick={() => setIsExportModalOpen(false)}
            >
              Cancel
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
};

export default ExpertTable;