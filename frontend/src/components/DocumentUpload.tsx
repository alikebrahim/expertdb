import { useState } from 'react';
import { documentApi } from '../services/api';
import Button from './ui/Button';
import { Document } from '../types';

interface DocumentUploadProps {
  expertId: number;
  onSuccess: (document: Document) => void;
}

type DocumentType = 'cv' | 'certificate' | 'research' | 'publication' | 'other';

const DocumentUpload = ({ expertId, onSuccess }: DocumentUploadProps) => {
  const [file, setFile] = useState<File | null>(null);
  const [documentType, setDocumentType] = useState<DocumentType>('other');
  const [isUploading, setIsUploading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const documentTypeOptions = [
    { value: 'cv', label: 'Curriculum Vitae' },
    { value: 'certificate', label: 'Certification' },
    { value: 'research', label: 'Research Paper' },
    { value: 'publication', label: 'Publication' },
    { value: 'other', label: 'Other Document' },
  ];
  
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0] || null;
    setFile(selectedFile);
    setError(null);
  };
  
  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!file) {
      setError('Please select a file to upload');
      return;
    }
    
    setIsUploading(true);
    setError(null);
    
    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('expertId', expertId.toString());
      formData.append('documentType', documentType);
      
      const response = await documentApi.uploadDocument(formData);
      
      if (response.success) {
        setFile(null);
        setDocumentType('other');
        onSuccess(response.data);
      } else {
        setError(response.message || 'Failed to upload document');
      }
    } catch (error) {
      console.error('Error uploading document:', error);
      setError('An error occurred while uploading the document');
    } finally {
      setIsUploading(false);
    }
  };
  
  return (
    <form onSubmit={handleSubmit} className="p-4 border rounded-md bg-white">
      <h3 className="text-lg font-medium mb-4">Upload Document</h3>
      
      {error && (
        <div className="mb-4 p-3 text-sm bg-red-50 text-red-600 rounded">
          {error}
        </div>
      )}
      
      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Document Type
          </label>
          <select
            value={documentType}
            onChange={(e) => setDocumentType(e.target.value as DocumentType)}
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            required
          >
            {documentTypeOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            File
          </label>
          <input
            type="file"
            onChange={handleFileChange}
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            accept=".pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.jpg,.jpeg,.png"
            required
          />
          <p className="mt-1 text-sm text-gray-500">
            Accepted formats: PDF, Word, Excel, PowerPoint, JPEG, PNG (max 10MB)
          </p>
        </div>
        
        {file && (
          <div className="text-sm text-green-600">
            Selected file: {file.name} ({(file.size / 1024 / 1024).toFixed(2)} MB)
          </div>
        )}
      </div>
      
      <div className="mt-4">
        <Button
          type="submit"
          isLoading={isUploading}
          disabled={!file || isUploading}
          fullWidth
        >
          Upload Document
        </Button>
      </div>
    </form>
  );
};

export default DocumentUpload;