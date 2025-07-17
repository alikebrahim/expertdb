
import React, { useState, useEffect, useCallback } from 'react';
import { Phase, PhaseListResponse } from '../types';
import { getPhases } from '../api/phases';
import CreatePhaseModal from '../components/forms/CreatePhaseModal';

const PhasePlanningPage: React.FC = () => {
  const [phases, setPhases] = useState<Phase[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);

  const fetchPhases = useCallback(async () => {
    setIsLoading(true);
    try {
      const response = await getPhases();
      if (response.success && response.data) {
        setPhases(response.data.phases);
      }
      setError(null);
    } catch (_err) {
      setError('Failed to fetch phases');
    } finally {
      setIsLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchPhases();
  }, [fetchPhases]);

  const handleOpenModal = () => {
    setIsModalOpen(true);
  };

  const handleCloseModal = () => {
    setIsModalOpen(false);
  };

  const handlePhaseCreated = () => {
    handleCloseModal();
    fetchPhases();
  };

  if (isLoading) {
    return <div>Loading...</div>;
  }

  if (error) {
    return <div>{error}</div>;
  }

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-2xl font-bold">Phase Planning</h1>
        <button
          onClick={handleOpenModal}
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          Create New Phase
        </button>
      </div>
      <CreatePhaseModal 
        isOpen={isModalOpen} 
        onClose={handleCloseModal} 
        onPhaseCreated={handlePhaseCreated} 
      />
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
        {phases.map((phase) => (
          <div key={phase.id} className="bg-white p-4 rounded-lg shadow">
            <h2 className="text-xl font-bold">{phase.title}</h2>
            <p>Status: {phase.status}</p>
            <p>Planner: {phase.plannerName}</p>
            <p>Applications: {phase.applications.length}</p>
          </div>
        ))}
      </div>
    </div>
  );
};

export default PhasePlanningPage;
