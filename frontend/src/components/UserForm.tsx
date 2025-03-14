import { useState } from 'react';
import { useForm } from 'react-hook-form';
import { User } from '../types';
import { usersApi } from '../services/api';
import Input from './ui/Input';
import Button from './ui/Button';

interface UserFormData {
  name: string;
  email: string;
  password?: string;
  role: 'admin' | 'user';
  isActive: boolean;
}

interface UserFormProps {
  user?: User;
  onSuccess: () => void;
  onCancel: () => void;
}

const UserForm = ({ user, onSuccess, onCancel }: UserFormProps) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  const { 
    register, 
    handleSubmit, 
    formState: { errors }
  } = useForm<UserFormData>({
    defaultValues: user ? {
      name: user.name,
      email: user.email,
      role: user.role as 'admin' | 'user',
      isActive: user.isActive
    } : {
      role: 'user',
      isActive: true
    }
  });
  
  const isEditMode = !!user;
  
  const onSubmit = async (data: UserFormData) => {
    setIsSubmitting(true);
    setError(null);
    
    try {
      // Remove password if empty (for edit mode)
      if (data.password === '') {
        delete data.password;
      }
      
      let response;
      
      if (isEditMode && user) {
        response = await usersApi.updateUser(user.id, data);
      } else {
        if (!data.password) {
          setError('Password is required for new users');
          setIsSubmitting(false);
          return;
        }
        response = await usersApi.createUser(data);
      }
      
      if (response.success) {
        onSuccess();
      } else {
        setError(response.message || `Failed to ${isEditMode ? 'update' : 'create'} user`);
      }
    } catch (error) {
      console.error(`Error ${isEditMode ? 'updating' : 'creating'} user:`, error);
      setError(`An error occurred while ${isEditMode ? 'updating' : 'creating'} the user`);
    } finally {
      setIsSubmitting(false);
    }
  };
  
  return (
    <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
      {error && (
        <div className="bg-secondary bg-opacity-10 text-secondary p-3 rounded">
          {error}
        </div>
      )}
      
      <Input
        label="Name *"
        error={errors.name?.message}
        {...register('name', { required: 'Name is required' })}
      />
      
      <Input
        label="Email *"
        type="email"
        error={errors.email?.message}
        {...register('email', { 
          required: 'Email is required',
          pattern: {
            value: /\S+@\S+\.\S+/,
            message: 'Invalid email format',
          }
        })}
      />
      
      <Input
        label={isEditMode ? 'Password (leave blank to keep current)' : 'Password *'}
        type="password"
        error={errors.password?.message}
        {...register('password', { 
          ...(isEditMode ? {} : { required: 'Password is required' }),
          minLength: {
            value: 6,
            message: 'Password must be at least 6 characters',
          }
        })}
      />
      
      <div className="mb-4">
        <label className="block text-sm font-medium text-neutral-700 mb-1">
          Role *
        </label>
        <select
          className="w-full px-3 py-2 bg-white border border-neutral-300 rounded focus:outline-none focus:ring-1 focus:ring-primary focus:border-primary"
          {...register('role', { required: 'Role is required' })}
        >
          <option value="user">User</option>
          <option value="admin">Admin</option>
        </select>
        {errors.role && (
          <p className="mt-1 text-sm text-secondary">{errors.role.message}</p>
        )}
      </div>
      
      <div className="flex items-center mb-4">
        <input
          type="checkbox"
          id="isActive"
          className="h-4 w-4 text-primary focus:ring-primary border-neutral-300 rounded"
          {...register('isActive')}
        />
        <label htmlFor="isActive" className="ml-2 block text-sm text-neutral-700">
          Active Account
        </label>
      </div>
      
      <div className="flex justify-end space-x-3">
        <Button 
          type="button" 
          variant="outline" 
          onClick={onCancel}
        >
          Cancel
        </Button>
        <Button 
          type="submit" 
          isLoading={isSubmitting}
        >
          {isEditMode ? 'Update User' : 'Create User'}
        </Button>
      </div>
    </form>
  );
};

export default UserForm;