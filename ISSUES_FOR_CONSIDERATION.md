# Issues for Consideration

This document tracks implementation issues and configuration decisions that need to be addressed during the ExpertDB development process.

## Overview

This document is organized into the following categories:
1. **Critical Issues** - Data integrity and security concerns requiring immediate attention
2. **Missing Implementations** - Features referenced but not implemented
3. **Architectural Concerns** - Design patterns and structural issues
4. **Configuration Requirements** - System settings and deployment needs
5. **Code Quality Issues** - Maintainability and consistency concerns
6. **Enhancement Opportunities** - Non-critical improvements

---

## CRITICAL ISSUES

### 1. Expert Request Endpoint Duplication (HIGH PRIORITY)
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
- Empty form submission protection exists but is a symptom of the duplication issue

**Recommended Approach:**
Consolidate into single endpoint that handles both data updates and status changes based on request payload, with unified access control logic.

### 2. Transaction Rollback in Expert Request Creation (MEDIUM PRIORITY)
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

### 3. SQL Injection Vulnerabilities (MEDIUM PRIORITY)
**Location:** Multiple locations in storage layer

**Issues Found:**
1. `internal/storage/sqlite/expert.go:1024` - Dynamic SQL with `fmt.Sprintf`:
   ```go
   query := fmt.Sprintf("UPDATE %s SET %s = ? WHERE id = ?", table, column)
   ```

2. `internal/storage/sqlite/expert_request.go:768` - Dynamic SQL for IN clause:
   ```go
   query := fmt.Sprintf(`UPDATE experts SET approval_document_id = ? WHERE id IN (%s)`,
       strings.Join(placeholders, ","))
   ```

**Risk:** While parameters are properly bound, table/column names in dynamic SQL could be exploited if not properly validated.

**Recommendation:** 
- Validate table/column names against whitelist before use
- Consider using prepared statements with fixed SQL
- Add input validation layer for all dynamic SQL components

---

## MISSING IMPLEMENTATIONS

### 4. Application Ratings Table (HIGH PRIORITY)
**Location:** `internal/api/handlers/phase/phase_handler.go:806`
**TODO Comment:** `// TODO: Store rating in application_ratings table (when implemented)`

**Issue:** Expert rating functionality mentions storing ratings in `application_ratings` table but implementation status is unclear.

**Investigation Results:**
- Database schema review confirms `application_ratings` table does NOT exist
- Rating API endpoints exist but don't persist data
- This represents missing core functionality

**Required Implementation:**
- Create database migration for `application_ratings` table
- Implement storage methods for rating CRUD operations
- Update API handlers to use new storage methods
- Add rating retrieval endpoints for reporting

**Risk Level:** High - Critical for rating system functionality

### 5. Expert Names and Rating Status Retrieval (MEDIUM PRIORITY)
**Location:** `internal/api/handlers/phase/phase_handler.go:871`
**TODO Comment:** `// TODO: Get expert names and rating request status from database`

**Issue:** Manager tasks functionality currently returns incomplete information, missing expert names and rating request status.

**Current Behavior:** Returns basic task information without full expert details or rating status.

**Implementation Requirements:**
- Enhance database queries to join expert information
- Include expert names in response payload
- Add rating request status tracking
- Maintain backward compatibility with existing API consumers

**Risk Level:** Medium - Affects user experience and functionality completeness

---

## ARCHITECTURAL CONCERNS

### 6. Storage Interface Anti-patterns
**Issue:** Multiple architectural issues in the storage layer

**Problems Identified:**
1. **Large Interface Violation**: Single `Storage` interface with 50+ methods violates Interface Segregation Principle
2. **Type Safety**: Methods use `interface{}` for documentService parameter instead of proper interface
3. **Tight Coupling**: Direct SQLite dependencies throughout storage layer

**Examples:**
```go
ApproveExpertRequestWithDocument(requestID, reviewedBy int64, documentService interface{}) (int64, error)
BatchApproveExpertRequestsWithFileMove(..., documentService interface{}) ([]int64, []int64, map[int64]error)
```

**Recommendations:**
- Split Storage interface into domain-specific interfaces (UserStorage, ExpertStorage, etc.)
- Define proper DocumentService interface instead of using interface{}
- Implement repository pattern with database abstraction layer

### 7. Dead Code and Empty Files
**Files Identified:**
- `internal/auth/jwt.go` - Empty file with consolidation comment
- `internal/auth/password.go` - Empty file

**Issue:** These files create confusion and clutter in the codebase

**Recommendation:** Remove empty files and update any references

---

## CONFIGURATION REQUIREMENTS

### 8. Email Notification System
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

### 12. Security Enhancements Required
**Missing Security Features:**
1. **Authentication Security:**
   - No rate limiting on login endpoints
   - No account lockout after failed attempts
   - JWT tokens valid for 24 hours with no refresh mechanism
   - No multi-factor authentication support

