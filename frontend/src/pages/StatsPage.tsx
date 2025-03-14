import { useState, useEffect } from 'react';
import { statisticsApi } from '../services/api';
import { NationalityStats, GrowthStats, IscedStats } from '../types';
import { NationalityChart, GrowthChart, IscedChart } from '../components/StatsCharts';

const StatsPage = () => {
  // Nationality stats
  const [nationalityStats, setNationalityStats] = useState<NationalityStats[]>([]);
  const [loadingNationality, setLoadingNationality] = useState(true);
  const [nationalityError, setNationalityError] = useState<string | null>(null);
  
  // Growth stats
  const [growthStats, setGrowthStats] = useState<GrowthStats[]>([]);
  const [loadingGrowth, setLoadingGrowth] = useState(true);
  const [growthError, setGrowthError] = useState<string | null>(null);
  
  // ISCED stats
  const [iscedStats, setIscedStats] = useState<IscedStats[]>([]);
  const [loadingIsced, setLoadingIsced] = useState(true);
  const [iscedError, setIscedError] = useState<string | null>(null);
  
  // Fetch all stats on mount
  useEffect(() => {
    const fetchNationalityStats = async () => {
      setLoadingNationality(true);
      setNationalityError(null);
      
      try {
        const response = await statisticsApi.getNationalityStats();
        
        if (response.success) {
          // Transform data for pie chart if needed
          const stats = response.data.stats || response.data;
          setNationalityStats(stats);
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
    
    const fetchIscedStats = async () => {
      setLoadingIsced(true);
      setIscedError(null);
      
      try {
        const response = await statisticsApi.getIscedStats();
        
        if (response.success) {
          setIscedStats(response.data);
        } else {
          setIscedError(response.message || 'Failed to fetch ISCED statistics');
        }
      } catch (error) {
        console.error('Error fetching ISCED stats:', error);
        setIscedError('An error occurred while fetching ISCED statistics');
      } finally {
        setLoadingIsced(false);
      }
    };
    
    // Fetch all stats in parallel
    fetchNationalityStats();
    fetchGrowthStats();
    fetchIscedStats();
  }, []);
  
  return (
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
              data={nationalityStats} 
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
        
        {/* ISCED Categories */}
        <div className="bg-white rounded-md shadow p-6">
          <h2 className="text-xl font-semibold text-primary mb-4">
            ISCED Field Distribution
          </h2>
          
          {iscedError ? (
            <div className="bg-secondary bg-opacity-10 text-secondary p-4 rounded">
              <p>Error: {iscedError}</p>
            </div>
          ) : (
            <IscedChart 
              data={iscedStats} 
              isLoading={loadingIsced} 
            />
          )}
        </div>
      </div>
    </div>
  );
};

export default StatsPage;