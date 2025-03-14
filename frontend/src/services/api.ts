import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { 
  ApiResponse, User, Expert, ExpertRequest, 
  NationalityStats, GrowthStats, IscedStats 
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

// Generic request function
const request = async <T>(config: AxiosRequestConfig): Promise<ApiResponse<T>> => {
  try {
    console.log(`Making request: ${config.method} ${config.url}`, config.data || config.params || '');
    const response = await api(config);
    console.log(`Response from ${config.url}:`, response.data);
    return response.data;
  } catch (error) {
    console.error(`Request failed for ${config.url}:`, error);
    
    if (axios.isAxiosError(error)) {
      const axiosError = error as AxiosError<ApiResponse<null>>;
      if (axiosError.response?.data) {
        console.error('Error response data:', axiosError.response.data);
        return axiosError.response.data;
      }
      
      // Network error or no response
      if (axiosError.code === 'ECONNABORTED' || !axiosError.response) {
        console.error('Network error or timeout');
      }
      
      return {
        success: false,
        message: axiosError.message,
        data: null as unknown as T,
      };
    }
    
    // Non-Axios error (should be rare)
    console.error('Unexpected non-Axios error:', error);
    return {
      success: false,
      message: 'An unexpected error occurred',
      data: null as unknown as T,
    };
  }
};

// Auth API
export const authApi = {
  login: async (email: string, password: string) => {
    console.log('Sending login request to:', `${import.meta.env.VITE_API_URL}/api/auth/login`);
    console.log('Login data:', { email, password });
    
    try {
      // Use direct axios call to get the raw response
      const response = await api({
        url: '/api/auth/login',
        method: 'POST',
        data: { email, password },
      });
      
      // Return the raw data as it comes from the backend
      return response.data;
    } catch (error) {
      console.error('Login error in API service:', error);
      if (axios.isAxiosError(error) && error.response) {
        // Return the error response data
        return error.response.data;
      }
      
      // Generic error
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
  getExperts: (params?: Record<string, string | boolean>) => 
    request<Expert[]>({
      url: '/api/experts',
      method: 'GET',
      params,
    }),

  getExpertById: (id: string) => 
    request<Expert>({
      url: `/api/experts/${id}`,
      method: 'GET',
    }),

  downloadExpertPdf: (id: string) => 
    api({
      url: `/api/experts/${id}/approval-pdf`,
      method: 'GET',
      responseType: 'blob',
    }),
};

// Expert Requests API
export const expertRequestsApi = {
  getExpertRequests: (params?: Record<string, string | boolean>) => 
    request<ExpertRequest[]>({
      url: '/api/expert-requests',
      method: 'GET',
      params,
    }),

  getExpertRequestById: (id: string) => 
    request<ExpertRequest>({
      url: `/api/expert-requests/${id}`,
      method: 'GET',
    }),

  createExpertRequest: (data: FormData) => 
    request<ExpertRequest>({
      url: '/api/expert-requests',
      method: 'POST',
      data,
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    }),

  updateExpertRequest: (id: string, data: Partial<ExpertRequest>) => 
    request<ExpertRequest>({
      url: `/api/expert-requests/${id}`,
      method: 'PUT',
      data,
    }),

  deleteExpertRequest: (id: string) => 
    request<void>({
      url: `/api/expert-requests/${id}`,
      method: 'DELETE',
    }),
};

// Users API
export const usersApi = {
  getUsers: (params?: Record<string, string | boolean>) => 
    request<User[]>({
      url: '/api/users',
      method: 'GET',
      params,
    }),

  getUserById: (id: string) => 
    request<User>({
      url: `/api/users/${id}`,
      method: 'GET',
    }),

  createUser: (data: Partial<User>) => 
    request<User>({
      url: '/api/users',
      method: 'POST',
      data,
    }),

  updateUser: (id: string, data: Partial<User>) => 
    request<User>({
      url: `/api/users/${id}`,
      method: 'PUT',
      data,
    }),

  deleteUser: (id: string) => 
    request<void>({
      url: `/api/users/${id}`,
      method: 'DELETE',
    }),
};

// Statistics API
export const statisticsApi = {
  getNationalityStats: () => 
    request<{ stats: NationalityStats[] }>({
      url: '/api/statistics/nationality',
      method: 'GET',
    }),

  getGrowthStats: () => 
    request<GrowthStats[]>({
      url: '/api/statistics/growth',
      method: 'GET',
    }),

  getIscedStats: () => 
    request<IscedStats[]>({
      url: '/api/statistics/isced',
      method: 'GET',
    }),
};

export default {
  auth: authApi,
  experts: expertsApi,
  expertRequests: expertRequestsApi,
  users: usersApi,
  statistics: statisticsApi,
};