import React from 'react';

interface LoadingSpinnerProps {
  size?: 'sm' | 'md' | 'lg' | 'xl';
  color?: 'primary' | 'secondary' | 'success' | 'danger' | 'warning' | 'info' | 'light' | 'dark';
  className?: string;
  label?: string;
  fullScreen?: boolean;
  fullPage?: boolean;
}

export const LoadingSpinner: React.FC<LoadingSpinnerProps> = ({
  size = 'md',
  color = 'primary',
  className = '',
  label,
  fullScreen = false,
  fullPage = false,
}) => {
  // Size mapping
  const sizeMap = {
    sm: 'h-4 w-4',
    md: 'h-8 w-8',
    lg: 'h-12 w-12',
    xl: 'h-16 w-16',
  };

  // Color mapping
  const colorMap = {
    primary: 'text-blue-600',
    secondary: 'text-gray-600',
    success: 'text-green-600',
    danger: 'text-red-600',
    warning: 'text-yellow-600',
    info: 'text-cyan-600',
    light: 'text-gray-300',
    dark: 'text-gray-800',
  };

  const spinnerSize = sizeMap[size];
  const spinnerColor = colorMap[color];

  const spinner = (
    <svg
      className={`animate-spin ${spinnerSize} ${spinnerColor} ${className}`}
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
      data-testid="loading-spinner"
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
  );

  if (fullScreen) {
    return (
      <div className="fixed inset-0 flex items-center justify-center z-50 bg-black bg-opacity-50">
        <div className="flex flex-col items-center p-6 bg-white dark:bg-gray-800 rounded-lg shadow-lg">
          {spinner}
          {label && <p className="mt-4 text-gray-700 dark:text-gray-300">{label}</p>}
        </div>
      </div>
    );
  }

  if (fullPage) {
    return (
      <div className="flex flex-col items-center justify-center w-full h-full min-h-[200px]">
        {spinner}
        {label && <p className="mt-4 text-gray-700 dark:text-gray-300">{label}</p>}
      </div>
    );
  }

  return (
    <div className="inline-flex items-center">
      {spinner}
      {label && <span className="ml-2 text-gray-700 dark:text-gray-300">{label}</span>}
    </div>
  );
};

export const LoadingOverlay: React.FC<{
  isLoading: boolean;
  children: React.ReactNode;
  spinner?: React.ReactNode;
  label?: string;
  className?: string;
}> = ({ isLoading, children, spinner, label, className = '' }) => {
  if (!isLoading) return <>{children}</>;

  return (
    <div className={`relative ${className}`}>
      <div className="absolute inset-0 flex items-center justify-center bg-white/70 dark:bg-gray-800/70 z-10 rounded-lg">
        {spinner || <LoadingSpinner label={label} />}
      </div>
      <div className="opacity-50 pointer-events-none">{children}</div>
    </div>
  );
};