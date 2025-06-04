# Expert Database Frontend Development Plan

## Project Overview

This document outlines the comprehensive plan for developing the frontend of the Expert Database system. The design will be based on the BQA (Education & Training Quality Authority) website design system to maintain a professional, government-appropriate aesthetic while implementing the specific functionality required for the expert request and admin review process.

## Design System Reference

Based on the analysis of the BQA website (`bqa_website/bqa_website_analysis/`), we will adopt their proven design patterns and visual language.

### Visual Identity & Branding
- **Primary Colors**:
  - BQA Green: `#397b26` (Primary brand color, for approvals and positive actions)
  - Deep Blue: `#1c4679` (Navigation, headers, and authority elements)
  - Accent Red: `#e64125` (Call-to-action, warnings, and rejections)
  - Orange: `#e68835` (Secondary actions and highlights)

- **Typography System**:
  - Primary Font: "Graphik" family (Regular 400, Medium 500, Semibold 600, Bold 700)
  - Fallback: "Segoe UI", Helvetica Neue, Arial, sans-serif
  - Arabic Support: "Graphik Arabic" for potential future localization
  - Hierarchy: h1(42px), h2(38px), h3(35px), h4(26px), body(16px)

- **Layout Framework**:
  - Bootstrap-based responsive grid system
  - Container max-width: 1212px
  - Responsive breakpoints: xs(<576px), sm(≥576px), md(≥768px), lg(≥992px), xl(≥1200px)

## Required Pages & Components

### 1. Login Page (Landing Page)

**Purpose**: Main entry point for the system, providing authentication for both users and administrators.

**Design Reference**: Based on BQA's clean header design (`01_header_hero_section.webp`) and form styling patterns.

**Layout Structure**:
```
┌─────────────────────────────────────────┐
│             Header Section              │
│   [Logo] [System Title] [Language]      │
├─────────────────────────────────────────┤
│                                         │
│             Hero Section                │
│        [Background Image/Pattern]       │
│                                         │
│    ┌─────────────────────────────────┐   │
│    │        Login Form Card          │   │
│    │  ┌─────────────────────────────┐ │   │
│    │  │      Email Field           │ │   │
│    │  │      Password Field        │ │   │
│    │  │   [Remember Me] [Login]    │ │   │
│    │  │                            │ │   │
│    │  │    [Forgot Password?]      │ │   │
│    │  └─────────────────────────────┘ │   │
│    └─────────────────────────────────┘   │
│                                         │
├─────────────────────────────────────────┤
│              Footer Section             │
│    [Links] [Contact] [Copyright]        │
└─────────────────────────────────────────┘
```

**Component Implementation**:
- **Header Component**: Clean navigation bar with system logo and title
- **Login Card**: Centered form with BQA styling patterns
- **Form Elements**: Custom-styled inputs following BQA form patterns
- **Footer**: Minimal footer with essential links

**Styling References**:
- Header styling from `bqa_main_style.css` lines 800-950 (header section)
- Form styling from `bqa_main_style.css` lines 600-700 (form elements)
- Button styling from `bqa_main_style.css` lines 1200-1350 (btn classes)

### 2. Expert Request Submission Page

**Purpose**: Allow users/planners to submit expert requests with all required details and CV upload.

**Design Reference**: Based on BQA's form layouts and structured content sections.

**Form Fields Based on Backend Schema**:

#### **Section 1: Personal Information**
- `name` (text, required): Expert's full name
- `designation` (text, required): Professional title or position
- `institution` (text, required): Organization affiliation
- `phone` (text, required): Contact phone number
- `email` (text, required): Contact email address

#### **Section 2: Professional Details**
- `isBahraini` (boolean, required): Nationality checkbox
- `isAvailable` (boolean, required): Availability status
- `rating` (text, required): Performance rating
- `role` (select, required): Options: "evaluator", "validator", "evaluator/validator"
- `employmentType` (select, required): Options: "academic", "employer"
- `isTrained` (boolean, required): Training completion status
- `isPublished` (boolean, optional): Publication preference (defaults to false)

