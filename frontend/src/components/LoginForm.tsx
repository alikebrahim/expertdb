import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { useAuth } from '../hooks/useAuth';
import { useFormWithNotifications } from '../hooks/useForm';
import { loginSchema } from '../utils/formSchemas';
import { Form } from './ui/Form';
import { FormField } from './ui/FormField';
import { LoadingOverlay } from './ui/LoadingSpinner';

interface LoginFormData {
  email: string;
  password: string;
}

const LoginForm = () => {
  const [isLoading, setIsLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();
  
  const form = useFormWithNotifications<LoginFormData>({
    schema: loginSchema as any,
    defaultValues: {
      email: '',
      password: '',
    },
  });
  
  const onSubmit = async (data: LoginFormData): Promise<void> => {
    setIsLoading(true);
    
    try {
      const success = await login(data.email, data.password);
      
      if (success) {
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
        
        // Success is handled by navigation
      } else {
        // Get the error message from AuthContext if available
        const authError = localStorage.getItem('auth_error');
        if (authError) {
          localStorage.removeItem('auth_error');
          throw new Error(authError);
        } else {
          throw new Error('Invalid email or password. Please check your credentials and try again.');
        }
      }
    } catch (error) {
      console.error('Login error:', error);
      
      let errorMessage = 'An unexpected error occurred during login. Please try again.';
      if (error instanceof Error) {
        errorMessage = `Login error: ${error.message}`;
      }
      
      // Error will be handled by form's error handling
      throw error;
    } finally {
      setIsLoading(false);
    }
  };
  
  return (
    <LoadingOverlay isLoading={isLoading} className="w-full max-w-md">
      <Form
        form={form}
        onSubmit={onSubmit}
        className="space-y-4 w-full max-w-md"
        submitText="Log In"
        submitButtonPosition="center"
      >
        <FormField
          form={form}
          name="email"
          label="Email"
          type="email"
          placeholder="Enter your email"
          required
        />
        
        <FormField
          form={form}
          name="password"
          label="Password"
          type="password"
          placeholder="Enter your password"
          required
        />
      </Form>
    </LoadingOverlay>
  );
};

export default LoginForm;