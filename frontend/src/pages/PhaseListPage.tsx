
import { useState, useEffect, useCallback } from 'react';
import { phaseApi } from '../services/api';
import { Phase } from '../types';
import { Button, Table, Input } from '../components/ui';
import { useNavigate } from 'react-router-dom';

const PhaseListPage = () => {
  const [phases, setPhases] = useState<Phase[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [pagination, setPagination] = useState({
    page: 1,
    limit: 10,
    total: 0,
    totalPages: 0
  });
  const [filters, setFilters] = useState({
    status: '',
    search: ''
  });
  const navigate = useNavigate();

  const fetchPhases = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const params: Record<string, string> = {};
      if (filters.status) params.status = filters.status;
      if (filters.search) params.search = filters.search;
      
      const response = await phaseApi.getPhases(pagination.page, pagination.limit, params);
      
      if (response.success) {
        setPhases(response.data.data);
        setPagination({
          page: pagination.page,
          limit: pagination.limit,
          total: response.data.total,
          totalPages: response.data.totalPages
        });
      } else {
        setError(response.message || 'Failed to load phases');
      }
    } catch (error) {
      console.error('Error fetching phases:', error);
      setError('An error occurred while loading phases');
    } finally {
      setIsLoading(false);
    }
  }, [pagination.page, pagination.limit, filters]);

  useEffect(() => {
    fetchPhases();
  }, [fetchPhases]);

  const handlePageChange = (newPage: number) => {
    setPagination(prev => ({ ...prev, page: newPage }));
  };

  const handleFilterChange = (e: React.ChangeEvent<HTMLSelectElement | HTMLInputElement>) => {
    const { name, value } = e.target;
    setFilters(prev => ({ ...prev, [name]: value }));
    setPagination(prev => ({ ...prev, page: 1 }));
  };

  const handleSearch = (e: React.FormEvent) => {
    e.preventDefault();
    fetchPhases();
  };

  const resetFilters = () => {
    setFilters({
      status: '',
      search: ''
    });
    setPagination(prev => ({ ...prev, page: 1 }));
  };

  const handleViewPhase = (phaseId: number) => {
    navigate(`/phases/${phaseId}`);
  };
  
  const handleCreatePhase = () => {
    navigate('/phases/create');
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString();
  };

  const getStatusBadgeColor = (status: string) => {
    switch (status) {
      case 'planning':
        return 'bg-yellow-100 text-yellow-800';
      case 'in_progress':
        return 'bg-blue-100 text-blue-800';
      case 'completed':
        return 'bg-green-100 text-green-800';
      case 'archived':
        return 'bg-gray-100 text-gray-800';
      default:
        return 'bg-gray-100 text-gray-800';
    }
  };

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <div className="mb-6">
        <h1 className="text-2xl font-bold mb-2">Phase Planning</h1>
        <p className="text-gray-600">
          Manage phase plans, track their status, and create new phases.
        </p>
      </div>
      
      <div className="bg-white rounded-md shadow-sm p-6 mb-6">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-xl font-semibold">Phases</h2>
          <Button
            variant="primary"
            onClick={handleCreatePhase}
          >
            Create Phase
          </Button>
        </div>
        
        <form onSubmit={handleSearch} className="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Status</label>
            <select
              name="status"
              value={filters.status}
              onChange={handleFilterChange}
              className="w-full px-3 py-2 border border-gray-300 rounded-md"
            >
              <option value="">All Statuses</option>
              <option value="planning">Planning</option>
              <option value="in_progress">In Progress</option>
              <option value="completed">Completed</option>
              <option value="archived">Archived</option>
            </select>
          </div>
          
          <div>
            <label className="block text-sm font-medium text-gray-700 mb-1">Search</label>
            <Input
              name="search"
              value={filters.search}
              onChange={handleFilterChange}
              placeholder="Search by title or planner..."
            />
          </div>
          
          <div className="flex items-end gap-2">
            <Button
              type="submit"
              variant="secondary"
            >
              Apply Filters
            </Button>
            <Button
              type="button"
              variant="outline"
              onClick={resetFilters}
            >
              Reset
            </Button>
          </div>
        </form>
        
        {error && (
          <div className="bg-red-50 border border-red-200 text-red-600 p-4 rounded-md mb-4">
            {error}
          </div>
        )}
        
        {isLoading ? (
          <div className="text-center p-8">
            <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary mx-auto"></div>
            <p className="mt-4 text-gray-600">Loading phases...</p>
          </div>
        ) : phases.length === 0 ? (
          <div className="text-center p-8 bg-gray-50 rounded-md">
            <p className="text-gray-500">No phases found.</p>
            <p className="text-gray-500 mt-2">Try adjusting your filters or create a new phase.</p>
          </div>
        ) : (
          <>
            <Table>
              <Table.Header>
                <Table.Row>
                  <Table.HeaderCell>Title</Table.HeaderCell>
                  <Table.HeaderCell>Planner</Table.HeaderCell>
                  <Table.HeaderCell>Status</Table.HeaderCell>
                  <Table.HeaderCell>Created At</Table.HeaderCell>
                  <Table.HeaderCell>Applications</Table.HeaderCell>
                  <Table.HeaderCell>Actions</Table.HeaderCell>
                </Table.Row>
              </Table.Header>
              <Table.Body>
                {phases.map(phase => (
                  <Table.Row key={phase.id}>
                    <Table.Cell>{phase.title}</Table.Cell>
                    <Table.Cell>{phase.plannerName}</Table.Cell>
                    <Table.Cell>
                      <span className={`px-2 py-1 rounded-full text-xs ${getStatusBadgeColor(phase.status)}`}>
                        {phase.status.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())}
                      </span>
                    </Table.Cell>
                    <Table.Cell>{formatDate(phase.createdAt)}</Table.Cell>
                    <Table.Cell>{phase.applications.length}</Table.Cell>
                    <Table.Cell>
                      <Button 
                        variant="outline" 
                        size="sm" 
                        onClick={() => handleViewPhase(phase.id)}
                      >
                        View
                      </Button>
                    </Table.Cell>
                  </Table.Row>
                ))}
              </Table.Body>
            </Table>
            
            {pagination.totalPages > 1 && (
              <div className="flex justify-center mt-6">
                <div className="flex space-x-2">
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={pagination.page === 1}
                    onClick={() => handlePageChange(pagination.page - 1)}
                  >
                    Previous
                  </Button>
                  <div className="flex items-center px-3">
                    Page {pagination.page} of {pagination.totalPages}
                  </div>
                  <Button
                    variant="outline"
                    size="sm"
                    disabled={pagination.page === pagination.totalPages}
                    onClick={() => handlePageChange(pagination.page + 1)}
                  >
                    Next
                  </Button>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
};

export default PhaseListPage;
