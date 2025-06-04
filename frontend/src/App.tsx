import { BrowserRouter, Route, Routes, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import { UIProvider } from './contexts/UIContext';
import LoginPage from './pages/LoginPage';
import SearchPage from './pages/SearchPage';
import ExpertRequestPage from './pages/ExpertRequestPage';
import ExpertApplicationPage from './pages/ExpertApplicationPage';
import AdminPage from './pages/AdminPage';
import AdminPanelPage from './pages/AdminPanelPage';
import StatsPage from './pages/StatsPage';
import ExpertDetailPage from './pages/ExpertDetailPage';
import ExpertManagementPage from './pages/ExpertManagementPage';
import EngagementManagementPage from './pages/EngagementManagementPage';
import ProtectedRoute from './components/ProtectedRoute';

function App() {
  return (
    <AuthProvider>
      <UIProvider>
        <BrowserRouter>
          <Routes>
            {/* Public route */}
            <Route path="/" element={<LoginPage />} />
            
            {/* Protected routes */}
            <Route 
              path="/search" 
              element={
                <ProtectedRoute>
                  <SearchPage />
                </ProtectedRoute>
              } 
            />
            
            <Route 
              path="/requests" 
              element={
                <ProtectedRoute>
                  <ExpertRequestPage />
                </ProtectedRoute>
              } 
            />
            
            <Route 
              path="/apply" 
              element={
                <ProtectedRoute>
                  <ExpertApplicationPage />
                </ProtectedRoute>
              } 
            />
            
            <Route 
              path="/stats" 
              element={
                <ProtectedRoute>
                  <StatsPage />
                </ProtectedRoute>
              } 
            />
            
            <Route 
              path="/experts/:id" 
              element={
                <ProtectedRoute>
                  <ExpertDetailPage />
                </ProtectedRoute>
              } 
            />
            
            {/* Admin-only routes */}
            <Route 
              path="/admin" 
              element={
                <ProtectedRoute requiredRole="admin">
                  <AdminPage />
                </ProtectedRoute>
              } 
            />
            
            <Route 
              path="/admin/requests" 
              element={
                <ProtectedRoute requiredRole="admin">
                  <AdminPanelPage />
                </ProtectedRoute>
              } 
            />
            
            <Route 
              path="/experts/manage" 
              element={
                <ProtectedRoute requiredRole="admin">
                  <ExpertManagementPage />
                </ProtectedRoute>
              } 
            />
            
            <Route 
              path="/engagements" 
              element={
                <ProtectedRoute requiredRole="admin">
                  <EngagementManagementPage />
                </ProtectedRoute>
              } 
            />
            
            {/* Fallback route */}
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </BrowserRouter>
      </UIProvider>
    </AuthProvider>
  );
}

export default App;