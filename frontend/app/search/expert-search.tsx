'use client';

import { useState, useEffect } from 'react';
import Link from 'next/link';
import { Card, CardContent, CardFooter } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Button } from '@/components/ui/button';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Expert, expertAPI } from '@/lib/api';

export function ExpertSearch() {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedArea, setSelectedArea] = useState('');
  const [selectedRole, setSelectedRole] = useState('');
  const [availabilityFilter, setAvailabilityFilter] = useState('');
  const [iscedFilter, setIscedFilter] = useState('');
  const [ratingFilter, setRatingFilter] = useState('');
  
  // Pagination state
  const [currentPage, setCurrentPage] = useState(0);
  const [pageSize, setPageSize] = useState(10);
  const [totalExperts, setTotalExperts] = useState(0);
  const [sortBy, setSortBy] = useState('name');
  const [sortOrder, setSortOrder] = useState('asc');
  
  const [experts, setExperts] = useState<Expert[]>([]);
  const [iscedFields, setIscedFields] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Load ISCED fields for filtering
  useEffect(() => {
    const loadIscedFields = async () => {
      try {
        setIsLoading(true);
        const fields = await expertAPI.getISCEDFields();
        setIscedFields(fields);
      } catch (error) {
        console.error('Failed to load ISCED fields:', error);
        setError('Failed to load filtering options. Please refresh the page.');
      } finally {
        // We still want to attempt to load experts even if ISCED fields fail
        await searchExperts(true);
      }
    };
    
    loadIscedFields();
    // We'll load experts after ISCED fields, so don't call searchExperts here
  }, []);
  
  // Search experts based on filters and pagination
  const searchExperts = async (resetPage = false) => {
    setIsLoading(true);
    setError(null);
    
    // Delay minimum loading time to prevent flickering for fast responses
    const loadingStartTime = Date.now();
    const MIN_LOADING_TIME = 500; // milliseconds
    
    // Reset to first page when filters change
    const pageToUse = resetPage ? 0 : currentPage;
    if (resetPage) {
      setCurrentPage(0);
    }
    
    try {
      const filters: Record<string, any> = {};
      
      if (searchQuery) {
        filters.name = searchQuery;
      }
      
      if (selectedArea) {
        filters.area = selectedArea;
      }
      
      if (selectedRole) {
        filters.role = selectedRole;
      }
      
      if (availabilityFilter) {
        filters.is_available = availabilityFilter === 'available';
      }
      
      if (iscedFilter) {
        filters.isced_field_id = parseInt(iscedFilter);
      }
      
      if (ratingFilter) {
        filters.min_rating = ratingFilter;
      }
      
      // Calculate offset
      const offset = pageToUse * pageSize;
      
      const response = await expertAPI.getAllExperts(
        filters,
        pageSize,
        offset,
        sortBy,
        sortOrder
      );
      
      // Calculate how much time has passed
      const loadingElapsed = Date.now() - loadingStartTime;
      
      // If loading was too fast, wait the remaining time to prevent flickering
      if (loadingElapsed < MIN_LOADING_TIME) {
        await new Promise(resolve => setTimeout(resolve, MIN_LOADING_TIME - loadingElapsed));
      }
      
      setExperts(response.experts || []);
      setTotalExperts(response.pagination.total);
    } catch (error: any) {
      // Enhanced error handling
      if (error.response) {
        // Server responded with an error status
        const errorMessage = error.response.data?.error || `Server error: ${error.response.status}`;
        setError(errorMessage);
        console.error('API error response:', errorMessage, error.response);
      } else if (error.request) {
        // Request was made but got no response
        setError('Unable to connect to the server. Please check your internet connection and try again.');
        console.error('API request error (no response):', error.request);
      } else {
        // Something else caused the error
        setError('An unexpected error occurred. Please try again.');
        console.error('API unexpected error:', error.message);
      }
    } finally {
      setIsLoading(false);
    }
  };
  
  // Handle form submission
  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    searchExperts(true); // Reset to first page when submitting form
  };
  
  // Handle page change
  const handlePageChange = (newPage: number) => {
    setCurrentPage(newPage);
  };
  
  // Handle sort change
  const handleSortChange = (field: string) => {
    if (sortBy === field) {
      // Toggle order if clicking the same field
      setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc');
    } else {
      setSortBy(field);
      setSortOrder('asc');
    }
  };
  
  // Handle page size change
  const handlePageSizeChange = (size: number) => {
    setPageSize(size);
    setCurrentPage(0); // Reset to first page when changing page size
  };
  
  // Fetch experts when pagination or sorting changes
  useEffect(() => {
    // Skip the initial render since we're already loading experts after ISCED fields
    if (experts.length > 0 || totalExperts > 0) {
      searchExperts();
    }
  }, [currentPage, pageSize, sortBy, sortOrder]);
  
  return (
    <div className="grid grid-cols-1 lg:grid-cols-4 gap-6">
      {/* Search filters */}
      <Card className="lg:col-span-1 h-fit sticky top-4">
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4 pt-6">
            <div className="space-y-2">
              <Label htmlFor="search">Search by name or keyword</Label>
              <Input 
                id="search" 
                placeholder="Search experts..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
              />
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="area">Area of expertise</Label>
              <Input 
                id="area" 
                placeholder="e.g., Computer Science, Medicine"
                value={selectedArea}
                onChange={(e) => setSelectedArea(e.target.value)}
              />
            </div>
            
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-1 gap-4">
              <div className="space-y-2">
                <Label htmlFor="role">Role</Label>
                <Select 
                  value={selectedRole} 
                  onValueChange={setSelectedRole}
                >
                  <SelectTrigger id="role">
                    <SelectValue placeholder="Any role" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="">Any</SelectItem>
                    <SelectItem value="Validator">Validator</SelectItem>
                    <SelectItem value="Evaluator">Evaluator</SelectItem>
                    <SelectItem value="Both">Both</SelectItem>
                  </SelectContent>
                </Select>
              </div>
              
              <div className="space-y-2">
                <Label htmlFor="availability">Availability</Label>
                <Select 
                  value={availabilityFilter} 
                  onValueChange={setAvailabilityFilter}
                >
                  <SelectTrigger id="availability">
                    <SelectValue placeholder="Any availability" />
                  </SelectTrigger>
                  <SelectContent>
                    <SelectItem value="">Any</SelectItem>
                    <SelectItem value="available">Available</SelectItem>
                    <SelectItem value="unavailable">Unavailable</SelectItem>
                  </SelectContent>
                </Select>
              </div>
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="isced">ISCED Field</Label>
              <Select 
                value={iscedFilter} 
                onValueChange={setIscedFilter}
              >
                <SelectTrigger id="isced">
                  <SelectValue placeholder="Any field" />
                </SelectTrigger>
                <SelectContent className="max-h-[200px]">
                  <SelectItem value="">Any</SelectItem>
                  {iscedFields.map((field) => (
                    <SelectItem key={field.id} value={field.id ? field.id.toString() : `field-${field.broadCode || Math.random()}`}>
                      {field.broadName}
                    </SelectItem>
                  ))}
                </SelectContent>
              </Select>
            </div>
            
            <div className="space-y-2">
              <Label htmlFor="rating">Minimum Rating</Label>
              <Select 
                value={ratingFilter} 
                onValueChange={setRatingFilter}
              >
                <SelectTrigger id="rating">
                  <SelectValue placeholder="Any rating" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="">Any</SelectItem>
                  <SelectItem value="5">★★★★★ (5.0)</SelectItem>
                  <SelectItem value="4">★★★★☆ (4.0+)</SelectItem>
                  <SelectItem value="3">★★★☆☆ (3.0+)</SelectItem>
                  <SelectItem value="2">★★☆☆☆ (2.0+)</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </CardContent>
          <CardFooter className="border-t pt-6">
            <Button type="submit" className="w-full" disabled={isLoading}>
              {isLoading ? (
                <div className="flex items-center gap-2">
                  <div className="animate-spin rounded-full h-4 w-4 border-2 border-current border-t-transparent"></div>
                  <span>Searching...</span>
                </div>
              ) : 'Search Experts'}
            </Button>
          </CardFooter>
        </form>
      </Card>
      
      {/* Search results */}
      <div className="lg:col-span-3">
        {error && (
          <div className="bg-destructive/10 border border-destructive text-destructive p-4 rounded-md mb-6">
            {error}
          </div>
        )}
        
        {isLoading ? (
          <div className="space-y-4">
            {/* Skeleton loading for results */}
            <div className="flex justify-between items-center mb-4">
              <div className="h-7 w-32 bg-muted animate-pulse rounded"></div>
              <div className="flex gap-2">
                <div className="h-9 w-20 bg-muted animate-pulse rounded"></div>
                <div className="h-9 w-32 bg-muted animate-pulse rounded"></div>
              </div>
            </div>
            
            {Array.from({ length: 3 }).map((_, i) => (
              <Card key={i} className="overflow-hidden">
                <CardContent className="pt-6">
                  <div className="flex flex-col md:flex-row justify-between gap-4">
                    <div className="space-y-3 w-full">
                      <div className="h-6 bg-muted animate-pulse rounded w-1/3"></div>
                      <div className="h-4 bg-muted animate-pulse rounded w-1/2"></div>
                      <div className="flex gap-2">
                        <div className="h-6 bg-muted animate-pulse rounded w-20"></div>
                        <div className="h-6 bg-muted animate-pulse rounded w-24"></div>
                      </div>
                    </div>
                    <div className="flex flex-col gap-2 items-end">
                      <div className="h-4 bg-muted animate-pulse rounded w-16"></div>
                      <div className="h-4 bg-muted animate-pulse rounded w-20"></div>
                      <div className="h-8 bg-muted animate-pulse rounded w-28"></div>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        ) : (
          <>
            {/* Results header with sort and page size options */}
            <div className="flex flex-col sm:flex-row justify-between items-start sm:items-center mb-4 gap-2">
              <h2 className="text-xl font-semibold">
                {totalExperts} {totalExperts === 1 ? 'Expert' : 'Experts'} Found
              </h2>
              
              <div className="flex flex-wrap gap-2">
                <div className="flex items-center gap-2">
                  <Label htmlFor="pageSize" className="text-sm whitespace-nowrap">Show:</Label>
                  <Select
                    value={String(pageSize)}
                    onValueChange={(value) => handlePageSizeChange(Number(value))}
                  >
                    <SelectTrigger id="pageSize" className="w-[70px]">
                      <SelectValue placeholder="10" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="5">5</SelectItem>
                      <SelectItem value="10">10</SelectItem>
                      <SelectItem value="25">25</SelectItem>
                      <SelectItem value="50">50</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
                
                <div className="flex items-center gap-2">
                  <Label htmlFor="sortBy" className="text-sm whitespace-nowrap">Sort by:</Label>
                  <Select
                    value={sortBy}
                    onValueChange={(value) => handleSortChange(value)}
                  >
                    <SelectTrigger id="sortBy" className="w-[120px]">
                      <SelectValue placeholder="Name" />
                    </SelectTrigger>
                    <SelectContent>
                      <SelectItem value="name">Name</SelectItem>
                      <SelectItem value="institution">Institution</SelectItem>
                      <SelectItem value="role">Role</SelectItem>
                      <SelectItem value="created_at">Date Added</SelectItem>
                    </SelectContent>
                  </Select>
                  
                  <Button 
                    variant="outline" 
                    size="icon"
                    onClick={() => setSortOrder(sortOrder === 'asc' ? 'desc' : 'asc')}
                    title={sortOrder === 'asc' ? 'Ascending' : 'Descending'}
                    className="h-9 w-9"
                  >
                    {sortOrder === 'asc' ? '↑' : '↓'}
                  </Button>
                </div>
              </div>
            </div>
            
            {/* Expert results list */}
            {experts.length === 0 ? (
              <Card>
                <CardContent className="pt-6 text-center py-12">
                  <h3 className="text-lg font-medium mb-2">No experts found</h3>
                  <p className="text-muted-foreground">
                    Try adjusting your search filters to find experts in our database.
                  </p>
                </CardContent>
              </Card>
            ) : (
              <div className="grid grid-cols-1 gap-4">
                {experts.map((expert) => (
                  <Card key={expert.id} className="hover:bg-muted/50 transition-colors overflow-hidden">
                    <CardContent className="p-4 sm:p-6 sm:pt-6">
                      <div className="flex flex-col md:flex-row justify-between gap-4">
                        <div>
                          <h3 className="text-lg font-semibold">{expert.name}</h3>
                          <p className="text-sm text-muted-foreground">
                            {expert.designation} at {expert.institution}
                          </p>
                          
                          <div className="mt-2 flex flex-wrap gap-2">
                            {expert.generalArea && (
                              <span className="inline-flex items-center rounded-md bg-blue-50 px-2 py-1 text-xs font-medium text-blue-700 ring-1 ring-inset ring-blue-700/10">
                                {expert.generalArea}
                              </span>
                            )}
                            {expert.specializedArea && (
                              <span className="inline-flex items-center rounded-md bg-green-50 px-2 py-1 text-xs font-medium text-green-700 ring-1 ring-inset ring-green-700/10">
                                {expert.specializedArea}
                              </span>
                            )}
                            {expert.iscedField && (
                              <span className="inline-flex items-center rounded-md bg-purple-50 px-2 py-1 text-xs font-medium text-purple-700 ring-1 ring-inset ring-purple-700/10">
                                {expert.iscedField.broadName}
                              </span>
                            )}
                            {expert.role && (
                              <span className="inline-flex items-center rounded-md bg-amber-50 px-2 py-1 text-xs font-medium text-amber-700 ring-1 ring-inset ring-amber-700/10">
                                {expert.role}
                              </span>
                            )}
                          </div>
                        </div>
                        
                        <div className="flex md:flex-col justify-between md:justify-center md:items-end gap-2">
                          <div className="flex md:flex-col items-center md:items-end text-sm gap-1">
                            <span className={expert.isAvailable ? 'text-green-600 font-medium' : 'text-red-600 font-medium'}>
                              {expert.isAvailable ? 'Available' : 'Unavailable'}
                            </span>
                            {expert.rating && expert.rating !== 'N/A' && (
                              <span className="flex items-center gap-1">
                                <span className="text-amber-500">{'★'.repeat(Math.floor(parseFloat(expert.rating)))}</span>
                                <span>{expert.rating}/5</span>
                              </span>
                            )}
                          </div>
                          <Button size="sm" asChild className="whitespace-nowrap">
                            <Link href={`/expert/${expert.id}`}>View Profile</Link>
                          </Button>
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))}
              </div>
            )}
            
            {/* Pagination controls */}
            {totalExperts > pageSize && (
              <div className="flex flex-col sm:flex-row justify-between items-center mt-6 gap-4">
                <div className="text-sm text-muted-foreground order-2 sm:order-1">
                  Showing {Math.min(currentPage * pageSize + 1, totalExperts)} to {Math.min((currentPage + 1) * pageSize, totalExperts)} of {totalExperts} experts
                </div>
                
                <div className="flex gap-1 order-1 sm:order-2">
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handlePageChange(0)}
                    disabled={currentPage === 0 || isLoading}
                  >
                    First
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handlePageChange(currentPage - 1)}
                    disabled={currentPage === 0 || isLoading}
                  >
                    Previous
                  </Button>
                  
                  <div className="flex items-center px-2 sm:px-4 text-sm sm:text-base">
                    <span className="hidden sm:inline">Page </span>{currentPage + 1} of {Math.max(1, Math.ceil(totalExperts / pageSize))}
                  </div>
                  
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handlePageChange(currentPage + 1)}
                    disabled={currentPage >= Math.ceil(totalExperts / pageSize) - 1 || isLoading}
                  >
                    Next
                  </Button>
                  <Button
                    variant="outline"
                    size="sm"
                    onClick={() => handlePageChange(Math.ceil(totalExperts / pageSize) - 1)}
                    disabled={currentPage >= Math.ceil(totalExperts / pageSize) - 1 || isLoading}
                  >
                    Last
                  </Button>
                </div>
              </div>
            )}
          </>
        )}
      </div>
    </div>
  );
}