#### **Section 3: Expertise Areas**
- `generalArea` (select, required): ID from expert_areas table
- `specializedArea` (text, required): Specific field of specialization
- `skills` (text/JSON array, required): Expert's skills and competencies

#### **Section 4: Biography & Documents**
- `biography` (textarea, required, max 1000 chars): Professional summary
- `cv` (file upload, required): PDF document upload

**Layout Structure**:
```
┌─────────────────────────────────────────┐
│              Navigation                 │
├─────────────────────────────────────────┤
│           Page Header Section           │
│     "Submit Expert Request"             │
│         [Progress Indicator]            │
├─────────────────────────────────────────┤
│                                         │
│     Multi-Section Request Form         │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │    Section 1: Personal Info        │ │
│  │  ┌───────────┐ ┌───────────────┐   │ │
│  │  │Name*      │ │Designation*   │   │ │
│  │  └───────────┘ └───────────────┘   │ │
│  │  ┌───────────┐ ┌───────────────┐   │ │
│  │  │Phone*     │ │Email*         │   │ │
│  │  └───────────┘ └───────────────┘   │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │Institution*                   │ │ │
│  │  └─────────────────────────────────┘ │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │   Section 2: Professional Details  │ │
│  │  ┌───────────┐ ┌───────────────┐   │ │
│  │  │☐ Bahraini │ │☐ Available    │   │ │
│  │  └───────────┘ └───────────────┘   │ │
│  │  ┌───────────┐ ┌───────────────┐   │ │
│  │  │Role*      │ │Employment*    │   │ │
│  │  │[Dropdown] │ │[Dropdown]     │   │ │
│  │  └───────────┘ └───────────────┘   │ │
│  │  ┌───────────┐ ┌───────────────┐   │ │
│  │  │Rating*    │ │☐ Trained      │   │ │
│  │  └───────────┘ └───────────────┘   │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │    Section 3: Expertise Areas      │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │General Area* [Dropdown]         │ │ │
│  │  └─────────────────────────────────┘ │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │Specialized Area*                │ │ │
│  │  └─────────────────────────────────┘ │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │Skills* (comma-separated)        │ │ │
│  │  └─────────────────────────────────┘ │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │   Section 4: Biography & Documents │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │Biography* (max 1000 chars)      │ │ │
│  │  │[Rich Text Editor/Textarea]      │ │ │
│  │  └─────────────────────────────────┘ │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │CV Upload* (PDF only)            │ │ │
│  │  │[Drag & Drop Area]               │ │ │
│  │  └─────────────────────────────────┘ │ │
│  └─────────────────────────────────────┘ │
│                                         │
│    [Save Draft] [Submit Request]        │
├─────────────────────────────────────────┤
│              Footer                     │
└─────────────────────────────────────────┘
```

**Key Components**:

1. **Form Sections with Progressive Disclosure**:
   - Collapsible/expandable sections for better UX
   - Form validation with real-time feedback
   - Required field indicators (* asterisks)
   - Character count for biography field

2. **Smart Form Elements**:
   - **General Area Dropdown**: Populated from `/api/expert/areas` endpoint
   - **Role Select**: "evaluator", "validator", "evaluator/validator" options
   - **Employment Type Select**: "academic", "employer" options
   - **Skills Input**: Tag-style input for multiple skills
   - **Biography Textarea**: Character counter (1000 max)

3. **Document Upload Enhancement**:
   - PDF-only validation for CV uploads
   - File size limits and progress indication
   - Preview functionality after upload
   - Error handling for upload failures

4. **Form State Management**:
   - Auto-save draft functionality
   - Form validation before submission
   - Success/error feedback upon submission

**Styling References**:
- Progressive form styling from BQA patterns
- File upload styling (`bqa_main_style.css` lines 400-450)
- Multi-section form layout with BQA card styling

### 3. Admin Panel - Expert Request Review

**Purpose**: Administrative interface for reviewing, approving, or rejecting expert requests with detailed workflow management.

**Design Reference**: Based on BQA's service cards and data table patterns for government applications.

**Data Table Structure Based on Backend API**:

