import React, { ButtonHTMLAttributes, ReactNode } from 'react';

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  variant?: 'primary' | 'secondary' | 'accent' | 'highlight' | 'outline' | 'ghost' | 'danger';
  size?: 'sm' | 'md' | 'lg';
  fullWidth?: boolean;
  isLoading?: boolean;
  icon?: ReactNode;
}

const Button: React.FC<ButtonProps> = ({
  children,
  className = '',
  variant = 'primary',
  size = 'md',
  fullWidth = false,
  isLoading = false,
  disabled,
  icon,
  ...props
}) => {
  const baseClasses = 'font-medium rounded focus:outline-none focus:ring-2 transition-all duration-200 ease-in-out';
  
  const variantClasses = {
    primary: 'bg-primary hover:bg-primary-light text-white focus:ring-primary focus:ring-opacity-50',
    secondary: 'bg-secondary hover:bg-secondary-light text-white focus:ring-secondary focus:ring-opacity-50',
    accent: 'bg-accent hover:bg-accent-light text-white focus:ring-accent focus:ring-opacity-50',
    highlight: 'bg-highlight hover:bg-highlight-light text-white focus:ring-highlight focus:ring-opacity-50',
    outline: 'border border-primary text-primary hover:bg-primary hover:text-white focus:ring-primary focus:ring-opacity-50',
    ghost: 'text-primary hover:bg-primary hover:bg-opacity-10 focus:ring-primary focus:ring-opacity-50',
    danger: 'bg-accent hover:bg-accent-dark text-white focus:ring-accent focus:ring-opacity-50',
  };
  
  const sizeClasses = {
    sm: 'py-1 px-3 text-sm',
    md: 'py-2 px-4',
    lg: 'py-3 px-6 text-lg',
  };
  
  const widthClass = fullWidth ? 'w-full' : '';
  const disabledClass = disabled || isLoading ? 'opacity-60 cursor-not-allowed' : '';
  
  return (
    <button
      className={`
        ${baseClasses}
        ${variantClasses[variant]}
        ${sizeClasses[size]}
        ${widthClass}
        ${disabledClass}
        ${className}
      `}
      disabled={disabled || isLoading}
      {...props}
    >
      {isLoading ? (
        <div className="flex items-center justify-center">
          <svg className="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
            <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          Loading...
        </div>
      ) : (
        <div className="flex items-center justify-center">
          {icon && <span className="mr-2">{icon}</span>}
          {children}
        </div>
      )}
    </button>
  );
};

export default Button;