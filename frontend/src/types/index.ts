// User Types
export interface User {
  id: number;
  email: string;
  name: string;
  role: string;
  isActive: boolean;
  createdAt: string;
  lastLogin: string;
}

export interface AuthState {
  token: string | null;
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

// Expert Types
export interface Expert {
  id: number;
  name: string;
  affiliation: string;
  primaryContact: string;
  contactType: string;
  skills: string[];
  role: string;
  employmentType: string;
  generalArea: number;
  cvPath: string;
  biography: string;
  isBahraini: boolean;
  availability: string;
  rating: number;
  created_at: string;
  updated_at: string;
}

// Expert Request Types
export interface ExpertRequest {
  id: number;
  requestorId: number;
  requestorName: string;
  requestorEmail: string;
  organizationName: string;
  projectName: string;
  projectDescription: string;
  expertiseRequired: string;
  timeframe: string;
  status: string;
  notes: string;
  createdAt: string;
  updatedAt: string;
}

// Document Types
export interface Document {
  id: number;
  expertId: number;
  filename: string;
  originalFilename: string;
  documentType: string;
  contentType: string;
  size: number;
  uploadedBy: number;
  uploadedAt: string;
}

// Statistics Types
export interface NationalityStats {
  bahraini: number;
  international: number;
  percentage: number;
}

export interface GrowthStats {
  month: string;
  newExperts: number;
  totalExperts: number;
}

// Engagement Types
export interface Engagement {
  id: number;
  expertId: number;
  requestId: number | null;
  title: string;
  description: string;
  engagementType: string;
  status: string;
  startDate: string;
  endDate: string;
  contactPerson: string;
  contactEmail: string;
  organizationName: string;
  notes: string;
  createdAt: string;
  updatedAt: string;
}

// API Response Types
export interface ApiResponse<T> {
  data: T;
  message: string;
  success: boolean;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  limit: number;
  totalPages: number;
}