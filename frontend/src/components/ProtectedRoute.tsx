import { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import Header from './layout/Header';
import Sidebar from './layout/Sidebar';
import Footer from './layout/Footer';

interface ProtectedRouteProps {
  children: ReactNode;
  requiredRole?: 'admin' | 'user';
}

const ProtectedRoute = ({ children, requiredRole }: ProtectedRouteProps) => {
  const { isAuthenticated, user, isLoading } = useAuth();
  const location = useLocation();

  // Show loading state while auth is being checked
  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="animate-spin h-10 w-10 border-4 border-primary border-t-transparent rounded-full"></div>
      </div>
    );
  }

  // Check if user is authenticated
  if (!isAuthenticated || !user) {
    // Redirect to login page with the return url
    return <Navigate to="/" state={{ from: location }} replace />;
  }

  // Check if user has required role (if specified)
  if (requiredRole && user.role !== requiredRole) {
    // Redirect to appropriate page based on role
    if (user.role === 'admin') {
      return <Navigate to="/admin" replace />;
    } else {
      return <Navigate to="/search" replace />;
    }
  }

  // Render the protected route with layout
  return (
    <div className="flex flex-col min-h-screen">
      <Header />
      <div className="flex flex-1">
        <Sidebar />
        <main className="flex-1 p-6 bg-neutral-50">
          {children}
        </main>
      </div>
      <Footer />
    </div>
  );
};

export default ProtectedRoute;