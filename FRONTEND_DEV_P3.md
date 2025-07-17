# ExpertDB Frontend Development - Phase 3 Implementation Plan

## Executive Summary

This document provides a comprehensive implementation plan for Phase 3 of the ExpertDB frontend development, focusing on the **Expert Request System**. This phase establishes the core workflow for non-admin users to submit expert requests and for admin users to review, approve, or reject these requests.

## Phase 3 Scope Overview

### Core Components to Implement

1. **Expert Request Submission System** (User Role)
   - Multi-section form with progressive disclosure
   - Real-time validation and auto-save
   - Biography editor with structured input
   - File upload system with drag-and-drop
   - Draft management system

2. **Admin Review Interface** (Admin Role)
   - Advanced data table with filtering
   - Request detail modal with comprehensive information
   - Batch operations for efficiency
   - Document management with approval workflow

3. **Supporting Infrastructure**
   - Form validation schemas
   - File upload handling
   - State management architecture
   - Error handling and user feedback

## Detailed Implementation Requirements

### 1. Expert Request Submission Form

#### 1.1 Form Structure and Navigation

**Route:** `/expert-requests/new`
**Component:** `ExpertRequestForm.tsx`

```typescript
interface ExpertRequestFormProps {
  initialData?: Partial<ExpertRequestForm>;
  isEditing?: boolean;
  onSubmit: (data: ExpertRequestForm) => Promise<void>;
  onSaveDraft: (data: Partial<ExpertRequestForm>) => Promise<void>;
}
```

**Form Sections:**
1. Personal Information
2. Professional Details  
3. Expertise Areas
4. Biography & Documents

#### 1.2 Section 1: Personal Information

**Required Fields:**
- Full Name: `string` (min 2 chars, required)
- Email: `string` (email format validation, required)
- Phone Number: `string` (phone format validation, required)
- Designation: `string` (professional title, required)
- Institution/Affiliation: `string` (min 2 chars, required)

**Validation Rules:**
```typescript
const personalInfoSchema = {
  name: yup.string().min(2, 'Name must be at least 2 characters').required('Name is required'),
  email: yup.string().email('Invalid email format').required('Email is required'),
  phone: yup.string().matches(/^[+]?[\d\s-()]+$/, 'Invalid phone format').required('Phone is required'),
  designation: yup.string().min(2, 'Designation must be at least 2 characters').required('Designation is required'),
  institution: yup.string().min(2, 'Institution must be at least 2 characters').required('Institution is required')
};
```

#### 1.3 Section 2: Professional Details

**Field Specifications:**
- **Nationality:** Dropdown with options (Bahraini, Non-Bahraini, Unknown)
- **Is Bahraini:** Checkbox (boolean flag)
- **Role:** Checkboxes (Evaluator, Validator - can select both)
- **Employment Type:** Checkboxes (Academic, Employer - can select both)
- **Is Available:** Checkbox (availability status)
- **Is Trained:** Checkbox (training completion status)
- **Is Published:** Checkbox (publication preference, optional)

**Component Implementation:**
```typescript
interface ProfessionalDetailsProps {
  data: ProfessionalDetails;
  onChange: (field: keyof ProfessionalDetails, value: any) => void;
  errors: Record<string, string>;
}

const ProfessionalDetails: React.FC<ProfessionalDetailsProps> = ({ data, onChange, errors }) => {
  return (
    <div className="space-y-4">
      <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
        <Select
          label="Nationality"
          options={[
            { value: 'bahraini', label: 'Bahraini' },
            { value: 'non-bahraini', label: 'Non-Bahraini' },
            { value: 'unknown', label: 'Unknown' }
          ]}
          value={data.nationality}
          onChange={(value) => onChange('nationality', value)}
          error={errors.nationality}
        />
        <div className="flex items-center space-x-2">
          <input
            type="checkbox"
            id="isBahraini"
            checked={data.isBahraini}
            onChange={(e) => onChange('isBahraini', e.target.checked)}
          />
          <label htmlFor="isBahraini">Bahraini National</label>
        </div>
      </div>
      {/* Additional fields... */}
    </div>
  );
};
```

#### 1.4 Section 3: Expertise Areas

**API Integration Requirements:**
- **General Area:** Dropdown populated from `/api/expert/areas` endpoint
- **Specialized Area:** Free text input for specific expertise
- **Skills:** Tag-style input converted to JSON array

**Implementation:**
```typescript
interface ExpertiseAreasProps {
  data: ExpertiseAreas;
  onChange: (field: keyof ExpertiseAreas, value: any) => void;
  errors: Record<string, string>;
}

const ExpertiseAreas: React.FC<ExpertiseAreasProps> = ({ data, onChange, errors }) => {
  const [areas, setAreas] = useState<ExpertArea[]>([]);
  
  useEffect(() => {
    fetchExpertAreas().then(setAreas);
  }, []);

  return (
    <div className="space-y-4">
      <Select
        label="General Area *"
        options={areas.map(area => ({ value: area.id, label: area.name }))}
        value={data.generalArea}
        onChange={(value) => onChange('generalArea', value)}
        error={errors.generalArea}
        required
      />
      <FormField
        label="Specialized Area *"
        value={data.specializedArea}
        onChange={(value) => onChange('specializedArea', value)}
        error={errors.specializedArea}
        required
      />
      <TagInput
        label="Skills *"
        tags={data.skills}
        onChange={(tags) => onChange('skills', tags)}
        error={errors.skills}
        placeholder="Add a skill and press Enter"
        required
      />
    </div>
  );
};
```

#### 1.5 Section 4: Biography & Documents

**Biography Editor Specifications:**
- **Structured Format:** No rich text, plain text only
- **Education Section:** Multiple entries with date ranges
- **Experience Section:** Multiple entries with date ranges
- **Real-time Preview:** Live formatted output display

**Entry Format:**
```typescript
interface BiographyEntry {
  id: string;
  dateFrom: string; // YYYY-MM format
  dateTo: string; // YYYY-MM format
  description: string;
}

interface Biography {
  education: BiographyEntry[];
  experience: BiographyEntry[];
}
```

