# Frontend-Backend Integration Plan (Revised)

This document outlines the adjusted phased plan for integrating the frontend and backend of the ExpertDB application, focusing on implementing the core functionalities and developing the UI for them.

## Current Status Assessment

### Backend
- REST API built with Go
- Modular architecture with domain, storage, service, and API layers
- Authentication with JWT
- Endpoints available for experts, expert requests, documents, engagements, statistics, and phases
- Role-based access control implemented

### Frontend
- React/TypeScript application with Tailwind CSS
- API client and type definitions in place
- Basic components and pages implemented
- Some integration with backend APIs already exists

## Integration Priorities

Considering the project's identified characteristics:
- Small scale application (10-12 users, ~2000 expert entries)
- Focus on simplicity and maintainability
- Modest performance requirements (responses under 2 seconds)
- Internal use only with organizational security
- Distinct user roles with clear workflow responsibilities

## Refactoring Needs

Before implementing the main functionalities, the frontend needs focused refactoring:

1. **Authentication & Authorization Foundation**
   - Implement complete JWT authentication flow
   - Create role-based component visibility system
   - Add visual indicators for current user role

2. **Component Structure Simplification**
   - Standardize component folder structure
   - Implement consistent prop and state patterns
   - Create simpler reusable components for common UI patterns

3. **API Integration Standardization**
   - Ensure all API calls follow the same pattern
   - Implement proper error handling and loading states
   - Use appropriate local storage caching for frequently accessed data

## Stage 1: Search & Filter Foundation

**Goal**: Create a focused search interface optimized for internal users with essential filtering and sorting.

### Tasks:

1. **Backend Integration**
   - Implement core filters for expert search API
   - Add basic sorting parameters (name, role, institution)
   - Ensure pagination works correctly

2. **Frontend Implementation**
   - Build efficient `ExpertFilters` component with only essential filters
   - Implement information-dense table layout for expert results
   - Add basic localStorage caching for recent searches
   - Implement role-specific viewable columns
   - Focus on keyboard navigation for power users

3. **Performance Optimization**
   - Implement debouncing for search inputs
   - Add loading indicators for search operations
   - Ensure response times remain under 2 seconds

**Estimated Duration**: 1 week

**Success Criteria**:
- Search results return within 1-2 seconds
- All roles can effectively filter and find experts
- Interface prioritizes information density over aesthetics
- Basic role-based permissions restrict data visibility

## Stage 2: Expert Creation Workflow

**Goal**: Implement the complete expert creation workflow from request to profile creation, with focus on document management.

### Tasks:

1. **Admin Review Interface**
   - Build admin dashboard for reviewing expert requests
   - Create status indicators for pending/approved/rejected requests
   - Implement batch operations for administrators
   - Add document preview capabilities

2. **Expert Request Submission**
   - Create streamlined request form with CV upload
   - Implement mandatory field validation
   - Add progress indicators for submission steps

3. **Document Management**
   - Implement secure document upload for CVs
   - Create document preview functionality
   - Add approval document attachment for admins
   - Implement document categorization

4. **Profile Creation**
   - Build profile creation form from approved requests
   - Add specialization area assignment
   - Implement rating and availability settings

**Estimated Duration**: 2 weeks

**Success Criteria**:
- Complete workflow from request to expert profile creation
- Document uploads and previews work correctly
- Admin review process is efficient and intuitive
- Notifications for status changes work properly

## Stage 3: Phase Planning Workflow

**Goal**: Implement the phase planning workflow with focus on assignment and approval processes.

### Tasks:

1. **Phase Management**
   - Create phase creation and management interface for admins
   - Implement assignment to planners
   - Build timeline visualization for phase planning

2. **Expert Assignment**
   - Develop expert selection interface optimized for planners
   - Implement efficient expert filtering within assignment context
   - Create conflict detection for expert assignments
   - Add batch assignment capabilities

3. **Approval Workflow**
   - Build admin approval dashboard for proposed experts
   - Implement approve/reject with comments functionality
   - Create engagement tracking upon approval

**Estimated Duration**: 2 weeks

**Success Criteria**:
- Phase creation and management works efficiently for admins
- Planners can easily assign appropriate experts
- Approval workflow maintains data integrity
- Engagements are correctly created upon approval

## Stage 4: Statistics & Reporting

**Goal**: Develop focused statistics and reporting capabilities tailored for internal decision-making.

### Tasks:

1. **Core Statistics Dashboard**
   - Implement essential charts for expert distribution
   - Create nationality and role distribution visuals
   - Build simple trend analysis for expert growth

2. **Report Generation**
   - Create exportable reports for key metrics
   - Implement scheduled report generation
   - Add custom report parameters

3. **Operational Analytics**
   - Build phase completion statistics
   - Implement expert utilization metrics
   - Create workflow efficiency indicators

**Estimated Duration**: 1 week

**Success Criteria**:
- Statistics load quickly (under 2 seconds)
- Reports can be easily generated and exported
- Key metrics are clearly visualized
- Data provides actionable insights for planning

## Implementation Approach

### Development Methodology
- Implement features in small, testable increments
- Prioritize admin and core workflows first
- Focus on functionality over visual aesthetics
- Regular testing with real user workflows

### Technical Simplification
- Use React Context and useReducer for state management
- Implement localStorage caching where appropriate
- Favor server-side processing over complex client-side operations
- Keep component hierarchy shallow and focused

### Role-Based Implementation
- Develop admin interfaces first
- Build planner capabilities second
- Implement regular user features last
- Maintain consistent role visibility throughout

## Technical Considerations

### Lightweight State Management
- Use simple Context providers without complex libraries
- Implement optimistic updates for better UX
- Use localStorage for non-sensitive data persistence

### Performance Optimization
- Implement efficient server-side filtering and sorting
- Limit client-side processing
- Optimize initial page load times
- Minimize dependencies

### Security Considerations
- Ensure proper authentication for all API requests
- Implement strict input validation
- Apply appropriate role-based access controls
- Secure document handling

## Next Steps

1. Begin with authentication and authorization foundation
2. Implement Stage 1 (Search & Filter Foundation)
3. Continue through stages in order, with testing at each stage
4. Prioritize admin workflows first, then planner and regular user features

## Risks and Mitigation

| Risk | Mitigation |
|------|------------|
| API changes affecting frontend | Document API contracts, automated tests |
| Role-based permission complexity | Start with admin-only view, add role restrictions gradually |
| Document security concerns | Implement proper validation and scanning |
| Workflow edge cases | Test with real-world scenarios frequently |

## Conclusion

This revised phased approach aligns with the project's focus on simplicity, maintainability, and internal use. By prioritizing core workflows and starting with admin-centric features, we ensure the most critical functionality is delivered first while maintaining security and performance expectations.