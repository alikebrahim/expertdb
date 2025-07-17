# ExpertDB Frontend Phase 2 Implementation Plan
# Expert Database Browsing Interface

## Overview

Phase 2 focuses on implementing the core functionality of the ExpertDB application - the expert database browsing interface. This phase builds upon the foundation established in Phase 1 (BQA Design System & Authentication) and creates the primary user interface that all users will interact with upon login.

**Key Principles:**
- **Unified Access**: All users have equal browsing privileges regardless of role
- **Performance First**: Full database load with efficient client-side operations
- **Customizable Display**: Users can select which columns to display
- **Real-time Operations**: Instant sorting, filtering, and searching without server requests

## Phase 2 Scope & Objectives

### Primary Goals
1. **Main Expert Search Interface**: Comprehensive search and filter system
2. **Expert Results Display**: Customizable table with full database loading
3. **Expert Profile Detail Views**: Complete expert information display

### Success Criteria
- Expert search response time < 1 second for all operations
- Support for 2000+ expert records with smooth performance
- 90%+ search success rate (searches leading to profile views)
- Mobile-responsive interface with consistent UX

## Detailed Implementation Specifications

## 1. Expert Search Interface Enhancement

### 1.1 Search Page Architecture

**File Structure:**
```
frontend/src/pages/SearchPage.tsx (Enhanced)
frontend/src/components/ExpertSearch/
├── SearchHeader.tsx
├── FilterPanel.tsx
├── ResultsSummary.tsx
├── ExpertTable.tsx
├── ColumnSelector.tsx
└── index.ts
```

**State Management:**
```typescript
interface ExpertSearchState {
  // Core data
  experts: Expert[];
  filteredExperts: Expert[];
  loading: boolean;
  error: string | null;
  
  // Filter state
  filters: ExpertFilters;
  activeFilterCount: number;
  
  // Table state
  sortConfig: SortConfig;
  selectedColumns: ColumnConfig[];
  selectedExperts: number[];
  
  // UI state
  showAdvancedFilters: boolean;
  showColumnSelector: boolean;
  currentStep: number;
  
  // Performance
  lastUpdated: Date;
  searchTerm: string;
  searchDebounce: NodeJS.Timeout | null;
}
```

### 1.2 Full Database Loading Strategy

**Implementation Approach:**
```typescript
// Enhanced SearchPage.tsx
const SearchPage = () => {
  const [experts, setExperts] = useState<Expert[]>([]);
  const [filteredExperts, setFilteredExperts] = useState<Expert[]>([]);
  const [loading, setLoading] = useState(true);
  
  // Load all experts on component mount
  useEffect(() => {
    const loadAllExperts = async () => {
      setLoading(true);
      try {
        // Fetch all experts without pagination
        const response = await expertsApi.getAllExperts();
        if (response.success) {
          setExperts(response.data);
          setFilteredExperts(response.data);
          
          // Cache data for performance
          localStorage.setItem('expertsCache', JSON.stringify({
            data: response.data,
            timestamp: Date.now()
          }));
        }
      } catch (error) {
        setError('Failed to load expert database');
      } finally {
        setLoading(false);
      }
    };
    
    // Check cache first
    const cached = getCachedExperts();
    if (cached && isCacheValid(cached.timestamp)) {
      setExperts(cached.data);
      setFilteredExperts(cached.data);
      setLoading(false);
    } else {
      loadAllExperts();
    }
  }, []);
  
  // Client-side filtering
  useEffect(() => {
    const filtered = applyFilters(experts, filters, searchTerm);
    setFilteredExperts(filtered);
  }, [experts, filters, searchTerm]);
};
```

**Cache Strategy:**
```typescript
// utils/expertCache.ts
interface ExpertCache {
  data: Expert[];
  timestamp: number;
}

const CACHE_DURATION = 5 * 60 * 1000; // 5 minutes

export const getCachedExperts = (): ExpertCache | null => {
  try {
    const cached = localStorage.getItem('expertsCache');
    return cached ? JSON.parse(cached) : null;
  } catch {
    return null;
  }
};

export const isCacheValid = (timestamp: number): boolean => {
  return Date.now() - timestamp < CACHE_DURATION;
};
```

### 1.3 Advanced Filter Panel

