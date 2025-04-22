import { request } from './client';
import { Expert, ExpertListResponse } from '../types';

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
  }>({
    url: '/experts',
    method: 'POST',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });

export const updateExpert = (id: string, data: FormData) => 
  request<null>({
    url: `/experts/${id}`,
    method: 'PUT',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  
export const deleteExpert = (id: string) => 
  request<null>({
    url: `/experts/${id}`,
    method: 'DELETE',
  });