**Biography Component:**
```typescript
const BiographyEditor: React.FC<BiographyEditorProps> = ({ data, onChange, errors }) => {
  const addEducationEntry = () => {
    const newEntry: BiographyEntry = {
      id: generateId(),
      dateFrom: '',
      dateTo: '',
      description: ''
    };
    onChange('education', [...data.education, newEntry]);
  };

  const updateEducationEntry = (id: string, field: keyof BiographyEntry, value: string) => {
    const updated = data.education.map(entry => 
      entry.id === id ? { ...entry, [field]: value } : entry
    );
    onChange('education', updated);
  };

  const formatBiography = (biography: Biography): string => {
    let formatted = '';
    
    if (biography.education.length > 0) {
      formatted += 'Education:\n';
      biography.education.forEach(entry => {
        formatted += `[${entry.dateFrom} - ${entry.dateTo}] ${entry.description}\n`;
      });
    }
    
    if (biography.experience.length > 0) {
      formatted += '\nExperience:\n';
      biography.experience.forEach(entry => {
        formatted += `[${entry.dateFrom} - ${entry.dateTo}] ${entry.description}\n`;
      });
    }
    
    return formatted;
  };

  return (
    <div className="space-y-6">
      <div>
        <h3 className="text-lg font-semibold mb-4">Education</h3>
        {data.education.map((entry, index) => (
          <BiographyEntryForm
            key={entry.id}
            entry={entry}
            onUpdate={(field, value) => updateEducationEntry(entry.id, field, value)}
            onRemove={() => removeEducationEntry(entry.id)}
            placeholder="Degree, Institution, Location/country[optional]"
          />
        ))}
        <button
          type="button"
          onClick={addEducationEntry}
          className="mt-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Add Education Entry
        </button>
      </div>

      <div>
        <h3 className="text-lg font-semibold mb-4">Experience</h3>
        {data.experience.map((entry, index) => (
          <BiographyEntryForm
            key={entry.id}
            entry={entry}
            onUpdate={(field, value) => updateExperienceEntry(entry.id, field, value)}
            onRemove={() => removeExperienceEntry(entry.id)}
            placeholder="Role/Position, Organization, Location/country[optional]"
          />
        ))}
        <button
          type="button"
          onClick={addExperienceEntry}
          className="mt-2 px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
        >
          Add Experience Entry
        </button>
      </div>

      <div>
        <h3 className="text-lg font-semibold mb-4">Biography Preview</h3>
        <div className="bg-gray-50 p-4 rounded border">
          <pre className="whitespace-pre-wrap text-sm">{formatBiography(data)}</pre>
        </div>
      </div>
    </div>
  );
};
```

#### 1.6 File Upload System

**Requirements:**
- **File Type:** PDF only
- **Size Limit:** 20MB maximum
- **Interface:** Drag-and-drop with fallback file picker
- **Features:** Real-time validation, upload progress, preview

**Implementation:**
```typescript
interface FileUploadProps {
  onFileSelect: (file: File) => void;
  onUpload: (file: File) => Promise<string>; // Returns file URL
  error?: string;
  required?: boolean;
}

const FileUpload: React.FC<FileUploadProps> = ({ onFileSelect, onUpload, error, required }) => {
  const [dragActive, setDragActive] = useState(false);
  const [uploading, setUploading] = useState(false);
  const [uploadProgress, setUploadProgress] = useState(0);
  const [uploadedFile, setUploadedFile] = useState<File | null>(null);

  const handleDrop = (e: React.DragEvent) => {
    e.preventDefault();
    setDragActive(false);
    const files = e.dataTransfer.files;
    if (files.length > 0) {
      handleFileSelect(files[0]);
    }
  };

  const handleFileSelect = async (file: File) => {
    // Validation
    if (file.type !== 'application/pdf') {
      setError('Please upload a PDF file only');
      return;
    }
    
    if (file.size > 20 * 1024 * 1024) { // 20MB
      setError('File size must be less than 20MB');
      return;
    }

    setUploading(true);
    setUploadProgress(0);
    
    try {
      const url = await onUpload(file);
      setUploadedFile(file);
      onFileSelect(file);
    } catch (error) {
      setError('Upload failed. Please try again.');
    } finally {
      setUploading(false);
    }
  };

  return (
    <div className="space-y-4">
      <div
        className={`border-2 border-dashed rounded-lg p-8 text-center transition-colors ${
          dragActive ? 'border-blue-400 bg-blue-50' : 'border-gray-300'
        } ${error ? 'border-red-400' : ''}`}
        onDragOver={(e) => { e.preventDefault(); setDragActive(true); }}
        onDragLeave={() => setDragActive(false)}
        onDrop={handleDrop}
      >
        {uploading ? (
          <div className="space-y-2">
            <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600 mx-auto"></div>
            <p>Uploading... {uploadProgress}%</p>
          </div>
        ) : uploadedFile ? (
          <div className="space-y-2">
            <div className="text-green-600">
              <svg className="w-8 h-8 mx-auto" fill="currentColor" viewBox="0 0 20 20">
                <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <p className="text-sm text-gray-600">
              {uploadedFile.name} ({(uploadedFile.size / 1024 / 1024).toFixed(1)}MB)
            </p>
          </div>
        ) : (
          <div className="space-y-2">
            <svg className="w-8 h-8 mx-auto text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <p className="text-sm text-gray-600">
              Drop your CV here or{' '}
              <button
                type="button"
                className="text-blue-600 hover:text-blue-800"
                onClick={() => document.getElementById('file-input')?.click()}
              >
                browse files
              </button>
            </p>
            <p className="text-xs text-gray-500">PDF files only, maximum 20MB</p>
          </div>
        )}
      </div>
      
      <input
        id="file-input"
        type="file"
        accept=".pdf"
        onChange={(e) => e.target.files?.[0] && handleFileSelect(e.target.files[0])}
        className="hidden"
      />
      
      {error && (
        <p className="text-red-600 text-sm">{error}</p>
      )}
    </div>
  );
};
```

#### 1.7 Form State Management

