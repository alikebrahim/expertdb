import { z } from 'zod';

// Filter schemas
export const expertFilterSchema = z.object({
  name: z.string().optional(),
  role: z.string().optional(),
  type: z.string().optional(),
  affiliation: z.string().optional(),
  expertArea: z.string().optional(),
  nationality: z.string().optional(),
  isAvailable: z.boolean().optional(),
  rating: z.string().optional(),
  isBahraini: z.boolean().optional(),
});

// User form schemas
export const loginSchema = z.object({
  email: z.string().email('Please enter a valid email address'),
  password: z.string().min(6, 'Password must be at least 6 characters'),
});

export const userSchema = z.object({
  id: z.number().optional(),
  name: z.string().min(2, 'Name must be at least 2 characters'),
  email: z.string().email('Please enter a valid email address'),
  role: z.enum(['admin', 'user', 'manager']),
  password: z.string().min(6, 'Password must be at least 6 characters').optional(),
  confirmPassword: z.string().optional(),
}).refine(data => data.password === data.confirmPassword, {
  message: "Passwords don't match",
  path: ['confirmPassword'],
});

// Expert form schemas
export const expertSchema = z.object({
  id: z.number().optional(),
  name: z.string().min(2, 'Name is required and must be at least 2 characters'),
  affiliation: z.string().min(2, 'Affiliation is required'),
  primaryContact: z.string().min(5, 'Primary contact is required'),
  contactType: z.enum(['email', 'phone', 'linkedin']),
  skills: z.string().min(3, 'Skills are required'),
  role: z.string().min(2, 'Role is required'),
  employmentType: z.enum(['full-time', 'part-time', 'consultant', 'retired', 'other']),
  generalArea: z.string().min(1, 'General area is required'),
  biography: z.string().optional(),
  isBahraini: z.boolean().default(false),
  availability: z.enum(['Available', 'Limited', 'Unavailable']),
  cvFile: z.any().optional(),
});

// Expert request form schemas
export const expertRequestSchema = z.object({
  title: z.string().min(5, 'Title must be at least 5 characters'),
  description: z.string().min(10, 'Description must be at least 10 characters'),
  requiredExpertise: z.array(z.number()).min(1, 'Please select at least one area of expertise'),
  deadline: z.string().min(1, 'Please select a deadline'),
  priority: z.enum(['low', 'medium', 'high']),
});

// Engagement form schemas
export const engagementSchema = z.object({
  id: z.number().optional(),
  expertId: z.number(),
  requestId: z.number(),
  status: z.enum(['pending', 'active', 'completed', 'cancelled']).default('pending'),
  startDate: z.string().min(1, 'Please select a start date'),
  endDate: z.string().min(1, 'Please select an end date').optional(),
  rate: z.number().min(0, 'Rate must be a positive number'),
  notes: z.string().max(500, 'Notes must not exceed 500 characters').optional(),
});

// Document upload schema
export const documentUploadSchema = z.object({
  title: z.string().min(2, 'Title must be at least 2 characters'),
  file: z.any().refine(file => file instanceof File, 'Please select a file'),
  type: z.enum(['cv', 'contract', 'report', 'other']),
  expertId: z.number().optional(),
  requestId: z.number().optional(),
  engagementId: z.number().optional(),
});

// Area schema
export const areaSchema = z.object({
  id: z.number().optional(),
  name: z.string().min(2, 'Name must be at least 2 characters'),
  description: z.string().min(5, 'Description must be at least 5 characters').max(200, 'Description must not exceed 200 characters').optional(),
});

// Phase schema
export const phaseSchema = z.object({
  id: z.number().optional(),
  name: z.string().min(2, 'Name must be at least 2 characters'),
  description: z.string().min(5, 'Description must be at least 5 characters').max(200, 'Description must not exceed 200 characters').optional(),
  engagementId: z.number(),
  startDate: z.string().min(1, 'Please select a start date'),
  endDate: z.string().min(1, 'Please select an end date').optional(),
  status: z.enum(['pending', 'active', 'completed', 'cancelled']).default('pending'),
});