#### **Main Request Table Columns** (from `GET /api/expert-requests`):
- **ID**: Request identifier for tracking
- **Expert Name**: `name` field from request
- **Institution**: `institution` field 
- **Specialization**: `specializedArea` field
- **Status**: `status` field ("pending", "approved", "rejected")
- **Submitted By**: `createdBy` field (user who submitted)
- **Date Submitted**: `createdAt` field
- **Actions**: View, Approve, Reject buttons

#### **Status Filter Tabs**:
- **All Requests**: Show all statuses
- **Pending Review** (default): `status: "pending"`
- **Approved**: `status: "approved"` 
- **Rejected**: `status: "rejected"`

**Layout Structure**:
```
┌─────────────────────────────────────────┐
│           Admin Navigation              │
│  [Dashboard] [Requests] [Experts] [...]  │
├─────────────────────────────────────────┤
│          Admin Header Section           │
│      "Expert Request Management"        │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │    Status Filter Tabs               │ │
│  │ [All] [Pending] [Approved] [Rejected]│ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │     Search & Filter Controls        │ │
│  │ [Search] [Institution] [Area] [Sort] │ │
│  └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│                                         │
│             Data Table                  │
│                                         │
│ ┌─────┬──────────┬─────────────┬────────┐ │
│ │ ID  │ Name     │ Institution │ Status │ │
│ ├─────┼──────────┼─────────────┼────────┤ │
│ │ 001 │ Dr. Ali  │ University  │ [🟡]   │ │
│ │     │ Ahmed    │ of Bahrain  │ Pending│ │
│ │     │          │             │ [View] │ │
│ ├─────┼──────────┼─────────────┼────────┤ │
│ │ 002 │ Sarah    │ Ministry    │ [🔵]   │ │
│ │     │ Hassan   │ of Health   │ Apprvd │ │
│ │     │          │             │ [View] │ │
│ ├─────┼──────────┼─────────────┼────────┤ │
│ │ 003 │ Ahmed    │ BCCI        │ [🔴]   │ │
│ │     │ Al-Said  │             │ Reject │ │
│ │     │          │             │ [View] │ │
│ └─────┴──────────┴─────────────┴────────┘ │
│                                         │
│         [Batch Actions] [Pagination]     │
├─────────────────────────────────────────┤
│              Footer                     │
└─────────────────────────────────────────┘
```

**Request Detail Modal** (triggered by View button):
```
┌─────────────────────────────────────────┐
│          Request Details Modal          │
│                                    [✕]  │
├─────────────────────────────────────────┤
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │         Expert Information          │ │
│  │  Name: Dr. Ali Ahmed                │ │
│  │  Designation: Senior Engineer       │ │
│  │  Institution: University of Bahrain │ │
│  │  Email: ali.ahmed@uob.edu.bh        │ │
│  │  Phone: +973 1234 5678              │ │
│  │  Bahraini: ✓ Yes   Available: ✓ Yes │ │
│  │  Role: evaluator/validator          │ │
│  │  Employment: academic               │ │
│  │  General Area: Engineering          │ │
│  │  Specialized: Civil Engineering     │ │
│  │  Skills: [Structural] [Design] [...] │ │
│  │  Training: ✓ Completed              │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │            Biography                │ │
│  │  [Scrollable text area with full    │ │
│  │   biography content up to 1000      │ │
│  │   characters...]                    │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │           Documents                 │ │
│  │  📄 CV Document                     │ │
│  │     [View PDF] [Download]           │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │         Admin Actions               │ │
│  │  Approval Document: [Choose File]   │ │
│  │                                     │ │
│  │  [✅ Approve] [❌ Reject] [✏️ Edit]  │ │
│  │                                     │ │
│  │  Rejection Reason (if rejecting):   │ │
│  │  [Text area for feedback]           │ │
│  └─────────────────────────────────────┘ │
│                                         │
├─────────────────────────────────────────┤
│        [Cancel] [Save Changes]          │
└─────────────────────────────────────────┘
```

**Key Components**:

1. **Advanced Data Table**:
   - **Sortable Columns**: Name, Institution, Date, Status
   - **Status Indicators**: Color-coded badges (🟡 Pending, 🔵 Approved, 🔴 Rejected)
   - **Quick Actions**: View details, bulk select checkboxes
   - **Pagination**: Standard pagination with page size options
   - **Search Functionality**: Real-time search across name, institution, specialization

