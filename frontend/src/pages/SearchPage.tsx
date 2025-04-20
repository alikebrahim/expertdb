import { useState } from 'react';
import { Expert } from '../types';
import * as expertsApi from '../api/experts';
import { useFetch } from '../hooks';
import Layout from '../components/layout/Layout';
import { ExpertFilters } from '../components/filters';
import { ExpertTable } from '../components/tables';
import { ProgressStepper } from '../components/ui';

interface ExpertFiltersType {
  name?: string;
  role?: string;
  type?: string;
  affiliation?: string;
  expertAreaId?: string;
  nationality?: string;
  rating?: string;
  isAvailable?: boolean;
  isBahraini?: boolean;
}

const SearchPage = () => {
  const [filters, setFilters] = useState<ExpertFiltersType>({
    isAvailable: true // Default to show only available experts
  });
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [currentStep, setCurrentStep] = useState(0);
  
  // Create fetch function based on current filters and pagination
  const fetchExperts = async () => {
    // Convert filters to API-friendly params
    const params: Record<string, string | boolean | number> = {};
    
    if (filters.name) params.name = filters.name;
    if (filters.role) params.role = filters.role;
    if (filters.type) params.employmentType = filters.type;
    if (filters.affiliation) params.institution = filters.affiliation;
    if (filters.expertAreaId) params.generalArea = parseInt(filters.expertAreaId);
    if (filters.nationality) params.nationality = filters.nationality;
    if (filters.rating) params.minRating = parseInt(filters.rating);
    if (filters.isAvailable !== undefined) params.isAvailable = filters.isAvailable;
    if (filters.isBahraini !== undefined) params.isBahraini = filters.isBahraini;
    
    const offset = (page - 1) * limit;
    const response = await expertsApi.getExperts(limit, offset, params);
    
    if (response.success && response.data) {
      return {
        experts: response.data.experts,
        totalPages: response.data.pagination.totalPages
      };
    } else {
      throw new Error(response.message || 'Failed to fetch experts');
    }
  };
  
  // Use the fetch hook
  const { 
    data, 
    isLoading, 
    error, 
    refetch 
  } = useFetch(fetchExperts, {
    errorMessage: 'Failed to fetch experts',
    deps: [filters, page, limit],
  });
  
  const handleFilterChange = (newFilters: ExpertFiltersType) => {
    setFilters(newFilters);
    setPage(1); // Reset to first page when filters change
  };
  
  const handlePageChange = (newPage: number) => {
    setPage(newPage);
  };
  
  // Define search process steps
  const searchSteps = [
    { id: 'filter', label: 'Define filters', description: 'Set search criteria' },
    { id: 'results', label: 'View results', description: 'Browse matching experts' },
    { id: 'contact', label: 'Contact experts', description: 'Reach out to selected experts' },
  ];

  return (
    <Layout>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-primary">Expert Search</h1>
        <p className="text-neutral-600">
          Search and filter experts based on various criteria
        </p>
        
        <div className="my-6">
          <ProgressStepper 
            steps={searchSteps} 
            currentStep={currentStep}
            onStepClick={setCurrentStep}
            showPercentage
          />
        </div>
      </div>
      
      <ExpertFilters onFilterChange={handleFilterChange} initialFilters={filters} />
      
      <div className="bg-white shadow rounded-lg overflow-hidden mt-6">
        <ExpertTable 
          experts={data?.experts || []} 
          isLoading={isLoading}
          error={error ? error.message : null}
          pagination={{
            currentPage: page,
            totalPages: data?.totalPages || 1,
            onPageChange: handlePageChange
          }}
        />
      </div>
    </Layout>
  );
};

export default SearchPage;