**Auto-save Implementation:**
```typescript
interface FormState {
  formData: ExpertRequestForm;
  isDraft: boolean;
  isSubmitting: boolean;
  validationErrors: Record<string, string>;
  generalAreas: ExpertArea[];
  uploadProgress: number;
  currentStep: number;
  lastSaved: Date | null;
}

const useExpertRequestForm = (initialData?: Partial<ExpertRequestForm>) => {
  const [state, setState] = useState<FormState>({
    formData: {
      ...defaultFormData,
      ...initialData
    },
    isDraft: true,
    isSubmitting: false,
    validationErrors: {},
    generalAreas: [],
    uploadProgress: 0,
    currentStep: 0,
    lastSaved: null
  });

  // Auto-save every 10 seconds
  useEffect(() => {
    const interval = setInterval(() => {
      if (state.isDraft && hasUnsavedChanges()) {
        saveDraft();
      }
    }, 10000);

    return () => clearInterval(interval);
  }, [state.formData, state.isDraft]);

  const saveDraft = async () => {
    try {
      await api.saveExpertRequestDraft(state.formData);
      setState(prev => ({ ...prev, lastSaved: new Date() }));
      toast.success('Draft saved', { duration: 2000 });
    } catch (error) {
      toast.error('Failed to save draft');
    }
  };

  const submitForm = async () => {
    setState(prev => ({ ...prev, isSubmitting: true }));
    
    try {
      const validationErrors = await validateForm(state.formData);
      if (Object.keys(validationErrors).length > 0) {
        setState(prev => ({ ...prev, validationErrors }));
        return;
      }

      await api.submitExpertRequest(state.formData);
      toast.success('Expert request submitted successfully');
      
      // Redirect to requests list
      router.push('/expert-requests');
    } catch (error) {
      if (error.response?.status === 400) {
        setState(prev => ({ ...prev, validationErrors: error.response.data.errors }));
      } else {
        toast.error('Failed to submit request. Please try again.');
      }
    } finally {
      setState(prev => ({ ...prev, isSubmitting: false }));
    }
  };

  return {
    ...state,
    updateField: (field: keyof ExpertRequestForm, value: any) => {
      setState(prev => ({
        ...prev,
        formData: { ...prev.formData, [field]: value }
      }));
    },
    saveDraft,
    submitForm
  };
};
```

### 2. Admin Review Interface

#### 2.1 Expert Requests List Page

**Route:** `/admin/expert-requests`
**Component:** `AdminRequestTable.tsx`

**Features:**
- Status filter tabs (All, Pending, Approved, Rejected)
- Advanced search across name, institution, specialization
- Sortable columns with visual indicators
- Bulk selection for batch operations
- Pagination with page size options

**Implementation:**
```typescript
interface AdminRequestTableProps {
  requests: ExpertRequest[];
  onStatusChange: (requestId: string, status: RequestStatus) => void;
  onBulkAction: (action: BulkAction, requestIds: string[]) => void;
}

const AdminRequestTable: React.FC<AdminRequestTableProps> = ({ requests, onStatusChange, onBulkAction }) => {
  const [selectedRequests, setSelectedRequests] = useState<string[]>([]);
  const [currentFilter, setCurrentFilter] = useState<RequestStatus | 'all'>('all');
  const [sortConfig, setSortConfig] = useState<{ key: string; direction: 'asc' | 'desc' }>({
    key: 'createdAt',
    direction: 'desc'
  });

  const filteredRequests = useMemo(() => {
    let filtered = requests;
    
    if (currentFilter !== 'all') {
      filtered = filtered.filter(req => req.status === currentFilter);
    }
    
    // Apply sorting
    filtered.sort((a, b) => {
      const aValue = a[sortConfig.key as keyof ExpertRequest];
      const bValue = b[sortConfig.key as keyof ExpertRequest];
      
      if (sortConfig.direction === 'asc') {
        return aValue > bValue ? 1 : -1;
      } else {
        return aValue < bValue ? 1 : -1;
      }
    });
    
    return filtered;
  }, [requests, currentFilter, sortConfig]);

  const getStatusColor = (status: RequestStatus) => {
    switch (status) {
      case 'pending': return 'bg-yellow-100 text-yellow-800';
      case 'approved': return 'bg-green-100 text-green-800';
      case 'rejected': return 'bg-red-100 text-red-800';
      case 'archived': return 'bg-gray-100 text-gray-800';
      default: return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="space-y-4">
      {/* Status Filter Tabs */}
      <div className="flex space-x-1 bg-gray-100 p-1 rounded-lg">
        {['all', 'pending', 'approved', 'rejected', 'archived'].map(status => (
          <button
            key={status}
            onClick={() => setCurrentFilter(status as RequestStatus | 'all')}
            className={`px-3 py-2 rounded-md text-sm font-medium transition-colors ${
              currentFilter === status 
                ? 'bg-white text-gray-900 shadow-sm' 
                : 'text-gray-600 hover:text-gray-900'
            }`}
          >
            {status.charAt(0).toUpperCase() + status.slice(1)}
            {status !== 'all' && (
              <span className="ml-1 bg-gray-200 text-gray-600 px-2 py-1 rounded-full text-xs">
                {requests.filter(r => r.status === status).length}
              </span>
            )}
          </button>
        ))}
      </div>

      {/* Bulk Actions */}
      {selectedRequests.length > 0 && (
        <div className="flex items-center space-x-2 p-3 bg-blue-50 rounded-lg">
          <span className="text-sm text-blue-800">
            {selectedRequests.length} selected
          </span>
          <button
            onClick={() => onBulkAction('approve', selectedRequests)}
            className="px-3 py-1 bg-green-600 text-white rounded text-sm hover:bg-green-700"
          >
            Batch Approve
          </button>
          <button
            onClick={() => onBulkAction('reject', selectedRequests)}
            className="px-3 py-1 bg-red-600 text-white rounded text-sm hover:bg-red-700"
          >
            Batch Reject
          </button>
          <button
            onClick={() => setSelectedRequests([])}
            className="px-3 py-1 bg-gray-600 text-white rounded text-sm hover:bg-gray-700"
          >
            Clear Selection
          </button>
        </div>
      )}

      {/* Data Table */}
      <div className="bg-white shadow rounded-lg overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                <input
                  type="checkbox"
                  checked={selectedRequests.length === filteredRequests.length}
                  onChange={(e) => {
                    if (e.target.checked) {
                      setSelectedRequests(filteredRequests.map(r => r.id));
                    } else {
                      setSelectedRequests([]);
                    }
                  }}
                />
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Name
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Institution
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Specialization
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Date
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {filteredRequests.map((request) => (
              <tr key={request.id} className="hover:bg-gray-50">
                <td className="px-6 py-4 whitespace-nowrap">
                  <input
                    type="checkbox"
                    checked={selectedRequests.includes(request.id)}
                    onChange={(e) => {
                      if (e.target.checked) {
                        setSelectedRequests([...selectedRequests, request.id]);
                      } else {
                        setSelectedRequests(selectedRequests.filter(id => id !== request.id));
                      }
                    }}
                  />
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="text-sm font-medium text-gray-900">{request.name}</div>
                  <div className="text-sm text-gray-500">{request.designation}</div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {request.institution}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {request.specializedArea}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  <span className={`inline-flex px-2 py-1 text-xs font-semibold rounded-full ${getStatusColor(request.status)}`}>
                    {request.status}
                  </span>
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                  {new Date(request.createdAt).toLocaleDateString()}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <button
                    onClick={() => openRequestModal(request)}
                    className="text-blue-600 hover:text-blue-900"
                  >
                    View
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};
```

#### 2.2 Request Detail Modal

**Component:** `RequestDetailModal.tsx`

**Features:**
- Comprehensive expert information display
- Integrated PDF viewer for documents
- Admin action buttons with validation
- Skills display in tag format
- Real-time status updates

**Implementation:**
```typescript
interface RequestDetailModalProps {
  request: ExpertRequest;
  isOpen: boolean;
  onClose: () => void;
  onApprove: (requestId: string, approvalDocument: File) => void;
  onReject: (requestId: string, comments: string) => void;
  onRequestAmendment: (requestId: string, comments: string) => void;
  onArchive: (requestId: string) => void;
}

const RequestDetailModal: React.FC<RequestDetailModalProps> = ({
  request,
  isOpen,
  onClose,
  onApprove,
  onReject,
  onRequestAmendment,
  onArchive
}) => {
  const [activeTab, setActiveTab] = useState<'info' | 'biography' | 'documents' | 'actions'>('info');
  const [approvalDocument, setApprovalDocument] = useState<File | null>(null);
  const [comments, setComments] = useState('');
  const [actionLoading, setActionLoading] = useState(false);

  const handleApprove = async () => {
    if (!approvalDocument) {
      toast.error('Please upload an approval document');
      return;
    }

    setActionLoading(true);
    try {
      await onApprove(request.id, approvalDocument);
      onClose();
    } catch (error) {
      toast.error('Failed to approve request');
    } finally {
      setActionLoading(false);
    }
  };

  const handleReject = async () => {
    if (!comments.trim()) {
      toast.error('Please provide rejection comments');
      return;
    }

    setActionLoading(true);
    try {
      await onReject(request.id, comments);
      onClose();
    } catch (error) {
      toast.error('Failed to reject request');
    } finally {
      setActionLoading(false);
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 z-50 overflow-y-auto">
      <div className="flex items-center justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div className="fixed inset-0 transition-opacity" onClick={onClose}>
          <div className="absolute inset-0 bg-gray-500 opacity-75"></div>
        </div>

        <div className="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-4xl sm:w-full">
          <div className="bg-white px-4 pt-5 pb-4 sm:p-6 sm:pb-4">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-lg font-medium text-gray-900">
                Expert Request Details
              </h3>
              <button
                onClick={onClose}
                className="text-gray-400 hover:text-gray-600"
              >
                <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            {/* Tab Navigation */}
            <div className="border-b border-gray-200 mb-4">
              <nav className="-mb-px flex space-x-8">
                {[
                  { id: 'info', label: 'Information' },
                  { id: 'biography', label: 'Biography' },
                  { id: 'documents', label: 'Documents' },
                  { id: 'actions', label: 'Actions' }
                ].map(tab => (
                  <button
                    key={tab.id}
                    onClick={() => setActiveTab(tab.id as any)}
                    className={`py-2 px-1 border-b-2 font-medium text-sm ${
                      activeTab === tab.id
                        ? 'border-blue-500 text-blue-600'
                        : 'border-transparent text-gray-500 hover:text-gray-700'
                    }`}
                  >
                    {tab.label}
                  </button>
                ))}
              </nav>
            </div>

            {/* Tab Content */}
            <div className="space-y-4">
              {activeTab === 'info' && (
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <h4 className="font-medium text-gray-900 mb-2">Personal Information</h4>
                    <dl className="space-y-1 text-sm">
                      <div className="flex justify-between">
                        <dt className="text-gray-500">Name:</dt>
                        <dd className="text-gray-900">{request.name}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-gray-500">Designation:</dt>
                        <dd className="text-gray-900">{request.designation}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-gray-500">Institution:</dt>
                        <dd className="text-gray-900">{request.institution}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-gray-500">Email:</dt>
                        <dd className="text-gray-900">{request.email}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-gray-500">Phone:</dt>
                        <dd className="text-gray-900">{request.phone}</dd>
                      </div>
                    </dl>
                  </div>
                  
                  <div>
                    <h4 className="font-medium text-gray-900 mb-2">Professional Details</h4>
                    <dl className="space-y-1 text-sm">
                      <div className="flex justify-between">
                        <dt className="text-gray-500">Bahraini:</dt>
                        <dd className="text-gray-900">{request.isBahraini ? 'Yes' : 'No'}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-gray-500">Available:</dt>
                        <dd className="text-gray-900">{request.isAvailable ? 'Yes' : 'No'}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-gray-500">Role:</dt>
                        <dd className="text-gray-900">{request.role}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-gray-500">Employment:</dt>
                        <dd className="text-gray-900">{request.employmentType}</dd>
                      </div>
                      <div className="flex justify-between">
                        <dt className="text-gray-500">Trained:</dt>
                        <dd className="text-gray-900">{request.isTrained ? 'Yes' : 'No'}</dd>
                      </div>
                    </dl>
                  </div>
                </div>
              )}

              {activeTab === 'biography' && (
                <div className="space-y-4">
                  <div>
                    <h4 className="font-medium text-gray-900 mb-2">Specialization</h4>
                    <p className="text-sm text-gray-600 mb-2">
                      <strong>General Area:</strong> {request.generalArea}
                    </p>
                    <p className="text-sm text-gray-600 mb-2">
                      <strong>Specialized Area:</strong> {request.specializedArea}
                    </p>
                    <div className="flex flex-wrap gap-2">
                      {request.skills.map((skill, index) => (
                        <span
                          key={index}
                          className="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded-full"
                        >
                          {skill}
                        </span>
                      ))}
                    </div>
                  </div>
                  
                  <div>
                    <h4 className="font-medium text-gray-900 mb-2">Biography</h4>
                    <div className="bg-gray-50 p-4 rounded border">
                      <pre className="whitespace-pre-wrap text-sm">{request.biography}</pre>
                    </div>
                  </div>
                </div>
              )}

              {activeTab === 'documents' && (
                <div className="space-y-4">
                  <div>
                    <h4 className="font-medium text-gray-900 mb-2">CV Document</h4>
                    <div className="border rounded-lg p-4">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center space-x-2">
                          <svg className="w-8 h-8 text-red-600" fill="currentColor" viewBox="0 0 20 20">
                            <path fillRule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z" clipRule="evenodd" />
                          </svg>
                          <div>
                            <p className="text-sm font-medium">CV Document</p>
                            <p className="text-xs text-gray-500">
                              Uploaded: {new Date(request.createdAt).toLocaleDateString()}
                            </p>
                          </div>
                        </div>
                        <div className="flex space-x-2">
                          <button
                            onClick={() => window.open(request.cvUrl, '_blank')}
                            className="px-3 py-1 bg-blue-600 text-white rounded text-sm hover:bg-blue-700"
                          >
                            View PDF
                          </button>
                          <button
                            onClick={() => downloadFile(request.cvUrl, 'CV.pdf')}
                            className="px-3 py-1 bg-gray-600 text-white rounded text-sm hover:bg-gray-700"
                          >
                            Download
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              )}

              {activeTab === 'actions' && (
                <div className="space-y-4">
                  <div>
                    <h4 className="font-medium text-gray-900 mb-2">Admin Actions</h4>
                    
                    <div className="space-y-4">
                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                          Approval Document (required for approval)
                        </label>
                        <FileUpload
                          onFileSelect={setApprovalDocument}
                          accept=".pdf"
                          error={!approvalDocument && actionLoading ? 'Approval document required' : ''}
                        />
                      </div>

                      <div>
                        <label className="block text-sm font-medium text-gray-700 mb-2">
                          Comments (required for rejection/amendment)
                        </label>
                        <textarea
                          value={comments}
                          onChange={(e) => setComments(e.target.value)}
                          rows={3}
                          className="mt-1 block w-full border border-gray-300 rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                          placeholder="Provide detailed feedback..."
                        />
                      </div>

                      <div className="flex space-x-2">
                        <button
                          onClick={handleApprove}
                          disabled={actionLoading}
                          className="px-4 py-2 bg-green-600 text-white rounded hover:bg-green-700 disabled:opacity-50"
                        >
                          {actionLoading ? 'Processing...' : 'Approve'}
                        </button>
                        <button
                          onClick={handleReject}
                          disabled={actionLoading}
                          className="px-4 py-2 bg-red-600 text-white rounded hover:bg-red-700 disabled:opacity-50"
                        >
                          {actionLoading ? 'Processing...' : 'Reject'}
                        </button>
                        <button
                          onClick={() => onRequestAmendment(request.id, comments)}
                          disabled={actionLoading}
                          className="px-4 py-2 bg-yellow-600 text-white rounded hover:bg-yellow-700 disabled:opacity-50"
                        >
                          Request Amendment
                        </button>
                        <button
                          onClick={() => onArchive(request.id)}
                          disabled={actionLoading}
                          className="px-4 py-2 bg-gray-600 text-white rounded hover:bg-gray-700 disabled:opacity-50"
                        >
                          Archive
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
```