2. **Filtering & Search System**:
   - **Status Tabs**: Filter by pending/approved/rejected status
   - **Institution Filter**: Dropdown of common institutions
   - **General Area Filter**: Filter by expertise area
   - **Date Range**: Filter by submission date
   - **Clear Filters**: Reset all filters option

3. **Request Detail Modal** (for `GET /api/expert-requests/{id}`):
   - **Complete Expert Profile**: All fields displayed in organized sections
   - **Document Viewer**: Integrated PDF viewer for CV documents
   - **Biography Section**: Formatted text display with full content
   - **Skills Display**: Tag-style skill representation
   - **Admin Action Panel**: Approval/rejection workflow controls

4. **Admin Workflow Actions** (for `PUT /api/expert-requests/{id}`):
   - **Approve Request**:
     - Requires approval document upload (mandatory)
     - Creates expert record automatically
     - Updates request status to "approved"
   - **Reject Request**:
     - Requires rejection reason (optional but recommended)
     - Updates request status to "rejected"
     - Notifies submitter for amendment
   - **Edit Before Approval** (for `PUT /api/expert-requests/{id}/edit`):
     - Admin can modify request details
     - Then proceed to approve/reject

5. **Batch Operations** (for `POST /api/expert-requests/batch-approve`):
   - **Multi-select Functionality**: Checkbox selection
   - **Batch Approval**: Single approval document for multiple requests
   - **Batch Actions**: Bulk status changes
   - **Export Options**: CSV export of filtered results

**Status Workflow Implementation**:
- **Pending → Approved**: Expert created with EXP-ID, documents linked
- **Pending → Rejected**: Status updated, rejection reason stored
- **Rejected → Pending**: User can edit and resubmit
- **Approved → Archive**: Approved requests moved to separate view

**Real-time Features**:
- **Live Status Updates**: WebSocket or polling for status changes
- **Notification System**: Alert admins of new submissions
- **Activity Log**: Track admin actions and decisions
- **Performance Metrics**: Dashboard showing approval rates and processing times

**Styling References**:
- Data table styling from BQA's structured data displays
- Status badges using BQA color system (Green/Blue/Red/Orange)
- Modal overlay styling (`bqa_main_style.css` lines 1800-1900)
- Card layout for request details (`bqa_main_style.css` lines 2000-2200)
- Card layout from BQA service cards (`bqa_main_style.css` lines 2000-2200)
- Status indicators and badges (`bqa_main_style.css` lines 900-950)
- Modal and overlay styling (`bqa_main_style.css` lines 1800-1900)

## Backend Integration Specifications

### API Endpoints Integration

#### **Authentication Endpoints**
- `POST /api/auth/login`: User authentication with JWT token response
- Role-based routing: admin/super_user → Admin Panel, user/planner → Request Form

#### **Expert Request Endpoints**
- `POST /api/expert-requests`: Submit new expert request with form-data (CV upload)
- `GET /api/expert-requests`: Retrieve paginated requests with status filtering
- `GET /api/expert-requests/{id}`: Get specific request details for modal view
- `PUT /api/expert-requests/{id}`: Approve/reject request with approval document
- `PUT /api/expert-requests/{id}/edit`: Admin edit request before approval
- `POST /api/expert-requests/batch-approve`: Bulk approve multiple requests

#### **Supporting Data Endpoints**
- `GET /api/expert/areas`: Populate General Area dropdown options
- `GET /api/documents/{id}`: Retrieve and display uploaded documents
- `POST /api/documents`: Handle file uploads for CV and approval documents

### Form Validation Schema

