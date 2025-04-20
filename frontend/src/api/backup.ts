import { apiClient } from './client';

export const generateBackup = () => 
  apiClient({
    url: '/backup',
    method: 'GET',
    responseType: 'blob',
  });