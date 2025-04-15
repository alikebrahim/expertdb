import { useState, useEffect } from 'react';
import { ExpertRequest } from '../types';
import { expertRequestsApi } from '../services/api';
import { useAuth } from '../hooks/useAuth';
import ExpertRequestForm from '../components/ExpertRequestForm';
import ExpertRequestTable from '../components/ExpertRequestTable';
import Button from '../components/ui/Button';

const ExpertRequestPage = () => {
  const { user } = useAuth();
  const [requests, setRequests] = useState<ExpertRequest[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [resubmittingRequest, setResubmittingRequest] = useState<ExpertRequest | null>(null);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [limit] = useState(10);
  
  // Fetch user's expert requests
  useEffect(() => {
    const fetchRequests = async () => {
      if (!user) return;
      
      setIsLoading(true);
      setError(null);
      
      try {
        const response = await expertRequestsApi.getExpertRequests(page, limit, {
          userId: user.id
        });
        
        if (response.success) {
          setRequests(response.data.data);
          setTotalPages(response.data.totalPages);
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
  }, [user, page, limit]);
  
  const handleNewRequest = () => {
    setResubmittingRequest(null);
    setShowForm(true);
    setSuccessMessage(null);
    
    // Scroll to form
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };
  
  const handleResubmit = (request: ExpertRequest) => {
    setResubmittingRequest(request);
    setShowForm(true);
    setSuccessMessage(null);
    
    // Scroll to form
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };
  
  const handleSuccess = () => {
    setShowForm(false);
    setSuccessMessage(
      resubmittingRequest
        ? 'Your expert request has been resubmitted successfully!'
        : 'Your expert request has been submitted successfully!'
    );
    
    // Refresh the list of requests
    if (user) {
      setIsLoading(true);
      expertRequestsApi.getExpertRequests(page, limit, { userId: user.id })
        .then(response => {
          if (response.success) {
            setRequests(response.data.data);
            setTotalPages(response.data.totalPages);
          }
        })
        .finally(() => setIsLoading(false));
    }
    
    setResubmittingRequest(null);
  };
  
  const handlePageChange = (newPage: number) => {
    setPage(newPage);
  };
  
  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-primary">Expert Requests</h1>
        <p className="text-neutral-600">
          Submit new expert requests or view your previous submissions
        </p>
      </div>
      
      {successMessage && (
        <div className="bg-green-100 text-green-800 p-4 rounded mb-6">
          {successMessage}
        </div>
      )}
      
      {showForm ? (
        <div className="bg-white rounded-md shadow p-6 mb-8">
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-semibold text-primary">
              {resubmittingRequest ? 'Resubmit Expert Request' : 'Submit New Expert Request'}
            </h2>
            <Button
              variant="outline"
              onClick={() => setShowForm(false)}
            >
              Cancel
            </Button>
          </div>
          
          <ExpertRequestForm onSuccess={handleSuccess} />
        </div>
      ) : (
        <div className="flex justify-end mb-6">
          <Button onClick={handleNewRequest}>
            Submit New Request
          </Button>
        </div>
      )}
      
      <div className="bg-white rounded-md shadow p-6">
        <h2 className="text-xl font-semibold text-primary mb-4">
          Your Request History
        </h2>
        
        <ExpertRequestTable
          requests={requests}
          isLoading={isLoading}
          error={error}
          onResubmit={handleResubmit}
          pagination={{
            currentPage: page,
            totalPages: totalPages,
            onPageChange: handlePageChange
          }}
        />
      </div>
    </div>
  );
};

export default ExpertRequestPage;