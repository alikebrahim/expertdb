import React from 'react';
import { Expert } from '../types';

interface ExpertProfileModalProps {
  expert: Expert;
  isOpen: boolean;
  onClose: () => void;
}

const ExpertProfileModal: React.FC<ExpertProfileModalProps> = ({ expert, isOpen, onClose }) => {
  if (!isOpen) return null;

  const handleCVDownload = () => {
    if (expert.cvPath) {
      // Create a download link for the CV
      const link = document.createElement('a');
      link.href = `/api/documents/download/${expert.cvPath}`;
      link.download = `${expert.name}_CV.pdf`;
      link.click();
    }
  };

  const handleBackgroundClick = (e: React.MouseEvent) => {
    if (e.target === e.currentTarget) {
      onClose();
    }
  };

  return (
    <div 
      className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      onClick={handleBackgroundClick}
    >
      <div className="bg-white rounded-lg shadow-xl max-w-4xl max-h-[90vh] w-full mx-4 overflow-y-auto">
        {/* Header */}
        <div className="sticky top-0 bg-white border-b border-neutral-200 px-6 py-4 flex justify-between items-center">
          <h2 className="text-xl font-semibold text-neutral-900">Expert Profile</h2>
          <button
            onClick={onClose}
            className="text-neutral-400 hover:text-neutral-600 transition-colors"
          >
            <svg className="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        {/* Content */}
        <div className="p-6 space-y-6">
          {/* Basic Information */}
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-4">
              <div>
                <h3 className="text-lg font-medium text-neutral-900 mb-3">Basic Information</h3>
                <div className="space-y-2">
                  <div className="flex items-center">
                    <span className="font-medium text-neutral-700 w-24">Name:</span>
                    <span className="text-neutral-900">{expert.name}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="font-medium text-neutral-700 w-24">ID:</span>
                    <span className="text-neutral-600">{expert.expertId}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="font-medium text-neutral-700 w-24">Role:</span>
                    <span className="text-neutral-900">{expert.role}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="font-medium text-neutral-700 w-24">Rating:</span>
                    <div className="flex items-center">
                      <span className="text-neutral-900 mr-2">{expert.rating === 0 ? 'N/A' : expert.rating}</span>
                      {expert.rating > 0 && (
                        <svg className="h-4 w-4 text-yellow-500">
                          <path fill="currentColor" d="M12 2l3.09 6.26L22 9.27l-5 4.87 1.18 6.88L12 17.77l-6.18 3.25L7 14.14 2 9.27l6.91-1.01L12 2z" />
                        </svg>
                      )}
                    </div>
                  </div>
                </div>
              </div>

              <div>
                <h4 className="font-medium text-neutral-900 mb-2">Status</h4>
                <div className="space-y-2">
                  <div className="flex items-center">
                    <span className="font-medium text-neutral-700 w-24">Available:</span>
                    <span className={`px-2 py-1 text-xs rounded-full ${
                      expert.isAvailable ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'
                    }`}>
                      {expert.isAvailable ? 'Yes' : 'No'}
                    </span>
                  </div>
                  <div className="flex items-center">
                    <span className="font-medium text-neutral-700 w-24">Bahraini:</span>
                    <span className={`px-2 py-1 text-xs rounded-full ${
                      expert.isBahraini ? 'bg-blue-100 text-blue-800' : 'bg-gray-100 text-gray-800'
                    }`}>
                      {expert.isBahraini ? 'Yes' : 'No'}
                    </span>
                  </div>
                  <div className="flex items-center">
                    <span className="font-medium text-neutral-700 w-24">Trained:</span>
                    <span className={`px-2 py-1 text-xs rounded-full ${
                      expert.isTrained ? 'bg-green-100 text-green-800' : 'bg-yellow-100 text-yellow-800'
                    }`}>
                      {expert.isTrained ? 'Yes' : 'No'}
                    </span>
                  </div>
                  <div className="flex items-center">
                    <span className="font-medium text-neutral-700 w-24">Published:</span>
                    <span className={`px-2 py-1 text-xs rounded-full ${
                      expert.isPublished ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                    }`}>
                      {expert.isPublished ? 'Yes' : 'No'}
                    </span>
                  </div>
                </div>
              </div>
            </div>

            <div className="space-y-4">
              <div>
                <h4 className="font-medium text-neutral-900 mb-2">Professional Details</h4>
                <div className="space-y-2">
                  <div className="flex items-center">
                    <span className="font-medium text-neutral-700 w-32">Institution:</span>
                    <span className="text-neutral-900">{expert.institution}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="font-medium text-neutral-700 w-32">Designation:</span>
                    <span className="text-neutral-900">{expert.designation}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="font-medium text-neutral-700 w-32">Employment:</span>
                    <span className="text-neutral-900">{expert.employmentType}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="font-medium text-neutral-700 w-32">General Area:</span>
                    <span className="text-neutral-900">{expert.generalAreaName}</span>
                  </div>
                  <div className="flex items-start">
                    <span className="font-medium text-neutral-700 w-32">Specialized:</span>
                    <span className="text-neutral-900">{expert.specializedArea}</span>
                  </div>
                </div>
              </div>

              <div>
                <h4 className="font-medium text-neutral-900 mb-2">Contact Information</h4>
                <div className="space-y-2">
                  <div className="flex items-center">
                    <span className="font-medium text-neutral-700 w-20">Phone:</span>
                    <span className="text-neutral-900">{expert.phone || 'Not provided'}</span>
                  </div>
                  <div className="flex items-center">
                    <span className="font-medium text-neutral-700 w-20">Email:</span>
                    <span className="text-neutral-900">{expert.email || 'Not provided'}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          {/* Biography */}
          {expert.biography && (
            <div>
              <h4 className="font-medium text-neutral-900 mb-2">Biography</h4>
              <p className="text-neutral-700 text-sm leading-relaxed">{expert.biography}</p>
            </div>
          )}

          {/* Skills */}
          {expert.skills && expert.skills.length > 0 && (
            <div>
              <h4 className="font-medium text-neutral-900 mb-2">Skills</h4>
              <div className="flex flex-wrap gap-2">
                {expert.skills.map((skill, index) => (
                  <span
                    key={index}
                    className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full"
                  >
                    {skill}
                  </span>
                ))}
              </div>
            </div>
          )}

          {/* CV Download */}
          <div>
            <h4 className="font-medium text-neutral-900 mb-2">Documents</h4>
            <div className="space-y-2">
              {expert.cvPath && (
                <button
                  onClick={handleCVDownload}
                  className="inline-flex items-center px-4 py-2 bg-primary text-white rounded-md hover:bg-primary-dark transition-colors"
                >
                  <svg className="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  Download CV
                </button>
              )}
              {!expert.cvPath && (
                <p className="text-neutral-500 text-sm">No CV available</p>
              )}
            </div>
          </div>

          {/* Metadata */}
          <div className="pt-4 border-t border-neutral-200">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm text-neutral-600">
              <div>
                <span className="font-medium">Created:</span> {new Date(expert.createdAt).toLocaleDateString()}
              </div>
              <div>
                <span className="font-medium">Updated:</span> {new Date(expert.updatedAt).toLocaleDateString()}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export default ExpertProfileModal;