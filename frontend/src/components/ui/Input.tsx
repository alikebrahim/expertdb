import { InputHTMLAttributes, forwardRef } from 'react';
import { Skeleton } from './Skeleton';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  fullWidth?: boolean;
  isLoading?: boolean;
  hint?: string;
}

const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, className = '', fullWidth = true, isLoading = false, hint, ...props }, ref) => {
    const widthClass = fullWidth ? 'w-full' : '';
    const errorClass = error ? 'border-secondary' : 'border-neutral-300';
    
    if (isLoading) {
      return (
        <div className={`mb-4 ${widthClass}`}>
          {label && <Skeleton height={20} width={120} className="mb-1" />}
          <Skeleton height={40} width="100%" className="rounded" />
        </div>
      );
    }
    
    return (
      <div className={`mb-4 ${widthClass}`}>
        {label && (
          <label className="block text-sm font-medium text-neutral-700 mb-1">
            {label}
            {props.required && <span className="ml-1 text-secondary">*</span>}
          </label>
        )}
        
        <input
          ref={ref}
          className={`
            px-3 py-2 bg-white border ${errorClass} rounded focus:outline-none
            focus:ring-1 focus:ring-primary focus:border-primary
            ${className}
            ${props.disabled ? 'bg-gray-100 cursor-not-allowed' : ''}
          `}
          {...props}
        />
        
        {hint && !error && (
          <p className="mt-1 text-xs text-gray-500">{hint}</p>
        )}
        
        {error && (
          <p className="mt-1 text-sm text-secondary">{error}</p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';

export default Input;