import { useState, useEffect, useCallback } from 'react';
import { User, ExpertRequest } from '../types';
import { usersApi, expertRequestsApi } from '../services/api';
import UserTable from '../components/UserTable';
import UserForm from '../components/UserForm';
import Button from '../components/ui/Button';
import ExpertRequestTable from '../components/ExpertRequestTable';
import { RoleAssignmentForm } from '../components/forms/RoleAssignmentForm';
import Layout from '../components/layout/Layout';

type ActiveTab = 'users' | 'requests';

const AdminPage = () => {
  const [activeTab, setActiveTab] = useState<ActiveTab>('users');
  
  // Users state
  const [users, setUsers] = useState<User[]>([]);
  const [isLoadingUsers, setIsLoadingUsers] = useState(true);
  const [userError, setUserError] = useState<string | null>(null);
  const [showUserForm, setShowUserForm] = useState(false);
  const [editingUser, setEditingUser] = useState<User | undefined>(undefined);
  const [userPage, setUserPage] = useState(1);
  const [userTotalPages, setUserTotalPages] = useState(1);
  const [userLimit] = useState(10);
  const [showRoleAssignmentForm, setShowRoleAssignmentForm] = useState(false);
  const [userForRoleAssignment, setUserForRoleAssignment] = useState<User | null>(null);
  
  // Expert requests state
  const [requests, setRequests] = useState<ExpertRequest[]>([]);
  const [isLoadingRequests, setIsLoadingRequests] = useState(true);
  const [requestError, setRequestError] = useState<string | null>(null);
  const [requestPage, setRequestPage] = useState(1);
  const [requestTotalPages, setRequestTotalPages] = useState(1);
  const [requestLimit] = useState(10);
  
  // Fetch users on mount and when needed
  const fetchUsers = useCallback(async () => {
    setIsLoadingUsers(true);
    setUserError(null);
    
    try {
      const response = await usersApi.getUsers(userLimit, (userPage - 1) * userLimit);
      
      if (response.success && response.data) {
        setUsers(response.data.data || []);
        setUserTotalPages(response.data.totalPages || 1);
      } else {
        // Check if this is likely a "no users" situation or a real error
        if (response.message?.includes("not found") || 
            response.message?.toLowerCase().includes("empty") ||
            response.message?.toLowerCase().includes("no users")) {
          // This is likely just an empty database, not an error
          setUsers([]);
          setUserTotalPages(1);
        } else {
          setUserError(response.message || 'Failed to fetch users');
        }
      }
    } catch (error) {
      console.error('Error fetching users:', error);
      setUserError('An error occurred while fetching users');
    } finally {
      setIsLoadingUsers(false);
    }
  }, [userPage, userLimit]);
  
  const handleUserPageChange = (page: number) => {
    setUserPage(page);
  };
  
  // Fetch expert requests on mount and when needed
  const fetchRequests = useCallback(async () => {
    setIsLoadingRequests(true);
    setRequestError(null);
    
    try {
      const response = await expertRequestsApi.getExpertRequests(requestLimit, (requestPage - 1) * requestLimit);
      
      if (response.success && response.data) {
        // Handle case where response.data is an array directly
        const requests = Array.isArray(response.data) ? response.data : [];
        setRequests(requests);
        setRequestTotalPages(1); // Since it's a simple array, no pagination
      } else {
        // Check if this is likely a "no requests" situation or a real error
        if (response.message?.includes("not found") || 
            response.message?.toLowerCase().includes("empty") ||
            response.message?.toLowerCase().includes("no requests")) {
          // This is likely just an empty database, not an error
          setRequests([]);
          setRequestTotalPages(1);
        } else {
          setRequestError(response.message || 'Failed to fetch expert requests');
        }
      }
    } catch (error) {
      console.error('Error fetching expert requests:', error);
      setRequestError('An error occurred while fetching requests');
    } finally {
      setIsLoadingRequests(false);
    }
  }, [requestPage, requestLimit]);
  
  const handleRequestPageChange = (page: number) => {
    setRequestPage(page);
  };
  
  // Initial fetch
  useEffect(() => {
    if (activeTab === 'users') {
      fetchUsers();
    } else {
      fetchRequests();
    }
  }, [activeTab, fetchUsers, fetchRequests]);
  
  // Fetch data when pagination changes
  useEffect(() => {
    if (activeTab === 'users') {
      fetchUsers();
    }
  }, [userPage, activeTab, fetchUsers]);
  
  useEffect(() => {
    if (activeTab === 'requests') {
      fetchRequests();
    }
  }, [requestPage, activeTab, fetchRequests]);
  
  // Tab switching
  const handleTabClick = (tab: ActiveTab) => {
    setActiveTab(tab);
  };
  
  // User form handlers
  const handleNewUser = () => {
    setEditingUser(undefined);
    setShowUserForm(true);
  };
  
  const handleEditUser = (user: User) => {
    setEditingUser(user);
    setShowUserForm(true);
  };
  
  const handleUserFormSuccess = () => {
    setShowUserForm(false);
    fetchUsers();
  };
  
  const handleUserFormCancel = () => {
    setShowUserForm(false);
  };

  // Role assignment handlers
  const handleManageRoles = (user: User) => {
    setUserForRoleAssignment(user);
    setShowRoleAssignmentForm(true);
  };

  const handleRoleAssignmentSuccess = () => {
    setShowRoleAssignmentForm(false);
    setUserForRoleAssignment(null);
    // No need to refresh users list as assignments don't change user data
  };

  const handleRoleAssignmentCancel = () => {
    setShowRoleAssignmentForm(false);
    setUserForRoleAssignment(null);
  };
  
  // Request update function no longer needed with ExpertRequestTable component
  
  // Request actions no longer needed with ExpertRequestTable component
  
  return (
    <Layout>
      <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-primary">Admin Panel</h1>
        <p className="text-neutral-600">
          Manage users and expert requests
        </p>
      </div>
      
      {/* Tabs */}
      <div className="mb-6 border-b border-neutral-200">
        <nav className="-mb-px flex">
          <button
            onClick={() => handleTabClick('users')}
            className={`mr-8 py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'users'
                ? 'border-primary text-primary'
                : 'border-transparent text-neutral-500 hover:text-neutral-700 hover:border-neutral-300'
            }`}
          >
            User Management
          </button>
          <button
            onClick={() => handleTabClick('requests')}
            className={`mr-8 py-4 px-1 border-b-2 font-medium text-sm ${
              activeTab === 'requests'
                ? 'border-primary text-primary'
                : 'border-transparent text-neutral-500 hover:text-neutral-700 hover:border-neutral-300'
            }`}
          >
            Expert Requests
          </button>
        </nav>
      </div>
      
      {/* User Management Tab */}
      {activeTab === 'users' && (
        <div>
          {showUserForm ? (
            <div className="bg-white rounded-md shadow p-6 mb-6">
              <h2 className="text-xl font-semibold text-primary mb-4">
                {editingUser ? 'Edit User' : 'Create New User'}
              </h2>
              <UserForm 
                user={editingUser} 
                onSuccess={handleUserFormSuccess}
                onCancel={handleUserFormCancel}
              />
            </div>
          ) : (
            <div className="mb-4 flex justify-end">
              <Button onClick={handleNewUser}>
                Create New User
              </Button>
            </div>
          )}
          
          <div className="bg-white rounded-md shadow p-6">
            <h2 className="text-xl font-semibold text-primary mb-4">
              Users
            </h2>
            <UserTable 
              users={users}
              isLoading={isLoadingUsers}
              error={userError}
              onEditUser={handleEditUser}
              onManageRoles={handleManageRoles}
              onRefresh={fetchUsers}
              pagination={{
                currentPage: userPage,
                totalPages: userTotalPages,
                onPageChange: handleUserPageChange
              }}
            />
          </div>
        </div>
      )}
      
      {/* Role Assignment Modal */}
      {showRoleAssignmentForm && userForRoleAssignment && (
        <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
          <div className="relative top-20 mx-auto p-5 border w-11/12 max-w-4xl shadow-lg rounded-md bg-white">
            <RoleAssignmentForm
              user={userForRoleAssignment}
              onSuccess={handleRoleAssignmentSuccess}
              onCancel={handleRoleAssignmentCancel}
            />
          </div>
        </div>
      )}
      
      {/* Expert Requests Tab */}
      {activeTab === 'requests' && (
        <div className="bg-white rounded-md shadow p-6">
          <h2 className="text-xl font-semibold text-primary mb-4">
            Expert Requests
          </h2>
          
          <ExpertRequestTable
            requests={requests}
            isLoading={isLoadingRequests}
            error={requestError}
            pagination={{
              currentPage: requestPage,
              totalPages: requestTotalPages,
              onPageChange: handleRequestPageChange
            }}
          />
        </div>
      )}
      </div>
    </Layout>
  );
};

export default AdminPage;