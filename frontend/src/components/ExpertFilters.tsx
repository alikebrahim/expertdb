import { useState } from 'react';
import { useForm } from 'react-hook-form';
import Input from './ui/Input';
import Button from './ui/Button';

interface FiltersFormData {
  name?: string;
  role?: string;
  type?: string;
  affiliation?: string;
  expertArea?: string;
  nationality?: string;
  isAvailable?: boolean;
}

interface ExpertFiltersProps {
  onFilterChange: (filters: FiltersFormData) => void;
}

const ExpertFilters = ({ onFilterChange }: ExpertFiltersProps) => {
  const [isExpanded, setIsExpanded] = useState(false);
  const { register, handleSubmit, reset } = useForm<FiltersFormData>();
  
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

  const toggleExpanded = () => {
    setIsExpanded(!isExpanded);
  };
  
  const onSubmit = (data: FiltersFormData) => {
    // Convert empty strings to undefined
    Object.keys(data).forEach(key => {
      const k = key as keyof FiltersFormData;
      if (data[k] === '') {
        data[k] = undefined;
      }
    });
    
    onFilterChange(data);
  };
  
  const handleReset = () => {
    reset({
      name: '',
      role: '',
      type: '',
      affiliation: '',
      expertArea: '',
      nationality: '',
      isAvailable: undefined
    });
    
    onFilterChange({});
  };
  
  return (
    <div className="bg-white rounded-md shadow p-4 mb-6">
      <div className="flex justify-between items-center mb-4">
        <h2 className="text-lg font-semibold text-primary">Filter Experts</h2>
        <button 
          onClick={toggleExpanded}
          className="text-primary hover:text-primary-light"
        >
          {isExpanded ? 'Hide Filters' : 'Show All Filters'}
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
              className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
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
              className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
              {...register('type')}
            >
              {types.map((option) => (
                <option key={option.value} value={option.value}>
                  {option.label}
                </option>
              ))}
            </select>
          </div>
          
          {/* Advanced filters */}
          {isExpanded && (
            <>
              <Input
                label="Affiliation"
                placeholder="Institution or company"
                {...register('affiliation')}
              />
              
              <Input
                label="Expert Area"
                placeholder="Area of expertise"
                {...register('expertArea')}
              />
              
              <Input
                label="Nationality"
                placeholder="Expert nationality"
                {...register('nationality')}
              />
              
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
            </>
          )}
        </div>
        
        <div className="flex justify-end space-x-3 mt-4">
          <Button 
            type="button" 
            variant="outline" 
            onClick={handleReset}
          >
            Reset
          </Button>
          <Button type="submit">
            Apply Filters
          </Button>
        </div>
      </form>
    </div>
  );
};

export default ExpertFilters;