### 3. Supporting Components

#### 3.1 TagInput Component

```typescript
interface TagInputProps {
  label: string;
  tags: string[];
  onChange: (tags: string[]) => void;
  error?: string;
  placeholder?: string;
  required?: boolean;
}

const TagInput: React.FC<TagInputProps> = ({ label, tags, onChange, error, placeholder, required }) => {
  const [inputValue, setInputValue] = useState('');
  const [focused, setFocused] = useState(false);

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && inputValue.trim()) {
      e.preventDefault();
      if (!tags.includes(inputValue.trim())) {
        onChange([...tags, inputValue.trim()]);
      }
      setInputValue('');
    } else if (e.key === 'Backspace' && !inputValue && tags.length > 0) {
      onChange(tags.slice(0, -1));
    }
  };

  const removeTag = (index: number) => {
    onChange(tags.filter((_, i) => i !== index));
  };

  return (
    <div className="space-y-1">
      <label className="block text-sm font-medium text-gray-700">
        {label}
        {required && <span className="text-red-500 ml-1">*</span>}
      </label>
      <div
        className={`min-h-[2.5rem] w-full border rounded-md px-3 py-2 focus-within:ring-2 focus-within:ring-blue-500 focus-within:border-blue-500 ${
          error ? 'border-red-400' : 'border-gray-300'
        }`}
      >
        <div className="flex flex-wrap gap-2">
          {tags.map((tag, index) => (
            <span
              key={index}
              className="inline-flex items-center px-2 py-1 bg-blue-100 text-blue-800 text-sm rounded-full"
            >
              {tag}
              <button
                type="button"
                onClick={() => removeTag(index)}
                className="ml-1 text-blue-600 hover:text-blue-800"
              >
                <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </span>
          ))}
          <input
            type="text"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
            onKeyDown={handleKeyDown}
            onFocus={() => setFocused(true)}
            onBlur={() => setFocused(false)}
            placeholder={tags.length === 0 ? placeholder : ''}
            className="flex-1 min-w-[120px] border-none outline-none bg-transparent"
          />
        </div>
      </div>
      {error && (
        <p className="text-red-600 text-sm">{error}</p>
      )}
    </div>
  );
};
```

#### 3.2 Select Component

```typescript
interface SelectOption {
  value: string | number;
  label: string;
}

interface SelectProps {
  label: string;
  options: SelectOption[];
  value: string | number;
  onChange: (value: string | number) => void;
  error?: string;
  placeholder?: string;
  required?: boolean;
}

const Select: React.FC<SelectProps> = ({ label, options, value, onChange, error, placeholder, required }) => {
  return (
    <div className="space-y-1">
      <label className="block text-sm font-medium text-gray-700">
        {label}
        {required && <span className="text-red-500 ml-1">*</span>}
      </label>
      <select
        value={value}
        onChange={(e) => onChange(e.target.value)}
        className={`mt-1 block w-full border rounded-md px-3 py-2 text-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500 ${
          error ? 'border-red-400' : 'border-gray-300'
        }`}
      >
        {placeholder && (
          <option value="" disabled>
            {placeholder}
          </option>
        )}
        {options.map((option) => (
          <option key={option.value} value={option.value}>
            {option.label}
          </option>
        ))}
      </select>
      {error && (
        <p className="text-red-600 text-sm">{error}</p>
      )}
    </div>
  );
};
```

