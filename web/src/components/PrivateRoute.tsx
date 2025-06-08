import React from 'react';
import { Navigate } from 'react-router-dom';
import authService from '../services/auth';

interface PrivateRouteProps {
  children: React.ReactNode;
}

const PrivateRoute: React.FC<PrivateRouteProps> = ({ children }) => {
  const isAuthenticated = authService.isAuthenticated();
  
  if (!isAuthenticated) {
    // Redirect to login if not authenticated
    return <Navigate to="/login" />;
  }
  
  return <>{children}</>;
};

export default PrivateRoute;