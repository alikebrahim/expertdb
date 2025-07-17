# ExpertDB Frontend Implementation Plan

## Overview

We need to investigate and plan the frontend work to ensure it matches the design criteria as described in @bqa_website/. The frontend must support three core workflows with clear user experiences. The design will focus on the colorscheme, menus, and implementing a professional, government-grade application that maintains the visual standards and user experience quality of the BQA website.

## Design System Foundation

### Visual Identity & Branding (from BQA Analysis)
- **Primary Colors**:
  - BQA Green: `#397b26` (Primary brand color, for approvals and positive actions)
  - Deep Blue: `#1c4679` (Navigation, headers, and authority elements)
  - Accent Red: `#e64125` (Call-to-action, warnings, and rejections)
  - Orange: `#e68835` (Secondary actions and highlights)

- **Typography System**:
  - Primary Font: "Graphik" family (Regular 400, Medium 500, Semibold 600, Bold 700)
  - Fallback: "Segoe UI", Helvetica Neue, Arial, sans-serif
  - Hierarchy: h1(42px), h2(38px), h3(35px), h4(26px), body(16px)

- **Layout Framework**:
  - Bootstrap-based responsive grid system
  - Container max-width: 1212px
  - Responsive breakpoints: xs(<576px), sm(≥576px), md(≥768px), lg(≥992px), xl(≥1200px)

## Core Workflows

### A. Expert Database Browsing Interface

#### 1. Main Expert Search & Browse Interface (All User Roles)

**Navigation:** Upon Login → Automatic redirect to `/search` → Main Expert Database Interface

**Primary Interface Features:**

The expert database browsing interface serves as the main entry point and core functionality of the application. All users, regardless of role, are automatically directed here upon login to access the comprehensive expert database.

**Advanced Search & Filter Interface:**

```
┌─────────────────────────────────────────┐
│              Navigation                 │
├─────────────────────────────────────────┤
│           Expert Search Header          │
│      "Search and filter experts"        │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │        Progress Stepper             │ │
│  │ [1] Define filters → [2] View       │ │
│  │ results → [3] Contact experts       │ │
│  └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│                                         │
│        Advanced Filter Panel           │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │     Basic Search Section            │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │Name/Institution Search          │ │ │
│  │  │[🔍 Search experts...]           │ │ │
│  │  └─────────────────────────────────┘ │ │
│  │                                     │ │
│  │  ┌───────────┐ ┌───────────────┐   │ │
│  │  │Role       │ │Employment Type│   │ │
│  │  │[Dropdown] │ │[Dropdown]     │   │ │
│  │  └───────────┘ └───────────────┘   │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │     Advanced Filters (Expandable)  │ │
│  │  ┌───────────┐ ┌───────────────┐   │ │
│  │  │Nationality│ │Expert Area    │   │ │
│  │  │[Dropdown] │ │[Dynamic API]  │   │ │
│  │  └───────────┘ └───────────────┘   │ │
│  │  ┌───────────┐ ┌───────────────┐   │ │
│  │  │☐ Available│ │☐ Bahraini     │   │ │
│  │  │☐ Published│ │Rating [Min ⬇] │   │ │
│  │  └───────────┘ └───────────────┘   │ │
│  │                                     │ │
│  │  [Clear All Filters] [Apply Filters]│ │
│  └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│                                         │
│           Results Summary               │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │ Found 127 experts (All loaded)      │ │
│  │ Filter results: 127 of 127 total    │ │
│  │                    [Contact Selected]│ │
│  └─────────────────────────────────────┘ │
│                                         │
├─────────────────────────────────────────┤
│                                         │
│         Expert Results Table            │
│                                         │
│  ┌─────┬──────────┬─────────────┬──────┬────────┬──────────┬─────────────┬──────┐ │
│  │Name │Institution│Specialized │Rating│Role    │Employment│General Area │Action│ │
│  │     │          │Area         │      │        │Type      │             │      │ │
│  ├─────┼──────────┼─────────────┼──────┼────────┼──────────┼─────────────┼──────┤ │
│  │Dr.  │University│Civil        │★★★★☆│Eval/   │Academic  │Engineering  │[View]│ │
│  │Ahmed│of Bahrain│Engineering  │ 4.2  │Valid   │          │             │      │ │
│  ├─────┼──────────┼─────────────┼──────┼────────┼──────────┼─────────────┼──────┤ │
│  │Sara │Ministry  │Healthcare   │★★★★★│Eval    │Employer  │Health       │[View]│ │
│  │Hassan│of Health │Management   │ 4.8  │        │          │Sciences     │      │ │
│  └─────┴──────────┴─────────────┴──────┴────────┴──────────┴─────────────┴──────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │         Table Controls              │ │
│  │ [Export] [Column Settings] [Sort]   │ │
│  │ Showing all 127 experts             │ │
│  └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│              Footer                     │
└─────────────────────────────────────────┘
```

**Filter System Specifications:**

**Basic Filters (Always Visible):**
- **Name/Institution Search**: Real-time text search across expert names and affiliations
- **Role Filter**: Dropdown with options (Evaluator, Validator, Consultant, Trainer, Expert)
- **Employment Type**: Dropdown with options (Academic, Employer, Freelance, Government, Other)

**Advanced Filters (Expandable Section):**
- **Nationality**: Dropdown with country options
- **Expert Area**: Dynamic dropdown populated from expert_areas API endpoint
- **Availability Status**: Checkbox for available experts only
- **Bahraini Status**: Checkbox for Bahraini nationals only
- **Published Status**: Checkbox for experts who consent to publication
- **Rating Filter**: Minimum rating selector (1-5 stars)

**Filter State Management:**
- **Persistent Filters**: Saved to localStorage and restored on page load
- **Default State**: Available experts checkbox enabled by default
- **Filter Indicators**: Visual badges showing active filters
- **Quick Clear**: One-click filter reset functionality

