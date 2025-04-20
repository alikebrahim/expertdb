import { request } from './client';
import { Phase } from '../types';

export const createPhase = (data: {
  title: string;
  assignedSchedulerId: number;
  status: string;
  applications: Array<{
    type: string;
    institutionName: string;
    qualificationName: string;
    expert1: number;
    expert2: number;
    status: string;
  }>;
}) => 
  request<{
    id: number;
    success: boolean;
    message: string;
  }>({
    url: '/phases',
    method: 'POST',
    data,
  });

export const getPhases = (limit: number = 10, offset: number = 0, params?: Record<string, string | number>) => 
  request<Phase[]>({
    url: '/phases',
    method: 'GET',
    params: {
      ...params,
      limit,
      offset
    },
  });

export const proposeExperts = (phaseId: number, applicationId: number, data: { expert1: number; expert2: number }) => 
  request<{
    success: boolean;
    message: string;
  }>({
    url: `/phases/${phaseId}/applications/${applicationId}`,
    method: 'PUT',
    data,
  });

export const reviewApplication = (phaseId: number, applicationId: number, data: { status: string; rejectionNotes?: string }) => 
  request<{
    success: boolean;
    message: string;
  }>({
    url: `/phases/${phaseId}/applications/${applicationId}/review`,
    method: 'PUT',
    data,
  });