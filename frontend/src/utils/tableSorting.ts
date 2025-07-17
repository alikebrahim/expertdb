import { Expert } from '../types';

export interface SortConfig {
  field: string;
  direction: 'asc' | 'desc';
}

export const sortExperts = (
  experts: Expert[],
  sortConfig: SortConfig
): Expert[] => {
  if (!sortConfig.field) return experts;
  
  return [...experts].sort((a, b) => {
    const aValue = getNestedValue(a, sortConfig.field);
    const bValue = getNestedValue(b, sortConfig.field);
    
    // Handle null/undefined values
    if (aValue == null && bValue == null) return 0;
    if (aValue == null) return 1;
    if (bValue == null) return -1;
    
    // Handle different data types
    if (typeof aValue === 'string' && typeof bValue === 'string') {
      const comparison = aValue.localeCompare(bValue);
      return sortConfig.direction === 'asc' ? comparison : -comparison;
    }
    
    if (typeof aValue === 'number' && typeof bValue === 'number') {
      const comparison = aValue - bValue;
      return sortConfig.direction === 'asc' ? comparison : -comparison;
    }
    
    if (aValue instanceof Date && bValue instanceof Date) {
      const comparison = aValue.getTime() - bValue.getTime();
      return sortConfig.direction === 'asc' ? comparison : -comparison;
    }
    
    // Boolean values
    if (typeof aValue === 'boolean' && typeof bValue === 'boolean') {
      const comparison = aValue === bValue ? 0 : aValue ? 1 : -1;
      return sortConfig.direction === 'asc' ? comparison : -comparison;
    }
    
    // Convert to string for comparison as fallback
    const aString = String(aValue).toLowerCase();
    const bString = String(bValue).toLowerCase();
    const comparison = aString.localeCompare(bString);
    return sortConfig.direction === 'asc' ? comparison : -comparison;
  });
};

const getNestedValue = (obj: any, path: string): any => {
  return path.split('.').reduce((current, key) => {
    if (current && typeof current === 'object' && key in current) {
      return current[key];
    }
    return undefined;
  }, obj);
};

export const getDefaultSortConfig = (): SortConfig => ({
  field: 'name',
  direction: 'asc'
});

export const toggleSortDirection = (currentConfig: SortConfig, field: string): SortConfig => {
  if (currentConfig.field === field) {
    return {
      field,
      direction: currentConfig.direction === 'asc' ? 'desc' : 'asc'
    };
  }
  
  return {
    field,
    direction: 'asc'
  };
};

export const getSortIndicator = (sortConfig: SortConfig, field: string): string => {
  if (sortConfig.field !== field) return '';
  return sortConfig.direction === 'asc' ? '↑' : '↓';
};

export const getSortableFields = (): Array<{ key: string; label: string }> => [
  { key: 'name', label: 'Name' },
  { key: 'role', label: 'Role' },
  { key: 'employmentType', label: 'Employment Type' },
  { key: 'affiliation', label: 'Affiliation' },
  { key: 'rating', label: 'Rating' },
  { key: 'nationality', label: 'Nationality' },
  { key: 'created_at', label: 'Date Added' },
  { key: 'updated_at', label: 'Last Updated' }
];