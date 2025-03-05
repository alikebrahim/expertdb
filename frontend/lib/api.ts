import axios from 'axios';

// Types based on backend structures
export interface Engagement {
  id: number;
  expertId: number;
  engagementType: string;
  startDate: string;
  endDate?: string;
  projectName?: string;
  status: string;
  feedbackScore?: number;
  notes?: string;
  createdAt: string;
}

export interface Expert {
  id: number;
  expertId: string;
  name: string;
  designation: string;
  institution: string;
  isBahraini: boolean;
  nationality: string;
  isAvailable: boolean;
  rating: string;
  role: string;
  employmentType: string;
  generalArea: string;
  specializedArea: string;
  isTrained: boolean;
  cvPath: string;
  phone: string;
  email: string;
  isPublished: boolean;
  iscedLevel?: {
    id: number;
    code: string;
    name: string;
    description?: string;
  };
  iscedField?: {
    id: number;
    broadCode: string;
    broadName: string;
    narrowCode?: string;
    narrowName?: string;
    detailedCode?: string;
    detailedName?: string;
    description?: string;
  };
  areas?: Array<{
    id: number;
    name: string;
  }>;
  engagements?: Engagement[];
  createdAt: string;
  updatedAt?: string;
}

export interface ExpertRequest {
  id?: number;
  expertId?: string;
  name: string;
  designation: string;
  institution: string;
  isBahraini: boolean;
  isAvailable: boolean;
  rating?: string;
  role?: string;
  employmentType?: string;
  generalArea?: string;
  specializedArea?: string;
  isTrained?: boolean;
  cvPath?: string;
  phone: string;
  email: string;
  isPublished?: boolean;
  status?: string;
  createdAt?: string;
  reviewedAt?: string;
  reviewedBy?: number;
}

// API endpoints
const API_URL = '/api';

// Get the API URL from environment or use default
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Axios instance
const apiClient = axios.create({
  baseURL: API_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  // Add timeout to prevent hanging requests
  timeout: 15000,
});

// Add Authorization header with JWT token if available
apiClient.interceptors.request.use(
  (config) => {
    if (typeof window !== 'undefined') {
      // Only add token if we're not on the login page
      if (window.location.pathname !== '/login') {
        const token = localStorage.getItem('token');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        } else if (window.location.pathname !== '/') {
          // If no token and not on login or homepage, redirect to login
          window.location.href = '/login';
          return Promise.reject(new Error('No authentication token found'));
        }
      }
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Expert API functions
export const expertAPI = {
  // Get all experts with optional filters, pagination and sorting
  getAllExperts: async (filters?: Record<string, any>, limit: number = 10, offset: number = 0, sortBy: string = 'name', sortOrder: string = 'asc') => {
    const params = new URLSearchParams();
    
    // Add pagination parameters
    params.append('limit', String(limit));
    params.append('offset', String(offset));
    
    // Add sorting parameters
    if (sortBy) {
      params.append('sort_by', sortBy);
      params.append('sort_order', sortOrder);
    }
    
    // Add filters
    if (filters) {
      Object.entries(filters).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== '') {
          params.append(key, String(value));
        }
      });
    }
    
    const response = await apiClient.get(`/experts?${params.toString()}`);
    return {
      experts: response.data,
      pagination: {
        limit,
        offset,
        total: parseInt(response.headers['x-total-count'] || '0'), // Get total count from header
      }
    };
  },
  
  // Get expert by ID
  getExpertById: async (id: number): Promise<Expert> => {
    const response = await apiClient.get(`/experts/${id}`);
    return response.data;
  },
  
  // Create expert request
  createExpertRequest: async (expertRequest: ExpertRequest) => {
    const response = await apiClient.post('/expert-requests', expertRequest);
    return response.data;
  },
  
  // Get ISCED classification data
  getISCEDLevels: async () => {
    const response = await apiClient.get('/isced/levels');
    return response.data;
  },
  
  getISCEDFields: async () => {
    const response = await apiClient.get('/isced/fields');
    return response.data;
  },
  
  // Engagement operations
  getExpertEngagements: async (expertId: number): Promise<Engagement[]> => {
    const response = await apiClient.get(`/experts/${expertId}/engagements`);
    return response.data;
  },
  
  createEngagement: async (engagement: Omit<Engagement, 'id' | 'createdAt'>) => {
    const response = await apiClient.post('/engagements', engagement);
    return response.data;
  },
  
  updateEngagement: async (id: number, engagement: Partial<Engagement>) => {
    const response = await apiClient.put(`/engagements/${id}`, engagement);
    return response.data;
  },
  
  deleteEngagement: async (id: number) => {
    const response = await apiClient.delete(`/engagements/${id}`);
    return response.data;
  }
};

