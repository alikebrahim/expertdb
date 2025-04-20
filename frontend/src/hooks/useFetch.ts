import { useState, useEffect, useCallback } from 'react';
import { useUI } from './useUI';

interface UseFetchOptions<T> {
  initialData?: T;
  onSuccess?: (data: T) => void;
  onError?: (error: Error) => void;
  showErrorNotifications?: boolean;
  errorMessage?: string;
  loadingDelay?: number;
  deps?: any[];
  autoFetch?: boolean;
}

export function useFetch<T>(
  fetchFn: () => Promise<T>,
  options: UseFetchOptions<T> = {}
) {
  const {
    initialData,
    onSuccess,
    onError,
    showErrorNotifications = true,
    errorMessage = 'Failed to fetch data',
    loadingDelay = 300,
    deps = [],
    autoFetch = true,
  } = options;

  const [data, setData] = useState<T | undefined>(initialData);
  const [isLoading, setIsLoading] = useState(autoFetch);
  const [error, setError] = useState<Error | null>(null);
  const { addNotification } = useUI();
  
  // This allows us to delay showing loading indicators to prevent flicker
  const [shouldShowLoading, setShouldShowLoading] = useState(false);

  const fetchData = useCallback(async () => {
    setError(null);
    setIsLoading(true);
    
    // Set a timeout to show loading state if fetch takes longer than loadingDelay
    const loadingTimer = setTimeout(() => {
      setShouldShowLoading(true);
    }, loadingDelay);
    
    try {
      const result = await fetchFn();
      setData(result);
      
      if (onSuccess) {
        onSuccess(result);
      }
      
      return result;
    } catch (err) {
      const error = err instanceof Error ? err : new Error(errorMessage);
      setError(error);
      
      if (showErrorNotifications) {
        addNotification({
          type: 'error',
          message: errorMessage || error.message,
          duration: 5000,
        });
      }
      
      if (onError) {
        onError(error);
      }
      
      return undefined;
    } finally {
      clearTimeout(loadingTimer);
      setShouldShowLoading(false);
      setIsLoading(false);
    }
  }, [fetchFn, onSuccess, onError, errorMessage, showErrorNotifications, addNotification, loadingDelay]);

  const refetch = useCallback(() => {
    return fetchData();
  }, [fetchData]);

  const reset = useCallback(() => {
    setData(initialData);
    setError(null);
    setIsLoading(false);
    setShouldShowLoading(false);
  }, [initialData]);

  // Auto-fetch on mount and when deps change if autoFetch is true
  useEffect(() => {
    if (autoFetch) {
      fetchData();
    }
  }, [...deps, autoFetch]); // eslint-disable-line react-hooks/exhaustive-deps

  return {
    data,
    setData,
    isLoading: isLoading && shouldShowLoading,
    isLoadingAny: isLoading,
    error,
    refetch,
    reset,
  };
}