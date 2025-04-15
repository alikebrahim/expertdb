import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { expertRequestsApi } from '../services/api';
import Input from './ui/Input';
import Button from './ui/Button';

// Updated form data interface to match API documentation
interface ExpertRequestFormData {
  organizationName: string;
  projectName: string;
  projectDescription: string;
  expertiseRequired: string;
  timeframe: string;
  notes: string;
}

interface ExpertRequestFormProps {
  onSuccess: () => void;
}

const ExpertRequestForm = ({ onSuccess }: ExpertRequestFormProps) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [attachmentFile, setAttachmentFile] = useState<File | null>(null);
  
  const { 
    register, 
    handleSubmit, 
    formState: { errors },
    reset
  } = useForm<ExpertRequestFormData>();
  
  const timeframeOptions = [
    { value: 'urgent', label: 'Urgent (1-2 weeks)' },
    { value: 'short', label: 'Short-term (1-3 months)' },
    { value: 'medium', label: 'Medium-term (3-6 months)' },
    { value: 'long', label: 'Long-term (6+ months)' }
  ];
  
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    setAttachmentFile(file);
  };
  
  const onSubmit = async (data: ExpertRequestFormData) => {
    setIsSubmitting(true);
    setError(null);
    
    try {
      // Create FormData for file upload
      const formData = new FormData();
      
      // Add all form fields to FormData
      Object.entries(data).forEach(([key, value]) => {
        formData.append(key, value.toString());
      });
      
      // Add attachment file if available
      if (attachmentFile) {
        formData.append('attachmentFile', attachmentFile);
      }
      
      const response = await expertRequestsApi.createExpertRequest(formData);
      
      if (response.success) {
        reset();
        setAttachmentFile(null);
        onSuccess();
      } else {
        setError(response.message || 'Failed to submit expert request');
      }
    } catch (error) {
      console.error('Error submitting expert request:', error);
      setError('An error occurred while submitting the request');
    } finally {
      setIsSubmitting(false);
    }
  };
  
  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      {error && (
        <div className="bg-secondary bg-opacity-10 text-secondary p-3 rounded">
          {error}
        </div>
      )}
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Input
          label="Organization Name *"
          error={errors.organizationName?.message}
          {...register('organizationName', { required: 'Organization name is required' })}
        />
        
        <Input
          label="Project Name *"
          error={errors.projectName?.message}
          {...register('projectName', { required: 'Project name is required' })}
        />
        
        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-neutral-700 mb-1">
            Project Description *
          </label>
          <textarea
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            rows={4}
            {...register('projectDescription', { required: 'Project description is required' })}
          ></textarea>
          {errors.projectDescription && (
            <p className="mt-1 text-sm text-secondary">{errors.projectDescription.message}</p>
          )}
        </div>
        
        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-neutral-700 mb-1">
            Expertise Required *
          </label>
          <textarea
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            rows={3}
            placeholder="Describe the expertise you are looking for..."
            {...register('expertiseRequired', { required: 'Required expertise information is required' })}
          ></textarea>
          {errors.expertiseRequired && (
            <p className="mt-1 text-sm text-secondary">{errors.expertiseRequired.message}</p>
          )}
        </div>
        
        <div className="mb-4">
          <label className="block text-sm font-medium text-neutral-700 mb-1">
            Timeframe *
          </label>
          <select
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            {...register('timeframe', { required: 'Timeframe is required' })}
          >
            <option value="">Select timeframe</option>
            {timeframeOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
          {errors.timeframe && (
            <p className="mt-1 text-sm text-secondary">{errors.timeframe.message}</p>
          )}
        </div>
        
        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-neutral-700 mb-1">
            Additional Notes
          </label>
          <textarea
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            rows={3}
            placeholder="Any additional information that might help us match you with the right expert..."
            {...register('notes')}
          ></textarea>
        </div>
        
        <div className="mb-4">
          <label className="block text-sm font-medium text-neutral-700 mb-1">
            Supporting Document (Optional)
          </label>
          <input
            type="file"
            accept=".pdf,.doc,.docx,.ppt,.pptx"
            onChange={handleFileChange}
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
          />
          <p className="mt-1 text-sm text-neutral-500">
            Upload relevant project documents (max 10MB)
          </p>
          {attachmentFile && (
            <p className="mt-1 text-sm text-green-600">
              File selected: {attachmentFile.name}
            </p>
          )}
        </div>
      </div>
      
      <div className="mt-6 flex justify-end space-x-3">
        <Button 
          type="button" 
          variant="outline" 
          onClick={() => reset()}
        >
          Reset
        </Button>
        <Button 
          type="submit" 
          isLoading={isSubmitting}
        >
          Submit Request
        </Button>
      </div>
    </form>
  );
};

export default ExpertRequestForm;