// AI Panel Suggestion API
export const aiAPI = {
  suggestExpertPanel: async (projectName: string, iscedFieldId?: number, numExperts: number = 3) => {
    const response = await apiClient.post('/ai/suggest-panel', {
      projectName,
      iscedFieldId: iscedFieldId || undefined,
      numExperts
    });
    return response.data;
  }
};

// Statistics API types
export interface AreaStat {
  name: string;
  count: number;
  percentage: number;
}

export interface GrowthStat {
  period: string;
  count: number;
  growthRate: number;
}

export interface ExpertStat {
  expertId: string;
  name: string;
  count: number;
}

export interface NationalityStats {
  total: number;
  bahraini: {
    count: number;
    percentage: number;
  };
  nonBahraini: {
    count: number;
    percentage: number;
  };
}

export interface Statistics {
  totalExperts: number;
  bahrainiPercentage: number;
  topAreas: AreaStat[];
  expertsByISCEDField: AreaStat[];
  engagementsByType: AreaStat[];
  monthlyGrowth: GrowthStat[];
  mostRequestedExperts: ExpertStat[];
  lastUpdated: string;
}

// Statistics API
export const statisticsAPI = {
  // Get all statistics
  getAllStatistics: async (): Promise<Statistics> => {
    const response = await apiClient.get('/statistics');
    return response.data;
  },
  
  // Get nationality statistics
  getNationalityStats: async (): Promise<NationalityStats> => {
    const response = await apiClient.get('/statistics/nationality');
    return response.data;
  },
  
  // Get ISCED statistics
  getISCEDStats: async (): Promise<AreaStat[]> => {
    const response = await apiClient.get('/statistics/isced');
    return response.data;
  },
  
  // Get engagement statistics
  getEngagementStats: async (): Promise<AreaStat[]> => {
    const response = await apiClient.get('/statistics/engagements');
    return response.data;
  },
  
  // Get growth statistics
  getGrowthStats: async (months: number = 12): Promise<GrowthStat[]> => {
    const response = await apiClient.get(`/statistics/growth?months=${months}`);
    return response.data;
  }
};

// Authentication API
export const authAPI = {
  login: async (email: string, password: string) => {
    try {
      // Try direct URL if API proxy fails
      const fallbackClient = axios.create({
        baseURL: API_BASE_URL,
        headers: { 'Content-Type': 'application/json' },
        timeout: 15000
      });
      
      // First try with the regular API client (using proxy)
      try {
        const response = await apiClient.post('/auth/login', { email, password });
        return response.data;
      } catch (proxyError) {
        // Fallback to direct URL if proxy fails
        const directResponse = await fallbackClient.post('/api/auth/login', { email, password });
        return directResponse.data;
      }
    } catch (error) {
      throw error;
    }
  },
  
  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    // Trigger storage event for cross-tab communication
    window.dispatchEvent(new Event('storage'));
  },
  
  getUser: () => {
    if (typeof window === 'undefined') return null;
    try {
      const user = localStorage.getItem('user');
      return user ? JSON.parse(user) : null;
    } catch (error) {
      console.error('Error parsing user from localStorage:', error);
      return null;
    }
  },
  
  isAuthenticated: () => {
    if (typeof window === 'undefined') return false;
    return !!localStorage.getItem('token');
  }
};

// Error handler middleware for Axios
apiClient.interceptors.response.use(
  (response) => response,
  (error) => {
    // Log detailed error information for debugging
    if (error.response) {
      // The request was made and the server responded with a status code
      // that falls out of the range of 2xx
      console.error('API Error Response:', {
        status: error.response.status,
        headers: error.response.headers,
        data: error.response.data,
      });
    } else if (error.request) {
      // The request was made but no response was received
      console.error('API No Response:', error.request);
    } else {
      // Something happened in setting up the request that triggered an Error
      console.error('API Request Error:', error.message);
    }
    
    // Handle authentication errors
    if (error.response?.status === 401 && typeof window !== 'undefined') {
      console.log('Authentication error detected, logging out user');
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      
      // Use a more controlled redirect that doesn't interfere with the current operation
      setTimeout(() => {
        window.location.href = '/login';
      }, 100);
    }
    
    return Promise.reject(error);
  }
);