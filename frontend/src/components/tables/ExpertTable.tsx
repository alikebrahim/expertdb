import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Expert } from '../../types';
import { useUI } from '../../hooks/useUI';
import { Table, TableRow, TableCell } from '../ui/Table';
import Button from '../ui/Button';
import Modal from '../Modal';
import { formatDate } from '../../utils/formatters';

interface ExpertTableProps {
  experts: Expert[];
  isLoading: boolean;
  error: string | null;
  pagination?: {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
  };
  onEdit?: (expert: Expert) => void;
  onDelete?: (expert: Expert) => void;
  enableBatchActions?: boolean;
  onBatchExport?: (expertIds: number[]) => void;
  onBatchPublish?: (expertIds: number[]) => void;
  onBatchUnpublish?: (expertIds: number[]) => void;
}

const ExpertTable = ({ 
  experts, 
  isLoading, 
  error, 
  pagination, 
  onEdit, 
  onDelete,
  enableBatchActions = false,
  onBatchExport,
  onBatchPublish,
  onBatchUnpublish
}: ExpertTableProps) => {
  const [selectedExpert, setSelectedExpert] = useState<Expert | null>(null);
  const [selectedExperts, setSelectedExperts] = useState<number[]>([]);
  const [isExportModalOpen, setIsExportModalOpen] = useState(false);
  const navigate = useNavigate();
  const { addNotification } = useUI();
  
  const headers = enableBatchActions 
    ? ['Select', 'Name', 'Role', 'Employment', 'Affiliation', 'Rating', 'Last Updated', 'Actions']
    : ['Name', 'Role', 'Employment', 'Affiliation', 'Contact', 'Rating', 'Actions'];
  
  const handleSelectExpert = (expert: Expert) => {
    setSelectedExpert(expert);
  };
  
  const handleViewExpertDetails = (expert: Expert, e: React.MouseEvent) => {
    e.stopPropagation();
    navigate(`/experts/${expert.id}`);
  };

  const handleSelectAll = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.checked) {
      setSelectedExperts(experts.map(expert => expert.id));
    } else {
      setSelectedExperts([]);
    }
  };

  const handleSelectSingle = (expertId: number, isChecked: boolean) => {
    if (isChecked) {
      setSelectedExperts(prev => [...prev, expertId]);
    } else {
      setSelectedExperts(prev => prev.filter(id => id !== expertId));
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

  return (
    <div className="space-y-4">
      {enableBatchActions && (
        <div className="flex flex-wrap gap-2 mb-4">
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
        <Table headers={headers} pagination={pagination}>
          {enableBatchActions && (
            <tr className="bg-neutral-50">
              <th className="px-4 py-2 text-left">
                <input
                  type="checkbox"
                  onChange={handleSelectAll}
                  checked={selectedExperts.length === experts.length && experts.length > 0}
                  className="h-4 w-4 text-primary focus:ring-primary border-neutral-300 rounded"
                />
              </th>
              <th colSpan={7} className="px-4 py-2 text-left text-sm font-medium text-neutral-500">
                {selectedExperts.length} of {experts.length} selected
              </th>
            </tr>
          )}
          {experts.map((expert) => (
            <TableRow 
              key={expert.id} 
              isClickable 
              onClick={() => handleSelectExpert(expert)}
            >
              {enableBatchActions && (
                <TableCell onClick={(e) => e.stopPropagation()}>
                  <input
                    type="checkbox"
                    checked={selectedExperts.includes(expert.id)}
                    onChange={(e) => handleSelectSingle(expert.id, e.target.checked)}
                    className="h-4 w-4 text-primary focus:ring-primary border-neutral-300 rounded"
                  />
                </TableCell>
              )}
              <TableCell>{expert.name}</TableCell>
              <TableCell>{expert.role}</TableCell>
              <TableCell>{expert.employmentType}</TableCell>
              <TableCell>{expert.affiliation}</TableCell>
              {!enableBatchActions && <TableCell>{expert.primaryContact}</TableCell>}
              <TableCell>
                <div className="flex items-center">
                  <span className="mr-1">{expert.rating}</span>
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-yellow-500" viewBox="0 0 20 20" fill="currentColor">
                    <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                  </svg>
                </div>
              </TableCell>
              {enableBatchActions && <TableCell>{formatDate(expert.updated_at)}</TableCell>}
              <TableCell>
                <div className="flex flex-wrap gap-2">
                  <Button 
                    variant="primary"
                    size="sm"
                    onClick={(e) => handleViewExpertDetails(expert, e)}
                  >
                    View
                  </Button>
                  {onEdit && (
                    <Button 
                      variant="secondary"
                      size="sm"
                      onClick={(e) => {
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
                      onClick={(e) => {
                        e.stopPropagation();
                        onDelete(expert);
                      }}
                    >
                      Delete
                    </Button>
                  )}
                </div>
              </TableCell>
            </TableRow>
          ))}
        </Table>
      )}
      
      {/* Expert details modal */}
      {selectedExpert && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
          <div className="bg-white rounded-lg shadow-lg w-full max-w-2xl overflow-hidden">
            <div className="p-6">
              <div className="flex justify-between items-start">
                <h2 className="text-2xl font-bold text-primary">{selectedExpert.name}</h2>
                <button
                  onClick={() => setSelectedExpert(null)}
                  className="text-neutral-500 hover:text-neutral-700"
                >
                  <span className="sr-only">Close</span>
                  <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
              
              <div className="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Role</h3>
                  <p className="mt-1">{selectedExpert.role}</p>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Employment Type</h3>
                  <p className="mt-1">{selectedExpert.employmentType}</p>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Affiliation</h3>
                  <p className="mt-1">{selectedExpert.affiliation}</p>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Contact</h3>
                  <p className="mt-1">{selectedExpert.primaryContact} ({selectedExpert.contactType})</p>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Skills</h3>
                  <div className="mt-1 flex flex-wrap gap-1">
                    {selectedExpert.skills.map((skill, index) => (
                      <span 
                        key={index}
                        className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-primary bg-opacity-10 text-primary"
                      >
                        {skill}
                      </span>
                    ))}
                  </div>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Rating</h3>
                  <div className="mt-1 flex items-center">
                    <span className="mr-1">{selectedExpert.rating}</span>
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-yellow-500" viewBox="0 0 20 20" fill="currentColor">
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                  </div>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Availability</h3>
                  <p className="mt-1">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      selectedExpert.availability === 'Available' 
                        ? 'bg-green-100 text-green-800'
                        : selectedExpert.availability === 'Limited'
                        ? 'bg-yellow-100 text-yellow-800'
                        : 'bg-red-100 text-red-800'
                    }`}>
                      {selectedExpert.availability}
                    </span>
                  </p>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Origin</h3>
                  <p className="mt-1">
                    {selectedExpert.isBahraini ? 'Bahraini' : 'International'}
                  </p>
                </div>
              </div>
              
              {selectedExpert.biography && (
                <div className="mt-6">
                  <h3 className="text-sm font-medium text-neutral-500">Biography</h3>
                  <p className="mt-1 text-neutral-800">{selectedExpert.biography}</p>
                </div>
              )}
              
              <div className="mt-4 text-sm text-neutral-500">
                <p>Created: {formatDate(selectedExpert.created_at)}</p>
                <p>Last updated: {formatDate(selectedExpert.updated_at)}</p>
              </div>
              
              <div className="mt-6 flex justify-end space-x-3">
                <Button 
                  variant="outline"
                  onClick={() => setSelectedExpert(null)}
                >
                  Close
                </Button>
                <Button 
                  variant="primary"
                  onClick={(e) => {
                    setSelectedExpert(null);
                    handleViewExpertDetails(selectedExpert, e as React.MouseEvent<HTMLButtonElement>);
                  }}
                >
                  View Full Profile
                </Button>
              </div>
            </div>
          </div>
        </div>
      )}

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