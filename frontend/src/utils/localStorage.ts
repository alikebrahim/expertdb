import { SortConfig } from '../components/tables';

// LocalStorage key prefixes
const FILTER_KEY = 'expertdb_filters';
const SORT_KEY = 'expertdb_sort';

/**
 * Save expert filters to localStorage
 * @param filters The filters to save
 */
export const saveFilters = (filters: Record<string, any>): void => {
  try {
    localStorage.setItem(FILTER_KEY, JSON.stringify(filters));
  } catch (error) {
    console.error('Error saving filters to localStorage:', error);
  }
};

/**
 * Load expert filters from localStorage
 * @returns The saved filters or an empty object
 */
export const loadFilters = (): Record<string, any> => {
  try {
    const savedFilters = localStorage.getItem(FILTER_KEY);
    return savedFilters ? JSON.parse(savedFilters) : {};
  } catch (error) {
    console.error('Error loading filters from localStorage:', error);
    return {};
  }
};

/**
 * Save sort configuration to localStorage
 * @param sortConfig The sort configuration to save
 */
export const saveSort = (sortConfig: SortConfig): void => {
  try {
    localStorage.setItem(SORT_KEY, JSON.stringify(sortConfig));
  } catch (error) {
    console.error('Error saving sort configuration to localStorage:', error);
  }
};

/**
 * Load sort configuration from localStorage
 * @returns The saved sort configuration or default (name, asc)
 */
export const loadSort = (): SortConfig => {
  try {
    const savedSort = localStorage.getItem(SORT_KEY);
    return savedSort ? JSON.parse(savedSort) : { field: 'name', direction: 'asc' };
  } catch (error) {
    console.error('Error loading sort configuration from localStorage:', error);
    return { field: 'name', direction: 'asc' };
  }
};

/**
 * Clear all saved filters and sort preferences
 */
export const clearSavedPreferences = (): void => {
  try {
    localStorage.removeItem(FILTER_KEY);
    localStorage.removeItem(SORT_KEY);
  } catch (error) {
    console.error('Error clearing saved preferences from localStorage:', error);
  }
};