import { request } from './client';
import { NationalityStats, GrowthStats } from '../types';

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
  request<{
    totalExperts: number;
    activeCount: number;
    bahrainiPercentage: number;
    publishedCount: number;
    publishedRatio: number;
    topAreas: Array<{ name: string; count: number; percentage: number }>;
    engagementsByType: Array<{ name: string; count: number; percentage: number }>;
    yearlyGrowth: Array<{ period: string; count: number; growthRate: number }>;
    mostRequestedExperts: Array<{ expertId: string; name: string; count: number }>;
    lastUpdated: string;
  }>({
    url: '/statistics',
    method: 'GET',
  });

export const getEngagementStats = () => 
  request<Array<{ name: string; count: number; percentage: number }>>({
    url: '/statistics/engagements',
    method: 'GET',
  });
  
export const getAreaStats = () => 
  request<{
    generalAreas: Array<{ name: string; count: number; percentage: number }>;
    topSpecializedAreas: Array<{ name: string; count: number; percentage: number }>;
    bottomSpecializedAreas: Array<{ name: string; count: number; percentage: number }>;
  }>({
    url: '/statistics/areas',
    method: 'GET',
  });