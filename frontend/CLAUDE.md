# ExpertDB Frontend Development Guidelines

## Build Commands
- `npm run dev` - Start development server with HMR
- `npm run build` - Compile TypeScript and build for production
- `npm run lint` - Run ESLint on codebase
- `npm run preview` - Preview production build locally

## Code Style Guidelines
- **TypeScript**: Use strict typing with no implicit any
- **Components**: Functional components with React hooks
- **Imports**: Group imports (React, external libs, internal components, styles)
- **Naming**:
  - PascalCase for components and interfaces
  - camelCase for variables, functions, and instances
- **Component Structure**:
  - Props interface above component declaration
  - State declarations at top of component
  - Helper functions before return statement
- **Error Handling**: Try/catch blocks with appropriate error messages
- **Styling**: Use shadcn/ui components with consistent design patterns

## Tech Stack
- React 18 with TypeScript
- Vite for building and development
- react-router-dom for routing
- react-hook-form for form handling
- shadcn/ui for component library
- axios for API requests

## Implementation Tracking
- /.system/IMPLEMENTATION.md - Tracks overall progress and outstanding tasks
- /.system/ISSUE_LOG.md - Tracks issues and their resolutions
- issues/ - Directory to document historical issues and their solutions
- Each implementation step is tagged with a status: TODO, IN_PROGRESS, COMPLETED