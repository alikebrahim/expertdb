import { useState, useEffect, useCallback } from 'react';
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
import { useUI } from '../hooks/useUI';
import FileUpload from './ui/FileUpload';

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
  const [lastSaved, setLastSaved] = useState<Date | null>(null);
  const [isAutoSaving, setIsAutoSaving] = useState(false);
  const [currentStep, setCurrentStep] = useState(0);
  const { showNotification } = useUI();
  
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
  
  const handleFileChange = (file: File | null) => {
    setCvFile(file);
    form.setValue('cv', file || undefined);
    
    if (file) {
      form.clearErrors('cv');
    }
  };

  const handleSkillsChange = (skills: string[]) => {
    setSkillTags(skills);
    form.setValue('skills', skills);
  };

  // Auto-save functionality
  const saveDraft = useCallback(async () => {
    if (isAutoSaving) return;
    
    setIsAutoSaving(true);
    try {
      const formData = form.getValues();
      const draftData = {
        ...formData,
        skills: skillTags,
        isDraft: true
      };
      
      // Save to localStorage as backup
      localStorage.setItem('expertRequestDraft', JSON.stringify(draftData));
      
      // TODO: Implement API call to save draft
      // await expertRequestsApi.saveDraft(draftData);
      
      setLastSaved(new Date());
      showNotification('Draft saved automatically', 'success', 2000);
    } catch (error) {
      console.error('Auto-save failed:', error);
    } finally {
      setIsAutoSaving(false);
    }
  }, [form, skillTags, isAutoSaving, showNotification]);

  // Auto-save every 10 seconds
  useEffect(() => {
    const interval = setInterval(() => {
      const formData = form.getValues();
      if (formData.name || formData.email || formData.phone || skillTags.length > 0) {
        saveDraft();
      }
    }, 10000);

    return () => clearInterval(interval);
  }, [saveDraft]);

  // Load draft on component mount
  useEffect(() => {
    const savedDraft = localStorage.getItem('expertRequestDraft');
    if (savedDraft) {
      try {
        const draftData = JSON.parse(savedDraft);
        form.reset(draftData);
        setSkillTags(draftData.skills || []);
        setLastSaved(new Date(draftData.lastSaved || Date.now()));
      } catch (error) {
        console.error('Failed to load draft:', error);
      }
    }
  }, [form]);

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
        // Clear draft from localStorage
        localStorage.removeItem('expertRequestDraft');
        setLastSaved(null);
        
        showNotification('Expert request submitted successfully!', 'success');
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

  const sections = [
    { id: 'personal', title: 'Personal Information', icon: '👤' },
    { id: 'professional', title: 'Professional Details', icon: '💼' },
    { id: 'expertise', title: 'Expertise Areas', icon: '🎯' },
    { id: 'biography', title: 'Biography & Documents', icon: '📄' }
  ];

  return (
    <div className="max-w-4xl mx-auto p-6 space-y-6">
      {/* Progress Stepper */}
      <div className="bg-white rounded-lg shadow-sm border p-6">
        <h1 className="text-2xl font-semibold text-gray-900 mb-6">Submit Expert Request</h1>
        
        <div className="flex items-center justify-between mb-6">
          {sections.map((section, index) => (
            <div key={section.id} className="flex items-center">
              <div className={`flex items-center justify-center w-10 h-10 rounded-full border-2 text-sm font-semibold ${
                index <= currentStep 
                  ? 'bg-blue-600 text-white border-blue-600' 
                  : 'bg-gray-100 text-gray-400 border-gray-300'
              }`}>
                {index < currentStep ? '✓' : index + 1}
              </div>
              <div className="ml-3 hidden sm:block">
                <div className={`text-sm font-medium ${
                  index <= currentStep ? 'text-blue-600' : 'text-gray-400'
                }`}>
                  {section.title}
                </div>
              </div>
              {index < sections.length - 1 && (
                <div className={`w-16 h-0.5 ml-4 ${
                  index < currentStep ? 'bg-blue-600' : 'bg-gray-300'
                }`} />
              )}
            </div>
          ))}
        </div>

        {/* Auto-save status */}
        <div className="flex items-center justify-between text-sm text-gray-500 mb-4">
          <div className="flex items-center space-x-2">
            {isAutoSaving ? (
              <>
                <div className="w-4 h-4 border-2 border-blue-600 border-t-transparent rounded-full animate-spin" />
                <span>Saving draft...</span>
              </>
            ) : lastSaved ? (
              <>
                <span className="text-green-600">✓</span>
                <span>Auto-saved {new Date(lastSaved).toLocaleTimeString()}</span>
              </>
            ) : (
              <span>Changes will be auto-saved</span>
            )}
          </div>
          <button
            type="button"
            onClick={saveDraft}
            disabled={isAutoSaving}
            className="text-blue-600 hover:text-blue-800 font-medium"
          >
            Save Draft Now
          </button>
        </div>
      </div>

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
            <p className="text-sm text-gray-600">
              Upload your curriculum vitae (CV) in PDF format
            </p>
          </CardHeader>
          <CardContent>
            <FileUpload
              onFileSelect={handleFileChange}
              accept=".pdf"
              maxSize={20}
              currentFile={cvFile}
              error={form.formState.errors.cv?.message}
              label="CV Document (PDF)"
              required
            />
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