#### 2. Expert Results Display System

**Advanced Data Table Features:**

**Column Configuration:**

**Standard Fields (Always Displayed):**
- **Expert ID**: Auto-generated EXP-#### identifier
- **Name**: Full name with sortable functionality
- **Institution/Affiliation**: Current workplace or organization
- **Specialized Area**: Primary area of expertise
- **Rating**: Visual star rating display with numerical value
- **Role**: Evaluator, Validator, etc.
- **Employment Type**: Academic, Employer, etc.
- **General Area**: Broad category of expertise
- **Actions**: View Profile, Contact Info buttons

**Optional Fields (User Selectable):**
- **Availability Status**: Visual indicator (Available/Busy)
- **Nationality**: Expert's nationality
- **Training Status**: Completion status
- **Date Added**: When expert was added to database

**Sorting Capabilities:**
- **Multi-field Sorting**: Name, Institution, Rating, Date Added, Specialization
- **Sort Direction**: Ascending/Descending with visual indicators
- **Sort Persistence**: Saved to localStorage for user convenience
- **Real-time Feedback**: Toast notifications for sort changes

**Data Display System:**
- **Full Database Load**: All expert entries loaded by default without pagination (single page table)
- **Client-side Operations**: Sorting, filtering, and searching performed on loaded dataset
- **Performance Optimized**: Efficient in-memory operations for up to 2000 expert records
- **Optional Column Selection**: Users can choose which fields to display to reduce table clutter

**Export Functionality:**
- **CSV Export**: Complete expert data with selected filters applied
- **Selective Export**: Checkbox selection for specific experts
- **Custom Fields**: Choose which columns to include in export
- **Progress Tracking**: Visual feedback for large export operations

#### 3. Expert Profile Detail View

**Navigation:** Expert Table → "View" Button → Individual Expert Profile (`/experts/:id`)

**Profile Interface Layout:**

```
┌─────────────────────────────────────────┐
│              Navigation                 │
│           [← Back to Search]            │
├─────────────────────────────────────────┤
│                                         │
│           Expert Profile Header         │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │     Expert ID: EXP-1247             │ │
│  │  ┌─────────┐                        │ │
│  │  │[Photo]  │  Dr. Ahmed Ali          │ │
│  │  │Placeholder│  Senior Engineer      │ │
│  │  │         │  University of Bahrain  │ │
│  │  │         │  ★★★★☆ 4.2/5.0         │ │
│  │  └─────────┘                        │ │
│  │     🟢 Available    📧 Contact      │ │
│  └─────────────────────────────────────┘ │
│                                         │
├─────────────────────────────────────────┤
│                                         │
│         Tabbed Information Panel        │
│                                         │
│  [Personal] [Expertise] [Biography]     │
│  [Documents] [Engagement History]       │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │        Personal Information         │ │
│  │  Email: ahmed.ali@uob.edu.bh        │ │
│  │  Phone: +973 1234 5678              │ │
│  │  Nationality: Bahraini              │ │
│  │  Role: Evaluator/Validator          │ │
│  │  Employment: Academic               │ │
│  │  Training Status: ✅ Completed      │ │
│  │  Publication Consent: ✅ Yes        │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │           Expertise Areas           │ │
│  │  General Area: Engineering          │ │
│  │  Specialized: Civil Engineering     │ │
│  │  Skills: [Structural] [Design]      │ │
│  │          [Project Management]       │ │
│  │          [Quality Assurance]        │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │           Biography                 │ │
│  │  Experience:                        │ │
│  │  [2020-Present] Senior Engineer,    │ │
│  │  Ministry of Works, Bahrain         │ │
│  │  [2018-2020] Project Engineer,      │ │
│  │  Gulf Construction Co, Bahrain      │ │
│  │                                     │ │
│  │  Education:                         │ │
│  │  [2015-2018] PhD Civil Engineering, │ │
│  │  University of Bahrain              │ │
│  │  [2011-2015] MSc Structural Eng,    │ │
│  │  King Fahd University, Saudi Arabia │ │
│  └─────────────────────────────────────┘ │
│                                         │
├─────────────────────────────────────────┤
│        [Edit Profile] [Download CV]     │
│         (Admin Only)     (All Users)    │
└─────────────────────────────────────────┘
```

**Profile Data Display:**
- **Personal Information Tab**: Contact details, nationality, role, employment type
- **Expertise Tab**: Specialized areas, skills tags, general classification
- **Biography Tab**: Formatted experience and education timeline
- **Documents Tab**: CV download, approval documents (if available)
- **Engagement History Tab**: Past assignments, ratings, project participation

**Role-Based Actions:**
- **All Users**: View profile, download CV, access contact information (same privileges for browsing functionality)
- **Contact Integration**: Direct email/phone contact capabilities

#### 4. Search Process Workflow

**Progressive Search Interface:**

**Step 1: Define Filters**
- **Filter Configuration**: Set search criteria using basic and advanced filters
- **Real-time Validation**: Immediate feedback on filter combinations
- **Saved Searches**: Option to save frequently used filter sets
- **Quick Filters**: Predefined filter combinations (Available Experts, Bahraini Experts, etc.)

**Step 2: View Results**
- **Results Display**: Comprehensive table with sorting and pagination
- **Relevance Scoring**: Results ordered by relevance to search criteria
- **Result Actions**: View profiles, select for contact, export data
- **Refinement Options**: Modify filters based on initial results

**Step 3: Contact Experts**
- **Contact Information**: Access to email and phone details
- **Batch Contact**: Select multiple experts for group communication
- **Contact History**: Track previous communications (Admin only)
- **Export Contact List**: Generate contact sheets for selected experts

#### 5. Data Import & Migration Planning

**CSV Data Integration Preparation:**

