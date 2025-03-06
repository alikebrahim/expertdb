# Phase 6: Polish and Deployment Prep for ExpertDB Frontend

## Overview
You are an AI tasked with finalizing the ExpertDB frontend by adding a layout, polishing the UI, and preparing for deployment. The project uses Vite, React, TypeScript, React Router, and shadcn/ui components. Logos are in `public/logos/`. Follow each step exactly, verify at each point, and avoid pitfalls.

## Objective
- Add a layout with a navbar and sidebar using BQA branding (logos, colors).
- Polish the UI with consistent styling.
- Prepare for local deployment with Docker (create Dockerfile, document docker-compose setup).

## Step-by-Step Instructions

1. **Create AppLayout Component**
   - **File**: Create `expertdb_grok/frontend/src/components/layout/AppLayout.tsx`.
   - **Code**:
     ```tsx
     import { useState } from "react";
     import { Link, useNavigate } from "react-router-dom";
     import { useAuth } from "../../context/AuthContext";
     import { Button } from "@/components/ui/button";

     export default function AppLayout({ children }: { children: React.ReactNode }) {
       const { user, logout } = useAuth();
       const navigate = useNavigate();
       const [isSidebarOpen, setIsSidebarOpen] = useState(false);

       const navItems = [
         { text: "Search Experts", path: "/search" },
         { text: "Submit Request", path: "/requests" },
         { text: "Statistics", path: "/statistics" },
         ...(user?.role === "admin" ? [{ text: "Admin Panel", path: "/admin" }] : []),
       ];

       const handleLogout = () => {
         logout();
         navigate("/login");
       };

       return (
         <div className="flex min-h-screen">
           {/* Sidebar */}
           <div
             className={`fixed inset-y-0 left-0 w-64 bg-navy text-white transform ${
               isSidebarOpen ? "translate-x-0" : "-translate-x-full"
             } md:translate-x-0 transition-transform duration-300 ease-in-out z-50`}
           >
             <div className="p-4">
               <img
                 src="/logos/BQA-Horizontal-Logo.svg"
                 alt="BQA Logo"
                 className="h-10 mb-6"
               />
               <nav>
                 {navItems.map((item) => (
                   <Link
                     key={item.text}
                     to={item.path}
                     className="block py-2 px-4 hover:bg-lightblue rounded"
                     onClick={() => setIsSidebarOpen(false)}
                   >
                     {item.text}
                   </Link>
                 ))}
               </nav>
             </div>
           </div>

           {/* Main Content */}
           <div className="flex-1 flex flex-col md:ml-64">
             {/* Navbar */}
             <header className="bg-navy text-white p-4 flex justify-between items-center">
               <div className="flex items-center">
                 <Button
                   variant="ghost"
                   className="md:hidden text-white"
                   onClick={() => setIsSidebarOpen(!isSidebarOpen)}
                 >
                   ☰
                 </Button>
                 <h1 className="text-xl ml-4">ExpertDB</h1>
               </div>
               {user && (
                 <Button variant="outline" onClick={handleLogout}>
                   Logout
                 </Button>
               )}
             </header>
             <main className="flex-1 p-6 bg-gray-50">{children}</main>
             <footer className="bg-navy text-white p-4 text-center">
               <p>© 2025 BQA. All rights reserved.</p>
             </footer>
           </div>
         </div>
       );
     }
     ```
   - **Verify**: Run `npx eslint src/components/layout/AppLayout.tsx`—should pass. Ensure `public/logos/BQA-Horizontal-Logo.svg` exists (placeholder from earlier).

2. **Update App.tsx with Layout**
   - **File**: Update `expertdb_grok/frontend/src/App.tsx`.
   - **Code**:
     ```tsx
     import { Routes, Route, Navigate } from "react-router-dom";
     import AppLayout from "./components/layout/AppLayout";
     import ProtectedRoute from "./components/ProtectedRoute";
     import Login from "./pages/Login";
     import Search from "./pages/Search";
     import Requests from "./pages/Requests";
     import Statistics from "./pages/Statistics";
     import Admin from "./pages/Admin";
     import { useAuth } from "./context/AuthContext";

     function App() {
       const { token } = useAuth();

       return (
         <Routes>
           <Route path="/login" element={<Login />} />
           <Route
             path="*"
             element={
               <AppLayout>
                 <Routes>
                   <Route path="/search" element={<ProtectedRoute><Search /></ProtectedRoute>} />
                   <Route path="/requests" element={<ProtectedRoute><Requests /></ProtectedRoute>} />
                   <Route path="/statistics" element={<ProtectedRoute><Statistics /></ProtectedRoute>} />
                   <Route path="/admin" element={<ProtectedRoute requireAdmin><Admin /></ProtectedRoute>} />
                   <Route path="/" element={<Navigate to={token ? "/search" : "/login"} replace />} />
                 </Routes>
               </AppLayout>
             }
           />
         </Routes>
       );
     }

     export default App;
     ```
   - **Verify**: Run `npx eslint src/App.tsx`—should pass. Run `npm run dev`, visit `http://localhost:5174/`—should redirect to `/login` if not logged in.

