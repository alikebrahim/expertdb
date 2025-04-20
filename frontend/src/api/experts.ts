import { request } from './client';
import { Expert } from '../types';

interface ExpertListResponse {
  experts: Expert[];
  pagination: {
    totalCount: number;
    totalPages: number;
    currentPage: number;
    pageSize: number;
    hasNextPage: boolean;
    hasPrevPage: boolean;
    hasMore: boolean;
  }
}

export const getExperts = (limit: number = 10, offset: number = 0, params?: Record<string, string | boolean | number>) => 
  request<ExpertListResponse>({
    url: '/experts',
    method: 'GET',
    params: {
      ...params,
      limit,
      offset
    },
  });

export const getExpertById = (id: string) => 
  request<Expert>({
    url: `/experts/${id}`,
    method: 'GET',
  });

export const createExpert = (data: FormData) => 
  request<{
    id: number;
    success: boolean;
    message: string;
  }>({
    url: '/experts',
    method: 'POST',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });

export const updateExpert = (id: string, data: FormData) => 
  request<{
    success: boolean;
    message: string;
  }>({
    url: `/experts/${id}`,
    method: 'PUT',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  
export const deleteExpert = (id: string) => 
  request<{
    success: boolean;
    message: string;
  }>({
    url: `/experts/${id}`,
    method: 'DELETE',
  });