**Filter Panel Component:**
```typescript
// components/ExpertSearch/FilterPanel.tsx
interface FilterPanelProps {
  filters: ExpertFilters;
  onFilterChange: (filters: ExpertFilters) => void;
  expertAreas: ExpertArea[];
}

const FilterPanel: React.FC<FilterPanelProps> = ({
  filters,
  onFilterChange,
  expertAreas
}) => {
  const [isExpanded, setIsExpanded] = useState(false);
  
  return (
    <div className="bg-white rounded-lg shadow-sm border border-neutral-200 p-6">
      {/* Basic Filters - Always Visible */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
        <SearchInput
          placeholder="Search name or institution..."
          value={filters.searchTerm}
          onChange={(value) => onFilterChange({ ...filters, searchTerm: value })}
          className="md:col-span-1"
        />
        
        <Select
          label="Role"
          value={filters.role}
          onChange={(value) => onFilterChange({ ...filters, role: value })}
          options={ROLE_OPTIONS}
        />
        
        <Select
          label="Employment Type"
          value={filters.employmentType}
          onChange={(value) => onFilterChange({ ...filters, employmentType: value })}
          options={EMPLOYMENT_OPTIONS}
        />
      </div>
      
      {/* Advanced Filters - Expandable */}
      <div className="border-t border-neutral-200 pt-4">
        <button
          onClick={() => setIsExpanded(!isExpanded)}
          className="flex items-center text-sm font-medium text-secondary hover:text-secondary-dark"
        >
          <span>Advanced Filters</span>
          <ChevronDownIcon className={`ml-2 h-4 w-4 transition-transform ${isExpanded ? 'rotate-180' : ''}`} />
        </button>
        
        {isExpanded && (
          <div className="mt-4 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <Select
              label="Expert Area"
              value={filters.expertAreaId}
              onChange={(value) => onFilterChange({ ...filters, expertAreaId: value })}
              options={expertAreas.map(area => ({ value: area.id, label: area.name }))}
            />
            
            <Select
              label="Nationality"
              value={filters.nationality}
              onChange={(value) => onFilterChange({ ...filters, nationality: value })}
              options={NATIONALITY_OPTIONS}
            />
            
            <Select
              label="Minimum Rating"
              value={filters.minRating}
              onChange={(value) => onFilterChange({ ...filters, minRating: value })}
              options={RATING_OPTIONS}
            />
            
            <div className="space-y-2">
              <label className="text-sm font-medium text-neutral-700">Status</label>
              <div className="space-y-1">
                <Checkbox
                  label="Available Only"
                  checked={filters.isAvailable}
                  onChange={(checked) => onFilterChange({ ...filters, isAvailable: checked })}
                />
                <Checkbox
                  label="Bahraini Only"
                  checked={filters.isBahraini}
                  onChange={(checked) => onFilterChange({ ...filters, isBahraini: checked })}
                />
                <Checkbox
                  label="Published Only"
                  checked={filters.isPublished}
                  onChange={(checked) => onFilterChange({ ...filters, isPublished: checked })}
                />
              </div>
            </div>
          </div>
        )}
      </div>
      
      {/* Filter Actions */}
      <div className="flex justify-between items-center mt-4 pt-4 border-t border-neutral-200">
        <FilterBadges filters={filters} onRemoveFilter={onFilterChange} />
        <div className="flex gap-2">
          <Button
            variant="outline"
            size="sm"
            onClick={() => onFilterChange(getDefaultFilters())}
          >
            Clear All
          </Button>
          <Button
            variant="primary"
            size="sm"
            onClick={() => saveFilterPreset(filters)}
          >
            Save Preset
          </Button>
        </div>
      </div>
    </div>
  );
};
```

### 1.4 Client-Side Filter Implementation

**Filter Logic:**
```typescript
// utils/expertFilters.ts
export const applyFilters = (
  experts: Expert[],
  filters: ExpertFilters,
  searchTerm: string
): Expert[] => {
  return experts.filter(expert => {
    // Text search
    if (searchTerm) {
      const searchLower = searchTerm.toLowerCase();
      const matchesSearch = 
        expert.name.toLowerCase().includes(searchLower) ||
        expert.institution.toLowerCase().includes(searchLower) ||
        expert.specializedArea.toLowerCase().includes(searchLower);
      
      if (!matchesSearch) return false;
    }
    
    // Role filter
    if (filters.role && expert.role !== filters.role) {
      return false;
    }
    
    // Employment type filter
    if (filters.employmentType && expert.employmentType !== filters.employmentType) {
      return false;
    }
    
    // Expert area filter
    if (filters.expertAreaId && expert.expertAreaId !== parseInt(filters.expertAreaId)) {
      return false;
    }
    
    // Nationality filter
    if (filters.nationality && expert.nationality !== filters.nationality) {
      return false;
    }
    
    // Rating filter
    if (filters.minRating && expert.rating < parseFloat(filters.minRating)) {
      return false;
    }
    
    // Boolean filters
    if (filters.isAvailable && !expert.isAvailable) return false;
    if (filters.isBahraini && !expert.isBahraini) return false;
    if (filters.isPublished && !expert.isPublished) return false;
    
    return true;
  });
};

// Debounced search for performance
export const useDebounceCallback = <T extends any[]>(
  callback: (...args: T) => void,
  delay: number
) => {
  const timeoutRef = useRef<NodeJS.Timeout>();
  
  return useCallback((...args: T) => {
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }
    
    timeoutRef.current = setTimeout(() => {
      callback(...args);
    }, delay);
  }, [callback, delay]);
};
```

## 2. Customizable Expert Table Display

### 2.1 Column Selection System

**Column Configuration:**
```typescript
// types/tableConfig.ts
interface ColumnConfig {
  key: string;
  label: string;
  width?: string;
  sortable: boolean;
  required: boolean;
  visible: boolean;
  type: 'text' | 'number' | 'date' | 'boolean' | 'rating' | 'status';
}

const DEFAULT_COLUMNS: ColumnConfig[] = [
  // Standard columns (always displayed)
  { key: 'expertId', label: 'Expert ID', sortable: true, required: true, visible: true, type: 'text' },
  { key: 'name', label: 'Name', sortable: true, required: true, visible: true, type: 'text' },
  { key: 'institution', label: 'Institution', sortable: true, required: true, visible: true, type: 'text' },
  { key: 'specializedArea', label: 'Specialized Area', sortable: true, required: true, visible: true, type: 'text' },
  { key: 'rating', label: 'Rating', sortable: true, required: true, visible: true, type: 'rating' },
  { key: 'role', label: 'Role', sortable: true, required: true, visible: true, type: 'text' },
  { key: 'employmentType', label: 'Employment Type', sortable: true, required: true, visible: true, type: 'text' },
  { key: 'generalArea', label: 'General Area', sortable: true, required: true, visible: true, type: 'text' },
  { key: 'actions', label: 'Actions', sortable: false, required: true, visible: true, type: 'text' },
  // Optional columns (user selectable)
  { key: 'isAvailable', label: 'Available', sortable: true, required: false, visible: false, type: 'status' },
  { key: 'nationality', label: 'Nationality', sortable: true, required: false, visible: false, type: 'text' },
  { key: 'isTrained', label: 'Trained', sortable: true, required: false, visible: false, type: 'status' },
  { key: 'dateAdded', label: 'Date Added', sortable: true, required: false, visible: false, type: 'date' }
];
```

**Column Selector Component:**
```typescript
// components/ExpertSearch/ColumnSelector.tsx
const ColumnSelector: React.FC<{
  columns: ColumnConfig[];
  onColumnChange: (columns: ColumnConfig[]) => void;
}> = ({ columns, onColumnChange }) => {
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
    <Popover open={isOpen} onOpenChange={setIsOpen}>
      <PopoverTrigger asChild>
        <Button variant="outline" size="sm">
          <SettingsIcon className="h-4 w-4 mr-2" />
          Column Settings
        </Button>
      </PopoverTrigger>
      
      <PopoverContent className="w-80" align="end">
        <div className="space-y-4">
          <div className="flex justify-between items-center">
            <h4 className="font-medium">Table Columns</h4>
            <Button variant="ghost" size="sm" onClick={resetToDefault}>
              Reset
            </Button>
          </div>
          
          <div className="space-y-2 max-h-64 overflow-y-auto">
            {columns.map(column => (
              <div key={column.key} className="flex items-center space-x-2">
                <Checkbox
                  id={column.key}
                  checked={column.visible}
                  disabled={column.required}
                  onCheckedChange={() => toggleColumn(column.key)}
                />
                <label 
                  htmlFor={column.key}
                  className={`text-sm ${column.required ? 'text-neutral-500' : 'text-neutral-900'}`}
                >
                  {column.label}
                  {column.required && <span className="ml-1 text-xs">(Required)</span>}
                </label>
              </div>
            ))}
          </div>
          
          <div className="text-xs text-neutral-500">
            {columns.filter(c => c.visible).length} of {columns.length} columns visible
          </div>
        </div>
      </PopoverContent>
    </Popover>
  );
};
```

