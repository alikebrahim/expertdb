import { createContext, useEffect, useState, ReactNode } from 'react';
import { authApi } from '../services/api';
import { User, AuthState } from '../types';

interface AuthContextType extends AuthState {
  login: (email: string, password: string) => Promise<boolean>;
  logout: () => void;
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
});

export const AuthProvider = ({ children }: { children: ReactNode }) => {
  const [state, setState] = useState<AuthState>(initialState);

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
      console.log('Login attempt:', { email });
      setState({
        ...state,
        isLoading: true,
        error: null,
      });

      console.log('API URL:', import.meta.env.VITE_API_URL);
      const response = await authApi.login(email, password);
      console.log('Login response:', response);
      
      // Check if the response includes token directly (backend format)
      if (response.token && response.user) {
        console.log('Direct token format detected - login successful');
        const { token, user } = response;
        
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
      // Check if the response has data property with token (expected format)
      else if (response.success && response.data) {
        console.log('Success with data property format detected');
        const { token, user } = response.data;
        console.log('Login successful:', { user });
        
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
      // Check for success message in the response
      else if (response.message && response.message.toLowerCase().includes("success")) {
        console.log('Success message detected in response');
        // Handle case where token might be in a different structure
        const token = response.token || '';
        const user = response.user || {};
        
        if (token && user) {
          console.log('Extracted token and user from response');
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
      }
      else {
        console.log('Login failed with response:', response);
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
  };

  return (
    <AuthContext.Provider value={contextValue}>
      {children}
    </AuthContext.Provider>
  );
};