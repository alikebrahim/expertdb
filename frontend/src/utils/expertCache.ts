import { Expert } from '../types';

interface ExpertCache {
  data: Expert[];
  timestamp: number;
}

const CACHE_DURATION = 5 * 60 * 1000; // 5 minutes
const CACHE_KEY = 'expertsCache';

export const getCachedExperts = (): ExpertCache | null => {
  try {
    const cached = localStorage.getItem(CACHE_KEY);
    return cached ? JSON.parse(cached) : null;
  } catch (error) {
    console.warn('Failed to parse cached experts data:', error);
    return null;
  }
};

export const setCachedExperts = (experts: Expert[]): void => {
  try {
    const cacheData: ExpertCache = {
      data: experts,
      timestamp: Date.now()
    };
    localStorage.setItem(CACHE_KEY, JSON.stringify(cacheData));
  } catch (error) {
    console.warn('Failed to cache experts data:', error);
  }
};

export const isCacheValid = (timestamp: number): boolean => {
  return Date.now() - timestamp < CACHE_DURATION;
};

export const clearExpertCache = (): void => {
  try {
    localStorage.removeItem(CACHE_KEY);
  } catch (error) {
    console.warn('Failed to clear expert cache:', error);
  }
};

export const getCacheStatus = (): { exists: boolean; valid: boolean; age: number } => {
  const cached = getCachedExperts();
  if (!cached) {
    return { exists: false, valid: false, age: 0 };
  }
  
  const age = Date.now() - cached.timestamp;
  const valid = isCacheValid(cached.timestamp);
  
  return { exists: true, valid, age };
};