#### **Expert Request Form Validation**
```typescript
interface ExpertRequestForm {
  // Personal Information (required)
  name: string;           // min 2 chars, required
  designation: string;    // min 2 chars, required  
  institution: string;    // min 2 chars, required
  phone: string;          // phone format validation, required
  email: string;          // email format validation, required
  
  // Professional Details (required)
  isBahraini: boolean;    // required checkbox
  isAvailable: boolean;   // required checkbox
  rating: string;         // required field
  role: 'evaluator' | 'validator' | 'evaluator/validator'; // required select
  employmentType: 'academic' | 'employer'; // required select
  isTrained: boolean;     // required checkbox
  isPublished?: boolean;  // optional, defaults to false
  
  // Expertise (required)
  generalArea: number;    // required, foreign key to expert_areas
  specializedArea: string; // required text field
  skills: string[];       // required array, converted to JSON
  
  // Biography & Documents (required)
  biography: string;      // required, max 1000 characters
  cv: File;              // required PDF file upload
}
```

#### **Admin Review Form Validation**
```typescript
interface ReviewForm {
  status: 'approved' | 'rejected'; // required
  rejectionReason?: string;        // optional for rejection
  approvalDocument?: File;         // required for approval (PDF)
}
```

### Data Transformation Layer

#### **Request Submission Payload**
```javascript
// Transform form data for API submission
const prepareRequestPayload = (formData) => {
  const payload = new FormData();
  
  // Personal Information
  payload.append('name', formData.name);
  payload.append('designation', formData.designation);
  payload.append('institution', formData.institution);
  payload.append('phone', formData.phone);
  payload.append('email', formData.email);
  
  // Professional Details
  payload.append('isBahraini', formData.isBahraini);
  payload.append('isAvailable', formData.isAvailable);
  payload.append('rating', formData.rating);
  payload.append('role', formData.role);
  payload.append('employmentType', formData.employmentType);
  payload.append('isTrained', formData.isTrained);
  payload.append('isPublished', formData.isPublished || false);
  
  // Expertise Areas
  payload.append('generalArea', formData.generalArea);
  payload.append('specializedArea', formData.specializedArea);
  payload.append('skills', JSON.stringify(formData.skills)); // Array to JSON
  
  // Biography & Documents
  payload.append('biography', formData.biography);
  payload.append('cv', formData.cv); // File object
  
  return payload;
};
```

#### **Admin Table Data Mapping**
```javascript
// Transform API response for table display
const transformRequestsForTable = (apiResponse) => {
  return apiResponse.data.requests.map(request => ({
    id: request.id,
    expertName: request.name,
    institution: request.institution,
    specialization: request.specializedArea,
    status: request.status,
    submittedBy: request.createdBy, // Requires user lookup
    dateSubmitted: new Date(request.createdAt).toLocaleDateString(),
    actions: ['view', 'approve', 'reject'], // Based on status
    // Additional data for modal
    fullDetails: request
  }));
};
```

### State Management Architecture

#### **Authentication State**
```typescript
interface AuthState {
  user: {
    id: number;
    name: string;
    email: string;
    role: 'super_user' | 'admin' | 'planner' | 'user';
    isActive: boolean;
  } | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
}
```

#### **Request Form State**
```typescript
interface RequestFormState {
  formData: ExpertRequestForm;
  isDraft: boolean;
  isSubmitting: boolean;
  validationErrors: Record<string, string>;
  generalAreas: Array<{id: number, name: string}>; // From API
  uploadProgress: number;
  currentStep: number; // For multi-step form
}
```

#### **Admin Panel State**
```typescript
interface AdminPanelState {
  requests: ExpertRequest[];
  selectedStatus: 'all' | 'pending' | 'approved' | 'rejected';
  filters: {
    search: string;
    institution: string;
    generalArea: string;
    dateRange: [Date, Date];
  };
  pagination: {
    page: number;
    limit: number;
    total: number;
  };
  selectedRequests: number[]; // For batch operations
  isLoading: boolean;
}
```

### Error Handling Strategy

#### **Form Submission Errors**
```javascript
const handleSubmissionError = (error) => {
  if (error.status === 400) {
    // Validation errors from backend
    const errors = error.response.data.errors;
    setValidationErrors(errors);
  } else if (error.status === 413) {
    // File too large
    setFileError('CV file size exceeds limit (5MB)');
  } else if (error.status === 415) {
    // Invalid file type
    setFileError('Please upload a PDF file');
  } else {
    // Generic server error
    setSubmissionError('Failed to submit request. Please try again.');
  }
};
```

