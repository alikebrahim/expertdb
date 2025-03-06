/* eslint-disable react-refresh/only-export-components */
import React, { createContext, useContext, useState, useEffect } from "react";
import axios from "axios";

interface User {
  id: string;
  name: string;
  email: string;
  role: string;
}

interface AuthContextType {
  token: string | null;
  user: User | null;
  login: (token: string, user: User) => void;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const [authState, setAuthState] = useState<{
    token: string | null;
    user: User | null;
  }>({
    token: localStorage.getItem("token"),
    user: null,
  });

  // Check token validity on mount
  useEffect(() => {
    const validateToken = async () => {
      const storedToken = localStorage.getItem("token");
      if (storedToken && !authState.user) {
        try {
          const response = await axios.get("http://localhost:8080/api/auth/me", {
            headers: { Authorization: `Bearer ${storedToken}` },
          });
          setAuthState({ token: storedToken, user: response.data.user });
        } catch {
          // Invalid token, clear it
          localStorage.removeItem("token");
          setAuthState({ token: null, user: null });
        }
      }
    };
    validateToken();
  }, [authState.user]);

  const login = (token: string, user: User) => {
    localStorage.setItem("token", token);
    setAuthState({ token, user });
  };

  const logout = () => {
    localStorage.removeItem("token");
    setAuthState({ token: null, user: null });
  };

  const isAuthenticated = !!authState.token && !!authState.user;

  return (
    <AuthContext.Provider value={{ ...authState, login, logout, isAuthenticated }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used within an AuthProvider");
  return context;
};