### 2.2 Enhanced Expert Table

**Table Component with Performance Optimization:**
```typescript
// components/ExpertSearch/ExpertTable.tsx
const ExpertTable: React.FC<{
  experts: Expert[];
  columns: ColumnConfig[];
  sortConfig: SortConfig;
  onSort: (field: string) => void;
  selectedExperts: number[];
  onSelectExpert: (expertId: number, selected: boolean) => void;
}> = ({ experts, columns, sortConfig, onSort, selectedExperts, onSelectExpert }) => {
  // Virtual scrolling for large datasets
  const tableRef = useRef<HTMLTableElement>(null);
  const [visibleRange, setVisibleRange] = useState({ start: 0, end: 50 });
  
  // Virtualization for performance with large datasets
  useEffect(() => {
    const updateVisibleRange = () => {
      if (!tableRef.current) return;
      
      const scrollTop = tableRef.current.scrollTop;
      const rowHeight = 60; // Approximate row height
      const containerHeight = tableRef.current.clientHeight;
      
      const start = Math.floor(scrollTop / rowHeight);
      const end = Math.min(start + Math.ceil(containerHeight / rowHeight) + 5, experts.length);
      
      setVisibleRange({ start, end });
    };
    
    if (experts.length > 100) {
      const table = tableRef.current;
      table?.addEventListener('scroll', updateVisibleRange);
      return () => table?.removeEventListener('scroll', updateVisibleRange);
    }
  }, [experts.length]);
  
  const visibleColumns = columns.filter(col => col.visible);
  const visibleExperts = experts.length > 100 
    ? experts.slice(visibleRange.start, visibleRange.end)
    : experts;
  
  return (
    <div className="bg-white rounded-lg shadow-sm border border-neutral-200 overflow-hidden">
      {/* Table Header */}
      <div className="px-6 py-4 border-b border-neutral-200">
        <div className="flex justify-between items-center">
          <h3 className="text-lg font-medium text-neutral-900">
            Expert Database ({experts.length} experts)
          </h3>
          
          <div className="flex items-center gap-2">
            <ExportButton experts={experts} selectedExperts={selectedExperts} />
            <ColumnSelector columns={columns} onColumnChange={onColumnChange} />
          </div>
        </div>
      </div>
      
      {/* Table */}
      <div 
        ref={tableRef}
        className="overflow-auto max-h-[600px]"
        style={{ height: experts.length > 100 ? '600px' : 'auto' }}
      >
        <table className="w-full">
          <thead className="bg-neutral-50 sticky top-0 z-10">
            <tr>
              {visibleColumns.map(column => (
                <th 
                  key={column.key}
                  className="px-4 py-3 text-left text-xs font-medium text-neutral-500 uppercase tracking-wider"
                  style={{ width: column.width }}
                >
                  {column.sortable ? (
                    <button
                      onClick={() => onSort(column.key)}
                      className="flex items-center space-x-1 hover:text-neutral-700"
                    >
                      <span>{column.label}</span>
                      <SortIcon 
                        field={column.key}
                        sortConfig={sortConfig}
                        className="h-4 w-4"
                      />
                    </button>
                  ) : (
                    column.label
                  )}
                </th>
              ))}
            </tr>
          </thead>
          
          <tbody className="bg-white divide-y divide-neutral-200">
            {experts.length > 100 && visibleRange.start > 0 && (
              <tr style={{ height: visibleRange.start * 60 }}>
                <td colSpan={visibleColumns.length} />
              </tr>
            )}
            
            {visibleExperts.map((expert, index) => (
              <ExpertTableRow 
                key={expert.id}
                expert={expert}
                columns={visibleColumns}
                selected={selectedExperts.includes(expert.id)}
                onSelect={(selected) => onSelectExpert(expert.id, selected)}
                style={{ position: 'absolute', top: (visibleRange.start + index) * 60 }}
              />
            ))}
            
            {experts.length > 100 && visibleRange.end < experts.length && (
              <tr style={{ height: (experts.length - visibleRange.end) * 60 }}>
                <td colSpan={visibleColumns.length} />
              </tr>
            )}
          </tbody>
        </table>
      </div>
      
      {experts.length === 0 && (
        <div className="text-center py-12">
          <UsersIcon className="h-12 w-12 text-neutral-400 mx-auto mb-4" />
          <h3 className="text-lg font-medium text-neutral-900 mb-2">No experts found</h3>
          <p className="text-neutral-500">Try adjusting your search criteria</p>
        </div>
      )}
    </div>
  );
};
```

