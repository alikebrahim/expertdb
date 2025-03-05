'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useParams, useRouter } from 'next/navigation';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle, CardFooter } from '@/components/ui/card';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle, DialogTrigger, DialogClose } from '@/components/ui/dialog';
import { Label } from '@/components/ui/label';
import { Input } from '@/components/ui/input';
import { Textarea } from '@/components/ui/textarea';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { Expert, Engagement, expertAPI } from '@/lib/api';

export default function ExpertDetailPage() {
  const params = useParams();
  const router = useRouter();
  const [expert, setExpert] = useState<Expert | null>(null);
  const [engagements, setEngagements] = useState<Engagement[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [isLoadingEngagements, setIsLoadingEngagements] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [isAddEngagementOpen, setIsAddEngagementOpen] = useState(false);
  
  // New engagement form state
  const [engagementForm, setEngagementForm] = useState({
    engagementType: '',
    startDate: '',
    endDate: '',
    projectName: '',
    status: 'pending',
    feedbackScore: 0,
    notes: '',
  });
  
  useEffect(() => {
    const fetchExpert = async () => {
      try {
        setIsLoading(true);
        setError(null);
        
        if (!params.id || Array.isArray(params.id)) {
          throw new Error('Invalid expert ID');
        }
        
        const expertId = parseInt(params.id);
        const expertData = await expertAPI.getExpertById(expertId);
        setExpert(expertData);
        
        // Fetch engagements after expert data is loaded
        setIsLoadingEngagements(true);
        const engagementData = await expertAPI.getExpertEngagements(expertId);
        setEngagements(engagementData);
      } catch (error: any) {
        console.error('Error fetching expert:', error);
        setError(error.response?.data?.error || 'Failed to load expert details');
      } finally {
        setIsLoading(false);
        setIsLoadingEngagements(false);
      }
    };
    
    fetchExpert();
  }, [params.id]);
  
  // Handle engagement form input changes
  const handleEngagementFormChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement> | { target: { name: string; value: string | number } }) => {
    const { name, value } = e.target;
    setEngagementForm(prev => ({
      ...prev,
      [name]: value
    }));
  };
  
  // Handle select input changes
  const handleSelectChange = (name: string, value: string) => {
    setEngagementForm(prev => ({
      ...prev,
      [name]: value
    }));
  };
  
  // Handle form submission for new engagement
  const handleSubmitEngagement = async (e: React.FormEvent) => {
    e.preventDefault();
    
    if (!expert) return;
    
    try {
      const newEngagement = {
        ...engagementForm,
        expertId: expert.id,
        feedbackScore: engagementForm.feedbackScore || undefined
      };
      
      await expertAPI.createEngagement(newEngagement);
      
      // Refresh engagements
      const refreshedEngagements = await expertAPI.getExpertEngagements(expert.id);
      setEngagements(refreshedEngagements);
      
      // Reset form and close dialog
      setEngagementForm({
        engagementType: '',
        startDate: '',
        endDate: '',
        projectName: '',
        status: 'pending',
        feedbackScore: 0,
        notes: '',
      });
      setIsAddEngagementOpen(false);
      
    } catch (error) {
      console.error('Error creating engagement:', error);
    }
  };
  
  // Function to calculate average rating from engagements
  const calculateAverageRating = (engagements: Engagement[]) => {
    const ratingsWithFeedback = engagements.filter(e => e.feedbackScore && e.feedbackScore > 0);
    if (ratingsWithFeedback.length === 0) return 'N/A';
    
    const sum = ratingsWithFeedback.reduce((acc, curr) => acc + (curr.feedbackScore || 0), 0);
    const average = (sum / ratingsWithFeedback.length).toFixed(1);
    return average;
  };
  
  if (isLoading) {
    return (
      <div className="container py-10">
        <div className="flex justify-center py-20">
          <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary"></div>
        </div>
      </div>
    );
  }
  
  if (error || !expert) {
    return (
      <div className="container py-10">
        <Card>
          <CardContent className="pt-6 text-center py-12">
            <h2 className="text-xl font-semibold mb-2">Expert Not Found</h2>
            <p className="text-muted-foreground mb-6">
              {error || 'The expert you are looking for could not be found.'}
            </p>
            <Button asChild>
              <Link href="/search">Back to Search</Link>
            </Button>
          </CardContent>
        </Card>
      </div>
    );
  }
  
  const averageRating = calculateAverageRating(engagements);

  return (
    <div className="container py-10">
      <div className="mb-6 flex justify-between items-center">
        <Button variant="outline" onClick={() => router.back()}>
          ← Back
        </Button>
        
        <Dialog open={isAddEngagementOpen} onOpenChange={setIsAddEngagementOpen}>
          <DialogTrigger asChild>
            <Button>Add Engagement</Button>
          </DialogTrigger>
          <DialogContent className="sm:max-w-[525px]">
            <DialogHeader>
              <DialogTitle>Add Expert Engagement</DialogTitle>
              <DialogDescription>
                Record a new engagement or project for this expert. Include feedback score to contribute to their rating.
              </DialogDescription>
            </DialogHeader>
            
            <form onSubmit={handleSubmitEngagement}>
              <div className="grid gap-4 py-4">
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="engagementType" className="text-right">Type</Label>
                  <div className="col-span-3">
                    <Select 
                      onValueChange={(value) => handleSelectChange('engagementType', value)}
                      value={engagementForm.engagementType}
                    >
                      <SelectTrigger id="engagementType">
                        <SelectValue placeholder="Select type" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="evaluation">Evaluation</SelectItem>
                        <SelectItem value="validation">Validation</SelectItem>
                        <SelectItem value="panel">Panel Review</SelectItem>
                        <SelectItem value="consultation">Consultation</SelectItem>
                        <SelectItem value="other">Other</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="projectName" className="text-right">Project</Label>
                  <Input
                    id="projectName"
                    name="projectName"
                    placeholder="Project name"
                    value={engagementForm.projectName}
                    onChange={handleEngagementFormChange}
                    className="col-span-3"
                  />
                </div>
                
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="startDate" className="text-right">Start Date</Label>
                  <Input
                    id="startDate"
                    name="startDate"
                    type="date"
                    value={engagementForm.startDate}
                    onChange={handleEngagementFormChange}
                    className="col-span-3"
                    required
                  />
                </div>
                
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="endDate" className="text-right">End Date</Label>
                  <Input
                    id="endDate"
                    name="endDate"
                    type="date"
                    value={engagementForm.endDate}
                    onChange={handleEngagementFormChange}
                    className="col-span-3"
                  />
                </div>
                
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="status" className="text-right">Status</Label>
                  <div className="col-span-3">
                    <Select 
                      onValueChange={(value) => handleSelectChange('status', value)}
                      value={engagementForm.status}
                    >
                      <SelectTrigger id="status">
                        <SelectValue placeholder="Select status" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="pending">Pending</SelectItem>
                        <SelectItem value="active">Active</SelectItem>
                        <SelectItem value="completed">Completed</SelectItem>
                        <SelectItem value="cancelled">Cancelled</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="feedbackScore" className="text-right">Rating (1-5)</Label>
                  <div className="col-span-3">
                    <Select 
                      onValueChange={(value) => handleSelectChange('feedbackScore', value)}
                      value={String(engagementForm.feedbackScore)}
                    >
                      <SelectTrigger id="feedbackScore">
                        <SelectValue placeholder="Select rating" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="0">No rating</SelectItem>
                        <SelectItem value="1">1 - Poor</SelectItem>
                        <SelectItem value="2">2 - Fair</SelectItem>
                        <SelectItem value="3">3 - Good</SelectItem>
                        <SelectItem value="4">4 - Very Good</SelectItem>
                        <SelectItem value="5">5 - Excellent</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>
                
                <div className="grid grid-cols-4 items-center gap-4">
                  <Label htmlFor="notes" className="text-right">Notes</Label>
                  <Textarea
                    id="notes"
                    name="notes"
                    placeholder="Add notes about this engagement"
                    value={engagementForm.notes}
                    onChange={handleEngagementFormChange}
                    className="col-span-3"
                  />
                </div>
              </div>
              
              <div className="flex justify-end gap-3">
                <DialogClose asChild>
                  <Button type="button" variant="outline">Cancel</Button>
                </DialogClose>
                <Button type="submit">Save Engagement</Button>
              </div>
            </form>
          </DialogContent>
        </Dialog>
      </div>
      
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="md:col-span-2">
          <Card>
            <CardHeader className="flex flex-row justify-between items-center">
              <div>
                <CardTitle>{expert.name}</CardTitle>
                <p className="text-sm text-muted-foreground mt-1">
                  {expert.expertId}
                </p>
              </div>
              <div className="flex flex-col items-end">
                <div className="flex items-center gap-1">
                  <span className="text-sm font-medium">Rating:</span>
                  <span className={`text-lg font-bold ${
                    averageRating === 'N/A' 
                      ? 'text-gray-400' 
                      : parseFloat(averageRating) >= 4 
                        ? 'text-green-600' 
                        : parseFloat(averageRating) >= 3 
                          ? 'text-amber-600' 
                          : 'text-red-600'
                  }`}>
                    {averageRating} {averageRating !== 'N/A' && '/ 5.0'}
                  </span>
                </div>
                <p className="text-xs text-muted-foreground mt-1">
                  {engagements.filter(e => e.feedbackScore && e.feedbackScore > 0).length} ratings
                </p>
              </div>
            </CardHeader>
            <CardContent>
              <dl className="grid grid-cols-1 sm:grid-cols-2 gap-x-4 gap-y-6">
                <div>
                  <dt className="text-sm font-medium text-muted-foreground">Designation</dt>
                  <dd className="mt-1">{expert.designation}</dd>
                </div>
                
                <div>
                  <dt className="text-sm font-medium text-muted-foreground">Institution</dt>
                  <dd className="mt-1">{expert.institution}</dd>
                </div>
                
                <div>
                  <dt className="text-sm font-medium text-muted-foreground">Nationality</dt>
                  <dd className="mt-1">{expert.nationality || (expert.isBahraini ? 'Bahraini' : 'Non-Bahraini')}</dd>
                </div>
                
                <div>
                  <dt className="text-sm font-medium text-muted-foreground">Availability</dt>
                  <dd className="mt-1">
                    <span className={expert.isAvailable ? 'text-green-600' : 'text-red-600'}>
                      {expert.isAvailable ? 'Available' : 'Unavailable'}
                    </span>
                  </dd>
                </div>
                
                <div>
                  <dt className="text-sm font-medium text-muted-foreground">General Area</dt>
                  <dd className="mt-1">{expert.generalArea || 'Not specified'}</dd>
                </div>
                
                <div>
                  <dt className="text-sm font-medium text-muted-foreground">Specialized Area</dt>
                  <dd className="mt-1">{expert.specializedArea || 'Not specified'}</dd>
                </div>
                
                <div>
                  <dt className="text-sm font-medium text-muted-foreground">Role</dt>
                  <dd className="mt-1">{expert.role || 'Not specified'}</dd>
                </div>
                
                <div>
                  <dt className="text-sm font-medium text-muted-foreground">Employment Type</dt>
                  <dd className="mt-1">{expert.employmentType || 'Not specified'}</dd>
                </div>
                
                <div>
                  <dt className="text-sm font-medium text-muted-foreground">Contact Email</dt>
                  <dd className="mt-1">{expert.email || 'Not available'}</dd>
                </div>
                
                <div>
                  <dt className="text-sm font-medium text-muted-foreground">Contact Phone</dt>
                  <dd className="mt-1">{expert.phone || 'Not available'}</dd>
                </div>
              </dl>
            </CardContent>
          </Card>
          
          {/* Engagement History */}
          <Card className="mt-6">
            <CardHeader>
              <CardTitle>Engagement History</CardTitle>
            </CardHeader>
            <CardContent>
              {isLoadingEngagements ? (
                <div className="flex justify-center py-6">
                  <div className="animate-spin rounded-full h-8 w-8 border-t-2 border-b-2 border-primary"></div>
                </div>
              ) : engagements.length === 0 ? (
                <div className="text-center py-6">
                  <p className="text-muted-foreground">No engagements yet</p>
                </div>
              ) : (
                <div className="space-y-4">
                  {engagements.map((engagement) => (
                    <div key={engagement.id} className="border rounded-lg p-4 hover:bg-muted/50 transition">
                      <div className="flex justify-between items-start">
                        <div>
                          <h4 className="font-medium capitalize">
                            {engagement.engagementType.replace('-', ' ')}
                          </h4>
                          <p className="text-sm">{engagement.projectName}</p>
                        </div>
                        <div className="text-right">
                          <span className={`text-xs px-2 py-1 rounded-full ${
                            engagement.status === 'completed' 
                              ? 'bg-green-100 text-green-800' 
                              : engagement.status === 'active' 
                                ? 'bg-blue-100 text-blue-800' 
                                : engagement.status === 'cancelled' 
                                  ? 'bg-red-100 text-red-800' 
                                  : 'bg-amber-100 text-amber-800'
                          }`}>
                            {engagement.status}
                          </span>
                          
                          {engagement.feedbackScore && engagement.feedbackScore > 0 && (
                            <div className="mt-2">
                              <span className="text-sm font-semibold">
                                Rating: {engagement.feedbackScore}/5
                              </span>
                            </div>
                          )}
                        </div>
                      </div>
                      
                      <div className="text-sm text-muted-foreground mt-2 flex gap-2">
                        <span>
                          Started: {new Date(engagement.startDate).toLocaleDateString()}
                        </span>
                        {engagement.endDate && (
                          <>
                            <span>•</span>
                            <span>
                              Ended: {new Date(engagement.endDate).toLocaleDateString()}
                            </span>
                          </>
                        )}
                      </div>
                      
                      {engagement.notes && (
                        <p className="text-sm mt-2 border-t pt-2">{engagement.notes}</p>
                      )}
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
            <CardFooter className="flex justify-center border-t">
              <Button 
                variant="outline" 
                className="w-full" 
                onClick={() => setIsAddEngagementOpen(true)}
              >
                Add New Engagement
              </Button>
            </CardFooter>
          </Card>
        </div>
        
        <div>
          <Card>
            <CardHeader>
              <CardTitle>ISCED Classification</CardTitle>
            </CardHeader>
            <CardContent>
              {expert.iscedField ? (
                <div className="space-y-4">
                  <div>
                    <h4 className="text-sm font-medium text-muted-foreground">Field</h4>
                    <p>{expert.iscedField.broadName}</p>
                    {expert.iscedField.narrowName && (
                      <p className="text-sm text-muted-foreground mt-1">
                        {expert.iscedField.narrowName}
                      </p>
                    )}
                  </div>
                  
                  {expert.iscedLevel && (
                    <div>
                      <h4 className="text-sm font-medium text-muted-foreground">Level</h4>
                      <p>{expert.iscedLevel.name}</p>
                    </div>
                  )}
                </div>
              ) : (
                <p className="text-muted-foreground">No ISCED classification available</p>
              )}
            </CardContent>
          </Card>
          
          {expert.areas && expert.areas.length > 0 && (
            <Card className="mt-6">
              <CardHeader>
                <CardTitle>Specialization Areas</CardTitle>
              </CardHeader>
              <CardContent>
                <ul className="space-y-1">
                  {expert.areas.map((area) => (
                    <li key={area.id}>{area.name}</li>
                  ))}
                </ul>
              </CardContent>
            </Card>
          )}
          
          {/* Rating Summary Card */}
          <Card className="mt-6">
            <CardHeader>
              <CardTitle>Rating Summary</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="text-center">
                <div className="text-4xl font-bold my-2">
                  <span className={`${
                    averageRating === 'N/A' 
                      ? 'text-gray-400' 
                      : parseFloat(averageRating) >= 4 
                        ? 'text-green-600' 
                        : parseFloat(averageRating) >= 3 
                          ? 'text-amber-600' 
                          : 'text-red-600'
                  }`}>
                    {averageRating}
                  </span>
                  {averageRating !== 'N/A' && ' / 5.0'}
                </div>
                <p className="text-sm text-muted-foreground">
                  Based on {engagements.filter(e => e.feedbackScore && e.feedbackScore > 0).length} ratings
                </p>
                
                {/* Rating distribution */}
                {averageRating !== 'N/A' && (
                  <div className="mt-4 space-y-2">
                    {[5, 4, 3, 2, 1].map(rating => {
                      const count = engagements.filter(e => e.feedbackScore === rating).length;
                      const percentage = engagements.filter(e => e.feedbackScore && e.feedbackScore > 0).length > 0
                        ? (count / engagements.filter(e => e.feedbackScore && e.feedbackScore > 0).length) * 100
                        : 0;
                      
                      return (
                        <div key={rating} className="flex items-center gap-2">
                          <span className="text-sm w-8">{rating} star</span>
                          <div className="h-2 bg-gray-200 rounded-full flex-1">
                            <div 
                              className="h-2 bg-primary rounded-full" 
                              style={{ width: `${percentage}%` }}
                            ></div>
                          </div>
                          <span className="text-sm w-8">{count}</span>
                        </div>
                      );
                    })}
                  </div>
                )}
              </div>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}