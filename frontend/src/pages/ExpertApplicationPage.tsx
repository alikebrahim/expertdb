import React, { useState } from 'react';
import { useAuth } from '../hooks/useAuth';
import ExpertRequestSubmissionForm from '../components/ExpertRequestSubmissionForm';
import { Card, CardHeader, CardContent } from '../components/ui/Card';
import { Alert } from '../components/ui/Alert';
import Button from '../components/ui/Button';

const ExpertApplicationPage: React.FC = () => {
  const { user } = useAuth();
  const [showForm, setShowForm] = useState(false);
  const [successMessage, setSuccessMessage] = useState<string | null>(null);

  const handleSuccess = () => {
    setShowForm(false);
    setSuccessMessage('Your expert application has been submitted successfully! You will be contacted once your application is reviewed.');
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  const handleStartApplication = () => {
    setShowForm(true);
    setSuccessMessage(null);
    window.scrollTo({ top: 0, behavior: 'smooth' });
  };

  if (!user) {
    return (
      <div className="max-w-4xl mx-auto">
        <Alert variant="warning">
          <h3 className="font-medium">Authentication Required</h3>
          <p className="mt-1">You must be logged in to apply as an expert. Please log in and try again.</p>
        </Alert>
      </div>
    );
  }

  return (
    <div className="max-w-4xl mx-auto">
      <div className="mb-8">
        <h1 className="text-3xl font-bold text-primary mb-2">Apply to Become an Expert</h1>
        <p className="text-lg text-gray-600">
          Join the BQA expert network and contribute your expertise to quality assurance initiatives in Bahrain.
        </p>
      </div>

      {successMessage && (
        <div className="mb-8">
          <Alert variant="success">
            <h3 className="font-medium">Application Submitted</h3>
            <p className="mt-1">{successMessage}</p>
          </Alert>
        </div>
      )}

      {!showForm && !successMessage && (
        <div className="space-y-6">
          <Card>
            <CardHeader>
              <h2 className="text-xl font-semibold text-primary">Expert Application Process</h2>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <p className="text-gray-700">
                  As a BQA expert, you will contribute to maintaining and improving quality standards across 
                  various sectors in Bahrain. Our experts play a crucial role in assessments, reviews, 
                  and consultations.
                </p>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
                  <div>
                    <h3 className="font-semibold text-primary mb-2">Requirements</h3>
                    <ul className="space-y-1 text-sm text-gray-600">
                      <li>• Relevant professional experience</li>
                      <li>• Academic qualifications in your field</li>
                      <li>• Strong communication skills</li>
                      <li>• Commitment to quality standards</li>
                      <li>• Availability for BQA activities</li>
                    </ul>
                  </div>
                  
                  <div>
                    <h3 className="font-semibold text-primary mb-2">Application Steps</h3>
                    <ul className="space-y-1 text-sm text-gray-600">
                      <li>1. Complete personal information</li>
                      <li>2. Provide professional details</li>
                      <li>3. Define expertise areas</li>
                      <li>4. Submit biography and CV</li>
                    </ul>
                  </div>
                </div>
                
                <div className="bg-blue-50 p-4 rounded-lg">
                  <h4 className="font-medium text-blue-900 mb-2">What happens next?</h4>
                  <p className="text-sm text-blue-800">
                    Once you submit your application, our team will review your qualifications and experience. 
                    If your profile matches our requirements, we will contact you for the next steps in the 
                    onboarding process.
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>

          <div className="text-center">
            <Button onClick={handleStartApplication} size="lg">
              Start Your Application
            </Button>
          </div>
        </div>
      )}

      {showForm && (
        <div>
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-xl font-semibold text-primary">Expert Application Form</h2>
            <Button
              variant="outline"
              onClick={() => setShowForm(false)}
            >
              Cancel
            </Button>
          </div>
          
          <ExpertRequestSubmissionForm onSuccess={handleSuccess} />
        </div>
      )}
    </div>
  );
};

export default ExpertApplicationPage;
