import { request } from './client';
import { EngagementListResponse } from '../types';

export const getEngagements = (limit: number = 10, offset: number = 0, params?: Record<string, string | boolean>) => 
  request<EngagementListResponse>({
    url: '/expert-engagements',
    method: 'GET',
    params: {
      ...params,
      limit,
      offset
    },
  });

export const importEngagements = (data: FormData) => 
  request<{
    imported: number;
    failed: number;
    details: Array<{
      expertId: number;
      status: 'success' | 'failed';
      error?: string;
    }>;
  }>({
    url: '/engagements/import',
    method: 'POST',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });