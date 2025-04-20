import { request } from './client';
import { User } from '../types';

interface LoginResponse {
  token: string;
  user: User;
}

export const login = async (email: string, password: string) => {
  console.log('Sending login request');
  
  try {
    const response = await request<LoginResponse>({
      url: '/api/auth/login',
      method: 'POST',
      data: { email, password },
    });
    
    return response;
  } catch (error) {
    console.error('Login error in API service:', error);
    return {
      success: false,
      message: 'Failed to connect to authentication service',
      data: null,
    };
  }
};

export const logout = () => {
  localStorage.removeItem('token');
  localStorage.removeItem('user');
  localStorage.removeItem('auth_error');
};

export const refreshToken = async () => {
  try {
    const response = await request<{ token: string }>({
      url: '/api/auth/refresh',
      method: 'POST',
    });
    
    if (response.success && response.data?.token) {
      localStorage.setItem('token', response.data.token);
      return response;
    }
    
    return {
      success: false,
      message: 'Failed to refresh token',
      data: null,
    };
  } catch (error) {
    console.error('Token refresh error:', error);
    return {
      success: false,
      message: 'Failed to refresh authentication token',
      data: null,
    };
  }
};