### 2.3 Client-Side Sorting Implementation

**Sorting Utilities:**
```typescript
// utils/tableSorting.ts
export const sortExperts = (
  experts: Expert[],
  sortConfig: SortConfig
): Expert[] => {
  if (!sortConfig.field) return experts;
  
  return [...experts].sort((a, b) => {
    const aValue = getNestedValue(a, sortConfig.field);
    const bValue = getNestedValue(b, sortConfig.field);
    
    // Handle different data types
    if (typeof aValue === 'string' && typeof bValue === 'string') {
      const comparison = aValue.localeCompare(bValue);
      return sortConfig.direction === 'asc' ? comparison : -comparison;
    }
    
    if (typeof aValue === 'number' && typeof bValue === 'number') {
      const comparison = aValue - bValue;
      return sortConfig.direction === 'asc' ? comparison : -comparison;
    }
    
    if (aValue instanceof Date && bValue instanceof Date) {
      const comparison = aValue.getTime() - bValue.getTime();
      return sortConfig.direction === 'asc' ? comparison : -comparison;
    }
    
    // Boolean values
    if (typeof aValue === 'boolean' && typeof bValue === 'boolean') {
      const comparison = aValue === bValue ? 0 : aValue ? 1 : -1;
      return sortConfig.direction === 'asc' ? comparison : -comparison;
    }
    
    return 0;
  });
};

const getNestedValue = (obj: any, path: string): any => {
  return path.split('.').reduce((current, key) => current?.[key], obj);
};
```

## 3. Expert Profile Detail Enhancement

### 3.1 Enhanced Profile Page

**Profile Page Structure:**
```typescript
// pages/ExpertDetailPage.tsx (Enhanced)
const ExpertDetailPage = () => {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const [expert, setExpert] = useState<Expert | null>(null);
  const [loading, setLoading] = useState(true);
  const [activeTab, setActiveTab] = useState('personal');
  
  // Load expert data
  useEffect(() => {
    const loadExpert = async () => {
      if (!id) return;
      
      try {
        const response = await expertsApi.getExpert(parseInt(id));
        if (response.success) {
          setExpert(response.data);
        }
      } catch (error) {
        // Handle error
      } finally {
        setLoading(false);
      }
    };
    
    loadExpert();
  }, [id]);
  
  if (loading) return <LoadingSpinner />;
  if (!expert) return <NotFound />;
  
  return (
    <Layout>
      <div className="max-w-4xl mx-auto">
        {/* Breadcrumb Navigation */}
        <nav className="mb-6">
          <button
            onClick={() => navigate('/search')}
            className="flex items-center text-sm text-secondary hover:text-secondary-dark"
          >
            <ArrowLeftIcon className="h-4 w-4 mr-2" />
            Back to Expert Search
          </button>
        </nav>
        
        {/* Profile Header */}
        <ExpertProfileHeader expert={expert} />
        
        {/* Tabbed Content */}
        <div className="mt-6">
          <TabNavigation 
            activeTab={activeTab}
            onTabChange={setActiveTab}
            tabs={PROFILE_TABS}
          />
          
          <div className="mt-6">
            {activeTab === 'personal' && <PersonalInfoTab expert={expert} />}
            {activeTab === 'expertise' && <ExpertiseTab expert={expert} />}
            {activeTab === 'biography' && <BiographyTab expert={expert} />}
            {activeTab === 'documents' && <DocumentsTab expert={expert} />}
            {activeTab === 'engagement' && <EngagementTab expert={expert} />}
          </div>
        </div>
      </div>
    </Layout>
  );
};
```

