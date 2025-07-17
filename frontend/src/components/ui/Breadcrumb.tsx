import React from 'react';
import { Link, useLocation } from 'react-router-dom';

interface BreadcrumbItem {
  label: string;
  path: string;
  isActive?: boolean;
}

interface BreadcrumbProps {
  items?: BreadcrumbItem[];
  className?: string;
}

const Breadcrumb: React.FC<BreadcrumbProps> = ({ 
  items, 
  className = '' 
}) => {
  const location = useLocation();
  
  // Auto-generate breadcrumbs from route if items not provided
  const generateBreadcrumbs = (): BreadcrumbItem[] => {
    const pathSegments = location.pathname.split('/').filter(Boolean);
    const breadcrumbs: BreadcrumbItem[] = [
      { label: 'Home', path: '/' }
    ];
    
    let currentPath = '';
    
    pathSegments.forEach((segment, index) => {
      currentPath += `/${segment}`;
      const isLast = index === pathSegments.length - 1;
      
      // Convert segment to readable label
      let label = segment.charAt(0).toUpperCase() + segment.slice(1);
      label = label.replace(/-/g, ' ').replace(/_/g, ' ');
      
      // Handle special cases
      if (segment === 'admin') label = 'Administration';
      if (segment === 'experts') label = 'Expert Management';
      if (segment === 'requests') label = 'Requests';
      if (segment === 'phases') label = 'Phase Planning';
      if (segment === 'stats') label = 'Statistics';
      if (segment === 'areas') label = 'Area Management';
      if (segment === 'data') label = 'Data Management';
      if (segment === 'engagements') label = 'Engagements';
      if (segment === 'manage') label = 'Management';
      if (segment === 'search') label = 'Expert Search';
      
      breadcrumbs.push({
        label,
        path: currentPath,
        isActive: isLast
      });
    });
    
    return breadcrumbs;
  };
  
  const breadcrumbItems = items || generateBreadcrumbs();
  
  // Don't show breadcrumbs on home page or if only one item
  if (breadcrumbItems.length <= 1 || location.pathname === '/') {
    return null;
  }
  
  return (
    <nav 
      aria-label="Breadcrumb" 
      className={`bg-neutral-100 py-3 border-b border-neutral-200 ${className}`}
    >
      <div className="container">
        <ol className="flex items-center space-x-2 text-sm">
          {breadcrumbItems.map((item, index) => (
            <li key={item.path} className="flex items-center">
              {index > 0 && (
                <svg 
                  className="w-4 h-4 text-neutral-400 mr-2" 
                  fill="currentColor" 
                  viewBox="0 0 20 20"
                  aria-hidden="true"
                >
                  <path 
                    fillRule="evenodd" 
                    d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" 
                    clipRule="evenodd" 
                  />
                </svg>
              )}
              
              {item.isActive || index === breadcrumbItems.length - 1 ? (
                <span 
                  className="text-secondary font-medium"
                  aria-current="page"
                >
                  {item.label}
                </span>
              ) : (
                <Link
                  to={item.path}
                  className="text-neutral-600 hover:text-primary transition-colors duration-200 hover:underline focus:outline-none focus:ring-2 focus:ring-primary focus:ring-opacity-50 rounded px-1"
                >
                  {item.label}
                </Link>
              )}
            </li>
          ))}
        </ol>
      </div>
    </nav>
  );
};

export default Breadcrumb;