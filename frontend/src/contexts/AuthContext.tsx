import { createContext, useEffect, useState, ReactNode } from 'react';
import { User, AuthState } from '../types';
import * as authApi from '../api/auth';

interface AuthContextType extends AuthState {
  login: (email: string, password: string) => Promise<boolean>;
  logout: () => void;
  refreshAuth: () => Promise<boolean>;
}

const initialState: AuthState = {
  token: null,
  user: null,
  isAuthenticated: false,
  isLoading: true,
  error: null,
};

export const AuthContext = createContext<AuthContextType>({
  ...initialState,
  login: async () => false,
  logout: () => {},
  refreshAuth: async () => false,
});

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [state, setState] = useState<AuthState>(initialState);
  const isDebugMode = import.meta.env.VITE_DEBUG_MODE === 'true';

  // Check for existing auth on mount
  useEffect(() => {
    const initializeAuth = async () => {
      try {
        const storedToken = localStorage.getItem('token');
        const storedUser = localStorage.getItem('user');
        
        if (storedToken && storedUser) {
          setState({
            token: storedToken,
            user: JSON.parse(storedUser) as User,
            isAuthenticated: true,
            isLoading: false,
            error: null,
          });
        } else {
          setState({
            ...initialState,
            isLoading: false,
          });
        }
      } catch {
        setState({
          ...initialState,
          isLoading: false,
          error: 'Failed to restore authentication state',
        });
      }
    };

    initializeAuth();
  }, []);

  // Login function
  const login = async (email: string, password: string): Promise<boolean> => {
    try {
      if (isDebugMode) {
        console.log('Login attempt:', { email });
        console.log('API URL:', import.meta.env.VITE_API_URL);
      }
      
      setState({
        ...state,
        isLoading: true,
        error: null,
      });

      const response = await authApi.login(email, password);
      
      if (isDebugMode) {
        console.log('Login response:', response);
      }
      
      // Check if the response includes token directly (backend format)
      if (response.data?.token && response.data?.user) {
        if (isDebugMode) {
          console.log('Direct token format detected - login successful');
        }
        
        const { token, user } = response.data;
        
        // Store auth data
        localStorage.setItem('token', token);
        localStorage.setItem('user', JSON.stringify(user));
        
        setState({
          token,
          user,
          isAuthenticated: true,
          isLoading: false,
          error: null,
        });
        
        return true;
      }
      else {
        if (isDebugMode) {
          console.log('Login failed with response:', response);
        }
        
        const errorMessage = response.message || 'Authentication failed';
        
        // Store error message for UI display
        localStorage.setItem('auth_error', errorMessage);
        
        setState({
          ...initialState,
          isLoading: false,
          error: errorMessage,
        });
        
        return false;
      }
    } catch (error) {
      console.error('Login error:', error);
      
      // Store detailed error for UI display
      let errorMessage = 'Authentication failed';
      if (error instanceof Error) {
        errorMessage = `Authentication error: ${error.message}`;
      }
      
      localStorage.setItem('auth_error', errorMessage);
      
      setState({
        ...initialState,
        isLoading: false,
        error: errorMessage,
      });
      
      return false;
    }
  };

  // Token refresh function
  const refreshAuth = async (): Promise<boolean> => {
    try {
      setState({
        ...state,
        isLoading: true,
      });
      
      const response = await authApi.refreshToken();
      
      if (response.success && response.data) {
        setState({
          ...state,
          token: response.data.token,
          isLoading: false,
        });
        return true;
      }
      
      return false;
    } catch {
      return false;
    } finally {
      setState(prev => ({
        ...prev,
        isLoading: false,
      }));
    }
  };

  // Logout function
  const logout = () => {
    authApi.logout();
    setState({
      ...initialState,
      isLoading: false,
    });
  };

  const contextValue: AuthContextType = {
    ...state,
    login,
    logout,
    refreshAuth,
  };

  return (
    <AuthContext.Provider value={contextValue}>
      {children}
    </AuthContext.Provider>
  );
};