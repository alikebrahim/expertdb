import { request } from './client';
import { Document, DocumentListResponse } from '../types';

export const uploadDocument = (data: FormData) => 
  request<{
    id: number;
  }>({
    url: '/api/documents',
    method: 'POST',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });

export const getDocument = (id: number) => 
  request<Document>({
    url: `/api/documents/${id}`,
    method: 'GET',
  });

export const deleteDocument = (id: number) => 
  request<null>({
    url: `/api/documents/${id}`,
    method: 'DELETE',
  });

export const getExpertDocuments = (expertId: number) => 
  request<DocumentListResponse>({
    url: `/api/experts/${expertId}/documents`,
    method: 'GET',
  });