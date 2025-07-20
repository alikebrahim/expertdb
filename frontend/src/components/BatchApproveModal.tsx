import React, { useState } from 'react';
import { ExpertRequest } from '../types';
import { expertRequestsApi } from '../services/api';
import Button from './ui/Button';
import { Alert } from './ui/Alert';
import FileUpload from './ui/FileUpload';
import { CheckCircle, XCircle, Clock } from 'lucide-react';

interface BatchApproveModalProps {
  requests: ExpertRequest[];
  onClose: () => void;
  onSuccess: (approvedIds: number[]) => void;
}

interface BatchResult {
  id: number;
  name: string;
  status: 'pending' | 'success' | 'failed';
  error?: string;
}

const BatchApproveModal: React.FC<BatchApproveModalProps> = ({
  requests,
  onClose,
  onSuccess
}) => {
  const [approvalFile, setApprovalFile] = useState<File | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [results, setResults] = useState<BatchResult[]>([]);
  const [error, setError] = useState<string | null>(null);

  const handleFileChange = (file: File | null) => {
    setApprovalFile(file);
    setError(null);
  };

  const handleBatchApprove = async () => {
    if (!approvalFile) {
      setError('Please upload an approval document');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    // Initialize results
    const initialResults: BatchResult[] = requests.map(req => ({
      id: req.id,
      name: req.name,
      status: 'pending'
    }));
    setResults(initialResults);

    try {
      const formData = new FormData();
      formData.append('approvalDocument', approvalFile);
      formData.append('requestIds', JSON.stringify(requests.map(r => r.id)));

      const response = await expertRequestsApi.batchApprove(formData);

      if (response.success) {
        // Update results based on API response
        const updatedResults = initialResults.map(result => {
          const apiResult = response.data.results.find(r => r.id === result.id);
          return {
            ...result,
            status: apiResult?.status || 'failed',
            error: apiResult?.error
          };
        });
        setResults(updatedResults);
        
        // Call success callback with approved IDs
        onSuccess(response.data.approvedIds);
      } else {
        setError(response.message || 'Batch approval failed');
      }
    } catch (error) {
      console.error('Error in batch approval:', error);
      setError('An error occurred during batch approval');
      
      // Mark all as failed
      const failedResults = initialResults.map(result => ({
        ...result,
        status: 'failed' as const,
        error: 'Network error'
      }));
      setResults(failedResults);
    } finally {
      setIsSubmitting(false);
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'success':
        return <CheckCircle className="w-5 h-5 text-green-600" />;
      case 'failed':
        return <XCircle className="w-5 h-5 text-red-600" />;
      case 'pending':
        return <Clock className="w-5 h-5 text-yellow-600" />;
      default:
        return null;
    }
  };

  const successCount = results.filter(r => r.status === 'success').length;
  const failedCount = results.filter(r => r.status === 'failed').length;
  const pendingCount = results.filter(r => r.status === 'pending').length;

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-2xl w-full max-h-[90vh] overflow-hidden">
        {/* Modal Header */}
        <div className="bg-green-600 text-white px-6 py-4 flex items-center justify-between">
          <h2 className="text-xl font-semibold">
            Batch Approve {requests.length} Request{requests.length > 1 ? 's' : ''}
          </h2>
          <button
            onClick={onClose}
            className="text-white hover:text-gray-200 text-2xl font-bold"
          >
            ×
          </button>
        </div>

        {/* Modal Content */}
        <div className="p-6 max-h-[70vh] overflow-y-auto">
          {error && (
            <Alert variant="error" className="mb-4" onClose={() => setError(null)}>
              {error}
            </Alert>
          )}

          {!isSubmitting && results.length === 0 && (
            <div className="space-y-6">
              <div className="bg-blue-50 p-4 rounded-lg border border-blue-200">
                <h3 className="font-medium text-blue-900 mb-2">Requests to Approve</h3>
                <div className="space-y-2">
                  {requests.map((request) => (
                    <div key={request.id} className="flex items-center space-x-3 text-sm">
                      <div className="w-2 h-2 bg-blue-600 rounded-full"></div>
                      <span className="font-medium">{request.name}</span>
                      <span className="text-gray-600">({request.institution})</span>
                    </div>
                  ))}
                </div>
              </div>

              <div className="space-y-4">
                <FileUpload
                  onFileSelect={handleFileChange}
                  accept=".pdf"
                  maxSize={20}
                  currentFile={approvalFile}
                  error={!approvalFile && error?.includes('document') ? 'Approval document is required' : ''}
                  label="Approval Document"
                  required
                />
                <p className="text-sm text-gray-600">
                  This approval document will be applied to all selected requests.
                  Make sure it's a valid approval certificate in PDF format.
                </p>
              </div>
            </div>
          )}

          {/* Results Display */}
          {results.length > 0 && (
            <div className="space-y-4">
              <div className="grid grid-cols-3 gap-4 text-center">
                <div className="bg-green-50 p-3 rounded-lg">
                  <div className="text-2xl font-bold text-green-600">{successCount}</div>
                  <div className="text-sm text-green-700">Approved</div>
                </div>
                <div className="bg-red-50 p-3 rounded-lg">
                  <div className="text-2xl font-bold text-red-600">{failedCount}</div>
                  <div className="text-sm text-red-700">Failed</div>
                </div>
                <div className="bg-yellow-50 p-3 rounded-lg">
                  <div className="text-2xl font-bold text-yellow-600">{pendingCount}</div>
                  <div className="text-sm text-yellow-700">Pending</div>
                </div>
              </div>

              <div className="border rounded-lg">
                <div className="bg-gray-50 px-4 py-3 border-b">
                  <h4 className="font-medium text-gray-900">Processing Results</h4>
                </div>
                <div className="divide-y">
                  {results.map((result) => (
                    <div key={result.id} className="flex items-center justify-between p-4">
                      <div className="flex items-center space-x-3">
                        {getStatusIcon(result.status)}
                        <div>
                          <div className="font-medium text-gray-900">{result.name}</div>
                          {result.error && (
                            <div className="text-sm text-red-600">{result.error}</div>
                          )}
                        </div>
                      </div>
                      <div className="text-sm text-gray-500">
                        {result.status === 'success' ? 'Approved' : 
                         result.status === 'failed' ? 'Failed' : 'Processing...'}
                      </div>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          )}
        </div>

        {/* Modal Footer */}
        <div className="bg-gray-50 px-6 py-4 flex justify-end space-x-3">
          {results.length === 0 ? (
            <>
              <Button variant="outline" onClick={onClose}>
                Cancel
              </Button>
              <Button
                variant="primary"
                onClick={handleBatchApprove}
                disabled={isSubmitting || !approvalFile}
                className="bg-green-600 hover:bg-green-700"
              >
                {isSubmitting ? 'Processing...' : 'Approve All'}
              </Button>
            </>
          ) : (
            <Button variant="outline" onClick={onClose}>
              {pendingCount > 0 ? 'Close' : 'Done'}
            </Button>
          )}
        </div>
      </div>
    </div>
  );
};

export default BatchApproveModal;