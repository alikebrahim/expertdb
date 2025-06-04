import { request } from './client';
import { Engagement, EngagementListResponse, ApiResponse } from '../types';

export const getEngagements = (limit: number = 10, offset: number = 0, params?: Record<string, string | boolean>) => 
  request<EngagementListResponse>({
    url: '/api/engagements',
    method: 'GET',
    params: {
      ...params,
      limit,
      offset
    },
  });

export const getExpertEngagements = (expertId: number, limit: number = 10, offset: number = 0) => 
  request<EngagementListResponse>({
    url: `/api/experts/${expertId}/engagements`,
    method: 'GET',
    params: {
      limit,
      offset
    },
  });

export const createEngagement = (engagement: Partial<Engagement>) => 
  request<{id: number}>({
    url: '/api/engagements',
    method: 'POST',
    data: engagement,
  });

export const updateEngagement = (id: number, engagement: Partial<Engagement>) => 
  request<{id: number}>({
    url: `/api/engagements/${id}`,
    method: 'PUT',
    data: engagement,
  });

export const deleteEngagement = (id: number) => 
  request<null>({
    url: `/api/engagements/${id}`,
    method: 'DELETE',
  });

export const importEngagements = (data: FormData) => 
  request<{
    success: boolean;
    successCount: number;
    failureCount: number;
    totalCount: number;
    errors: Record<string, string>;
  }>({
    url: '/api/engagements/import',
    method: 'POST',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });