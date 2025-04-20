import React from 'react';
import { useUI } from '../../hooks/useUI';
import { ToastContainer } from './Toast';

export const NotificationContainer: React.FC = () => {
  const { notifications, dismissNotification } = useUI();

  return (
    <ToastContainer
      notifications={notifications}
      onDismiss={dismissNotification}
      position="top-right"
    />
  );
};

export default NotificationContainer;