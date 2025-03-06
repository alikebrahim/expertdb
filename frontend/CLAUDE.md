# ExpertDB Frontend Development Guidelines

## Implementation Strategy
- Review all phases (1-6) before beginning implementation to understand dependencies and plan ahead
- Identify gaps in instructions and ask clarifying questions when critical details are missing
- Keep track of implementation decisions and rationale in this document
- Pay special attention to cross-phase concerns like authentication flow, role-based access, and data structures
- Ensure each phase smoothly integrates with both previous and upcoming phases

## Version Management

The project follows semantic versioning with phase-based milestones:
- **0.0.0**: Initial project setup (Phase 0)
- **0.1.0**: Authentication implementation (Phase 1)
- **0.2.0**: Expert Database Search (Phase 2)
- **0.3.0**: Request Submission (Phase 3)
- **0.4.0**: Statistics Dashboard (Phase 4)
- **0.5.0**: Admin Panel (Phase 5)
- **1.0.0**: Polish and Production Deployment (Phase 6)

For each completed phase:
1. Update version in package.json
2. Update IMPLEMENTATION.md with phase status
3. Update TESTING.md with testing instructions
4. Commit with conventional commit message format

## Implementation Notes

### Authentication (Phase 1)
- User state persists via localStorage with "token" key
- Role-based access: regular users -> /search, admins -> /admin
- ProtectedRoute redirects unauthenticated users to /login
- Initial shadcn/ui components not used due to import alias configuration issues
- Using standard HTML elements with Tailwind classes for now

### Future Considerations
- Plan to implement proper shadcn/ui components in Phase 2
- Error handling consistent across all requests
- API response types consistent with backend responses in ENDPOINTS.md

## Project Documentation
- [MASTER_PLAN.md](/MASTER_PLAN.md) - Main implementation plan with phase overview
- [PHASE_1.md](/PHASE_1.md) - Authentication implementation details
- [PHASE_2.md](/PHASE_2.md) - Expert Database Searching implementation
- [PHASE_3.md](/PHASE_3.md) - Request Submission implementation
- [PHASE_4.md](/PHASE_4.md) - Statistics Dashboard implementation
- [PHASE_5.md](/PHASE_5.md) - Admin Panel implementation
- [PHASE_6.md](/PHASE_6.md) - Polish and Deployment Prep
- [ENDPOINTS.md](/ENDPOINTS.md) - Backend API reference (important for integration)

## Build Commands
- `npm run dev` - Start the development server (Vite)
- `npm run build` - Build TypeScript and create production build
- `npm run lint` - Run ESLint on all files
- `npm run lint src/components/MyComponent.tsx` - Lint a specific file
- `npm run preview` - Preview the production build locally

## Code Style
- **TypeScript**: Strict mode enabled; explicit typing required
- **Imports**: Use absolute imports with `@/*` path alias (e.g., `import Component from '@/components/Component'`)
- **Components**: Use PascalCase for component filenames (e.g., `AuthContext.tsx`, `Login.tsx`)
- **Non-Components**: Use camelCase for utilities, hooks, and other files (e.g., `api.ts`, `utils.ts`)
- **Formatting**: 2-space indentation; semicolons required
- **Error Handling**: Use try/catch with proper error typing in async functions
- **State Management**: Use React Context for global state, React Query for API data

## Conventions
- **File Structure**: Components in `src/components/`, pages in `src/pages/`, utilities in `src/utils/`
- **Component Structure**: Props interface above component definition
- **API Calls**: Use axios with the configured instance from `src/api/api.ts`
- **Routing**: React Router with declared routes in App.tsx
- **Styling**: Tailwind CSS with consistent class ordering; no inline styles
- **ESM Format**: Use ESM format (`export default`) for all JS files, not CommonJS
- **Error States**: Always handle loading, error, and empty states in UI components