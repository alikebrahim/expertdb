'use client';

import { useState, useEffect } from 'react';
import { ExpertSearch } from './expert-search';
import { Navbar } from '@/components/layout/navbar';
import RequireAuth from '@/components/auth/require-auth';

export default function SearchPage() {
  const [isPageLoading, setIsPageLoading] = useState(true);

  useEffect(() => {
    // Simulate checking data and permissions
    const timer = setTimeout(() => {
      setIsPageLoading(false);
    }, 300);
    
    return () => clearTimeout(timer);
  }, []);

  return (
    <RequireAuth>
      <div className="container py-6 md:py-10">
        <div className="mx-auto">
          {isPageLoading ? (
            <div className="space-y-4">
              <div className="h-10 w-60 bg-muted animate-pulse rounded"></div>
              <div className="h-5 w-full max-w-2xl bg-muted animate-pulse rounded"></div>
              <div className="h-[600px] w-full bg-muted animate-pulse rounded mt-8"></div>
            </div>
          ) : (
            <>
              <h1 className="text-3xl font-bold mb-4 md:mb-6">Search Experts</h1>
              <p className="text-muted-foreground mb-6 md:mb-8 max-w-2xl">
                Search our database of experts by name, institution, area of expertise, or other criteria. 
                Use the filters to narrow down your results and find the perfect match for your needs.
              </p>
              <ExpertSearch />
            </>
          )}
        </div>
      </div>
    </RequireAuth>
  );
}