import { useState, useEffect } from 'react';
import { Expert } from '../types';
import { expertsApi } from '../services/api';
import { useFormWithNotifications } from '../hooks/useForm';
import { expertSchema } from '../utils/formSchemas';
import { Form, FormField, LoadingOverlay } from './ui';

interface ExpertFormProps {
  expert?: Expert;
  onSuccess: (expert: Expert) => void;
  onCancel: () => void;
}

interface ExpertFormData {
  name: string;
  affiliation: string;
  primaryContact: string;
  contactType: 'email' | 'phone' | 'linkedin';
  skills: string;
  role: string;
  employmentType: 'full-time' | 'part-time' | 'consultant' | 'retired' | 'other';
  generalArea: string;
  biography: string;
  isBahraini: boolean;
  availability: 'Available' | 'Limited' | 'Unavailable';
  cvFile?: File;
}

const ExpertForm = ({ expert, onSuccess, onCancel }: ExpertFormProps) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [cvFile, setCvFile] = useState<File | null>(null);
  
  const isEditMode = !!expert;
  
  const form = useFormWithNotifications<ExpertFormData>({
    schema: expertSchema,
    defaultValues: {
      name: '',
      affiliation: '',
      primaryContact: '',
      contactType: 'email',
      skills: '',
      role: '',
      employmentType: 'full-time',
      generalArea: '',
      biography: '',
      isBahraini: false,
      availability: 'Available',
    }
  });
  
  // Set form values when expert data is available (for edit mode)
  useEffect(() => {
    if (expert) {
      form.reset({
        name: expert.name,
        affiliation: expert.affiliation,
        primaryContact: expert.primaryContact,
        contactType: expert.contactType as 'email' | 'phone' | 'linkedin',
        skills: expert.skills.join(', '),
        role: expert.role,
        employmentType: expert.employmentType as 'full-time' | 'part-time' | 'consultant' | 'retired' | 'other',
        generalArea: expert.generalArea.toString(),
        biography: expert.biography,
        isBahraini: expert.isBahraini,
        availability: expert.availability as 'Available' | 'Limited' | 'Unavailable',
      });
    }
  }, [expert, form]);
  
  const contactTypeOptions = [
    { label: 'Email', value: 'email' },
    { label: 'Phone', value: 'phone' },
    { label: 'LinkedIn', value: 'linkedin' },
  ];
  
  const employmentTypeOptions = [
    { label: 'Full-Time', value: 'full-time' },
    { label: 'Part-Time', value: 'part-time' },
    { label: 'Consultant', value: 'consultant' },
    { label: 'Retired', value: 'retired' },
    { label: 'Other', value: 'other' },
  ];
  
  const availabilityOptions = [
    { label: 'Available', value: 'Available' },
    { label: 'Limited Availability', value: 'Limited' },
    { label: 'Currently Unavailable', value: 'Unavailable' },
  ];
  
  const roleOptions = [
    { label: 'Professor', value: 'Professor' },
    { label: 'Doctor', value: 'Doctor' },
    { label: 'Engineer', value: 'Engineer' },
    { label: 'Researcher', value: 'Researcher' },
    { label: 'Consultant', value: 'Consultant' },
    { label: 'Teacher', value: 'Teacher' },
    { label: 'Other', value: 'Other' },
  ];
  
  const handleCvFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    setCvFile(file);
  };
  
  const onSubmit = async (data: ExpertFormData) => {
    setIsSubmitting(true);
    
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
        form.reset();
        setCvFile(null);
        onSuccess(response.data);
        return { success: true, message: `Expert ${isEditMode ? 'updated' : 'created'} successfully` };
      } else {
        return { 
          success: false, 
          message: response.message || `Failed to ${isEditMode ? 'update' : 'create'} expert` 
        };
      }
    } catch (error) {
      console.error(`Error ${isEditMode ? 'updating' : 'creating'} expert:`, error);
      return { 
        success: false, 
        message: `An error occurred while ${isEditMode ? 'updating' : 'creating'} the expert` 
      };
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <LoadingOverlay isLoading={isSubmitting}>
      <Form
        form={form}
        onSubmit={form.handleSubmitWithNotifications(onSubmit)}
        className="space-y-4"
        resetOnSuccess={false}
        submitText={isEditMode ? 'Update Expert' : 'Create Expert'}
        showResetButton={true}
        resetText="Cancel"
        onReset={onCancel}
      >
        <h2 className="text-xl font-bold mb-4">
          {isEditMode ? 'Edit Expert' : 'Create New Expert'}
        </h2>
        
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          <FormField
            form={form}
            name="name"
            label="Name"
            placeholder="Enter expert's full name"
            required
          />
          
          <FormField
            form={form}
            name="affiliation"
            label="Affiliation"
            placeholder="Enter expert's organization"
            required
          />
          
          <FormField
            form={form}
            name="primaryContact"
            label="Primary Contact"
            placeholder="Email, phone or LinkedIn profile"
            required
          />
          
          <FormField
            form={form}
            name="contactType"
            label="Contact Type"
            type="select"
            options={contactTypeOptions}
            required
          />
          
          <FormField
            form={form}
            name="role"
            label="Role"
            type="select"
            options={roleOptions}
            required
          />
          
          <FormField
            form={form}
            name="employmentType"
            label="Employment Type"
            type="select"
            options={employmentTypeOptions}
            required
          />
          
          <FormField
            form={form}
            name="generalArea"
            label="General Area ID"
            type="number"
            placeholder="Enter area ID"
            required
          />
          
          <FormField
            form={form}
            name="availability"
            label="Availability"
            type="select"
            options={availabilityOptions}
            required
          />
          
          <FormField
            form={form}
            name="isBahraini"
            label="Is Bahraini Citizen"
            type="checkbox"
          />
          
          <div className="md:col-span-2">
            <FormField
              form={form}
              name="skills"
              label="Skills (comma-separated)"
              type="textarea"
              rows={2}
              placeholder="E.g., Machine Learning, Data Science, Python, SQL"
              required
            />
          </div>
          
          <div className="md:col-span-2">
            <FormField
              form={form}
              name="biography"
              label="Biography"
              type="textarea"
              rows={4}
              placeholder="Professional biography..."
            />
          </div>
          
          <div className="md:col-span-2">
            <label className="block text-sm font-medium text-gray-700 mb-1">
              {isEditMode ? 'Update CV (Optional)' : 'CV File *'}
            </label>
            <input
              type="file"
              accept=".pdf,.doc,.docx"
              onChange={handleCvFileChange}
              className="w-full px-3 py-2 bg-white border border-gray-300 rounded focus:outline-none focus:ring-2 focus:ring-blue-500"
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
      </Form>
    </LoadingOverlay>
  );
};

export default ExpertForm;