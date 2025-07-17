import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { phaseApi, usersApi, expertsApi } from '../services/api';
import { User, Expert, PhaseApplication } from '../types';
import { Button, Input, Select } from '../components/ui';

const CreatePhasePage = () => {
  const [title, setTitle] = useState('');
  const [assignedPlannerId, setAssignedPlannerId] = useState<number | null>(null);
  const [applications, setApplications] = useState<Partial<PhaseApplication>[]>([
    { type: 'QP', institutionName: '', qualificationName: '', expert1: undefined, expert2: undefined },
  ]);
  const [planners, setPlanners] = useState<User[]>([]);
  const [experts, setExperts] = useState<Expert[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const navigate = useNavigate();

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [plannersRes, expertsRes] = await Promise.all([
          usersApi.getUsers(1, 100, { role: 'planner' }),
          expertsApi.getExperts(1, 1000), // Fetch a large number of experts
        ]);

        if (plannersRes.success) {
          setPlanners(plannersRes.data.data);
        } else {
          setError(plannersRes.message || 'Failed to load planners');
        }

        if (expertsRes.success) {
          setExperts(expertsRes.data.experts);
        } else {
          setError(expertsRes.message || 'Failed to load experts');
        }
      } catch (error) {
        setError('An error occurred while fetching data');
      } finally {
        setIsLoading(false);
      }
    };

    fetchData();
  }, []);

  const handleApplicationChange = (index: number, field: keyof PhaseApplication, value: any) => {
    const newApplications = [...applications];
    newApplications[index] = { ...newApplications[index], [field]: value };
    setApplications(newApplications);
  };

  const addApplication = () => {
    setApplications([...applications, { type: 'QP', institutionName: '', qualificationName: '', expert1: undefined, expert2: undefined }]);
  };

  const removeApplication = (index: number) => {
    const newApplications = applications.filter((_, i) => i !== index);
    setApplications(newApplications);
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);

    if (!title || !assignedPlannerId || applications.length === 0) {
      setError('Please fill in all required fields.');
      return;
    }

    try {
      const response = await phaseApi.createPhase({
        title,
        assignedSchedulerId: assignedPlannerId,
        status: 'planning',
        applications: applications.map(app => ({ ...app, status: 'pending' })) as any,
      });

      if (response.success) {
        navigate('/phases');
      } else {
        setError(response.message || 'Failed to create phase');
      }
    } catch (error) {
      setError('An error occurred while creating the phase');
    }
  };

  if (isLoading) {
    return <div>Loading...</div>;
  }

  return (
    <div className="max-w-4xl mx-auto px-4 py-8">
      <h1 className="text-2xl font-bold mb-6">Create New Phase</h1>
      <form onSubmit={handleSubmit} className="space-y-6">
        <div className="bg-white p-6 rounded-md shadow-sm">
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label htmlFor="title" className="block text-sm font-medium text-gray-700">Phase Title</label>
              <Input
                id="title"
                value={title}
                onChange={(e) => setTitle(e.target.value)}
                required
              />
            </div>
            <div>
              <label htmlFor="planner" className="block text-sm font-medium text-gray-700">Assigned Planner</label>
              <Select
                id="planner"
                value={assignedPlannerId || ''}
                onChange={(e) => setAssignedPlannerId(Number(e.target.value))}
                required
              >
                <option value="" disabled>Select a planner</option>
                {planners.map(planner => (
                  <option key={planner.id} value={planner.id}>{planner.name}</option>
                ))}
              </Select>
            </div>
          </div>
        </div>

        <div className="bg-white p-6 rounded-md shadow-sm">
          <h2 className="text-xl font-semibold mb-4">Applications</h2>
          {applications.map((app, index) => (
            <div key={index} className="border-b pb-6 mb-6">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                <Input
                  placeholder="Institution Name"
                  value={app.institutionName || ''}
                  onChange={(e) => handleApplicationChange(index, 'institutionName', e.target.value)}
                  required
                />
                <Input
                  placeholder="Qualification Name"
                  value={app.qualificationName || ''}
                  onChange={(e) => handleApplicationChange(index, 'qualificationName', e.target.value)}
                  required
                />
                <Select
                  value={app.type || 'QP'}
                  onChange={(e) => handleApplicationChange(index, 'type', e.target.value)}
                >
                  <option value="QP">Qualification Placement</option>
                  <option value="IL">Institutional Listing</option>
                </Select>
                <Select
                  value={app.expert1 || ''}
                  onChange={(e) => handleApplicationChange(index, 'expert1', Number(e.target.value))}
                  required
                >
                  <option value="" disabled>Select Expert 1</option>
                  {experts.map(expert => (
                    <option key={expert.id} value={expert.id}>{expert.name}</option>
                  ))}
                </Select>
                <Select
                  value={app.expert2 || ''}
                  onChange={(e) => handleApplicationChange(index, 'expert2', Number(e.target.value))}
                  required
                >
                  <option value="" disabled>Select Expert 2</option>
                  {experts.map(expert => (
                    <option key={expert.id} value={expert.id}>{expert.name}</option>
                  ))}
                </Select>
              </div>
              <div className="mt-4 flex justify-end">
                <Button type="button" variant="danger" onClick={() => removeApplication(index)}>
                  Remove
                </Button>
              </div>
            </div>
          ))}
          <Button type="button" variant="secondary" onClick={addApplication}>
            Add Application
          </Button>
        </div>

        {error && <div className="text-red-500">{error}</div>}

        <div className="flex justify-end space-x-4">
          <Button type="button" variant="outline" onClick={() => navigate('/phases')}>Cancel</Button>
          <Button type="submit" variant="primary">Create Phase</Button>
        </div>
      </form>
    </div>
  );
};

export default CreatePhasePage;