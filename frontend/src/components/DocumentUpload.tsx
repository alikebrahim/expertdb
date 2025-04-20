import { useState } from 'react';
import { documentApi } from '../services/api';
import { Document } from '../types';
import { useFormWithNotifications } from '../hooks/useForm';
import { z } from 'zod';
import { Form } from './ui/Form';
import { FormField } from './ui/FormField';
import { LoadingOverlay } from './ui/LoadingSpinner';

interface DocumentUploadProps {
  expertId: number;
  onSuccess: (document: Document) => void;
}

// Document types
const DOCUMENT_TYPES = ['cv', 'certificate', 'research', 'publication', 'other'] as const;
type DocumentType = typeof DOCUMENT_TYPES[number];

// Document upload schema
const documentUploadSchema = z.object({
  documentType: z.enum(DOCUMENT_TYPES),
  description: z.string().optional(),
});

type DocumentUploadFormData = z.infer<typeof documentUploadSchema>;

const DocumentUpload = ({ expertId, onSuccess }: DocumentUploadProps) => {
  const [file, setFile] = useState<File | null>(null);
  const [isUploading, setIsUploading] = useState(false);
  
  // Document type options for the select field
  const documentTypeOptions = [
    { value: 'cv', label: 'Curriculum Vitae' },
    { value: 'certificate', label: 'Certification' },
    { value: 'research', label: 'Research Paper' },
    { value: 'publication', label: 'Publication' },
    { value: 'other', label: 'Other Document' },
  ];
  
  const form = useFormWithNotifications<DocumentUploadFormData>({
    schema: documentUploadSchema,
    defaultValues: {
      documentType: 'other',
      description: '',
    }
  });
  
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const selectedFile = e.target.files?.[0] || null;
    setFile(selectedFile);
  };
  
  const handleFormReset = () => {
    setFile(null);
  };
  
  const onSubmit = async (data: DocumentUploadFormData) => {
    if (!file) {
      return { success: false, message: 'Please select a file to upload' };
    }
    
    setIsUploading(true);
    
    try {
      const formData = new FormData();
      formData.append('file', file);
      formData.append('expertId', expertId.toString());
      formData.append('documentType', data.documentType);
      
      if (data.description) {
        formData.append('description', data.description);
      }
      
      const response = await documentApi.uploadDocument(formData);
      
      if (response.success) {
        setFile(null);
        form.reset();
        onSuccess(response.data);
        return { success: true, message: 'Document uploaded successfully!' };
      } else {
        return { success: false, message: response.message || 'Failed to upload document' };
      }
    } catch (error) {
      console.error('Error uploading document:', error);
      return { 
        success: false, 
        message: error instanceof Error 
          ? error.message 
          : 'An error occurred while uploading the document'
      };
    } finally {
      setIsUploading(false);
    }
  };
  
  return (
    <LoadingOverlay 
      isLoading={isUploading}
      className="p-4 border rounded-md bg-white w-full"
      label="Uploading document..."
    >
      <h3 className="text-lg font-medium mb-4">Upload Document</h3>
      
      <Form
        form={form}
        onSubmit={form.handleSubmitWithNotifications(onSubmit)}
        className="space-y-4"
        onReset={handleFormReset}
        showResetButton={!!file}
        resetText="Clear"
        submitText="Upload Document"
      >
        <FormField
          form={form}
          name="documentType"
          label="Document Type"
          type="select"
          options={documentTypeOptions}
          required
        />
        
        <FormField
          form={form}
          name="description"
          label="Description (Optional)"
          type="textarea"
          rows={3}
          placeholder="Brief description of the document..."
          hint="Add a short description to help identify the document"
        />
        
        <div className="mb-4">
          <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
            File <span className="text-red-500">*</span>
          </label>
          <input
            type="file"
            onChange={handleFileChange}
            className="w-full px-3 py-2 bg-white border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary dark:bg-gray-800 dark:border-gray-600 dark:text-white"
            accept=".pdf,.doc,.docx,.xls,.xlsx,.ppt,.pptx,.jpg,.jpeg,.png"
            required
          />
          <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
            Accepted formats: PDF, Word, Excel, PowerPoint, JPEG, PNG (max 10MB)
          </p>
          {file && (
            <p className="mt-1 text-sm text-green-600 dark:text-green-400">
              Selected file: {file.name} ({(file.size / 1024 / 1024).toFixed(2)} MB)
            </p>
          )}
        </div>
      </Form>
    </LoadingOverlay>
  );
};

export default DocumentUpload;