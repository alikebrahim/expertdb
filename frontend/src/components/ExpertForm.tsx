import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { Expert } from '../types';
import { expertsApi } from '../services/api';
import Input from './ui/Input';
import Button from './ui/Button';

interface ExpertFormProps {
  expert?: Expert;
  onSuccess: (expert: Expert) => void;
  onCancel: () => void;
}

interface ExpertFormData {
  name: string;
  affiliation: string;
  primaryContact: string;
  contactType: string;
  skills: string;
  role: string;
  employmentType: string;
  generalArea: string;
  biography: string;
  isBahraini: boolean;
  availability: string;
}

const ExpertForm = ({ expert, onSuccess, onCancel }: ExpertFormProps) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [cvFile, setCvFile] = useState<File | null>(null);
  
  const isEditMode = !!expert;
  
  const { 
    register, 
    handleSubmit, 
    formState: { errors },
    reset,
    setValue,
  } = useForm<ExpertFormData>();
  
  useEffect(() => {
    if (expert) {
      setValue('name', expert.name);
      setValue('affiliation', expert.affiliation);
      setValue('primaryContact', expert.primaryContact);
      setValue('contactType', expert.contactType);
      setValue('skills', expert.skills.join(', '));
      setValue('role', expert.role);
      setValue('employmentType', expert.employmentType);
      setValue('generalArea', expert.generalArea.toString());
      setValue('biography', expert.biography);
      setValue('isBahraini', expert.isBahraini);
      setValue('availability', expert.availability);
    }
  }, [expert, setValue]);
  
  const contactTypeOptions = [
    { value: 'email', label: 'Email' },
    { value: 'phone', label: 'Phone' },
    { value: 'linkedin', label: 'LinkedIn' },
  ];
  
  const employmentTypeOptions = [
    { value: 'full-time', label: 'Full-Time' },
    { value: 'part-time', label: 'Part-Time' },
    { value: 'consultant', label: 'Consultant' },
    { value: 'retired', label: 'Retired' },
    { value: 'other', label: 'Other' },
  ];
  
  const availabilityOptions = [
    { value: 'Available', label: 'Available' },
    { value: 'Limited', label: 'Limited Availability' },
    { value: 'Unavailable', label: 'Currently Unavailable' },
  ];
  
  const roleOptions = [
    { value: 'Professor', label: 'Professor' },
    { value: 'Doctor', label: 'Doctor' },
    { value: 'Engineer', label: 'Engineer' },
    { value: 'Researcher', label: 'Researcher' },
    { value: 'Consultant', label: 'Consultant' },
    { value: 'Teacher', label: 'Teacher' },
    { value: 'Other', label: 'Other' },
  ];
  
  const handleCvFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    setCvFile(file);
  };
  
  const onSubmit = async (data: ExpertFormData) => {
    setIsSubmitting(true);
    setError(null);
    
    try {
      // Create FormData for file upload
      const formData = new FormData();
      
      // Process skills from comma-separated string to array
      const skillsArray = data.skills.split(',').map(skill => skill.trim());
      
      // Add all form fields to FormData
      Object.entries(data).forEach(([key, value]) => {
        if (key === 'skills') {
          // Skip skills as we'll handle them separately
          return;
        }
        formData.append(key, value.toString());
      });
      
      // Add skills as JSON array
      formData.append('skills', JSON.stringify(skillsArray));
      
      // Add CV file if available (for new experts or if updating the CV)
      if (cvFile) {
        formData.append('cvFile', cvFile);
      }
      
      let response;
      if (isEditMode && expert) {
        response = await expertsApi.updateExpert(expert.id.toString(), formData);
      } else {
        response = await expertsApi.createExpert(formData);
      }
      
      if (response.success) {
        reset();
        setCvFile(null);
        onSuccess(response.data);
      } else {
        setError(response.message || `Failed to ${isEditMode ? 'update' : 'create'} expert`);
      }
    } catch (error) {
      console.error(`Error ${isEditMode ? 'updating' : 'creating'} expert:`, error);
      setError(`An error occurred while ${isEditMode ? 'updating' : 'creating'} the expert`);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      <h2 className="text-xl font-bold mb-4">
        {isEditMode ? 'Edit Expert' : 'Create New Expert'}
      </h2>
      
      {error && (
        <div className="bg-red-50 text-red-600 p-3 rounded">
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
          label="Affiliation *"
          error={errors.affiliation?.message}
          {...register('affiliation', { required: 'Affiliation is required' })}
        />
        
        <Input
          label="Primary Contact *"
          error={errors.primaryContact?.message}
          {...register('primaryContact', { required: 'Primary contact is required' })}
        />
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Contact Type *
          </label>
          <select
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            {...register('contactType', { required: 'Contact type is required' })}
          >
            <option value="">Select contact type</option>
            {contactTypeOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
          {errors.contactType && (
            <p className="mt-1 text-sm text-red-600">{errors.contactType.message}</p>
          )}
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Role *
          </label>
          <select
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            {...register('role', { required: 'Role is required' })}
          >
            <option value="">Select role</option>
            {roleOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
          {errors.role && (
            <p className="mt-1 text-sm text-red-600">{errors.role.message}</p>
          )}
        </div>
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Employment Type *
          </label>
          <select
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            {...register('employmentType', { required: 'Employment type is required' })}
          >
            <option value="">Select employment type</option>
            {employmentTypeOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
          {errors.employmentType && (
            <p className="mt-1 text-sm text-red-600">{errors.employmentType.message}</p>
          )}
        </div>
        
        <Input
          label="General Area ID *"
          type="number"
          min="1"
          error={errors.generalArea?.message}
          {...register('generalArea', { required: 'General area is required' })}
        />
        
        <div>
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Availability *
          </label>
          <select
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            {...register('availability', { required: 'Availability is required' })}
          >
            <option value="">Select availability</option>
            {availabilityOptions.map((option) => (
              <option key={option.value} value={option.value}>
                {option.label}
              </option>
            ))}
          </select>
          {errors.availability && (
            <p className="mt-1 text-sm text-red-600">{errors.availability.message}</p>
          )}
        </div>
        
        <div className="flex items-center">
          <input
            type="checkbox"
            id="isBahraini"
            className="h-4 w-4 text-primary focus:ring-primary border-gray-300 rounded"
            {...register('isBahraini')}
          />
          <label htmlFor="isBahraini" className="ml-2 block text-sm text-gray-700">
            Is Bahraini Citizen
          </label>
        </div>
        
        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Skills * (comma-separated)
          </label>
          <textarea
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            rows={2}
            placeholder="E.g., Machine Learning, Data Science, Python, SQL"
            {...register('skills', { required: 'Skills are required' })}
          ></textarea>
          {errors.skills && (
            <p className="mt-1 text-sm text-red-600">{errors.skills.message}</p>
          )}
        </div>
        
        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-gray-700 mb-1">
            Biography
          </label>
          <textarea
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            rows={4}
            placeholder="Professional biography..."
            {...register('biography')}
          ></textarea>
        </div>
        
        <div className="md:col-span-2">
          <label className="block text-sm font-medium text-gray-700 mb-1">
            {isEditMode ? 'Update CV (Optional)' : 'CV File *'}
          </label>
          <input
            type="file"
            accept=".pdf,.doc,.docx"
            onChange={handleCvFileChange}
            className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
            required={!isEditMode}
          />
          <p className="mt-1 text-sm text-gray-500">
            Upload CV in PDF or Word format (max 10MB)
          </p>
          {cvFile && (
            <p className="mt-1 text-sm text-green-600">
              File selected: {cvFile.name} ({(cvFile.size / 1024 / 1024).toFixed(2)} MB)
            </p>
          )}
          {isEditMode && !cvFile && (
            <p className="mt-1 text-sm text-gray-600">
              Current CV: {expert?.cvPath ? expert.cvPath.split('/').pop() : 'None'}
            </p>
          )}
        </div>
      </div>
      
      <div className="flex justify-end space-x-3 mt-6">
        <Button 
          type="button" 
          variant="outline" 
          onClick={onCancel}
        >
          Cancel
        </Button>
        <Button 
          type="submit" 
          isLoading={isSubmitting}
        >
          {isEditMode ? 'Update Expert' : 'Create Expert'}
        </Button>
      </div>
    </form>
  );
};

export default ExpertForm;