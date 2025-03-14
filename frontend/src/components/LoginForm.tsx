import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import Input from './ui/Input';
import Button from './ui/Button';

interface LoginFormData {
  email: string;
  password: string;
}

const LoginForm = () => {
  const [loginError, setLoginError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();
  
  const { 
    register, 
    handleSubmit, 
    formState: { errors } 
  } = useForm<LoginFormData>();
  
  const onSubmit = async (data: LoginFormData) => {
    setLoginError(null);
    setIsLoading(true);
    
    try {
      console.log('Attempting login with:', data.email);
      const success = await login(data.email, data.password);
      
      if (success) {
        // Check user role from auth context to determine redirect
        const userStr = localStorage.getItem('user');
        if (userStr) {
          try {
            const user = JSON.parse(userStr);
            console.log('Login successful, redirecting based on role:', user.role);
            
            if (user.role === 'admin') {
              navigate('/admin');
            } else {
              navigate('/search');
            }
          } catch (e) {
            console.error('Error parsing user data:', e);
            navigate('/search'); // default fallback
          }
        } else {
          console.warn('No user data found in localStorage after successful login');
          navigate('/search'); // default fallback
        }
      } else {
        // Get the error message from AuthContext if available
        const authError = localStorage.getItem('auth_error');
        if (authError) {
          setLoginError(authError);
          localStorage.removeItem('auth_error');
        } else {
          setLoginError('Invalid email or password. Please check your credentials and try again.');
        }
      }
    } catch (error) {
      console.error('Login error:', error);
      if (error instanceof Error) {
        setLoginError(`Login error: ${error.message}`);
      } else {
        setLoginError('An unexpected error occurred during login. Please try again.');
      }
    } finally {
      setIsLoading(false);
    }
  };
  
  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 w-full max-w-md">
      {loginError && (
        <div className="bg-secondary bg-opacity-10 text-secondary p-3 rounded">
          {loginError}
        </div>
      )}
      
      <Input
        label="Email"
        type="email"
        error={errors.email?.message}
        {...register('email', { 
          required: 'Email is required', 
          pattern: {
            value: /\S+@\S+\.\S+/,
            message: 'Invalid email address',
          },
        })}
      />
      
      <Input
        label="Password"
        type="password"
        error={errors.password?.message}
        {...register('password', { 
          required: 'Password is required',
          minLength: {
            value: 6,
            message: 'Password must be at least 6 characters',
          },
        })}
      />
      
      <Button 
        type="submit" 
        variant="primary" 
        fullWidth 
        isLoading={isLoading}
      >
        Log In
      </Button>
    </form>
  );
};

export default LoginForm;