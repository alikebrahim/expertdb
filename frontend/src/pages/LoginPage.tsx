import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { LoginForm } from '../components/forms';
import NotificationContainer from '../components/ui/NotificationContainer';

const LoginPage = () => {
  const { isAuthenticated, user } = useAuth();
  const navigate = useNavigate();
  
  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated && user) {
      // All users go to search page as it's the main function of the app
      navigate('/search');
    }
  }, [isAuthenticated, user, navigate]);
  
  return (
    <div className="min-h-screen flex items-center justify-center bg-accent p-4">
      <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
        <div className="text-center mb-8">
          <img 
            src="/BQA - Horizontal Logo with Descriptor.svg" 
            alt="BQA Logo" 
            className="mx-auto h-16 mb-4"
          />
          <h1 className="text-2xl font-bold text-primary">
            Expert Database
          </h1>
          <p className="text-neutral-600">
            Log in to access the expert database management system
          </p>
        </div>
        
        <LoginForm />
      </div>
      <NotificationContainer />
    </div>
  );
};

export default LoginPage;