While the application already supports CSV export functionality through the backup system, the new phase includes planning for importing existing expert data from CSV format into the SQL database.

**Import Process Specifications:**
- **Data Mapping**: Map CSV columns to database fields automatically
- **Validation System**: Comprehensive data validation before import
- **Conflict Resolution**: Handle duplicate entries and data conflicts
- **Progress Tracking**: Real-time import progress with error reporting
- **Rollback Capability**: Ability to revert imports if issues occur

**Data Migration Workflow:**
1. **Pre-Import Analysis**: Analyze CSV structure and data quality
2. **Field Mapping**: Map CSV columns to database schema
3. **Validation Phase**: Check data integrity and format compliance
4. **Import Execution**: Batch import with progress monitoring
5. **Post-Import Verification**: Validate imported data accuracy

**Quality Assurance:**
- **Data Integrity Checks**: Ensure all required fields are populated
- **Duplicate Detection**: Identify and merge duplicate expert entries
- **Reference Validation**: Verify foreign key relationships (expert areas, etc.)
- **Error Reporting**: Comprehensive log of import issues and resolutions

#### 6. Performance & User Experience

**Search Performance Optimization:**
- **Indexed Search**: Database indexing for name, institution, and specialization fields
- **Caching Strategy**: Client-side caching of expert areas and filter options
- **Lazy Loading**: Progressive data loading for large result sets
- **Debounced Search**: Optimized real-time search with request throttling

**User Experience Features:**
- **Responsive Design**: Mobile-optimized interface for tablet and phone access
- **Keyboard Navigation**: Full keyboard accessibility for power users
- **Search Shortcuts**: Quick filter hotkeys and keyboard shortcuts
- **State Persistence**: Maintain search state across sessions

**Accessibility Compliance:**
- **ARIA Labels**: Comprehensive screen reader support
- **Color Contrast**: High contrast mode for visually impaired users
- **Focus Management**: Clear focus indicators and logical tab order
- **Voice Navigation**: Support for voice commands and dictation

### B. New Expert Creation Workflow

#### 1. Expert Request Submission (User Role)

**Navigation:** User Dashboard → "Expert Requests" menu item → Expert Requests List Page → "New Expert Request" button → Expert Request Form

**Expert Requests List Page Features:**
- Display all user's expert requests with status (pending/approved/rejected/archived)
- Status indicators using BQA color system (🟡 Pending, 🔵 Approved, 🔴 Rejected, 📁 Archived)
- "New Expert Request" button to create new request
- View request details and admin feedback
- Search functionality across request details

**Multi-Section Form Structure:**

```
┌─────────────────────────────────────────┐
│              Navigation                 │
├─────────────────────────────────────────┤
│           Page Header Section           │
│     "Submit Expert Request"             │
│    ████████████████████████████████     │
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
│  │  │Nationality│ │☐ Bahraini     │   │ │
│  │  │[Dropdown] │ │☐ Available    │   │ │
│  │  └───────────┘ └───────────────┘   │ │
│  │  ┌───────────┐ ┌───────────────┐   │ │
│  │  │Role*      │ │Employment*    │   │ │
│  │  │☐Evaluator │ │☐Academic      │   │ │
│  │  │☐Validator │ │☐Employer      │   │ │
│  │  └───────────┘ └───────────────┘   │ │
│  │  ┌───────────┐ ┌───────────────┐   │ │
│  │  │☐ Trained  │ │☐ Published    │   │ │
│  │  └───────────┘ └───────────────┘   │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │    Section 3: Expertise Areas      │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │General Area* [Dropdown]         │ │
│  │  │(Populated from API)             │ │
│  │  └─────────────────────────────────┘ │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │Specialized Area*                │ │
│  │  └─────────────────────────────────┘ │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │Skills* [Tag Input Interface]    │ │
│  │  │[Skill 1] [Skill 2] [+ Add]      │ │
│  │  └─────────────────────────────────┘ │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │   Section 4: Biography & Documents │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │Experience* [Add Experience +]   │ │ │
│  │  │┌─ Experience Entry 1 ─────────┐│ │ │
│  │  ││ From: [YYYY-MM] To: [YYYY-MM]││ │ │
│  │  ││ [DateFrom - DateTo] Role,    ││ │ │
│  │  ││ Organization, Location       ││ │ │
│  │  ││ [Remove Entry]               ││ │ │
│  │  │└─────────────────────────────┘│ │ │
│  │  └─────────────────────────────────┘ │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │Education* [Add Education +]     │ │ │
│  │  │┌─ Education Entry 1 ──────────┐│ │ │
│  │  ││ From: [YYYY-MM] To: [YYYY-MM]││ │ │
│  │  ││ [DateFrom - DateTo] Degree,  ││ │ │
│  │  ││ Institution, Location        ││ │ │
│  │  ││ [Remove Entry]               ││ │ │
│  │  │└─────────────────────────────┘│ │ │
│  │  └─────────────────────────────────┘ │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │Biography Preview                │ │ │
│  │  │[Real-time formatted display]    │ │ │
│  │  └─────────────────────────────────┘ │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │CV Upload* (PDF, max 20MB)       │ │ │
│  │  │[Drag & Drop Area]               │ │ │
│  │  │📎 Drop files here or click      │ │ │
│  │  │   [Browse Files]                │ │ │
│  │  └─────────────────────────────────┘ │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │        Form State Management       │ │
│  │  💾 Auto-saved 2 minutes ago       │ │
│  │  [Save Draft] [Submit for Review]  │ │
│  └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│              Footer                     │
└─────────────────────────────────────────┘
```

**Form Fields (Complete Schema):**

**Section 1: Personal Information**
- Full Name (required) - text input with validation
- Email (required) - email format validation
- Phone Number (required) - phone format validation
- Designation (professional title) - text input
- Institution/Affiliation - text input
- Nationality (dropdown: Bahraini/Non-Bahraini/Unknown)
- Is Bahraini (checkbox) - boolean flag

