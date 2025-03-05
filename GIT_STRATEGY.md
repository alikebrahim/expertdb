# ExpertDB Git Strategy

## Branching Model
- **main**: Production-ready code
- **develop**: Integration branch for features
- **feature/[feature-name]**: Individual feature development
- **bugfix/[issue-ID]**: Bug fixes
- **release/[version]**: Release preparation
- **hotfix/[issue-ID]**: Emergency production fixes

## Commit Message Convention
Format: `[type]([scope]): [description]`

Types:
- **feat**: New feature
- **fix**: Bug fix
- **docs**: Documentation changes
- **style**: Formatting, styling changes
- **refactor**: Code restructuring without functionality change
- **test**: Test additions/modifications
- **chore**: Build process or tooling changes

Example: `feat(auth): implement JWT token validation`

## Development Workflow
1. Create feature branch from develop
2. Implement changes with atomic commits
3. Open PR to develop branch
4. Code review and testing
5. Merge to develop upon approval
6. Periodically merge develop to main as releases

## Initial Commit Strategy
Given the current project state, the following commits are recommended:

1. `docs: add project documentation and guidelines`
   - README.md, IMPLEMENTATION.md, ISSUES.md, STATUS.md, FIRST_RUN.md
   - CLAUDE.md files (project, frontend, backend)

2. `chore: initial project configuration`
   - Docker configurations
   - Build scripts
   - Run scripts

3. `feat(backend): implement core API structure`
   - Go API foundation
   - Authentication system
   - Database schema and migrations

4. `feat(frontend): set up Next.js application structure`
   - Component framework
   - Authentication flow
   - UI component library integration

5. `feat(integration): connect frontend and backend services`
   - API integration
   - Authentication flow

After initial setup, follow the standard branching model for ongoing development.