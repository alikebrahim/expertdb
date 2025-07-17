import { useState, useEffect } from 'react';
import { Engagement, EngagementListResponse } from '../types';
import { engagementApi, expertRequestsApi } from '../services/api';
import { z } from 'zod';
import { useFormWithNotifications } from '../hooks/useForm';
import { Form } from './ui/Form';
import { FormField } from './ui/FormField';
import { LoadingOverlay } from './ui/LoadingSpinner';
import Button from './ui/Button';

interface EngagementFormProps {
  expertId: number;
  engagement?: Engagement;
  onSuccess: (engagement: Engagement) => void;
  onCancel: () => void;
}

// Engagement form schema
const engagementFormSchema = z.object({
  title: z.string().min(2, 'Title is required'),
  description: z.string().min(10, 'Description must be at least 10 characters'),
  engagementType: z.enum(['validator', 'evaluator']),
  status: z.enum(['pending', 'confirmed', 'in_progress', 'completed', 'cancelled']),
  startDate: z.string().min(1, 'Start date is required'),
  endDate: z.string().min(1, 'End date is required'),
  contactPerson: z.string().min(2, 'Contact person is required'),
  contactEmail: z.string().email('Please provide a valid email address'),
  organizationName: z.string().min(2, 'Organization name is required'),
  notes: z.string().optional(),
  requestId: z.string().optional(),
});

type EngagementFormData = z.infer<typeof engagementFormSchema>;

const EngagementForm = ({ expertId, engagement, onSuccess, onCancel }: EngagementFormProps) => {
  const isEditMode = !!engagement;
  const [isLoading, setIsLoading] = useState(false);
  const [requestOptions, setRequestOptions] = useState<{ value: string; label: string }[]>([]);
  
  // Define engagement type options
  const engagementTypeOptions = [
    { value: 'validator', label: 'Validator' },
    { value: 'evaluator', label: 'Evaluator' },
  ];
  
  // Define status options
  const statusOptions = [
    { value: 'pending', label: 'Pending' },
    { value: 'confirmed', label: 'Confirmed' },
    { value: 'in_progress', label: 'In Progress' },
    { value: 'completed', label: 'Completed' },
    { value: 'cancelled', label: 'Cancelled' }
  ];
  
  // Initialize form with default values
  const form = useFormWithNotifications<EngagementFormData>({
    schema: engagementFormSchema,
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
          engagementType: 'validator',
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
        if (response.success && response.data) {
          // Handle both array and paginated response formats
          const requests = Array.isArray(response.data) ? response.data : (response.data as EngagementListResponse).engagements || [];
          const options = requests.map((request: any) => ({
            value: request.id.toString(),
            label: `${request.name || request.organizationName} (${request.institution || request.projectName})`
          }));
          setRequestOptions([{ value: '', label: 'None (Direct Engagement)' }, ...options]);
        }
      } catch (error) {
        console.error('Error fetching expert requests:', error);
      }
    };
    
    fetchRequests();
  }, []);

  const onSubmit = async (data: EngagementFormData): Promise<void> => {
    setIsLoading(true);
    
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
        onSuccess(response.data);
        // Success is handled through onSuccess callback
      } else {
        throw new Error(response.message || `Failed to ${isEditMode ? 'update' : 'create'} engagement`);
      }
    } catch (error) {
      console.error(`Error ${isEditMode ? 'updating' : 'creating'} engagement:`, error);
      // Error will be handled by form's error handling
      throw error;
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <LoadingOverlay 
      isLoading={isLoading} 
      className="w-full"
      label={isEditMode ? "Updating engagement..." : "Creating engagement..."}
    >
      <Form
        form={form}
        onSubmit={onSubmit}
        className="space-y-4"
        submitText={isEditMode ? 'Update Engagement' : 'Create Engagement'}
      >
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField
            form={form}
            name="title"
            label="Title"
            required
          />
          
          <FormField
            form={form}
            name="organizationName"
            label="Organization"
            required
          />
        </div>
        
        <FormField
          form={form}
          name="description"
          label="Description"
          type="textarea"
          rows={3}
          required
        />
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField
            form={form}
            name="engagementType"
            label="Engagement Type"
            type="select"
            options={engagementTypeOptions}
            required
          />
          
          <FormField
            form={form}
            name="status"
            label="Status"
            type="select"
            options={statusOptions}
            required
          />
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField
            form={form}
            name="startDate"
            label="Start Date"
            type="date"
            required
          />
          
          <FormField
            form={form}
            name="endDate"
            label="End Date"
            type="date"
            required
          />
        </div>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField
            form={form}
            name="contactPerson"
            label="Contact Person"
            required
          />
          
          <FormField
            form={form}
            name="contactEmail"
            label="Contact Email"
            type="email"
            required
          />
        </div>
        
        <FormField
          form={form}
          name="requestId"
          label="Related Expert Request (Optional)"
          type="select"
          options={requestOptions}
        />
        
        <FormField
          form={form}
          name="notes"
          label="Notes (Optional)"
          type="textarea"
          rows={3}
          hint="Any additional information about this engagement"
        />
        
        <div className="flex justify-end space-x-3 mt-6">
          <Button
            type="button"
            variant="outline"
            onClick={onCancel}
          >
            Cancel
          </Button>
        </div>
      </Form>
    </LoadingOverlay>
  );
};

export default EngagementForm;