import { useState } from 'react';
import { Document } from '../types';
import DocumentUpload from './DocumentUpload';
import DocumentList from './DocumentList';

interface DocumentManagerProps {
  expertId: number;
}

const DocumentManager = ({ expertId }: DocumentManagerProps) => {
  const [showUploadForm, setShowUploadForm] = useState(false);
  const [refreshCounter, setRefreshCounter] = useState(0);

  const handleDocumentUploaded = (_document: Document) => {
    setShowUploadForm(false);
    setRefreshCounter(prev => prev + 1);
  };

  const handleRefreshNeeded = () => {
    setRefreshCounter(prev => prev + 1);
  };

  return (
    <div className="bg-white rounded-md shadow-sm">
      <div className="border-b px-4 py-3 flex justify-between items-center">
        <h2 className="text-lg font-medium">Expert Documents</h2>
        <button
          onClick={() => setShowUploadForm(!showUploadForm)}
          className="text-primary hover:text-primary-dark text-sm font-medium underline"
        >
          {showUploadForm ? 'Cancel Upload' : '+ Add Document'}
        </button>
      </div>

      {showUploadForm && (
        <div className="p-4 border-b">
          <DocumentUpload 
            expertId={expertId} 
            onSuccess={handleDocumentUploaded} 
          />
        </div>
      )}

      <div key={refreshCounter}>
        <DocumentList 
          expertId={expertId} 
          onRefreshNeeded={handleRefreshNeeded} 
        />
      </div>
    </div>
  );
};

export default DocumentManager;