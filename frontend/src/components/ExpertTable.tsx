import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Expert } from '../types';
import { Table, Button } from './ui';

interface ExpertTableProps {
  experts: Expert[];
  isLoading: boolean;
  error: string | null;
  onEdit?: (expert: Expert) => void;
  onDelete?: (expert: Expert) => void;
}

const ExpertTable = ({ experts, isLoading, error, onEdit, onDelete }: ExpertTableProps) => {
  const [selectedExpert, setSelectedExpert] = useState<Expert | null>(null);
  const navigate = useNavigate();
  
  const handleSelectExpert = (expert: Expert) => {
    setSelectedExpert(expert);
  };
  
  const handleViewExpertDetails = (expert: Expert, e: React.MouseEvent) => {
    e.stopPropagation();
    navigate(`/experts/${expert.id}`);
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
        <Table>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Name</Table.HeaderCell>
              <Table.HeaderCell>Role</Table.HeaderCell>
              <Table.HeaderCell>Employment</Table.HeaderCell>
              <Table.HeaderCell>Affiliation</Table.HeaderCell>
              <Table.HeaderCell>Contact</Table.HeaderCell>
              <Table.HeaderCell>Rating</Table.HeaderCell>
              <Table.HeaderCell>Actions</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {experts.map((expert) => (
              <Table.Row 
                key={expert.id} 
                onClick={() => handleSelectExpert(expert)}
              >
                <Table.Cell>{expert.name}</Table.Cell>
                <Table.Cell>{expert.role}</Table.Cell>
                <Table.Cell>{expert.employmentType}</Table.Cell>
                <Table.Cell>{expert.institution}</Table.Cell>
                <Table.Cell>{expert.phone}</Table.Cell>
                <Table.Cell>{expert.rating}</Table.Cell>
                <Table.Cell>
                  <div className="flex space-x-2">
                    <Button 
                      variant="outline"
                      size="sm"
                      onClick={(e: React.MouseEvent) => {
                        e.stopPropagation();
                        handleSelectExpert(expert);
                      }}
                    >
                      Quick View
                    </Button>
                    <Button 
                      variant="primary"
                      size="sm"
                      onClick={(e: React.MouseEvent) => handleViewExpertDetails(expert, e)}
                    >
                      Full Profile
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
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
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
                  <h3 className="text-sm font-medium text-neutral-500">Institution</h3>
                  <p className="mt-1">{selectedExpert.institution}</p>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Contact</h3>
                  <p className="mt-1">{selectedExpert.phone} ({selectedExpert.email})</p>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Skills</h3>
                  <div className="mt-1 flex flex-wrap gap-1">
                    {selectedExpert.skills.map((skill: string, index: number) => (
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
                  <p className="mt-1">{selectedExpert.rating}/5</p>
                </div>
                
                <div>
                  <h3 className="text-sm font-medium text-neutral-500">Availability</h3>
                  <p className="mt-1">
                    <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                      selectedExpert.isAvailable
                        ? 'bg-green-100 text-green-800'
                        : 'bg-red-100 text-red-800'
                    }`}>
                      {selectedExpert.isAvailable ? 'Available' : 'Unavailable'}
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
                <p>Created: {new Date(selectedExpert.createdAt).toLocaleDateString()}</p>
                <p>Last updated: {new Date(selectedExpert.updatedAt).toLocaleDateString()}</p>
              </div>
              
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