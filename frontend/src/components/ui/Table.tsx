import React, { ReactNode } from 'react';

interface TableProps {
  headers: string[];
  children: ReactNode;
  className?: string;
}

export const Table: React.FC<TableProps> = ({ headers, children, className = '' }) => {
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
          {children}
        </tbody>
      </table>
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