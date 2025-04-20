import { request } from './client';
import { User, PaginatedResponse } from '../types';

export const getUsers = (page: number = 1, limit: number = 10, params?: Record<string, string | boolean>) => 
  request<PaginatedResponse<User>>({
    url: '/users',
    method: 'GET',
    params: {
      ...params,
      page,
      limit
    },
  });

export const getUserById = (id: string) => 
  request<User>({
    url: `/users/${id}`,
    method: 'GET',
  });

export const createUser = (data: Partial<User>) => 
  request<User>({
    url: '/users',
    method: 'POST',
    data,
  });

export const updateUser = (id: string, data: Partial<User>) => 
  request<User>({
    url: `/users/${id}`,
    method: 'PUT',
    data,
  });

export const deleteUser = (id: string) => 
  request<void>({
    url: `/users/${id}`,
    method: 'DELETE',
  });