**Section 2: Professional Details**
- Role (checkboxes: Evaluator/Validator - can select both)
- Employment Type (checkboxes: Academic/Employer - can select both)
- Is Available (checkbox) - availability status
- Is Trained (checkbox) - training completion status
- Is Published (checkbox) - publication preference

**Section 3: Specialization Areas**
- General Area (dropdown from expert_areas table via API)
- Specialized Area (text field) - specific expertise
- Skills (tag-style input) - array of competencies converted to JSON

**Section 4: Biography & Documents**
- Biography Template Editor (Structured Format):
  - Education Section (multiple entries):
    - Date Range: From (YYYY-MM) - To (YYYY-MM)
    - Description: "[DateFrom - DateTo] Degree, Institution, Location/country[optional]"
    - Add/Remove functionality for multiple entries
  - Experience Section (multiple entries):
    - Date Range: From (YYYY-MM) - To (YYYY-MM)
    - Description: "[DateFrom - DateTo] Role/Position, Organization, Location/country[optional]"
    - Add/Remove functionality for multiple entries
  - Real-time biography preview showing formatted output
  - No rich text formatting, plain text only
- Document Upload:
  - CV (mandatory, PDF only, max 20MB)
  - Drag-and-drop interface with fallback file picker
  - Real-time file validation and upload progress
  - File preview functionality

**Advanced Form Features:**
- **Progressive Disclosure**: Collapsible/expandable sections for better UX
- **Real-time Validation**: Field-level error display with helpful messaging
- **Auto-save Functionality**: Every 10 seconds with "Draft saved" indicator
- **Form State Management**: Comprehensive state tracking and error handling
- **Responsive Design**: Mobile-optimized layouts with touch-friendly interactions

**User Actions:**
- Save as Draft:
  - Auto-save every 10 seconds (with "Draft saved" indicator)
  - Manual "Save Draft" button
  - Drafts stored as expert_requests with 'pending' status
- Submit for Review - full form validation before submission
- View My Requests (navigate back to Expert Requests List)

#### 2. Admin Review Process

**Navigation:** Admin Dashboard → "Pending Expert Requests" → Review Queue

**Advanced Data Table Interface:**

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
│  │ [🔍 Search] [Institution ⬇] [Area ⬇]│ │
│  │ [Clear Filters] [Sort: Date ⬇]      │ │
│  └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│                                         │
│             Advanced Data Table         │
│                                         │
│ ┌─────┬──────────┬─────────────┬────────┐ │
│ │ ☐ID │ Name     │ Institution │ Status │ │
│ ├─────┼──────────┼─────────────┼────────┤ │
│ │☐001 │ Dr. Ali  │ University  │ 🟡     │ │
│ │     │ Ahmed    │ of Bahrain  │ Pending│ │
│ │     │          │ 2024-01-15  │ [View] │ │
│ ├─────┼──────────┼─────────────┼────────┤ │
│ │☐002 │ Sarah    │ Ministry    │ 🔵     │ │
│ │     │ Hassan   │ of Health   │ Apprvd │ │
│ │     │          │ 2024-01-14  │ [View] │ │
│ ├─────┼──────────┼─────────────┼────────┤ │
│ │☐003 │ Ahmed    │ BCCI        │ 🔴     │ │
│ │     │ Al-Said  │             │ Reject │ │
│ │     │          │ 2024-01-13  │ [View] │ │
│ └─────┴──────────┴─────────────┴────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │        Batch Operations             │ │
│  │ [Batch Approve] [Batch Reject]      │ │
│  │ [Export CSV] [← Prev] [Next →]      │ │
│  └─────────────────────────────────────┘ │
├─────────────────────────────────────────┤
│              Footer                     │
└─────────────────────────────────────────┘
```

**Review Interface Features:**
- **Status Filter Tabs**: Filter by pending/approved/rejected/archived status
- **Advanced Search**: Real-time search across name, institution, specialization
- **Sortable Columns**: Name, Institution, Date, Status with visual indicators
- **Bulk Selection**: Checkboxes for batch operations
- **Pagination**: Standard pagination with page size options
- **Export Options**: CSV export of filtered results

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
│  │  Experience:                        │ │
│  │  [2020-Present] Senior Engineer,    │ │
│  │  Ministry of Works, Bahrain         │ │
│  │                                     │ │
│  │  Education:                         │ │
│  │  [2015-2018] PhD Civil Engineering, │ │
│  │  University of Bahrain              │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │           Documents                 │ │
│  │  📄 CV Document                     │ │
│  │     [View PDF] [Download]           │ │
│  │     Size: 2.3MB | Uploaded: Jan 15  │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │         Admin Actions               │ │
│  │  Approval Document: [Choose File]   │ │
│  │  📎 Upload PDF approval document    │ │
│  │                                     │ │
│  │  [✅ Approve] [❌ Reject] [✏️ Edit]  │ │
│  │                                     │ │
│  │  Amendment Comments (if needed):    │ │
│  │  [Text area for feedback]           │ │
│  └─────────────────────────────────────┘ │
│                                         │
├─────────────────────────────────────────┤
│        [Cancel] [Save Changes]          │
└─────────────────────────────────────────┘
```

**Admin Actions:**

**✅ Approve → Create expert profile automatically**
- **Single Approval**: Upload approval document (PDF) dialog
- **Batch Approval**: 
  - Select multiple requests via checkboxes
  - "Batch Approve" button → Upload single approval document
  - Dialog shows: "Apply this approval document to X selected requests"
  - Single document applied to all selected expert profiles
  - Progress tracking for bulk operations

**🔄 Request Amendment → Return to user with mandatory comments**
- Text area for detailed feedback
- Email notification sent to user automatically
- Status updated to "amendment_requested"

