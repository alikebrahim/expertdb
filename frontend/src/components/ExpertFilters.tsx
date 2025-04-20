import { useState } from 'react';
import { useFormWithNotifications } from '../hooks/useForm';
import { expertFilterSchema } from '../utils/formSchemas';
import { Form, FormField } from './ui';

interface FiltersFormData {
  name?: string;
  role?: string;
  type?: string;
  affiliation?: string;
  expertArea?: string;
  nationality?: string;
  isAvailable?: boolean;
  rating?: string;
  isBahraini?: boolean;
}

interface ExpertFiltersProps {
  onFilterChange: (filters: FiltersFormData) => void;
  initialFilters?: FiltersFormData;
}

const ExpertFilters = ({ onFilterChange, initialFilters = {} }: ExpertFiltersProps) => {
  const [isExpanded, setIsExpanded] = useState(false);
  
  const form = useFormWithNotifications<FiltersFormData>({
    schema: expertFilterSchema,
    defaultValues: {
      name: initialFilters.name || '',
      role: initialFilters.role || '',
      type: initialFilters.type || '',
      affiliation: initialFilters.affiliation || '',
      expertArea: initialFilters.expertArea || '',
      nationality: initialFilters.nationality || '',
      isAvailable: initialFilters.isAvailable || false,
      rating: initialFilters.rating || '',
      isBahraini: initialFilters.isBahraini || false,
    }
  });
  
  const roles = [
    { value: '', label: 'All Roles' },
    { value: 'evaluator', label: 'Evaluator' },
    { value: 'validator', label: 'Validator' },
    { value: 'consultant', label: 'Consultant' },
    { value: 'trainer', label: 'Trainer' },
    { value: 'expert', label: 'Expert' }
  ];
  
  const types = [
    { value: '', label: 'All Types' },
    { value: 'academic', label: 'Academic' },
    { value: 'employer', label: 'Employer' },
    { value: 'freelance', label: 'Freelance' },
    { value: 'government', label: 'Government' },
    { value: 'other', label: 'Other' }
  ];
  
  const ratingOptions = [
    { value: '', label: 'Any Rating' },
    { value: '1', label: '1 Star & Above' },
    { value: '2', label: '2 Stars & Above' },
    { value: '3', label: '3 Stars & Above' },
    { value: '4', label: '4 Stars & Above' },
    { value: '5', label: '5 Stars Only' },
  ];

  const toggleExpanded = () => {
    setIsExpanded(!isExpanded);
  };
  
  const onSubmit = (data: FiltersFormData) => {
    // Convert empty strings to undefined
    const cleanedData = Object.fromEntries(
      Object.entries(data).filter(([_, value]) => {
        if (typeof value === 'string') {
          return value !== '';
        }
        return value !== undefined;
      })
    ) as FiltersFormData;
    
    onFilterChange(cleanedData);
    return { success: true };
  };
  
  const handleReset = () => {
    form.reset({
      name: '',
      role: '',
      type: '',
      affiliation: '',
      expertArea: '',
      nationality: '',
      isAvailable: false,
      rating: '',
      isBahraini: false,
    });
    
    onFilterChange({});
  };
  
  return (
    <div className="bg-white rounded-md shadow p-4 mb-6">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-lg font-semibold text-primary">Filter Experts</h2>
        <button 
          onClick={toggleExpanded}
          className="text-primary hover:text-primary-light flex items-center"
          type="button"
        >
          {isExpanded ? (
            <>
              <span>Hide Filters</span>
              <svg className="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 15l7-7 7 7" />
              </svg>
            </>
          ) : (
            <>
              <span>Show All Filters</span>
              <svg className="w-4 h-4 ml-1" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
              </svg>
            </>
          )}
        </button>
      </div>
      
      <Form
        form={form}
        onSubmit={form.handleSubmitWithNotifications(onSubmit)}
        showResetButton={true}
        resetText="Reset"
        submitText="Apply Filters"
        submitButtonPosition="right"
      >
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {/* Basic filters always visible */}
          <FormField
            form={form}
            name="name"
            label="Expert Name"
            placeholder="Search by name"
          />
          
          <FormField
            form={form}
            name="role"
            label="Role"
            type="select"
            options={roles}
          />
          
          <FormField
            form={form}
            name="type"
            label="Type"
            type="select"
            options={types}
          />
          
          {/* Advanced filters */}
          {isExpanded && (
            <>
              <FormField
                form={form}
                name="affiliation"
                label="Affiliation"
                placeholder="Institution or company"
              />
              
              <FormField
                form={form}
                name="expertArea"
                label="Expert Area"
                placeholder="Area of expertise"
              />
              
              <FormField
                form={form}
                name="nationality"
                label="Nationality"
                placeholder="Expert nationality"
              />
              
              <FormField
                form={form}
                name="rating"
                label="Minimum Rating"
                type="select"
                options={ratingOptions}
              />
              
              <FormField
                form={form}
                name="isAvailable"
                label="Available experts only"
                type="checkbox"
              />
              
              <FormField
                form={form}
                name="isBahraini"
                label="Bahraini citizens only"
                type="checkbox"
              />
            </>
          )}
        </div>
      </Form>
    </div>
  );
};

export default ExpertFilters;