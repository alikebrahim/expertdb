# Phase 0: Setup for ExpertDB Frontend

## Overview
This phase documents the completed setup for the ExpertDB frontend. It’s already done but included for reference.

## Steps Completed
1. Initialized project with Vite, React, TypeScript.
2. Installed dependencies: shadcn/ui, Tailwind CSS (v4), Recharts, Axios, jsPDF, Mammoth.
3. Configured Tailwind CSS (`tailwind.config.js`, `src/index.css`).
4. Set up shadcn/ui (`components.json`, `src/components/ui/`).
5. Created API client (`src/api/api.ts`) with Axios and JWT interceptor.
6. Configured linting/formatting (ESLint, Prettier).
7. Set up routing (`src/App.tsx`, `src/main.tsx`) with React Router.
8. Created directory structure (`src/api/`, `src/components/layout/`, `src/pages/`, `public/logos/`).

## Final State
- App runs at `http://localhost:5174/` with placeholder routes (`/login`, `/search`, etc.).
- API client is ready in `src/api/api.ts`.
- Poppins font is used (`src/index.css`).
