import React, { useState, useEffect } from 'react';
import { Control, useFieldArray, UseFormSetValue, UseFormWatch } from 'react-hook-form';
import { Plus, Trash2 } from 'lucide-react';
import { FormField } from './ui/FormField';
import Button from './ui/Button';
import { Card, CardHeader, CardContent } from './ui/Card';

interface ExperienceEntry {
  dateFrom: string;
  dateTo: string;
  description: string;
}

interface EducationEntry {
  dateFrom: string;
  dateTo: string;
  description: string;
}

interface Biography {
  experience: ExperienceEntry[];
  education: EducationEntry[];
}

interface BiographyFormProps {
  control: Control<any>;
  setValue: UseFormSetValue<any>;
  watch: UseFormWatch<any>;
  errors?: any;
}

const BiographyForm: React.FC<BiographyFormProps> = ({
  control,
  setValue,
  watch,
  errors,
}) => {
  const {
    fields: experienceFields,
    append: appendExperience,
    remove: removeExperience,
  } = useFieldArray({
    control,
    name: 'biography.experience',
  });

  const {
    fields: educationFields,
    append: appendEducation,
    remove: removeEducation,
  } = useFieldArray({
    control,
    name: 'biography.education',
  });

  const addExperienceEntry = () => {
    appendExperience({
      dateFrom: '',
      dateTo: '',
      description: '',
    });
  };

  const addEducationEntry = () => {
    appendEducation({
      dateFrom: '',
      dateTo: '',
      description: '',
    });
  };

  // Initialize with one entry of each type if none exist
  useEffect(() => {
    if (experienceFields.length === 0) {
      addExperienceEntry();
    }
    if (educationFields.length === 0) {
      addEducationEntry();
    }
  }, []);

  return (
    <div className="space-y-6">
      {/* Experience Section */}
      <Card>
        <CardHeader>
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold">Professional Experience</h3>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={addExperienceEntry}
              className="flex items-center gap-2"
            >
              <Plus className="h-4 w-4" />
              Add Experience
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {experienceFields.map((field, index) => (
              <div key={field.id} className="border rounded-lg p-4 space-y-4">
                <div className="flex justify-between items-center">
                  <h4 className="font-medium">Experience {index + 1}</h4>
                  {experienceFields.length > 1 && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => removeExperience(index)}
                      className="text-red-600 hover:text-red-800"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <FormField
                    name={`biography.experience.${index}.dateFrom`}
                    label="Start Date (YYYY-MM)"
                    control={control}
                    error={errors?.biography?.experience?.[index]?.dateFrom}
                    placeholder="e.g., 2020-01"
                    required
                  />
                  <FormField
                    name={`biography.experience.${index}.dateTo`}
                    label="End Date (YYYY-MM)"
                    control={control}
                    error={errors?.biography?.experience?.[index]?.dateTo}
                    placeholder="e.g., 2023-12 or Present"
                    required
                  />
                </div>
                
                <FormField
                  name={`biography.experience.${index}.description`}
                  label="Experience Description"
                  type="textarea"
                  control={control}
                  error={errors?.biography?.experience?.[index]?.description}
                  placeholder="Format: Role/Position, Organization, Location/country[optional]"
                  required
                  rows={2}
                />
                
                <div className="text-sm text-gray-600 bg-gray-50 p-3 rounded">
                  <p><strong>Format:</strong> Role/Position, Organization, Location/country[optional]</p>
                  <p><strong>Example:</strong> Senior Engineer, Ministry of Works, Bahrain</p>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Education Section */}
      <Card>
        <CardHeader>
          <div className="flex justify-between items-center">
            <h3 className="text-lg font-semibold">Education</h3>
            <Button
              type="button"
              variant="outline"
              size="sm"
              onClick={addEducationEntry}
              className="flex items-center gap-2"
            >
              <Plus className="h-4 w-4" />
              Add Education
            </Button>
          </div>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {educationFields.map((field, index) => (
              <div key={field.id} className="border rounded-lg p-4 space-y-4">
                <div className="flex justify-between items-center">
                  <h4 className="font-medium">Education {index + 1}</h4>
                  {educationFields.length > 1 && (
                    <Button
                      type="button"
                      variant="ghost"
                      size="sm"
                      onClick={() => removeEducation(index)}
                      className="text-red-600 hover:text-red-800"
                    >
                      <Trash2 className="h-4 w-4" />
                    </Button>
                  )}
                </div>
                
                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <FormField
                    name={`biography.education.${index}.dateFrom`}
                    label="Start Date (YYYY-MM)"
                    control={control}
                    error={errors?.biography?.education?.[index]?.dateFrom}
                    placeholder="e.g., 2015-09"
                    required
                  />
                  <FormField
                    name={`biography.education.${index}.dateTo`}
                    label="End Date (YYYY-MM)"
                    control={control}
                    error={errors?.biography?.education?.[index]?.dateTo}
                    placeholder="e.g., 2018-06"
                    required
                  />
                </div>
                
                <FormField
                  name={`biography.education.${index}.description`}
                  label="Education Description"
                  type="textarea"
                  control={control}
                  error={errors?.biography?.education?.[index]?.description}
                  placeholder="Format: Degree, Institution, Location/country[optional]"
                  required
                  rows={2}
                />
                
                <div className="text-sm text-gray-600 bg-gray-50 p-3 rounded">
                  <p><strong>Format:</strong> Degree, Institution, Location/country[optional]</p>
                  <p><strong>Example:</strong> PhD Civil Engineering, University of Bahrain</p>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Biography Preview */}
      <Card>
        <CardHeader>
          <h3 className="text-lg font-semibold">Biography Preview</h3>
        </CardHeader>
        <CardContent>
          <BiographyPreview watch={watch} />
        </CardContent>
      </Card>
    </div>
  );
};

// Biography Preview Component
interface BiographyPreviewProps {
  watch: UseFormWatch<any>;
}

const BiographyPreview: React.FC<BiographyPreviewProps> = ({ watch }) => {
  const biography = watch('biography');
  
  if (!biography || (!biography.experience?.length && !biography.education?.length)) {
    return <p className="text-gray-500 italic">Biography preview will appear here as you fill in the form.</p>;
  }

  const formatBiography = (biography: Biography): string => {
    let formatted = '';
    
    if (biography.education?.length > 0) {
      formatted += 'Education:\n';
      biography.education.forEach((entry: EducationEntry) => {
        if (entry.dateFrom && entry.dateTo && entry.description) {
          formatted += `[${entry.dateFrom} - ${entry.dateTo}] ${entry.description}\n`;
        }
      });
    }
    
    if (biography.experience?.length > 0) {
      formatted += '\nExperience:\n';
      biography.experience.forEach((entry: ExperienceEntry) => {
        if (entry.dateFrom && entry.dateTo && entry.description) {
          formatted += `[${entry.dateFrom} - ${entry.dateTo}] ${entry.description}\n`;
        }
      });
    }
    
    return formatted;
  };

  return (
    <div className="bg-gray-50 p-4 rounded border">
      <pre className="whitespace-pre-wrap text-sm font-mono">
        {formatBiography(biography) || 'Biography preview will appear here as you fill in the form.'}
      </pre>
    </div>
  );
};

export default BiographyForm;