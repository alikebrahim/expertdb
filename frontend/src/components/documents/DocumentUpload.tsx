import React, { useState, useRef } from 'react';
import { useUI } from '../../hooks/useUI';
import Button from '../ui/Button';

interface DocumentUploadProps {
  expertId: number;
  onUploadSuccess: (documentId: number) => void;
  allowedTypes?: string[];
  maxSizeMB?: number;
}

const DocumentUpload: React.FC<DocumentUploadProps> = ({
  expertId,
  onUploadSuccess,
  allowedTypes = ['application/pdf', 'image/jpeg', 'image/png', 'application/msword', 'application/vnd.openxmlformats-officedocument.wordprocessingml.document'],
  maxSizeMB = 10
}) => {
  const [file, setFile] = useState<File | null>(null);
  const [docType, setDocType] = useState<string>('cv');
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [progress, setProgress] = useState<number>(0);
  const fileInputRef = useRef<HTMLInputElement>(null);
  const { addNotification } = useUI();

  const documentTypes = [
    { value: 'cv', label: 'CV/Resume' },
    { value: 'certificate', label: 'Certificate' },
    { value: 'publication', label: 'Publication' },
    { value: 'reference', label: 'Reference Letter' },
    { value: 'identification', label: 'ID Document' },
    { value: 'other', label: 'Other' }
  ];

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0] || null;
    setError(null);
    
    if (!selectedFile) {
      return;
    }
    
    // Validate file type
    if (!allowedTypes.includes(selectedFile.type)) {
      setError(`File type not allowed. Please upload one of the following: ${allowedTypes.join(', ')}`);
      setFile(null);
      return;
    }
    
    // Validate file size
    const maxSizeBytes = maxSizeMB * 1024 * 1024;
    if (selectedFile.size > maxSizeBytes) {
      setError(`File too large. Maximum size is ${maxSizeMB}MB.`);
      setFile(null);
      return;
    }
    
    setFile(selectedFile);
  };

  const handleUpload = async () => {
    if (!file || !expertId) {
      setError('Please select a file to upload.');
      return;
    }

    setIsUploading(true);
    setProgress(0);
    setError(null);

    const formData = new FormData();
    formData.append('file', file);
    formData.append('expertId', expertId.toString());
    formData.append('documentType', docType);

    let interval: number;
    
    try {
      // Simulate progress (in a real app, you'd use XMLHttpRequest with progress events)
      interval = setInterval(() => {
        setProgress(prev => {
          const newProgress = prev + 10;
          if (newProgress >= 90) {
            clearInterval(interval);
            return 90;
          }
          return newProgress;
        });
      }, 300);

      // This would be an actual API call in a real application
      // const response = await documentsApi.uploadDocument(formData);
      
      // For demo purposes, simulate a successful upload after a delay
      setTimeout(() => {
        clearInterval(interval);
        setProgress(100);
        
        // Simulate a successful response with document ID
        const mockDocumentId = Math.floor(Math.random() * 10000);
        onUploadSuccess(mockDocumentId);
        
        addNotification({
          type: 'success',
          message: 'Document uploaded successfully',
          duration: 3000,
        });
        
        // Reset form
        setFile(null);
        setProgress(0);
        setIsUploading(false);
        if (fileInputRef.current) {
          fileInputRef.current.value = '';
        }
      }, 2000);
    } catch (error) {
      if (interval) clearInterval(interval);
      console.error('Upload error:', error);
      setError('An error occurred while uploading the document. Please try again.');
      setIsUploading(false);
      setProgress(0);
      
      addNotification({
        type: 'error',
        message: 'Failed to upload document',
        duration: 5000,
      });
    }
  };

  return (
    <div className="bg-white p-4 rounded-md shadow">
      <h3 className="text-lg font-medium text-gray-900 mb-4">Upload Document</h3>
      
      {error && (
        <div className="mb-4 p-3 bg-red-50 text-red-700 rounded-md text-sm">
          {error}
        </div>
      )}
      
      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Document Type
          </label>
          <select
            value={docType}
            onChange={(e) => setDocType(e.target.value)}
            className="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-primary focus:border-primary"
            disabled={isUploading}
          >
            {documentTypes.map((type) => (
              <option key={type.value} value={type.value}>
                {type.label}
              </option>
            ))}
          </select>
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Select File
          </label>
          <div className="flex items-center">
            <input
              type="file"
              ref={fileInputRef}
              onChange={handleFileChange}
              disabled={isUploading}
              className="hidden"
              accept={allowedTypes.join(',')}
            />
            <Button
              type="button"
              variant="outline"
              onClick={() => fileInputRef.current?.click()}
              disabled={isUploading}
              className="mr-3"
            >
              Browse...
            </Button>
            <span className="text-sm text-gray-500 truncate max-w-xs">
              {file ? file.name : 'No file selected'}
            </span>
          </div>
          <p className="mt-1 text-xs text-gray-500">
            Allowed file types: PDF, JPEG, PNG, DOC, DOCX. Max size: {maxSizeMB}MB
          </p>
        </div>
        
        {isUploading && (
          <div className="mt-2">
            <div className="relative pt-1">
              <div className="flex mb-2 items-center justify-between">
                <div>
                  <span className="text-xs font-semibold inline-block py-1 px-2 uppercase rounded-full text-primary bg-primary bg-opacity-10">
                    Uploading
                  </span>
                </div>
                <div className="text-right">
                  <span className="text-xs font-semibold inline-block text-primary">
                    {progress}%
                  </span>
                </div>
              </div>
              <div className="overflow-hidden h-2 mb-4 text-xs flex rounded bg-primary bg-opacity-10">
                <div
                  style={{ width: `${progress}%` }}
                  className="shadow-none flex flex-col text-center whitespace-nowrap text-white justify-center bg-primary transition-all duration-500"
                ></div>
              </div>
            </div>
          </div>
        )}
        
        <div className="flex justify-end">
          <Button
            type="button"
            variant="primary"
            onClick={handleUpload}
            disabled={!file || isUploading}
            isLoading={isUploading}
          >
            Upload Document
          </Button>
        </div>
      </div>
    </div>
  );
};

export default DocumentUpload;