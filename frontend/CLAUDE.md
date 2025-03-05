# ExpertDB Frontend Guidelines

## Build Commands
```bash
# Install dependencies
npm install

# Run development server
npm run dev

# Build for production
npm run build

# Run production build locally
npm run start

# Run linting
npm run lint

# Type checking
npm run typecheck
```

## Code Style Guidelines
- **TypeScript**: Use proper type annotations for all components and functions
- **Components**: Use functional components with hooks
- **State Management**: Use React Context API for global state, React Query for data fetching
- **Naming**: Use PascalCase for components, camelCase for variables and functions
- **Imports**: Group imports by type (React, libraries, components, utilities)
- **CSS**: Use Tailwind CSS for styling, avoid inline styles
- **Error Handling**: Implement proper error boundaries and fallbacks
- **Loading States**: Always provide loading indicators for asynchronous operations

## UI/UX Design Guidelines

### Brand Colors
Follow the BQA (Bahrain Education & Training Quality Authority) color palette:
- **Primary Blue**: #133566 (navy blue) - Used for headers, primary buttons, and key UI elements
- **Secondary Blue**: #1B4882 (light blue) - Used for secondary elements, hover states
- **Green**: #192012 - Used as an accent color for success states or subtle effects
- **Orange**: #DC8335 - Used for warnings, call-to-action elements, or subtle accents
- **Red**: #FF4040 - Used for errors, alerts, or critical information

### Typography
- **Primary Font**: Graphik for English text
- **Headings**: Use appropriate heading levels (h1-h6) for semantic structure
- **Text Sizes**:
  - Headings: 1.5rem - 3rem
  - Body text: 1rem
  - Small text: 0.875rem
- **Font Weights**:
  - Headings: 600-700
  - Body: 400
  - Emphasized text: 500-600

### Layout Guidelines
- Maintain consistent padding and margins throughout the application
- Use a grid-based layout for alignment
- Ensure responsive design with appropriate breakpoints:
  - Mobile: < 640px
  - Tablet: 640px - 1024px
  - Desktop: > 1024px
- Use white space effectively to create visual hierarchy

### Component Design
- Use the shadcn/ui component library for consistent UI elements
- Implement accessible designs with proper ARIA attributes
- Follow these button styles:
  - Primary: Solid navy blue (#133566) with white text
  - Secondary: Light blue (#1B4882) or outline style
  - Danger: Red (#FF4040) for destructive actions
- Form inputs should have:
  - Clear labels
  - Visible focus states
  - Appropriate error messages
  - Consistent sizing

### Icons and Images
- Use SVG icons for scalability
- BQA logos are available in the guidelines/logo directory
- Maintain proper spacing around logos according to brand guidelines
- Use appropriate image formats and optimize for web

### Animation and Effects
- Use subtle animations for state changes and transitions
- Avoid excessive animation that could distract users
- Consider reduced motion preferences for accessibility

### User Experience Guidelines
- Provide clear feedback for all user actions
- Implement consistent navigation patterns
- Use breadcrumbs for deep navigation structures
- Create intuitive form flows with proper validation
- Design with accessibility in mind (keyboard navigation, screen readers)

## Frontend Architecture

### Directory Structure
- `/app`: Next.js App Router pages and layouts
- `/components`: Reusable UI components
  - `/layout`: Layout components (navbar, footer, etc.)
  - `/ui`: UI components from shadcn
- `/lib`: Utility functions and API client
- `/public`: Static assets

### State Management
- **Authentication**: React Context for user authentication state
- **Data Fetching**: React Query for API data with caching
- **Forms**: React Hook Form for form state and validation
- **UI State**: Local state with useState or useReducer

### API Integration
- Use the API client in `/lib/api.ts` for all backend communication
- Handle loading states and errors consistently
- Implement proper authentication header management
- Use React Query for data fetching, caching, and synchronization

## Implementing UI/UX Design
The application design should follow the BQA guidelines as demonstrated in the sample implementation artifacts found in the guidelines directory. Key aspects include:

1. **Brand Consistency**: Use the BQA color palette and typography
2. **Professional Appearance**: Create a clean, professional interface suitable for an internal system
3. **Intuitive Navigation**: Design clear navigation paths for different user roles
4. **Responsive Design**: Ensure the application works well on all screen sizes
5. **Accessibility**: Make the application accessible to all users

## Logo Usage Guidelines
The BQA logos are available in the frontend/guidelines/logo directory:
- `BQA - Horizontal Logo.svg`: Use in the navbar and header sections
- `BQA - Horizontal Logo with Descriptor.svg`: Use on the login page (first point of contact for users)
- `BQA - Vertical Logo - With descriptor.svg`: Use in the footer or sidebar
- `Icon Logo - Color.svg`: Use for favicon and smaller UI elements

Please follow appropriate brand guidelines for logo usage, ensuring proper spacing and visibility.

## Current Implementation Status
See the [Frontend Implementation Status](/frontend/IMPLEMENTATION.md) document for details on what has been implemented and what is still pending.