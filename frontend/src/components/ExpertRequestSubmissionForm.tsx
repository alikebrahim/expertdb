import React, { useState, useEffect } from 'react';
import { useFormWithNotifications } from '../hooks/useForm';
import { expertRequestSchema } from '../utils/formSchemas';
import { getExpertAreas } from '../api/areas';
import { expertsApi } from '../services/api';
import { Form } from './ui/Form';
import { FormField } from './ui/FormField';
import { ProgressStepper } from './ui/ProgressStepper';
import { Card, CardHeader, CardContent } from './ui/Card';
import { Alert } from './ui/Alert';
import { LoadingOverlay } from './ui/LoadingSpinner';
import { z } from 'zod';

type ExpertRequestFormData = z.infer<typeof expertRequestSchema>;

interface ExpertRequestSubmissionFormProps {
  onSuccess: () => void;
  initialData?: Partial<ExpertRequestFormData>;
}

interface ExpertArea {
  id: number;
  name: string;
}

const ExpertRequestSubmissionForm: React.FC<ExpertRequestSubmissionFormProps> = ({
  onSuccess,
  initialData
}) => {
  const [currentStep, setCurrentStep] = useState(1);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [expertAreas, setExpertAreas] = useState<ExpertArea[]>([]);
  const [loadingAreas, setLoadingAreas] = useState(true);
  const [cvFile, setCvFile] = useState<File | null>(null);

  const steps = [
    { id: 1, label: 'Personal Information', description: 'Basic contact details' },
    { id: 2, label: 'Professional Details', description: 'Employment and availability' },
    { id: 3, label: 'Expertise Areas', description: 'Skills and specializations' },
    { id: 4, label: 'Biography & Documents', description: 'Profile and CV upload' }
  ];

  const form = useFormWithNotifications<ExpertRequestFormData>({
    schema: expertRequestSchema,
    defaultValues: {
      // Personal Information
      name: initialData?.name || '',
      designation: initialData?.designation || '',
      institution: initialData?.institution || '',
      phone: initialData?.phone || '',
      email: initialData?.email || '',
      
      // Professional Details
      isBahraini: initialData?.isBahraini || false,
      isAvailable: initialData?.isAvailable || true,
      rating: initialData?.rating?.toString() || '0',
      role: initialData?.role || 'evaluator',
      employmentType: initialData?.employmentType || 'academic',
      isTrained: initialData?.isTrained || false,
      isPublished: initialData?.isPublished || false,
      
      // Expertise Areas
      generalArea: typeof initialData?.generalArea === 'number' ? initialData.generalArea : parseInt(initialData?.generalArea || '1') || 1,
      specializedArea: initialData?.specializedArea || '',
      skills: initialData?.skills || '',
      
      // Biography & Documents
      biography: initialData?.biography || ''
    },
  });

  // Load expert areas on component mount
  useEffect(() => {
    const loadExpertAreas = async () => {
      try {
        setLoadingAreas(true);
        const response = await getExpertAreas();
        if (response.success && response.data) {
          setExpertAreas(response.data);
        }
      } catch (error) {
        console.error('Failed to load expert areas:', error);
      } finally {
        setLoadingAreas(false);
      }
    };

    loadExpertAreas();
  }, []);

  const roleOptions = [
    { value: 'consultant', label: 'Consultant' },
    { value: 'reviewer', label: 'Reviewer' },
    { value: 'auditor', label: 'Auditor' },
    { value: 'assessor', label: 'Assessor' },
    { value: 'trainer', label: 'Trainer' }
  ];

  const employmentTypeOptions = [
    { value: 'full_time', label: 'Full Time' },
    { value: 'part_time', label: 'Part Time' },
    { value: 'contract', label: 'Contract' },
    { value: 'freelance', label: 'Freelance' }
  ];

  const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0] || null;
    setCvFile(file);
  };

  const validateCurrentStep = async (): Promise<boolean> => {
    const values = form.getValues();
    
    try {
      switch (currentStep) {
        case 1:
          // Validate personal information fields
          await expertRequestSchema.pick({
            name: true,
            designation: true,
            institution: true,
            phone: true,
            email: true
          }).parseAsync(values);
          break;
        case 2:
          // Validate professional details fields
          await expertRequestSchema.pick({
            isBahraini: true,
            isAvailable: true,
            role: true,
            employmentType: true,
            isTrained: true,
            isPublished: true
          }).parseAsync(values);
          break;
        case 3:
          // Validate expertise areas fields
          await expertRequestSchema.pick({
            generalArea: true,
            specializedArea: true,
            skills: true
          }).parseAsync(values);
          break;
        case 4:
          // Validate biography and documents
          await expertRequestSchema.pick({
            biography: true
          }).parseAsync(values);
          break;
      }
      return true;
    } catch (error) {
      // Trigger form validation to show errors
      form.trigger();
      return false;
    }
  };

  const handleNext = async () => {
    const isValid = await validateCurrentStep();
    if (isValid && currentStep < 4) {
      setCurrentStep(currentStep + 1);
    }
  };

  const handlePrevious = () => {
    if (currentStep > 1) {
      setCurrentStep(currentStep - 1);
    }
  };

  const onSubmit = async (data: ExpertRequestFormData): Promise<void> => {
    setIsSubmitting(true);
    
    try {
      // Create FormData for file upload
      const formData = new FormData();
      
      // Add all form fields
      Object.entries(data).forEach(([key, value]) => {
        if (value !== null && value !== undefined) {
          formData.append(key, value.toString());
        }
      });
      
      // Add CV file if available
      if (cvFile) {
        formData.append('cv', cvFile);
      }
      
      const response = await expertsApi.createExpert(formData);
      
      if (response.success) {
        onSuccess();
      } else {
        throw new Error(response.message || 'Failed to submit expert application');
      }
    } catch (error) {
      console.error('Error submitting expert application:', error);
      throw error;
    } finally {
      setIsSubmitting(false);
    }
  };

  if (loadingAreas) {
    return (
      <LoadingOverlay isLoading={true} label="Loading form..." className="min-h-96">
        <div className="min-h-96 flex items-center justify-center">
          <p>Loading...</p>
        </div>
      </LoadingOverlay>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-8">
        <ProgressStepper steps={steps} currentStep={currentStep} />
      </div>

      <LoadingOverlay isLoading={isSubmitting} label="Submitting application...">
        <Form form={form} onSubmit={onSubmit} className="space-y-6">
          {/* Step 1: Personal Information */}
          {currentStep === 1 && (
            <Card>
              <CardHeader>
                <h3 className="text-lg font-semibold text-primary">Personal Information</h3>
                <p className="text-sm text-gray-600">Please provide your basic contact details</p>
              </CardHeader>
              <CardContent>
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div className="md:col-span-2">
                    <FormField
                      form={form}
                      name="name"
                      label="Full Name"
                      placeholder="Enter your full name"
                      required
                    />
                  </div>
                  
                  <FormField
                    form={form}
                    name="designation"
                    label="Designation/Title"
                    placeholder="e.g., Senior Quality Manager"
                    required
                  />
                  
                  <FormField
                    form={form}
                    name="institution"
                    label="Institution/Company"
                    placeholder="Enter your organization"
                    required
                  />
                  
                  <FormField
                    form={form}
                    name="phone"
                    label="Phone Number"
                    placeholder="e.g., +973 XXXX XXXX"
                    required
                  />
                  
                  <FormField
                    form={form}
                    name="email"
                    label="Email Address"
                    type="email"
                    placeholder="your.email@example.com"
                    required
                  />
                </div>
              </CardContent>
            </Card>
          )}

          {/* Step 2: Professional Details */}
          {currentStep === 2 && (
            <Card>
              <CardHeader>
                <h3 className="text-lg font-semibold text-primary">Professional Details</h3>
                <p className="text-sm text-gray-600">Tell us about your professional background</p>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <FormField
                      form={form}
                      name="role"
                      label="Preferred Role"
                      type="select"
                      options={[
                        { value: '', label: 'Select a role' },
                        ...roleOptions
                      ]}
                      required
                    />
                    
                    <FormField
                      form={form}
                      name="employmentType"
                      label="Employment Type"
                      type="select"
                      options={[
                        { value: '', label: 'Select employment type' },
                        ...employmentTypeOptions
                      ]}
                      required
                    />
                  </div>
                  
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <label className="flex items-center space-x-2">
                        <FormField
                          form={form}
                          name="isBahraini"
                          type="checkbox"
                          label=""
                        />
                        <span className="text-sm font-medium">Bahraini National</span>
                      </label>
                    </div>
                    
                    <div className="space-y-2">
                      <label className="flex items-center space-x-2">
                        <FormField
                          form={form}
                          name="isAvailable"
                          type="checkbox"
                          label=""
                        />
                        <span className="text-sm font-medium">Currently Available</span>
                      </label>
                    </div>
                  </div>
                  
                  <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                    <div className="space-y-2">
                      <label className="flex items-center space-x-2">
                        <FormField
                          form={form}
                          name="isTrained"
                          type="checkbox"
                          label=""
                        />
                        <span className="text-sm font-medium">Received BQA Training</span>
                      </label>
                    </div>
                    
                    <div className="space-y-2">
                      <label className="flex items-center space-x-2">
                        <FormField
                          form={form}
                          name="isPublished"
                          type="checkbox"
                          label=""
                        />
                        <span className="text-sm font-medium">Has Published Work</span>
                      </label>
                    </div>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Step 3: Expertise Areas */}
          {currentStep === 3 && (
            <Card>
              <CardHeader>
                <h3 className="text-lg font-semibold text-primary">Expertise Areas</h3>
                <p className="text-sm text-gray-600">Define your areas of expertise and skills</p>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <FormField
                    form={form}
                    name="generalArea"
                    label="General Area of Expertise"
                    type="select"
                    options={[
                      { value: '', label: 'Select general area' },
                      ...expertAreas.map(area => ({
                        value: area.name,
                        label: area.name
                      }))
                    ]}
                    required
                  />
                  
                  <FormField
                    form={form}
                    name="specializedArea"
                    label="Specialized Area"
                    placeholder="Enter your specific area of specialization"
                    required
                  />
                  
                  <FormField
                    form={form}
                    name="skills"
                    label="Key Skills"
                    type="textarea"
                    placeholder="List your key skills and competencies (comma-separated)"
                    rows={4}
                    required
                    hint="e.g., Quality Assurance, Process Improvement, Risk Management"
                  />
                </div>
              </CardContent>
            </Card>
          )}

          {/* Step 4: Biography & Documents */}
          {currentStep === 4 && (
            <Card>
              <CardHeader>
                <h3 className="text-lg font-semibold text-primary">Biography & Documents</h3>
                <p className="text-sm text-gray-600">Provide your professional biography and upload your CV</p>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <FormField
                    form={form}
                    name="biography"
                    label="Professional Biography"
                    type="textarea"
                    rows={8}
                    placeholder="Write a brief professional biography highlighting your experience, achievements, and expertise..."
                    required
                    hint={`${form.watch('biography')?.length || 0}/1000 characters`}
                  />
                  
                  <div>
                    <label className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">
                      Curriculum Vitae (CV) *
                    </label>
                    <input
                      type="file"
                      accept=".pdf"
                      onChange={handleFileChange}
                      className="w-full px-3 py-2 bg-white border border-gray-300 rounded-md focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary dark:bg-gray-800 dark:border-gray-600 dark:text-white"
                      required
                    />
                    <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">
                      Upload your CV in PDF format (max 5MB)
                    </p>
                    {cvFile && (
                      <p className="mt-1 text-sm text-green-600 dark:text-green-400">
                        File selected: {cvFile.name}
                      </p>
                    )}
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Navigation Buttons */}
          <div className="flex justify-between pt-6">
            <button
              type="button"
              onClick={handlePrevious}
              disabled={currentStep === 1}
              className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Previous
            </button>
            
            {currentStep < 4 ? (
              <button
                type="button"
                onClick={handleNext}
                className="px-4 py-2 text-sm font-medium text-white bg-primary border border-transparent rounded-md hover:bg-primary-dark focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary"
              >
                Next
              </button>
            ) : (
              <button
                type="submit"
                disabled={isSubmitting}
                className="px-6 py-2 text-sm font-medium text-white bg-primary border border-transparent rounded-md hover:bg-primary-dark focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-primary disabled:opacity-50 disabled:cursor-not-allowed"
              >
                {isSubmitting ? 'Submitting...' : 'Submit Application'}
              </button>
            )}
          </div>
        </Form>
      </LoadingOverlay>
    </div>
  );
};

export default ExpertRequestSubmissionForm;
