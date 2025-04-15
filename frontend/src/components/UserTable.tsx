import { useState } from 'react';
import { User } from '../types';
import { usersApi } from '../services/api';
import { Table, TableRow, TableCell } from './ui/Table';
import Button from './ui/Button';

interface UserTableProps {
  users: User[];
  isLoading: boolean;
  error: string | null;
  onEditUser: (user: User) => void;
  onRefresh: () => void;
  pagination?: {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
  };
}

const UserTable = ({ 
  users, 
  isLoading, 
  error, 
  onEditUser,
  onRefresh,
  pagination
}: UserTableProps) => {
  const [deletingUserId, setDeletingUserId] = useState<string | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  
  const headers = [
    'Name',
    'Email',
    'Role',
    'Status',
    'Last Login',
    'Actions'
  ];
  
  const formatDate = (dateString: string | null | undefined) => {
    if (!dateString) return 'Never';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };
  
  const handleDelete = async (userId: string) => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      setDeletingUserId(userId);
      setIsDeleting(true);
      
      try {
        const response = await usersApi.deleteUser(userId);
        
        if (response.success) {
          onRefresh();
        } else {
          alert(`Failed to delete user: ${response.message}`);
        }
      } catch (error) {
        console.error('Error deleting user:', error);
        alert('An error occurred while deleting user');
      } finally {
        setDeletingUserId(null);
        setIsDeleting(false);
      }
    }
  };
  
  if (error) {
    // Check if the error message indicates no data vs an actual error
    if (error.toLowerCase().includes("not found") || 
        error.toLowerCase().includes("empty") || 
        error.toLowerCase().includes("no users")) {
      return (
        <div className="bg-accent p-6 rounded text-center">
          <p className="text-neutral-600">No users have been added to the system yet.</p>
          <p className="text-sm text-neutral-500 mt-2">Use the "Create New User" button to add users.</p>
        </div>
      );
    } else {
      return (
        <div className="bg-secondary bg-opacity-10 text-secondary p-4 rounded">
          <p>Error loading users: {error}</p>
          <p className="text-sm mt-2">Please try refreshing the page or contact support if the problem persists.</p>
        </div>
      );
    }
  }
  
  return (
    <div className="space-y-4">
      {isLoading ? (
        <div className="flex justify-center py-8">
          <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full"></div>
          <span className="sr-only">Loading...</span>
        </div>
      ) : users.length === 0 ? (
        <div className="bg-accent p-6 rounded text-center">
          <p className="text-neutral-600">No users found.</p>
        </div>
      ) : (
        <Table headers={headers} pagination={pagination}>
          {users.map((user) => (
            <TableRow key={user.id}>
              <TableCell>{user.name}</TableCell>
              <TableCell>{user.email}</TableCell>
              <TableCell>
                <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                  user.role === 'admin' 
                    ? 'bg-purple-100 text-purple-800' 
                    : 'bg-blue-100 text-blue-800'
                }`}>
                  {user.role}
                </span>
              </TableCell>
              <TableCell>
                <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                  user.isActive 
                    ? 'bg-green-100 text-green-800' 
                    : 'bg-red-100 text-red-800'
                }`}>
                  {user.isActive ? 'Active' : 'Inactive'}
                </span>
              </TableCell>
              <TableCell>{formatDate(user.lastLogin)}</TableCell>
              <TableCell>
                <div className="flex space-x-2">
                  <Button 
                    variant="outline"
                    size="sm"
                    onClick={() => onEditUser(user)}
                  >
                    Edit
                  </Button>
                  <Button 
                    variant="outline"
                    size="sm"
                    onClick={() => handleDelete(user.id)}
                    isLoading={isDeleting && deletingUserId === user.id}
                    className="text-secondary hover:bg-secondary hover:text-white"
                  >
                    Delete
                  </Button>
                </div>
              </TableCell>
            </TableRow>
          ))}
        </Table>
      )}
    </div>
  );
};

export default UserTable;