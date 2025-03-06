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
  }>(() => {
    const token = localStorage.getItem("token");
    const userStr = localStorage.getItem("user");
    return {
      token,
      user: userStr ? JSON.parse(userStr) as User : null,
    };
  });

  // Check token validity on mount
  useEffect(() => {
    const validateToken = async () => {
      const storedToken = localStorage.getItem("token");
      const storedUser = localStorage.getItem("user");
      if (storedToken && !storedUser) {
        try {
          const response = await axios.get("http://localhost:8080/api/auth/me", {
            headers: { Authorization: `Bearer ${storedToken}` },
          });
          const user = response.data.user;
          localStorage.setItem("user", JSON.stringify(user));
          setAuthState({ token: storedToken, user });
        } catch {
          // If backend isn't running, rely on stored data if available
          localStorage.removeItem("token");
          localStorage.removeItem("user");
          setAuthState({ token: null, user: null });
        }
      }
    };
    validateToken();
  }, []);

  const login = (token: string, user: User) => {
    localStorage.setItem("token", token);
    localStorage.setItem("user", JSON.stringify(user));
    setAuthState({ token, user });
  };

  const logout = () => {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
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