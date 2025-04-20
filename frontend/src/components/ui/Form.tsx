import React from 'react';
import { FieldValues, UseFormReturn } from 'react-hook-form';

interface FormProps<T extends FieldValues> extends React.FormHTMLAttributes<HTMLFormElement> {
  form: UseFormReturn<T>;
  onSubmit: (data: T) => void | Promise<void>;
  children: React.ReactNode;
  className?: string;
  resetOnSuccess?: boolean;
  submitText?: string;
  isSubmitting?: boolean;
  showSubmitButton?: boolean;
  submitButtonPosition?: 'left' | 'center' | 'right';
  resetText?: string;
  showResetButton?: boolean;
  onReset?: () => void;
}

export const Form = <T extends FieldValues>({
  form,
  onSubmit,
  children,
  className = '',
  resetOnSuccess = false,
  submitText = 'Submit',
  isSubmitting: externalIsSubmitting,
  showSubmitButton = true,
  submitButtonPosition = 'right',
  resetText = 'Reset',
  showResetButton = false,
  onReset,
  ...rest
}: FormProps<T>) => {
  const {
    handleSubmit,
    reset,
    formState: { isSubmitting: internalIsSubmitting },
  } = form;

  const isSubmitting = externalIsSubmitting !== undefined ? externalIsSubmitting : internalIsSubmitting;

  const handleFormSubmit = async (data: T) => {
    await onSubmit(data);
    if (resetOnSuccess) {
      reset();
    }
  };

  const positionClassMap = {
    left: 'justify-start',
    center: 'justify-center',
    right: 'justify-end',
  };

  return (
    <form
      className={`space-y-4 ${className}`}
      onSubmit={handleSubmit(handleFormSubmit)}
      {...rest}
    >
      {children}

      {(showSubmitButton || showResetButton) && (
        <div className={`flex space-x-4 mt-6 ${positionClassMap[submitButtonPosition]}`}>
          {showResetButton && (
            <button
              type="button"
              onClick={() => {
                reset();
                if (onReset) onReset();
              }}
              className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 dark:bg-gray-700 dark:text-gray-200 dark:border-gray-600 dark:hover:bg-gray-600"
              disabled={isSubmitting}
            >
              {resetText}
            </button>
          )}

          {showSubmitButton && (
            <button
              type="submit"
              className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 dark:bg-blue-700 dark:hover:bg-blue-800 disabled:opacity-50 disabled:cursor-not-allowed"
              disabled={isSubmitting}
            >
              {isSubmitting ? (
                <div className="flex items-center">
                  <svg
                    className="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle
                      className="opacity-25"
                      cx="12"
                      cy="12"
                      r="10"
                      stroke="currentColor"
                      strokeWidth="4"
                    ></circle>
                    <path
                      className="opacity-75"
                      fill="currentColor"
                      d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                    ></path>
                  </svg>
                  Processing...
                </div>
              ) : (
                submitText
              )}
            </button>
          )}
        </div>
      )}
    </form>
  );
};