import { request } from './client';
import { Document, DocumentListResponse } from '../types';

export const uploadDocument = (data: FormData) => 
  request<{
    id: number;
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
  request<null>({
    url: `/documents/${id}`,
    method: 'DELETE',
  });

export const getExpertDocuments = (expertId: number) => 
  request<DocumentListResponse>({
    url: `/experts/${expertId}/documents`,
    method: 'GET',
  });