### 3.2 Profile Header Component

**Enhanced Profile Header:**
```typescript
// components/ExpertProfile/ProfileHeader.tsx
const ExpertProfileHeader: React.FC<{ expert: Expert }> = ({ expert }) => {
  return (
    <div className="bg-white rounded-lg shadow-sm border border-neutral-200 p-6">
      <div className="flex flex-col lg:flex-row lg:items-start lg:space-x-6">
        {/* Profile Image */}
        <div className="flex-shrink-0 mb-4 lg:mb-0">
          <div className="h-24 w-24 bg-neutral-200 rounded-lg flex items-center justify-center">
            <UserIcon className="h-12 w-12 text-neutral-400" />
          </div>
        </div>
        
        {/* Basic Info */}
        <div className="flex-grow">
          <div className="flex flex-col lg:flex-row lg:justify-between lg:items-start">
            <div>
              <h1 className="text-2xl font-bold text-neutral-900">{expert.name}</h1>
              <p className="text-lg text-neutral-600 mt-1">{expert.designation}</p>
              <p className="text-neutral-500 mt-1">{expert.institution}</p>
              
              <div className="flex items-center mt-3 space-x-4">
                <RatingDisplay rating={expert.rating} showValue />
                <AvailabilityBadge isAvailable={expert.isAvailable} />
                {expert.isBahraini && <BahrainiFlag />}
              </div>
            </div>
            
            {/* Actions */}
            <div className="flex space-x-2 mt-4 lg:mt-0">
              <Button
                variant="primary"
                onClick={() => window.open(`mailto:${expert.email}`)}
              >
                <MailIcon className="h-4 w-4 mr-2" />
                Contact
              </Button>
              
              <Button
                variant="outline"
                onClick={() => window.open(`tel:${expert.phone}`)}
              >
                <PhoneIcon className="h-4 w-4 mr-2" />
                Call
              </Button>
              
              <Button
                variant="outline"
                onClick={() => downloadCV(expert.id)}
              >
                <DocumentIcon className="h-4 w-4 mr-2" />
                Download CV
              </Button>
            </div>
          </div>
        </div>
      </div>
      
      {/* Expert ID */}
      <div className="mt-4 pt-4 border-t border-neutral-200">
        <span className="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-secondary bg-opacity-10 text-secondary">
          Expert ID: {expert.expertId}
        </span>
      </div>
    </div>
  );
};
```

## 4. Performance Optimization

### 4.1 Memory Management

**Efficient Data Handling:**
```typescript
// utils/memoryOptimization.ts
export class ExpertDataManager {
  private data: Expert[] = [];
  private filteredData: Expert[] = [];
  private cache: Map<string, Expert[]> = new Map();
  
  constructor(private maxCacheSize: number = 10) {}
  
  setData(experts: Expert[]) {
    this.data = experts;
    this.filteredData = experts;
    this.clearCache();
  }
  
  filter(filters: ExpertFilters): Expert[] {
    const cacheKey = this.generateCacheKey(filters);
    
    if (this.cache.has(cacheKey)) {
      return this.cache.get(cacheKey)!;
    }
    
    const filtered = this.applyFilters(this.data, filters);
    
    // Cache management
    if (this.cache.size >= this.maxCacheSize) {
      const firstKey = this.cache.keys().next().value;
      this.cache.delete(firstKey);
    }
    
    this.cache.set(cacheKey, filtered);
    return filtered;
  }
  
  private generateCacheKey(filters: ExpertFilters): string {
    return JSON.stringify(filters);
  }
  
  private clearCache() {
    this.cache.clear();
  }
}
```

### 4.2 Virtual Scrolling for Large Datasets

