import { ApiResponse } from '../types';
import { apiClient } from './client';

export interface UserAssignments {
  user_id: number;
  planner_applications: number[];
  manager_applications: number[];
}

export interface AssignmentRequest {
  application_ids: number[];
}

export interface AssignmentResponse {
  message: string;
  user_id: number;
  assigned_applications: number;
}

export interface RemovalResponse {
  message: string;
  user_id: number;
  removed_applications: number;
}

/**
 * Assigns a user as planner to multiple applications
 */
export const assignUserToPlannerApplications = async (
  userId: number,
  applicationIds: number[]
): Promise<ApiResponse<AssignmentResponse>> => {
  const response = await apiClient.post<ApiResponse<AssignmentResponse>>(
    `/api/users/${userId}/planner-assignments`,
    { application_ids: applicationIds }
  );
  return response.data;
};

/**
 * Assigns a user as manager to multiple applications
 */
export const assignUserToManagerApplications = async (
  userId: number,
  applicationIds: number[]
): Promise<ApiResponse<AssignmentResponse>> => {
  const response = await apiClient.post<ApiResponse<AssignmentResponse>>(
    `/api/users/${userId}/manager-assignments`,
    { application_ids: applicationIds }
  );
  return response.data;
};

/**
 * Removes planner assignments for a user from specific applications
 */
export const removeUserPlannerAssignments = async (
  userId: number,
  applicationIds: number[]
): Promise<ApiResponse<RemovalResponse>> => {
  const response = await apiClient.delete<ApiResponse<RemovalResponse>>(
    `/api/users/${userId}/planner-assignments`,
    { data: { application_ids: applicationIds } }
  );
  return response.data;
};

/**
 * Removes manager assignments for a user from specific applications
 */
export const removeUserManagerAssignments = async (
  userId: number,
  applicationIds: number[]
): Promise<ApiResponse<RemovalResponse>> => {
  const response = await apiClient.delete<ApiResponse<RemovalResponse>>(
    `/api/users/${userId}/manager-assignments`,
    { data: { application_ids: applicationIds } }
  );
  return response.data;
};

/**
 * Gets all planner and manager assignments for a user
 */
export const getUserAssignments = async (
  userId: number
): Promise<ApiResponse<UserAssignments>> => {
  const response = await apiClient.get<ApiResponse<UserAssignments>>(
    `/api/users/${userId}/assignments`
  );
  return response.data;
};