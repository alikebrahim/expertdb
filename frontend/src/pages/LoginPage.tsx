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
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-br from-secondary to-secondary-dark p-4">
      <div className="bg-white p-8 rounded-xl shadow-2xl w-full max-w-md border border-neutral-200">
        <div className="text-center mb-8">
          <img 
            src="/BQA - Horizontal Logo with Descriptor.svg" 
            alt="BQA Logo" 
            className="mx-auto h-16 mb-6"
          />
          <h1 className="text-2xl font-bold text-secondary mb-2">
            Expert Database
          </h1>
          <p className="text-neutral-600 text-sm">
            Log in to access the expert database management system
          </p>
        </div>
        
        <LoginForm />
        
        <div className="mt-6 text-center">
          <p className="text-xs text-neutral-500">
            Powered by Bahrain Qualifications Authority
          </p>
        </div>
      </div>
      <NotificationContainer />
    </div>
  );
};

export default LoginPage;