import { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import LoginForm from '../components/LoginForm';

const LoginPage = () => {
  const { isAuthenticated, user } = useAuth();
  const navigate = useNavigate();
  
  // Redirect if already authenticated
  useEffect(() => {
    if (isAuthenticated && user) {
      if (user.role === 'admin') {
        navigate('/admin');
      } else {
        navigate('/search');
      }
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
    </div>
  );
};

export default LoginPage;