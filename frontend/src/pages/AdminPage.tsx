import { useState, useEffect } from 'react';
import { User, ExpertRequest } from '../types';
import { usersApi, expertRequestsApi } from '../services/api';
import UserTable from '../components/UserTable';
import UserForm from '../components/UserForm';
import Button from '../components/ui/Button';

type ActiveTab = 'users' | 'requests';

const AdminPage = () => {
  const [activeTab, setActiveTab] = useState<ActiveTab>('users');
  
  // Users state
  const [users, setUsers] = useState<User[]>([]);
  const [isLoadingUsers, setIsLoadingUsers] = useState(true);
  const [userError, setUserError] = useState<string | null>(null);
  const [showUserForm, setShowUserForm] = useState(false);
  const [editingUser, setEditingUser] = useState<User | undefined>(undefined);
  
  // Expert requests state
  const [requests, setRequests] = useState<ExpertRequest[]>([]);
  const [isLoadingRequests, setIsLoadingRequests] = useState(true);
  const [requestError, setRequestError] = useState<string | null>(null);
  
  // Fetch users on mount and when needed
  const fetchUsers = async () => {
    setIsLoadingUsers(true);
    setUserError(null);
    
    try {
      const response = await usersApi.getUsers();
      
      if (response.success) {
        setUsers(response.data);
      } else {
        // Check if this is likely a "no users" situation or a real error
        if (response.message?.includes("not found") || 
            response.message?.toLowerCase().includes("empty") ||
            response.message?.toLowerCase().includes("no users")) {
          // This is likely just an empty database, not an error
          setUsers([]);
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
  };
  
  // Fetch expert requests on mount and when needed
  const fetchRequests = async () => {
    setIsLoadingRequests(true);
    setRequestError(null);
    
    try {
      const response = await expertRequestsApi.getExpertRequests();
      
      if (response.success) {
        setRequests(response.data);
      } else {
        // Check if this is likely a "no requests" situation or a real error
        if (response.message?.includes("not found") || 
            response.message?.toLowerCase().includes("empty") ||
            response.message?.toLowerCase().includes("no requests")) {
          // This is likely just an empty database, not an error
          setRequests([]);
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
  };
  
  // Initial fetch
  useEffect(() => {
    if (activeTab === 'users') {
      fetchUsers();
    } else {
      fetchRequests();
    }
  }, [activeTab]);
  
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
  
  // Handle approving or rejecting an expert request
  const handleUpdateRequest = async (requestId: string, status: 'approved' | 'rejected', rejectionReason?: string) => {
    try {
      const response = await expertRequestsApi.updateExpertRequest(requestId, {
        status,
        rejectionReason
      });
      
      if (response.success) {
        fetchRequests();
      } else {
        alert(`Failed to update request: ${response.message}`);
      }
    } catch (error) {
      console.error('Error updating request:', error);
      alert('An error occurred while updating the request');
    }
  };
  
  // Handle request actions
  const handleApproveRequest = (request: ExpertRequest) => {
    if (window.confirm(`Are you sure you want to approve the expert request from ${request.name}?`)) {
      handleUpdateRequest(request.id, 'approved');
    }
  };
  
  const handleRejectRequest = (request: ExpertRequest) => {
    const reason = window.prompt('Please provide a reason for rejection:');
    if (reason !== null) {
      handleUpdateRequest(request.id, 'rejected', reason);
    }
  };
  
  return (
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
              onRefresh={fetchUsers}
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
          
          <div className="space-y-4">
            {isLoadingRequests ? (
              <div className="flex justify-center py-8">
                <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full"></div>
                <span className="sr-only">Loading...</span>
              </div>
            ) : requestError ? (
              // Check if the error indicates no data vs an actual error
              requestError.toLowerCase().includes("not found") || 
              requestError.toLowerCase().includes("empty") || 
              requestError.toLowerCase().includes("no requests") ? (
                <div className="bg-accent p-6 rounded text-center">
                  <p className="text-neutral-600">No expert requests have been submitted yet.</p>
                  <p className="text-sm text-neutral-500 mt-2">
                    Users can submit expert requests from the Requests page.
                  </p>
                </div>
              ) : (
                <div className="bg-secondary bg-opacity-10 text-secondary p-4 rounded">
                  <p>Error loading requests: {requestError}</p>
                  <p className="text-sm mt-2">Please try refreshing the page or contact support if the problem persists.</p>
                </div>
              )
            ) : requests.length === 0 ? (
              <div className="bg-accent p-6 rounded text-center">
                <p className="text-neutral-600">No expert requests found.</p>
                <p className="text-sm text-neutral-500 mt-2">
                  When users submit expert requests, they will appear here for your review.
                </p>
              </div>
            ) : (
              <table className="min-w-full bg-white border border-neutral-200 rounded-md overflow-hidden">
                <thead className="bg-primary text-white">
                  <tr>
                    <th className="py-3 px-4 text-left font-medium text-sm">Name</th>
                    <th className="py-3 px-4 text-left font-medium text-sm">Affiliation</th>
                    <th className="py-3 px-4 text-left font-medium text-sm">Role</th>
                    <th className="py-3 px-4 text-left font-medium text-sm">Status</th>
                    <th className="py-3 px-4 text-left font-medium text-sm">Submitted By</th>
                    <th className="py-3 px-4 text-left font-medium text-sm">Date</th>
                    <th className="py-3 px-4 text-left font-medium text-sm">Actions</th>
                  </tr>
                </thead>
                <tbody className="divide-y divide-neutral-200">
                  {requests.map((request) => (
                    <tr key={request.id}>
                      <td className="py-3 px-4 text-sm text-neutral-800">{request.name}</td>
                      <td className="py-3 px-4 text-sm text-neutral-800">{request.affiliation}</td>
                      <td className="py-3 px-4 text-sm text-neutral-800">{request.role}</td>
                      <td className="py-3 px-4 text-sm text-neutral-800">
                        <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${
                          request.status === 'approved' 
                            ? 'bg-green-100 text-green-800' 
                            : request.status === 'rejected'
                              ? 'bg-red-100 text-red-800'
                              : 'bg-yellow-100 text-yellow-800'
                        }`}>
                          {request.status.charAt(0).toUpperCase() + request.status.slice(1)}
                        </span>
                      </td>
                      <td className="py-3 px-4 text-sm text-neutral-800">{request.userId}</td>
                      <td className="py-3 px-4 text-sm text-neutral-800">
                        {new Date(request.createdAt).toLocaleDateString()}
                      </td>
                      <td className="py-3 px-4 text-sm text-neutral-800">
                        {request.status === 'pending' && (
                          <div className="flex space-x-2">
                            <Button 
                              variant="outline"
                              size="sm"
                              className="bg-green-50 text-green-700 border-green-300 hover:bg-green-100"
                              onClick={() => handleApproveRequest(request)}
                            >
                              Approve
                            </Button>
                            <Button 
                              variant="outline"
                              size="sm"
                              className="bg-red-50 text-red-700 border-red-300 hover:bg-red-100"
                              onClick={() => handleRejectRequest(request)}
                            >
                              Reject
                            </Button>
                          </div>
                        )}
                        {request.status === 'rejected' && request.rejectionReason && (
                          <span className="text-xs text-neutral-500">
                            Reason: {request.rejectionReason}
                          </span>
                        )}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            )}
          </div>
        </div>
      )}
    </div>
  );
};

export default AdminPage;