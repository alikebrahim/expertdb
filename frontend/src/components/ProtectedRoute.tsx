import { ReactNode } from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';

interface ProtectedRouteProps {
  children: ReactNode;
  requiredRole?: 'super_user' | 'admin' | 'planner' | 'regular';
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
  if (requiredRole) {
    const roleHierarchy = {
      'super_user': 4,
      'admin': 3,
      'planner': 2,
      'regular': 1
    };
    
    const userRoleLevel = roleHierarchy[user.role as keyof typeof roleHierarchy] || 0;
    const requiredRoleLevel = roleHierarchy[requiredRole] || 0;
    
    // If the user's role level is less than the required role level
    if (userRoleLevel < requiredRoleLevel) {
      // Redirect to appropriate page based on role
      if (userRoleLevel >= 3) { // admin or super_user
        return <Navigate to="/admin" replace />;
      } else if (userRoleLevel === 2) { // planner
        return <Navigate to="/search" replace />;
      } else {
        return <Navigate to="/search" replace />;
      }
    }
  }

  // Render the protected route - layout is handled by individual pages
  return <>{children}</>;
};

export default ProtectedRoute;