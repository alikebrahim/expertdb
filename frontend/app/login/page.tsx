'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import Image from 'next/image';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { authAPI } from '@/lib/api';

export default function LoginPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const router = useRouter();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setIsLoading(true);
    setError(null);

    try {
      const response = await authAPI.login(email, password);

      // Store token in localStorage
      localStorage.setItem('token', response.token);
      localStorage.setItem('user', JSON.stringify(response.user));

      // Trigger storage event for other tabs
      window.dispatchEvent(new Event('storage'));

      // Redirect to expert search page
      router.push('/search');
    } catch (error: any) {
      // Improved error handling with more specific messages
      if (error.response) {
        // Server responded with an error
        setError(error.response.data?.error || `Server error: ${error.response.status}`);
      } else if (error.request) {
        // Request was made but no response received
        setError('Server not responding. Please try again later.');
      } else {
        // Request setup error
        setError('Login failed. Please check your credentials and try again.');
      }
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen p-4 bg-gradient-to-b from-background to-muted">
      <Card className="w-full max-w-md border border-[#E5E7EB] shadow-[0_2px_4px_rgba(0,0,0,0.05)]">
        <CardHeader className="space-y-1">
          <div className="flex justify-center mb-4">
            <Image 
              src="/images/logo/BQA - Horizontal Logo with Descriptor.svg"
              alt="BQA Logo"
              width={180}
              height={40}
              className="h-12 w-auto"
              priority
            />
          </div>
          <CardTitle className="text-2xl font-bold text-center text-[#133566]">Expert Database Login</CardTitle>
          <CardDescription className="text-center">
            Enter your credentials to access the expert database
          </CardDescription>
        </CardHeader>
        <form onSubmit={handleSubmit}>
          <CardContent className="space-y-4">
            <div className="space-y-2">
              <Label htmlFor="email" className="text-primary/80">Email</Label>
              <Input
                id="email"
                type="email"
                placeholder="admin@expertdb.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="border-[#D1D5DB] focus-visible:border-[#1B4882] focus-visible:border-2 rounded-[4px] px-2 py-1"
                aria-label="Email address"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="password" className="text-primary/80">Password</Label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="border-[#D1D5DB] focus-visible:border-[#1B4882] focus-visible:border-2 rounded-[4px] px-2 py-1"
                aria-label="Password"
              />
            </div>
            {error && (
              <div 
                className="bg-[#FFEBEB] border border-[#FF4040] text-[#FF4040] p-3 rounded-md text-sm"
                role="alert"
                aria-live="assertive"
              >
                {error}
              </div>
            )}
          </CardContent>
          <CardFooter>
            <Button 
              type="submit" 
              className="w-full bg-[#133566] hover:bg-[#1B4882] text-white rounded-[4px] h-10 px-4 py-2" 
              disabled={isLoading}
            >
              {isLoading ? 'Logging in...' : 'Sign In'}
            </Button>
          </CardFooter>
        </form>
      </Card>
    </div>
  );
}