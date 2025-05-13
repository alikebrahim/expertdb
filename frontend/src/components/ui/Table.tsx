import React, { ReactNode } from 'react';
import { SkeletonTable } from './Skeleton';
import { SortConfig } from '../tables/ExpertTable';

export interface TableHeader {
  label: string;
  field?: string;
  sortable?: boolean;
}

interface TableProps {
  headers: (string | TableHeader)[];
  children: ReactNode;
  className?: string;
  pagination?: {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
  };
  isLoading?: boolean;
  loadingRows?: number;
  emptyState?: ReactNode;
  isDataEmpty?: boolean;
  sortConfig?: SortConfig;
  onSort?: (field: string) => void;
}

export const Table: React.FC<TableProps> = ({ 
  headers, 
  children, 
  className = '',
  pagination,
  isLoading = false,
  loadingRows = 5,
  emptyState,
  isDataEmpty = false,
  sortConfig,
  onSort
}) => {
  if (isLoading) {
    return <SkeletonTable rows={loadingRows} columns={headers.length} className={className} />;
  }
  
  const renderSortIcon = (field: string) => {
    if (!sortConfig || sortConfig.field !== field) {
      // Neutral icon for unsorted column
      return (
        <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-white opacity-50" viewBox="0 0 20 20" fill="currentColor">
          <path fillRule="evenodd" d="M10 3a1 1 0 01.707.293l3 3a1 1 0 01-1.414 1.414L10 5.414 7.707 7.707a1 1 0 01-1.414-1.414l3-3A1 1 0 0110 3zm-3.707 9.293a1 1 0 011.414 0L10 14.586l2.293-2.293a1 1 0 011.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clipRule="evenodd" />
        </svg>
      );
    } else if (sortConfig.direction === 'asc') {
      // Up arrow for ascending sort
      return (
        <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-white" viewBox="0 0 20 20" fill="currentColor">
          <path fillRule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z" clipRule="evenodd" />
        </svg>
      );
    } else {
      // Down arrow for descending sort
      return (
        <svg xmlns="http://www.w3.org/2000/svg" className="h-4 w-4 text-white" viewBox="0 0 20 20" fill="currentColor">
          <path fillRule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clipRule="evenodd" />
        </svg>
      );
    }
  };
  
  return (
    <div className={`overflow-x-auto ${className}`}>
      <table className="min-w-full bg-white border border-neutral-200 rounded-md overflow-hidden">
        <thead className="bg-primary text-white">
          <tr>
            {headers.map((header, index) => {
              // Handle string headers (legacy support)
              if (typeof header === 'string') {
                return (
                  <th
                    key={index}
                    className="py-3 px-4 text-left font-medium text-sm"
                  >
                    {header}
                  </th>
                );
              }
              
              // Handle object headers with sortable functionality
              const { label, field, sortable } = header;
              
              return (
                <th
                  key={index}
                  className={`py-3 px-4 text-left font-medium text-sm ${
                    sortable && onSort ? 'cursor-pointer hover:bg-primary-dark' : ''
                  }`}
                  onClick={() => sortable && field && onSort && onSort(field)}
                >
                  <div className="flex items-center space-x-1">
                    <span>{label}</span>
                    {sortable && field && onSort && (
                      <span className="inline-block ml-1">
                        {renderSortIcon(field)}
                      </span>
                    )}
                  </div>
                </th>
              );
            })}
          </tr>
        </thead>
        <tbody className="divide-y divide-neutral-200">
          {isDataEmpty ? (
            <tr>
              <td colSpan={headers.length} className="py-8 px-4 text-center text-gray-500">
                {emptyState || "No data available"}
              </td>
            </tr>
          ) : (
            children
          )}
        </tbody>
      </table>
      
      {pagination && (
        <div className="flex items-center justify-between border-t border-neutral-200 bg-white px-4 py-3 mt-2">
          <div className="flex flex-1 justify-between sm:hidden">
            <button
              onClick={() => pagination.onPageChange(pagination.currentPage - 1)}
              disabled={pagination.currentPage === 1}
              className={`relative inline-flex items-center rounded-md px-4 py-2 text-sm font-medium ${
                pagination.currentPage === 1
                  ? 'text-neutral-400 cursor-not-allowed'
                  : 'text-primary hover:bg-primary hover:text-white'
              }`}
            >
              Previous
            </button>
            <button
              onClick={() => pagination.onPageChange(pagination.currentPage + 1)}
              disabled={pagination.currentPage === pagination.totalPages}
              className={`relative ml-3 inline-flex items-center rounded-md px-4 py-2 text-sm font-medium ${
                pagination.currentPage === pagination.totalPages
                  ? 'text-neutral-400 cursor-not-allowed'
                  : 'text-primary hover:bg-primary hover:text-white'
              }`}
            >
              Next
            </button>
          </div>
          <div className="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
            <div>
              <p className="text-sm text-neutral-700">
                Showing page <span className="font-medium">{pagination.currentPage}</span> of{' '}
                <span className="font-medium">{pagination.totalPages}</span>
              </p>
            </div>
            <div>
              <nav className="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
                <button
                  onClick={() => pagination.onPageChange(1)}
                  disabled={pagination.currentPage === 1}
                  className={`relative inline-flex items-center rounded-l-md px-2 py-2 text-neutral-400 ${
                    pagination.currentPage === 1
                      ? 'cursor-not-allowed'
                      : 'hover:bg-primary hover:text-white'
                  }`}
                >
                  <span className="sr-only">First</span>
                  <svg className="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fillRule="evenodd" d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" clipRule="evenodd" />
                    <path fillRule="evenodd" d="M6.79 5.23a.75.75 0 01-.02 1.06L2.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" clipRule="evenodd" />
                  </svg>
                </button>
                <button
                  onClick={() => pagination.onPageChange(pagination.currentPage - 1)}
                  disabled={pagination.currentPage === 1}
                  className={`relative inline-flex items-center px-2 py-2 text-neutral-400 ${
                    pagination.currentPage === 1
                      ? 'cursor-not-allowed'
                      : 'hover:bg-primary hover:text-white'
                  }`}
                >
                  <span className="sr-only">Previous</span>
                  <svg className="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fillRule="evenodd" d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" clipRule="evenodd" />
                  </svg>
                </button>
                
                {Array.from({ length: pagination.totalPages }, (_, i) => i + 1)
                  .filter(page => (
                    page === 1 || 
                    page === pagination.totalPages || 
                    (page >= pagination.currentPage - 1 && page <= pagination.currentPage + 1)
                  ))
                  .map((page, index, array) => {
                    const showEllipsis = index > 0 && array[index - 1] !== page - 1;
                    
                    return (
                      <React.Fragment key={page}>
                        {showEllipsis && (
                          <span className="relative inline-flex items-center px-4 py-2 text-sm font-semibold text-neutral-700 cursor-default">
                            ...
                          </span>
                        )}
                        <button
                          onClick={() => pagination.onPageChange(page)}
                          className={`relative inline-flex items-center px-4 py-2 text-sm font-semibold ${
                            page === pagination.currentPage
                              ? 'bg-primary text-white focus:z-20'
                              : 'text-neutral-900 hover:bg-primary hover:text-white'
                          }`}
                        >
                          {page}
                        </button>
                      </React.Fragment>
                    );
                  })}
                
                <button
                  onClick={() => pagination.onPageChange(pagination.currentPage + 1)}
                  disabled={pagination.currentPage === pagination.totalPages}
                  className={`relative inline-flex items-center px-2 py-2 text-neutral-400 ${
                    pagination.currentPage === pagination.totalPages
                      ? 'cursor-not-allowed'
                      : 'hover:bg-primary hover:text-white'
                  }`}
                >
                  <span className="sr-only">Next</span>
                  <svg className="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fillRule="evenodd" d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clipRule="evenodd" />
                  </svg>
                </button>
                <button
                  onClick={() => pagination.onPageChange(pagination.totalPages)}
                  disabled={pagination.currentPage === pagination.totalPages}
                  className={`relative inline-flex items-center rounded-r-md px-2 py-2 text-neutral-400 ${
                    pagination.currentPage === pagination.totalPages
                      ? 'cursor-not-allowed'
                      : 'hover:bg-primary hover:text-white'
                  }`}
                >
                  <span className="sr-only">Last</span>
                  <svg className="h-5 w-5" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fillRule="evenodd" d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clipRule="evenodd" />
                    <path fillRule="evenodd" d="M13.21 14.77a.75.75 0 01.02-1.06L17.168 10 13.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clipRule="evenodd" />
                  </svg>
                </button>
              </nav>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};

interface TableRowProps {
  children: ReactNode;
  onClick?: () => void;
  isClickable?: boolean;
}

export const TableRow: React.FC<TableRowProps> = ({ 
  children, 
  onClick, 
  isClickable = false 
}) => {
  const clickableClass = isClickable 
    ? 'cursor-pointer hover:bg-neutral-50' 
    : '';
  
  return (
    <tr 
      className={clickableClass}
      onClick={isClickable ? onClick : undefined}
    >
      {children}
    </tr>
  );
};

interface TableCellProps {
  children: ReactNode;
  className?: string;
}

export const TableCell: React.FC<TableCellProps> = ({ children, className = '' }) => {
  return (
    <td className={`py-3 px-4 text-sm text-neutral-800 ${className}`}>
      {children}
    </td>
  );
};

export default { Table, TableRow, TableCell };