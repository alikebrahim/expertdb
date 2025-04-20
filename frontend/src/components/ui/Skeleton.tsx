import React from 'react';

interface SkeletonProps {
  height?: string | number;
  width?: string | number;
  className?: string;
  variant?: 'text' | 'circular' | 'rectangular';
  animation?: 'pulse' | 'wave' | 'none';
}

export const Skeleton: React.FC<SkeletonProps> = ({
  height = '1rem',
  width = '100%',
  className = '',
  variant = 'rectangular',
  animation = 'pulse',
}) => {
  // Convert number to pixel value if needed
  const h = typeof height === 'number' ? `${height}px` : height;
  const w = typeof width === 'number' ? `${width}px` : width;
  
  // Define base classes
  let baseClasses = 'bg-gray-200 dark:bg-gray-700 ';
  
  // Add animation classes
  if (animation === 'pulse') {
    baseClasses += 'animate-pulse ';
  } else if (animation === 'wave') {
    baseClasses += 'relative overflow-hidden before:absolute before:inset-0 before:-translate-x-full before:animate-[wave_1.5s_infinite] before:bg-gradient-to-r before:from-transparent before:via-white/20 before:to-transparent ';
  }
  
  // Add variant classes
  if (variant === 'circular') {
    baseClasses += 'rounded-full ';
  } else if (variant === 'rectangular') {
    baseClasses += 'rounded ';
  } else if (variant === 'text') {
    baseClasses += 'rounded w-full h-4 ';
  }
  
  return (
    <div 
      className={`${baseClasses} ${className}`}
      style={{ height: h, width: w }}
      aria-hidden="true"
    />
  );
};

export const SkeletonText: React.FC<{ lines?: number; className?: string }> = ({ 
  lines = 3,
  className = ''
}) => {
  return (
    <div className={`space-y-2 ${className}`}>
      {Array.from({ length: lines }).map((_, i) => (
        <Skeleton 
          key={i} 
          variant="text" 
          width={i === lines - 1 && lines > 1 ? '80%' : '100%'} 
        />
      ))}
    </div>
  );
};

export const SkeletonCard: React.FC<{ className?: string }> = ({ className = '' }) => {
  return (
    <div className={`p-4 border rounded-lg shadow ${className}`}>
      <Skeleton height={150} className="mb-4" />
      <SkeletonText lines={3} />
    </div>
  );
};

export const SkeletonTable: React.FC<{ 
  rows?: number; 
  columns?: number;
  className?: string;
}> = ({ 
  rows = 5, 
  columns = 4,
  className = ''
}) => {
  return (
    <div className={`overflow-hidden rounded-lg border ${className}`}>
      <div className="grid grid-cols-1 divide-y divide-gray-200 dark:divide-gray-700">
        {/* Header */}
        <div className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4 p-4 bg-gray-50 dark:bg-gray-800">
          {Array.from({ length: columns }).map((_, i) => (
            <Skeleton key={`header-${i}`} height={24} />
          ))}
        </div>
        
        {/* Rows */}
        {Array.from({ length: rows }).map((_, rowIndex) => (
          <div 
            key={`row-${rowIndex}`} 
            className="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-4 lg:grid-cols-6 gap-4 p-4"
          >
            {Array.from({ length: columns }).map((_, colIndex) => (
              <Skeleton 
                key={`cell-${rowIndex}-${colIndex}`} 
                height={20} 
                width="100%" 
              />
            ))}
          </div>
        ))}
      </div>
    </div>
  );
};

export const SkeletonAvatar: React.FC<{ 
  size?: 'sm' | 'md' | 'lg';
  className?: string;
}> = ({ 
  size = 'md',
  className = ''
}) => {
  const sizeMap = {
    sm: 8, // 2rem
    md: 12, // 3rem
    lg: 16, // 4rem
  };
  
  const sizePx = sizeMap[size] * 4; // Convert to pixels (tailwind's rem * 4)
  
  return (
    <Skeleton 
      variant="circular" 
      height={sizePx} 
      width={sizePx} 
      className={className}
    />
  );
};