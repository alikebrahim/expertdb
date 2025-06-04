import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import * as expertsApi from '../api/experts';
import * as documentsApi from '../api/documents';
import { Expert, Document } from '../types';
import Layout from '../components/layout/Layout';
import Button from '../components/ui/Button';
import { DocumentList, DocumentUpload } from '../components/documents';
import EngagementList from '../components/EngagementList';
import Modal from '../components/Modal';
import ExpertForm from '../components/ExpertForm';
import { useAuth } from '../hooks/useAuth';
import { useUI } from '../hooks/useUI';
import { formatDate } from '../utils/formatters';

const ExpertDetailPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const { addNotification } = useUI();
  const [expert, setExpert] = useState<Expert | null>(null);
  const [documents, setDocuments] = useState<Document[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isDocumentsLoading, setIsDocumentsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isEditModalOpen, setIsEditModalOpen] = useState(false);
  const [isUploadModalOpen, setIsUploadModalOpen] = useState(false);

  const fetchExpert = useCallback(async () => {
    if (!id) return;
    
    setIsLoading(true);
    setError(null);
    
    try {
      const expertId = id;
      const response = await expertsApi.getExpertById(expertId);
      
      if (response.success && response.data) {
        setExpert(response.data);
      } else {
        setError(response.message || 'Failed to load expert details');
        addNotification({
          type: 'error',
          message: response.message || 'Failed to load expert details',
        });
      }
    } catch (error) {
      console.error('Error fetching expert details:', error);
      setError('An error occurred while loading expert details');
      addNotification({
        type: 'error',
        message: 'An error occurred while loading expert details',
      });
    } finally {
      setIsLoading(false);
    }
  }, [id, addNotification]);

  const fetchDocuments = useCallback(async () => {
    if (!id) return;
    
    setIsDocumentsLoading(true);
    
    try {
      // In a real implementation, you'd use the actual API call
      // const response = await documentsApi.getExpertDocuments(parseInt(id));
      
      // For demo purposes, create mock documents
      const mockDocuments: Document[] = [
        {
          id: 1,
          expertId: parseInt(id),
          filename: 'expert_cv.pdf',
          originalFilename: 'John_Smith_CV_2025.pdf',
          documentType: 'cv',
          contentType: 'application/pdf',
          size: 1245678,
          createdAt: new Date().toISOString(),
          filePath: '/uploads/expert_cv.pdf',
        },
        {
          id: 2,
          expertId: parseInt(id),
          filename: 'certificate.jpg',
          originalFilename: 'Teaching_Certificate.jpg',
          documentType: 'certificate',
          contentType: 'image/jpeg',
          size: 587341,
          createdAt: new Date(Date.now() - 7 * 24 * 60 * 60 * 1000).toISOString(),
          filePath: '/uploads/certificate.jpg',
        },
      ];
      
      setDocuments(mockDocuments);
    } catch (error) {
      console.error('Error fetching expert documents:', error);
      addNotification({
        type: 'error',
        message: 'Failed to load expert documents',
      });
    } finally {
      setIsDocumentsLoading(false);
    }
  }, [id, addNotification]);

  useEffect(() => {
    fetchExpert();
    fetchDocuments();
  }, [fetchExpert, fetchDocuments]);

  const handleEditExpert = (updatedExpert: Expert) => {
    setExpert(updatedExpert);
    setIsEditModalOpen(false);
    addNotification({
      type: 'success',
      message: 'Expert updated successfully',
      duration: 3000,
    });
  };

  const handleDocumentUpload = (documentId: number) => {
    fetchDocuments();
    setIsUploadModalOpen(false);
    addNotification({
      type: 'success',
      message: 'Document uploaded successfully',
      duration: 3000,
    });
  };

  const handleDocumentDownload = (document: Document) => {
    // In a real implementation, you'd use the API to download the file
    addNotification({
      type: 'info',
      message: `Downloading ${document.originalFilename}...`,
      duration: 3000,
    });
    
    // Simulate a download
    setTimeout(() => {
      addNotification({
        type: 'success',
        message: `Download complete: ${document.originalFilename}`,
        duration: 3000,
      });
    }, 1500);
  };

  const handleDocumentDelete = (document: Document) => {
    // In a real implementation, you'd call the API to delete the document
    setDocuments(documents.filter(d => d.id !== document.id));
    addNotification({
      type: 'success',
      message: 'Document deleted successfully',
      duration: 3000,
    });
  };

  if (isLoading) {
    return (
      <Layout>
        <div className="flex justify-center items-center min-h-screen">
          <div className="text-center">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary mx-auto"></div>
            <p className="mt-4 text-gray-600">Loading expert details...</p>
          </div>
        </div>
      </Layout>
    );
  }

  if (error || !expert) {
    return (
      <Layout>
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
      </Layout>
    );
  }

  return (
    <Layout>
      <div className="mb-6 flex flex-wrap justify-between items-center gap-4">
        <Button 
          variant="outline" 
          onClick={() => navigate(-1)}
          icon={<span>&larr;</span>}
        >
          Back
        </Button>
        
        <div className="flex flex-wrap gap-3">
          {user?.role === 'admin' && (
            <Button 
              variant="secondary"
              onClick={() => setIsEditModalOpen(true)}
            >
              Edit Expert
            </Button>
          )}
          <Button 
            onClick={() => {
              addNotification({
                type: 'success',
                message: 'Expert PDF downloaded successfully',
                duration: 3000,
              });
            }}
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
              <div className="grid grid-cols-2 gap-4">
                <p className="mb-1">
                  <strong>Availability:</strong>{' '}
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    expert.availability === 'Available' 
                      ? 'bg-green-100 text-green-800'
                      : expert.availability === 'Limited'
                      ? 'bg-yellow-100 text-yellow-800'
                      : 'bg-red-100 text-red-800'
                  }`}>
                    {expert.availability}
                  </span>
                </p>
                <p>
                  <strong>Rating:</strong>{' '}
                  <span className="inline-flex items-center">
                    {expert.rating}
                    <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-yellow-500 ml-1" viewBox="0 0 20 20" fill="currentColor">
                      <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                    </svg>
                  </span>
                </p>
              </div>
            </div>
            
            <div>
              <h3 className="text-sm uppercase tracking-wider text-gray-500 mb-2">Biography</h3>
              <div className="prose max-w-none">
                <p>{expert.biography || 'No biography available.'}</p>
              </div>
            </div>
            
            <div className="mt-6 border-t pt-4 text-sm text-gray-500">
              <p>Last updated: {formatDate(expert.updated_at)}</p>
            </div>
          </div>
          
          {/* Engagement List Panel */}
          <div className="bg-white rounded-md shadow-sm p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-lg font-semibold text-primary">Engagement History</h2>
              {user?.role === 'admin' && (
                <Button variant="outline" size="sm">
                  Add Engagement
                </Button>
              )}
            </div>
            <EngagementList expertId={expert.id} />
          </div>
        </div>
        
        {/* Side Panel */}
        <div className="lg:col-span-1 space-y-6">
          {/* Document Management Panel */}
          <div className="bg-white rounded-md shadow-sm p-6">
            <div className="flex justify-between items-center mb-4">
              <h2 className="text-lg font-semibold text-primary">Documents</h2>
              {user?.role === 'admin' && (
                <Button 
                  variant="outline" 
                  size="sm"
                  onClick={() => setIsUploadModalOpen(true)}
                >
                  Upload New
                </Button>
              )}
            </div>
            <DocumentList 
              documents={documents}
              isLoading={isDocumentsLoading}
              onDownload={handleDocumentDownload}
              onDelete={user?.role === 'admin' ? handleDocumentDelete : undefined}
            />
          </div>
          
          {/* Additional side panels can be added here */}
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
      
      {/* Upload Document Modal */}
      <Modal
        isOpen={isUploadModalOpen}
        onClose={() => setIsUploadModalOpen(false)}
        title="Upload Document"
        size="md"
      >
        {expert && (
          <DocumentUpload
            expertId={expert.id}
            onUploadSuccess={handleDocumentUpload}
          />
        )}
      </Modal>
    </Layout>
  );
};

export default ExpertDetailPage;