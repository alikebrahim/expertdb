# Phase 3: Request Submission for ExpertDB Frontend

## Overview
You are an AI tasked with implementing request submission for the ExpertDB frontend using React, TypeScript, and shadcn/ui components. The project uses Vite, React Router, and Axios (in `src/api/api.ts`) to submit requests to `http://localhost:8080/api/expert-requests`. Use Mammoth for `.docx` parsing and jsPDF for PDF generation (both installed). The page is protected (requires login). Follow each step exactly, verify at each point, and avoid common pitfalls.

## Objective
- Create a `Requests` page with a form to select an expert, enter request details (`designation`, `name`, `institution`, `primaryContact`, `contactType`, `skills`, `biography`, `availability`), upload a CV (`.docx`), and generate a PDF client-side.
- Submit the request to `POST /api/expert-requests`.
- Use shadcn/ui components (`Select`, `Input`, `Button`, `Textarea`) for the form.

## Step-by-Step Instructions

1. **Update api.ts with Request Endpoints**
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

     export default api;
     ```
   - **Verify**: Run `npx eslint src/api/api.ts`—should pass.

2. **Create Requests Page**
   - **File**: Create `expertdb_grok/frontend/src/pages/Requests.tsx`.
   - **Code**:
     ```tsx
     import { useState, useEffect } from "react";
     import { useNavigate } from "react-router-dom";
     import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from "@/components/ui/select";
     import { Input } from "@/components/ui/input";
     import { Textarea } from "@/components/ui/textarea";
     import { Button } from "@/components/ui/button";
     import { Expert, getExperts, submitExpertRequest } from "../api/api";
     import mammoth from "mammoth";
     import { jsPDF } from "jspdf";

     export default function Requests() {
       const [experts, setExperts] = useState<Expert[]>([]);
       const [selectedExpert, setSelectedExpert] = useState("");
       const [requestDetails, setRequestDetails] = useState({
         designation: "",
         name: "",
         institution: "",
         primaryContact: "",
         contactType: "email",
         skills: [] as string[],
         biography: "",
         availability: "full-time",
       });
       const [skillsInput, setSkillsInput] = useState("");
       const [cvFile, setCvFile] = useState<File | null>(null);
       const [cvText, setCvText] = useState("");
       const [loading, setLoading] = useState(true);
       const [error, setError] = useState("");
       const navigate = useNavigate();

       useEffect(() => {
         const fetchExperts = async () => {
           try {
             const data = await getExperts();
             setExperts(data);
           } catch (err) {
             setError("Failed to load experts");
           } finally {
             setLoading(false);
           }
         };
         fetchExperts();
       }, []);

       const handleFileChange = async (e: React.ChangeEvent<HTMLInputElement>) => {
         const file = e.target.files?.[0];
         if (file) {
           setCvFile(file);
           const arrayBuffer = await file.arrayBuffer();
           const result = await mammoth.extractRawText({ arrayBuffer });
           setCvText(result.value);
         }
       };

       const generatePDF = () => {
         const doc = new jsPDF();
         doc.text(`Expert Request`, 10, 10);
         doc.text(`Designation: ${requestDetails.designation}`, 10, 20);
         doc.text(`Name: ${requestDetails.name}`, 10, 30);
         doc.text(`Institution: ${requestDetails.institution}`, 10, 40);
         doc.text(`Primary Contact: ${requestDetails.primaryContact} (${requestDetails.contactType})`, 10, 50);
         doc.text(`Skills: ${requestDetails.skills.join(", ")}`, 10, 60);
         doc.text(`Availability: ${requestDetails.availability}`, 10, 70);
         doc.text(`Biography: ${requestDetails.biography.slice(0, 200)}`, 10, 80);
         if (cvText) doc.text(`CV Content: ${cvText.slice(0, 200)}`, 10, 90);
         return doc.output("blob");
       };

       const handleSubmit = async (e: React.FormEvent) => {
         e.preventDefault();
         if (!selectedExpert) {
           setError("Please select an expert");
           return;
         }

         try {
           const selected = experts.find((expert) => expert.id === selectedExpert);
           if (!selected) throw new Error("Invalid expert selected");

           await submitExpertRequest({
             designation: requestDetails.designation,
             name: selected.name,
             institution: selected.affiliation,
             bio: requestDetails.biography,
             primaryContact: requestDetails.primaryContact,
             contactType: requestDetails.contactType,
             skills: requestDetails.skills,
             availability: requestDetails.availability,
             cv: cvFile || undefined,
           });

           // Generate PDF (client-side)
           const pdf = generatePDF();
           const pdfFile = new File([pdf], "request.pdf", { type: "application/pdf" });
           const formData = new FormData();
           formData.append("type", "request");
           formData.append("expert_id", selectedExpert);
           formData.append("file", pdfFile);

           await api.post("/api/documents", formData, {
             headers: { "Content-Type": "multipart/form-data" },
           });

           navigate("/search");
         } catch (err) {
           setError("Failed to submit request");
         }
       };

       const addSkill = () => {
         if (skillsInput.trim()) {
           setRequestDetails({
             ...requestDetails,
             skills: [...requestDetails.skills, skillsInput.trim()],
           });
           setSkillsInput("");
         }
       };

       return (
         <div className="p-6">
           <h2 className="text-2xl mb-4">Submit Expert Request</h2>
           {loading && <p>Loading...</p>}
           {error && <p className="text-red-500 mb-4">{error}</p>}
           {!loading && !error && (
             <form onSubmit={handleSubmit} className="space-y-4">
               <Select onValueChange={setSelectedExpert} value={selectedExpert}>
                 <SelectTrigger className="w-[300px]">
                   <SelectValue placeholder="Select an expert" />
                 </SelectTrigger>
                 <SelectContent>
                   {experts.map((expert) => (
                     <SelectItem key={expert.id} value={expert.id}>
                       {expert.name}
                     </SelectItem>
                   ))}
                 </SelectContent>
               </Select>
               <Input
                 placeholder="Designation"
                 value={requestDetails.designation}
                 onChange={(e) => setRequestDetails({ ...requestDetails, designation: e.target.value })}
                 className="max-w-md"
               />
               <Input
                 placeholder="Primary Contact (Email/Phone)"
                 value={requestDetails.primaryContact}
                 onChange={(e) => setRequestDetails({ ...requestDetails, primaryContact: e.target.value })}
                 className="max-w-md"
               />
               <Select
                 value={requestDetails.contactType}
                 onValueChange={(value) => setRequestDetails({ ...requestDetails, contactType: value })}
               >
                 <SelectTrigger className="w-[200px]">
                   <SelectValue placeholder="Contact Type" />
                 </SelectTrigger>
                 <SelectContent>
                   <SelectItem value="email">Email</SelectItem>
                   <SelectItem value="phone">Phone</SelectItem>
                 </SelectContent>
               </Select>
               <div className="flex space-x-2">
                 <Input
                   placeholder="Add a skill (e.g., JavaScript)"
                   value={skillsInput}
                   onChange={(e) => setSkillsInput(e.target.value)}
                   className="max-w-md"
                 />
                 <Button type="button" onClick={addSkill}>Add Skill</Button>
               </div>
               <div>
                 {requestDetails.skills.map((skill, index) => (
                   <span key={index} className="inline-block bg-gray-200 px-2 py-1 mr-2 mb-2 rounded">
                     {skill}
                   </span>
                 ))}
               </div>
               <Select
                 value={requestDetails.availability}
                 onValueChange={(value) => setRequestDetails({ ...requestDetails, availability: value })}
               >
                 <SelectTrigger className="w-[200px]">
                   <SelectValue placeholder="Availability" />
                 </SelectTrigger>
                 <SelectContent>
                   <SelectItem value="full-time">Full-Time</SelectItem>
                   <SelectItem value="part-time">Part-Time</SelectItem>
                   <SelectItem value="weekends">Weekends</SelectItem>
                 </SelectContent>
               </Select>
               <Textarea
                 placeholder="Biography"
                 value={requestDetails.biography}
                 onChange={(e) => setRequestDetails({ ...requestDetails, biography: e.target.value })}
                 className="max-w-md"
               />
               <Input type="file" accept=".doc,.docx" onChange={handleFileChange} />
               {cvText && <p className="text-sm text-gray-600">CV Preview: {cvText.slice(0, 100)}...</p>}
               <Button type="submit">Submit Request</Button>
             </form>
           )}
         </div>
       );
     }
     ```
   - **Verify**: Run `npx eslint src/pages/Requests.tsx`—should pass. Run `npm run dev`, visit `http://localhost:5174/requests`—should redirect to `/login` if not logged in.

