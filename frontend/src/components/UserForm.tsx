import { useState } from 'react';
import { User } from '../types';
import { usersApi } from '../services/api';
import { useFormWithNotifications } from '../hooks/useForm';
import { userSchema } from '../utils/formSchemas';
import { Form, FormField, LoadingOverlay } from './ui';

interface UserFormData {
  name: string;
  email: string;
  password?: string;
  confirmPassword?: string;
  role: 'admin' | 'user' | 'manager';
  isActive: boolean;
}

interface UserFormProps {
  user?: User;
  onSuccess: () => void;
  onCancel: () => void;
}

const UserForm = ({ user, onSuccess, onCancel }: UserFormProps) => {
  const [isSubmitting, setIsSubmitting] = useState(false);
  
  const isEditMode = !!user;
  
  const form = useFormWithNotifications<UserFormData>({
    schema: userSchema,
    defaultValues: user ? {
      name: user.name,
      email: user.email,
      role: user.role as 'admin' | 'user' | 'manager',
      isActive: user.isActive,
      password: '',
      confirmPassword: '',
    } : {
      name: '',
      email: '',
      role: 'user',
      isActive: true,
      password: '',
      confirmPassword: '',
    }
  });
  
  const roleOptions = [
    { label: 'User', value: 'user' },
    { label: 'Admin', value: 'admin' },
    { label: 'Manager', value: 'manager' },
  ];
  
  const onSubmit = async (data: UserFormData) => {
    setIsSubmitting(true);
    
    try {
      // Remove password if empty (for edit mode)
      if (data.password === '') {
        delete data.password;
      }
      
      // Remove confirmPassword before sending to API
      const { confirmPassword, ...apiData } = data;
      
      let response;
      
      if (isEditMode && user) {
        response = await usersApi.updateUser(user.id, apiData);
      } else {
        if (!data.password) {
          return { 
            success: false, 
            message: 'Password is required for new users' 
          };
        }
        response = await usersApi.createUser(apiData);
      }
      
      if (response.success) {
        onSuccess();
        return { 
          success: true, 
          message: `User ${isEditMode ? 'updated' : 'created'} successfully` 
        };
      } else {
        return { 
          success: false, 
          message: response.message || `Failed to ${isEditMode ? 'update' : 'create'} user` 
        };
      }
    } catch (error) {
      console.error(`Error ${isEditMode ? 'updating' : 'creating'} user:`, error);
      return { 
        success: false, 
        message: `An error occurred while ${isEditMode ? 'updating' : 'creating'} the user` 
      };
    } finally {
      setIsSubmitting(false);
    }
  };
  
  return (
    <LoadingOverlay isLoading={isSubmitting}>
      <Form
        form={form}
        onSubmit={form.handleSubmitWithNotifications(onSubmit)}
        className="space-y-4"
        resetOnSuccess={false}
        submitText={isEditMode ? 'Update User' : 'Create User'}
        showResetButton={true}
        resetText="Cancel"
        submitButtonPosition="right"
        onReset={onCancel}
      >
        <h2 className="text-xl font-bold mb-4">
          {isEditMode ? 'Edit User' : 'Create New User'}
        </h2>
        
        <FormField
          form={form}
          name="name"
          label="Name"
          placeholder="Enter user's full name"
          required
        />
        
        <FormField
          form={form}
          name="email"
          label="Email"
          type="email"
          placeholder="Enter user's email address"
          required
        />
        
        <FormField
          form={form}
          name="password"
          label={isEditMode ? 'Password (leave blank to keep current)' : 'Password'}
          type="password"
          placeholder="Enter password"
          required={!isEditMode}
        />
        
        <FormField
          form={form}
          name="confirmPassword"
          label="Confirm Password"
          type="password"
          placeholder="Confirm password"
          required={!isEditMode}
          hint={isEditMode ? "Only required when changing password" : ""}
        />
        
        <FormField
          form={form}
          name="role"
          label="Role"
          type="select"
          options={roleOptions}
          required
        />
        
        <FormField
          form={form}
          name="isActive"
          label="Active Account"
          type="checkbox"
          hint="Inactive users cannot log in"
        />
      </Form>
    </LoadingOverlay>
  );
};

export default UserForm;