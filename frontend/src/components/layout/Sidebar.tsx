import { NavLink } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

interface NavItem {
  to: string;
  label: string;
  roles: string[];
}

const Sidebar = () => {
  const { user } = useAuth();
  
  const navItems: NavItem[] = [
    { to: '/search', label: 'Expert Search', roles: ['user', 'admin'] },
    { to: '/requests', label: 'Expert Requests', roles: ['user', 'admin'] },
    { to: '/stats', label: 'Statistics', roles: ['user', 'admin'] },
    { to: '/admin', label: 'Admin Panel', roles: ['admin'] },
  ];

  const filteredNavItems = navItems.filter(
    item => user && item.roles.includes(user.role)
  );

  return (
    <aside className="bg-neutral-100 w-64 min-h-screen p-4">
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
                  `block px-4 py-2 rounded transition-colors ${
                    isActive
                      ? 'bg-primary text-white'
                      : 'text-primary hover:bg-primary hover:bg-opacity-10'
                  }`
                }
              >
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