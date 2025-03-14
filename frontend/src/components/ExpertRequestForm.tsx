import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { expertRequestsApi } from '../services/api';
import Input from './ui/Input';
import Button from './ui/Button';

interface ExpertRequestFormData {
  name: string;
  designation: string;
  institution: string;
  role: string;
  employmentType: string;
  generalArea: string;
  specializedArea: string;
  isBahraini: boolean;
  isAvailable: boolean;
  isTrained: boolean;
  email: string;
  phone: string;
}

interface ExpertRequestFormProps {
  onSuccess: () => void;
}

const ExpertRequestForm = ({ onSuccess }: ExpertRequestFormProps) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [cvFile, setCvFile] = useState<File | null>(null);
  
  const { 
    register, 
    handleSubmit, 
    formState: { errors },
    reset
  } = useForm<ExpertRequestFormData>();
  
  const roles = [
    { value: 'evaluator', label: 'Evaluator' },
    { value: 'validator', label: 'Validator' },
    { value: 'consultant', label: 'Consultant' },
    { value: 'trainer', label: 'Trainer' },
    { value: 'expert', label: 'Expert' }
  ];
  
  const employmentTypes = [
    { value: 'academic', label: 'Academic' },
    { value: 'employer', label: 'Employer' },
    { value: 'freelance', label: 'Freelance' },
    { value: 'government', label: 'Government' },
    { value: 'other', label: 'Other' }
  ];
  
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    setCvFile(file);
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
      
      // Add CV file if available
      if (cvFile) {
        formData.append('cvFile', cvFile);
      } else {
        setError('CV file is required');
        setIsSubmitting(false);
        return;
      }
      
      const response = await expertRequestsApi.createExpertRequest(formData);
      
      if (response.success) {
        reset();
        setCvFile(null);
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
          label="Name *"
          error={errors.name?.message}
          {...register('name', { required: 'Name is required' })}
        />
        
        <Input
          label="Designation/Title *"
          error={errors.designation?.message}
          {...register('designation', { required: 'Designation is required' })}
        />
        
        <Input
          label="Institution/Affiliation *"
          error={errors.institution?.message}
          {...register('institution', { required: 'Institution is required' })}
        />
        
        <div className="mb-4">
          <label className="block text-sm font-medium text-neutral-700 mb-1">
            Role *
          </label>
          <select
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            {...register('role', { required: 'Role is required' })}
          >
            <option value="">Select a role</option>
            {roles.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
          {errors.role && (
            <p className="mt-1 text-sm text-secondary">{errors.role.message}</p>
          )}
        </div>
        
        <div className="mb-4">
          <label className="block text-sm font-medium text-neutral-700 mb-1">
            Employment Type *
          </label>
          <select
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            {...register('employmentType', { required: 'Employment type is required' })}
          >
            <option value="">Select employment type</option>
            {employmentTypes.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
          {errors.employmentType && (
            <p className="mt-1 text-sm text-secondary">{errors.employmentType.message}</p>
          )}
        </div>
        
        <Input
          label="General Area *"
          placeholder="e.g. Education, Engineering, Medicine"
          error={errors.generalArea?.message}
          {...register('generalArea', { required: 'General area is required' })}
        />
        
        <Input
          label="Specialized Area *"
          placeholder="e.g. Math Education, Civil Engineering"
          error={errors.specializedArea?.message}
          {...register('specializedArea', { required: 'Specialized area is required' })}
        />
        
        <div className="mb-4">
          <label className="block text-sm font-medium text-neutral-700 mb-1">
            Contact Email *
          </label>
          <input
            type="email"
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            {...register('email', { 
              required: 'Email is required',
              pattern: {
                value: /\S+@\S+\.\S+/,
                message: 'Invalid email format',
              },
            })}
          />
          {errors.email && (
            <p className="mt-1 text-sm text-secondary">{errors.email.message}</p>
          )}
        </div>
        
        <Input
          label="Contact Phone"
          type="tel"
          error={errors.phone?.message}
          {...register('phone')}
        />
        
        <div className="mb-4">
          <label className="block text-sm font-medium text-neutral-700 mb-1">
            CV Upload *
          </label>
          <input
            type="file"
            accept=".pdf,.doc,.docx"
            onChange={handleFileChange}
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
          />
          {!cvFile && (
            <p className="mt-1 text-sm text-neutral-500">
              Upload CV in PDF, DOC, or DOCX format (max 10MB)
            </p>
          )}
          {cvFile && (
            <p className="mt-1 text-sm text-green-600">
              File selected: {cvFile.name}
            </p>
          )}
        </div>
      </div>
      
      <div className="space-y-3">
        <div className="flex items-center">
          <input
            type="checkbox"
            id="isBahraini"
            className="h-4 w-4 text-primary focus:ring-primary border-neutral-300 rounded"
            {...register('isBahraini')}
          />
          <label htmlFor="isBahraini" className="ml-2 block text-sm text-neutral-700">
            Bahraini Citizen
          </label>
        </div>
        
        <div className="flex items-center">
          <input
            type="checkbox"
            id="isAvailable"
            className="h-4 w-4 text-primary focus:ring-primary border-neutral-300 rounded"
            {...register('isAvailable')}
          />
          <label htmlFor="isAvailable" className="ml-2 block text-sm text-neutral-700">
            Available for Engagements
          </label>
        </div>
        
        <div className="flex items-center">
          <input
            type="checkbox"
            id="isTrained"
            className="h-4 w-4 text-primary focus:ring-primary border-neutral-300 rounded"
            {...register('isTrained')}
          />
          <label htmlFor="isTrained" className="ml-2 block text-sm text-neutral-700">
            Has Received BQA Training
          </label>
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