import { ReactNode } from 'react';
import Header from './Header';
import Sidebar from './Sidebar';
import Footer from './Footer';
import NotificationContainer from '../ui/NotificationContainer';
import { Breadcrumb } from '../ui';
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
        <div className="flex-1 flex flex-col">
          <Breadcrumb />
          <main className="flex-1 p-6 bg-neutral-50">
            {children}
          </main>
        </div>
      </div>
      <Footer />
      <NotificationContainer />
    </div>
  );
};

export default Layout;