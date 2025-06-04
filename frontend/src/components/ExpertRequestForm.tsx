import { useState, useEffect } from 'react';
import { expertRequestsApi, expertAreasApi } from '../services/api';
import { useFormWithNotifications } from '../hooks/useForm';
import { expertRequestSchema } from '../utils/formSchemas';
import { z } from 'zod';
import { Form } from './ui/Form';
import { FormField } from './ui/FormField';
import { LoadingOverlay } from './ui/LoadingSpinner';

type ExpertRequestFormData = z.infer<typeof expertRequestSchema>;

interface ExpertArea {
  id: number;
  name: string;
}

interface ExpertRequestFormProps {
  onSuccess: () => void;
}

const ExpertRequestForm = ({ onSuccess }: ExpertRequestFormProps) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [cvFile, setCvFile] = useState<File | null>(null);
  const [expertAreas, setExpertAreas] = useState<ExpertArea[]>([]);
  const [loadingAreas, setLoadingAreas] = useState(true);
  
  // Load expert areas for dropdown
  useEffect(() => {
    const fetchExpertAreas = async () => {
      try {
        const response = await expertAreasApi.getExpertAreas();
        if (response.success && response.data) {
          setExpertAreas(response.data);
        }
      } catch (error) {
        console.error('Failed to fetch expert areas:', error);
      } finally {
        setLoadingAreas(false);
      }
    };
    
    fetchExpertAreas();
  }, []);

  // Dropdown options
  const roleOptions = [
    { value: '', label: 'Select role' },
    { value: 'evaluator', label: 'Evaluator' },
    { value: 'validator', label: 'Validator' },
    { value: 'evaluator/validator', label: 'Evaluator/Validator' }
  ];

  const employmentTypeOptions = [
    { value: '', label: 'Select employment type' },
    { value: 'academic', label: 'Academic' },
    { value: 'employer', label: 'Employer' }
  ];

  const generalAreaOptions = [
    { value: 0, label: 'Select general area' },
    ...expertAreas.map(area => ({ value: area.id, label: area.name }))
  ];
  
  const form = useFormWithNotifications<ExpertRequestFormData>({
    schema: expertRequestSchema,
    defaultValues: {
      name: '',
      designation: '',
      institution: '',
      phone: '',
      email: '',
      isBahraini: false,
      isAvailable: false,
      rating: '',
      role: undefined,
      employmentType: undefined,
      isTrained: false,
      isPublished: false,
      generalArea: 0,
      specializedArea: '',
      skills: '',
      biography: '',
      cv: undefined
    },
  });
  
  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    
    if (file) {
      // Validate file type
      if (file.type !== 'application/pdf') {
        form.setError('cv', {
          type: 'manual',
          message: 'Only PDF files are allowed'
        });
        setCvFile(null);
        e.target.value = ''; // Clear the input
        return;
      }
      
      // Validate file size (10MB = 10 * 1024 * 1024 bytes)
      const maxSize = 10 * 1024 * 1024;
      if (file.size > maxSize) {
        form.setError('cv', {
          type: 'manual',
          message: 'File size must be less than 10MB'
        });
        setCvFile(null);
        e.target.value = ''; // Clear the input
        return;
      }
      
      // Clear any previous errors
      form.clearErrors('cv');
      setCvFile(file);
      form.setValue('cv', file);
    } else {
      setCvFile(null);
      form.setValue('cv', undefined);
    }
  };
  
  const handleFormReset = () => {
    setCvFile(null);
    form.setValue('cv', undefined);
  };
  
  const onSubmit = async (data: ExpertRequestFormData): Promise<void> => {
    setIsSubmitting(true);
    
    try {
      // Create FormData for file upload
      const formData = new FormData();
      
      // Add all form fields to FormData
      formData.append('name', data.name);
      formData.append('designation', data.designation);
      formData.append('institution', data.institution);
      formData.append('phone', data.phone);
      formData.append('email', data.email);
      formData.append('isBahraini', data.isBahraini.toString());
      formData.append('isAvailable', data.isAvailable.toString());
      formData.append('rating', data.rating);
      formData.append('role', data.role);
      formData.append('employmentType', data.employmentType);
      formData.append('isTrained', data.isTrained.toString());
      formData.append('isPublished', data.isPublished.toString());
      formData.append('generalArea', data.generalArea.toString());
      formData.append('specializedArea', data.specializedArea);
      formData.append('skills', data.skills);
      formData.append('biography', data.biography);
      
      // Add CV file
      if (cvFile) {
        formData.append('cv', cvFile);
      }
      
      const response = await expertRequestsApi.createExpertRequest(formData);
      
      if (response.success) {
        form.reset();
        setCvFile(null);
        onSuccess();
      } else {
        throw new Error(response.message || 'Failed to submit expert request');
      }
    } catch (error) {
      console.error('Error submitting expert request:', error);
      throw error;
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <LoadingOverlay 
      isLoading={isSubmitting}
      className="w-full"
      label="Submitting expert request..."
    >
      <Form
        form={form}
        onSubmit={onSubmit}
        className="space-y-6"
        showResetButton={true}
        resetText="Reset Form"
        onReset={handleFormReset}
        submitText="Submit Expert Profile"
      >
        {/* Section 1: Personal Information */}
        <div className="bg-gray-50 p-4 rounded-lg">
          <h3 className="text-lg font-semibold text-primary mb-4">Personal Information</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <FormField
              form={form}
              name="name"
              label="Expert Name"
              placeholder="Enter expert's full name"
              required
            />
            
            <FormField
              form={form}
              name="designation"
              label="Designation"
              placeholder="Enter professional title"
              required
            />
            
            <FormField
              form={form}
              name="institution"
              label="Institution"
              placeholder="Enter affiliated organization"
              required
            />
            
            <FormField
              form={form}
              name="phone"
              label="Phone Number"
              placeholder="Enter contact phone"
              required
            />
            
            <div className="md:col-span-2">
              <FormField
                form={form}
                name="email"
                label="Email Address"
                type="email"
                placeholder="Enter contact email"
                required
              />
            </div>
          </div>
        </div>

        {/* Section 2: Professional Details */}
        <div className="bg-gray-50 p-4 rounded-lg">
          <h3 className="text-lg font-semibold text-primary mb-4">Professional Details</h3>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div className="flex items-center space-x-4">
              <FormField
                form={form}
                name="isBahraini"
                label="Bahraini National"
                type="checkbox"
              />
              
              <FormField
                form={form}
                name="isAvailable"
                label="Currently Available"
                type="checkbox"
              />
            </div>
            
            <FormField
              form={form}
              name="rating"
              label="Performance Rating"
              placeholder="Enter rating/score"
              required
            />
            
            <FormField
              form={form}
              name="role"
              label="Expert Role"
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
            
            <div className="flex items-center space-x-4">
              <FormField
                form={form}
                name="isTrained"
                label="Training Completed"
                type="checkbox"
              />
              
              <FormField
                form={form}
                name="isPublished"
                label="Allow Publishing"
                type="checkbox"
                hint="Allow profile to be published publicly"
              />
            </div>
          </div>
        </div>

        {/* Section 3: Expertise Areas */}
        <div className="bg-gray-50 p-4 rounded-lg">
          <h3 className="text-lg font-semibold text-primary mb-4">Expertise Areas</h3>
          <div className="grid grid-cols-1 gap-4">
            <FormField
              form={form}
              name="generalArea"
              label="General Area"
              type="select"
              options={generalAreaOptions}
              required
              disabled={loadingAreas}
              hint={loadingAreas ? "Loading areas..." : "Select the primary expertise area"}
            />
            
            <FormField
              form={form}
              name="specializedArea"
              label="Specialized Area"
              placeholder="Enter specific field of specialization"
              required
            />
            
            <FormField
              form={form}
              name="skills"
              label="Skills & Competencies"
              type="textarea"
              placeholder="Enter skills separated by commas (e.g., Project Management, Quality Assurance, Data Analysis)"
              rows={3}
              required
              hint="List the expert's key skills and competencies"
            />
          </div>
        </div>

        {/* Section 4: Biography & Documents */}
        <div className="bg-gray-50 p-4 rounded-lg">
          <h3 className="text-lg font-semibold text-primary mb-4">Biography & Documents</h3>
          <div className="space-y-4">
            <FormField
              form={form}
              name="biography"
              label="Professional Biography"
              type="textarea"
              placeholder="Enter a comprehensive professional summary..."
              rows={6}
              required
              hint="Maximum 1000 characters - include education, experience, and achievements"
            />
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                CV Document <span className="text-red-500">*</span>
              </label>
              <input
                type="file"
                accept=".pdf"
                onChange={handleFileChange}
                className={`w-full px-3 py-2 bg-white border rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary ${
                  form.formState.errors.cv ? 'border-red-500' : 'border-gray-300'
                }`}
                required
              />
              <p className="mt-1 text-sm text-gray-500">
                Upload expert's CV in PDF format (max 10MB)
              </p>
              {form.formState.errors.cv && (
                <p className="mt-1 text-sm text-red-600">
                  {form.formState.errors.cv.message}
                </p>
              )}
              {cvFile && !form.formState.errors.cv && (
                <p className="mt-1 text-sm text-green-600">
                  File selected: {cvFile.name} ({(cvFile.size / 1024 / 1024).toFixed(2)} MB)
                </p>
              )}
            </div>
          </div>
        </div>
      </Form>
    </LoadingOverlay>
  );
};

export default ExpertRequestForm;