import { request } from './client';
import { ExpertRequest } from '../types';

export const getExpertRequests = (limit: number = 10, offset: number = 0, params?: Record<string, string | boolean>) => 
  request<ExpertRequest[]>({
    url: '/expert-requests',
    method: 'GET',
    params: {
      ...params,
      limit,
      offset
    },
  });

export const getExpertRequestById = (id: string) => 
  request<ExpertRequest>({
    url: `/expert-requests/${id}`,
    method: 'GET',
  });

export const createExpertRequest = (data: FormData) => 
  request<ExpertRequest>({
    url: '/expert-requests',
    method: 'POST',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });

export const updateExpertRequest = (id: string, data: FormData) => 
  request<{
    success: boolean;
    message: string;
  }>({
    url: `/expert-requests/${id}`,
    method: 'PUT',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  
export const editExpertRequest = (id: string, data: FormData) => 
  request<{
    success: boolean;
    message: string;
  }>({
    url: `/expert-requests/${id}/edit`,
    method: 'PUT',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  
export const batchApprove = (data: FormData) => 
  request<{
    success: boolean;
    message: string;
    results: Array<{
      id: number;
      status: 'success' | 'failed';
      error?: string;
    }>;
  }>({
    url: '/expert-requests/batch-approve',
    method: 'POST',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });