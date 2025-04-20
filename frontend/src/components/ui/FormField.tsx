import React from 'react';
import { FieldValues, UseFormReturn, Path, FieldError } from 'react-hook-form';

interface FormFieldProps<T extends FieldValues> {
  form: UseFormReturn<T>;
  name: Path<T>;
  label: string;
  type?: 'text' | 'email' | 'password' | 'number' | 'date' | 'textarea' | 'select' | 'checkbox' | 'radio';
  placeholder?: string;
  className?: string;
  disabled?: boolean;
  required?: boolean;
  options?: { label: string; value: string | number }[];
  children?: React.ReactNode;
  rows?: number;
  hint?: string;
}

export const FormField = <T extends FieldValues>({
  form,
  name,
  label,
  type = 'text',
  placeholder = '',
  className = '',
  disabled = false,
  required = false,
  options = [],
  children,
  rows = 3,
  hint,
}: FormFieldProps<T>) => {
  const {
    register,
    formState: { errors, isSubmitting },
  } = form;

  const error = errors[name] as FieldError | undefined;
  const isCheckboxOrRadio = type === 'checkbox' || type === 'radio';

  return (
    <div className={`mb-4 ${className}`}>
      {!isCheckboxOrRadio && (
        <label
          className="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1"
          htmlFor={name}
        >
          {label}
          {required && <span className="text-red-500 ml-1">*</span>}
        </label>
      )}

      {type === 'textarea' ? (
        <textarea
          id={name}
          rows={rows}
          className={`block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-600 dark:text-white dark:placeholder-gray-400 ${
            error ? 'border-red-500 focus:border-red-500 focus:ring-red-500' : ''
          }`}
          placeholder={placeholder}
          disabled={disabled || isSubmitting}
          {...register(name)}
        />
      ) : type === 'select' ? (
        <select
          id={name}
          className={`block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-600 dark:text-white ${
            error ? 'border-red-500 focus:border-red-500 focus:ring-red-500' : ''
          }`}
          disabled={disabled || isSubmitting}
          {...register(name)}
        >
          {options.map((option) => (
            <option key={option.value} value={option.value}>
              {option.label}
            </option>
          ))}
        </select>
      ) : type === 'checkbox' ? (
        <div className="flex items-start">
          <div className="flex items-center h-5">
            <input
              id={name}
              type="checkbox"
              className={`w-4 h-4 border-gray-300 rounded focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 ${
                error ? 'border-red-500 focus:ring-red-500' : ''
              }`}
              disabled={disabled || isSubmitting}
              {...register(name)}
            />
          </div>
          <div className="ml-3 text-sm">
            <label
              className="font-medium text-gray-700 dark:text-gray-300"
              htmlFor={name}
            >
              {label}
              {required && <span className="text-red-500 ml-1">*</span>}
            </label>
            {hint && <p className="text-gray-500 dark:text-gray-400">{hint}</p>}
          </div>
        </div>
      ) : type === 'radio' ? (
        <div>
          {options.map((option) => (
            <div key={option.value} className="flex items-center mb-1">
              <input
                id={`${name}-${option.value}`}
                type="radio"
                value={option.value}
                className={`w-4 h-4 border-gray-300 focus:ring-2 focus:ring-blue-500 dark:bg-gray-700 dark:border-gray-600 ${
                  error ? 'border-red-500 focus:ring-red-500' : ''
                }`}
                disabled={disabled || isSubmitting}
                {...register(name)}
              />
              <label
                className="ml-2 block text-sm font-medium text-gray-700 dark:text-gray-300"
                htmlFor={`${name}-${option.value}`}
              >
                {option.label}
              </label>
            </div>
          ))}
        </div>
      ) : (
        <>
          <input
            id={name}
            type={type}
            className={`block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 dark:bg-gray-800 dark:border-gray-600 dark:text-white dark:placeholder-gray-400 ${
              error ? 'border-red-500 focus:border-red-500 focus:ring-red-500' : ''
            }`}
            placeholder={placeholder}
            disabled={disabled || isSubmitting}
            {...register(name)}
          />
          {children}
        </>
      )}

      {hint && !isCheckboxOrRadio && (
        <p className="mt-1 text-sm text-gray-500 dark:text-gray-400">{hint}</p>
      )}

      {error && (
        <p className="mt-1 text-sm text-red-500">{error.message as string}</p>
      )}
    </div>
  );
};