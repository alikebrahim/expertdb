import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { useUI } from '../../hooks/useUI';

const Header = () => {
  const { user, logout } = useAuth();
  const { toggleSidebar } = useUI();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <header className="bg-white shadow-md">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <div className="flex items-center">
          <button 
            onClick={toggleSidebar}
            className="mr-4 text-gray-600 hover:text-gray-900 focus:outline-none"
            aria-label="Toggle sidebar"
          >
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          <img 
            src="/BQA - Horizontal Logo.svg" 
            alt="BQA Logo" 
            className="h-10"
          />
        </div>
        
        {user && (
          <div className="flex items-center space-x-4">
            <span className="text-primary font-medium">
              {user.name} ({user.role})
            </span>
            <button 
              onClick={handleLogout}
              className="btn-outline text-sm py-1 px-3"
            >
              Logout
            </button>
          </div>
        )}
      </div>
    </header>
  );
};

export default Header;