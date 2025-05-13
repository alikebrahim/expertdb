import { useState, useEffect } from 'react';
import { useForm } from 'react-hook-form';
import { useUI } from '../../hooks/useUI';
import * as areasApi from '../../api/areas';
import Input from '../ui/Input';
import Button from '../ui/Button';
import { LoadingSpinner } from '../ui';

interface FiltersFormData {
  name?: string;
  role?: string;
  type?: string;
  affiliation?: string;
  expertAreaId?: string;
  nationality?: string;
  rating?: string;
  isAvailable?: boolean;
  isBahraini?: boolean;
  isPublished?: boolean;
}

interface ExpertFiltersProps {
  onFilterChange: (filters: FiltersFormData) => void;
  initialFilters?: Partial<FiltersFormData>;
}

const ExpertFilters = ({ onFilterChange, initialFilters = {} }: ExpertFiltersProps) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const [expertAreas, setExpertAreas] = useState<Array<{ id: number; name: string }>>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [activeFilterCount, setActiveFilterCount] = useState(0);
  const { addNotification } = useUI();
  const { register, handleSubmit, reset, setValue, watch } = useForm<FiltersFormData>({
    defaultValues: initialFilters
  });
  
  // Watch form values to update active filter count
  const formValues = watch();
  
  // Update active filter count when form values change
  useEffect(() => {
    const count = Object.entries(formValues).filter(([_, value]) => {
      if (typeof value === 'string') {
        return value !== '';
      }
      return value !== undefined;
    }).length;
    
    setActiveFilterCount(count);
  }, [formValues]);
  
  // Fetch expert areas on mount
  useEffect(() => {
    const fetchExpertAreas = async () => {
      setIsLoading(true);
      try {
        const response = await areasApi.getExpertAreas();
        if (response.success && response.data) {
          setExpertAreas(response.data);
        } else {
          addNotification({
            type: 'error',
            message: 'Failed to load expert areas',
            duration: 5000,
          });
        }
      } catch (error) {
        console.error('Error fetching expert areas:', error);
        addNotification({
          type: 'error',
          message: 'Error loading expert areas',
          duration: 5000,
        });
      } finally {
        setIsLoading(false);
      }
    };

    fetchExpertAreas();
  }, [addNotification]);

  // Set initial filters
  useEffect(() => {
    if (initialFilters && Object.keys(initialFilters).length > 0) {
      Object.entries(initialFilters).forEach(([key, value]) => {
        setValue(key as keyof FiltersFormData, value);
      });
    }
  }, [initialFilters, setValue]);
  
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
  
  const ratings = [
    { value: '', label: 'Any Rating' },
    { value: '5', label: '5 Stars' },
    { value: '4', label: '4+ Stars' },
    { value: '3', label: '3+ Stars' },
    { value: '2', label: '2+ Stars' }
  ];
  
  const nationalities = [
    { value: '', label: 'All Nationalities' },
    { value: 'bahraini', label: 'Bahraini' },
    { value: 'international', label: 'International' }
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
    
    // Show notification of applied filters
    if (Object.keys(cleanedData).length > 0) {
      addNotification({
        type: 'success',
        message: `${Object.keys(cleanedData).length} filter(s) applied`,
        duration: 2000,
      });
    }
  };
  
  const handleReset = () => {
    reset({
      name: '',
      role: '',
      type: '',
      affiliation: '',
      expertAreaId: '',
      nationality: '',
      rating: '',
      isAvailable: undefined,
      isBahraini: undefined
    });
    
    onFilterChange({});
    
    addNotification({
      type: 'info',
      message: 'Filters reset',
      duration: 2000,
    });
  };
  
  return (
    <div className="bg-white rounded-lg shadow-md p-4 mb-6">
      <div className="flex justify-between items-center mb-4">
        <div className="flex items-center">
          <h2 className="text-lg font-semibold text-primary">Filter Experts</h2>
          {activeFilterCount > 0 && (
            <span className="ml-2 inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-primary text-white">
              {activeFilterCount} active
            </span>
          )}
        </div>
        <button 
          onClick={toggleExpanded}
          className="text-primary hover:text-primary-dark transition-colors px-3 py-1 rounded-md border border-primary hover:bg-primary-50 text-sm font-medium"
        >
          {isExpanded ? 'Hide Advanced Filters' : 'Show Advanced Filters'}
        </button>
      </div>
      
      <form onSubmit={handleSubmit(onSubmit)}>
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {/* Basic filters always visible */}
          <Input
            label="Expert Name"
            placeholder="Search by name"
            {...register('name')}
          />
          
          <div className="mb-4">
            <label className="block text-sm font-medium text-neutral-700 mb-1">
              Role
            </label>
            <select
              className="w-full px-3 py-2 bg-white border border-neutral-300 rounded-md focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
              {...register('role')}
            >
              {roles.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>
          
          <div className="mb-4">
            <label className="block text-sm font-medium text-neutral-700 mb-1">
              Type
            </label>
            <select
              className="w-full px-3 py-2 bg-white border border-neutral-300 rounded-md focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
              {...register('type')}
            >
              {types.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>
          
          {/* Quick filter checkboxes */}
          <div className="mb-4 flex items-center">
            <input
              type="checkbox"
              id="isAvailable"
              className="h-4 w-4 text-primary focus:ring-primary border-neutral-300 rounded"
              {...register('isAvailable')}
            />
            <label htmlFor="isAvailable" className="ml-2 block text-sm text-neutral-700">
              Available experts only
            </label>
          </div>
          
          {/* Advanced filters */}
          {isExpanded && (
            <>
              <Input
                label="Affiliation"
                placeholder="Institution or company"
                {...register('affiliation')}
              />
              
              <div className="mb-4">
                <label className="block text-sm font-medium text-neutral-700 mb-1">
                  Expert Area
                </label>
                {isLoading ? (
                  <div className="flex items-center justify-center py-1">
                    <LoadingSpinner size="sm" />
                    <span className="ml-2 text-sm text-neutral-500">Loading areas...</span>
                  </div>
                ) : (
                  <select
                    className="w-full px-3 py-2 bg-white border border-neutral-300 rounded-md focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
                    {...register('expertAreaId')}
                  >
                    <option value="">All Areas</option>
                    {expertAreas.map((area) => (
                      <option key={area.id} value={area.id.toString()}>
                        {area.name}
                      </option>
                    ))}
                  </select>
                )}
              </div>
              
              <div className="mb-4">
                <label className="block text-sm font-medium text-neutral-700 mb-1">
                  Nationality
                </label>
                <select
                  className="w-full px-3 py-2 bg-white border border-neutral-300 rounded-md focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
                  {...register('nationality')}
                >
                  {nationalities.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
              </div>
              
              <div className="mb-4">
                <label className="block text-sm font-medium text-neutral-700 mb-1">
                  Minimum Rating
                </label>
                <select
                  className="w-full px-3 py-2 bg-white border border-neutral-300 rounded-md focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
                  {...register('rating')}
                >
                  {ratings.map((option) => (
                    <option key={option.value} value={option.value}>
                      {option.label}
                    </option>
                  ))}
                </select>
              </div>
              
              <div className="mb-4 flex items-center">
                <input
                  type="checkbox"
                  id="isBahraini"
                  className="h-4 w-4 text-primary focus:ring-primary border-neutral-300 rounded"
                  {...register('isBahraini')}
                />
                <label htmlFor="isBahraini" className="ml-2 block text-sm text-neutral-700">
                  Bahraini experts only
                </label>
              </div>
            </>
          )}
        </div>
        
        <div className="flex justify-end space-x-3 mt-4">
          <Button 
            type="button" 
            variant="outline" 
            onClick={handleReset}
          >
            Reset Filters
          </Button>
          <Button 
            type="submit"
            variant="primary"
          >
            Apply Filters
          </Button>
        </div>
      </form>
    </div>
  );
};

export default ExpertFilters;