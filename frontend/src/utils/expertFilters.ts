import { Expert } from '../types';

export interface ExpertFilters {
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
  searchTerm?: string;
}

export const applyFilters = (
  experts: Expert[],
  filters: ExpertFilters,
  searchTerm?: string
): Expert[] => {
  return experts.filter(expert => {
    // Global search term
    if (searchTerm) {
      const searchLower = searchTerm.toLowerCase();
      const searchableFields = [
        expert.name,
        expert.affiliation,
        expert.specializedArea,
        expert.role,
        expert.employmentType,
        expert.institution
      ];
      
      const matchesSearch = searchableFields.some(field => 
        field?.toLowerCase().includes(searchLower)
      );
      
      if (!matchesSearch) return false;
    }
    
    // Name filter
    if (filters.name) {
      const nameLower = filters.name.toLowerCase();
      if (!expert.name.toLowerCase().includes(nameLower)) {
        return false;
      }
    }
    
    // Role filter
    if (filters.role && expert.role !== filters.role) {
      return false;
    }
    
    // Employment type filter
    if (filters.type && expert.employmentType !== filters.type) {
      return false;
    }
    
    // Affiliation filter
    if (filters.affiliation) {
      const affiliationLower = filters.affiliation.toLowerCase();
      if (!expert.affiliation?.toLowerCase().includes(affiliationLower)) {
        return false;
      }
    }
    
    // Expert area filter
    if (filters.expertAreaId) {
      const areaId = parseInt(filters.expertAreaId);
      if (expert.generalArea !== areaId) {
        return false;
      }
    }
    
    // Nationality filter - skip for now since nationality is not in Expert type
    // if (filters.nationality && expert.nationality !== filters.nationality) {
    //   return false;
    // }
    
    // Rating filter (minimum rating)
    if (filters.rating) {
      const minRating = parseFloat(filters.rating);
      const expertRating = parseFloat(expert.rating);
      if (expertRating < minRating) {
        return false;
      }
    }
    
    // Boolean filters
    if (filters.isAvailable !== undefined && expert.isAvailable !== filters.isAvailable) {
      return false;
    }
    
    if (filters.isBahraini !== undefined && expert.isBahraini !== filters.isBahraini) {
      return false;
    }
    
    if (filters.isPublished !== undefined && expert.isPublished !== filters.isPublished) {
      return false;
    }
    
    return true;
  });
};

export const getDefaultFilters = (): ExpertFilters => ({
  name: '',
  role: '',
  type: '',
  affiliation: '',
  expertAreaId: '',
  nationality: '',
  rating: '',
  isAvailable: undefined,
  isBahraini: undefined,
  isPublished: undefined,
  searchTerm: ''
});

export const getActiveFilterCount = (filters: ExpertFilters): number => {
  let count = 0;
  
  if (filters.name) count++;
  if (filters.role) count++;
  if (filters.type) count++;
  if (filters.affiliation) count++;
  if (filters.expertAreaId) count++;
  if (filters.nationality) count++;
  if (filters.rating) count++;
  if (filters.isAvailable !== undefined) count++;
  if (filters.isBahraini !== undefined) count++;
  if (filters.isPublished !== undefined) count++;
  if (filters.searchTerm) count++;
  
  return count;
};

export const getFilterSummary = (filters: ExpertFilters): string[] => {
  const summary: string[] = [];
  
  if (filters.name) summary.push(`Name: ${filters.name}`);
  if (filters.role) summary.push(`Role: ${filters.role}`);
  if (filters.type) summary.push(`Type: ${filters.type}`);
  if (filters.affiliation) summary.push(`Affiliation: ${filters.affiliation}`);
  if (filters.nationality) summary.push(`Nationality: ${filters.nationality}`);
  if (filters.rating) summary.push(`Min Rating: ${filters.rating}`);
  if (filters.isAvailable !== undefined) summary.push(`Available: ${filters.isAvailable ? 'Yes' : 'No'}`);
  if (filters.isBahraini !== undefined) summary.push(`Bahraini: ${filters.isBahraini ? 'Yes' : 'No'}`);
  if (filters.isPublished !== undefined) summary.push(`Published: ${filters.isPublished ? 'Yes' : 'No'}`);
  
  return summary;
};