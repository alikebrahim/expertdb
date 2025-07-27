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