#### **Admin Panel Error Handling**
```javascript
const handleReviewError = (error, action) => {
  if (error.status === 400 && action === 'approve') {
    setReviewError('Approval document is required');
  } else if (error.status === 404) {
    setReviewError('Request not found');
  } else {
    setReviewError(`Failed to ${action} request. Please try again.`);
  }
};
```

## Component Library Architecture

### Shared Components

1. **Layout Components**:
   - `Header`: Main navigation with user info
   - `Footer`: System footer with links
   - `Sidebar`: Admin panel navigation
   - `Container`: Content wrapper with responsive grid

2. **Form Components**:
   - `Input`: Styled text inputs following BQA patterns
   - `Textarea`: Multi-line text input with formatting
   - `Select`: Dropdown with custom styling
   - `FileUpload`: Drag-and-drop file upload component
   - `FormSection`: Grouped form elements with headers

3. **UI Components**:
   - `Button`: Various button styles (primary, secondary, success, danger)
   - `Card`: Content containers for requests and information
   - `Modal`: Overlay dialogs for detailed views
   - `StatusBadge`: Visual status indicators
   - `Pagination`: Navigation for large datasets

4. **Data Components**:
   - `RequestCard`: Expert request display component
   - `ExpertProfile`: Expert information display
   - `DocumentViewer`: PDF/document viewing component

### CSS Architecture

**Base Styles** (`src/styles/base/`):
```
base/
├── reset.css          # Normalize and reset
├── typography.css     # Font definitions and text styles
├── variables.css      # CSS custom properties for colors, spacing
└── utilities.css      # Utility classes
```

**Component Styles** (`src/styles/components/`):
```
components/
├── header.css         # Navigation and header styles
├── forms.css          # Form element styling
├── buttons.css        # Button variations
├── cards.css          # Card component styles
├── modals.css         # Modal and overlay styles
└── tables.css         # Data table styling
```

**Layout Styles** (`src/styles/layout/`):
```
layout/
├── grid.css           # Bootstrap-style grid system
├── containers.css     # Content container styles
└── responsive.css     # Media queries and responsive rules
```

## Technical Implementation Stack

### Framework & Libraries
- **React 18**: Core frontend framework
- **TypeScript**: Type safety and development experience
- **Vite**: Build tool and development server
- **React Router**: Client-side routing
- **React Hook Form**: Form management and validation
- **Tailwind CSS**: Utility-first CSS framework (customized with BQA theme)
- **Headless UI**: Unstyled, accessible UI components

### State Management
- **React Context**: Authentication and global state
- **React Query**: Server state management and caching
- **Zustand**: Local state management for complex components

### UI Enhancement
- **Framer Motion**: Animations and transitions
- **React Dropzone**: File upload functionality
- **React PDF**: Document viewing
- **React Icons**: Icon library

## Development Phases

### Phase 1: Foundation Setup (Week 1)
1. **Project Initialization**:
   - Vite React TypeScript project setup
   - Configure Tailwind CSS with BQA custom theme
   - Set up directory structure following component architecture
   - Configure ESLint, Prettier, and TypeScript with strict settings

2. **BQA Design System Implementation**:
   - Create CSS custom properties for BQA color palette
   - Implement Graphik font typography system
   - Build base component library (Button, Input, Card, Modal)
   - Establish responsive grid system matching BQA patterns

3. **Authentication Framework**:
   - JWT token management setup
   - React Context for authentication state
   - Protected route components
   - Role-based access control implementation

### Phase 2: Authentication & Core UI (Week 2)
1. **Login Page Development**:
   - Responsive login page matching BQA header design
   - Form validation with React Hook Form
   - Authentication API integration (`POST /api/auth/login`)
   - Role-based redirection logic (admin → Admin Panel, user → Request Form)

2. **Navigation Framework**:
   - Header component with user session display
   - Breadcrumb navigation system
   - Footer component with BQA styling
   - Responsive navigation for mobile/tablet

3. **Core UI Components**:
   - Status badges with BQA color coding
   - Loading states and error boundaries
   - Toast notification system
   - File upload component foundation