2. **Audit Logging:**
   - No comprehensive audit trail for sensitive operations
   - Missing logging for data modifications
   - No compliance-ready audit reports

3. **API Security:**
   - No API rate limiting
   - Missing request validation middleware
   - No API versioning strategy

**Recommendations:**
- Implement rate limiting using middleware
- Add account lockout mechanism after N failed attempts
- Implement JWT refresh token flow
- Add comprehensive audit logging for all data modifications
- Consider API gateway for advanced security features

---

## CODE QUALITY ISSUES

### 13. Inconsistent Error Handling
**Issues Identified:**
1. Mix of domain errors and string errors throughout codebase
2. Different error formats returned by various endpoints
3. Error messages sometimes expose internal implementation details
4. Inconsistent use of error wrapping and context

**Examples:**
- Some endpoints use `domain.ErrNotFound`, others return custom error strings
- Database errors sometimes exposed directly to API responses
- Validation errors have inconsistent structure

**Recommendations:**
- Standardize on domain error types for all common scenarios
- Implement consistent error wrapping with context
- Create error response middleware for uniform API errors
- Never expose internal error details in production

### 14. Logging Inconsistencies
**Issues Identified:**
1. Inconsistent log levels across the application
2. Some critical operations lack logging entirely
3. Debug logs sometimes contain sensitive information (user IDs, request data)
4. No structured logging format for analysis

**Recommendations:**
- Define and document log level guidelines
- Implement structured logging (JSON format)
- Add request ID tracking for correlation
- Sanitize sensitive data from all log output
- Ensure all critical operations are logged

### 15. Documentation and Implementation Gaps
**Issues Found:**
1. **SRS Misalignment**: Several features marked "implemented" aren't complete
2. **API Documentation**: Endpoints documented but not all parameters explained
3. **Role System**: Documentation doesn't match implementation complexity
4. **Statistics**: Calculations may not match documented requirements

**Recommendations:**
- Audit all "implemented" features against actual code
- Update documentation to reflect current implementation
- Add OpenAPI/Swagger documentation generation
- Create developer onboarding documentation

---

## ENHANCEMENT OPPORTUNITIES

### 16. Advanced Filtering Implementation
**Location:** `internal/api/handlers/phase/phase_handler.go:951`
**TODO Comment:** `// TODO: Implement more sophisticated filtering in storage layer`

**Current State:** Basic filtering functionality that may not meet all user requirements.

**Enhancement Opportunities:**
- Multi-criteria filtering with AND/OR logic
- Date range filtering for time-based queries
- Full-text search across multiple fields
- Sorting by multiple columns
- Saved filter preferences
- Export filtered results

**Risk Level:** Low - Enhancement rather than critical functionality

### 17. Performance Optimizations
**Areas for Improvement:**
1. **Database Queries**: Some endpoints make multiple queries that could be joined
2. **Caching**: No caching layer for frequently accessed data
3. **Pagination**: Default limits may be too high for large datasets
4. **File Handling**: Document uploads could benefit from streaming

**Recommendations:**
- Implement query optimization with proper joins
- Add Redis caching for hot data paths
- Review and optimize pagination defaults
- Implement streaming for large file operations

---

## RESOLUTION PRIORITY MATRIX

### Critical Priority (Immediate Action Required):
1. **Expert Request Endpoint Duplication** - Data integrity risk
2. **Application Ratings Table Missing** - Core functionality incomplete
3. **SQL Injection Vulnerabilities** - Security risk

### High Priority (Next Sprint):
4. **Transaction Rollback Implementation** - Data consistency
5. **Email Notification System** - Required functionality
6. **Security Enhancements** - Authentication and audit logging

### Medium Priority (Next Quarter):
7. **Storage Interface Refactoring** - Technical debt
8. **Expert Names Retrieval** - Feature completeness
9. **Error Handling Standardization** - Code quality
10. **Starting Phase Configuration** - Deployment requirement

### Low Priority (Future Enhancements):
11. **Advanced Filtering** - Nice to have
12. **Performance Optimizations** - Proactive improvement
13. **Dead Code Cleanup** - Housekeeping
14. **Documentation Updates** - Ongoing task

---

## RECOMMENDED ACTION PLAN

### Phase 1 - Critical Fixes (Week 1-2):
1. Consolidate expert request endpoints into single REST endpoint
2. Create and implement `application_ratings` table
3. Fix SQL injection vulnerabilities with proper validation

### Phase 2 - Core Features (Week 3-4):
1. Implement transaction support for file operations
2. Set up email notification infrastructure
3. Add rate limiting and account lockout

### Phase 3 - Architecture (Week 5-6):
1. Refactor storage interface into smaller interfaces
2. Implement proper audit logging
3. Standardize error handling

### Phase 4 - Polish (Week 7-8):
1. Complete all TODO implementations
2. Update documentation
3. Performance testing and optimization
4. Security audit and penetration testing