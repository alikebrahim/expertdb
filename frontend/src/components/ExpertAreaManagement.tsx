import { useState, useEffect, useCallback } from 'react';
import { useForm } from 'react-hook-form';
import { useUI } from '../hooks/useUI';
import * as areasApi from '../api/areas';
import Button from './ui/Button';
import Input from './ui/Input';

interface ExpertArea {
  id: number;
  name: string;
}

interface ExpertAreaFormData {
  name: string;
}

const ExpertAreaManagement = () => {
  const [areas, setAreas] = useState<ExpertArea[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [editingAreaId, setEditingAreaId] = useState<number | null>(null);
  const { addNotification } = useUI();
  const { register, handleSubmit, reset, setValue, formState: { errors } } = useForm<ExpertAreaFormData>();

  const fetchAreas = useCallback(async () => {
    setIsLoading(true);
    try {
      const response = await areasApi.getExpertAreas();
      if (response.success && response.data) {
        setAreas(response.data);
      } else {
        addNotification({
          type: 'error',
          message: 'Failed to load expert areas',
          duration: 5000,
        });
      }
    } catch (error) {
      console.error('Error fetching expert areas:', error);
      addNotification({
        type: 'error',
        message: 'Error loading expert areas',
        duration: 5000,
      });
    } finally {
      setIsLoading(false);
    }
  }, [addNotification]);

  // Fetch areas on mount
  useEffect(() => {
    fetchAreas();
  }, [fetchAreas]);

  const handleCreateArea = async (data: ExpertAreaFormData) => {
    try {
      setIsLoading(true);
      const response = await areasApi.createExpertArea({ name: data.name });
      
      if (response.success) {
        addNotification({
          type: 'success',
          message: 'Expert area created successfully',
          duration: 3000,
        });
        fetchAreas();
        reset();
      } else {
        addNotification({
          type: 'error',
          message: response.message || 'Failed to create expert area',
          duration: 5000,
        });
      }
    } catch (error) {
      console.error('Error creating expert area:', error);
      addNotification({
        type: 'error',
        message: 'Error creating expert area',
        duration: 5000,
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleUpdateArea = async (data: ExpertAreaFormData) => {
    if (!editingAreaId) return;
    
    try {
      setIsLoading(true);
      const response = await areasApi.updateExpertArea(editingAreaId, { name: data.name });
      
      if (response.success) {
        addNotification({
          type: 'success',
          message: 'Expert area updated successfully',
          duration: 3000,
        });
        fetchAreas();
        reset();
        setEditingAreaId(null);
      } else {
        addNotification({
          type: 'error',
          message: response.message || 'Failed to update expert area',
          duration: 5000,
        });
      }
    } catch (error) {
      console.error('Error updating expert area:', error);
      addNotification({
        type: 'error',
        message: 'Error updating expert area',
        duration: 5000,
      });
    } finally {
      setIsLoading(false);
    }
  };

  const handleEditClick = (area: ExpertArea) => {
    setEditingAreaId(area.id);
    setValue('name', area.name);
  };

  const handleCancelEdit = () => {
    setEditingAreaId(null);
    reset();
  };

  const onSubmit = (data: ExpertAreaFormData) => {
    if (editingAreaId) {
      handleUpdateArea(data);
    } else {
      handleCreateArea(data);
    }
  };

  return (
    <div className="bg-white rounded-md shadow p-6">
      <h2 className="text-xl font-semibold text-primary mb-4">Expert Areas Management</h2>
      
      <div className="mb-6">
        <form onSubmit={handleSubmit(onSubmit)} className="flex items-end gap-2">
          <div className="flex-grow">
            <Input
              label={editingAreaId ? 'Edit Area Name' : 'New Area Name'}
              error={errors.name?.message}
              {...register('name', { 
                required: 'Area name is required',
                minLength: {
                  value: 2,
                  message: 'Area name must be at least 2 characters'
                } 
              })}
            />
          </div>
          <div className="flex gap-2">
            {editingAreaId ? (
              <>
                <Button
                  type="button"
                  variant="outline"
                  onClick={handleCancelEdit}
                  disabled={isLoading}
                >
                  Cancel
                </Button>
                <Button
                  type="submit"
                  variant="primary"
                  disabled={isLoading}
                >
                  Update
                </Button>
              </>
            ) : (
              <Button
                type="submit"
                variant="primary"
                disabled={isLoading}
              >
                Add Area
              </Button>
            )}
          </div>
        </form>
      </div>
      
      <div className="border rounded overflow-hidden">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Area Name
              </th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {isLoading && !areas.length ? (
              <tr>
                <td colSpan={2} className="px-6 py-4 text-center text-sm text-gray-500">
                  Loading areas...
                </td>
              </tr>
            ) : areas.length === 0 ? (
              <tr>
                <td colSpan={2} className="px-6 py-4 text-center text-sm text-gray-500">
                  No expert areas found. Add your first area above.
                </td>
              </tr>
            ) : (
              areas.map((area) => (
                <tr key={area.id} className="hover:bg-gray-50">
                  <td className="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                    {area.name}
                  </td>
                  <td className="px-6 py-4 whitespace-nowrap text-right text-sm">
                    <button
                      onClick={() => handleEditClick(area)}
                      className="text-indigo-600 hover:text-indigo-900 mr-4"
                      disabled={isLoading || editingAreaId === area.id}
                    >
                      Edit
                    </button>
                  </td>
                </tr>
              ))
            )}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default ExpertAreaManagement;