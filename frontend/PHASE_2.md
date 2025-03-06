# Phase 2: Expert Database Searching for ExpertDB Frontend

## Overview
You are an AI tasked with implementing the search functionality for the ExpertDB frontend using React, TypeScript, and shadcn/ui components. The project uses Vite, React Router, and Axios (in `src/api/api.ts`) to fetch data from `http://localhost:8080/api/experts` and `http://localhost:8080/api/isced/fields`. The `AuthContext` from Phase 1 manages authentication. Use `.tsx` for React components and `.ts` for other files. Follow each step exactly, verify at each point, and avoid common pitfalls.

## Objective
- Create a `Search` page to fetch and display a list of experts.
- Add search by `name` and `isced_field_id` using an input and dropdown.
- Add filters for `affiliation`, `is_bahraini`, `isced_field_id`, `is_available` (default to `is_available=true`).
- Implement pagination (10 experts per page) and sorting by `name`.
- Use shadcn/ui components (`Input`, `Select`, `Table`, `Button`) for the UI.
- Ensure the page is protected (requires login).

## Step-by-Step Instructions

1. **Update api.ts with Additional Endpoints**
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

     export interface IscedField {
       id: string;
       name: string;
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

     export const getIscedFields = async (): Promise<IscedField[]> => {
       const response = await api.get("/api/isced/fields");
       return response.data;
     };

     export default api;
     ```
   - **Verify**: Run `npx eslint src/api/api.ts`—should pass.

2. **Create Search Page**
   - **File**: Create `expertdb_grok/frontend/src/pages/Search.tsx`.
   - **Code**:
     ```tsx
     import { useState, useEffect } from "react";
     import { Input } from "@/components/ui/input";
     import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
     import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
     import { Button } from "@/components/ui/button";
     import { Expert, IscedField, getExperts, getIscedFields } from "../api/api";

     export default function Search() {
       const [experts, setExperts] = useState<Expert[]>([]);
       const [iscedFields, setIscedFields] = useState<IscedField[]>([]);
       const [filters, setFilters] = useState({
         name: "",
         affiliation: "",
         is_bahraini: undefined as boolean | undefined,
         isced_field_id: "",
         is_available: true,
         page: 1,
         limit: 10,
       });
       const [sortOrder, setSortOrder] = useState<"asc" | "desc">("asc");
       const [loading, setLoading] = useState(true);
       const [error, setError] = useState("");

       useEffect(() => {
         const fetchData = async () => {
           try {
             const [expertsData, fieldsData] = await Promise.all([
               getExperts(filters),
               getIscedFields(),
             ]);
             setExperts(expertsData);
             setIscedFields(fieldsData);
           } catch (err) {
             setError("Failed to load data");
           } finally {
             setLoading(false);
           }
         };
         fetchData();
       }, [filters]);

       const handleSort = () => {
         const newOrder = sortOrder === "asc" ? "desc" : "asc";
         setSortOrder(newOrder);
         setExperts([...experts].sort((a, b) =>
           newOrder === "asc"
             ? a.name.localeCompare(b.name)
             : b.name.localeCompare(a.name)
         ));
       };

       return (
         <div className="p-6">
           <h2 className="text-2xl mb-4">Search Experts</h2>
           <div className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
             <Input
               placeholder="Search by name..."
               value={filters.name}
               onChange={(e) => setFilters({ ...filters, name: e.target.value, page: 1 })}
             />
             <Input
               placeholder="Filter by affiliation..."
               value={filters.affiliation}
               onChange={(e) => setFilters({ ...filters, affiliation: e.target.value, page: 1 })}
             />
             <Select
               value={filters.isced_field_id || ""}
               onValueChange={(value) => setFilters({ ...filters, isced_field_id: value, page: 1 })}
             >
               <SelectTrigger>
                 <SelectValue placeholder="Select ISCED Field" />
               </SelectTrigger>
               <SelectContent>
                 <SelectItem value="">All Fields</SelectItem>
                 {iscedFields.map((field) => (
                   <SelectItem key={field.id} value={field.id}>
                     {field.name}
                   </SelectItem>
                 ))}
               </SelectContent>
             </Select>
             <Select
               value={filters.is_bahraini === undefined ? "" : filters.is_bahraini ? "true" : "false"}
               onValueChange={(value) => setFilters({
                 ...filters,
                 is_bahraini: value === "" ? undefined : value === "true",
                 page: 1,
               })}
             >
               <SelectTrigger>
                 <SelectValue placeholder="Bahraini Status" />
               </SelectTrigger>
               <SelectContent>
                 <SelectItem value="">All</SelectItem>
                 <SelectItem value="true">Bahraini</SelectItem>
                 <SelectItem value="false">Non-Bahraini</SelectItem>
               </SelectContent>
             </Select>
             <Select
               value={filters.is_available ? "true" : "false"}
               onValueChange={(value) => setFilters({
                 ...filters,
                 is_available: value === "true",
                 page: 1,
               })}
             >
               <SelectTrigger>
                 <SelectValue placeholder="Availability" />
               </SelectTrigger>
               <SelectContent>
                 <SelectItem value="true">Available</SelectItem>
                 <SelectItem value="false">Not Available</SelectItem>
               </SelectContent>
             </Select>
           </div>
           {loading && <p>Loading...</p>}
           {error && <p className="text-red-500">{error}</p>}
           {!loading && !error && (
             <>
               <Table>
                 <TableHeader>
                   <TableRow>
                     <TableHead onClick={handleSort} className="cursor-pointer">
                       Name {sortOrder === "asc" ? "↑" : "↓"}
                     </TableHead>
                     <TableHead>Affiliation</TableHead>
                     <TableHead>Bahraini</TableHead>
                     <TableHead>ISCED Field</TableHead>
                     <TableHead>Available</TableHead>
                   </TableRow>
                 </TableHeader>
                 <TableBody>
                   {experts.map((expert) => (
                     <TableRow key={expert.id}>
                       <TableCell>{expert.name}</TableCell>
                       <TableCell>{expert.affiliation}</TableCell>
                       <TableCell>{expert.is_bahraini ? "Yes" : "No"}</TableCell>
                       <TableCell>{expert.isced_field_id}</TableCell>
                       <TableCell>{expert.is_available ? "Yes" : "No"}</TableCell>
                     </TableRow>
                   ))}
                 </TableBody>
               </Table>
               <div className="flex justify-between mt-4">
                 <Button
                   disabled={filters.page === 1}
                   onClick={() => setFilters({ ...filters, page: filters.page - 1 })}
                 >
                   Previous
                 </Button>
                 <span>Page {filters.page}</span>
                 <Button
                   disabled={experts.length < filters.limit}
                   onClick={() => setFilters({ ...filters, page: filters.page + 1 })}
                 >
                   Next
                 </Button>
               </div>
             </>
           )}
         </div>
       );
     }
     ```
   - **Verify**: Run `npx eslint src/pages/Search.tsx`—should pass. Run `npm run dev`, visit `http://localhost:5174/search`—should redirect to `/login` if not logged in (expected).

3. **Add shadcn/ui Components**
   - **Command**: Run `npx shadcn@latest add table select`.
   - **Verify**: Check `src/components/ui/table.tsx` and `src/components/ui/select.tsx` exist.

## Common Pitfalls and How to Avoid Them

1. **Missing shadcn/ui Components**:
   - **Pitfall**: `Table` or `Select` components might not exist if not installed.
   - **Fix**: Ensure `npx shadcn@latest add table select` runs successfully.

2. **API Errors**:
   - **Pitfall**: `/api/experts` or `/api/isced/fields` might fail if the backend isn’t running.
   - **Fix**: It’s okay—focus on UI. The `loading` and `error` states will handle failures.

3. **Case Sensitivity**:
   - **Pitfall**: Importing `getExperts` from `"../api/api"` must match the file name.
   - **Fix**: Verify `src/api/api.ts` exists and the import path is correct.

4. **Auth Redirects**:
   - **Pitfall**: If not logged in, `/search` redirects to `/login`.
   - **Fix**: Expected behavior—mock a login in the console if needed for testing.

## Final Checks
Provide the following instructions for the user to execute and verify the app’s state:
- Run `npm run dev` in the terminal.
- Log in (mock if needed: `localStorage.setItem("token", "mock-token")` in the browser console).
- Visit `http://localhost:5174/search`—should show a search input, filters, table, pagination, and sorting.
- Filter by `name`, `isced_field_id`, etc.—should update the table dynamically.
- Run `npm run build` in the terminal—should succeed.

## Questions
- Should the table include additional columns (e.g., `bio`) or actions (e.g., “View Details” button)?
- What styling preferences do you have for the pagination buttons (e.g., colors, sizes)?
- Should sorting be server-side instead of client-side for better performance with large datasets?
