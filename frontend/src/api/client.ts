import axios, { AxiosError, AxiosRequestConfig } from 'axios';
import { ApiResponse } from '../types';

// Check if we're in debug mode
const isDebugMode = import.meta.env?.VITE_DEBUG_MODE === 'true' || false;

// Create and export the API client
export const createApiClient = (baseURL = '/api') => {
  const client = axios.create({
    baseURL,
    headers: {
      'Content-Type': 'application/json',
    },
    withCredentials: false, // Set to true only if needed for cookie-based auth
    timeout: 10000, // Add a 10-second timeout
  });
  
  // Debug logging
  if (isDebugMode) {
    console.log('API baseURL:', baseURL);
    console.log('Debug mode: enabled');
    
    // Add more debug info
    client.interceptors.request.use(function (config) {
      console.log('API CLIENT: Making request:', config.method, config.url, config.data || config.params || '');
      return config;
    });

    client.interceptors.response.use(function (response) {
      console.log('API CLIENT: Response received:', response.status, response.config.url, response.data);
      return response;
    });
  }

  // Request interceptor to add auth token
  client.interceptors.request.use(
    (config) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.set('Authorization', `Bearer ${token}`);
      }
      return config;
    },
    (error: unknown) => Promise.reject(error)
  );

  // Response interceptor for error handling
  client.interceptors.response.use(
    (response) => response,
    (error: AxiosError) => {
      if (isDebugMode) {
        console.error('API Error:', error.message);
        console.error('Response data:', error.response?.data);
        console.error('Status:', error.response?.status);
      }
      
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

  return client;
};

// Default client instance
export const apiClient = createApiClient('');

// Generic request function for standard API responses
export const request = async <T>(config: AxiosRequestConfig): Promise<ApiResponse<T>> => {
  try {
    console.log(`API CLIENT: Making request: ${config.method} ${config.url}`, config.data || config.params || '');
    
    const response = await apiClient(config);
    
    console.log(`API CLIENT: Raw response from ${config.url}:`, response);
    console.log(`API CLIENT: Response data:`, response.data);
    console.log(`API CLIENT: Response status:`, response.status);
    
    return response.data;
  } catch (error) {
    console.error(`API CLIENT: Request failed for ${config.url}:`, error);
    
    if (axios.isAxiosError(error)) {
      const axiosError = error as AxiosError<ApiResponse<null>>;
      
      // Log specific error information based on status codes
      if (axiosError.response) {
        const status = axiosError.response.status;
        
        console.error(`API CLIENT: Error ${status} for ${config.url}:`, axiosError.response.data);
        
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
          data: null,
        };
      }
      
      // Network error or timeout
      if (axiosError.code === 'ECONNABORTED') {
        return {
          success: false,
          message: 'Request timed out. Please try again.',
          data: null,
        };
      }
      
      if (!axiosError.response) {
        return {
          success: false,
          message: 'Network connection error. Please check your connection.',
          data: null,
        };
      }
      
      return {
        success: false,
        message: axiosError.message,
        data: null,
      };
    }
    
    // Non-Axios error (should be rare)
    if (isDebugMode) {
      console.error('Unexpected non-Axios error:', error);
    }
    
    return {
      success: false,
      message: 'An unexpected error occurred',
      data: null,
    };
  }
};