'use client';

import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { Card, CardContent, CardFooter } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { ExpertRequest, expertAPI } from '@/lib/api';

// Form validation schema
const requestFormSchema = z.object({
  name: z.string().min(2, 'Name must be at least 2 characters'),
  designation: z.string().min(2, 'Designation is required'),
  institution: z.string().min(2, 'Institution is required'),
  email: z.string().email('Valid email is required'),
  phone: z.string().optional(),
  isBahraini: z.boolean(),
  generalArea: z.string().optional(),
  specializedArea: z.string().optional(),
  employmentType: z.string().optional(),
  role: z.string().optional(),
  isAvailable: z.boolean().default(true),
});

type RequestFormValues = z.infer<typeof requestFormSchema>;

export function RequestForm() {
  const router = useRouter();
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);
  const [formSuccess, setFormSuccess] = useState<boolean>(false);

  const {
    register,
    handleSubmit,
    setValue,
    formState: { errors },
    reset
  } = useForm<RequestFormValues>({
    resolver: zodResolver(requestFormSchema),
    defaultValues: {
      name: '',
      designation: '',
      institution: '',
      email: '',
      phone: '',
      isBahraini: false,
      generalArea: '',
      specializedArea: '',
      employmentType: '',
      role: '',
      isAvailable: true,
    },
  });

  const onSubmit = async (data: RequestFormValues) => {
    try {
      setIsSubmitting(true);
      setFormError(null);

      // Convert form data to the expected API format
      const requestData: ExpertRequest = {
        name: data.name,
        designation: data.designation,
        institution: data.institution,
        email: data.email,
        phone: data.phone || '',
        isBahraini: data.isBahraini,
        generalArea: data.generalArea,
        specializedArea: data.specializedArea,
        employmentType: data.employmentType,
        role: data.role,
        isAvailable: data.isAvailable,
        status: 'pending', // Always set status to pending for new requests
      };

      // Submit the request
      await expertAPI.createExpertRequest(requestData);
      
      // Show success state
      setFormSuccess(true);
      reset();
      
      // Redirect to success page after a short delay
      setTimeout(() => {
        router.push('/request/success');
      }, 2000);
    } catch (error: any) {
      setFormError(error.response?.data?.error || 'An error occurred during submission. Please try again.');
      console.error('Request submission error:', error);
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <Card>
      <form onSubmit={handleSubmit(onSubmit)}>
        <CardContent className="space-y-6 pt-6">
          <div className="space-y-2">
            <Label htmlFor="name">Full Name <span className="text-red-500">*</span></Label>
            <Input
              id="name"
              placeholder="Enter the expert's full name"
              {...register('name')}
            />
            {errors.name && (
              <p className="text-red-500 text-sm">{errors.name.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="designation">Designation <span className="text-red-500">*</span></Label>
            <Input
              id="designation"
              placeholder="e.g., Professor, Director, PhD, etc."
              {...register('designation')}
            />
            {errors.designation && (
              <p className="text-red-500 text-sm">{errors.designation.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="institution">Institution <span className="text-red-500">*</span></Label>
            <Input
              id="institution"
              placeholder="University, organization, or company"
              {...register('institution')}
            />
            {errors.institution && (
              <p className="text-red-500 text-sm">{errors.institution.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="email">Email <span className="text-red-500">*</span></Label>
            <Input
              id="email"
              type="email"
              placeholder="expert@example.com"
              {...register('email')}
            />
            {errors.email && (
              <p className="text-red-500 text-sm">{errors.email.message}</p>
            )}
          </div>

          <div className="space-y-2">
            <Label htmlFor="phone">Phone Number</Label>
            <Input
              id="phone"
              placeholder="+973 1234 5678"
              {...register('phone')}
            />
            {errors.phone && (
              <p className="text-red-500 text-sm">{errors.phone.message}</p>
            )}
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-2">
              <Label htmlFor="nationality">Nationality</Label>
              <Select
                onValueChange={(value) => setValue('isBahraini', value === 'bahraini')}
                defaultValue="non-bahraini"
              >
                <SelectTrigger id="nationality">
                  <SelectValue placeholder="Select nationality" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="bahraini">Bahraini</SelectItem>
                  <SelectItem value="non-bahraini">Non-Bahraini</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="availability">Availability</Label>
              <Select 
                onValueChange={(value) => setValue('isAvailable', value === 'available')}
                defaultValue="available"
              >
                <SelectTrigger id="availability">
                  <SelectValue placeholder="Select availability" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="available">Available</SelectItem>
                  <SelectItem value="unavailable">Unavailable</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="generalArea">General Area</Label>
            <Input
              id="generalArea"
              placeholder="e.g., Computer Science, Medicine, Finance"
              {...register('generalArea')}
            />
          </div>

          <div className="space-y-2">
            <Label htmlFor="specializedArea">Specialized Area</Label>
            <Input
              id="specializedArea"
              placeholder="e.g., Machine Learning, Cardiology, Investment Banking"
              {...register('specializedArea')}
            />
          </div>

          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div className="space-y-2">
              <Label htmlFor="employmentType">Employment Type</Label>
              <Select onValueChange={(value) => setValue('employmentType', value)}>
                <SelectTrigger id="employmentType">
                  <SelectValue placeholder="Select type" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Academic">Academic</SelectItem>
                  <SelectItem value="Employer">Employer</SelectItem>
                  <SelectItem value="Both">Both</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div className="space-y-2">
              <Label htmlFor="role">Role</Label>
              <Select onValueChange={(value) => setValue('role', value)}>
                <SelectTrigger id="role">
                  <SelectValue placeholder="Select role" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="Validator">Validator</SelectItem>
                  <SelectItem value="Evaluator">Evaluator</SelectItem>
                  <SelectItem value="Both">Both</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          {formError && (
            <div className="bg-destructive/10 border border-destructive text-destructive text-sm p-3 rounded-md">
              {formError}
            </div>
          )}

          {formSuccess && (
            <div className="bg-green-50 border border-green-200 text-green-700 text-sm p-3 rounded-md">
              Expert request submitted successfully! Redirecting...
            </div>
          )}
        </CardContent>
        
        <CardFooter className="flex justify-between border-t p-6">
          <Button 
            variant="outline"
            onClick={() => reset()}
            type="button"
            disabled={isSubmitting}
          >
            Cancel
          </Button>
          <Button type="submit" disabled={isSubmitting}>
            {isSubmitting ? 'Submitting...' : 'Submit Request'}
          </Button>
        </CardFooter>
      </form>
    </Card>
  );
}