### 4. API Integration

#### 4.1 API Service Functions

```typescript
import axios from 'axios';

const api = axios.create({
  baseURL: '/api',
  headers: {
    'Content-Type': 'application/json',
  },
});

// Request interceptor for auth token
api.interceptors.request.use((config) => {
  const token = localStorage.getItem('token');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export const expertRequestApi = {
  // Get all expert requests (admin only)
  getExpertRequests: async (params?: {
    status?: string;
    search?: string;
    page?: number;
    limit?: number;
  }) => {
    const response = await api.get('/expert-requests', { params });
    return response.data;
  },

  // Get user's expert requests
  getMyExpertRequests: async () => {
    const response = await api.get('/expert-requests/mine');
    return response.data;
  },

  // Get single expert request
  getExpertRequest: async (id: string) => {
    const response = await api.get(`/expert-requests/${id}`);
    return response.data;
  },

  // Create new expert request
  createExpertRequest: async (data: ExpertRequestForm) => {
    const formData = new FormData();
    
    // Append all fields
    Object.entries(data).forEach(([key, value]) => {
      if (key === 'cv' && value instanceof File) {
        formData.append('cv', value);
      } else if (key === 'skills' && Array.isArray(value)) {
        formData.append('skills', JSON.stringify(value));
      } else if (key === 'biography' && typeof value === 'object') {
        formData.append('biography', JSON.stringify(value));
      } else {
        formData.append(key, String(value));
      }
    });

    const response = await api.post('/expert-requests', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  // Save draft
  saveDraft: async (data: Partial<ExpertRequestForm>) => {
    const response = await api.post('/expert-requests/draft', data);
    return response.data;
  },

  // Update expert request
  updateExpertRequest: async (id: string, data: Partial<ExpertRequestForm>) => {
    const response = await api.put(`/expert-requests/${id}`, data);
    return response.data;
  },

  // Admin actions
  approveRequest: async (id: string, approvalDocument: File) => {
    const formData = new FormData();
    formData.append('approvalDocument', approvalDocument);
    
    const response = await api.post(`/expert-requests/${id}/approve`, formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  rejectRequest: async (id: string, comments: string) => {
    const response = await api.post(`/expert-requests/${id}/reject`, { comments });
    return response.data;
  },

  requestAmendment: async (id: string, comments: string) => {
    const response = await api.post(`/expert-requests/${id}/request-amendment`, { comments });
    return response.data;
  },

  archiveRequest: async (id: string) => {
    const response = await api.post(`/expert-requests/${id}/archive`);
    return response.data;
  },

  // Batch operations
  batchApprove: async (ids: string[], approvalDocument: File) => {
    const formData = new FormData();
    formData.append('approvalDocument', approvalDocument);
    formData.append('requestIds', JSON.stringify(ids));
    
    const response = await api.post('/expert-requests/batch-approve', formData, {
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
    return response.data;
  },

  batchReject: async (ids: string[], comments: string) => {
    const response = await api.post('/expert-requests/batch-reject', {
      requestIds: ids,
      comments
    });
    return response.data;
  },
};

// Expert areas API
export const expertAreasApi = {
  getExpertAreas: async () => {
    const response = await api.get('/expert/areas');
    return response.data;
  },
};
```

#### 4.2 Form Validation Schema

```typescript
import * as yup from 'yup';

export const expertRequestSchema = yup.object({
  // Personal Information
  name: yup.string()
    .min(2, 'Name must be at least 2 characters')
    .required('Name is required'),
  
  email: yup.string()
    .email('Invalid email format')
    .required('Email is required'),
  
  phone: yup.string()
    .matches(/^[+]?[\d\s-()]+$/, 'Invalid phone format')
    .required('Phone number is required'),
  
  designation: yup.string()
    .min(2, 'Designation must be at least 2 characters')
    .required('Designation is required'),
  
  institution: yup.string()
    .min(2, 'Institution must be at least 2 characters')
    .required('Institution is required'),

  // Professional Details
  nationality: yup.string()
    .oneOf(['bahraini', 'non-bahraini', 'unknown'])
    .required('Nationality is required'),
  
  isBahraini: yup.boolean()
    .required('Bahraini status is required'),
  
  role: yup.string()
    .oneOf(['evaluator', 'validator', 'evaluator/validator'])
    .required('Role is required'),
  
  employmentType: yup.string()
    .oneOf(['academic', 'employer', 'academic/employer'])
    .required('Employment type is required'),
  
  isAvailable: yup.boolean()
    .required('Availability status is required'),
  
  isTrained: yup.boolean()
    .required('Training status is required'),
  
  isPublished: yup.boolean()
    .default(false),

  // Expertise Areas
  generalArea: yup.number()
    .integer('General area must be a valid selection')
    .required('General area is required'),
  
  specializedArea: yup.string()
    .min(2, 'Specialized area must be at least 2 characters')
    .required('Specialized area is required'),
  
  skills: yup.array()
    .of(yup.string().min(1, 'Skill cannot be empty'))
    .min(1, 'At least one skill is required')
    .required('Skills are required'),

  // Biography
  biography: yup.object({
    education: yup.array()
      .of(yup.object({
        dateFrom: yup.string()
          .matches(/^\d{4}-\d{2}$/, 'Date must be in YYYY-MM format')
          .required('Start date is required'),
        dateTo: yup.string()
          .matches(/^\d{4}-\d{2}$/, 'Date must be in YYYY-MM format')
          .required('End date is required'),
        description: yup.string()
          .min(5, 'Description must be at least 5 characters')
          .required('Description is required')
      }))
      .min(1, 'At least one education entry is required')
      .required('Education is required'),
    
    experience: yup.array()
      .of(yup.object({
        dateFrom: yup.string()
          .matches(/^\d{4}-\d{2}$/, 'Date must be in YYYY-MM format')
          .required('Start date is required'),
        dateTo: yup.string()
          .matches(/^\d{4}-\d{2}$/, 'Date must be in YYYY-MM format')
          .required('End date is required'),
        description: yup.string()
          .min(5, 'Description must be at least 5 characters')
          .required('Description is required')
      }))
      .min(1, 'At least one experience entry is required')
      .required('Experience is required')
  }).required('Biography is required'),

  // Documents
  cv: yup.mixed()
    .test('fileSize', 'File size must be less than 20MB', (value) => {
      if (!value) return false;
      return value instanceof File && value.size <= 20 * 1024 * 1024;
    })
    .test('fileType', 'Only PDF files are allowed', (value) => {
      if (!value) return false;
      return value instanceof File && value.type === 'application/pdf';
    })
    .required('CV document is required')
});
```

