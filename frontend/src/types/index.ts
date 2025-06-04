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
  expertId: string;
  name: string;
  designation: string;
  institution: string;
  isBahraini: boolean;
  isAvailable: boolean;
  rating: string;
  role: string;
  employmentType: string;
  generalArea: number;
  generalAreaName: string;
  specializedArea: string;
  isTrained: boolean;
  cvPath: string;
  phone: string;
  email: string;
  isPublished: boolean;
  biography: string;
  approvalDocumentPath: string;
  skills: string[];
  createdAt: string;
  updatedAt: string;
  created_at?: string; // For API compatibility
  updated_at?: string; // For API compatibility
  
  // Computed properties
  availability?: 'Available' | 'Limited' | 'Unavailable';
  primaryContact?: string;
  contactType?: string;
  affiliation?: string; // Alias for institution
}

// Expert Request Types
export interface ExpertRequest {
  id: number;
  name: string;
  status: string;
  cvPath: string;
  approvalDocumentPath: string;
  designation: string;
  institution: string;
  isBahraini: boolean;
  isAvailable: boolean;
  rating: string;
  role: string;
  employmentType: string;
  generalArea: number;
  specializedArea: string;
  isTrained: boolean;
  phone: string;
  email: string;
  biography: string;
  skills: string[];
  isPublished: boolean;
  createdAt: string;
  updatedAt: string;
  rejectionReason?: string;
  createdBy?: number;
  reviewedBy?: number;
  reviewedAt?: string;
  expertId?: string;
}

// Document Types
export interface Document {
  id: number;
  expertId: number;
  documentType: 'cv' | 'certificate' | 'publication' | 'research' | 'other';
  filePath: string;
  originalFilename?: string;
  filename?: string; // For compatibility
  contentType?: string;
  size?: number;
  uploadedAt?: string;
  createdAt: string;
}

// Statistics Types
export interface StatItem {
  name: string;
  count: number;
  percentage: number;
}

export interface NationalityStats {
  total: number;
  stats: StatItem[];
  bahraini?: number;
  international?: number;
  percentage?: number;
}

export interface GrowthStats {
  period: string;
  count: number;
  growthRate: number;
}

export interface AreaStats {
  generalAreas: StatItem[];
  topSpecializedAreas: StatItem[];
  bottomSpecializedAreas: StatItem[];
}

export interface EngagementStats {
  total: number;
  byType: StatItem[];
  byStatus: StatItem[];
}

export interface ExpertStats {
  totalExperts: number;
  activeCount: number;
  bahrainiPercentage: number;
  publishedCount: number;
  publishedRatio: number;
  topAreas: StatItem[];
  engagementsByType: StatItem[];
  yearlyGrowth: GrowthStats[];
  mostRequestedExperts: {
    expertId: string;
    name: string;
    count: number;
  }[];
  lastUpdated: string;
}

// Phase Types
export interface PhaseApplication {
  id: number;
  phaseId: number;
  type: string;
  institutionName: string;
  qualificationName: string;
  expert1: number;
  expert1Name: string;
  expert2: number;
  expert2Name: string;
  status: string;
  rejectionNotes: string;
  createdAt: string;
  updatedAt: string;
}

export interface Phase {
  id: number;
  phaseId: string;
  title: string;
  assignedSchedulerId: number;
  schedulerName: string;
  status: string;
  createdAt: string;
  updatedAt: string;
  applications: PhaseApplication[];
}

// Engagement Types
export interface Engagement {
  id: number;
  expertId: number;
  expertName: string;
  title: string;
  description: string;
  engagementType: 'consulting' | 'consultation' | 'training' | 'mentoring' | 'research' | 'project' | 'workshop' | 'other';
  startDate: string;
  endDate: string;
  projectName: string;
  organizationName: string;
  contactPerson: string;
  contactEmail: string;
  status: 'pending' | 'confirmed' | 'in_progress' | 'completed' | 'cancelled';
  notes?: string;
  requestId?: number;
  createdAt: string;
}

// API Response Types
export interface ApiResponse<T> {
  success: boolean;
  message: string;
  data: T;
}

export interface PaginatedResponse<T> {
  data: T[];
  page: number;
  limit: number;
  total: number;
  totalPages: number;
}

export interface PaginationInfo {
  totalCount: number;
  totalPages: number;
  currentPage: number;
  pageSize: number;
  hasNextPage: boolean;
  hasPrevPage: boolean;
  hasMore?: boolean;
}

// Standard paginated responses for different entity types
export interface ExpertListResponse {
  experts: Expert[];
  pagination: PaginationInfo;
}

export interface RequestListResponse {
  requests: ExpertRequest[];
  pagination: PaginationInfo;
}

export interface EngagementListResponse {
  engagements: Engagement[];
  pagination: PaginationInfo;
  filters: {
    expertId?: number;
    type?: string;
  };
}

export interface PhaseListResponse {
  phases: Phase[];
  pagination: PaginationInfo;
  filters: {
    status?: string;
    schedulerId?: number;
  };
}

export interface DocumentListResponse {
  expertId: number;
  count: number;
  documents: Document[];
}

export interface BatchApproveResponse {
  results: {
    id: number;
    status: 'success' | 'failed';
    error?: string;
  }[];
  approvedIds: number[];
  errorCount: number;
}