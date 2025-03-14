import { useState } from 'react';
import { Expert } from '../types';
import { Table, TableRow, TableCell } from './ui/Table';
import Button from './ui/Button';

interface ExpertTableProps {
  experts: Expert[];
  isLoading: boolean;
  error: string | null;
}

const ExpertTable = ({ experts, isLoading, error }: ExpertTableProps) => {
  const [selectedExpert, setSelectedExpert] = useState<Expert | null>(null);
  
  const headers = [
    'Name',
    'Role',
    'Type',
    'Affiliation',
    'Specialization',
    'ISCED',
    'Actions',
  ];
  
  const handleSelectExpert = (expert: Expert) => {
    setSelectedExpert(expert);
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
        <Table headers={headers}>
          {experts.map((expert) => (
            <TableRow 
              key={expert.id} 
              isClickable 
              onClick={() => handleSelectExpert(expert)}
            >
              <TableCell>{expert.name}</TableCell>
              <TableCell>{expert.role}</TableCell>
              <TableCell>{expert.type}</TableCell>
              <TableCell>{expert.affiliation}</TableCell>
              <TableCell>{expert.specialization}</TableCell>
              <TableCell>{expert.isced}</TableCell>
              <TableCell>
                <Button 
                  variant="outline"
                  size="sm"
                  onClick={(e) => {
                    e.stopPropagation();
                    handleSelectExpert(expert);
                  }}
                >
                  Details
                </Button>
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
                  <h3 className="text-sm font-medium text-neutral-500">Type</h3>
                  <p className="mt-1">{selectedExpert.type}</p>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Affiliation</h3>
                  <p className="mt-1">{selectedExpert.affiliation}</p>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Nationality</h3>
                  <p className="mt-1">{selectedExpert.nationality}</p>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Specialization</h3>
                  <p className="mt-1">{selectedExpert.specialization}</p>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">ISCED</h3>
                  <p className="mt-1">{selectedExpert.isced}</p>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Status</h3>
                  <p className="mt-1">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      selectedExpert.status === 'available' 
                        ? 'bg-green-100 text-green-800'
                        : 'bg-red-100 text-red-800'
                    }`}>
                      {selectedExpert.status === 'available' ? 'Available' : 'Unavailable'}
                    </span>
                  </p>
                </div>
              </div>
              
              {selectedExpert.biography && (
                <div className="mt-6">
                  <h3 className="text-sm font-medium text-neutral-500">Biography</h3>
                  <p className="mt-1 text-neutral-800">{selectedExpert.biography}</p>
                </div>
              )}
              
              <div className="mt-6 flex justify-end space-x-3">
                <Button 
                  variant="outline"
                  onClick={() => setSelectedExpert(null)}
                >
                  Close
                </Button>
              </div>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

export default ExpertTable;