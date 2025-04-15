import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { Engagement } from '../types';
import { engagementApi, expertRequestsApi } from '../services/api';
import Button from './ui/Button';
import Input from './ui/Input';

interface EngagementFormProps {
  expertId: number;
  engagement?: Engagement;
  onSuccess: (engagement: Engagement) => void;
  onCancel: () => void;
}

type EngagementFormData = {
  title: string;
  description: string;
  engagementType: string;
  status: string;
  startDate: string;
  endDate: string;
  contactPerson: string;
  contactEmail: string;
  organizationName: string;
  notes: string;
  requestId?: string;
};

const EngagementForm = ({ expertId, engagement, onSuccess, onCancel }: EngagementFormProps) => {
  const isEditMode = !!engagement;
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [requestOptions, setRequestOptions] = useState<{ value: string; label: string }[]>([]);
  
  const { register, handleSubmit, formState: { errors }, reset } = useForm<EngagementFormData>({
    defaultValues: isEditMode
      ? {
          title: engagement.title,
          description: engagement.description,
          engagementType: engagement.engagementType,
          status: engagement.status,
          startDate: engagement.startDate.split('T')[0],
          endDate: engagement.endDate.split('T')[0],
          contactPerson: engagement.contactPerson,
          contactEmail: engagement.contactEmail,
          organizationName: engagement.organizationName,
          notes: engagement.notes,
          requestId: engagement.requestId?.toString() || ''
        }
      : {
          title: '',
          description: '',
          engagementType: 'consultation',
          status: 'pending',
          startDate: new Date().toISOString().split('T')[0],
          endDate: new Date(Date.now() + 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
          contactPerson: '',
          contactEmail: '',
          organizationName: '',
          notes: '',
          requestId: ''
        }
  });

  // Fetch expert requests for dropdown
  useEffect(() => {
    const fetchRequests = async () => {
      try {
        const response = await expertRequestsApi.getExpertRequests();
        if (response.success) {
          const options = response.data.data.map(request => ({
            value: request.id.toString(),
            label: `${request.projectName} (${request.organizationName})`
          }));
          setRequestOptions([{ value: '', label: 'None (Direct Engagement)' }, ...options]);
        }
      } catch (error) {
        console.error('Error fetching expert requests:', error);
      }
    };
    
    fetchRequests();
  }, []);

  const onSubmit = async (data: EngagementFormData) => {
    setIsSubmitting(true);
    setError(null);
    
    try {
      // Transform form data to API format
      const engagementData: Partial<Engagement> = {
        ...data,
        expertId,
        requestId: data.requestId ? Number(data.requestId) : null
      };
      
      let response;
      if (isEditMode && engagement) {
        response = await engagementApi.updateEngagement(engagement.id.toString(), engagementData);
      } else {
        response = await engagementApi.createEngagement(engagementData);
      }
      
      if (response.success) {
        reset();
        onSuccess(response.data);
      } else {
        setError(response.message || `Failed to ${isEditMode ? 'update' : 'create'} engagement`);
      }
    } catch (error) {
      console.error(`Error ${isEditMode ? 'updating' : 'creating'} engagement:`, error);
      setError(`An error occurred while ${isEditMode ? 'updating' : 'creating'} the engagement`);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      {error && (
        <div className="bg-red-50 border border-red-200 text-red-600 p-4 rounded-md mb-4">
          {error}
        </div>
      )}
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Title</label>
          <Input
            {...register('title', { required: 'Title is required' })}
            error={errors.title?.message}
          />
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Organization</label>
          <Input
            {...register('organizationName', { required: 'Organization is required' })}
            error={errors.organizationName?.message}
          />
        </div>
      </div>
      
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Description</label>
        <textarea
          {...register('description', { required: 'Description is required' })}
          className={`w-full px-3 py-2 border rounded-md ${
            errors.description ? 'border-red-500' : 'border-gray-300'
          }`}
          rows={3}
        />
        {errors.description && (
          <p className="mt-1 text-sm text-red-600">{errors.description.message}</p>
        )}
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Engagement Type</label>
          <select
            {...register('engagementType', { required: 'Engagement type is required' })}
            className={`w-full px-3 py-2 border rounded-md ${
              errors.engagementType ? 'border-red-500' : 'border-gray-300'
            }`}
          >
            <option value="consultation">Consultation</option>
            <option value="project">Project</option>
            <option value="workshop">Workshop</option>
            <option value="training">Training</option>
            <option value="research">Research</option>
            <option value="other">Other</option>
          </select>
          {errors.engagementType && (
            <p className="mt-1 text-sm text-red-600">{errors.engagementType.message}</p>
          )}
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
          <select
            {...register('status', { required: 'Status is required' })}
            className={`w-full px-3 py-2 border rounded-md ${
              errors.status ? 'border-red-500' : 'border-gray-300'
            }`}
          >
            <option value="pending">Pending</option>
            <option value="confirmed">Confirmed</option>
            <option value="in_progress">In Progress</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
          </select>
          {errors.status && (
            <p className="mt-1 text-sm text-red-600">{errors.status.message}</p>
          )}
        </div>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Start Date</label>
          <Input
            type="date"
            {...register('startDate', { required: 'Start date is required' })}
            error={errors.startDate?.message}
          />
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">End Date</label>
          <Input
            type="date"
            {...register('endDate', { required: 'End date is required' })}
            error={errors.endDate?.message}
          />
        </div>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Contact Person</label>
          <Input
            {...register('contactPerson', { required: 'Contact person is required' })}
            error={errors.contactPerson?.message}
          />
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">Contact Email</label>
          <Input
            type="email"
            {...register('contactEmail', { 
              required: 'Email is required',
              pattern: {
                value: /^[A-Z0-9._%+-]+@[A-Z0-9.-]+\.[A-Z]{2,}$/i,
                message: 'Invalid email address'
              }
            })}
            error={errors.contactEmail?.message}
          />
        </div>
      </div>
      
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Related Expert Request (Optional)</label>
        <select
          {...register('requestId')}
          className="w-full px-3 py-2 border border-gray-300 rounded-md"
        >
          {requestOptions.map(option => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      </div>
      
      <div>
        <label className="block text-sm font-medium text-gray-700 mb-1">Notes (Optional)</label>
        <textarea
          {...register('notes')}
          className="w-full px-3 py-2 border border-gray-300 rounded-md"
          rows={3}
        />
      </div>
      
      <div className="flex justify-end space-x-3">
        <Button
          type="button"
          variant="outline"
          onClick={onCancel}
          disabled={isSubmitting}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          variant="primary"
          disabled={isSubmitting}
        >
          {isSubmitting 
            ? (isEditMode ? 'Updating...' : 'Creating...') 
            : (isEditMode ? 'Update Engagement' : 'Create Engagement')}
        </Button>
      </div>
    </form>
  );
};

export default EngagementForm;