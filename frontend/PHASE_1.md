# Phase 1: Authentication for ExpertDB Frontend

## Overview
You are an AI tasked with implementing authentication for the ExpertDB frontend using React, TypeScript, and React Router. The project uses Vite and has `"type": "module"` in `package.json`, so all `.js` files are ESM. Use `.tsx` for React components and `.ts` for other files. The backend API is at `http://localhost:8080`, and you’ll use Axios (already installed) to make requests. The app has shadcn/ui for UI components, and the API client is in `src/api/api.ts`. Follow each step exactly, verify at each point, and avoid common pitfalls.

## Objective
- Create an `AuthContext` to manage user authentication state (token, user, role).
- Build a `Login` page to authenticate users via `POST /api/auth/login`.
- Implement a `ProtectedRoute` component to restrict access based on login status and role (`/admin` for admins only, `/requests`, `/search`, `/statistics` for logged-in users).
- Update routing to use `AuthProvider`, protect routes, and redirect `/` to `/login` if unauthenticated or `/search` if authenticated.

## Step-by-Step Instructions

1. **Create AuthContext**
   - **File**: Create `expertdb_grok/frontend/src/context/AuthContext.tsx`.
   - **Code**:
     ```tsx
     import React, { createContext, useContext, useState } from "react";

     interface AuthContextType {
       token: string | null;
       user: { id: string; name: string; email: string; role: string } | null;
       login: (token: string, user: any) => void;
       logout: () => void;
     }

     const AuthContext = createContext<AuthContextType | undefined>(undefined);

     export function AuthProvider({ children }: { children: React.ReactNode }) {
       const [authState, setAuthState] = useState<{
         token: string | null;
         user: { id: string; name: string; email: string; role: string } | null;
       }>({
         token: localStorage.getItem("token"),
         user: null,
       });

       const login = (token: string, user: any) => {
         localStorage.setItem("token", token);
         setAuthState({ token, user });
       };

       const logout = () => {
         localStorage.removeItem("token");
         setAuthState({ token: null, user: null });
       };

       return (
         <AuthContext.Provider value={{ ...authState, login, logout }}>
           {children}
         </AuthContext.Provider>
       );
     }

     export const useAuth = () => {
       const context = useContext(AuthContext);
       if (!context) throw new Error("useAuth must be used within an AuthProvider");
       return context;
     };
     ```
   - **Verify**: Run `npx eslint src/context/AuthContext.tsx`—should pass with no errors. If errors, fix them (e.g., add `// eslint-disable-next-line` for unused vars).

2. **Update main.tsx with AuthProvider**
   - **File**: Update `expertdb_grok/frontend/src/main.tsx`.
   - **Code**:
     ```tsx
     import React from "react";
     import ReactDOM from "react-dom/client";
     import { BrowserRouter } from "react-router-dom";
     import App from "./App.tsx";
     import "./index.css";
     import { AuthProvider } from "./context/AuthContext";

     ReactDOM.createRoot(document.getElementById("root")!).render(
       <React.StrictMode>
         <BrowserRouter>
           <AuthProvider>
             <App />
           </AuthProvider>
         </BrowserRouter>
       </React.StrictMode>
     );
     ```
   - **Verify**: Run `npx eslint src/main.tsx`—should pass. Then run `npm run dev` and visit `http://localhost:5174/`—should load without errors (same as before).

3. **Create ProtectedRoute Component**
   - **File**: Create `expertdb_grok/frontend/src/components/ProtectedRoute.tsx`.
   - **Code**:
     ```tsx
     import { Navigate } from "react-router-dom";
     import { useAuth } from "../context/AuthContext";

     interface ProtectedRouteProps {
       children: React.ReactNode;
       requireAdmin?: boolean;
     }

     export default function ProtectedRoute({ children, requireAdmin = false }: ProtectedRouteProps) {
       const { token, user } = useAuth();

       if (!token) {
         return <Navigate to="/login" replace />;
       }

       if (requireAdmin && user?.role !== "admin") {
         return <Navigate to="/search" replace />;
       }

       return <>{children}</>;
     }
     ```
   - **Verify**: Run `npx eslint src/components/ProtectedRoute.tsx`—should pass. If errors, fix them.

4. **Update api.ts with Login Call**
   - **File**: Update `expertdb_grok/frontend/src/api/api.ts`.
   - **Code**:
     ```typescript
     import axios from "axios";

     const api = axios.create({
       baseURL: "http://localhost:8080",
       headers: { "Content-Type": "application/json" },
     });

     api.interceptors.request.use(
       (config) => {
         const token = localStorage.getItem("token");
         if (token) {
           config.headers.Authorization = `Bearer ${token}`;
         }
         return config;
       },
       (error) => Promise.reject(error)
     );

     export interface Expert {
       id: string;
       name: string;
       affiliation: string;
       is_bahraini: boolean;
       isced_field_id: string;
       is_available: boolean;
       bio: string;
     }

     export interface AuthResponse {
       token: string;
       user: {
         id: string;
         name: string;
         email: string;
         role: string;
       };
     }

     export const login = async (email: string, password: string): Promise<AuthResponse> => {
       const response = await api.post("/api/auth/login", { email, password });
       return response.data;
     };

     export const getExperts = async (filters: {
       name?: string;
       affiliation?: string;
       is_bahraini?: boolean;
       isced_field_id?: string;
       is_available?: boolean;
       page?: number;
       limit?: number;
     } = {}): Promise<Expert[]> => {
       const params = new URLSearchParams();
       if (filters.name) params.append("name", filters.name);
       if (filters.affiliation) params.append("affiliation", filters.affiliation);
       if (filters.is_bahraini !== undefined) params.append("is_bahraini", String(filters.is_bahraini));
       if (filters.isced_field_id) params.append("isced_field_id", filters.isced_field_id);
       if (filters.is_available !== undefined) params.append("is_available", String(filters.is_available));
       if (filters.page) params.append("page", String(filters.page));
       if (filters.limit) params.append("limit", String(filters.limit));

       const response = await api.get("/api/experts", { params });
       return response.data;
     };

     export default api;
     ```
   - **Verify**: Run `npx eslint src/api/api.ts`—should pass.

