import { ExpertRequest } from '../types';
import { Table, Button } from './ui';

interface ExpertRequestTableProps {
  requests: ExpertRequest[];
  isLoading: boolean;
  error: string | null;
  onResubmit?: (request: ExpertRequest) => void;
  pagination?: {
    currentPage: number;
    totalPages: number;
    onPageChange: (page: number) => void;
  };
}

const ExpertRequestTable = ({ 
  requests, 
  isLoading, 
  error, 
  onResubmit
}: ExpertRequestTableProps) => {

  
  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };
  
  const getStatusBadgeClass = (status: string) => {
    switch (status) {
      case 'approved':
        return 'bg-green-100 text-green-800';
      case 'rejected':
        return 'bg-red-100 text-red-800';
      default:
        return 'bg-yellow-100 text-yellow-800';
    }
  };
  
  if (error) {
    return (
      <div className="bg-secondary bg-opacity-10 text-secondary p-4 rounded">
        <p>Error loading requests: {error}</p>
      </div>
    );
  }
  
  return (
    <div className="space-y-4">
      {isLoading ? (
        <div className="flex justify-center py-8">
          <div className="animate-spin h-8 w-8 border-4 border-primary border-t-transparent rounded-full"></div>
          <span className="sr-only">Loading...</span>
        </div>
      ) : requests.length === 0 ? (
        <div className="bg-accent p-6 rounded text-center">
          <p className="text-neutral-600">No expert requests found.</p>
          <p className="text-sm text-neutral-500 mt-1">
            You have not submitted any expert requests yet.
          </p>
        </div>
      ) : (
        <Table>
          <Table.Header>
            <Table.Row>
              <Table.HeaderCell>Name</Table.HeaderCell>
              <Table.HeaderCell>Institution</Table.HeaderCell>
              <Table.HeaderCell>Role</Table.HeaderCell>
              <Table.HeaderCell>Status</Table.HeaderCell>
              <Table.HeaderCell>Request Date</Table.HeaderCell>
              <Table.HeaderCell>Actions</Table.HeaderCell>
            </Table.Row>
          </Table.Header>
          <Table.Body>
            {requests.map((request) => (
              <Table.Row key={request.id}>
                <Table.Cell>{request.name}</Table.Cell>
                <Table.Cell>{request.institution}</Table.Cell>
                <Table.Cell>{request.role}</Table.Cell>
                <Table.Cell>
                  <span className={`inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium ${getStatusBadgeClass(request.status)}`}>
                    {request.status.charAt(0).toUpperCase() + request.status.slice(1)}
                  </span>
                </Table.Cell>
                <Table.Cell>{formatDate(request.createdAt)}</Table.Cell>
                <Table.Cell>
                  {request.status === 'rejected' && onResubmit && (
                    <Button 
                      variant="outline"
                      size="sm"
                      onClick={() => onResubmit(request)}
                    >
                      Resubmit
                    </Button>
                  )}
                  
                  {request.status === 'approved' && (
                    <a 
                      href={`/api/experts/${request.id}/approval-pdf`}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="inline-block"
                    >
                      <Button 
                        variant="outline"
                        size="sm"
                      >
                        Download PDF
                      </Button>
                    </a>
                  )}
                </Table.Cell>
              </Table.Row>
            ))}
          </Table.Body>
        </Table>
      )}
    </div>
  );
};

export default ExpertRequestTable;