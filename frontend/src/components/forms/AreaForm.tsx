import React, { useState, useEffect } from 'react';
import { ExpertArea } from '../../types';

interface AreaFormProps {
  area: ExpertArea | null;
  onSubmit: (name: string) => void;
  onCancel: () => void;
}

const AreaForm: React.FC<AreaFormProps> = ({ area, onSubmit, onCancel }) => {
  const [name, setName] = useState('');

  useEffect(() => {
    if (area) {
      setName(area.name);
    } else {
      setName('');
    }
  }, [area]);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    onSubmit(name);
  };

  return (
    <form onSubmit={handleSubmit}>
      <div className="mt-4">
        <label htmlFor="name" className="block text-sm font-medium text-gray-700">Area Name</label>
        <input
          type="text"
          id="name"
          value={name}
          onChange={(e) => setName(e.target.value)}
          className="mt-1 block w-full px-3 py-2 bg-white border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
        />
      </div>
      <div className="mt-6 flex justify-end">
        <button
          type="button"
          onClick={onCancel}
          className="bg-gray-200 hover:bg-gray-300 text-gray-800 font-bold py-2 px-4 rounded mr-2"
        >
          Cancel
        </button>
        <button
          type="submit"
          className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded"
        >
          Save
        </button>
      </div>
    </form>
  );
};

export default AreaForm;