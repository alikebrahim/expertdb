import { useState, useCallback } from 'react';
import { useUI } from './useUI';

interface OptimisticUpdateOptions<T> {
  onSuccess?: (data: T) => void;
  onError?: (error: Error, originalData: T) => void;
  successMessage?: string;
  errorMessage?: string;
}

/**
 * A hook for performing optimistic updates
 * @param updateFn - The function to call to perform the actual update
 * @param options - Options for handling success and errors
 */
export function useOptimisticUpdate<T, R = unknown>(
  updateFn: (data: T) => Promise<R>,
  options: OptimisticUpdateOptions<T> = {}
) {
  const { addNotification } = useUI();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);

  const execute = useCallback(
    async (data: T, optimisticUpdateFn?: (data: T) => void) => {
      setIsLoading(true);
      setError(null);
      
      // Store the original state before making optimistic updates
      const originalData = { ...data };
      
      // Perform optimistic update if provided
      if (optimisticUpdateFn) {
        optimisticUpdateFn(data);
      }
      
      try {
        // Perform the actual update
        const result = await updateFn(data);
        
        // Handle success
        if (options.successMessage) {
          addNotification({
            type: 'success',
            message: options.successMessage,
            duration: 5000,
          });
        }
        
        if (options.onSuccess) {
          options.onSuccess(data);
        }
        
        return result;
      } catch (err) {
        // Handle error
        const error = err instanceof Error ? err : new Error('Unknown error occurred');
        setError(error);
        
        // Show error notification
        addNotification({
          type: 'error',
          message: options.errorMessage || error.message,
          duration: 5000,
        });
        
        // Call onError callback
        if (options.onError) {
          options.onError(error, originalData);
        }
        
        throw error;
      } finally {
        setIsLoading(false);
      }
    },
    [updateFn, options, addNotification]
  );

  return {
    execute,
    isLoading,
    error,
  };
}

/**
 * A hook for optimistically updating a collection of items
 */
export function useOptimisticCollection<T extends { id: string | number }>(
  fetchItems: () => Promise<T[]>,
  options: {
    onSuccess?: (items: T[]) => void;
    onError?: (error: Error) => void;
  } = {}
) {
  const [items, setItems] = useState<T[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<Error | null>(null);
  const { addNotification } = useUI();

  // Load items
  const loadItems = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const data = await fetchItems();
      setItems(data);
      
      if (options.onSuccess) {
        options.onSuccess(data);
      }
      
      return data;
    } catch (err) {
      const error = err instanceof Error ? err : new Error('Failed to load items');
      setError(error);
      
      addNotification({
        type: 'error',
        message: `Error loading items: ${error.message}`,
        duration: 5000,
      });
      
      if (options.onError) {
        options.onError(error);
      }
      
      return [];
    } finally {
      setIsLoading(false);
    }
  }, [fetchItems, options, addNotification]);

  // Add item optimistically
  const addItem = useCallback(
    async (
      item: T,
      addFn: (item: T) => Promise<T>,
      options: {
        successMessage?: string;
        errorMessage?: string;
      } = {}
    ) => {
      // Add optimistically
      setItems((prev) => [...prev, item]);
      
      try {
        const result = await addFn(item);
        
        // Update with actual result
        setItems((prev) =>
          prev.map((i) => (i.id === item.id ? result : i))
        );
        
        if (options.successMessage) {
          addNotification({
            type: 'success',
            message: options.successMessage,
            duration: 5000,
          });
        }
        
        return result;
      } catch (err) {
        // Revert optimistic update
        setItems((prev) => prev.filter((i) => i.id !== item.id));
        
        const error = err instanceof Error ? err : new Error('Failed to add item');
        
        addNotification({
          type: 'error',
          message: options.errorMessage || `Error adding item: ${error.message}`,
          duration: 5000,
        });
        
        throw error;
      }
    },
    [addNotification]
  );

  // Update item optimistically
  const updateItem = useCallback(
    async (
      updatedItem: T,
      updateFn: (item: T) => Promise<T>,
      options: {
        successMessage?: string;
        errorMessage?: string;
      } = {}
    ) => {
      // Store original for rollback
      const originalItem = items.find((i) => i.id === updatedItem.id);
      
      if (!originalItem) {
        throw new Error(`Item with id ${updatedItem.id} not found`);
      }
      
      // Update optimistically
      setItems((prev) =>
        prev.map((i) => (i.id === updatedItem.id ? updatedItem : i))
      );
      
      try {
        const result = await updateFn(updatedItem);
        
        // Update with actual result
        setItems((prev) =>
          prev.map((i) => (i.id === updatedItem.id ? result : i))
        );
        
        if (options.successMessage) {
          addNotification({
            type: 'success',
            message: options.successMessage,
            duration: 5000,
          });
        }
        
        return result;
      } catch (err) {
        // Revert optimistic update
        setItems((prev) =>
          prev.map((i) => (i.id === updatedItem.id ? originalItem : i))
        );
        
        const error = err instanceof Error ? err : new Error('Failed to update item');
        
        addNotification({
          type: 'error',
          message: options.errorMessage || `Error updating item: ${error.message}`,
          duration: 5000,
        });
        
        throw error;
      }
    },
    [items, addNotification]
  );

  // Delete item optimistically
  const deleteItem = useCallback(
    async (
      id: string | number,
      deleteFn: (id: string | number) => Promise<void>,
      options: {
        successMessage?: string;
        errorMessage?: string;
      } = {}
    ) => {
      // Store original for rollback
      const originalItem = items.find((i) => i.id === id);
      
      if (!originalItem) {
        throw new Error(`Item with id ${id} not found`);
      }
      
      // Delete optimistically
      setItems((prev) => prev.filter((i) => i.id !== id));
      
      try {
        await deleteFn(id);
        
        if (options.successMessage) {
          addNotification({
            type: 'success',
            message: options.successMessage,
            duration: 5000,
          });
        }
      } catch (err) {
        // Revert optimistic update
        setItems((prev) => [...prev, originalItem]);
        
        const error = err instanceof Error ? err : new Error('Failed to delete item');
        
        addNotification({
          type: 'error',
          message: options.errorMessage || `Error deleting item: ${error.message}`,
          duration: 5000,
        });
        
        throw error;
      }
    },
    [items, addNotification]
  );

  return {
    items,
    setItems,
    isLoading,
    error,
    loadItems,
    addItem,
    updateItem,
    deleteItem,
  };
}