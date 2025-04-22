import { useState, useEffect, useCallback } from 'react';
import { documentApi } from '../services/api';
import { Document } from '../types';
import Button from './ui/Button';

interface DocumentListProps {
  expertId: number;
  onRefreshNeeded?: () => void;
}

const DocumentList = ({ expertId, onRefreshNeeded }: DocumentListProps) => {
  const [documents, setDocuments] = useState<Document[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isDeleting, setIsDeleting] = useState<number | null>(null);

  const fetchDocuments = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    
    try {
      const response = await documentApi.getExpertDocuments(expertId);
      
      if (response.success) {
        setDocuments(response.data.documents);
      } else {
        setError(response.message || 'Failed to load documents');
      }
    } catch (error) {
      console.error('Error fetching documents:', error);
      setError('An error occurred while loading documents');
    } finally {
      setIsLoading(false);
    }
  }, [expertId]);

  useEffect(() => {
    fetchDocuments();
  }, [expertId, fetchDocuments]);

  const handleDelete = async (id: number) => {
    if (!confirm('Are you sure you want to delete this document?')) {
      return;
    }
    
    setIsDeleting(id);
    
    try {
      const response = await documentApi.deleteDocument(id);
      
      if (response.success) {
        setDocuments(documents.filter(doc => doc.id !== id));
        if (onRefreshNeeded) {
          onRefreshNeeded();
        }
      } else {
        setError(response.message || 'Failed to delete document');
      }
    } catch (error) {
      console.error('Error deleting document:', error);
      setError('An error occurred while deleting the document');
    } finally {
      setIsDeleting(null);
    }
  };

  const formatFileSize = (bytes: number) => {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(2) + ' MB';
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  };

  const getDocumentTypeLabel = (type: string) => {
    const types: Record<string, string> = {
      'cv': 'Curriculum Vitae',
      'certificate': 'Certification',
      'research': 'Research Paper',
      'publication': 'Publication',
      'other': 'Other Document'
    };
    
    return types[type] || type;
  };

  const getDocumentIcon = (contentType: string) => {
    if (contentType.includes('pdf')) {
      return '📄';
    } else if (contentType.includes('word') || contentType.includes('doc')) {
      return '📝';
    } else if (contentType.includes('excel') || contentType.includes('sheet')) {
      return '📊';
    } else if (contentType.includes('powerpoint') || contentType.includes('presentation')) {
      return '📑';
    } else if (contentType.includes('image')) {
      return '🖼️';
    }
    return '📁';
  };

  if (isLoading) {
    return <div className="p-4 text-center">Loading documents...</div>;
  }

  if (error) {
    return (
      <div className="p-4 text-center text-red-600">
        {error}
        <button
          onClick={fetchDocuments}
          className="ml-2 text-primary underline"
        >
          Try again
        </button>
      </div>
    );
  }

  if (documents.length === 0) {
    return <div className="p-4 text-center text-gray-500">No documents available</div>;
  }

  return (
    <div className="overflow-hidden">
      <div className="bg-white overflow-x-auto">
        <table className="min-w-full divide-y divide-gray-200">
          <thead className="bg-gray-50">
            <tr>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Document
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Type
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Size
              </th>
              <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Uploaded
              </th>
              <th className="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody className="bg-white divide-y divide-gray-200">
            {documents.map((doc) => (
              <tr key={doc.id}>
                <td className="px-6 py-4 whitespace-nowrap">
                  <div className="flex items-center">
                    <span className="text-xl mr-2">📄</span>
                    <div className="truncate max-w-xs">{doc.filePath.split('/').pop()}</div>
                  </div>
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {getDocumentTypeLabel(doc.documentType)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {/* Size not available in updated API */}
                  N/A
                </td>
                <td className="px-6 py-4 whitespace-nowrap">
                  {formatDate(doc.createdAt)}
                </td>
                <td className="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                  <div className="flex justify-end space-x-2">
                    <Button
                      size="sm"
                      variant="outline"
                      onClick={() => window.open(`/api/documents/${doc.id}`, '_blank')}
                    >
                      View
                    </Button>
                    <Button
                      size="sm"
                      variant="danger"
                      isLoading={isDeleting === doc.id}
                      onClick={() => handleDelete(doc.id)}
                    >
                      Delete
                    </Button>
                  </div>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
};

export default DocumentList;