**Virtual Scrolling Implementation:**
```typescript
// hooks/useVirtualScrolling.ts
export const useVirtualScrolling = <T>(
  items: T[],
  containerHeight: number,
  itemHeight: number,
  overscan: number = 5
) => {
  const [scrollTop, setScrollTop] = useState(0);
  
  const visibleStart = Math.floor(scrollTop / itemHeight);
  const visibleEnd = Math.min(
    visibleStart + Math.ceil(containerHeight / itemHeight) + overscan,
    items.length
  );
  
  const totalHeight = items.length * itemHeight;
  const offsetY = visibleStart * itemHeight;
  
  const visibleItems = items.slice(visibleStart, visibleEnd);
  
  return {
    visibleItems,
    totalHeight,
    offsetY,
    onScroll: (e: React.UIEvent<HTMLDivElement>) => {
      setScrollTop(e.currentTarget.scrollTop);
    }
  };
};
```

## 5. Testing Strategy

### 5.1 Unit Tests

**Filter Logic Testing:**
```typescript
// __tests__/expertFilters.test.ts
describe('Expert Filtering', () => {
  const mockExperts: Expert[] = [
    { id: 1, name: 'Dr. Ahmed Ali', institution: 'UoB', rating: 4.5, isAvailable: true },
    { id: 2, name: 'Sarah Hassan', institution: 'MOH', rating: 3.8, isAvailable: false },
  ];
  
  test('filters by name correctly', () => {
    const filters = { searchTerm: 'Ahmed' };
    const result = applyFilters(mockExperts, filters, '');
    expect(result).toHaveLength(1);
    expect(result[0].name).toBe('Dr. Ahmed Ali');
  });
  
  test('filters by availability correctly', () => {
    const filters = { isAvailable: true };
    const result = applyFilters(mockExperts, filters, '');
    expect(result).toHaveLength(1);
    expect(result[0].isAvailable).toBe(true);
  });
});
```

### 5.2 Integration Tests

**Search Page Integration:**
```typescript
// __tests__/SearchPage.integration.test.tsx
describe('SearchPage Integration', () => {
  test('loads experts and applies filters', async () => {
    render(<SearchPage />);
    
    // Wait for data to load
    await waitFor(() => {
      expect(screen.getByText(/Found \d+ experts/)).toBeInTheDocument();
    });
    
    // Apply filter
    const searchInput = screen.getByPlaceholderText('Search name or institution...');
    fireEvent.change(searchInput, { target: { value: 'Ahmed' } });
    
    // Verify filtered results
    await waitFor(() => {
      expect(screen.getByText('Dr. Ahmed Ali')).toBeInTheDocument();
    });
  });
});
```

## 6. Deployment Checklist

### 6.1 Pre-deployment Tasks

- [ ] Performance testing with 2000+ expert records
- [ ] Cross-browser compatibility testing
- [ ] Mobile responsiveness verification
- [ ] Accessibility compliance testing (WCAG AA)
- [ ] Load testing for concurrent users
- [ ] Memory leak testing for extended usage

### 6.2 Production Configuration

```typescript
// config/production.ts
export const PRODUCTION_CONFIG = {
  cache: {
    expertDataTTL: 5 * 60 * 1000, // 5 minutes
    filterCacheSize: 20,
    enableVirtualScrolling: true,
    virtualScrollThreshold: 100
  },
  performance: {
    debounceDelay: 300,
    maxConcurrentRequests: 5,
    enableServiceWorker: true
  },
  monitoring: {
    enableAnalytics: true,
    trackUserInteractions: true,
    performanceMetrics: true
  }
};
```

## Success Metrics & KPIs

### Performance Metrics
- **Initial Load Time**: < 2 seconds for complete expert database
- **Search Response Time**: < 500ms for filter/sort operations
- **Memory Usage**: < 50MB for 2000 expert records
- **UI Responsiveness**: 60fps during scrolling and interactions

### User Experience Metrics
- **Search Success Rate**: > 90% of searches lead to profile views
- **Feature Adoption**: > 80% users utilize column customization
- **Task Completion**: < 30 seconds average time to find and contact expert
- **User Satisfaction**: > 4.5/5 rating in user feedback

This comprehensive Phase 2 implementation plan provides the detailed specifications needed to develop a high-performance, user-friendly expert database browsing interface that serves as the core functionality of the ExpertDB application.