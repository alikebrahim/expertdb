/**
 * Date and time formatting utilities
 */

// Format date as DD/MM/YYYY
export const formatDate = (dateString: string): string => {
  if (!dateString) return '-';
  
  const date = new Date(dateString);
  if (isNaN(date.getTime())) return '-';
  
  return date.toLocaleDateString('en-GB', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
  });
};

// Format date as DD/MM/YYYY HH:MM
export const formatDateTime = (dateString: string): string => {
  if (!dateString) return '-';
  
  const date = new Date(dateString);
  if (isNaN(date.getTime())) return '-';
  
  return date.toLocaleDateString('en-GB', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
  });
};

// Format number as currency with optional currency symbol
export const formatCurrency = (
  amount: number, 
  currency: string = 'BHD', 
  locale: string = 'en-BH'
): string => {
  if (amount === undefined || amount === null) return '-';
  
  return new Intl.NumberFormat(locale, {
    style: 'currency',
    currency,
    minimumFractionDigits: 2,
  }).format(amount);
};

// Format name with titles (Dr., Prof., etc.)
export const formatName = (
  name: string, 
  title?: string, 
  abbreviate: boolean = false
): string => {
  if (!name) return '-';
  
  if (!title) return name;
  
  if (abbreviate) {
    const parts = name.split(' ');
    const firstName = parts[0];
    const lastName = parts[parts.length - 1];
    
    if (parts.length === 1) {
      return `${title} ${firstName}`;
    }
    
    return `${title} ${firstName.charAt(0)}. ${lastName}`;
  }
  
  return `${title} ${name}`;
};

/**
 * Format status badge value
 * Returns the status text and a corresponding TailwindCSS class name
 */
export const formatStatusBadge = (
  status: string
): { text: string; className: string } => {
  const normalizedStatus = status?.toLowerCase() || '';
  
  switch (normalizedStatus) {
    case 'active':
    case 'approved':
    case 'completed':
    case 'published':
      return {
        text: status,
        className: 'bg-green-100 text-green-800 border-green-500',
      };
      
    case 'pending':
    case 'in review':
    case 'processing':
    case 'in progress':
      return {
        text: status,
        className: 'bg-blue-100 text-blue-800 border-blue-500',
      };
      
    case 'rejected':
    case 'cancelled':
    case 'failed':
    case 'inactive':
      return {
        text: status,
        className: 'bg-red-100 text-red-800 border-red-500',
      };
      
    case 'draft':
    case 'waiting':
    case 'on hold':
      return {
        text: status,
        className: 'bg-yellow-100 text-yellow-800 border-yellow-500',
      };
      
    default:
      return {
        text: status || 'Unknown',
        className: 'bg-gray-100 text-gray-800 border-gray-500',
      };
  }
};