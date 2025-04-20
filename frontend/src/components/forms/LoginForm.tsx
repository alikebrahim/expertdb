import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { useUI } from '../../hooks/useUI';
import Input from '../ui/Input';
import Button from '../ui/Button';

interface LoginFormData {
  email: string;
  password: string;
}

const LoginForm = () => {
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();
  const { addNotification } = useUI();
  const navigate = useNavigate();
  
  const { 
    register, 
    handleSubmit, 
    formState: { errors } 
  } = useForm<LoginFormData>();
  
  const onSubmit = async (data: LoginFormData) => {
    setIsLoading(true);
    
    try {
      const success = await login(data.email, data.password);
      
      if (success) {
        addNotification({
          type: 'success',
          message: 'Login successful!',
          duration: 3000,
        });
        
        // Check user role from auth context to determine redirect
        const userStr = localStorage.getItem('user');
        if (userStr) {
          try {
            const user = JSON.parse(userStr);
            
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
          addNotification({
            type: 'error',
            message: authError,
            duration: 5000,
          });
          localStorage.removeItem('auth_error');
        } else {
          addNotification({
            type: 'error',
            message: 'Invalid email or password. Please check your credentials and try again.',
            duration: 5000,
          });
        }
      }
    } catch (error) {
      console.error('Login error:', error);
      
      let errorMessage = 'An unexpected error occurred during login. Please try again.';
      if (error instanceof Error) {
        errorMessage = `Login error: ${error.message}`;
      }
      
      addNotification({
        type: 'error',
        message: errorMessage,
        duration: 5000,
      });
    } finally {
      setIsLoading(false);
    }
  };
  
  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4 w-full max-w-md">
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