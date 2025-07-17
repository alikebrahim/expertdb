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
    <header className="bg-secondary shadow-lg border-b-2 border-secondary-dark">
      <div className="container py-4 flex justify-between items-center">
        <div className="flex items-center">
          <button 
            onClick={toggleSidebar}
            className="mr-4 text-white hover:text-neutral-200 focus:outline-none focus:ring-2 focus:ring-white focus:ring-opacity-50 rounded p-1 transition-colors"
            aria-label="Toggle sidebar"
          >
            <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 12h16M4 18h16" />
            </svg>
          </button>
          <img 
            src="/BQA - Horizontal Logo.svg" 
            alt="BQA Logo" 
            className="h-10 filter brightness-0 invert"
          />
        </div>
        
        {user && (
          <div className="flex items-center space-x-4">
            <span className="text-white font-medium text-sm">
              {user.name} <span className="text-neutral-200">({user.role})</span>
            </span>
            <button 
              onClick={handleLogout}
              className="bg-white text-secondary hover:bg-neutral-100 hover:text-secondary-dark font-medium py-1 px-3 rounded transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-white focus:ring-opacity-50 text-sm"
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