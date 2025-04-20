import { zodResolver } from '@hookform/resolvers/zod';
import { useForm as useHookForm, UseFormProps, FieldValues, UseFormReturn } from 'react-hook-form';
import { ZodSchema, z } from 'zod';
import { useUI } from './useUI';

interface UseZodFormProps<T extends FieldValues> extends Omit<UseFormProps<T>, 'resolver'> {
  schema: ZodSchema<T>;
  onError?: (error: string) => void;
}

export function useZodForm<T extends FieldValues>({
  schema,
  onError,
  ...formProps
}: UseZodFormProps<T>): UseFormReturn<T> {
  const { addNotification } = useUI();
  
  return useHookForm<T>({
    ...formProps,
    resolver: zodResolver(schema, {
      errorMap: (error, ctx) => {
        const message = error.message || `Invalid ${error.path.join('.')}`;
        if (onError) {
          onError(message);
        } else {
          addNotification({
            type: 'error',
            message,
            duration: 5000,
          });
        }
        return { message };
      },
    }),
  });
}

export function useFormWithNotifications<T extends FieldValues>({
  schema,
  onSuccess,
  ...formProps
}: UseZodFormProps<T> & { onSuccess?: (data: T) => void }): UseFormReturn<T> & {
  handleSubmitWithNotifications: (callback: (data: T) => Promise<{ success: boolean; message?: string }>) => (e?: React.BaseSyntheticEvent) => Promise<void>;
} {
  const form = useZodForm<T>({ schema, ...formProps });
  const { addNotification } = useUI();

  const handleSubmitWithNotifications = (
    callback: (data: T) => Promise<{ success: boolean; message?: string }>
  ) => {
    return form.handleSubmit(async (data) => {
      try {
        const result = await callback(data);
        
        if (result.success) {
          addNotification({
            type: 'success',
            message: result.message || 'Operation completed successfully',
            duration: 5000,
          });
          
          if (onSuccess) {
            onSuccess(data);
          }
        } else {
          addNotification({
            type: 'error',
            message: result.message || 'Operation failed',
            duration: 5000,
          });
        }
      } catch (error) {
        addNotification({
          type: 'error',
          message: error instanceof Error ? error.message : 'An unexpected error occurred',
          duration: 5000,
        });
      }
    });
  };

  return {
    ...form,
    handleSubmitWithNotifications,
  };
}
