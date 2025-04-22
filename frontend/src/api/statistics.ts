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
    url: '/statistics/nationality',
    method: 'GET',
  });

export const getGrowthStats = (years?: number) => 
  request<GrowthStats[]>({
    url: '/statistics/growth',
    method: 'GET',
    params: { years },
  });

export const getOverallStats = () => 
  request<ExpertStats>({
    url: '/statistics',
    method: 'GET',
  });

export const getEngagementStats = () => 
  request<EngagementStats>({
    url: '/statistics/engagements',
    method: 'GET',
  });
  
export const getAreaStats = () => 
  request<AreaStats>({
    url: '/statistics/areas',
    method: 'GET',
  });