### 5. State Management

#### 5.1 Context Provider

```typescript
import React, { createContext, useContext, useReducer } from 'react';

interface ExpertRequestState {
  requests: ExpertRequest[];
  loading: boolean;
  error: string | null;
  currentRequest: ExpertRequest | null;
  expertAreas: ExpertArea[];
}

type ExpertRequestAction = 
  | { type: 'SET_LOADING'; payload: boolean }
  | { type: 'SET_ERROR'; payload: string | null }
  | { type: 'SET_REQUESTS'; payload: ExpertRequest[] }
  | { type: 'SET_CURRENT_REQUEST'; payload: ExpertRequest | null }
  | { type: 'SET_EXPERT_AREAS'; payload: ExpertArea[] }
  | { type: 'ADD_REQUEST'; payload: ExpertRequest }
  | { type: 'UPDATE_REQUEST'; payload: ExpertRequest }
  | { type: 'DELETE_REQUEST'; payload: string };

const initialState: ExpertRequestState = {
  requests: [],
  loading: false,
  error: null,
  currentRequest: null,
  expertAreas: []
};

const expertRequestReducer = (state: ExpertRequestState, action: ExpertRequestAction): ExpertRequestState => {
  switch (action.type) {
    case 'SET_LOADING':
      return { ...state, loading: action.payload };
    
    case 'SET_ERROR':
      return { ...state, error: action.payload };
    
    case 'SET_REQUESTS':
      return { ...state, requests: action.payload };
    
    case 'SET_CURRENT_REQUEST':
      return { ...state, currentRequest: action.payload };
    
    case 'SET_EXPERT_AREAS':
      return { ...state, expertAreas: action.payload };
    
    case 'ADD_REQUEST':
      return { ...state, requests: [...state.requests, action.payload] };
    
    case 'UPDATE_REQUEST':
      return {
        ...state,
        requests: state.requests.map(req =>
          req.id === action.payload.id ? action.payload : req
        )
      };
    
    case 'DELETE_REQUEST':
      return {
        ...state,
        requests: state.requests.filter(req => req.id !== action.payload)
      };
    
    default:
      return state;
  }
};

const ExpertRequestContext = createContext<{
  state: ExpertRequestState;
  dispatch: React.Dispatch<ExpertRequestAction>;
  actions: {
    loadRequests: () => Promise<void>;
    loadExpertAreas: () => Promise<void>;
    createRequest: (data: ExpertRequestForm) => Promise<void>;
    updateRequest: (id: string, data: Partial<ExpertRequestForm>) => Promise<void>;
    approveRequest: (id: string, approvalDocument: File) => Promise<void>;
    rejectRequest: (id: string, comments: string) => Promise<void>;
  };
} | null>(null);

export const ExpertRequestProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [state, dispatch] = useReducer(expertRequestReducer, initialState);

  const actions = {
    loadRequests: async () => {
      dispatch({ type: 'SET_LOADING', payload: true });
      try {
        const requests = await expertRequestApi.getExpertRequests();
        dispatch({ type: 'SET_REQUESTS', payload: requests });
      } catch (error) {
        dispatch({ type: 'SET_ERROR', payload: 'Failed to load requests' });
      } finally {
        dispatch({ type: 'SET_LOADING', payload: false });
      }
    },

    loadExpertAreas: async () => {
      try {
        const areas = await expertAreasApi.getExpertAreas();
        dispatch({ type: 'SET_EXPERT_AREAS', payload: areas });
      } catch (error) {
        dispatch({ type: 'SET_ERROR', payload: 'Failed to load expert areas' });
      }
    },

    createRequest: async (data: ExpertRequestForm) => {
      dispatch({ type: 'SET_LOADING', payload: true });
      try {
        const request = await expertRequestApi.createExpertRequest(data);
        dispatch({ type: 'ADD_REQUEST', payload: request });
      } catch (error) {
        dispatch({ type: 'SET_ERROR', payload: 'Failed to create request' });
        throw error;
      } finally {
        dispatch({ type: 'SET_LOADING', payload: false });
      }
    },

    updateRequest: async (id: string, data: Partial<ExpertRequestForm>) => {
      dispatch({ type: 'SET_LOADING', payload: true });
      try {
        const request = await expertRequestApi.updateExpertRequest(id, data);
        dispatch({ type: 'UPDATE_REQUEST', payload: request });
      } catch (error) {
        dispatch({ type: 'SET_ERROR', payload: 'Failed to update request' });
        throw error;
      } finally {
        dispatch({ type: 'SET_LOADING', payload: false });
      }
    },

    approveRequest: async (id: string, approvalDocument: File) => {
      dispatch({ type: 'SET_LOADING', payload: true });
      try {
        const request = await expertRequestApi.approveRequest(id, approvalDocument);
        dispatch({ type: 'UPDATE_REQUEST', payload: request });
      } catch (error) {
        dispatch({ type: 'SET_ERROR', payload: 'Failed to approve request' });
        throw error;
      } finally {
        dispatch({ type: 'SET_LOADING', payload: false });
      }
    },

    rejectRequest: async (id: string, comments: string) => {
      dispatch({ type: 'SET_LOADING', payload: true });
      try {
        const request = await expertRequestApi.rejectRequest(id, comments);
        dispatch({ type: 'UPDATE_REQUEST', payload: request });
      } catch (error) {
        dispatch({ type: 'SET_ERROR', payload: 'Failed to reject request' });
        throw error;
      } finally {
        dispatch({ type: 'SET_LOADING', payload: false });
      }
    }
  };

  return (
    <ExpertRequestContext.Provider value={{ state, dispatch, actions }}>
      {children}
    </ExpertRequestContext.Provider>
  );
};

export const useExpertRequest = () => {
  const context = useContext(ExpertRequestContext);
  if (!context) {
    throw new Error('useExpertRequest must be used within ExpertRequestProvider');
  }
  return context;
};
```

