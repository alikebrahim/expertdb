import { request } from './client';
import { 
  NationalityStats, 
  GrowthStats, 
  ExpertStats, 
  EngagementStats, 
  AreaStats 
} from '../types';

export const getNationalityStats = () => 
  request<NationalityStats>({
    url: '/api/statistics/nationality',
    method: 'GET',
  });

export const getGrowthStats = (years?: number) => 
  request<GrowthStats[]>({
    url: '/api/statistics/growth',
    method: 'GET',
    params: { years },
  });

export const getOverallStats = () => 
  request<ExpertStats>({
    url: '/api/statistics',
    method: 'GET',
  });

export const getEngagementStats = () => 
  request<EngagementStats>({
    url: '/api/statistics/engagements',
    method: 'GET',
  });
  
export const getAreaStats = () => 
  request<AreaStats>({
    url: '/api/statistics/areas',
    method: 'GET',
  });