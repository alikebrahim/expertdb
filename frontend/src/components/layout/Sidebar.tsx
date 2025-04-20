import { NavLink } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';
import { useUI } from '../../hooks/useUI';

interface NavItem {
  to: string;
  label: string;
  roles: string[];
  icon: JSX.Element;
}

const Sidebar = () => {
  const { user } = useAuth();
  const { isSidebarOpen } = useUI();
  
  const navItems: NavItem[] = [
    { 
      to: '/search', 
      label: 'Expert Search', 
      roles: ['user', 'admin'],
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fillRule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clipRule="evenodd" />
        </svg>
      )
    },
    { 
      to: '/requests', 
      label: 'Expert Requests', 
      roles: ['user', 'admin'],
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fillRule="evenodd" d="M6 2a1 1 0 00-1 1v1H4a2 2 0 00-2 2v10a2 2 0 002 2h12a2 2 0 002-2V6a2 2 0 00-2-2h-1V3a1 1 0 10-2 0v1H7V3a1 1 0 00-1-1zm0 5a1 1 0 000 2h8a1 1 0 100-2H6z" clipRule="evenodd" />
        </svg>
      )
    },
    { 
      to: '/stats', 
      label: 'Statistics', 
      roles: ['user', 'admin'],
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path d="M2 11a1 1 0 011-1h2a1 1 0 011 1v5a1 1 0 01-1 1H3a1 1 0 01-1-1v-5zM8 7a1 1 0 011-1h2a1 1 0 011 1v9a1 1 0 01-1 1H9a1 1 0 01-1-1V7zM14 4a1 1 0 011-1h2a1 1 0 011 1v12a1 1 0 01-1 1h-2a1 1 0 01-1-1V4z" />
        </svg>
      )
    },
    { 
      to: '/experts/manage', 
      label: 'Expert Management', 
      roles: ['admin'],
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path d="M13 6a3 3 0 11-6 0 3 3 0 016 0zM18 8a2 2 0 11-4 0 2 2 0 014 0zM14 15a4 4 0 00-8 0v1h8v-1zM6 8a2 2 0 11-4 0 2 2 0 014 0zM16 18v-1a5.972 5.972 0 00-.75-2.906A3.005 3.005 0 0119 18v1h-3zM4.75 12.094A5.973 5.973 0 004 15v1H1v-1a3 3 0 013.75-2.906z" />
        </svg>
      )
    },
    { 
      to: '/engagements', 
      label: 'Engagements', 
      roles: ['admin'],
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fillRule="evenodd" d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z" clipRule="evenodd" />
        </svg>
      )
    },
    { 
      to: '/admin', 
      label: 'Admin Panel', 
      roles: ['admin'],
      icon: (
        <svg xmlns="http://www.w3.org/2000/svg" className="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fillRule="evenodd" d="M11.49 3.17c-.38-1.56-2.6-1.56-2.98 0a1.532 1.532 0 01-2.286.948c-1.372-.836-2.942.734-2.106 2.106.54.886.061 2.042-.947 2.287-1.561.379-1.561 2.6 0 2.978a1.532 1.532 0 01.947 2.287c-.836 1.372.734 2.942 2.106 2.106a1.532 1.532 0 012.287.947c.379 1.561 2.6 1.561 2.978 0a1.533 1.533 0 012.287-.947c1.372.836 2.942-.734 2.106-2.106a1.533 1.533 0 01.947-2.287c1.561-.379 1.561-2.6 0-2.978a1.532 1.532 0 01-.947-2.287c.836-1.372-.734-2.942-2.106-2.106a1.532 1.532 0 01-2.287-.947zM10 13a3 3 0 100-6 3 3 0 000 6z" clipRule="evenodd" />
        </svg>
      )
    },
  ];

  const filteredNavItems = navItems.filter(
    item => user && item.roles.includes(user.role)
  );

  if (!isSidebarOpen) {
    return (
      <aside className="bg-neutral-100 w-16 min-h-screen p-4 transition-all duration-300">
        <div className="mb-6 flex justify-center">
          <img 
            src="/Icon Logo - Color.svg" 
            alt="BQA Icon" 
            className="h-8"
          />
        </div>
        
        <nav>
          <ul className="space-y-6 flex flex-col items-center">
            {filteredNavItems.map((item) => (
              <li key={item.to} className="w-full">
                <NavLink
                  to={item.to}
                  className={({ isActive }) => 
                    `flex justify-center p-2 rounded transition-colors ${
                      isActive
                        ? 'bg-primary text-white'
                        : 'text-primary hover:bg-primary hover:bg-opacity-10'
                    }`
                  }
                  title={item.label}
                >
                  {item.icon}
                </NavLink>
              </li>
            ))}
          </ul>
        </nav>
      </aside>
    );
  }

  return (
    <aside className="bg-neutral-100 w-64 min-h-screen p-4 transition-all duration-300">
      <div className="mb-6">
        <img 
          src="/Icon Logo - Color.svg" 
          alt="BQA Icon" 
          className="h-8 mx-auto"
        />
        <h2 className="text-xl font-semibold text-center mt-2 text-primary">
          ExpertDB
        </h2>
      </div>
      
      <nav>
        <ul className="space-y-2">
          {filteredNavItems.map((item) => (
            <li key={item.to}>
              <NavLink
                to={item.to}
                className={({ isActive }) => 
                  `flex items-center px-4 py-2 rounded transition-colors ${
                    isActive
                      ? 'bg-primary text-white'
                      : 'text-primary hover:bg-primary hover:bg-opacity-10'
                  }`
                }
              >
                <span className="mr-3">{item.icon}</span>
                {item.label}
              </NavLink>
            </li>
          ))}
        </ul>
      </nav>
    </aside>
  );
};

export default Sidebar;