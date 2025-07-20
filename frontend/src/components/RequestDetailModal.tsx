import { useState } from 'react';
import { ExpertRequest } from '../types';
import { expertRequestsApi } from '../services/api';
import Button from './ui/Button';
import { Alert } from './ui/Alert';
import FileUpload from './ui/FileUpload';

interface RequestDetailModalProps {
  request: ExpertRequest;
  onClose: () => void;
  onRequestUpdate: (message: string) => void;
}

const RequestDetailModal = ({ request, onClose, onRequestUpdate }: RequestDetailModalProps) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [approvalFile, setApprovalFile] = useState<File | null>(null);
  const [rejectionReason, setRejectionReason] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState<'details' | 'documents' | 'actions'>('details');

  const handleApprove = async () => {
    if (!approvalFile) {
      setError('Approval document is required');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      const formData = new FormData();
      formData.append('status', 'approved');
      formData.append('approvalDocument', approvalFile);

      const response = await expertRequestsApi.updateExpertRequest(request.id, formData);

      if (response.success) {
        onRequestUpdate('Expert request approved successfully. Expert profile has been added to the database.');
      } else {
        setError(response.message || 'Failed to approve request');
      }
    } catch (error) {
      console.error('Error approving request:', error);
      setError('An error occurred while approving the request');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleReject = async () => {
    if (!rejectionReason.trim()) {
      setError('Rejection reason is required');
      return;
    }

    setIsSubmitting(true);
    setError(null);

    try {
      const formData = new FormData();
      formData.append('status', 'rejected');
      formData.append('rejectionReason', rejectionReason);

      const response = await expertRequestsApi.updateExpertRequest(request.id, formData);

      if (response.success) {
        onRequestUpdate('Expert request has been rejected. The submitter will be notified.');
      } else {
        setError(response.message || 'Failed to reject request');
      }
    } catch (error) {
      console.error('Error rejecting request:', error);
      setError('An error occurred while rejecting the request');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleFileChange = (file: File | null) => {
    setApprovalFile(file);
    setError(null);
  };

  const formatSkills = (skills: string | string[]) => {
    try {
      if (Array.isArray(skills)) {
        return skills.join(', ');
      }
      if (typeof skills === 'string' && skills.startsWith('[') && skills.endsWith(']')) {
        const skillsArray = JSON.parse(skills);
        return skillsArray.join(', ');
      }
      return skills;
    } catch {
      return Array.isArray(skills) ? skills.join(', ') : skills;
    }
  };

  return (
    <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div className="bg-white rounded-lg shadow-xl max-w-4xl w-full max-h-[90vh] overflow-hidden">
        {/* Modal Header */}
        <div className="bg-primary text-white px-6 py-4 flex items-center justify-between">
          <h2 className="text-xl font-semibold">Expert Request Details</h2>
          <button
            onClick={onClose}
            className="text-white hover:text-gray-200 text-2xl font-bold"
          >
            ×
          </button>
        </div>

        {/* Tab Navigation */}
        <div className="border-b border-gray-200">
          <nav className="flex">
            {[
              { id: 'details', label: 'Expert Information' },
              { id: 'documents', label: 'Documents' },
              { id: 'actions', label: 'Admin Actions' }
            ].map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id as 'details' | 'documents' | 'actions')}
                className={`px-6 py-3 text-sm font-medium border-b-2 ${
                  activeTab === tab.id
                    ? 'border-primary text-primary'
                    : 'border-transparent text-gray-500 hover:text-gray-700'
                }`}
              >
                {tab.label}
              </button>
            ))}
          </nav>
        </div>

        {/* Modal Content */}
        <div className="p-6 max-h-[60vh] overflow-y-auto">
          {error && (
            <Alert variant="error" className="mb-4" onClose={() => setError(null)}>
              {error}
            </Alert>
          )}

          {/* Expert Information Tab */}
          {activeTab === 'details' && (
            <div className="space-y-6">
              {/* Personal Information */}
              <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Personal Information</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Name</label>
                    <p className="text-gray-900">{request.name}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Designation</label>
                    <p className="text-gray-900">{request.designation}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Institution</label>
                    <p className="text-gray-900">{request.institution}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Email</label>
                    <p className="text-gray-900">{request.email}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Phone</label>
                    <p className="text-gray-900">{request.phone}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Nationality</label>
                    <p className="text-gray-900">{request.isBahraini ? 'Bahraini' : 'Non-Bahraini'}</p>
                  </div>
                </div>
              </div>

              {/* Professional Details */}
              <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Professional Details</h3>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Role</label>
                    <p className="text-gray-900">{request.role}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Employment Type</label>
                    <p className="text-gray-900">{request.employmentType}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Rating</label>
                    <p className="text-gray-900">{request.rating}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Availability</label>
                    <p className="text-gray-900">{request.isAvailable ? 'Available' : 'Not Available'}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Training Status</label>
                    <p className="text-gray-900">{request.isTrained ? 'Trained' : 'Not Trained'}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Publishing Permission</label>
                    <p className="text-gray-900">{request.isPublished ? 'Allowed' : 'Not Allowed'}</p>
                  </div>
                </div>
              </div>

              {/* Expertise */}
              <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Expertise Areas</h3>
                <div className="grid grid-cols-1 gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700">General Area</label>
                    <p className="text-gray-900">{request.generalArea}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Specialized Area</label>
                    <p className="text-gray-900">{request.specializedArea}</p>
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700">Skills</label>
                    <p className="text-gray-900">{formatSkills(request.skills)}</p>
                  </div>
                </div>
              </div>

              {/* Biography */}
              <div className="bg-gray-50 p-4 rounded-lg">
                <h3 className="text-lg font-semibold text-gray-900 mb-4">Biography</h3>
                <div className="prose max-w-none">
                  <p className="text-gray-900 whitespace-pre-wrap">{request.biography}</p>
                </div>
              </div>
            </div>
          )}

          {/* Documents Tab */}
          {activeTab === 'documents' && (
            <div className="space-y-4">
              <h3 className="text-lg font-semibold text-gray-900">Documents</h3>
              
              {/* CV Document */}
              <div className="bg-gray-50 p-4 rounded-lg">
                <div className="flex items-center justify-between">
                  <div>
                    <h4 className="font-medium text-gray-900">CV Document</h4>
                    <p className="text-sm text-gray-600">Expert's curriculum vitae</p>
                  </div>
                  <div className="flex space-x-2">
                    {request.cvPath && (
                      <a
                        href={`/api/documents/${request.cvPath}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-block"
                      >
                        <Button variant="outline" size="sm">
                          📄 View PDF
                        </Button>
                      </a>
                    )}
                  </div>
                </div>
              </div>

              {/* Approval Document (if approved) */}
              {request.status === 'approved' && request.approvalDocumentPath && (
                <div className="bg-green-50 p-4 rounded-lg border border-green-200">
                  <div className="flex items-center justify-between">
                    <div>
                      <h4 className="font-medium text-green-900">Approval Document</h4>
                      <p className="text-sm text-green-700">Official approval certificate</p>
                    </div>
                    <div className="flex space-x-2">
                      <a
                        href={`/api/documents/${request.approvalDocumentPath}`}
                        target="_blank"
                        rel="noopener noreferrer"
                        className="inline-block"
                      >
                        <Button variant="outline" size="sm">
                          📄 View Approval
                        </Button>
                      </a>
                    </div>
                  </div>
                </div>
              )}
            </div>
          )}

          {/* Admin Actions Tab */}
          {activeTab === 'actions' && (
            <div className="space-y-6">
              <h3 className="text-lg font-semibold text-gray-900">Administrative Actions</h3>
              
              {/* Current Status */}
              <div className="bg-gray-50 p-4 rounded-lg">
                <h4 className="font-medium text-gray-900 mb-2">Current Status</h4>
                <span className={`inline-flex items-center px-3 py-1 rounded-full text-sm font-medium ${
                  request.status === 'approved' ? 'bg-green-100 text-green-800' :
                  request.status === 'rejected' ? 'bg-red-100 text-red-800' :
                  'bg-yellow-100 text-yellow-800'
                }`}>
                  {request.status.charAt(0).toUpperCase() + request.status.slice(1)}
                </span>
                
                {request.status === 'rejected' && request.rejectionReason && (
                  <div className="mt-2">
                    <p className="text-sm text-gray-600">Rejection Reason:</p>
                    <p className="text-sm text-gray-900 bg-white p-2 rounded border">
                      {request.rejectionReason}
                    </p>
                  </div>
                )}
              </div>

              {/* Action Forms */}
              {request.status === 'pending' && (
                <div className="space-y-6">
                  {/* Approve Section */}
                  <div className="bg-green-50 p-4 rounded-lg border border-green-200">
                    <h4 className="font-medium text-green-900 mb-4">Approve Request</h4>
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
                      <p className="text-xs text-green-700">
                        Upload the official approval document (PDF format required)
                      </p>
                      <Button
                        variant="primary"
                        onClick={handleApprove}
                        disabled={isSubmitting || !approvalFile}
                        className="bg-green-600 hover:bg-green-700"
                      >
                        {isSubmitting ? 'Approving...' : '✓ Approve Request'}
                      </Button>
                    </div>
                  </div>

                  {/* Reject Section */}
                  <div className="bg-red-50 p-4 rounded-lg border border-red-200">
                    <h4 className="font-medium text-red-900 mb-4">Reject Request</h4>
                    <div className="space-y-4">
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                          Rejection Reason <span className="text-red-500">*</span>
                        </label>
                        <textarea
                          value={rejectionReason}
                          onChange={(e) => setRejectionReason(e.target.value)}
                          rows={4}
                          className="w-full px-3 py-2 border border-red-300 rounded focus:outline-none focus:ring-1 focus:ring-red-500"
                          placeholder="Please provide a detailed reason for rejection..."
                          required
                        />
                        <p className="text-xs text-red-700 mt-1">
                          This reason will be sent to the submitter for review and correction
                        </p>
                      </div>
                      <Button
                        variant="outline"
                        onClick={handleReject}
                        disabled={isSubmitting || !rejectionReason.trim()}
                        className="border-red-300 text-red-700 hover:bg-red-50"
                      >
                        {isSubmitting ? 'Rejecting...' : '✗ Reject Request'}
                      </Button>
                    </div>
                  </div>
                </div>
              )}

              {/* Already Processed */}
              {request.status !== 'pending' && (
                <div className="bg-gray-50 p-4 rounded-lg">
                  <p className="text-gray-600">
                    This request has already been {request.status}. No further actions are available.
                  </p>
                </div>
              )}
            </div>
          )}
        </div>

        {/* Modal Footer */}
        <div className="bg-gray-50 px-6 py-4 flex justify-end">
          <Button variant="outline" onClick={onClose}>
            Close
          </Button>
        </div>
      </div>
    </div>
  );
};

export default RequestDetailModal;