### Phase 3: Expert Request System (Week 3-4)
1. **Multi-Section Form Development**:
   - **Section 1**: Personal Information (name, designation, institution, phone, email)
   - **Section 2**: Professional Details (role, employment, availability, rating, training)
   - **Section 3**: Expertise Areas (general area dropdown, specialization, skills)
   - **Section 4**: Biography & Documents (biography textarea, CV upload)

2. **Form Enhancement Features**:
   - Real-time validation with field-level error display
   - General Areas API integration (`GET /api/expert/areas`)
   - Progressive form sections with step indicator
   - Auto-save draft functionality (localStorage)
   - Skills input with tag-style interface
   - Biography character counter (1000 max)

3. **File Upload System**:
   - Drag-and-drop CV upload area
   - PDF validation and file size limits
   - Upload progress indicator
   - File preview and removal functionality
   - Integration with document API (`POST /api/documents`)

4. **Form Submission**:
   - Form data transformation for API payload
   - API integration (`POST /api/expert-requests`)
   - Success/error feedback with toast notifications
   - Form reset and redirect after successful submission

### Phase 4: Admin Panel Development (Week 5-6)
1. **Data Table Implementation**:
   - Advanced data table with sorting and pagination
   - Status filter tabs (All, Pending, Approved, Rejected)
   - Search functionality across name, institution, specialization
   - Responsive table design for mobile devices
   - API integration (`GET /api/expert-requests`)

2. **Request Management Interface**:
   - Status indicator badges with BQA color system
   - Bulk selection checkboxes for batch operations
   - Filter controls (institution, area, date range)
   - Export functionality for filtered results
   - Real-time status updates

3. **Request Detail Modal**:
   - Complete expert information display
   - Integrated PDF viewer for CV documents
   - Biography section with formatted text display
   - Skills display with tag styling
   - Request timeline and submission details

4. **Admin Workflow Implementation**:
   - **Approval Workflow**: 
     - Approval document upload requirement
     - API integration (`PUT /api/expert-requests/{id}`)
     - Expert creation confirmation
   - **Rejection Workflow**:
     - Rejection reason input
     - Status update and user notification
   - **Edit Functionality**:
     - Admin edit form (`PUT /api/expert-requests/{id}/edit`)
     - Modified field highlighting

5. **Batch Operations**:
   - Multi-select request handling
   - Batch approval with single document upload
   - API integration (`POST /api/expert-requests/batch-approve`)
   - Progress tracking for bulk operations

### Phase 5: Advanced Features (Week 7)
1. **Enhanced User Experience**:
   - Optimistic UI updates for better responsiveness
   - Advanced search with autocomplete
   - Keyboard shortcuts for admin actions
   - Accessibility improvements (WCAG AA compliance)

2. **Performance Optimization**:
   - React.lazy for code splitting
   - Virtualized tables for large datasets
   - Image/document lazy loading
   - Bundle size optimization

3. **Error Handling & Monitoring**:
   - Comprehensive error boundary implementation
   - API error handling with user-friendly messages
   - Form validation error aggregation
   - Client-side logging and error reporting

### Phase 6: Testing & Deployment (Week 8)
1. **Comprehensive Testing**:
   - Unit tests for all components (Jest + React Testing Library)
   - Integration tests for form submission and admin workflows
   - E2E testing for critical user journeys (Playwright)
   - API integration testing

2. **Cross-browser & Device Testing**:
   - Desktop testing (Chrome, Firefox, Safari, Edge)
   - Mobile responsiveness testing (iOS Safari, Chrome Mobile)
   - Tablet layout testing
   - Performance testing on low-end devices

3. **Production Deployment**:
   - Production build optimization
   - Environment configuration setup
   - CI/CD pipeline configuration
   - Performance monitoring setup

## Success Metrics & KPIs

### User Experience Metrics
- **Form Completion Rate**: >95% for expert request submissions
- **Admin Processing Time**: <2 minutes average per request review
- **Error Rate**: <1% form submission failures
- **User Satisfaction**: >4.5/5 rating from admin users

### Technical Performance Metrics
- **Page Load Time**: <2 seconds on 3G connections
- **First Contentful Paint**: <1.5 seconds
- **Bundle Size**: <300KB gzipped main bundle
- **Lighthouse Score**: >90 for Performance, Accessibility, Best Practices, SEO

