/**
 * Email validation
 */
export const isValidEmail = (email: string): boolean => {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  return emailRegex.test(email);
};

/**
 * Password validation - checks if password meets the specified criteria:
 * - Minimum length (default: 8 characters)
 * - Requires uppercase (default: true)
 * - Requires lowercase (default: true)
 * - Requires numbers (default: true)
 * - Requires symbols (default: false)
 */
export const validatePassword = (
  password: string,
  options: {
    minLength?: number;
    requireUppercase?: boolean;
    requireLowercase?: boolean;
    requireNumbers?: boolean;
    requireSymbols?: boolean;
  } = {}
): { isValid: boolean; message: string } => {
  const {
    minLength = 8,
    requireUppercase = true,
    requireLowercase = true,
    requireNumbers = true,
    requireSymbols = false,
  } = options;

  if (!password) {
    return { isValid: false, message: 'Password is required' };
  }

  if (password.length < minLength) {
    return {
      isValid: false,
      message: `Password must be at least ${minLength} characters long`,
    };
  }

  if (requireUppercase && !/[A-Z]/.test(password)) {
    return {
      isValid: false,
      message: 'Password must contain at least one uppercase letter',
    };
  }

  if (requireLowercase && !/[a-z]/.test(password)) {
    return {
      isValid: false,
      message: 'Password must contain at least one lowercase letter',
    };
  }

  if (requireNumbers && !/\d/.test(password)) {
    return {
      isValid: false,
      message: 'Password must contain at least one number',
    };
  }

  if (requireSymbols && !/[!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]/.test(password)) {
    return {
      isValid: false,
      message: 'Password must contain at least one special character',
    };
  }

  return { isValid: true, message: 'Password is valid' };
};

/**
 * Validate required fields in a form
 * Returns an object with any missing required fields
 */
export const validateRequired = <T extends Record<string, any>>(
  data: T,
  requiredFields: Array<keyof T>
): Record<string, string> => {
  const errors: Record<string, string> = {};

  for (const field of requiredFields) {
    const value = data[field];
    
    if (value === undefined || value === null || value === '') {
      errors[field as string] = `${String(field)} is required`;
    }
  }

  return errors;
};

/**
 * Validates if a string is a valid phone number
 * Accepts various formats: +973-1234-5678, (973) 1234 5678, 973 12345678, etc.
 */
export const isValidPhoneNumber = (phone: string): boolean => {
  // Remove all non-digit characters except + (for country code)
  const digits = phone.replace(/[^\d+]/g, '');
  // Check if it's a reasonable length for a phone number (7-15 digits)
  return digits.length >= 7 && digits.length <= 15;
};

/**
 * Validates if a date string is in the future
 */
export const isFutureDate = (dateString: string): boolean => {
  const date = new Date(dateString);
  const now = new Date();
  
  // Clear time portion for date-only comparison
  now.setHours(0, 0, 0, 0);
  
  return date > now;
};

/**
 * Validates if a date string is in the past
 */
export const isPastDate = (dateString: string): boolean => {
  const date = new Date(dateString);
  const now = new Date();
  
  // Clear time portion for date-only comparison
  now.setHours(0, 0, 0, 0);
  
  return date < now;
};