import React from 'react';

export interface TableProps {
  children: React.ReactNode;
  className?: string;
}

export interface TableHeaderProps {
  children: React.ReactNode;
  className?: string;
}

export interface TableBodyProps {
  children: React.ReactNode;
  className?: string;
}

export interface TableRowProps {
  children: React.ReactNode;
  className?: string;
  onClick?: () => void;
}

export interface TableCellProps {
  children: React.ReactNode;
  className?: string;
  onClick?: (e: React.MouseEvent) => void;
}

export interface TableHeaderCellProps {
  children: React.ReactNode;
  className?: string;
  onClick?: () => void;
}

export const TableHeader = ({ children, className = '' }: TableHeaderProps) => (
  <thead className={`bg-gray-50 ${className}`}>
    {children}
  </thead>
);

export const TableBody = ({ children, className = '' }: TableBodyProps) => (
  <tbody className={`bg-white divide-y divide-gray-200 ${className}`}>
    {children}
  </tbody>
);

export const TableRow = ({ children, className = '', onClick }: TableRowProps) => (
  <tr className={className} onClick={onClick}>
    {children}
  </tr>
);

export const TableCell = ({ children, className = '', onClick }: TableCellProps) => (
  <td className={`px-6 py-4 whitespace-nowrap text-sm text-gray-900 ${className}`} onClick={onClick}>
    {children}
  </td>
);

export const TableHeaderCell = ({ children, className = '', onClick }: TableHeaderCellProps) => (
  <th 
    className={`px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider ${className}`} 
    onClick={onClick}
  >
    {children}
  </th>
);

// Define the compound component interface
interface TableComponent extends React.FC<TableProps> {
  Header: React.FC<TableHeaderProps>;
  Body: React.FC<TableBodyProps>;
  Row: React.FC<TableRowProps>;
  Cell: React.FC<TableCellProps>;
  HeaderCell: React.FC<TableHeaderCellProps>;
}

const Table = ({ children, className = '' }: TableProps) => (
  <table className={`divide-y divide-gray-200 ${className}`}>
    {children}
  </table>
);

// Attach sub-components to the main Table component
(Table as TableComponent).Header = TableHeader;
(Table as TableComponent).Body = TableBody;
(Table as TableComponent).Row = TableRow;
(Table as TableComponent).Cell = TableCell;
(Table as TableComponent).HeaderCell = TableHeaderCell;

export default Table as TableComponent;