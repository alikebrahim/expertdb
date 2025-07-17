import React, { useState, useEffect } from 'react';
import { User, Phase, PhaseApplication } from '../../types';
import { getUserAssignments, assignUserToPlannerApplications, assignUserToManagerApplications } from '../../api/roleAssignments';
import { getPhases } from '../../api/phases';

interface RoleAssignmentFormProps {
  user: User;
  onSuccess: () => void;
  onCancel: () => void;
}

interface PhaseWithApplications extends Phase {
  applications: PhaseApplication[];
}

export const RoleAssignmentForm: React.FC<RoleAssignmentFormProps> = ({
  user,
  onSuccess,
  onCancel,
}) => {
  const [phases, setPhases] = useState<PhaseWithApplications[]>([]);
  const [currentAssignments, setCurrentAssignments] = useState<{
    planner_applications: number[];
    manager_applications: number[];
  }>({ planner_applications: [], manager_applications: [] });
  const [selectedPlannerApps, setSelectedPlannerApps] = useState<Set<number>>(new Set());
  const [selectedManagerApps, setSelectedManagerApps] = useState<Set<number>>(new Set());
  const [loading, setLoading] = useState(true);
  const [saving, setSaving] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    const loadData = async () => {
      try {
        setLoading(true);
        
        // Load phases and user assignments in parallel
        const [phasesResponse, assignmentsResponse] = await Promise.all([
          getPhases(),
          getUserAssignments(user.id)
        ]);

        if (phasesResponse.success && assignmentsResponse.success) {
          setPhases(phasesResponse.data.phases);
          setCurrentAssignments(assignmentsResponse.data);
          
          // Initialize selected apps with current assignments
          setSelectedPlannerApps(new Set(assignmentsResponse.data.planner_applications));
          setSelectedManagerApps(new Set(assignmentsResponse.data.manager_applications));
        } else {
          setError('Failed to load data');
        }
      } catch (err) {
        setError('Error loading data');
        console.error('Error loading role assignment data:', err);
      } finally {
        setLoading(false);
      }
    };

    loadData();
  }, [user.id]);

  const handlePlannerAppToggle = (appId: number) => {
    const newSet = new Set(selectedPlannerApps);
    if (newSet.has(appId)) {
      newSet.delete(appId);
    } else {
      newSet.add(appId);
    }
    setSelectedPlannerApps(newSet);
  };

  const handleManagerAppToggle = (appId: number) => {
    const newSet = new Set(selectedManagerApps);
    if (newSet.has(appId)) {
      newSet.delete(appId);
    } else {
      newSet.add(appId);
    }
    setSelectedManagerApps(newSet);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    
    try {
      setSaving(true);
      setError(null);

      // Update planner assignments
      const plannerAppsArray = Array.from(selectedPlannerApps);
      const plannerResponse = await assignUserToPlannerApplications(user.id, plannerAppsArray);
      
      if (!plannerResponse.success) {
        throw new Error('Failed to update planner assignments');
      }

      // Update manager assignments
      const managerAppsArray = Array.from(selectedManagerApps);
      const managerResponse = await assignUserToManagerApplications(user.id, managerAppsArray);
      
      if (!managerResponse.success) {
        throw new Error('Failed to update manager assignments');
      }

      onSuccess();
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to update assignments');
      console.error('Error updating role assignments:', err);
    } finally {
      setSaving(false);
    }
  };

  if (loading) {
    return (
      <div className="p-6">
        <div className="flex items-center justify-center">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          <span className="ml-2">Loading assignments...</span>
        </div>
      </div>
    );
  }

  return (
    <div className="bg-white rounded-lg shadow-lg p-6">
      <div className="flex items-center justify-between mb-6">
        <h2 className="text-xl font-semibold text-gray-900">
          Manage Role Assignments for {user.name}
        </h2>
        <button
          onClick={onCancel}
          className="text-gray-400 hover:text-gray-600"
        >
          <span className="sr-only">Close</span>
          ✕
        </button>
      </div>

      {error && (
        <div className="mb-4 p-3 bg-red-100 border border-red-400 text-red-700 rounded">
          {error}
        </div>
      )}

      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="text-sm text-gray-600 mb-4">
          <p><strong>Current Email:</strong> {user.email}</p>
          <p><strong>Current Role:</strong> {user.role}</p>
        </div>

        {phases.map((phase) => (
          <div key={phase.id} className="border border-gray-200 rounded-lg p-4">
            <h3 className="text-lg font-medium text-gray-900 mb-3">
              {phase.title} ({phase.phaseId})
            </h3>
            <p className="text-sm text-gray-600 mb-4">Status: {phase.status}</p>

            {phase.applications && phase.applications.length > 0 ? (
              <div className="space-y-3">
                {phase.applications.map((app) => (
                  <div key={app.id} className="bg-gray-50 p-3 rounded border">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <h4 className="font-medium text-gray-900">
                          {app.qualificationName}
                        </h4>
                        <p className="text-sm text-gray-600">
                          Institution: {app.institutionName}
                        </p>
                        <p className="text-sm text-gray-600">
                          Type: {app.type} | Status: {app.status}
                        </p>
                      </div>
                      <div className="flex flex-col space-y-2 ml-4">
                        <label className="flex items-center">
                          <input
                            type="checkbox"
                            checked={selectedPlannerApps.has(app.id)}
                            onChange={() => handlePlannerAppToggle(app.id)}
                            className="h-4 w-4 text-blue-600 border-gray-300 rounded focus:ring-blue-500"
                          />
                          <span className="ml-2 text-sm text-gray-700">Planner</span>
                        </label>
                        <label className="flex items-center">
                          <input
                            type="checkbox"
                            checked={selectedManagerApps.has(app.id)}
                            onChange={() => handleManagerAppToggle(app.id)}
                            className="h-4 w-4 text-green-600 border-gray-300 rounded focus:ring-green-500"
                          />
                          <span className="ml-2 text-sm text-gray-700">Manager</span>
                        </label>
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            ) : (
              <p className="text-sm text-gray-500 italic">No applications in this phase</p>
            )}
          </div>
        ))}

        <div className="flex justify-end space-x-3 pt-4 border-t">
          <button
            type="button"
            onClick={onCancel}
            className="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
          >
            Cancel
          </button>
          <button
            type="submit"
            disabled={saving}
            className="px-4 py-2 text-sm font-medium text-white bg-blue-600 border border-transparent rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {saving ? 'Saving...' : 'Update Assignments'}
          </button>
        </div>
      </form>
    </div>
  );
};