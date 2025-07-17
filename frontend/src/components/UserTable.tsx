import React, { useState } from 'react';
import { User } from '../types';
import { usersApi } from '../services/api';
import { Table, Button } from './ui';

interface PaginationProps {
  currentPage: number;
  totalPages: number;
  onPageChange: (page: number) => void;
}

interface UserTableProps {
  users: User[];
  isLoading: boolean;
  error: string | null;
  onEditUser: (user: User) => void;
  onManageRoles?: (user: User) => void;
  onRefresh: () => void;
  pagination?: PaginationProps;
}

const UserTable = ({ 
  users, 
  isLoading, 
  error, 
  onEditUser,
  onManageRoles,
  onRefresh,
  pagination
}: UserTableProps) => {
  const [deletingUserId, setDeletingUserId] = useState<string | null>(null);
  const [isDeleting, setIsDeleting] = useState(false);
  
  const formatDate = (dateString: string | null | undefined) => {
    if (!dateString) return 'Never';
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };
  
  const handleDelete = async (userId: number) => {
    if (window.confirm('Are you sure you want to delete this user?')) {
      setDeletingUserId(userId.toString());
      setIsDeleting(true);
      
      try {
        const response = await usersApi.deleteUser(userId.toString());
        
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
        <Table>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Name</Table.HeaderCell>
              <Table.HeaderCell>Email</Table.HeaderCell>
              <Table.HeaderCell>Role</Table.HeaderCell>
              <Table.HeaderCell>Status</Table.HeaderCell>
              <Table.HeaderCell>Last Login</Table.HeaderCell>
              <Table.HeaderCell>Actions</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {users.map((user) => (
              <Table.Row key={user.id}>
                <Table.Cell>{user.name}</Table.Cell>
                <Table.Cell>{user.email}</Table.Cell>
                <Table.Cell>
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    user.role === 'admin' 
                      ? 'bg-purple-100 text-purple-800' 
                      : 'bg-blue-100 text-blue-800'
                  }`}>
                    {user.role}
                  </span>
                </Table.Cell>
                <Table.Cell>
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                    user.isActive 
                      ? 'bg-green-100 text-green-800' 
                      : 'bg-red-100 text-red-800'
                  }`}>
                    {user.isActive ? 'Active' : 'Inactive'}
                  </span>
                </Table.Cell>
                <Table.Cell>{formatDate(user.lastLogin)}</Table.Cell>
                <Table.Cell>
                  <div className="flex space-x-2">
                    <Button 
                      variant="outline"
                      size="sm"
                      onClick={() => onEditUser(user)}
                    >
                      Edit
                    </Button>
                    {onManageRoles && (
                      <Button 
                        variant="outline"
                        size="sm"
                        onClick={() => onManageRoles(user)}
                        className="text-blue-600 hover:bg-blue-50"
                      >
                        Roles
                      </Button>
                    )}
                    <Button 
                      variant="outline"
                      size="sm"
                      onClick={() => handleDelete(user.id)}
                      isLoading={isDeleting && deletingUserId === user.id.toString()}
                      className="text-secondary hover:bg-secondary hover:text-white"
                    >
                      Delete
                    </Button>
                  </div>
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      )}
      
      {/* Pagination controls */}
      {pagination && pagination.totalPages > 1 && !isLoading && users.length > 0 && (
        <div className="px-6 py-3 border-t border-gray-200 bg-gray-50 flex items-center justify-between">
          <div className="text-sm text-gray-700">
            Page {pagination.currentPage} of {pagination.totalPages}
          </div>
          <div className="flex space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => pagination.onPageChange(pagination.currentPage - 1)}
              disabled={pagination.currentPage <= 1}
            >
              Previous
            </Button>
            <Button
              variant="outline"
              size="sm"
              onClick={() => pagination.onPageChange(pagination.currentPage + 1)}
              disabled={pagination.currentPage >= pagination.totalPages}
            >
              Next
            </Button>
          </div>
        </div>
      )}
    </div>
  );
};

export default UserTable;