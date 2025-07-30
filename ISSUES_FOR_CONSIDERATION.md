# Issues for Consideration

This document tracks implementation issues and configuration decisions that need to be addressed during the ExpertDB development process.

## 1. Email Notification System
**Context:** All notifications will be sent via email (no in-app notifications)
**Required for:**
- Expert request status updates (approved/rejected/amendment needed)
- Assignment notifications for planners and managers
- Expert rating request notifications
- Automatic reminders for pending ratings

**Considerations:**
- Email service provider selection (SMTP, SendGrid, AWS SES, etc.)
- Email templates for different notification types
- Email queue management for reliability
- Unsubscribe/preference management
- Email delivery tracking and bounce handling

## 2. Starting Phase Number Configuration
**Context:** Phases represent academic semesters and are numbered sequentially (Phase 22, Phase 23, etc.)
**Issue:** The system will continue from existing phases that won't be imported, requiring a configurable starting number (likely Phase 23 or 24)

**Considerations:**
- Where to store the base phase number (config file, environment variable, database setting)
- How to ensure phase numbers are sequential and unique
- Migration strategy if the starting number needs to change

## 3. Automatic Reminder System
**Context:** Rating requests have a 2-week deadline and automatic reminders should be sent
**Requirements:**
- Send initial rating request notification
- Send reminder notifications (frequency to be determined)
- Track reminder status to avoid duplicate sends

**Considerations:**
- Reminder schedule (e.g., 7 days before deadline, 3 days, 1 day)
- Background job scheduling system
- Handling weekends and holidays
- Manager notification preferences
- Escalation path for overdue ratings

## 4. Honorarium PDF Generation
**Context:** Account managers need to generate honorarium reports for applications
**Requirements:**
- Create a facility for account managers to generate honorarium.pdf for applications
- Button interface to generate reports including panel details
- Calculate remuneration based on provided schedule

**Considerations:**
- PDF generation library selection
- Remuneration schedule configuration and storage
- Report template design and branding
- Panel member data aggregation
- Payment calculation logic based on roles and time commitments
- File generation and download handling

## 5. Documents and Database Backup Mechanism
**Context:** System needs automated backup capabilities for documents and database with super_user configuration
**Requirements:**
- Automated backup mechanism for SQLite database and document files
- API endpoints for super_user to configure backup schedules
- Interface for super_user to view backup files and status
- Integration with system settings in the admin panel for backup configuration
- Configurable backup retention policies

**Considerations:**
- Backup storage location (local, cloud storage integration)
- Backup frequency configuration (daily, weekly, monthly)
- Compression and encryption of backup files
- Database consistency during backup operations
- Document file synchronization with database backups
- Backup verification and integrity checks
- API design for backup configuration and monitoring
- Backup restoration procedures and testing

## 6. Expert Request Endpoint Duplication and Naming Confusion
**Context:** The system currently has two separate endpoints for updating expert requests with overlapping functionality and inconsistent access control patterns.

**Current Implementation:**
- `PUT /api/expert-requests/{id}` - Handled by `HandleUpdateExpertRequest`
- `PUT /api/expert-requests/{id}/edit` - Handled by `HandleEditExpertRequest`

**Functional Differences:**
1. **HandleUpdateExpertRequest** (PUT /api/expert-requests/{id}):
   - Primary purpose: Status management (approve/reject requests)
   - Can handle status changes: approved, rejected, pending
   - When status = "approved": Creates expert profile via `ApproveExpertRequestWithDocument()`
   - Complex approval logic with document requirements
   - Access control: Admin can edit any request, users can edit only their own rejected requests

2. **HandleEditExpertRequest** (PUT /api/expert-requests/{id}/edit):
   - Primary purpose: Data editing (modify request fields)
   - Cannot change status, focuses on updating request data
   - Simple data update without approval logic
   - Auto-resets status: If user edits rejected request → status becomes "pending"
   - Access control: Admin can edit pending requests, users can edit rejected requests

**Issues Identified:**
- Overlapping responsibilities: Both can update request data
- Inconsistent access control rules between endpoints
- Non-standard REST naming with `/edit` suffix
- Confusing separation of concerns (status vs data updates)
- Potential for data inconsistency with different update paths

