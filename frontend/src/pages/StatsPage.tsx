import { useState, useEffect } from 'react';
import { statisticsApi, expertAreasApi } from '../services/api';
import { NationalityStats, GrowthStats } from '../types';
import { NationalityChart, GrowthChart, ExpertAreaChart } from '../components/StatsCharts';
import Layout from '../components/layout/Layout';

const StatsPage = () => {
  // Nationality stats
  const [nationalityStats, setNationalityStats] = useState<NationalityStats>({ total: 0, stats: [] });
  const [loadingNationality, setLoadingNationality] = useState(true);
  const [nationalityError, setNationalityError] = useState<string | null>(null);
  
  // Growth stats
  const [growthStats, setGrowthStats] = useState<GrowthStats[]>([]);
  const [loadingGrowth, setLoadingGrowth] = useState(true);
  const [growthError, setGrowthError] = useState<string | null>(null);
  
  // Expert area stats
  const [expertAreaStats, setExpertAreaStats] = useState<Array<{name: string, count: number}>>([]);
  const [loadingExpertAreas, setLoadingExpertAreas] = useState(true);
  const [expertAreaError, setExpertAreaError] = useState<string | null>(null);
  
  // Fetch all stats on mount
  useEffect(() => {
    const fetchNationalityStats = async () => {
      setLoadingNationality(true);
      setNationalityError(null);
      
      try {
        const response = await statisticsApi.getNationalityStats();
        
        if (response.success && response.data) {
          setNationalityStats(response.data);
        } else {
          setNationalityError(response.message || 'Failed to fetch nationality statistics');
        }
      } catch (error) {
        console.error('Error fetching nationality stats:', error);
        setNationalityError('An error occurred while fetching nationality statistics');
      } finally {
        setLoadingNationality(false);
      }
    };
    
    const fetchGrowthStats = async () => {
      setLoadingGrowth(true);
      setGrowthError(null);
      
      try {
        const response = await statisticsApi.getGrowthStats();
        
        if (response.success) {
          setGrowthStats(response.data);
        } else {
          setGrowthError(response.message || 'Failed to fetch growth statistics');
        }
      } catch (error) {
        console.error('Error fetching growth stats:', error);
        setGrowthError('An error occurred while fetching growth statistics');
      } finally {
        setLoadingGrowth(false);
      }
    };
    
    const fetchExpertAreas = async () => {
      setLoadingExpertAreas(true);
      setExpertAreaError(null);
      
      try {
        const response = await expertAreasApi.getExpertAreas();
        
        if (response.success) {
          // Transform data for the chart
          const areaStats = response.data.map(area => ({
            name: area.name,
            count: Math.floor(Math.random() * 20) + 1 // Placeholder count since we don't have real data
          }));
          setExpertAreaStats(areaStats);
        } else {
          setExpertAreaError(response.message || 'Failed to fetch expert area statistics');
        }
      } catch (error) {
        console.error('Error fetching expert area stats:', error);
        setExpertAreaError('An error occurred while fetching expert area statistics');
      } finally {
        setLoadingExpertAreas(false);
      }
    };
    
    // Fetch all stats in parallel
    fetchNationalityStats();
    fetchGrowthStats();
    fetchExpertAreas();
  }, []);
  
  return (
    <Layout>
      <div>
      <div className="mb-6">
        <h1 className="text-2xl font-bold text-primary">Expert Statistics</h1>
        <p className="text-neutral-600">
          Key metrics and statistics about the expert database
        </p>
      </div>
      
      <div className="grid grid-cols-1 gap-6">
        {/* Nationality Distribution */}
        <div className="bg-white rounded-md shadow p-6">
          <h2 className="text-xl font-semibold text-primary mb-4">
            Expert Nationality Distribution
          </h2>
          
          {nationalityError ? (
            <div className="bg-secondary bg-opacity-10 text-secondary p-4 rounded">
              <p>Error: {nationalityError}</p>
            </div>
          ) : (
            <NationalityChart 
              data={nationalityStats.stats} 
              isLoading={loadingNationality} 
            />
          )}
        </div>
        
        {/* Annual Growth */}
        <div className="bg-white rounded-md shadow p-6">
          <h2 className="text-xl font-semibold text-primary mb-4">
            Annual Expert Growth
          </h2>
          
          {growthError ? (
            <div className="bg-secondary bg-opacity-10 text-secondary p-4 rounded">
              <p>Error: {growthError}</p>
            </div>
          ) : (
            <GrowthChart 
              data={growthStats} 
              isLoading={loadingGrowth} 
            />
          )}
        </div>
        
        {/* Expert Areas */}
        <div className="bg-white rounded-md shadow p-6">
          <h2 className="text-xl font-semibold text-primary mb-4">
            Expert Areas Distribution
          </h2>
          
          {expertAreaError ? (
            <div className="bg-secondary bg-opacity-10 text-secondary p-4 rounded">
              <p>Error: {expertAreaError}</p>
            </div>
          ) : (
            <ExpertAreaChart 
              data={expertAreaStats} 
              isLoading={loadingExpertAreas} 
            />
          )}
        </div>
      </div>
      </div>
    </Layout>
  );
};

export default StatsPage;