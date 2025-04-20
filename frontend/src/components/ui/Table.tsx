import React, { ReactNode } from 'react';
import { SkeletonTable } from './Skeleton';

interface TableProps {
  headers: string[];
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
}

export const Table: React.FC<TableProps> = ({ 
  headers, 
  children, 
  className = '',
  pagination,
  isLoading = false,
  loadingRows = 5,
  emptyState,
  isDataEmpty = false
}) => {
  if (isLoading) {
    return <SkeletonTable rows={loadingRows} columns={headers.length} className={className} />;
  }
  
  return (
    <div className={`overflow-x-auto ${className}`}>
      <table className="min-w-full bg-white border border-neutral-200 rounded-md overflow-hidden">
        <thead className="bg-primary text-white">
          <tr>
            {headers.map((header, index) => (
              <th
                key={index}
                className="py-3 px-4 text-left font-medium text-sm"
              >
                {header}
              </th>
            ))}
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