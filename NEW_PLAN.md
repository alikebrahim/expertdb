# Marvelito Experiment Plan: Structured Context-Building for ExpertDB with Claude

## Overview
This plan guides Claude to work on *ExpertDB*, a web app at 60-70% completion, using a multi-file configuration to build structured context. The focus is on enabling Claude to understand the codebase, generate key files, and manage tasks effectively without explicit commands—all actions are Claude’s responsibility to determine and execute as needed.

**Project Context**:
- **Type**: Web app for managing an expert database
- **Tech Stack**: 
  - Backend: Go-based REST API, SQLite database, JWT authentication
  - Frontend: React with TypeScript, Vite, shadcn/ui components
  - Other: Docker, SQLite migrations
- **Main Features**: Expert database management, request submission/review, document management, advanced search, stats/reporting, user management, role-based access control
- **Progress**: 
  - Done: Backend APIs with auth, database schema, role-based security, CSV import, frontend auth/admin UI, Docker config
  - Remaining: Login page styling/functionality, expert search with filtering, expert bio page, request page, stats dashboard, user management UI, auth/role testing, integration tests, deployment
- **Pain Points**: 
  - Frontend-backend integration issues (e.g., users/experts not fetched/displayed)
  - Recent auth persistence fixes needing validation
  - Database migration management
  - Codebase reorganization (files marked for deletion/newly added)
- **Purpose**: Use structured context to improve Claude’s ability to manage ExpertDB’s larger codebase, starting with login page completion

**Experiment Goal**: Test how structured context-building—mapping endpoints, function signatures, and recording issues—helps Claude work on ExpertDB in a directed, focused manner.

---

## Instructions for Claude

This document is for you, Claude, to ingest and act upon. Your task is to:
1. Analyze the ExpertDB Go/React codebase.
2. Generate the files listed below in a `.system/` directory to reflect the project’s current state.
3. Use these files to build context and work on the project, starting with completing the login page styling and functionality.
4. Perform actions (e.g., coding, debugging, documenting) as you see fit, guided by the files and project state.

---

## Files to Generate

Generate these files in `.system/` to map ExpertDB’s current state and provide structured context. Use them to understand, document, and extend the codebase.

### 1. `CLAUDE.md` (UPDATE since it already exist)
- **Location**: Project root (`./CLAUDE.md`)
- **Purpose**: Your entry point to understand ExpertDB and guide your actions.
- **Contents**: 
  - Persona: “A meticulous craftsman refining ExpertDB’s Go backend and React frontend integration.”
  - Guidelines: “Analyze the codebase, prioritize frontend-backend integration, document gaps, and complete unfinished features like the login page. Use other `.system/` files for details.”
  - File Links: List all `.system/` files (e.g., `.system/ENDPOINTS.md`).
- **Use Case**: Start here to orient yourself and proceed with tasks.

### 2. Reference Documents
- **`ENDPOINTS.md`**:
  - **Purpose**: Map all Go API endpoints to reflect the current backend state.
  - **Contents**: List endpoints (e.g., `/experts GET` – returns expert list), request/response schemas (JSON), and status (e.g., “working” or “fetch failing”).
  - **Use Case**: Use this to fix integration issues (e.g., why experts aren’t displaying in the frontend).
- **`UI_UX_GUIDELINES.md`**:
  - **Purpose**: Outline React UI components and their current state.
  - **Contents**: List components (e.g., “Login: shadcn/ui button, needs styling”), expected behavior (e.g., “Login submits JWT”), and issues (e.g., “No spinner”).
  - **Use Case**: Guide your work on login page styling and other UI tasks.
- **`FUNCTION_SIGNATURES.md`**:
  - **Purpose**: Index Go/React functions for quick reference.
  - **Contents**: Function names, args, return types, locations (e.g., `fetchExperts() -> Expert[], file: api.ts, line: 20`).
  - **Use Case**: Speed up analysis of integration bugs (e.g., tracing `fetchExperts`).
- **`AUTH_GUIDELINES.md`**:
  - **Purpose**: Detail JWT auth and role-based logic.
  - **Contents**: Auth flow, role rules (e.g., “Admin views all users”), issues (e.g., “Persistence fixed?”).
  - **Use Case**: Validate auth fixes and test role access.

### 3. Progress Tracker
- **`IMPLEMENTATION.md`**:
  - **Purpose**: Track ExpertDB’s stages and current 60-70% state.
  - **Contents**: 
    - Stages: “Backend APIs + Auth (done), Admin/Auth UI (done), Login Styling + Integration (in progress), Search/Bio/Requests/Stats/User Mgmt (todo).”
    - Details: “Login API: POST /login (done), Frontend fetch failing; Search: filters pending.”
  - **Use Case**: Prioritize login completion and plan next steps.

### 4. Log Documents
- **`MASTER_RECORD.md`**:
  - **Purpose**: Summarize your actions for a project timeline.
  - **Contents**: Timestamps, actions, outcomes (e.g., “2025-03-12: Styled login button”).
  - **Use Case**: Track what you’ve done to maintain context.
- **`ISSUE_LOG.md`**:
  - **Purpose**: Record bugs and fixes consistently.
  - **Contents**: Issue, steps, status (e.g., “Bug: Experts not fetched, Fix: Updated api.ts, Status: closed”).
  - **Use Case**: Address integration/persistence issues and avoid repeating fixes.

---

## Workflow

1. **Setup**: Ingest this plan and the ExpertDB codebase.
2. **File Generation**: Create all listed files in `.system/` by analyzing the Go/React code to reflect the current state.
3. **Context Building**: Use the files to understand ExpertDB’s structure, progress, and issues.
4. **First Task**: Complete the login page styling (shadcn/ui) and functionality (fix fetch issues), updating relevant files (e.g., `UI_UX_GUIDELINES.md`, `ISSUE_LOG.md`).
5. **Ongoing Work**: Proceed with tasks like fixing expert fetch or documenting endpoints as needed, keeping files updated.

---

## Expectations

- **Your Role**: Generate files, build context, and work autonomously. Errors (e.g., misreading Go structs) may occur—I’ll refine your outputs.
- **Maintenance**: You generate and update files as you work; I’ll check accuracy periodically.
- **Context Goal**: Structured files (endpoints, signatures, logs) should help you manage ExpertDB more effectively than a single prompt.

---

## Experiment Goal

Test how mapping endpoints, function signatures, and recording issues in structured files helps you stay focused and directed on ExpertDB, starting with the login page.

---

## Notes
- **Customization**: Files reflect Go/React/SQLite specifics.
- **Outcome**: Success if you style the login page, fix a fetch bug, or document endpoints usefully.
- **Support**: I’ll adjust post-experiment based on your results.
