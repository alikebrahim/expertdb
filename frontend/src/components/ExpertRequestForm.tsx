import { useState, useEffect } from 'react';
import { expertRequestsApi, expertAreasApi } from '../services/api';
import { useFormWithNotifications } from '../hooks/useForm';
import { expertRequestSchema } from '../utils/formSchemas';
import { z } from 'zod';
import { Form } from './ui/Form';
import { FormField } from './ui/FormField';
import { LoadingOverlay } from './ui/LoadingSpinner';
import { Card, CardHeader, CardContent } from './ui/Card';
import Button from './ui/Button';
import BiographyForm from './BiographyForm';
import { TagInput } from './ui/TagInput';

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
  const [skillTags, setSkillTags] = useState<string[]>([]);
  
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
  const designationOptions = [
    { value: '', label: 'Select designation' },
    { value: 'Prof.', label: 'Prof.' },
    { value: 'Dr.', label: 'Dr.' },
    { value: 'Mr.', label: 'Mr.' },
    { value: 'Ms.', label: 'Ms.' },
    { value: 'Mrs.', label: 'Mrs.' },
    { value: 'Miss', label: 'Miss' },
    { value: 'Eng.', label: 'Eng.' }
  ];

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

  const ratingOptions = [
    { value: 0, label: 'Select rating' },
    { value: 1, label: '1 - Basic' },
    { value: 2, label: '2 - Fair' },
    { value: 3, label: '3 - Good' },
    { value: 4, label: '4 - Very Good' },
    { value: 5, label: '5 - Excellent' }
  ];

  const generalAreaOptions = [
    { value: 0, label: 'Select general area' },
    ...expertAreas.map(area => ({ value: area.id, label: area.name }))
  ];
  
  const form = useFormWithNotifications<ExpertRequestFormData>({
    schema: expertRequestSchema,
    defaultValues: {
      name: '',
      designation: undefined,
      affiliation: '',
      phone: '',
      email: '',
      isBahraini: false,
      isAvailable: false,
      rating: 0,
      role: undefined,
      employmentType: undefined,
      isTrained: false,
      isPublished: false,
      generalArea: 0,
      specializedArea: '',
      skills: [],
      biography: {
        experience: [],
        education: []
      },
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

      // Validate file size (5MB)
      if (file.size > 5 * 1024 * 1024) {
        form.setError('cv', {
          type: 'manual',
          message: 'File size must be less than 5MB'
        });
        setCvFile(null);
        e.target.value = ''; // Clear the input
        return;
      }

      // File is valid
      setCvFile(file);
      form.setValue('cv', file);
      form.clearErrors('cv');
    } else {
      setCvFile(null);
      form.setValue('cv', undefined);
    }
  };

  const handleSkillsChange = (skills: string[]) => {
    setSkillTags(skills);
    form.setValue('skills', skills);
  };

  const onSubmit = async (data: ExpertRequestFormData) => {
    try {
      setIsSubmitting(true);
      
      // Create FormData for file upload
      const formData = new FormData();
      
      // Add all form fields
      formData.append('name', data.name);
      formData.append('designation', data.designation);
      formData.append('affiliation', data.affiliation);
      formData.append('phone', data.phone);
      formData.append('email', data.email);
      formData.append('isBahraini', data.isBahraini.toString());
      formData.append('isAvailable', data.isAvailable.toString());
      formData.append('rating', data.rating.toString());
      formData.append('role', data.role);
      formData.append('employmentType', data.employmentType);
      formData.append('generalArea', data.generalArea.toString());
      formData.append('specializedArea', data.specializedArea);
      formData.append('isTrained', data.isTrained.toString());
      formData.append('isPublished', data.isPublished.toString());
      
      // Handle skills array
      formData.append('skills', JSON.stringify(data.skills));
      
      // Handle biography object
      formData.append('biography', JSON.stringify(data.biography));
      
      // Add CV file
      if (cvFile) {
        formData.append('cv', cvFile);
      }

      const response = await expertRequestsApi.createExpertRequest(formData);
      
      if (response.success) {
        onSuccess();
        form.reset();
        setCvFile(null);
        setSkillTags([]);
      }
    } catch (error) {
      console.error('Error submitting expert request:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  if (loadingAreas) {
    return <LoadingOverlay isLoading={true}>Loading areas...</LoadingOverlay>;
  }

  return (
    <div className="max-w-4xl mx-auto p-6 space-y-6">
      <Form
        form={form}
        onSubmit={onSubmit}
        className="space-y-6"
      >
        {/* Personal Information Section */}
        <Card>
          <CardHeader>
            <h2 className="text-xl font-semibold">Personal Information</h2>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <FormField
                name="name"
                label="Full Name"
                control={form.control}
                error={form.formState.errors.name}
                required
                placeholder="Enter full name"
              />
              <FormField
                name="designation"
                label="Designation"
                type="select"
                options={designationOptions}
                control={form.control}
                error={form.formState.errors.designation}
                required
              />
            </div>
            
            <FormField
              name="affiliation"
              label="Affiliation"
              control={form.control}
              error={form.formState.errors.affiliation}
              required
              placeholder="Enter organization or institution"
            />
            
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <FormField
                name="phone"
                label="Phone Number"
                control={form.control}
                error={form.formState.errors.phone}
                required
                placeholder="+973 XXXX XXXX"
              />
              <FormField
                name="email"
                label="Email Address"
                type="email"
                control={form.control}
                error={form.formState.errors.email}
                required
                placeholder="email@example.com"
              />
            </div>
          </CardContent>
        </Card>

        {/* Professional Details Section */}
        <Card>
          <CardHeader>
            <h2 className="text-xl font-semibold">Professional Details</h2>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div className="flex items-center space-x-2">
                <FormField
                  name="isBahraini"
                  label="Bahraini Citizen"
                  type="checkbox"
                  control={form.control}
                />
              </div>
              <div className="flex items-center space-x-2">
                <FormField
                  name="isAvailable"
                  label="Currently Available"
                  type="checkbox"
                  control={form.control}
                />
              </div>
              <div className="flex items-center space-x-2">
                <FormField
                  name="isTrained"
                  label="BQA Trained"
                  type="checkbox"
                  control={form.control}
                />
              </div>
            </div>
            
            <div className="flex items-center space-x-2">
              <FormField
                name="isPublished"
                label="Has Published Work"
                type="checkbox"
                control={form.control}
              />
            </div>
            
            <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
              <FormField
                name="rating"
                label="Performance Rating"
                type="select"
                options={ratingOptions}
                control={form.control}
                error={form.formState.errors.rating}
                required
              />
              <FormField
                name="role"
                label="Expert Role"
                type="select"
                options={roleOptions}
                control={form.control}
                error={form.formState.errors.role}
                required
              />
              <FormField
                name="employmentType"
                label="Employment Type"
                type="select"
                options={employmentTypeOptions}
                control={form.control}
                error={form.formState.errors.employmentType}
                required
              />
            </div>
          </CardContent>
        </Card>

        {/* Expertise Areas Section */}
        <Card>
          <CardHeader>
            <h2 className="text-xl font-semibold">Expertise Areas</h2>
          </CardHeader>
          <CardContent className="space-y-4">
            <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
              <FormField
                name="generalArea"
                label="General Area"
                type="select"
                options={generalAreaOptions}
                control={form.control}
                error={form.formState.errors.generalArea}
                required
              />
              <FormField
                name="specializedArea"
                label="Specialized Area"
                control={form.control}
                error={form.formState.errors.specializedArea}
                required
                placeholder="Enter specific field of specialization"
              />
            </div>
            
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                Skills & Competencies *
              </label>
              <TagInput
                value={skillTags}
                onChange={handleSkillsChange}
                placeholder="Type a skill and press Enter"
                className="w-full"
              />
              {form.formState.errors.skills && (
                <p className="mt-1 text-sm text-red-600">
                  {form.formState.errors.skills.message}
                </p>
              )}
            </div>
          </CardContent>
        </Card>

        {/* Biography Section */}
        <Card>
          <CardHeader>
            <h2 className="text-xl font-semibold">Biography</h2>
            <p className="text-sm text-gray-600">
              Add your professional experience and educational background
            </p>
          </CardHeader>
          <CardContent>
            <BiographyForm
              control={form.control}
              setValue={form.setValue}
              watch={form.watch}
              errors={form.formState.errors}
            />
          </CardContent>
        </Card>

        {/* CV Upload Section */}
        <Card>
          <CardHeader>
            <h2 className="text-xl font-semibold">CV Upload</h2>
          </CardHeader>
          <CardContent>
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-1">
                CV Document (PDF) *
              </label>
              <input
                type="file"
                accept=".pdf"
                onChange={handleFileChange}
                className="mt-1 block w-full text-sm text-gray-500
                          file:mr-4 file:py-2 file:px-4
                          file:rounded-full file:border-0
                          file:text-sm file:font-semibold
                          file:bg-blue-50 file:text-blue-700
                          hover:file:bg-blue-100"
              />
              {cvFile && (
                <p className="mt-2 text-sm text-green-600">
                  Selected: {cvFile.name} ({Math.round(cvFile.size / 1024)} KB)
                </p>
              )}
              {form.formState.errors.cv && (
                <p className="mt-1 text-sm text-red-600">
                  {form.formState.errors.cv.message}
                </p>
              )}
              <p className="mt-1 text-xs text-gray-500">
                Maximum file size: 5MB. Only PDF files are allowed.
              </p>
            </div>
          </CardContent>
        </Card>

        {/* Submit Button */}
        <div className="flex justify-end space-x-4">
          <Button
            type="button"
            variant="outline"
            onClick={() => {
              form.reset();
              setCvFile(null);
              setSkillTags([]);
            }}
            disabled={isSubmitting}
          >
            Reset Form
          </Button>
          <Button
            type="submit"
            disabled={isSubmitting}
            className="px-8"
          >
            {isSubmitting ? 'Submitting...' : 'Submit Expert Request'}
          </Button>
        </div>
      </Form>
      
      {isSubmitting && <LoadingOverlay isLoading={true}>Submitting...</LoadingOverlay>}
    </div>
  );
};

export default ExpertRequestForm;