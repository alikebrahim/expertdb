import { request } from './client';
import { Engagement } from '../types';

export const getEngagements = (limit: number = 10, offset: number = 0, params?: Record<string, string | boolean>) => 
  request<Engagement[]>({
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
    success: boolean;
    message: string;
    imported: number;
    failed: number;
  }>({
    url: '/engagements/import',
    method: 'POST',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });