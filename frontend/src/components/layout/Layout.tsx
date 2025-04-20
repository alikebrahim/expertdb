import { ReactNode } from 'react';
import Header from './Header';
import Sidebar from './Sidebar';
import Footer from './Footer';
import NotificationContainer from '../ui/NotificationContainer';
import { useUI } from '../../hooks/useUI';

interface LayoutProps {
  children: ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { isSidebarOpen } = useUI();
  
  return (
    <div className="flex flex-col min-h-screen">
      <Header />
      <div className="flex flex-1">
        <Sidebar />
        <main className={`flex-1 p-6 transition-all duration-300 ${isSidebarOpen ? 'ml-64' : 'ml-16'}`}>
          {children}
        </main>
      </div>
      <Footer />
      <NotificationContainer />
    </div>
  );
};

export default Layout;