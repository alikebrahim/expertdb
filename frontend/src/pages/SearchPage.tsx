import React, { useState, useEffect } from 'react';
import * as expertsApi from '../api/experts';
import { useUI } from '../hooks/useUI';
import Layout from '../components/layout/Layout';
import { ExpertTable, SortConfig } from '../components/tables';
import { Button } from '../components/ui';
import { getCachedExperts, setCachedExperts, isCacheValid } from '../utils/expertCache';
import { fuzzySearch } from '../utils/fuzzySearch';
import { sortExperts } from '../utils/tableSorting';
import { Expert } from '../types';

const SearchPage = () => {
  const [searchTerm, setSearchTerm] = useState('');
  const [allExperts, setAllExperts] = useState<Expert[]>([]);
  const [filteredExperts, setFilteredExperts] = useState<Expert[]>([]);
  const [dataLoaded, setDataLoaded] = useState(false);
  const [sortConfig, setSortConfig] = useState<SortConfig>({ field: '', direction: 'asc' });
  const [showAdvancedOptions, setShowAdvancedOptions] = useState(false);
  
  const { addNotification } = useUI();
  
  // Load all experts once on component mount
  useEffect(() => {
    const loadAllExperts = async () => {
      try {
        // Check cache first
        const cached = getCachedExperts();
        if (cached && isCacheValid(cached.timestamp)) {
          setAllExperts(cached.data);
          setFilteredExperts(cached.data);
          setDataLoaded(true);
          return;
        }
        
        // Fetch all experts without pagination
        const response = await expertsApi.getExperts(10000, 0, {});
        if (response.success && response.data) {
          // Handle different possible response structures
          let experts: Expert[] = [];
          
          if (response.data.experts) {
            experts = response.data.experts;
          } else if (Array.isArray(response.data)) {
            experts = response.data;
          } else if (response.data.data && Array.isArray(response.data.data)) {
            experts = response.data.data;
          }
          
          setAllExperts(experts);
          setFilteredExperts(experts);
          
          // Cache the data
          setCachedExperts(experts);
        }
      } catch (error) {
        console.error('Failed to load experts:', error);
        addNotification({
          type: 'error',
          message: 'Failed to load expert database. Please try again.',
          duration: 5000
        });
      } finally {
        setDataLoaded(true);
      }
    };
    
    loadAllExperts();
  }, [addNotification]);

  // Apply search and sorting
  useEffect(() => {
    if (!dataLoaded) return;
    
    let filtered = allExperts;
    
    // Apply fuzzy search if search term exists
    if (searchTerm.trim()) {
      filtered = fuzzySearch(
        allExperts,
        searchTerm,
        (expert) => [
          expert.name,
          expert.institution || expert.affiliation || '',
          expert.generalAreaName || '',
          expert.specializedArea || ''
        ]
      );
    }
    
    // Apply sorting
    if (sortConfig.field) {
      filtered = sortExperts(filtered, sortConfig);
    }
    
    setFilteredExperts(filtered);
  }, [allExperts, searchTerm, sortConfig, dataLoaded]);

  const handleSort = (field: string) => {
    setSortConfig(prevConfig => {
      // If clicking on the same field, toggle direction
      if (prevConfig.field === field) {
        return {
          ...prevConfig,
          direction: prevConfig.direction === 'asc' ? 'desc' : 'asc'
        };
      } else {
        // If clicking on a new field, sort ascending by default
        return {
          field,
          direction: 'asc'
        };
      }
    });
  };

  const handleReset = () => {
    setSearchTerm('');
    setSortConfig({ field: '', direction: 'asc' });
    setFilteredExperts(allExperts);
    
    addNotification({
      type: 'info',
      message: 'Search and filters reset',
      duration: 2000
    });
  };

  const handleSearchChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSearchTerm(e.target.value);
  };

  return (
    <Layout>
      <div className="space-y-6">
        {/* Page Header */}
        <div>
          <h1 className="text-2xl font-bold text-primary">Expert Database</h1>
          <p className="text-neutral-600">
            Search and browse experts by name, institution, or area of expertise
          </p>
        </div>

        {/* Search Controls */}
        <div className="bg-white rounded-lg shadow-sm border border-neutral-200 p-4">
          <div className="space-y-4">
            {/* Global Search */}
            <div className="relative">
              <div className="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg className="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
              <input
                type="text"
                placeholder="Search experts by name, institution, or area..."
                value={searchTerm}
                onChange={handleSearchChange}
                className="w-full pl-10 pr-4 py-3 border border-neutral-300 rounded-lg focus:ring-2 focus:ring-primary focus:border-transparent"
              />
            </div>

            {/* Advanced Options Toggle */}
            <div className="flex items-center justify-between">
              <button
                onClick={() => setShowAdvancedOptions(!showAdvancedOptions)}
                className="flex items-center space-x-2 text-sm text-primary hover:text-primary-dark"
              >
                <span>Advanced Search Options</span>
                <svg 
                  className={`h-4 w-4 transition-transform ${showAdvancedOptions ? 'rotate-180' : ''}`}
                  fill="none" 
                  stroke="currentColor" 
                  viewBox="0 0 24 24"
                >
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                </svg>
              </button>
              
              <div className="flex items-center space-x-4">
                <span className="text-sm text-neutral-600">
                  {filteredExperts.length} of {allExperts.length} experts
                </span>
                <Button 
                  variant="outline" 
                  size="sm" 
                  onClick={handleReset}
                >
                  Reset
                </Button>
              </div>
            </div>

            {/* Advanced Options Panel */}
            {showAdvancedOptions && (
              <div className="border-t border-neutral-200 pt-4">
                <div className="bg-neutral-50 rounded-lg p-4">
                  <h3 className="text-sm font-medium text-neutral-900 mb-3">Column Display Options</h3>
                  <p className="text-sm text-neutral-600">
                    Use the column selector in the table below to customize which fields are displayed.
                    Available optional fields include contact details (phone, email) and CV download links.
                  </p>
                </div>
              </div>
            )}
          </div>
        </div>

        {/* Results Table */}
        <div className="bg-white rounded-lg shadow-sm border border-neutral-200">
          <div className="overflow-x-auto max-w-full">
            {dataLoaded ? (
              <ExpertTable 
                experts={filteredExperts} 
                isLoading={false}
                error={null}
                sortConfig={sortConfig}
                onSort={handleSort}
                showColumnSelector={true}
              />
            ) : (
              <div className="flex flex-col items-center justify-center py-16">
                <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full"></div>
                <p className="mt-4 text-neutral-600">Loading experts...</p>
              </div>
            )}
          </div>
        </div>
      </div>
    </Layout>
  );
};

export default SearchPage;