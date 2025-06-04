import { request } from './client';
import { PhaseListResponse } from '../types';

export const createPhase = (data: {
  title: string;
  assignedPlannerId: number;
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
  }>({
    url: '/api/phases',
    method: 'POST',
    data,
  });

export const getPhases = (limit: number = 10, offset: number = 0, params?: Record<string, string | number>) => 
  request<PhaseListResponse>({
    url: '/api/phases',
    method: 'GET',
    params: {
      ...params,
      limit,
      offset
    },
  });

export const proposeExperts = (phaseId: number, applicationId: number, data: { expert1: number; expert2: number }) => 
  request<null>({
    url: `/api/phases/${phaseId}/applications/${applicationId}`,
    method: 'PUT',
    data,
  });

export const reviewApplication = (phaseId: number, applicationId: number, data: { status: string; rejectionNotes?: string }) => 
  request<null>({
    url: `/api/phases/${phaseId}/applications/${applicationId}/review`,
    method: 'PUT',
    data,
  });