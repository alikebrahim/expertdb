/* eslint-disable react-refresh/only-export-components */
import React, { createContext, useContext, useState } from "react";

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

  const login = (token: string, user: User) => {
    localStorage.setItem("token", token);
    setAuthState({ token, user });
  };

  const logout = () => {
    localStorage.removeItem("token");
    setAuthState({ token: null, user: null });
  };

  return (
    <AuthContext.Provider value={{ ...authState, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) throw new Error("useAuth must be used within an AuthProvider");
  return context;
};