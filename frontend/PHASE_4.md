# Phase 4: Statistics Dashboard for ExpertDB Frontend

## Overview
You are an AI tasked with implementing a statistics dashboard for the ExpertDB frontend using React, TypeScript, and Recharts (already installed). The project uses Vite, React Router, and Axios (in `src/api/api.ts`) to fetch data from `http://localhost:8080/api/statistics/*` endpoints. The page is protected (requires login). Follow each step exactly, verify at each point, and avoid pitfalls.

## Objective
- Create a `Statistics` page with:
  - Pie chart for Bahraini vs. non-Bahraini experts.
  - Bar chart for experts by ISCED field.
  - Line chart for annual growth of the database.
  - Bar chart for top institutions by expert count.
- Fetch data from `/api/statistics/nationality`, `/api/statistics/isced`, `/api/statistics/growth`, and derive top institutions from `/api/experts`.
- Use Recharts for visualization.

## Step-by-Step Instructions

1. **Update api.ts with Statistics Endpoints**
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

     export interface ExpertRequest {
       id: string;
       designation: string;
       name: string;
       institution: string;
       bio: string;
       status: string;
       cv_path: string;
       user_id: string;
     }

     export interface NationalityStats {
       bahraini: number;
       non_bahraini: number;
     }

     export interface IscedStats {
       field_id: string;
       name: string;
       count: number;
     }

     export interface GrowthStats {
       period: string;
       count: number;
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

     export const submitExpertRequest = async (data: {
       designation: string;
       name: string;
       institution: string;
       bio: string;
       primaryContact: string;
       contactType: string;
       skills: string[];
       availability: string;
       cv?: File;
     }): Promise<ExpertRequest> => {
       const formData = new FormData();
       formData.append("designation", data.designation);
       formData.append("name", data.name);
       formData.append("institution", data.institution);
       formData.append("bio", data.bio);
       formData.append("primaryContact", data.primaryContact);
       formData.append("contactType", data.contactType);
       data.skills.forEach((skill, index) => formData.append(`skills[${index}]`, skill));
       formData.append("availability", data.availability);
       if (data.cv) formData.append("cv", data.cv);

       const response = await api.post("/api/expert-requests", formData, {
         headers: { "Content-Type": "multipart/form-data" },
       });
       return response.data;
     };

     export const getNationalityStats = async (): Promise<NationalityStats> => {
       const response = await api.get("/api/statistics/nationality");
       return response.data;
     };

     export const getIscedStats = async (): Promise<IscedStats[]> => {
       const response = await api.get("/api/statistics/isced");
       return response.data;
     };

     export const getGrowthStats = async (): Promise<GrowthStats[]> => {
       const response = await api.get("/api/statistics/growth");
       return response.data;
     };

     export default api;
     ```
   - **Verify**: Run `npx eslint src/api/api.ts`—should pass.

2. **Create Statistics Page**
   - **File**: Create `expertdb_grok/frontend/src/pages/Statistics.tsx`.
   - **Code**:
     ```tsx
     import { useState, useEffect } from "react";
     import { PieChart, Pie, Cell, BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, LineChart, Line } from "recharts";
     import { getNationalityStats, getIscedStats, getGrowthStats, getExperts } from "../api/api";

     interface TopInstitution {
       name: string;
       count: number;
     }

     export default function Statistics() {
       const [nationalityStats, setNationalityStats] = useState<{ name: string; value: number }[]>([]);
       const [iscedStats, setIscedStats] = useState<{ name: string; count: number }[]>([]);
       const [growthStats, setGrowthStats] = useState<{ period: string; count: number }[]>([]);
       const [topInstitutions, setTopInstitutions] = useState<TopInstitution[]>([]);
       const [loading, setLoading] = useState(true);
       const [error, setError] = useState("");

       useEffect(() => {
         const fetchData = async () => {
           try {
             const [natStats, iscedData, growthData, expertsData] = await Promise.all([
               getNationalityStats(),
               getIscedStats(),
               getGrowthStats(),
               getExperts(),
             ]);

             // Nationality stats
             setNationalityStats([
               { name: "Bahraini", value: natStats.bahraini },
               { name: "Non-Bahraini", value: natStats.non_bahraini },
             ]);

             // ISCED stats
             setIscedStats(iscedData);

             // Growth stats
             setGrowthStats(growthData);

             // Top institutions
             const institutionCounts: { [key: string]: number } = {};
             expertsData.forEach((expert) => {
               const inst = expert.affiliation || "Unknown";
               institutionCounts[inst] = (institutionCounts[inst] || 0) + 1;
             });
             const topInst = Object.entries(institutionCounts)
               .map(([name, count]) => ({ name, count }))
               .sort((a, b) => b.count - a.count)
               .slice(0, 5);
             setTopInstitutions(topInst);
           } catch (err) {
             setError("Failed to load statistics");
           } finally {
             setLoading(false);
           }
         };
         fetchData();
       }, []);

       return (
         <div className="p-6">
           <h2 className="text-2xl mb-4">Statistics Dashboard</h2>
           {loading && <p>Loading...</p>}
           {error && <p className="text-red-500">{error}</p>}
           {!loading && !error && (
             <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
               <div>
                 <h3 className="text-xl mb-2">Bahraini vs Non-Bahraini Experts</h3>
                 <PieChart width={400} height={300}>
                   <Pie
                     data={nationalityStats}
                     dataKey="value"
                     nameKey="name"
                     cx="50%"
                     cy="50%"
                     outerRadius={80}
                     fill="#8884d8"
                     label
                   >
                     {nationalityStats.map((entry, index) => (
                       <Cell key={`cell-${index}`} fill={index === 0 ? "#133566" : "#1B4882"} />
                     ))}
                   </Pie>
                   <Tooltip />
                   <Legend />
                 </PieChart>
               </div>
               <div>
                 <h3 className="text-xl mb-2">Experts by ISCED Field</h3>
                 <BarChart width={400} height={300} data={iscedStats}>
                   <CartesianGrid strokeDasharray="3 3" />
                   <XAxis dataKey="name" />
                   <YAxis />
                   <Tooltip />
                   <Legend />
                   <Bar dataKey="count" fill="#133566" />
                 </BarChart>
               </div>
               <div>
                 <h3 className="text-xl mb-2">Annual Growth</h3>
                 <LineChart width={400} height={300} data={growthStats}>
                   <CartesianGrid strokeDasharray="3 3" />
                   <XAxis dataKey="period" />
                   <YAxis />
                   <Tooltip />
                   <Legend />
                   <Line type="monotone" dataKey="count" stroke="#133566" />
                 </LineChart>
               </div>
               <div>
                 <h3 className="text-xl mb-2">Top Institutions</h3>
                 <BarChart width={400} height={300} data={topInstitutions}>
                   <CartesianGrid strokeDasharray="3 3" />
                   <XAxis dataKey="name" />
                   <YAxis />
                   <Tooltip />
                   <Legend />
                   <Bar dataKey="count" fill="#1B4882" />
                 </BarChart>
               </div>
             </div>
           )}
         </div>
       );
     }
     ```
   - **Verify**: Run `npx eslint src/pages/Statistics.tsx`—should pass.

## Common Pitfalls and How to Avoid Them

1. **Recharts Imports**:
   - **Pitfall**: Missing Recharts components (e.g., `BarChart`) cause errors.
   - **Fix**: Ensure all imports match the code (already installed).

2. **API Errors**:
   - **Pitfall**: `/api/statistics/*` might fail if the backend isn’t running.
   - **Fix**: Focus on UI—error state will handle failures.

3. **Data Format**:
   - **Pitfall**: API response might not match expected interfaces.
   - **Fix**: Adjust interfaces if backend response differs (assumed formats provided).

## Final Checks
Provide the following instructions for the user to execute and verify the app’s state:
- Run `npm run dev` in the terminal.
- Log in (mock if needed: `localStorage.setItem("token", "mock-token")`).
- Visit `http://localhost:5174/statistics`—should show charts for nationality, ISCED fields, growth, and top institutions.
- Run `npm run build` in the terminal—should succeed.

## Questions
- What specific colors should the charts use (e.g., for Bahraini vs. non-Bahraini)?
- Should the charts include interactive features (e.g., hover tooltips, click to filter)?
- How many top institutions should be displayed (currently set to 5)?
