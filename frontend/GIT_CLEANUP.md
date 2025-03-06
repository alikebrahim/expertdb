# Git Repository Cleanup Notes

## Current Git Status

The current git status shows:

1. **Deleted files**: Many `.next` directory files and various Next.js-related files are marked as deleted. This appears to be from a previous Next.js-based iteration of the frontend that has been replaced with the Vite-based implementation.

2. **Modified backend files**: Several backend files are showing as modified, which are outside the scope of the frontend directory.

3. **Untracked frontend files**: Several configuration and static files in the frontend directory that should be tracked.

## Cleanup Plan

### Step 1: Handle Deleted Next.js Files

These files are from the previous Next.js implementation that has been replaced with Vite. Since we're now using Vite, these files should be properly removed from git tracking:

```bash
# Remove deleted Next.js files from git tracking
git rm -r --cached .next
git rm -r --cached app
git rm -r --cached components
git rm -r --cached lib
git rm -r --cached public/images
git rm --cached Dockerfile
git rm --cached UI_UX_GUIDELINES.md
git rm --cached next-env.d.ts
git rm --cached next.config.js
git rm --cached postcss.config.js
git rm --cached debugging_notes.md
git rm -r --cached issues
```

### Step 2: Handle Modified Backend Files

Since these files are outside the frontend directory, they should be managed separately by someone working on the backend:

```bash
# Assuming we only want to commit frontend changes
git checkout -- ../backend/api.go
git checkout -- ../backend/expert_operations.go
git checkout -- ../backend/expert_request.go
git checkout -- ../backend/expert_request_operations.go
git checkout -- ../backend/import_experts.go
git checkout -- ../backend/storage.go
git checkout -- ../backend/types.go
git checkout -- ../notes.md
```

### Step 3: Add Untracked Frontend Files

Add all relevant untracked files for the frontend implementation:

```bash
# Add relevant configuration files
git add .eslintrc.json
git add .gitignore
git add .prettierrc
```

### Step 4: Commit Changes

```bash
git commit -m "chore: clean up repository and git tracking

- Remove deleted Next.js files from git tracking
- Add missing configuration files
- Clean up repository structure for Vite-based implementation"
```

## Recommended `.gitignore` Additions

Consider adding the following patterns to `.gitignore` to avoid future issues:

```
# Build directories
dist/
build/
out/
.next/

# Local environment files
.env.local
.env.development.local
.env.test.local
.env.production.local

# Tooling files
.vscode/
.idea/
*.sublime-project
*.sublime-workspace

# OS files
.DS_Store
Thumbs.db

# Logs
npm-debug.log*
yarn-debug.log*
yarn-error.log*
```

## Future Recommendations

1. **Branch Management**: Create feature branches for each phase instead of working directly on master
2. **Project Structure**: Maintain clear separation between frontend and backend in the repository
3. **Commit Guidelines**: Continue using conventional commit messages for clarity
4. **PR Process**: Consider implementing Pull Requests for code review before merging to master