**Considerations:**
- **Consolidation Option**: Merge into single `PUT /api/expert-requests/{id}` endpoint
- **Access Control Standardization**: Unified permission model across both functions
- **Clear Separation**: If keeping separate, clarify distinct responsibilities
- **API Consistency**: Follow standard RESTful patterns
- **Breaking Changes**: Impact on existing clients/documentation
- **Testing Requirements**: Ensure all functionality is preserved during consolidation

**Recommended Approach:**
Consolidate into single endpoint that handles both data updates and status changes based on request payload, with unified access control logic.

## 7. Outstanding TODO Items Requiring Resolution

**Context:** During code review, several TODO comments were identified that require implementation decisions or completion.

### 7.1 Transaction Rollback in Expert Request Creation
**Location:** `internal/api/handlers/expert_request.go:183`
**TODO Comment:** `// TODO: Consider rolling back the request creation here`

**Issue:** Expert request creation follows a two-step process:
1. Create expert request record in database
2. Upload CV document and link to request

If step 2 fails, step 1 remains committed, leaving an expert request without a CV in the database.

**Current Behavior:** Request remains in database with failed CV upload, potentially causing data inconsistency.

**Considerations:**
- **Option A**: Implement rollback by deleting the request if CV upload fails
- **Option B**: Use database transactions to ensure atomicity across both operations
- **Option C**: Accept current behavior as valid (requests can exist temporarily without CVs)
- **Option D**: Implement CV upload first, then create request with document reference

**Risk Level:** Medium - Affects data consistency and user experience

### 7.2 Rating Storage Implementation
**Location:** `internal/api/handlers/phase/phase_handler.go:806`
**TODO Comment:** `// TODO: Store rating in application_ratings table (when implemented)`

**Issue:** Expert rating functionality mentions storing ratings in `application_ratings` table but implementation status is unclear.

**Investigation Required:**
- Verify if `application_ratings` table exists in current database schema
- Check if ratings are being stored through alternative mechanism
- Determine if this represents missing functionality or outdated comment

**Current Behavior:** Rating submission appears to work but storage mechanism needs verification.

**Considerations:**
- Database schema review and potential migration needs
- Rating data model and relationships
- Historical rating data preservation
- Rating retrieval and reporting functionality

**Risk Level:** High - Critical for rating system functionality

### 7.3 Expert Names and Rating Status Retrieval
**Location:** `internal/api/handlers/phase/phase_handler.go:871`
**TODO Comment:** `// TODO: Get expert names and rating request status from database`

**Issue:** Manager tasks functionality currently returns incomplete information, missing expert names and rating request status.

**Current Behavior:** Returns basic task information without full expert details or rating status.

**Implementation Requirements:**
- Enhance database queries to join expert information
- Include expert names in response payload
- Add rating request status tracking
- Maintain backward compatibility with existing API consumers

**Considerations:**
- Database query performance with additional joins
- API response format changes and versioning
- Caching strategy for frequently accessed expert information
- Impact on mobile/frontend applications

**Risk Level:** Medium - Affects user experience and functionality completeness

### 7.4 Advanced Filtering Implementation
**Location:** `internal/api/handlers/phase/phase_handler.go:951`
**TODO Comment:** `// TODO: Implement more sophisticated filtering in storage layer`

**Issue:** Current filtering implementation is basic and could be enhanced with more advanced capabilities.

**Current Behavior:** Basic filtering functionality that may not meet all user requirements.

**Enhancement Opportunities:**
- Multi-criteria filtering with AND/OR logic
- Date range filtering for time-based queries
- Text search across multiple fields
- Sorting by multiple columns
- Pagination improvements

**Considerations:**
- Performance impact of complex queries
- Database indexing strategy
- API complexity vs. usability trade-offs
- Frontend filtering UI requirements
- Backward compatibility with existing filter parameters

**Risk Level:** Low - Enhancement rather than critical functionality

### 7.5 Resolution Priority and Timeline

**High Priority (Data Integrity Issues):**
1. Rating Storage Implementation (7.2) - Verify and complete rating persistence
2. Transaction Rollback (7.1) - Implement proper error handling and rollback

**Medium Priority (Feature Completeness):**
3. Expert Names Retrieval (7.3) - Complete API response information

**Low Priority (Enhancements):**
4. Advanced Filtering (7.4) - Can be deferred based on user feedback

**Recommended Next Steps:**
1. Database schema review to understand rating system implementation
2. Implement transaction rollback for expert request creation
3. Enhance manager tasks API with complete expert information
4. Evaluate filtering requirements against current implementation