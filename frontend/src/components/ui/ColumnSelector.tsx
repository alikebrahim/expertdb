import React, { useState } from 'react';
import Button from './Button';
import { Checkbox } from './Checkbox';

export interface ColumnConfig {
  key: string;
  label: string;
  width?: string;
  sortable: boolean;
  required: boolean;
  visible: boolean;
  type: 'text' | 'number' | 'date' | 'boolean' | 'rating' | 'status' | 'file';
}

export const DEFAULT_COLUMNS: ColumnConfig[] = [
  // Default column positions: 1. id 2. name 3. affiliation 4. general area 5. specialized area 6. rating 7. nationality 8. available 9. trained 10. phone 11. email 12. cv 13. published
  
  // Required columns (always visible)
  { key: 'id', label: 'Expert ID', sortable: true, required: true, visible: true, type: 'text' },
  { key: 'name', label: 'Name', sortable: true, required: true, visible: true, type: 'text' },
  { key: 'affiliation', label: 'Affiliation', sortable: true, required: true, visible: true, type: 'text' },
  { key: 'generalArea', label: 'General Area', sortable: true, required: true, visible: true, type: 'text' },
  { key: 'specializedArea', label: 'Specialized Area', sortable: true, required: true, visible: true, type: 'text' },
  { key: 'rating', label: 'Rating', sortable: true, required: true, visible: true, type: 'rating' },
  
  // Optional columns (user selectable) - in default position order
  { key: 'nationality', label: 'Nationality', sortable: true, required: false, visible: false, type: 'text' },
  { key: 'isAvailable', label: 'Available', sortable: true, required: false, visible: false, type: 'status' },
  { key: 'isTrained', label: 'Trained', sortable: true, required: false, visible: false, type: 'status' },
  { key: 'phone', label: 'Phone', sortable: true, required: false, visible: false, type: 'text' },
  { key: 'email', label: 'Email', sortable: true, required: false, visible: false, type: 'text' },
  { key: 'cvPath', label: 'CV', sortable: false, required: false, visible: false, type: 'file' },
  { key: 'isPublished', label: 'Published', sortable: true, required: false, visible: false, type: 'status' },
];

interface ColumnSelectorProps {
  columns: ColumnConfig[];
  onColumnChange: (columns: ColumnConfig[]) => void;
}

export const ColumnSelector: React.FC<ColumnSelectorProps> = ({ columns, onColumnChange }) => {
  const [isOpen, setIsOpen] = useState(false);
  
  const toggleColumn = (columnKey: string) => {
    const updated = columns.map(col => 
      col.key === columnKey && !col.required
        ? { ...col, visible: !col.visible }
        : col
    );
    onColumnChange(updated);
  };
  
  const resetToDefault = () => {
    onColumnChange(DEFAULT_COLUMNS);
  };
  
  return (
    <div className="relative">
      <Button
        variant="outline"
        size="sm"
        onClick={() => setIsOpen(!isOpen)}
        className="flex items-center gap-2"
      >
        <svg className="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4" />
        </svg>
        Column Settings
      </Button>
      
      {isOpen && (
        <div className="absolute left-0 mt-2 w-80 bg-white rounded-md shadow-lg border border-neutral-200 z-50">
          <div className="p-4 space-y-4">
            <div className="flex justify-between items-center">
              <h4 className="font-medium text-neutral-900">Table Columns</h4>
              <Button variant="ghost" size="sm" onClick={resetToDefault}>
                Reset
              </Button>
            </div>
            
            <div className="space-y-3 max-h-64 overflow-y-auto">
              {columns.map(column => (
                <div key={column.key} className="flex items-center space-x-3">
                  <Checkbox
                    id={column.key}
                    checked={column.visible}
                    disabled={column.required}
                    onChange={() => toggleColumn(column.key)}
                  />
                  <label 
                    htmlFor={column.key}
                    className={`text-sm cursor-pointer ${
                      column.required ? 'text-neutral-500' : 'text-neutral-900'
                    }`}
                  >
                    {column.label}
                    {column.required && <span className="ml-1 text-xs text-neutral-400">(Required)</span>}
                  </label>
                </div>
              ))}
            </div>
            
            <div className="pt-3 border-t border-neutral-200">
              <p className="text-xs text-neutral-500">
                {columns.filter(c => c.visible).length} of {columns.length} columns visible
              </p>
            </div>
            
            <div className="flex justify-end">
              <Button
                variant="outline"
                size="sm"
                onClick={() => setIsOpen(false)}
              >
                Close
              </Button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
};