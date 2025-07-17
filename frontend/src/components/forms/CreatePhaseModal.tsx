import React, { useState, useEffect } from 'react';
import { User } from '../../types';
import { getUsers } from '../../api/users';
import { createPhase } from '../../api/phases';

interface CreatePhaseModalProps {
  isOpen: boolean;
  onClose: () => void;
  onPhaseCreated: () => void;
}

const CreatePhaseModal: React.FC<CreatePhaseModalProps> = ({ isOpen, onClose, onPhaseCreated }) => {
  const [title, setTitle] = useState('');
  const [plannerId, setPlannerId] = useState<number | null>(null);
  const [planners, setPlanners] = useState<User[]>([]);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (isOpen) {
      const fetchPlanners = async () => {
        try {
          // Assuming getUsers can filter by role, or we filter on the client
          const response = await getUsers();
          if (response.success && response.data) {
            const plannerUsers = response.data.data.filter(user => user.role === 'planner');
            setPlanners(plannerUsers);
          }
        } catch (_err) {
          setError('Failed to fetch planners');
        }
      };
      fetchPlanners();
    }
  }, [isOpen]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!title || !plannerId) {
      setError('Title and planner are required');
      return;
    }

    try {
      await createPhase({ title, assignedPlannerId: plannerId, status: 'active', applications: [] });
      onPhaseCreated();
      onClose();
    } catch (_err) {
      setError('Failed to create phase');
    }
  };

  if (!isOpen) return null;

  return (
    <div className="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full">
      <div className="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
        <form onSubmit={handleSubmit}>
          <h3 className="text-lg leading-6 font-medium text-gray-900">Create New Phase</h3>
          {error && <p className="text-red-500">{error}</p>}
          <div className="mt-4">
            <label htmlFor="title" className="block text-sm font-medium text-gray-700">Title</label>
            <input
              type="text"
              id="title"
              value={title}
              onChange={(e) => setTitle(e.target.value)}
              className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
            />
          </div>
          <div className="mt-4">
            <label htmlFor="planner" className="block text-sm font-medium text-gray-700">Assign Planner</label>
            <select
              id="planner"
              value={plannerId ?? ''}
              onChange={(e) => setPlannerId(Number(e.target.value))}
              className="mt-1 block w-full pl-3 pr-10 py-2 text-base border-gray-300 focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
            >
              <option value="" disabled>Select a planner</option>
              {planners.map(planner => (
                <option key={planner.id} value={planner.id}>{planner.name}</option>
              ))}
            </select>
          </div>
          <div className="mt-6 flex justify-end">
            <button
              type="button"
              onClick={onClose}
              className="bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-2 px-4 rounded mr-2"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
            >
              Create
            </button>
          </div>
        </form>
      </div>
    </div>
  );
};

export default CreatePhaseModal;