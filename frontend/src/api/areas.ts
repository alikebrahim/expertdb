import { request } from './client';

export const getExpertAreas = () => 
  request<Array<{
    id: number;
    name: string;
  }>>({
    url: '/api/expert/areas',
    method: 'GET',
  });
  
export const createExpertArea = (data: { name: string }) => 
  request<{
    id: number;
    success: boolean;
    message: string;
  }>({
    url: '/api/expert/areas',
    method: 'POST',
    data,
  });
  
export const updateExpertArea = (id: number, data: { name: string }) => 
  request<{
    success: boolean;
    message: string;
  }>({
    url: `/api/expert/areas/${id}`,
    method: 'PUT',
    data,
  });