import { useState } from 'react';
import { expertRequestsApi } from '../services/api';
import { useFormWithNotifications } from '../hooks/useForm';
import { z } from 'zod';
import { Form } from './ui/Form';
import { FormField } from './ui/FormField';
import { LoadingOverlay } from './ui/LoadingSpinner';

// Expert request schema
const expertRequestFormSchema = z.object({
  organizationName: z.string().min(2, 'Organization name is required'),
  projectName: z.string().min(2, 'Project name is required'),
  projectDescription: z.string().min(10, 'Project description must be at least 10 characters'),
  expertiseRequired: z.string().min(10, 'Required expertise must be at least 10 characters'),
  timeframe: z.string().min(1, 'Timeframe is required'),
  notes: z.string().optional(),
});

type ExpertRequestFormData = z.infer<typeof expertRequestFormSchema>;

interface ExpertRequestFormProps {
  onSuccess: () => void;
}

const ExpertRequestForm = ({ onSuccess }: ExpertRequestFormProps) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [attachmentFile, setAttachmentFile] = useState<File | null>(null);
  
  // Time frame options
  const timeframeOptions = [
    { value: '', label: 'Select timeframe' },
    { value: 'urgent', label: 'Urgent (1-2 weeks)' },
    { value: 'short', label: 'Short-term (1-3 months)' },
    { value: 'medium', label: 'Medium-term (3-6 months)' },
    { value: 'long', label: 'Long-term (6+ months)' }
  ];
  
  const form = useFormWithNotifications<ExpertRequestFormData>({
    schema: expertRequestFormSchema,
    defaultValues: {
      organizationName: '',
      projectName: '',
      projectDescription: '',
      expertiseRequired: '',
      timeframe: '',
      notes: ''
    },
  });
  
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    setAttachmentFile(file);
  };
  
  const handleFormReset = () => {
    setAttachmentFile(null);
  };
  
  const onSubmit = async (data: ExpertRequestFormData) => {
    setIsSubmitting(true);
    
    try {
      // Create FormData for file upload
      const formData = new FormData();
      
      // Add all form fields to FormData
      Object.entries(data).forEach(([key, value]) => {
        if (value) {
          formData.append(key, value.toString());
        }
      });
      
      // Add attachment file if available
      if (attachmentFile) {
        formData.append('attachmentFile', attachmentFile);
      }
      
      const response = await expertRequestsApi.createExpertRequest(formData);
      
      if (response.success) {
        form.reset();
        setAttachmentFile(null);
        onSuccess();
        return { success: true, message: 'Expert request submitted successfully!' };
      } else {
        return { success: false, message: response.message || 'Failed to submit expert request' };
      }
    } catch (error) {
      console.error('Error submitting expert request:', error);
      return { 
        success: false, 
        message: error instanceof Error 
          ? error.message 
          : 'An unexpected error occurred while submitting the request' 
      };
    } finally {
      setIsSubmitting(false);
    }
  };
  
  return (
    <LoadingOverlay 
      isLoading={isSubmitting}
      className="w-full"
      label="Submitting request..."
    >
      <Form
        form={form}
        onSubmit={form.handleSubmitWithNotifications(onSubmit)}
        className="space-y-4"
        showResetButton={true}
        resetText="Reset"
        onReset={handleFormReset}
        submitText="Submit Request"
      >
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField
            form={form}
            name="organizationName"
            label="Organization Name"
            placeholder="Enter organization name"
            required
          />
          
          <FormField
            form={form}
            name="projectName"
            label="Project Name"
            placeholder="Enter project name"
            required
          />
          
          <div className="md:col-span-2">
            <FormField
              form={form}
              name="projectDescription"
              label="Project Description"
              type="textarea"
              placeholder="Describe your project..."
              rows={4}
              required
            />
          </div>
          
          <div className="md:col-span-2">
            <FormField
              form={form}
              name="expertiseRequired"
              label="Expertise Required"
              type="textarea"
              placeholder="Describe the expertise you are looking for..."
              rows={3}
              required
            />
          </div>
          
          <FormField
            form={form}
            name="timeframe"
            label="Timeframe"
            type="select"
            options={timeframeOptions}
            required
          />
          
          <div className="md:col-span-2">
            <FormField
              form={form}
              name="notes"
              label="Additional Notes"
              type="textarea"
              placeholder="Any additional information that might help us match you with the right expert..."
              rows={3}
              hint="Optional - include any other relevant details"
            />
          </div>
          
          <div className="mb-4">
            <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">
              Supporting Document (Optional)
            </label>
            <input
              type="file"
              accept=".pdf,.doc,.docx,.ppt,.pptx"
              onChange={handleFileChange}
              className="w-full px-3 py-2 bg-white border border-gray-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary dark:bg-gray-800 dark:border-gray-600 dark:text-white"
            />
            <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
              Upload relevant project documents (max 10MB)
            </p>
            {attachmentFile && (
              <p className="mt-1 text-sm text-green-600 dark:text-green-400">
                File selected: {attachmentFile.name}
              </p>
            )}
          </div>
        </div>
      </Form>
    </LoadingOverlay>
  );
};

export default ExpertRequestForm;