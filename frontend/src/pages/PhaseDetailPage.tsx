import { useState, useEffect, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { phaseApi, expertsApi } from '../services/api';
import { Phase, PhaseApplication, Expert } from '../types';
import { Button, Table, Select } from '../components/ui';
import Modal from '../components/Modal';
import { useAuth } from '../hooks/useAuth';

const PhaseDetailPage = () => {
  const { id } = useParams<{ id: string }>();
  const { user } = useAuth();
  const navigate = useNavigate();
  const [phase, setPhase] = useState<Phase | null>(null);
  const [experts, setExperts] = useState<Expert[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [selectedApplication, setSelectedApplication] = useState<PhaseApplication | null>(null);
  const [isReviewModalOpen, setIsReviewModalOpen] = useState(false);
  const [isProposeModalOpen, setIsProposeModalOpen] = useState(false);
  const [reviewStatus, setReviewStatus] = useState('approved');
  const [rejectionNotes, setRejectionNotes] = useState('');
  const [proposedExperts, setProposedExperts] = useState<{ expert1?: number; expert2?: number }>({});

  const fetchPhase = useCallback(async () => {
    if (!id) return;
    setIsLoading(true);
    setError(null);
    try {
      const response = await phaseApi.getPhaseById(id);
      if (response.success) {
        setPhase(response.data);
      } else {
        setError(response.message || 'Failed to load phase details');
      }
    } catch (error) {
      setError('An error occurred while fetching phase details');
    } finally {
      setIsLoading(false);
    }
  }, [id]);

  useEffect(() => {
    fetchPhase();
    expertsApi.getExperts(1, 1000).then(res => {
      if (res.success) setExperts(res.data.experts);
    });
  }, [fetchPhase]);

  const handleReviewSubmit = async () => {
    if (!selectedApplication || !id) return;
    try {
      const response = await phaseApi.reviewApplication(parseInt(id), selectedApplication.id, {
        status: reviewStatus,
        rejectionNotes: reviewStatus === 'rejected' ? rejectionNotes : undefined,
      });
      if (response.success) {
        setIsReviewModalOpen(false);
        fetchPhase();
      } else {
        setError(response.message || 'Failed to review application');
      }
    } catch (error) {
      setError('An error occurred while reviewing the application');
    }
  };

  const handleProposeSubmit = async () => {
    if (!selectedApplication || !id) return;
    try {
      const response = await phaseApi.proposeExperts(parseInt(id), selectedApplication.id, {
        expert1: proposedExperts.expert1 || selectedApplication.expert1,
        expert2: proposedExperts.expert2 || selectedApplication.expert2,
      });
      if (response.success) {
        setIsProposeModalOpen(false);
        fetchPhase();
      } else {
        setError(response.message || 'Failed to propose experts');
      }
    } catch (error) {
      setError('An error occurred while proposing experts');
    }
  };

  if (isLoading) return <div>Loading...</div>;
  if (error) return <div className="text-red-500">{error}</div>;
  if (!phase) return <div>Phase not found.</div>;

  return (
    <div className="max-w-7xl mx-auto px-4 py-8">
      <Button variant="outline" onClick={() => navigate('/phases')} className="mb-4">Back to Phases</Button>
      <div className="bg-white p-6 rounded-md shadow-sm mb-6">
        <h1 className="text-2xl font-bold">{phase.title}</h1>
        <p className="text-gray-600">Planner: {phase.plannerName}</p>
        <p className="text-gray-600">Status: {phase.status}</p>
      </div>

      <div className="bg-white p-6 rounded-md shadow-sm">
        <h2 className="text-xl font-semibold mb-4">Applications</h2>
        <Table>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Institution</Table.HeaderCell>
              <Table.HeaderCell>Qualification</Table.HeaderCell>
              <Table.HeaderCell>Type</Table.HeaderCell>
              <Table.HeaderCell>Expert 1</Table.HeaderCell>
              <Table.HeaderCell>Expert 2</Table.HeaderCell>
              <Table.HeaderCell>Status</Table.HeaderCell>
              <Table.HeaderCell>Actions</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {phase.applications.map(app => (
              <Table.Row key={app.id}>
                <Table.Cell>{app.institutionName}</Table.Cell>
                <Table.Cell>{app.qualificationName}</Table.Cell>
                <Table.Cell>{app.type}</Table.Cell>
                <Table.Cell>{experts.find(e => e.id === app.expert1)?.name || 'N/A'}</Table.Cell>
                <Table.Cell>{experts.find(e => e.id === app.expert2)?.name || 'N/A'}</Table.Cell>
                <Table.Cell>{app.status}</Table.Cell>
                <Table.Cell>
                  {user?.role === 'admin' && (
                    <Button size="sm" onClick={() => { setSelectedApplication(app); setIsReviewModalOpen(true); }}>Review</Button>
                  )}
                  {user?.role === 'planner' && (
                    <Button size="sm" variant="secondary" onClick={() => { setSelectedApplication(app); setProposedExperts({ expert1: app.expert1, expert2: app.expert2 }); setIsProposeModalOpen(true); }}>Propose</Button>
                  )}
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      </div>

      {selectedApplication && (
        <>
          <Modal isOpen={isReviewModalOpen} onClose={() => setIsReviewModalOpen(false)} title="Review Application">
            <div className="space-y-4">
              <Select value={reviewStatus} onChange={e => setReviewStatus(e.target.value)}>
                <option value="approved">Approve</option>
                <option value="rejected">Reject</option>
              </Select>
              {reviewStatus === 'rejected' && (
                <textarea value={rejectionNotes} onChange={e => setRejectionNotes(e.target.value)} placeholder="Rejection Notes" className="w-full p-2 border rounded" />
              )}
              <div className="flex justify-end space-x-2">
                <Button variant="outline" onClick={() => setIsReviewModalOpen(false)}>Cancel</Button>
                <Button onClick={handleReviewSubmit}>Submit</Button>
              </div>
            </div>
          </Modal>

          <Modal isOpen={isProposeModalOpen} onClose={() => setIsProposeModalOpen(false)} title="Propose Experts">
            <div className="space-y-4">
              <Select value={proposedExperts.expert1 || ''} onChange={e => setProposedExperts(p => ({ ...p, expert1: Number(e.target.value) }))}>
                <option value="">Select Expert 1</option>
                {experts.map(e => <option key={e.id} value={e.id}>{e.name}</option>)}
              </Select>
              <Select value={proposedExperts.expert2 || ''} onChange={e => setProposedExperts(p => ({ ...p, expert2: Number(e.target.value) }))}>
                <option value="">Select Expert 2</option>
                {experts.map(e => <option key={e.id} value={e.id}>{e.name}</option>)}
              </Select>
              <div className="flex justify-end space-x-2">
                <Button variant="outline" onClick={() => setIsProposeModalOpen(false)}>Cancel</Button>
                <Button onClick={handleProposeSubmit}>Submit</Button>
              </div>
            </div>
          </Modal>
        </>
      )}
    </div>
  );
};

export default PhaseDetailPage;