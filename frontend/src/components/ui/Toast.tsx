import React, { useEffect, useState } from 'react';
import { Notification } from '../../contexts/UIContext';

interface ToastProps extends Notification {
  onDismiss: () => void;
}

export const Toast: React.FC<ToastProps> = ({ id, type, message, duration = 5000, onDismiss }) => {
  const [isVisible, setIsVisible] = useState(false);
  const [isClosing, setIsClosing] = useState(false);
  
  // Set icon and color based on type
  const getIconAndColor = () => {
    switch (type) {
      case 'success':
        return {
          icon: (
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clipRule="evenodd"></path>
            </svg>
          ),
          color: 'bg-green-100 text-green-800 dark:bg-green-800 dark:text-green-100',
          closeColor: 'text-green-600 dark:text-green-100',
          progressColor: 'bg-green-600',
        };
      case 'error':
        return {
          icon: (
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
              <path fillRule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clipRule="evenodd"></path>
            </svg>
          ),
          color: 'bg-red-100 text-red-800 dark:bg-red-800 dark:text-red-100',
          closeColor: 'text-red-600 dark:text-red-100',
          progressColor: 'bg-red-600',
        };
      case 'warning':
        return {
          icon: (
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
              <path fillRule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clipRule="evenodd"></path>
            </svg>
          ),
          color: 'bg-yellow-100 text-yellow-800 dark:bg-yellow-800 dark:text-yellow-100',
          closeColor: 'text-yellow-600 dark:text-yellow-100',
          progressColor: 'bg-yellow-600',
        };
      case 'info':
      default:
        return {
          icon: (
            <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
              <path fillRule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clipRule="evenodd"></path>
            </svg>
          ),
          color: 'bg-blue-100 text-blue-800 dark:bg-blue-800 dark:text-blue-100',
          closeColor: 'text-blue-600 dark:text-blue-100',
          progressColor: 'bg-blue-600',
        };
    }
  };

  const { icon, color, closeColor, progressColor } = getIconAndColor();

  // Handle animation and auto-dismiss
  useEffect(() => {
    // Start entrance animation
    const enterTimer = setTimeout(() => {
      setIsVisible(true);
    }, 10);
    
    // Auto-dismiss after duration
    let dismissTimer: NodeJS.Timeout;
    if (duration !== 0) {
      dismissTimer = setTimeout(() => {
        handleDismiss();
      }, duration);
    }
    
    return () => {
      clearTimeout(enterTimer);
      if (dismissTimer) clearTimeout(dismissTimer);
    };
  }, [duration]);

  // Handle dismiss with exit animation
  const handleDismiss = () => {
    setIsClosing(true);
    setTimeout(() => {
      onDismiss();
    }, 300); // Match this with the CSS transition duration
  };

  return (
    <div 
      className={`
        max-w-xs w-full ${color} shadow-lg rounded-lg pointer-events-auto 
        transform transition-all duration-300 ease-in-out
        ${isVisible ? 'translate-x-0 opacity-100' : 'translate-x-full opacity-0'}
        ${isClosing ? 'translate-x-full opacity-0' : ''}
      `}
      role="alert"
    >
      <div className="relative overflow-hidden rounded-lg">
        {/* Progress bar */}
        {duration !== 0 && (
          <div 
            className={`absolute bottom-0 left-0 h-1 ${progressColor}`}
            style={{ 
              width: '100%', 
              animation: `shrink ${duration}ms linear forwards` 
            }}
          />
        )}
        
        <div className="p-3">
          <div className="flex items-start">
            <div className="flex-shrink-0">
              {icon}
            </div>
            <div className="ml-3 w-0 flex-1 pt-0.5">
              <p className="text-sm font-medium">
                {message}
              </p>
            </div>
            <div className="ml-4 flex-shrink-0 flex">
              <button
                type="button"
                className={`bg-transparent rounded-md inline-flex ${closeColor} focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500`}
                onClick={handleDismiss}
              >
                <span className="sr-only">Close</span>
                <svg className="h-5 w-5" fill="currentColor" viewBox="0 0 20 20" xmlns="http://www.w3.org/2000/svg">
                  <path fillRule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clipRule="evenodd"></path>
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};

export const ToastContainer: React.FC<{ 
  notifications: Notification[];
  onDismiss: (id: string) => void;
  position?: 'top-right' | 'top-left' | 'bottom-right' | 'bottom-left';
}> = ({ 
  notifications, 
  onDismiss,
  position = 'top-right' 
}) => {
  // Get position classes
  const getPositionClasses = () => {
    switch (position) {
      case 'top-left':
        return 'top-0 left-0';
      case 'bottom-right':
        return 'bottom-0 right-0';
      case 'bottom-left':
        return 'bottom-0 left-0';
      case 'top-right':
      default:
        return 'top-0 right-0';
    }
  };

  const positionClasses = getPositionClasses();

  return (
    <div 
      className={`fixed ${positionClasses} z-50 p-4 space-y-4 w-full max-w-xs pointer-events-none`}
      aria-live="assertive"
    >
      {notifications.map((notification) => (
        <Toast
          key={notification.id}
          {...notification}
          onDismiss={() => onDismiss(notification.id)}
        />
      ))}
    </div>
  );
};

// Add keyframes animation for the progress bar to global CSS
// This should be added to your tailwind.css or the head of your index.html
const style = document.createElement('style');
style.textContent = `
  @keyframes shrink {
    0% { width: 100%; }
    100% { width: 0%; }
  }
`;
document.head.appendChild(style);