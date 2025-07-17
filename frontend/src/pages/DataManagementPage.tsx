
import React, { useState } from 'react';
import { backupApi } from '../services/api';

const DataManagementPage: React.FC = () => {
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const handleBackup = async () => {
    setIsLoading(true);
    setError(null);
    setSuccess(null);

    try {
      const response = await backupApi.generateBackup();
      const url = window.URL.createObjectURL(new Blob([response.data]));
      const link = document.createElement('a');
      link.href = url;
      link.setAttribute('download', 'expertdb_backup.zip');
      document.body.appendChild(link);
      link.click();
      link.remove();
      setSuccess('Backup downloaded successfully!');
    } catch (_err) {
      setError('Failed to download backup');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Data Management</h1>
      <div className="bg-white p-6 rounded-lg shadow">
        <h2 className="text-xl font-bold mb-2">Database Backup</h2>
        <p className="mb-4">Click the button below to download a ZIP file containing CSV exports of the database tables.</p>
        {error && <p className="text-red-500">{error}</p>}
        {success && <p className="text-green-500">{success}</p>}
        <button
          onClick={handleBackup}
          disabled={isLoading}
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded disabled:bg-gray-400"
        >
          {isLoading ? 'Generating Backup...' : 'Download Backup'}
        </button>
      </div>
    </div>
  );
};

export default DataManagementPage;
