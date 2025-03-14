import { BrowserRouter, Route, Routes, Navigate } from 'react-router-dom';
import { AuthProvider } from './contexts/AuthContext';
import LoginPage from './pages/LoginPage';
import SearchPage from './pages/SearchPage';
import ExpertRequestPage from './pages/ExpertRequestPage';
import AdminPage from './pages/AdminPage';
import StatsPage from './pages/StatsPage';
import ProtectedRoute from './components/ProtectedRoute';


function App() {
  return (
    <AuthProvider>
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
            path="/stats" 
            element={
              <ProtectedRoute>
                <StatsPage />
              </ProtectedRoute>
            } 
          />
          
          {/* Admin-only route */}
          <Route 
            path="/admin" 
            element={
              <ProtectedRoute requiredRole="admin">
                <AdminPage />
              </ProtectedRoute>
            } 
          />
          
          {/* Fallback route */}
          <Route path="*" element={<Navigate to="/" replace />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;
