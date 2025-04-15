import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { expertsApi } from '../services/api';
import { Expert } from '../types';
import Button from '../components/ui/Button';
import DocumentManager from '../components/DocumentManager';
import EngagementList from '../components/EngagementList';
import Modal from '../components/Modal';
import ExpertForm from '../components/ExpertForm';
import { useAuth } from '../hooks/useAuth';

const ExpertDetailPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [expert, setExpert] = useState<Expert | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);

  const fetchExpert = useCallback(async () => {
    if (!id) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const expertId = id;
      const response = await expertsApi.getExpertById(expertId);
      
      if (response.success) {
        setExpert(response.data);
      } else {
        setError(response.message || 'Failed to load expert details');
      }
    } catch (error) {
      console.error('Error fetching expert details:', error);
      setError('An error occurred while loading expert details');
    } finally {
      setIsLoading(false);
    }
  }, [id]);

  useEffect(() => {
    fetchExpert();
  }, [fetchExpert]);

  if (isLoading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="text-center">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary mx-auto"></div>
          <p className="mt-4 text-gray-600">Loading expert details...</p>
        </div>
      </div>
    );
  }

  if (error || !expert) {
    return (
      <div className="max-w-4xl mx-auto mt-8 px-4">
        <div className="bg-red-50 border border-red-200 text-red-600 p-4 rounded-md">
          <h2 className="text-lg font-medium mb-2">Error</h2>
          <p>{error || 'Expert not found'}</p>
          <Button 
            className="mt-4" 
            variant="outline" 
            onClick={() => navigate(-1)}
          >
            Go Back
          </Button>
        </div>
      </div>
    );
  }

  const handleEditExpert = (updatedExpert: Expert) => {
    setExpert(updatedExpert);
    setIsEditModalOpen(false);
  };

  return (
    <div className="max-w-6xl mx-auto px-4 py-8">
      <div className="mb-6 flex justify-between items-center">
        <Button 
          variant="outline" 
          onClick={() => navigate(-1)}
          icon={<span>&larr;</span>}
        >
          Back
        </Button>
        
        <div className="flex space-x-3">
          {user?.role === 'admin' && (
            <Button 
              variant="secondary"
              onClick={() => setIsEditModalOpen(true)}
            >
              Edit Expert
            </Button>
          )}
          <Button 
            onClick={() => window.open(`/api/experts/${expert.id}/approval-pdf`, '_blank')}
          >
            Download Expert PDF
          </Button>
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Expert Info Panel */}
        <div className="lg:col-span-2 space-y-6">
          <div className="bg-white rounded-md shadow-sm p-6">
            <h1 className="text-2xl font-bold mb-1 text-primary">{expert.name}</h1>
            <p className="text-gray-500 mb-4">{expert.affiliation}</p>
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
              <div>
                <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-2">Contact Information</h3>
                <p className="mb-1"><strong>Contact:</strong> {expert.primaryContact}</p>
                <p><strong>Type:</strong> {expert.contactType}</p>
              </div>
              <div>
                <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-2">Employment Details</h3>
                <p className="mb-1"><strong>Role:</strong> {expert.role}</p>
                <p className="mb-1"><strong>Employment Type:</strong> {expert.employmentType}</p>
                <p><strong>Nationality:</strong> {expert.isBahraini ? 'Bahraini' : 'International'}</p>
              </div>
            </div>
            
            <div className="mb-6">
              <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-2">Expertise & Skills</h3>
              <div className="flex flex-wrap gap-2 mb-4">
                {expert.skills.map((skill, index) => (
                  <span 
                    key={index} 
                    className="bg-gray-100 px-3 py-1 rounded-full text-sm text-gray-700"
                  >
                    {skill}
                  </span>
                ))}
              </div>
              <p className="mb-1"><strong>Availability:</strong> {expert.availability}</p>
              <p><strong>Rating:</strong> {expert.rating} / 5</p>
            </div>
            
            <div>
              <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-2">Biography</h3>
              <div className="prose max-w-none">
                <p>{expert.biography || 'No biography available.'}</p>
              </div>
            </div>
          </div>
          
          {/* Engagement List */}
          <EngagementList expertId={expert.id} />
        </div>
        
        {/* Side Panel */}
        <div className="lg:col-span-1 space-y-6">
          {/* Document Manager Panel */}
          <DocumentManager expertId={expert.id} />
        </div>
      </div>
      
      {/* Edit Expert Modal */}
      <Modal
        isOpen={isEditModalOpen}
        onClose={() => setIsEditModalOpen(false)}
        title="Edit Expert"
        size="lg"
      >
        {expert && (
          <ExpertForm
            expert={expert}
            onSuccess={handleEditExpert}
            onCancel={() => setIsEditModalOpen(false)}
          />
        )}
      </Modal>
    </div>
  );
};

export default ExpertDetailPage;