3. **Add shadcn/ui Textarea Component**
   - **Command**: Run `npx shadcn@latest add textarea`.
   - **Verify**: Check `src/components/ui/textarea.tsx` exists.

## Common Pitfalls and How to Avoid Them

1. **Missing shadcn/ui Components**:
   - **Pitfall**: `Textarea` component might not exist.
   - **Fix**: Ensure `npx shadcn@latest add textarea` runs successfully.

2. **File Upload Issues**:
   - **Pitfall**: FormData might fail if `cvFile` is null.
   - **Fix**: Validate inputs before submitting (already in code).

3. **API Errors**:
   - **Pitfall**: `/api/expert-requests` might fail if the backend isn’t running.
   - **Fix**: Focus on UI—error state will handle failures.

4. **PDF Generation**:
   - **Pitfall**: jsPDF might throw errors if `cvText` is too long.
   - **Fix**: Truncate `cvText` in `generatePDF` (already done).

## Final Checks
Provide the following instructions for the user to execute and verify the app’s state:
- Run `npm run dev` in the terminal.
- Log in (mock if needed: `localStorage.setItem("token", "mock-token")` in the browser console).
- Visit `http://localhost:5174/requests`—should show a form with dropdowns, inputs, and file upload.
- Fill in details, upload a `.docx` file, and submit—should redirect to `/search` if successful, or show an error if the API fails.
- Run `npm run build` in the terminal—should succeed.

## Questions
- Should the form include validation for specific fields (e.g., email format for `primaryContact` if `contactType` is `email`)?
- What should the PDF layout look like (e.g., specific formatting, additional fields)?
- Should there be a confirmation dialog before submitting the request?
