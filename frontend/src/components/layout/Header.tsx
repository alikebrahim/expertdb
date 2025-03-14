import { useNavigate } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

const Header = () => {
  const { user, logout } = useAuth();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    navigate('/');
  };

  return (
    <header className="bg-white shadow-md">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4 flex justify-between items-center">
        <div className="flex items-center">
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