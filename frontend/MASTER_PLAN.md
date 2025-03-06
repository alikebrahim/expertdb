# Master Plan for ExpertDB Frontend Development

## Overview
This project is a frontend web app for ExpertDB, developed for the Bahrain Education & Training Quality Authority (BQA). The app allows users to authenticate, search experts, submit requests, view statistics, and manage admin tasks. It uses React, TypeScript, Vite, shadcn/ui (with Tailwind CSS), Recharts, Axios, jsPDF, and Mammoth. The backend API is at `http://localhost:8080`.

## Prerequisites
- **Setup Complete (Phase 0)**: The project setup is complete (Vite, React, TypeScript, shadcn/ui, Tailwind CSS, Axios, Recharts, jsPDF, Mammoth, routing, API client in `src/api/api.ts`, linting/formatting). See [PHASE_0.md](PHASE_0.md) for details.
- **Directory**: All work is in `expertdb_grok/frontend/`.
- **File Naming and Case Sensitivity**:
  - For React component files (`.tsx`), use camelCase (e.g., `AuthContext.tsx`, `Login.tsx`) to match React conventions.
  - For non-component files (e.g., utilities, API files), use lowercase (e.g., `api.ts`, `utils.ts`).
  - Linux is case-sensitive, so ensure imports match the file name exactly. For example, if the file is `AuthContext.tsx`, the import must be `import { AuthProvider } from "./context/AuthContext"`, not `import { AuthProvider } from "./context/authcontext"`.

## General Instructions for LLM
You are an AI tasked with implementing the ExpertDB frontend. Follow each phase in order, using the instructions in the respective `PHASE_X.md` files. Key guidelines:

- **File Paths**: Use exact paths (e.g., `expertdb_grok/frontend/src/context/AuthContext.tsx`).
- **ESM vs. CommonJS**: The project uses `"type": "module"` in `package.json`. All `.js` files must use ESM (`export default`). Config files (e.g., `postcss.config.cjs`) must use `.cjs` extension.
- **Backend Not Running**: If API calls fail (e.g., `http://localhost:8080` down), focus on UI and error handling. Network errors are expected.
- **Verify Each Step**: After each step, run `npx eslint <file>` and `npm run dev` to catch issues early.
- **Questions**: Analyze all `PHASE_X.md` files (PHASE_1.md to PHASE_6.md) in this directory. In each phase’s “Questions” section, formulate questions for the user if more information is needed (e.g., missing API details, UI styling preferences). Example questions:
  - “What should the redirect behavior be after a successful login?”
  - “What colors should the charts use in the statistics dashboard?”

## Phases
- [Phase 0: Setup (Complete)](PHASE_0.md)
- [Phase 1: Authentication](PHASE_1.md)
- [Phase 2: Expert Database Searching](PHASE_2.md)
- [Phase 3: Request Submission](PHASE_3.md)
- [Phase 4: Statistics Dashboard](PHASE_4.md)
- [Phase 5: Admin Panel](PHASE_5.md)
- [Phase 6: Polish and Deployment Prep](PHASE_6.md)

## Common Pitfalls
1. **Case Sensitivity**: Double-check file names and imports (e.g., `src/api/api.ts` must match exactly).
2. **ESM Errors**: Use `export default` in `.js` files, not `module.exports`.
3. **Tailwind Issues**: Ensure `theme.extend` is used in `tailwind.config.js` to preserve default utilities.
4. **API Failures**: If the backend isn’t running, UI should show error/loading states—don’t panic.
