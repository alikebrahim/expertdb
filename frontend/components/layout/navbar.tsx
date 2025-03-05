'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import Image from 'next/image';
import { Button } from '@/components/ui/button';
import { authAPI } from '@/lib/api';
import { useRouter } from 'next/navigation';

export function Navbar() {
  const [isAuthenticated, setIsAuthenticated] = useState(false);
  const [user, setUser] = useState<any>(null);
  const router = useRouter();

  useEffect(() => {
    // Check authentication status on client side
    const checkAuth = () => {
      setIsAuthenticated(authAPI.isAuthenticated());
      setUser(authAPI.getUser());
    };

    checkAuth();
    // Add event listener for storage changes (for multi-tab support)
    window.addEventListener('storage', checkAuth);
    
    return () => {
      window.removeEventListener('storage', checkAuth);
    };
  }, []);

  const handleLogout = () => {
    authAPI.logout();
    setIsAuthenticated(false);
    setUser(null);
    router.push('/');
  };

  return (
    <header className="bg-[#133566] text-white">
      <div className="container flex h-16 items-center justify-between px-6">
        <div className="flex items-center gap-6">
          <Link href="/" className="flex items-center">
            <Image 
              src="/images/logo/BQA - Horizontal Logo.svg"
              alt="BQA Logo"
              width={150}
              height={40}
              className="h-10 w-auto mx-5" 
              priority
            />
          </Link>
          <nav className="hidden md:flex gap-6">
            <Link href="/" className="text-sm font-medium text-white hover:bg-[#1B4882] px-3 py-2 rounded-sm transition-colors duration-200 border-b-2 border-transparent hover:border-b-2 hover:border-[#DC8335]">
              Home
            </Link>
            <Link href="/search" className="text-sm font-medium text-white hover:bg-[#1B4882] px-3 py-2 rounded-sm transition-colors duration-200 border-b-2 border-transparent hover:border-b-2 hover:border-[#DC8335]">
              Search Experts
            </Link>
            <Link href="/request" className="text-sm font-medium text-white hover:bg-[#1B4882] px-3 py-2 rounded-sm transition-colors duration-200 border-b-2 border-transparent hover:border-b-2 hover:border-[#DC8335]">
              Submit Request
            </Link>
            {isAuthenticated && (
              <>
                <Link href="/panel" className="text-sm font-medium text-white hover:bg-[#1B4882] px-3 py-2 rounded-sm transition-colors duration-200 border-b-2 border-transparent hover:border-b-2 hover:border-[#DC8335]">
                  AI Panel Suggestion
                </Link>
                <Link href="/statistics" className="text-sm font-medium text-white hover:bg-[#1B4882] px-3 py-2 rounded-sm transition-colors duration-200 border-b-2 border-transparent hover:border-b-2 hover:border-[#DC8335]">
                  Statistics
                </Link>
              </>
            )}
          </nav>
        </div>
        <div className="flex items-center gap-4">
          {isAuthenticated ? (
            <div className="flex items-center gap-4">
              <span className="text-sm font-medium hidden sm:inline-block text-white">
                {user?.name || 'User'}
              </span>
              <Button 
                variant="outlineInverse" 
                onClick={handleLogout}
              >
                Logout
              </Button>
            </div>
          ) : (
            <Button 
              variant="outlineInverse" 
              asChild
            >
              <Link href="/login">Login</Link>
            </Button>
          )}
        </div>
      </div>
    </header>
  );
}