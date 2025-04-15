import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { 
  ApiResponse, User, Expert, ExpertRequest, 
  NationalityStats, GrowthStats,
  PaginatedResponse
} from '../types';

// Check if we're in debug mode
const isDebugMode = import.meta.env.VITE_DEBUG_MODE === 'true';

// Create axios instance with default config
const api = axios.create({
  baseURL: import.meta.env.VITE_API_URL || '/api',
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: false, // Set to true only if needed for cookie-based auth
  timeout: 10000, // Add a 10-second timeout
});

console.log('API baseURL:', import.meta.env.VITE_API_URL || '/api');
console.log('Debug mode:', isDebugMode ? 'enabled' : 'disabled');

// Request interceptor to add auth token
api.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token) {
      config.headers.Authorization = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

// Response interceptor for error handling
api.interceptors.response.use(
  (response) => response,
  (error: AxiosError) => {
    console.error('API Error:', error.message);
    console.error('Response data:', error.response?.data);
    console.error('Status:', error.response?.status);
    
    // Handle 401 Unauthorized - redirect to login
    if (error.response?.status === 401) {
      localStorage.removeItem('token');
      localStorage.removeItem('user');
      window.location.href = '/';
    }
    
    // Handle CORS errors
    if (error.message.includes('Network Error') || error.message.includes('CORS')) {
      console.error('Possible CORS issue - check that backend allows cross-origin requests');
    }
    
    return Promise.reject(error);
  }
);

const request = async <T>(config: AxiosRequestConfig): Promise<ApiResponse<T>> => {
  try {
    if (isDebugMode) {
      console.log(`Making request: ${config.method} ${config.url}`, config.data || config.params || '');
    }
    
    const response = await api(config);
    
    if (isDebugMode) {
      console.log(`Response from ${config.url}:`, response.data);
    }
    
    return response.data;
  } catch (error) {
    if (isDebugMode) {
      console.error(`Request failed for ${config.url}:`, error);
    }
    
    if (axios.isAxiosError(error)) {
      const axiosError = error as AxiosError<ApiResponse<null>>;
      
      // Log specific error information based on status codes
      if (axiosError.response) {
        const status = axiosError.response.status;
        
        if (isDebugMode) {
          console.error(`Error ${status} for ${config.url}:`, axiosError.response.data);
        }
        
        // Return backend error response if available
        if (axiosError.response.data) {
          return axiosError.response.data;
        }
        
        // Generate appropriate error messages based on status
        let errorMessage = axiosError.message;
        if (status === 400) errorMessage = 'Invalid request data';
        else if (status === 403) errorMessage = 'Permission denied';
        else if (status === 404) errorMessage = 'Resource not found';
        else if (status === 500) errorMessage = 'Server error occurred';
        
        return {
          success: false,
          message: errorMessage,
          data: null as unknown as T,
        };
      }
      
      // Network error or timeout
      if (axiosError.code === 'ECONNABORTED') {
        return {
          success: false,
          message: 'Request timed out. Please try again.',
          data: null as unknown as T,
        };
      }
      
      if (!axiosError.response) {
        return {
          success: false,
          message: 'Network connection error. Please check your connection.',
          data: null as unknown as T,
        };
      }
      
      return {
        success: false,
        message: axiosError.message,
        data: null as unknown as T,
      };
    }
    
    // Non-Axios error (should be rare)
    if (isDebugMode) {
      console.error('Unexpected non-Axios error:', error);
    }
    
    return {
      success: false,
      message: 'An unexpected error occurred',
      data: null as unknown as T,
    };
  }
};

export const authApi = {
  login: async (email: string, password: string) => {
    console.log('Sending login request to:', `${api.defaults.baseURL}/api/auth/login`);
    if (isDebugMode) {
      console.log('Login data:', { email });
    }
    
    try {
      const response = await api({
        url: '/api/auth/login',
        method: 'POST',
        data: { email, password },
      });
      
      return response.data;
    } catch (error) {
      console.error('Login error in API service:', error);
      if (axios.isAxiosError(error) && error.response) {
        return error.response.data;
      }
      
      return {
        success: false,
        message: 'Failed to connect to authentication service',
      };
    }
  },

  logout: () => {
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    localStorage.removeItem('auth_error');
  },
};

// Experts API
export const expertsApi = {
  getExperts: (page: number = 1, limit: number = 10, params?: Record<string, string | boolean>) => 
    request<PaginatedResponse<Expert>>({
      url: '/experts',
      method: 'GET',
      params: {
        ...params,
        page,
        limit
      },
    }),

  getExpertById: (id: string) => 
    request<Expert>({
      url: `/experts/${id}`,
      method: 'GET',
    }),

  createExpert: (data: FormData) => 
    request<Expert>({
      url: '/experts',
      method: 'POST',
      data,
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }),

  updateExpert: (id: string, data: FormData) => 
    request<Expert>({
      url: `/experts/${id}`,
      method: 'PUT',
      data,
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }),
    
  deleteExpert: (id: string) => 
    request<void>({
      url: `/experts/${id}`,
      method: 'DELETE',
    }),

  downloadExpertPdf: (id: string) => 
    api({
      url: `/experts/${id}/approval-pdf`,
      method: 'GET',
      responseType: 'blob',
    }),
};

// Expert Requests API
export const expertRequestsApi = {
  getExpertRequests: (page: number = 1, limit: number = 10, params?: Record<string, string | boolean>) => 
    request<PaginatedResponse<ExpertRequest>>({
      url: '/expert-requests',
      method: 'GET',
      params: {
        ...params,
        page,
        limit
      },
    }),

  getExpertRequestById: (id: string) => 
    request<ExpertRequest>({
      url: `/expert-requests/${id}`,
      method: 'GET',
    }),

  createExpertRequest: (data: FormData) => 
    request<ExpertRequest>({
      url: '/expert-requests',
      method: 'POST',
      data,
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }),

  updateExpertRequest: (id: string, data: Partial<ExpertRequest>) => 
    request<ExpertRequest>({
      url: `/expert-requests/${id}`,
      method: 'PUT',
      data,
    }),

  deleteExpertRequest: (id: string) => 
    request<void>({
      url: `/expert-requests/${id}`,
      method: 'DELETE',
    }),
};

// Users API
export const usersApi = {
  getUsers: (page: number = 1, limit: number = 10, params?: Record<string, string | boolean>) => 
    request<PaginatedResponse<User>>({
      url: '/users',
      method: 'GET',
      params: {
        ...params,
        page,
        limit
      },
    }),

  getUserById: (id: string) => 
    request<User>({
      url: `/users/${id}`,
      method: 'GET',
    }),

  createUser: (data: Partial<User>) => 
    request<User>({
      url: '/users',
      method: 'POST',
      data,
    }),

  updateUser: (id: string, data: Partial<User>) => 
    request<User>({
      url: `/users/${id}`,
      method: 'PUT',
      data,
    }),

  deleteUser: (id: string) => 
    request<void>({
      url: `/users/${id}`,
      method: 'DELETE',
    }),
};

// Statistics API
export const statisticsApi = {
  getNationalityStats: () => 
    request<NationalityStats>({
      url: '/statistics/nationality',
      method: 'GET',
    }),

  getGrowthStats: (months?: number) => 
    request<GrowthStats[]>({
      url: '/statistics/growth',
      method: 'GET',
      params: { months },
    }),

  getOverallStats: () => 
    request<{
      totalExperts: number;
      totalBahraini: number;
      totalInternational: number;
      totalEngagements: number;
      byEmploymentType: Record<string, number>;
      byAvailability: Record<string, number>;
    }>({
      url: '/statistics',
      method: 'GET',
    }),

  getEngagementStats: () => 
    request<{
      total: number;
      byStatus: Record<string, number>;
      byType: Record<string, number>;
    }>({
      url: '/statistics/engagements',
      method: 'GET',
    }),
};

// Expert Areas API
export const expertAreasApi = {
  getExpertAreas: () => 
    request<{
      id: number;
      name: string;
      description: string;
    }[]>({
      url: '/expert/areas',
      method: 'GET',
    }),
};

// Document API
export const documentApi = {
  uploadDocument: (data: FormData) => 
    request<{
      id: number;
      expertId: number;
      filename: string;
      originalFilename: string;
      documentType: string;
      contentType: string;
      size: number;
      uploadedBy: number;
      uploadedAt: string;
    }>({
      url: '/documents',
      method: 'POST',
      data,
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }),

  getDocument: (id: number) => 
    request<Document>({
      url: `/documents/${id}`,
      method: 'GET',
    }),

  deleteDocument: (id: number) => 
    request<void>({
      url: `/documents/${id}`,
      method: 'DELETE',
    }),

  getExpertDocuments: (expertId: number) => 
    request<Document[]>({
      url: `/experts/${expertId}/documents`,
      method: 'GET',
    }),
};

// Engagement API
export const engagementApi = {
  getEngagements: (page: number = 1, limit: number = 10, params?: Record<string, string | boolean>) => 
    request<PaginatedResponse<Engagement>>({
      url: '/engagements',
      method: 'GET',
      params: {
        ...params,
        page,
        limit
      },
    }),

  getEngagementById: (id: string) => 
    request<Engagement>({
      url: `/engagements/${id}`,
      method: 'GET',
    }),

  createEngagement: (data: Partial<Engagement>) => 
    request<Engagement>({
      url: '/engagements',
      method: 'POST',
      data,
    }),

  updateEngagement: (id: string, data: Partial<Engagement>) => 
    request<Engagement>({
      url: `/engagements/${id}`,
      method: 'PUT',
      data,
    }),

  deleteEngagement: (id: string) => 
    request<void>({
      url: `/engagements/${id}`,
      method: 'DELETE',
    }),

  getExpertEngagements: (expertId: string, page: number = 1, limit: number = 10) => 
    request<PaginatedResponse<Engagement>>({
      url: `/experts/${expertId}/engagements`,
      method: 'GET',
      params: {
        page,
        limit
      },
    }),
};

export default {
  auth: authApi,
  experts: expertsApi,
  expertRequests: expertRequestsApi,
  users: usersApi,
  statistics: statisticsApi,
  expertAreas: expertAreasApi,
  documents: documentApi,
  engagements: engagementApi,
};