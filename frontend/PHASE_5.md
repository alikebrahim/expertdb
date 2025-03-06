# Phase 5: Admin Panel for ExpertDB Frontend

## Overview
You are an AI tasked with implementing an admin panel for the ExpertDB frontend using React, TypeScript, and shadcn/ui components. The project uses Vite, React Router, and Axios (in `src/api/api.ts`) to fetch data from `http://localhost:8080/api/requests`, `/api/experts`, `/api/users`, and other endpoints. The page is protected (requires admin role). Follow each step exactly, verify at each point, and avoid pitfalls.

## Objective
- Create an `Admin` page to:
  - List all experts.
  - List all requests with options to approve/reject.
  - Manage users (create, delete).
  - Show statistics for each user (e.g., number of requests submitted).
- Use shadcn/ui components (`Table`, `Button`, `Input`, `Select`) for the UI.

## Step-by-Step Instructions

1. **Update api.ts with Admin Endpoints**
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
       reason?: string;
     }

     export interface User {
       id: string;
       name: string;
       email: string;
       role: string;
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

     export const getRequests = async (status?: string): Promise<ExpertRequest[]> => {
       const params = new URLSearchParams();
       if (status) params.append("status", status);
       const response = await api.get("/api/expert-requests", { params });
       return response.data;
     };

     export const updateRequest = async (id: string, data: {
       designation?: string;
       name?: string;
       institution?: string;
       bio?: string;
       status: string;
       reason?: string;
     }): Promise<ExpertRequest> => {
       const response = await api.put(`/api/expert-requests/${id}`, data);
       return response.data;
     };

     export const getUsers = async (): Promise<User[]> => {
       const response = await api.get("/api/users");
       return response.data;
     };

     export const createUser = async (data: {
       name: string;
       email: string;
       password: string;
       role: string;
     }): Promise<User> => {
       const response = await api.post("/api/users", data);
       return response.data;
     };

     export const deleteUser = async (id: string): Promise<void> => {
       await api.delete(`/api/users/${id}`);
     };

     export default api;
     ```
   - **Verify**: Run `npx eslint src/api/api.ts`—should pass.

2. **Create Admin Page**
   - **File**: Create `expertdb_grok/frontend/src/pages/Admin.tsx`.
   - **Code**:
     ```tsx
     import { useState, useEffect } from "react";
     import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
     import { Button } from "@/components/ui/button";
     import { Input } from "@/components/ui/input";
     import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
     import { Expert, ExpertRequest, User, getExperts, getRequests, getUsers, createUser, deleteUser, updateRequest } from "../api/api";

     export default function Admin() {
       const [experts, setExperts] = useState<Expert[]>([]);
       const [requests, setRequests] = useState<ExpertRequest[]>([]);
       const [users, setUsers] = useState<User[]>([]);
       const [userStats, setUserStats] = useState<{ [key: string]: number }>({});
       const [newUser, setNewUser] = useState({ name: "", email: "", password: "", role: "user" });
       const [loading, setLoading] = useState(true);
       const [error, setError] = useState("");

       useEffect(() => {
         const fetchData = async () => {
           try {
             const [expertsData, requestsData, usersData] = await Promise.all([
               getExperts(),
               getRequests(),
               getUsers(),
             ]);
             setExperts(expertsData);
             setRequests(requestsData);
             setUsers(usersData);

             // Calculate user stats
             const stats: { [key: string]: number } = {};
             requestsData.forEach((req) => {
               stats[req.user_id] = (stats[req.user_id] || 0) + 1;
             });
             setUserStats(stats);
           } catch (err) {
             setError("Failed to load data");
           } finally {
             setLoading(false);
           }
         };
         fetchData();
       }, []);

       const handleApprove = async (requestId: string) => {
         try {
           const updated = await updateRequest(requestId, { status: "approved" });
           setRequests(requests.map((req) =>
             req.id === requestId ? updated : req
           ));
         } catch (err) {
           setError("Failed to approve request");
         }
       };

       const handleReject = async (requestId: string, reason: string) => {
         try {
           const updated = await updateRequest(requestId, { status: "rejected", reason });
           setRequests(requests.map((req) =>
             req.id === requestId ? updated : req
           ));
         } catch (err) {
           setError("Failed to reject request");
         }
       };

       const handleCreateUser = async (e: React.FormEvent) => {
         e.preventDefault();
         try {
           const created = await createUser(newUser);
           setUsers([...users, created]);
           setNewUser({ name: "", email: "", password: "", role: "user" });
         } catch (err) {
           setError("Failed to create user");
         }
       };

       const handleDeleteUser = async (userId: string) => {
         try {
           await deleteUser(userId);
           setUsers(users.filter((user) => user.id !== userId));
         } catch (err) {
           setError("Failed to delete user");
         }
       };

       return (
         <div className="p-6">
           <h2 className="text-2xl mb-4">Admin Panel</h2>
           {loading && <p>Loading...</p>}
           {error && <p className="text-red-500 mb-4">{error}</p>}
           {!loading && !error && (
             <>
               <h3 className="text-xl mb-2">Experts</h3>
               <Table className="mb-6">
                 <TableHeader>
                   <TableRow>
                     <TableHead>Name</TableHead>
                     <TableHead>Affiliation</TableHead>
                     <TableHead>Available</TableHead>
                   </TableRow>
                 </TableHeader>
                 <TableBody>
                   {experts.map((expert) => (
                     <TableRow key={expert.id}>
                       <TableCell>{expert.name}</TableCell>
                       <TableCell>{expert.affiliation}</TableCell>
                       <TableCell>{expert.is_available ? "Yes" : "No"}</TableCell>
                     </TableRow>
                   ))}
                 </TableBody>
               </Table>

               <h3 className="text-xl mb-2">Requests</h3>
               <Table className="mb-6">
                 <TableHeader>
                   <TableRow>
                     <TableHead>User</TableHead>
                     <TableHead>Expert Name</TableHead>
                     <TableHead>Institution</TableHead>
                     <TableHead>Status</TableHead>
                     <TableHead>Actions</TableHead>
                   </TableRow>
                 </TableHeader>
                 <TableBody>
                   {requests.map((request) => (
                     <TableRow key={request.id}>
                       <TableCell>{request.user_id}</TableCell>
                       <TableCell>{request.name}</TableCell>
                       <TableCell>{request.institution}</TableCell>
                       <TableCell>{request.status}</TableCell>
                       <TableCell>
                         {request.status === "pending" && (
                           <>
                             <Button onClick={() => handleApprove(request.id)} className="mr-2">
                               Approve
                             </Button>
                             <Button
                               variant="destructive"
                               onClick={() => {
                                 const reason = prompt("Enter rejection reason:");
                                 if (reason) handleReject(request.id, reason);
                               }}
                             >
                               Reject
                             </Button>
                           </>
                         )}
                       </TableCell>
                     </TableRow>
                   ))}
                 </TableBody>
               </Table>

               <h3 className="text-xl mb-2">Users</h3>
               <form onSubmit={handleCreateUser} className="mb-6 space-y-4">
                 <Input
                   placeholder="Name"
                   value={newUser.name}
                   onChange={(e) => setNewUser({ ...newUser, name: e.target.value })}
                 />
                 <Input
                   placeholder="Email"
                   type="email"
                   value={newUser.email}
                   onChange={(e) => setNewUser({ ...newUser, email: e.target.value })}
                 />
                 <Input
                   placeholder="Password"
                   type="password"
                   value={newUser.password}
                   onChange={(e) => setNewUser({ ...newUser, password: e.target.value })}
                 />
                 <Select
                   value={newUser.role}
                   onValueChange={(value) => setNewUser({ ...newUser, role: value })}
                 >
                   <SelectTrigger>
                     <SelectValue placeholder="Role" />
                   </SelectTrigger>
                   <SelectContent>
                     <SelectItem value="user">User</SelectItem>
                     <SelectItem value="admin">Admin</SelectItem>
                   </SelectContent>
                 </Select>
                 <Button type="submit">Create User</Button>
               </form>
               <Table>
                 <TableHeader>
                   <TableRow>
                     <TableHead>Name</TableHead>
                     <TableHead>Email</TableHead>
                     <TableHead>Role</TableHead>
                     <TableHead>Requests Submitted</TableHead>
                     <TableHead>Actions</TableHead>
                   </TableRow>
                 </TableHeader>
                 <TableBody>
                   {users.map((user) => (
                     <TableRow key={user.id}>
                       <TableCell>{user.name}</TableCell>
                       <TableCell>{user.email}</TableCell>
                       <TableCell>{user.role}</TableCell>
                       <TableCell>{userStats[user.id] || 0}</TableCell>
                       <TableCell>
                         <Button variant="destructive" onClick={() => handleDeleteUser(user.id)}>
                           Delete
                         </Button>
                       </TableCell>
                     </TableRow>
                   ))}
                 </TableBody>
               </Table>
             </>
           )}
         </div>
       );
     }
     ```
   - **Verify**: Run `npx eslint src/pages/Admin.tsx`—should pass.

## Common Pitfalls and How to Avoid Them

1. **Admin Access**:
   - **Pitfall**: `/admin` redirects if `role` isn’t `admin`.
   - **Fix**: Mock admin login in console if needed (`localStorage.setItem("token", "mock-token")`).

2. **API Errors**:
   - **Pitfall**: `/api/requests`, `/api/users`, etc., might fail if the backend isn’t running.
   - **Fix**: Focus on UI—error state will handle failures.

3. **Data Format**:
   - **Pitfall**: API response might not match `Request` or `User` interfaces.
   - **Fix**: Adjust interfaces if backend response differs.

## Final Checks
Provide the following instructions for the user to execute and verify the app’s state:
- Run `npm run dev` in the terminal.
- Log in as admin (mock: `localStorage.setItem("token", "mock-token")
