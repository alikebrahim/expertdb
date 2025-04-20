import { request } from './client';
import { Document } from '../types';

export const uploadDocument = (data: FormData) => 
  request<{
    id: number;
    success: boolean;
    message: string;
  }>({
    url: '/documents',
    method: 'POST',
    data,
    headers: {
      'Content-Type': 'multipart/form-data',
    },
  });

export const getDocument = (id: number) => 
  request<Document>({
    url: `/documents/${id}`,
    method: 'GET',
  });

export const deleteDocument = (id: number) => 
  request<{
    success: boolean;
    message: string;
  }>({
    url: `/documents/${id}`,
    method: 'DELETE',
  });

export const getExpertDocuments = (expertId: number) => 
  request<Document[]>({
    url: `/experts/${expertId}/documents`,
    method: 'GET',
  });