5. **Create Login Page**
   - **File**: Create `expertdb_grok/frontend/src/pages/Login.tsx`.
   - **Code**:
     ```tsx
     import { useState } from "react";
     import { useNavigate } from "react-router-dom";
     import { Input } from "@/components/ui/input";
     import { Button } from "@/components/ui/button";
     import { useAuth } from "../context/AuthContext";
     import { login } from "../api/api";

     export default function Login() {
       const [email, setEmail] = useState("");
       const [password, setPassword] = useState("");
       const [error, setError] = useState("");
       const { login: loginUser } = useAuth();
       const navigate = useNavigate();

       const handleSubmit = async (e: React.FormEvent) => {
         e.preventDefault();
         try {
           const { token, user } = await login(email, password);
           loginUser(token, user);
           navigate(user.role === "admin" ? "/admin" : "/search");
         } catch (err) {
           setError("Invalid credentials");
         }
       };

       return (
         <div className="flex items-center justify-center min-h-screen bg-gray-50">
           <form onSubmit={handleSubmit} className="p-6 bg-white rounded shadow-md">
             <h2 className="text-2xl mb-4">Login</h2>
             {error && <p className="text-red-500 mb-4">{error}</p>}
             <Input
               type="email"
               placeholder="Email"
               value={email}
               onChange={(e) => setEmail(e.target.value)}
               className="mb-4"
             />
             <Input
               type="password"
               placeholder="Password"
               value={password}
               onChange={(e) => setPassword(e.target.value)}
               className="mb-4"
             />
             <Button type="submit">Login</Button>
           </form>
         </div>
       );
     }
     ```
   - **Verify**: Run `npx eslint src/pages/Login.tsx`—should pass. Run `npm run dev`, visit `http://localhost:5174/login`—should see a login form. If the backend isn’t running, expect a network error on submit (normal).

6. **Update App.tsx with Protected Routes and Redirect Logic**
   - **File**: Update `expertdb_grok/frontend/src/App.tsx`.
   - **Code**:
     ```tsx
     import { Routes, Route, Navigate } from "react-router-dom";
     import ProtectedRoute from "./components/ProtectedRoute";
     import Login from "./pages/Login";
     import { useAuth } from "./context/AuthContext";

     function App() {
       const { token } = useAuth();

       return (
         <Routes>
           <Route path="/login" element={<Login />} />
           <Route path="/search" element={<ProtectedRoute><div>Search Page</div></ProtectedRoute>} />
           <Route path="/requests" element={<ProtectedRoute><div>Requests Page</div></ProtectedRoute>} />
           <Route path="/statistics" element={<ProtectedRoute><div>Statistics Page</div></ProtectedRoute>} />
           <Route path="/admin" element={<ProtectedRoute requireAdmin><div>Admin Panel</div></ProtectedRoute>} />
           <Route path="/" element={token ? <Navigate to="/search" replace /> : <Navigate to="/login" replace />} />
         </Routes>
       );
     }

     export default App;
     ```
   - **Verify**: Run `npx eslint src/App.tsx`—should pass. Run `npm run dev`. Visit `http://localhost:5174/`—should redirect to `/login` (since you’re not logged in).

## Common Pitfalls and How to Avoid Them

1. **File Case Sensitivity**:
   - **Pitfall**: Linux is case-sensitive. Mismatches like `src/testApi.ts` vs. `src/testAPI.ts` caused errors before.
   - **Fix**: Use camelCase for `.tsx` files (e.g., `AuthContext.tsx`) and lowercase for `.ts` files (e.g., `api.ts`). Ensure imports match file names exactly.

2. **ESM vs. CommonJS Mismatch**:
   - **Pitfall**: Using `module.exports` in `.js` files breaks with `"type": "module"`.
   - **Fix**: All `.js` files must use ESM (`export default`). Config files like `postcss.config.cjs` must stay `.cjs`.

3. **AuthContext Errors**:
   - **Pitfall**: Using `useAuth` outside `AuthProvider` causes "useAuth must be used within an AuthProvider".
   - **Fix**: Ensure `AuthProvider` wraps `<App />` in `main.tsx` as shown.

4. **Backend Not Running**:
   - **Pitfall**: `POST /api/auth/login` will fail if `http://localhost:8080` isn’t up.
   - **Fix**: It’s okay for now—test with a mock response or wait for backend setup. Focus on UI behavior (e.g., redirect to `/login`).

## Final Checks
Provide the following instructions for the user to execute and verify the app’s state:
- Run `npm run dev` in the terminal.
- Visit `http://localhost:5174/`—should redirect to `/login` since you’re not logged in.
- Visit `http://localhost:5174/login`—should see a login form.
- Visit `http://localhost:5174/search`—should redirect to `/login`.
- Visit `http://localhost:5174/admin`—should redirect to `/login`.
- Run `npm run build` in the terminal—should succeed with no TypeScript errors.

## Questions
- Should the login form include additional fields (e.g., "Remember Me" checkbox)?
- What should the UI look like if the login fails (e.g., specific styling for the error message)?
- After a successful login, should there be a loading state or animation before redirecting to `/search` or `/admin`?