### 6. Error Handling & User Feedback

#### 6.1 Toast Notification System

```typescript
import React, { createContext, useContext, useState, useCallback } from 'react';

interface Toast {
  id: string;
  message: string;
  type: 'success' | 'error' | 'warning' | 'info';
  duration?: number;
}

interface ToastContextType {
  toasts: Toast[];
  addToast: (message: string, type: Toast['type'], duration?: number) => void;
  removeToast: (id: string) => void;
  clearAll: () => void;
}

const ToastContext = createContext<ToastContextType | null>(null);

export const ToastProvider: React.FC<{ children: React.ReactNode }> = ({ children }) => {
  const [toasts, setToasts] = useState<Toast[]>([]);

  const addToast = useCallback((message: string, type: Toast['type'], duration = 5000) => {
    const id = Math.random().toString(36).substr(2, 9);
    const toast: Toast = { id, message, type, duration };
    
    setToasts(prev => [...prev, toast]);
    
    if (duration > 0) {
      setTimeout(() => removeToast(id), duration);
    }
  }, []);

  const removeToast = useCallback((id: string) => {
    setToasts(prev => prev.filter(toast => toast.id !== id));
  }, []);

  const clearAll = useCallback(() => {
    setToasts([]);
  }, []);

  return (
    <ToastContext.Provider value={{ toasts, addToast, removeToast, clearAll }}>
      {children}
      <ToastContainer />
    </ToastContext.Provider>
  );
};

const ToastContainer: React.FC = () => {
  const { toasts, removeToast } = useToast();

  return (
    <div className="fixed top-4 right-4 z-50 space-y-2">
      {toasts.map(toast => (
        <ToastItem key={toast.id} toast={toast} onClose={() => removeToast(toast.id)} />
      ))}
    </div>
  );
};

const ToastItem: React.FC<{ toast: Toast; onClose: () => void }> = ({ toast, onClose }) => {
  const getToastStyles = (type: Toast['type']) => {
    switch (type) {
      case 'success':
        return 'bg-green-500 text-white';
      case 'error':
        return 'bg-red-500 text-white';
      case 'warning':
        return 'bg-yellow-500 text-white';
      case 'info':
        return 'bg-blue-500 text-white';
      default:
        return 'bg-gray-500 text-white';
    }
  };

  return (
    <div className={`flex items-center justify-between p-4 rounded-lg shadow-lg max-w-md ${getToastStyles(toast.type)}`}>
      <span className="text-sm">{toast.message}</span>
      <button onClick={onClose} className="ml-4 text-white hover:text-gray-200">
        <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  );
};

export const useToast = () => {
  const context = useContext(ToastContext);
  if (!context) {
    throw new Error('useToast must be used within ToastProvider');
  }
  return context;
};

// Convenience hook for common toast actions
export const useToastActions = () => {
  const { addToast } = useToast();
  
  return {
    success: (message: string) => addToast(message, 'success'),
    error: (message: string) => addToast(message, 'error'),
    warning: (message: string) => addToast(message, 'warning'),
    info: (message: string) => addToast(message, 'info')
  };
};
```

#### 6.2 Error Boundary Component

```typescript
import React, { Component, ErrorInfo, ReactNode } from 'react';

interface Props {
  children: ReactNode;
  fallback?: ReactNode;
}

interface State {
  hasError: boolean;
  error?: Error;
}

class ErrorBoundary extends Component<Props, State> {
  public state: State = {
    hasError: false
  };

  public static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  public componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('Error caught by boundary:', error, errorInfo);
  }

  public render() {
    if (this.state.hasError) {
      return this.props.fallback || (
        <div className="flex items-center justify-center min-h-screen bg-gray-50">
          <div className="text-center">
            <h2 className="text-2xl font-bold text-gray-900 mb-4">Something went wrong</h2>
            <p className="text-gray-600 mb-4">
              An error occurred while loading this page. Please try refreshing or contact support.
            </p>
            <button
              onClick={() => window.location.reload()}
              className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700"
            >
              Refresh Page
            </button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}

export default ErrorBoundary;
```

## Implementation Timeline

### Week 1: Foundation Components
- [ ] Set up project structure and routing
- [ ] Implement basic form components (FormField, Select, TagInput)
- [ ] Create error handling and toast notification system
- [ ] Set up API service layer

### Week 2: Expert Request Form
- [ ] Build multi-section form structure
- [ ] Implement personal information section
- [ ] Implement professional details section
- [ ] Add form validation and error handling

### Week 3: Biography Editor & File Upload
- [ ] Create biography editor with structured input
- [ ] Implement file upload component with drag-and-drop
- [ ] Add real-time preview functionality
- [ ] Implement auto-save mechanism

### Week 4: Admin Review Interface
- [ ] Build admin requests data table
- [ ] Create request detail modal
- [ ] Implement batch operations
- [ ] Add sorting and filtering functionality

### Week 5: Integration & Testing
- [ ] Connect all components to backend API
- [ ] Implement state management
- [ ] Add comprehensive error handling
- [ ] Conduct user acceptance testing

### Week 6: Polish & Optimization
- [ ] Performance optimization
- [ ] Accessibility improvements
- [ ] Mobile responsiveness
- [ ] Documentation and deployment

## Success Criteria

### Functional Requirements
- [ ] Users can successfully submit expert requests with all required information
- [ ] Form validation prevents invalid submissions
- [ ] Auto-save functionality prevents data loss
- [ ] Admin can review, approve, reject, and manage requests efficiently
- [ ] File upload supports PDF files up to 20MB
- [ ] Batch operations work correctly for multiple requests

### Technical Requirements
- [ ] Page load times < 2 seconds
- [ ] Form submission success rate > 95%
- [ ] Mobile-responsive design
- [ ] WCAG AA accessibility compliance
- [ ] Comprehensive error handling
- [ ] Clean, maintainable code structure

### User Experience Requirements
- [ ] Intuitive navigation and workflow
- [ ] Clear feedback for all user actions
- [ ] Consistent visual design
- [ ] Efficient admin workflow
- [ ] Helpful error messages and guidance

This comprehensive implementation plan provides all the necessary specifications, code examples, and guidelines to successfully implement Phase 3 of the ExpertDB frontend development. The plan prioritizes user experience, technical excellence, and maintainability while ensuring the application meets all functional requirements.