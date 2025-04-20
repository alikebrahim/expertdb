/**
 * Role-based permissions for ExpertDB
 */

// Available user roles
export type UserRole = 'admin' | 'user' | 'guest';

// Permission types
export type Permission =
  | 'view_experts'
  | 'create_expert'
  | 'edit_expert'
  | 'delete_expert'
  | 'view_requests'
  | 'create_request'
  | 'approve_request'
  | 'reject_request'
  | 'view_stats'
  | 'view_engagements'
  | 'create_engagement'
  | 'edit_engagement'
  | 'delete_engagement'
  | 'view_documents'
  | 'upload_document'
  | 'delete_document'
  | 'manage_users'
  | 'manage_areas'
  | 'configure_system'
  | 'backup_data';

// Mapping of roles to permissions
const rolePermissions: Record<UserRole, Permission[]> = {
  admin: [
    'view_experts',
    'create_expert',
    'edit_expert',
    'delete_expert',
    'view_requests',
    'create_request',
    'approve_request',
    'reject_request',
    'view_stats',
    'view_engagements',
    'create_engagement',
    'edit_engagement',
    'delete_engagement',
    'view_documents',
    'upload_document',
    'delete_document',
    'manage_users',
    'manage_areas',
    'configure_system',
    'backup_data',
  ],
  
  user: [
    'view_experts',
    'view_requests',
    'create_request',
    'view_stats',
    'view_documents',
  ],
  
  guest: [
    'view_experts',
  ],
};

/**
 * Check if a user with the given role has the specified permission
 */
export const hasPermission = (role: UserRole, permission: Permission): boolean => {
  return rolePermissions[role]?.includes(permission) || false;
};

/**
 * Get all permissions for a given role
 */
export const getPermissionsForRole = (role: UserRole): Permission[] => {
  return [...(rolePermissions[role] || [])];
};

/**
 * Check if a user with the given role can access a specific route
 */
export const canAccessRoute = (role: UserRole, route: string): boolean => {
  // Simple route-based permission check
  switch (route) {
    case '/':
      return true; // Login page accessible to all
      
    case '/search':
      return hasPermission(role, 'view_experts');
      
    case '/requests':
      return hasPermission(role, 'view_requests');
      
    case '/stats':
      return hasPermission(role, 'view_stats');
      
    case '/experts/manage':
    case '/admin':
    case '/engagements':
      return role === 'admin';
      
    default:
      if (route.startsWith('/experts/')) {
        return hasPermission(role, 'view_experts');
      }
      return false;
  }
};