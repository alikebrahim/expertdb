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