3. **Polish UI**
   - **File**: Update `src/index.css` to add global styles:
     ```css
     @tailwind base;
     @tailwind components;
     @tailwind utilities;

     @import url("@fontsource/poppins/400.css");
     @import url("@fontsource/poppins/500.css");

     :root {
       font-family: "Poppins", -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, sans-serif;
       line-height: 1.5;
       font-weight: 400;
       color-scheme: light dark;
     }

     body {
       margin: 0;
       min-width: 320px;
       min-height: 100vh;
     }

     a {
       font-weight: 500;
       color: #1B4882;
       text-decoration: none;
     }
     a:hover {
       color: #133566;
     }

     @media (prefers-color-scheme: light) {
       :root {
         color: #111827;
         background-color: #ffffff;
       }
     }
     ```
   - **Verify**: Run `npm run dev`—all pages should have consistent typography and colors.

4. **Prepare for Local Deployment with Docker**
   - **File**: Create `expertdb_grok/frontend/Dockerfile`.
   - **Code**:
     ```dockerfile
     # Use Node.js 20 as the base image
     FROM node:20-alpine

     # Set working directory
     WORKDIR /app

     # Copy package.json and package-lock.json
     COPY package*.json ./

     # Install dependencies
     RUN npm install

     # Copy the rest of the application
     COPY . .

     # Build the app
     RUN npm run build

     # Expose the port Vite uses
     EXPOSE 5173

     # Command to serve the app
     CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0"]
     ```
   - **File**: Create `expertdb_grok/frontend/docker-compose.yml`.
   - **Code**:
     ```yaml
     version: "3.8"
     services:
       frontend:
         build: .
         ports:
           - "5173:5173"
         volumes:
           - .:/app
           - /app/node_modules
         environment:
           - NODE_ENV=development
     ```
   - **Documentation**: Add to `expertdb_grok/frontend/README.md`:
     ```
     ## Deployment with Docker

     To run the frontend locally with Docker:

     1. Ensure Docker and Docker Compose are installed.
     2. Build and run the container:
        ```bash
        docker-compose up --build
        ```
     3. Access the app at `http://localhost:5173/`.
     4. To stop, run:
        ```bash
        docker-compose down
        ```

     Note: This setup is for local development. For production, update the `Dockerfile` to serve the built app (e.g., with `vite preview` or an Nginx server).
     ```
   - **Verify**: Run `docker-compose up --build`—should start the app at `http://localhost:5173/`.

## Common Pitfalls and How to Avoid Them

1. **Logo Path Errors**:
   - **Pitfall**: `public/logos/BQA-Horizontal-Logo.svg` might not exist.
   - **Fix**: Ensure placeholders are in `public/logos/` (from setup).

2. **Auth Redirects**:
   - **Pitfall**: Layout might break if not logged in.
   - **Fix**: Ensure `ProtectedRoute` wraps all routes except `/login`.

3. **Docker Issues**:
   - **Pitfall**: Docker might fail if port 5173 is in use.
   - **Fix**: Stop other processes using port 5173 or change the port in `docker-compose.yml`.

## Final Checks
Provide the following instructions for the user to execute and verify the app’s state:
- Run `npm run dev` in the terminal (or `docker-compose up --build` if using Docker).
- Log in (mock if needed: `localStorage.setItem("token", "mock-token")`).
- Visit all routes (`/search`, `/requests`, `/statistics`, `/admin`)—layout should be consistent with navbar, sidebar, and footer.
- Run `npm run build` in the terminal—should create a `dist/` folder.
- Test the build: `npx vite preview`—visit `http://localhost:4173/` to ensure it works.

## Questions
- Should the sidebar include additional links or features (e.g., user profile)?
- What styling preferences do you have for the navbar and footer (e.g., colors, fonts)?
- For production deployment, will the backend be bundled in the same `docker-compose.yml`, or will it be separate?
