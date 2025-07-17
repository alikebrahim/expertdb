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

// Biography schemas for structured expert profiles
export const experienceEntrySchema = z.object({
  start_date: z.string().min(1, 'Start date is required'),
  end_date: z.string().min(1, 'End date is required'),
  title: z.string().min(1, 'Title is required'),
  organization: z.string().min(1, 'Organization is required'),
  description: z.string().min(1, 'Description is required'),
});

export const educationEntrySchema = z.object({
  start_date: z.string().min(1, 'Start date is required'),
  end_date: z.string().min(1, 'End date is required'),
  title: z.string().min(1, 'Title is required'),
  institution: z.string().min(1, 'Institution is required'),
});

export const biographySchema = z.object({
  experience: z.array(experienceEntrySchema).min(0, 'Experience entries'),
  education: z.array(educationEntrySchema).min(0, 'Education entries'),
}).refine(data => data.experience.length > 0 || data.education.length > 0, {
  message: 'At least one experience or education entry is required',
});

// Expert request form schemas - based on backend API requirements
export const expertRequestSchema = z.object({
  // Personal Information (required)
  name: z.string().min(2, 'Name is required and must be at least 2 characters'),
  designation: z.enum(['Prof.', 'Dr.', 'Mr.', 'Ms.', 'Mrs.', 'Miss', 'Eng.'], {
    errorMap: () => ({ message: 'Please select a valid designation' })
  }),
  affiliation: z.string().min(2, 'Affiliation is required'),
  phone: z.string().min(8, 'Valid phone number is required'),
  email: z.string().email('Valid email address is required'),
  
  // Professional Details (required)
  isBahraini: z.boolean(),
  isAvailable: z.boolean(),
  rating: z.number().int().min(1, 'Rating must be at least 1').max(5, 'Rating must not exceed 5'),
  role: z.enum(['evaluator', 'validator', 'evaluator/validator'], {
    errorMap: () => ({ message: 'Role must be evaluator, validator, or evaluator/validator' })
  }),
  employmentType: z.enum(['academic', 'employer'], {
    errorMap: () => ({ message: 'Employment type must be academic or employer' })
  }),
  isTrained: z.boolean(),
  isPublished: z.boolean().optional().default(false),
  
  // Expertise Areas (required)
  generalArea: z.number().min(1, 'General area is required'),
  specializedArea: z.string().min(2, 'Specialized area is required'),
  skills: z.array(z.string().min(1, 'Skill cannot be empty')).min(1, 'At least one skill is required'),
  
  // Biography & Documents (required)
  biography: biographySchema,
  cv: z.any().refine(file => file instanceof File, 'CV document is required')
    .refine(file => file?.type === 'application/pdf', 'CV must be a PDF file')
    .refine(file => file?.size <= 5 * 1024 * 1024, 'CV file size must be less than 5MB'),
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