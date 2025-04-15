import { useState, useEffect } from 'react';
import { Expert } from '../types';
import { expertsApi } from '../services/api';
import ExpertFilters from '../components/ExpertFilters';
import ExpertTable from '../components/ExpertTable';

interface ExpertFilters {
  name?: string;
  role?: string;
  type?: string;
  affiliation?: string;
  expertArea?: string;
  nationality?: string;
  isAvailable?: boolean;
}

const SearchPage = () => {
  const [experts, setExperts] = useState<Expert[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [filters, setFilters] = useState<ExpertFilters>({});
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [limit] = useState(10);
  
  // Fetch experts on mount and when filters or pagination change
  useEffect(() => {
    const fetchExperts = async () => {
      setIsLoading(true);
      setError(null);
      
      try {
        // Convert filters to API-friendly params
        const params: Record<string, string | boolean> = {};
        
        if (filters.name) params.name = filters.name;
        if (filters.role) params.role = filters.role;
        if (filters.type) params.employmentType = filters.type;
        if (filters.affiliation) params.institution = filters.affiliation;
        if (filters.expertArea) params.general_area = filters.expertArea;
        if (filters.nationality) params.nationality = filters.nationality;
        if (filters.isAvailable !== undefined) params.is_available = filters.isAvailable;
        
        // Add default filter to show only available experts
        if (filters.isAvailable === undefined) {
          params.is_available = true;
        }
        
        const response = await expertsApi.getExperts(page, limit, params);
        
        if (response.success) {
          setExperts(response.data.data);
          setTotalPages(response.data.totalPages);
        } else {
          setError(response.message || 'Failed to fetch experts');
          setExperts([]);
          setTotalPages(1);
        }
      } catch (error) {
        console.error('Error fetching experts:', error);
        setError('An error occurred while fetching experts');
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchExperts();
  }, [filters, page, limit]);
  
  const handleFilterChange = (newFilters: ExpertFilters) => {
    setFilters(newFilters);
    setPage(1); // Reset to first page when filters change
  };
  
  const handlePageChange = (newPage: number) => {
    setPage(newPage);
  };
  
  return (
    <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-primary">Expert Search</h1>
        <p className="text-neutral-600">
          Search and filter experts based on various criteria
        </p>
      </div>
      
      <ExpertFilters onFilterChange={handleFilterChange} />
      
      <ExpertTable 
        experts={experts} 
        isLoading={isLoading}
        error={error}
        pagination={{
          currentPage: page,
          totalPages: totalPages,
          onPageChange: handlePageChange
        }}
      />
    </div>
  );
};

export default SearchPage;