// User Types
export interface User {
  id: string;
  email: string;
  name: string;
  role: 'admin' | 'user';
  isActive: boolean;
  createdAt?: string;
  lastLogin?: string;
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
  id: string;
  name: string;
  affiliation: string;
  role: string;
  type: string;
  specialization: string;
  isced: string;
  nationality: string;
  status: 'available' | 'unavailable';
  biography?: string;
}

// Expert Request Types
export interface ExpertRequest {
  id: string;
  name: string;
  affiliation: string;
  role: string;
  type: string;
  specialization: string;
  isced: string;
  nationality: string;
  status: 'pending' | 'approved' | 'rejected';
  rejectionReason?: string;
  userId: string;
  createdAt: string;
  updatedAt: string;
  documents?: Document[];
}

// Document Types
export interface Document {
  id: string;
  name: string;
  type: string;
  url: string;
  expertId?: string;
  expertRequestId?: string;
}

// Statistics Types
export interface NationalityStats {
  nationality: string;
  count: number;
}

export interface GrowthStats {
  year: number;
  count: number;
}

export interface IscedStats {
  isced: string;
  count: number;
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