**📁 Archive → Mark as archived (kept in database, no further action)**
- Can be restored later via "Restore" action
- Archived items moved to separate view

**Advanced Admin Features:**
- **Integrated PDF Viewer**: In-modal document viewing capability
- **Skills Display**: Tag-style skill representation
- **Real-time Status Updates**: Live status indicators
- **Activity Log**: Track admin actions and decisions
- **Keyboard Shortcuts**: Efficient navigation for power users

### C. Phase Planning Workflow

#### 1. Phase Management (Admin Role)

**Navigation:** Admin Dashboard → "Phase Management" → Phase List

**Phase Management Interface:**

```
┌─────────────────────────────────────────┐
│           Phase Management              │
├─────────────────────────────────────────┤
│  ┌─────────────────────────────────────┐ │
│  │         Active Phases               │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │ Phase 23 (Active)               │ │ │
│  │  │ 📅 Jan 2024 - Jun 2024         │ │ │
│  │  │ 📊 Applications: 25/30          │ │ │
│  │  │ [View Details] [Edit]           │ │ │
│  │  └─────────────────────────────────┘ │ │
│  │                                     │ │
│  │  ┌─────────────────────────────────┐ │ │
│  │  │ Phase 24 (Draft)                │ │ │
│  │  │ 📅 Jul 2024 - Dec 2024         │ │ │
│  │  │ 📊 Applications: 0/30           │ │ │
│  │  │ [View Details] [Edit]           │ │ │
│  │  └─────────────────────────────────┘ │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  [+ Create New Phase]                   │
└─────────────────────────────────────────┘
```

**Phase Creation/Edit Form:**
- **Phase Details**:
  - Phase Name: "Phase [number]" (e.g., Phase 22, Phase 23)
  - Date Range with overlap validation (no overlaps allowed)
  - Status: Draft → Active → Closed (lifecycle management)
  - Description (optional) - detailed phase information
  - Target number of applications

#### 2. Application Management (Admin Role)

**Within Phase Context:** Phase Detail Page → "Applications" Tab

**Application Management Interface:**

```
┌─────────────────────────────────────────┐
│              Phase 23 Details           │
│           Applications Management       │
├─────────────────────────────────────────┤
│  ┌─────────────────────────────────────┐ │
│  │      Application Status Overview    │ │
│  │  🟢 Approved: 15  🟡 Pending: 8     │ │
│  │  ⬜ Unassigned: 2  📊 Total: 25     │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │        Applications List            │ │
│  │ ┌─────┬─────────────┬─────────────┐ │ │
│  │ │Type │ Title       │ Status      │ │ │
│  │ ├─────┼─────────────┼─────────────┤ │ │
│  │ │ IL  │ UoB Campus  │ 🟢 Approved │ │ │
│  │ │     │ Engineering │ Expert: Ali │ │ │
│  │ │     │ Program     │ Ahmed       │ │ │
│  │ ├─────┼─────────────┼─────────────┤ │ │
│  │ │ QP  │ Advanced    │ 🟡 Pending  │ │ │
│  │ │     │ Engineering │ Planner:    │ │ │
│  │ │     │ Diploma     │ Sarah Hassan│ │ │
│  │ ├─────┼─────────────┼─────────────┤ │ │
│  │ │ IL  │ BCCI Trade  │ ⬜ Unassigned│ │ │
│  │ │     │ Center      │ [Assign]    │ │ │
│  │ └─────┴─────────────┴─────────────┘ │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  [+ Create New Application]             │
└─────────────────────────────────────────┘
```

**Application Creation Form:**
- **Type Selection**: 
  - IL (Institutional Listing) - single institution title field
  - QP (Qualification Placement) - qualification title + institution title (two fields)
- **Title Fields** (Dynamic based on type):
  - IL Type: Institution Title only
  - QP Type: Qualification Title + Institution Title
- **Sector** (Dropdown): Higher Education (HE), Vocational Education and Training (VET), General Education (GE)
- **Assignments** (Optional at creation):
  - Planner dropdown (lists all active users)
  - Manager dropdown (lists all active users)
  - Note: Same user can be assigned as both planner and manager
- **Expert Requirements**:
  - Expert_1 slot (mandatory)
  - Expert_2 slot (optional)

**Note:** No bulk operations for application creation - single application creation only

#### 3. Assignment Management (Admin Role)

**Flexible Assignment Interface:**
- **During Creation**: Optional assignment via dropdowns
- **After Creation**: Single application edit with reassignment capability
- **Reassignment**: Allowed at any time with audit trail
- **No Constraints**: Flexible assignment rules allowing same user multiple roles

#### 4. Expert Proposal Process (Planner Role)

**Navigation:** Planner Dashboard → "My Assigned Applications"

**Planner Workspace:**

```
┌─────────────────────────────────────────┐
│          Planner Dashboard              │
│        My Assigned Applications         │
├─────────────────────────────────────────┤
│  ┌─────────────────────────────────────┐ │
│  │      Application Queue              │ │
│  │ ┌─────┬─────────────┬─────────────┐ │ │
│  │ │Phase│ Application │ Status      │ │ │
│  │ ├─────┼─────────────┼─────────────┤ │ │
│  │ │ 23  │ Engineering │ 🟡 Pending  │ │ │
│  │ │     │ Diploma     │ Expert Req  │ │ │
│  │ │     │ - UoB       │ [Select]    │ │ │
│  │ ├─────┼─────────────┼─────────────┤ │ │
│  │ │ 23  │ Trade Center│ 🟢 Submitted│ │ │
│  │ │     │ - BCCI      │ Expert: Ali │ │ │
│  │ │     │             │ Ahmed       │ │ │
│  │ └─────┴─────────────┴─────────────┘ │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  [Batch Submit Selected]                │
└─────────────────────────────────────────┘
```

