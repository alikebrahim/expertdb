'use client';

import { useState, useEffect } from 'react';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardFooter, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Textarea } from '@/components/ui/textarea';
import { aiAPI, expertAPI } from '@/lib/api';
import Link from 'next/link';
import { Navbar } from '@/components/layout/navbar';
import RequireAuth from '@/components/auth/require-auth';

export default function AdminPanelPage() {
  const [iscedFields, setIscedFields] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  
  // Fetch ISCED fields on component mount
  useEffect(() => {
    const fetchISCEDFields = async () => {
      try {
        const fields = await expertAPI.getISCEDFields();
        setIscedFields(fields);
      } catch (error) {
        console.error('Failed to fetch ISCED fields:', error);
      }
    };
    
    fetchISCEDFields();
  }, []);
  
  return (
    <RequireAuth>
      <>
        <Navbar />
        <div className="container py-10">
          <div className="mb-6">
            <h1 className="text-3xl font-bold">Admin Panel</h1>
            <p className="text-muted-foreground mt-2">
              Manage expert database and user access.
            </p>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <Card>
              <CardHeader>
                <CardTitle>Database Summary</CardTitle>
                <CardDescription>
                  Current status of the expert database
                </CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="grid grid-cols-2 gap-4">
                  <div className="bg-muted p-4 rounded-md">
                    <h3 className="font-medium text-sm text-muted-foreground">Total Experts</h3>
                    <p className="text-2xl font-bold">482</p>
                  </div>
                  <div className="bg-muted p-4 rounded-md">
                    <h3 className="font-medium text-sm text-muted-foreground">ISCED Fields</h3>
                    <p className="text-2xl font-bold">{iscedFields.length}</p>
                  </div>
                  <div className="bg-muted p-4 rounded-md">
                    <h3 className="font-medium text-sm text-muted-foreground">Pending Requests</h3>
                    <p className="text-2xl font-bold">8</p>
                  </div>
                  <div className="bg-muted p-4 rounded-md">
                    <h3 className="font-medium text-sm text-muted-foreground">Active Users</h3>
                    <p className="text-2xl font-bold">12</p>
                  </div>
                </div>
              </CardContent>
              <CardFooter>
                <Button asChild className="w-full">
                  <Link href="/statistics">View Detailed Statistics</Link>
                </Button>
              </CardFooter>
            </Card>
            
            <Card>
              <CardHeader>
                <CardTitle>Quick Actions</CardTitle>
                <CardDescription>
                  Common administrative tasks
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <Button asChild className="w-full">
                    <Link href="/search">Search Experts</Link>
                  </Button>
                  <Button asChild variant="outline" className="w-full">
                    <Link href="/request">Add New Expert</Link>
                  </Button>
                  {error && (
                    <div className="bg-destructive/10 border border-destructive text-destructive p-4 rounded-md">
                      {error}
                    </div>
                  )}
                </div>
              </CardContent>
            </Card>
          </div>
          
          <div className="mt-6">
            <Card>
              <CardHeader>
                <CardTitle>Recent Activity</CardTitle>
                <CardDescription>
                  Latest updates to the expert database
                </CardDescription>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="border-b pb-2">
                    <div className="flex justify-between items-start">
                      <div>
                        <h3 className="font-medium">New Expert Added</h3>
                        <p className="text-sm text-muted-foreground">Dr. James Wilson - Engineering</p>
                      </div>
                      <span className="text-xs text-muted-foreground">Today, 14:32</span>
                    </div>
                  </div>
                  <div className="border-b pb-2">
                    <div className="flex justify-between items-start">
                      <div>
                        <h3 className="font-medium">Expert Request Approved</h3>
                        <p className="text-sm text-muted-foreground">Prof. Maria Garcia - Education</p>
                      </div>
                      <span className="text-xs text-muted-foreground">Yesterday, 09:15</span>
                    </div>
                  </div>
                  <div className="border-b pb-2">
                    <div className="flex justify-between items-start">
                      <div>
                        <h3 className="font-medium">User Account Created</h3>
                        <p className="text-sm text-muted-foreground">ahmed.mohammed@example.com</p>
                      </div>
                      <span className="text-xs text-muted-foreground">Mar 3, 2025</span>
                    </div>
                  </div>
                  <div className="border-b pb-2">
                    <div className="flex justify-between items-start">
                      <div>
                        <h3 className="font-medium">Expert Profile Updated</h3>
                        <p className="text-sm text-muted-foreground">Dr. Sarah Johnson - Computer Science</p>
                      </div>
                      <span className="text-xs text-muted-foreground">Mar 2, 2025</span>
                    </div>
                  </div>
                </div>
              </CardContent>
              <CardFooter>
                <Button variant="outline" className="w-full">
                  View All Activity
                </Button>
              </CardFooter>
            </Card>
          </div>
        </div>
      </>
    </RequireAuth>
  );
}