### Accessibility & Compliance
- **WCAG AA Compliance**: 100% compliance verified
- **Keyboard Navigation**: Full keyboard accessibility
- **Screen Reader Support**: Comprehensive ARIA labeling
- **Color Contrast**: Minimum 4.5:1 ratio across all UI elements

## Risk Mitigation Strategies

### Technical Risks
1. **File Upload Limitations**: Implement chunked upload for large files
2. **Browser Compatibility**: Progressive enhancement and polyfills
3. **Performance Issues**: Implement virtualization and lazy loading
4. **API Response Time**: Client-side caching and optimistic updates

### User Experience Risks
1. **Complex Form Abandonment**: Progressive disclosure and auto-save
2. **Admin Workflow Confusion**: Clear status indicators and guided workflows
3. **Mobile Usability**: Touch-friendly design and responsive layouts
4. **Accessibility Barriers**: Comprehensive testing and ARIA implementation

This comprehensive frontend development plan provides a structured approach to implementing a professional, government-grade application that meets the specific requirements of the Expert Database system while maintaining the visual standards and user experience quality of the BQA website.

## File Structure

```
src/
├── components/           # Reusable UI components
│   ├── common/          # Shared components (Button, Input, etc.)
│   ├── forms/           # Form-specific components
│   ├── layout/          # Layout components (Header, Footer, etc.)
│   └── admin/           # Admin-specific components
├── pages/               # Page components
│   ├── Login.tsx
│   ├── ExpertRequest.tsx
│   └── admin/
│       └── RequestReview.tsx
├── hooks/               # Custom React hooks
├── contexts/            # React context providers
├── services/            # API integration
├── types/               # TypeScript type definitions
├── utils/               # Utility functions
├── styles/              # CSS and styling
│   ├── base/
│   ├── components/
│   └── layout/
└── assets/              # Static assets (images, icons, etc.)
```

## Responsive Design Strategy

Following BQA's mobile-first approach:

1. **Mobile (< 768px)**:
   - Single-column layouts
   - Touch-friendly interactions
   - Collapsed navigation
   - Simplified forms

2. **Tablet (768px - 1024px)**:
   - Two-column layouts
   - Expanded navigation
   - Side-by-side form elements

3. **Desktop (> 1024px)**:
   - Multi-column layouts
   - Full navigation display
   - Optimized form layouts
   - Enhanced interactions

## Accessibility Considerations

Based on BQA's WCAG AA compliance:

1. **Keyboard Navigation**: Full keyboard accessibility
2. **Screen Reader Support**: Proper ARIA labels and semantic HTML
3. **Color Contrast**: Minimum 4.5:1 contrast ratios
4. **Focus Management**: Clear focus indicators
5. **Alternative Text**: Comprehensive alt text for images
6. **Form Accessibility**: Proper labeling and error messages

## Performance Goals

1. **Initial Load**: < 3 seconds on 3G
2. **First Contentful Paint**: < 1.5 seconds
3. **Lighthouse Score**: > 90 for all metrics
4. **Bundle Size**: < 300KB gzipped

## Testing Strategy

1. **Unit Testing**: Jest + React Testing Library
2. **Integration Testing**: API integration tests
3. **E2E Testing**: Playwright for critical user journeys
4. **Visual Testing**: Percy for UI regression testing
5. **Accessibility Testing**: Axe-core integration

## Deployment & CI/CD

1. **Build Process**: Vite production build with optimization
2. **Environment Configuration**: Development, staging, production
3. **CI/CD Pipeline**: GitHub Actions for automated testing and deployment
4. **Performance Monitoring**: Web Vitals tracking and monitoring

## Documentation

1. **Component Storybook**: Interactive component documentation
2. **API Documentation**: Integration guides and examples
3. **User Guides**: End-user documentation for all features
4. **Developer Documentation**: Setup and contribution guides

This comprehensive plan provides a solid foundation for developing a professional, government-grade frontend application that maintains the visual standards and user experience quality demonstrated by the BQA website while serving the specific needs of the Expert Database system.
