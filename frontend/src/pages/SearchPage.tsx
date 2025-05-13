import { useState, useEffect } from 'react';
import * as expertsApi from '../api/experts';
import { useFetch } from '../hooks';
import { useUI } from '../hooks/useUI';
import Layout from '../components/layout/Layout';
import { ExpertFilters } from '../components/filters';
import { ExpertTable, SortConfig } from '../components/tables';
import { ProgressStepper, LoadingSpinner } from '../components/ui';
import { saveFilters, loadFilters, saveSort, loadSort } from '../utils/localStorage';

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
  isPublished?: boolean;
}

const SearchPage = () => {
  // Initialize filters from localStorage or use default
  const [filters, setFilters] = useState<ExpertFiltersType>(() => {
    const savedFilters = loadFilters();
    // Always include isAvailable if not explicitly set in saved filters
    return {
      ...savedFilters,
      isAvailable: savedFilters.isAvailable !== undefined ? savedFilters.isAvailable : true
    };
  });
  
  const [page, setPage] = useState(1);
  const [limit] = useState(10);
  const [currentStep, setCurrentStep] = useState(0);
  const [totalResults, setTotalResults] = useState(0);
  
  // Initialize sort config from localStorage or use default
  const [sortConfig, setSortConfig] = useState<SortConfig>(() => loadSort());
  
  const { addNotification } = useUI();
  
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
    if (filters.isPublished !== undefined) params.isPublished = filters.isPublished;
    
    // Add sorting parameters
    if (sortConfig) {
      params.sort_by = sortConfig.field;
      params.sort_order = sortConfig.direction;
    }
    
    const offset = (page - 1) * limit;
    const response = await expertsApi.getExperts(limit, offset, params);
    
    if (response.success && response.data) {
      // Update total results for display
      setTotalResults(response.data.pagination.totalCount);
      
      return {
        experts: response.data.experts,
        totalPages: response.data.pagination.totalPages,
        pagination: response.data.pagination
      };
    } else {
      throw new Error(response.message || 'Failed to fetch experts');
    }
  };
  
  // Use the fetch hook
  const { 
    data, 
    isLoading, 
    error
  } = useFetch(fetchExperts, {
    errorMessage: 'Failed to fetch experts',
    deps: [filters, page, limit, sortConfig],
  });
  
  // For the first load, show a loading state
  const [initialLoad, setInitialLoad] = useState(true);
  
  useEffect(() => {
    if (!isLoading && initialLoad) {
      setInitialLoad(false);
    }
  }, [isLoading, initialLoad]);
  
  // Update to step 2 if results are found
  useEffect(() => {
    if (data?.experts && data.experts.length > 0 && currentStep === 0) {
      setCurrentStep(1); // Move to results step
    }
  }, [data, currentStep]);
  
  const handleFilterChange = (newFilters: ExpertFiltersType) => {
    setFilters(newFilters);
    setPage(1); // Reset to first page when filters change
    
    // Move to step 1 (filter step) when changing filters
    setCurrentStep(0);
    
    // Save filters to localStorage
    saveFilters(newFilters);
  };
  
  const handlePageChange = (newPage: number) => {
    setPage(newPage);
    
    // Scroll to top when changing pages
    window.scrollTo({
      top: 0,
      behavior: 'smooth'
    });
  };
  
  const handleSort = (field: string) => {
    setSortConfig(prevConfig => {
      // If clicking on the same field, toggle direction
      let newConfig: SortConfig;
      
      if (prevConfig.field === field) {
        newConfig = {
          ...prevConfig,
          direction: prevConfig.direction === 'asc' ? 'desc' : 'asc'
        };
      } else {
        // If clicking on a new field, sort ascending by default
        newConfig = {
          field,
          direction: 'asc'
        };
      }
      
      // Save sort config to localStorage
      saveSort(newConfig);
      
      return newConfig;
    });
    
    // Reset to first page when sorting changes
    setPage(1);
    
    // Notify user about sort change
    addNotification({
      type: 'info',
      message: `Sorted by ${field}`,
      duration: 2000
    });
  };
  
  // Define search process steps
  const searchSteps = [
    { id: 'filter', label: 'Define filters', description: 'Set search criteria' },
    { id: 'results', label: 'View results', description: 'Browse matching experts' },
    { id: 'contact', label: 'Contact experts', description: 'Reach out to selected experts' },
  ];
  
  const handleRequestDetails = () => {
    // Move to step 3 (contact step)
    setCurrentStep(2);
    
    addNotification({
      type: 'info',
      message: 'Contact information for selected experts is available in their profiles',
      duration: 5000
    });
  };

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
      
      {/* Results summary */}
      {!initialLoad && data?.experts && (
        <div className="bg-accent bg-opacity-10 p-4 rounded-lg mb-4">
          <div className="flex justify-between items-center">
            <div>
              <h2 className="text-lg font-semibold text-primary">
                {totalResults === 0 ? 'No experts found' : `Found ${totalResults} expert${totalResults === 1 ? '' : 's'}`}
              </h2>
              <p className="text-sm text-neutral-600">
                {totalResults > 0 ? `Showing ${((page - 1) * limit) + 1} - ${Math.min(page * limit, totalResults)} of ${totalResults}` : 'Try adjusting your filters to find experts'}
              </p>
            </div>
            
            {totalResults > 0 && (
              <button
                onClick={handleRequestDetails}
                className="px-4 py-2 bg-primary text-white rounded-md hover:bg-primary-dark transition-colors duration-200"
              >
                Contact Selected Experts
              </button>
            )}
          </div>
        </div>
      )}
      
      {/* Loading state for initial load */}
      {initialLoad ? (
        <div className="flex flex-col items-center justify-center py-16 bg-white shadow rounded-lg">
          <LoadingSpinner size="lg" />
          <p className="mt-4 text-neutral-600">Loading experts...</p>
        </div>
      ) : (
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
            sortConfig={sortConfig}
            onSort={handleSort}
          />
        </div>
      )}
    </Layout>
  );
};

export default SearchPage;