import { InputHTMLAttributes, forwardRef } from 'react';

interface InputProps extends InputHTMLAttributes<HTMLInputElement> {
  label?: string;
  error?: string;
  fullWidth?: boolean;
}

const Input = forwardRef<HTMLInputElement, InputProps>(
  ({ label, error, className = '', fullWidth = true, ...props }, ref) => {
    const widthClass = fullWidth ? 'w-full' : '';
    const errorClass = error ? 'border-secondary' : 'border-neutral-300';
    
    return (
      <div className={`mb-4 ${widthClass}`}>
        {label && (
          <label className="block text-sm font-medium text-neutral-700 mb-1">
            {label}
          </label>
        )}
        
        <input
          ref={ref}
          className={`
            px-3 py-2 bg-white border ${errorClass} rounded focus:outline-none
            focus:ring-1 focus:ring-primary focus:border-primary
            ${className}
          `}
          {...props}
        />
        
        {error && (
          <p className="mt-1 text-sm text-secondary">{error}</p>
        )}
      </div>
    );
  }
);

Input.displayName = 'Input';

export default Input;