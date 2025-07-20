import { request } from './client';
import { ExpertRequest, RequestListResponse, BatchApproveResponse } from '../types';

export const getExpertRequests = (limit: number = 10, offset: number = 0, params?: Record<string, string | boolean>) => 
  request<RequestListResponse>({
    url: '/api/expert-requests',
    method: 'GET',
    params: {
      ...params,
      limit,
      offset
    },
  });

export const getExpertRequestById = (id: string) => 
  request<ExpertRequest>({
    url: `/api/expert-requests/${id}`,
    method: 'GET',
  });

export const createExpertRequest = (data: FormData) => 
  request<{id: number}>({
    url: '/api/expert-requests',
    method: 'POST',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });

export const updateExpertRequest = (id: string, data: FormData) => 
  request<null>({
    url: `/api/expert-requests/${id}`,
    method: 'PUT',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  
export const editExpertRequest = (id: string, data: FormData) => 
  request<null>({
    url: `/api/expert-requests/${id}/edit`,
    method: 'PUT',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });
  
export const batchApprove = (data: FormData) => 
  request<BatchApproveResponse>({
    url: '/api/expert-requests/batch-approve',
    method: 'POST',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });

export const batchReject = (data: { requestIds: number[]; comments: string }) => 
  request<BatchApproveResponse>({
    url: '/api/expert-requests/batch-reject',
    method: 'POST',
    data,
  });

export const saveDraft = (data: Partial<any>) => 
  request<{id: number}>({
    url: '/api/expert-requests/draft',
    method: 'POST',
    data,
  });

export const getDraft = (id: string) => 
  request<any>({
    url: `/api/expert-requests/draft/${id}`,
    method: 'GET',
  });

export const requestAmendment = (id: string, data: { comments: string }) => 
  request<null>({
    url: `/api/expert-requests/${id}/request-amendment`,
    method: 'POST',
    data,
  });

export const archiveRequest = (id: string) => 
  request<null>({
    url: `/api/expert-requests/${id}/archive`,
    method: 'POST',
  });