**Expert Selection Interface:**

```
┌─────────────────────────────────────────┐
│        Expert Selection Modal          │
│      Engineering Diploma - UoB         │
├─────────────────────────────────────────┤
│  ┌─────────────────────────────────────┐ │
│  │      Search & Filter Experts       │ │
│  │  [🔍 Search] [Area ⬇] [Available ⬇]│ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │         Expert Database             │ │
│  │ ┌─────┬─────────────┬─────────────┐ │ │
│  │ │     │ Name        │ Details     │ │ │
│  │ ├─────┼─────────────┼─────────────┤ │ │
│  │ │ ◯   │ Dr. Ali     │ Civil Eng   │ │ │
│  │ │     │ Ahmed       │ Rating: 4.5 │ │ │
│  │ │     │             │ ⚠️ Also in   │ │ │
│  │ │     │             │ App #15     │ │ │
│  │ ├─────┼─────────────┼─────────────┤ │ │
│  │ │ ◯   │ Sarah       │ Mechanical  │ │ │
│  │ │     │ Hassan      │ Rating: 4.2 │ │ │
│  │ │     │             │ Available   │ │ │
│  │ └─────┴─────────────┴─────────────┘ │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │        Selection Summary            │ │
│  │  Expert_1 (Required): [Select]      │ │
│  │  Expert_2 (Optional): [Select]      │ │
│  │                                     │ │
│  │  [Auto-Save Draft] [Submit]         │ │
│  └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

**Proposal Interface Features:**
- **Expert Database Search**: Advanced search/filter functionality
- **Expert Selection**: Select Expert_1 (mandatory) and Expert_2 (optional)
- **Auto-save**: Expert selections saved immediately (no explicit draft status)
- **Double-booking Warning**: Visual indicator when expert assigned to multiple applications in same phase (but allows assignment)
- **Batch Operations**: Submit single proposal or batch submit multiple proposals

#### 5. Proposal Review (Admin Role)

**Navigation:** Admin Dashboard → Phase Management → Select Phase → Applications List

**Enhanced Review Interface:**

```
┌─────────────────────────────────────────┐
│           Phase 23 - Applications       │
│            Expert Proposals Review      │
├─────────────────────────────────────────┤
│  ┌─────────────────────────────────────┐ │
│  │      Color-Coded Application List   │ │
│  │ ┌─────┬─────────────┬─────────────┐ │ │
│  │ │Type │ Title       │ Status      │ │ │
│  │ ├─────┼─────────────┼─────────────┤ │ │
│  │ │ IL  │ UoB Campus  │ 🟢 APPROVED │ │ │
│  │ │     │ Engineering │ Expert: Ali │ │ │
│  │ │     │ Program     │ Ahmed + 1   │ │ │
│  │ ├─────┼─────────────┼─────────────┤ │ │
│  │ │ QP  │ Advanced    │ 🟡 PENDING  │ │ │
│  │ │     │ Engineering │ REVIEW      │ │ │
│  │ │     │ Diploma     │ [Review]    │ │ │
│  │ ├─────┼─────────────┼─────────────┤ │ │
│  │ │ IL  │ BCCI Trade  │ ⬜ AWAITING │ │ │
│  │ │     │ Center      │ ASSIGNMENT  │ │ │
│  │ │     │             │ [Assign]    │ │ │
│  │ └─────┴─────────────┴─────────────┘ │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │           Legend                    │ │
│  │  🟢 Green: Approved applications    │ │
│  │  🟡 Yellow: Awaiting review         │ │
│  │  ⬜ No highlight: Pending assignment│ │
│  └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

**Review Interface Features:**
- **Color-coded Status System**:
  - 🟡 Yellow highlight: Expert proposals submitted, awaiting review
  - 🟢 Green highlight: Approved applications
  - ⬜ No highlight: Pending expert assignment
- **Quick Actions**: Approve/reject buttons per application
- **Detailed Review**: Modal with expert profiles and justifications

**Admin Actions:**
- **✅ Approve** → Experts assigned to application
  - Can approve Expert_1 only (partial approval allowed)
  - Can approve both experts
- **🔄 Reject** → Return to planner with mandatory comments
  - No revision limit - unlimited resubmissions allowed
- **Real-time Updates**: Status changes reflected immediately

### D. Expert Rating Workflow

#### 1. Rating Request Initiation (Admin Role)

**Navigation:** Admin Dashboard → Phase Management → Select Phase → Applications List → Select Applications

**Rating Request Interface:**

```
┌─────────────────────────────────────────┐
│          Rating Request Management      │
│               Phase 23                  │
├─────────────────────────────────────────┤
│  ┌─────────────────────────────────────┐ │
│  │      Application Selection          │ │
│  │ ┌─────┬─────────────┬─────────────┐ │ │
│  │ │ ☐   │ Application │ Rating      │ │ │
│  │ ├─────┼─────────────┼─────────────┤ │ │
│  │ │ ☐   │ UoB Campus  │ ✅ Rated    │ │ │
│  │ │     │ Engineering │ Dr. Ahmed   │ │ │
│  │ ├─────┼─────────────┼─────────────┤ │ │
│  │ │ ☐   │ Advanced    │ ⏳ Pending  │ │ │
│  │ │     │ Engineering │ Need Rating │ │ │
│  │ │     │ Diploma     │ Due: 7 days │ │ │
│  │ ├─────┼─────────────┼─────────────┤ │ │
│  │ │ ☐   │ BCCI Trade  │ ❌ Overdue  │ │ │
│  │ │     │ Center      │ 3 days late │ │ │
│  │ └─────┴─────────────┴─────────────┘ │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │        Rating Request Options       │ │
│  │  📅 Deadline: 2 weeks from today    │ │
│  │  📧 Auto-reminders: 7, 3, 1 days    │ │
│  │                                     │ │
│  │  [Send Rating Request]              │ │
│  └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

**Selection Interface Features:**
- **Multi-select Applications**: Checkbox selection for rating requests
- **Rating Status Indicators**: Visual indicators showing which applications have ratings
- **Deadline Management**: Automatic 2-week deadline set from request date
- **Progress Tracking**: Visual indication of rating completion status

#### 2. Rating Submission (Manager Role)

**Navigation:** Manager Dashboard → "Rating Requests" → Rating Form

**Rating Form Interface:**

```
┌─────────────────────────────────────────┐
│            Expert Rating Form           │
│              Phase 23                   │
├─────────────────────────────────────────┤
│  ┌─────────────────────────────────────┐ │
│  │      Application Summary            │ │
│  │  Phase: Phase 23                    │ │
│  │  Type: QP - Engineering Program     │ │
│  │  Institution: University of Bahrain │ │
│  │  Sector: Higher Education           │ │
│  │  Due: January 30, 2024              │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │        Expert Information           │ │
│  │  Name: Dr. Ali Ahmed                │ │
│  │  Specialization: Civil Engineering  │ │
│  │  Previous Rating: 4.2/5.0           │ │
│  │  Last Assignment: Phase 22          │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │         Rating Criteria             │ │
│  │  [GAP: What's the rating scale?     │ │
│  │   1-5, 1-10, or custom?]           │ │
│  │                                     │ │
│  │  [GAP: Are there multiple rating    │ │
│  │   categories?]                      │ │
│  │                                     │ │
│  │  [GAP: Is there a comments field?   │ │
│  │   Is it mandatory?]                 │ │
│  │                                     │ │
│  │  [GAP: Are there different criteria │ │
│  │   for IL vs QP?]                    │ │
│  └─────────────────────────────────────┘ │
│                                         │
│  ┌─────────────────────────────────────┐ │
│  │        Submission Options           │ │
│  │  [Submit Rating]                    │ │
│  │  Note: Ratings cannot be edited     │ │
│  │  after submission                   │ │
│  └─────────────────────────────────────┘ │
└─────────────────────────────────────────┘
```

**Rating Form Features:**
- **Application Context**: Full application details for informed rating
- **Expert Profile**: Complete expert information for reference
- **Rating Input**: Professional rating interface (specifics needed)
- **Immutable Submissions**: No editing after submission
- **No Draft Functionality**: Direct submission only

#### 3. Rating Management

**System Behavior:**
- **Storage**: Ratings stored against expert profile
- **Integration**: Linked to application and phase context
- **Audit Trail**: Complete history of rating submissions

**[GAP] Areas Requiring Specification:**
- **Rating Display**: How are ratings displayed? Average? Distribution?
- **Privacy Settings**: Who can view ratings? Role-based access?
- **Rating History**: Is there rating history/versioning?
- **Search Integration**: Are ratings used for expert search/filtering?

## Technical Implementation Specifications

### API Integration Requirements

#### Authentication & Authorization
- **JWT Token Management**: Secure token storage and refresh
- **Role-based Access Control**: Super_user > Admin > Planner > User
- **Session Management**: Automatic logout and session timeout

#### Form Validation Schema
```typescript
interface ExpertRequestForm {
  // Personal Information
  name: string;           // min 2 chars, required
  designation: string;    // min 2 chars, required  
  institution: string;    // min 2 chars, required
  phone: string;          // phone format validation
  email: string;          // email format validation
  
  // Professional Details
  isBahraini: boolean;    // required
  isAvailable: boolean;   // required
  role: 'evaluator' | 'validator' | 'evaluator/validator';
  employmentType: 'academic' | 'employer';
  isTrained: boolean;     // required
  isPublished?: boolean;  // optional, defaults to false
  
  // Expertise Areas
  generalArea: number;    // foreign key to expert_areas
  specializedArea: string; // required
  skills: string[];       // required array, converted to JSON
  
  // Biography & Documents
  biography: Biography;   // structured biography
  cv: File;              // required PDF file
}
```

#### State Management Architecture
```typescript
interface RequestFormState {
  formData: ExpertRequestForm;
  isDraft: boolean;
  isSubmitting: boolean;
  validationErrors: Record<string, string>;
  generalAreas: Array<{id: number, name: string}>;
  uploadProgress: number;
  currentStep: number;
}
```

### Performance & Accessibility Standards

#### Performance Goals
- **Initial Load**: < 3 seconds on 3G connections
- **First Contentful Paint**: < 1.5 seconds
- **Bundle Size**: < 300KB gzipped main bundle
- **Lighthouse Score**: > 90 for Performance, Accessibility, Best Practices, SEO

#### Accessibility Compliance
- **WCAG AA Compliance**: 100% compliance verified
- **Keyboard Navigation**: Full keyboard accessibility
- **Screen Reader Support**: Comprehensive ARIA labeling
- **Color Contrast**: Minimum 4.5:1 ratio across all UI elements
- **Focus Management**: Clear focus indicators and logical tab order

#### Responsive Design Strategy
- **Mobile (< 768px)**: Single-column layouts, touch-friendly interactions
- **Tablet (768px - 1024px)**: Two-column layouts, expanded navigation
- **Desktop (> 1024px)**: Multi-column layouts, full feature set

### Error Handling & User Feedback

#### Comprehensive Error Management
```javascript
const handleSubmissionError = (error) => {
  if (error.status === 400) {
    const errors = error.response.data.errors;
    setValidationErrors(errors);
    showToast('Please correct the highlighted fields', 'error');
  } else if (error.status === 413) {
    setFileError('CV file size exceeds limit (20MB)');
    showToast('File too large. Please upload a smaller file.', 'error');
  } else if (error.status === 415) {
    setFileError('Please upload a PDF file');
    showToast('Invalid file type. Only PDF files are allowed.', 'error');
  } else {
    setSubmissionError('Failed to submit request. Please try again.');
    showToast('Submission failed. Please try again later.', 'error');
  }
};
```

#### Toast Notification System
- **Success Notifications**: Green background with checkmark icon
- **Error Notifications**: Red background with error icon  
- **Warning Notifications**: Orange background with warning icon
- **Info Notifications**: Blue background with info icon
- **Auto-dismiss**: Configurable timeout with manual dismiss option

## Global UI/UX Specifications

### Navigation & Dashboard Design

#### Role-Based Dashboard Content

**Admin Dashboard:**
- **Statistics Panel**: Pending requests, monthly approvals, system health
- **Recent Activity**: Latest submissions, phase updates, rating reminders
- **Quick Actions**: Review requests, create phase, manage experts, view reports
- **Performance Metrics**: Processing times, approval rates, user activity

**User Dashboard:**
- **My Requests**: Status overview of submitted expert requests
- **Quick Actions**: Create new request, view request history
- **Help & Support**: User guides, contact information, FAQ

**Planner Dashboard:**
- **Assigned Applications**: Queue of applications requiring expert proposals
- **Submission Status**: Track proposal submissions and approvals
- **Quick Actions**: Submit proposals, view application details

**Manager Dashboard:**
- **Rating Requests**: List of pending expert ratings
- **Rating History**: Previously submitted ratings
- **Performance Analytics**: Rating patterns and trends

#### Global Search Functionality
- **Real-time Search**: Instant results across all entities
- **Search Categories**: Experts, requests, phases, applications
- **Advanced Filters**: Date ranges, status, institution, area
- **Keyboard Shortcuts**: Ctrl+K to open search, Esc to close
- **Search History**: Recent searches with quick access

### Notifications System

#### Email Notifications (All notifications via email only)
- **Expert Request Updates**: Approval, rejection, amendment requests
- **Assignment Notifications**: New planner/manager assignments
- **Rating Requests**: New rating assignments with deadline reminders
- **Automatic Reminders**: 7, 3, 1 days before rating deadline
- **System Notifications**: Phase status changes, system maintenance

#### Notification Preferences
- **Email Frequency**: Immediate, daily digest, weekly summary
- **Notification Types**: All, critical only, custom selection
- **Delivery Times**: Configurable quiet hours and timezone

### Data Management & Security

#### Data Retention & Export
- **Expert Data**: Permanent retention with archival after 5 years
- **Request History**: Complete audit trail maintained
- **Rating Data**: Permanent retention with versioning
- **Export Capabilities**: CSV, PDF reports for all data types

#### Security Features
- **Role-based Permissions**: Granular access control
- **Audit Logging**: Complete activity tracking
- **Data Encryption**: All sensitive data encrypted at rest
- **Session Security**: Automatic timeout and secure token handling

### Performance Optimization

#### Loading States & Feedback
- **Skeleton Loading**: Placeholder content while data loads
- **Progress Indicators**: Upload progress, form submission status
- **Optimistic Updates**: Immediate UI feedback with server confirmation
- **Error Recovery**: Graceful handling of network failures

#### Caching Strategy
- **Expert Areas**: Cached for 1 hour, refreshed on demand
- **User Profiles**: Cached for session duration
- **Request Data**: Real-time updates with optimistic caching
- **Static Assets**: Long-term caching with versioning

## Implementation Priorities

### Phase 1: Foundation & Authentication
1. **BQA Design System Implementation**: Colors, typography, layout framework
2. **Authentication System**: Login, role-based access, session management
3. **Navigation Framework**: Headers, menus, breadcrumbs
4. **Core Components**: Buttons, forms, modals, tables

### Phase 2: Expert Database Browsing Interface
1. **Advanced Search & Filter System**: Real-time search, expandable filters, persistent state
2. **Expert Results Display**: Sortable tables, pagination, export functionality
3. **Expert Profile Detail Views**: Comprehensive profile pages with tabbed information
4. **CSV Data Import Planning**: Data migration tools and validation systems

### Phase 3: Expert Request System
1. **Multi-section Form**: Progressive disclosure, validation, auto-save
2. **File Upload System**: Drag-and-drop, validation, preview
3. **Biography Editor**: Structured input, real-time preview
4. **Admin Review Interface**: Data tables, modals, batch operations

### Phase 4: Phase Planning System
1. **Phase Management**: Creation, editing, status tracking
2. **Application Management**: Creation, assignment, tracking
3. **Expert Selection**: Search, filtering, conflict detection
4. **Proposal Review**: Color-coded interface, approval workflow

### Phase 5: Rating System & Advanced Features
1. **Rating Request Management**: Selection, deadline tracking
2. **Rating Submission**: Professional rating interface
3. **Rating Analytics**: Display, history, reporting
4. **Performance Optimization**: Code splitting, caching, monitoring

## Related Documentation

See ISSUES_FOR_CONSIDERATION.md for:
- Email notification system implementation details
- Starting phase number configuration
- Automatic reminder system for ratings

## Success Metrics

### User Experience Metrics
- **Expert Search Response Time**: <1 second for filtered search results
- **Form Completion Rate**: >95% for expert request submissions
- **Admin Processing Time**: <2 minutes average per request review
- **Error Rate**: <1% form submission failures
- **User Satisfaction**: >4.5/5 rating from admin users
- **Search Success Rate**: >90% of searches result in expert profile views

### Technical Performance Metrics
- **Page Load Time**: <2 seconds on 3G connections
- **First Contentful Paint**: <1.5 seconds
- **Bundle Size**: <300KB gzipped main bundle
- **Lighthouse Score**: >90 for Performance, Accessibility, Best Practices, SEO

This comprehensive implementation plan provides detailed specifications for developing a professional, government-grade frontend application that maintains BQA design standards while delivering modern user experience and technical excellence.