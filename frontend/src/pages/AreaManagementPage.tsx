import React, { useState, useEffect, useCallback } from 'react';
import { ExpertArea } from '../types';
import { getExpertAreas, createExpertArea, updateExpertArea } from '../api/areas';
import AreaForm from '../components/forms/AreaForm';

const AreaManagementPage: React.FC = () => {
  const [areas, setAreas] = useState<ExpertArea[]>([]);
  const [selectedArea, setSelectedArea] = useState<ExpertArea | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const fetchAreas = useCallback(async () => {
    setIsLoading(true);
    try {
      const response = await getExpertAreas();
      setAreas(response.data);
      setError(null);
    } catch (_err) {
      setError('Failed to fetch areas');
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchAreas();
  }, [fetchAreas]);

  const handleFormSubmit = async (name: string) => {
    try {
      if (selectedArea) {
        await updateExpertArea(selectedArea.id, { name });
      } else {
        await createExpertArea({ name });
      }
      fetchAreas();
      setSelectedArea(null);
    } catch (_err) {
      setError('Failed to save area');
    }
  };

  const handleEdit = (area: ExpertArea) => {
    setSelectedArea(area);
  };

  const handleCancel = () => {
    setSelectedArea(null);
  };

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>{error}</div>;
  }

  return (
    <div className="container mx-auto p-4">
      <h1 className="text-2xl font-bold mb-4">Specialization Area Management</h1>
      <div className="grid grid-cols-1 md:grid-cols-2 gap-8">
        <div>
          <h2 className="text-xl font-bold mb-2">Existing Areas</h2>
          <ul>
            {areas.map((area) => (
              <li key={area.id} className="flex justify-between items-center p-2 border-b">
                <span>{area.name}</span>
                <button onClick={() => handleEdit(area)} className="text-blue-500">Edit</button>
              </li>
            ))}
          </ul>
        </div>
        <div>
          <h2 className="text-xl font-bold mb-2">{selectedArea ? 'Edit Area' : 'Create Area'}</h2>
          <AreaForm 
            area={selectedArea} 
            onSubmit={handleFormSubmit} 
            onCancel={handleCancel} 
          />
        </div>
      </div>
    </div>
  );
};

export default AreaManagementPage;