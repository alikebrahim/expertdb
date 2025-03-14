const Footer = () => {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="bg-primary text-white py-6 mt-auto">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex flex-col md:flex-row justify-between items-center">
          <div className="mb-4 md:mb-0">
            <img 
              src="/BQA - Vertical Logo - With descriptor.svg" 
              alt="BQA Logo" 
              className="h-16"
            />
          </div>
          
          <div className="text-center md:text-right">
            <p className="text-sm">
              &copy; {currentYear} Education & Training Quality Authority
            </p>
            <p className="text-xs mt-1 text-gray-300">
              All rights reserved
            </p>
          </div>
        </div>
      </div>
    </footer>
  );
};

export default Footer;