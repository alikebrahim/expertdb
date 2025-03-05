'use client';

import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { authAPI } from '@/lib/api';

export default function Home() {
  const router = useRouter();
  const [isChecking, setIsChecking] = useState(true);
  
  // Check authentication and redirect to appropriate page
  useEffect(() => {
    const checkAuth = async () => {
      try {
        // Set a timeout to ensure we don't hang indefinitely
        const timeoutId = setTimeout(() => {
          console.log('Authentication check timed out, redirecting to login');
          router.push('/login');
          setIsChecking(false);
        }, 2000); // 2 second timeout
        
        // Check authentication status (synchronous operation)
        const isAuthenticated = authAPI.isAuthenticated();
        
        // Clear the timeout since we've completed the check
        clearTimeout(timeoutId);
        
        // Redirect based on authentication status
        if (isAuthenticated) {
          router.push('/search');
        } else {
          router.push('/login');
        }
        
        // Always set checking to false once we've attempted redirection
        setIsChecking(false);
      } catch (error) {
        console.error('Error during authentication check:', error);
        // If any error occurs, redirect to login and stop checking
        router.push('/login');
        setIsChecking(false);
      }
    };
    
    // Start the authentication check
    checkAuth();
    
    // Cleanup function to ensure we don't have memory leaks
    return () => {
      setIsChecking(false);
    };
  }, [router]);
  
  // Return during authentication check to avoid flash of content
  return (
    <div className="flex items-center justify-center min-h-screen">
      {isChecking && (
        <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-primary"></div>
      )}
    </div>
  );
}