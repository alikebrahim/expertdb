#!/bin/bash

# Script to clean up Git repository based on .gitignore patterns

# Remove frontend node_modules from git tracking (but keep files on disk)
git rm -r --cached frontend/node_modules/

# Remove cached files for other .gitignore patterns
git rm -r --cached frontend/.next/
git rm -r --cached backend/logs/*.log
git rm -r --cached backend/db/sqlite/*.sqlite

# Add all changes and commit 
git status
echo "Run 'git add .' and then 'git commit -m \"chore: cleanup repository based on gitignore patterns\"' to complete cleanup"